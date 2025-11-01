package ssa_callgraph

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Silhouette-sophist/repo_profile/zap_log"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Program struct {
	Args        InitProgramArgs
	Graph       *Graph
	RootPkgPath string
	PackagePkgs []*packages.Package
	SsaPkgs     []*ssa.Package
}

type Algo string

type InitProgramArgs struct {
	Path      string `validate:"required"`
	Algorithm string `validate:"one of=cha rta vta pta"`
}

func (p *Program) Load(ctx context.Context, args InitProgramArgs) error {
	p.Args = args
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
		Dir:   args.Path,
		Logf:  nil,
	}
	// 1.加载所有包的类型
	pakcagesPkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		zap_log.CtxError(ctx, "failed to load pakcagesPkgs packages %v", err)
		return err
	}
	p.PackagePkgs = pakcagesPkgs
	// 2.检测错误
	if err = CheckErrors(ctx, pakcagesPkgs); err != nil {
		zap_log.CtxError(ctx, "failed to load pakcagesPkgs packages %v", err)
		return err
	}
	// 3.创建ssa program，并获取当前仓库主main包
	prog, ssaPkgs := ssautil.AllPackages(pakcagesPkgs, 0)
	p.RootPkgPath = GetCommonPkgPath(ssaPkgs)
	p.SsaPkgs = ssaPkgs
	// 4.Build calls Package.Build for each package in prog.
	prog.Build()
	// 5.获取所有main包
	mainPkgs := ssautil.MainPackages(ssaPkgs)
	// 6.获取静态调用图
	var g *callgraph.Graph
	switch args.Algorithm {
	case "cha":
		g = cha.CallGraph(prog)
	case "rta":
		var roots []*ssa.Function
		for _, mainPkg := range mainPkgs {
			roots = append(roots, mainPkg.Func("main"))
			roots = append(roots, mainPkg.Func("init"))
		}
		res := rta.Analyze(roots, true)
		g = res.CallGraph
	case "vta":
		g = vta.CallGraph(ssautil.AllFunctions(prog), cha.CallGraph(prog))
	case "pta":
		result, err := pointer.Analyze(&pointer.Config{
			Mains:          mainPkgs,
			BuildCallGraph: true,
		})
		if err != nil {
			return err
		}
		g = result.CallGraph
	}
	// 7.转为本地图
	p.ToGraph(g)
	// 8.移除无用节点
	p.removeMeaningLessNode()
	return nil
}

func (p *Program) ToGraph(g *callgraph.Graph) {
	g.DeleteSyntheticNodes()
	graph := &Graph{NodeMap: make(map[string]*Node)}
	for fc, n := range g.Nodes {
		node := &Node{
			In:  make(map[string]*Edge),
			Out: make(map[string]*Edge),
		}
		node.Func = fc
		node.ID = strconv.Itoa(n.ID)
		graph.NodeMap[node.ID] = node
		packagePath := getPackagePath(fc)
		name := getFullFunctionName(fc)
		functionName := getFunctionName(fc)
		var targetPackage *packages.Package
		for _, pkg := range p.PackagePkgs {
			if pkg.PkgPath == packagePath {
				targetPackage = pkg
				break
			}
		}
		functionFile := getFunctionFile(fc, targetPackage)
		file := getFunctionFile(fc, p.PackagePkgs[0])
		fmt.Println(packagePath, functionName, name, functionFile, file)
	}

	for _, n := range g.Nodes {
		for _, e := range n.Out {
			edge := &Edge{}
			edge.Site = e.Site
			edge.CallerID = strconv.Itoa(e.Caller.ID)
			edge.CalleeID = strconv.Itoa(e.Callee.ID)
			graph.NodeMap[edge.CallerID].Out[edge.CalleeID] = edge
			graph.NodeMap[edge.CalleeID].In[edge.CallerID] = edge
		}
	}
	p.Graph = graph
}

func GetCommonPkgPath(pkgs []*ssa.Package) string {
	var paths []string
	for _, pkg := range pkgs {
		paths = append(paths, pkg.Pkg.Path())
	}
	return GetCommonPrefix(paths)
}

func (p *Program) removeMeaningLessNode() {
	for _, node := range p.Graph.NodeMap {
		if node.Func.Pkg == nil {
			p.Graph.Delete(node.ID)
			continue
		}
		pkgPath := node.Func.Pkg.Pkg.Path()
		if p.IsTargetPkg(pkgPath) {
			continue
		}
		// 删除非目标包的函数
		p.Graph.Delete(node.ID)
	}
}

func (p *Program) IsTargetPkg(pkgPath string) bool {
	return strings.HasPrefix(pkgPath, p.RootPkgPath) || strings.HasPrefix(pkgPath, p.RootPkgPath+"/")
}

func CheckErrors(ctx context.Context, pkgs []*packages.Package) error {
	eMsg := ""
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		if len(pkg.Errors) > 0 {
			eMsg += "\npackage " + pkg.PkgPath + " contain errors"
			for _, err := range pkg.Errors {
				eMsg += "\n\t" + err.Error()
			}
		}
	})
	if eMsg != "" {
		return errors.New(eMsg)
	} else {
		return nil
	}
}

func GetCommonPrefix(paths []string) string {
	res := ""
	for i := 0; true; i++ {
		var preCh uint8
		for _, p := range paths {
			if i >= len(p) {
				return res
			}
			ch := p[i]
			if preCh == 0 {
				preCh = ch
			}
			if ch != preCh {
				return res
			}
		}
		res = res + string(preCh)
	}
	return res
}
