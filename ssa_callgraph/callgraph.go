package ssa_callgraph

import (
	"errors"
	"strconv"
	"strings"

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
}

type Algo string

type InitProgramArgs struct {
	Path      string `validate:"required"`
	Algorithm string `validate:"oneof=cha rta vta pta"`
}

func (p *Program) Load(args InitProgramArgs) error {
	p.Args = args
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
		Dir:   args.Path,
		Logf:  nil,
	}
	initial, err := packages.Load(cfg, "./...")
	if err != nil {
		return err
	}
	err = CheckErrors(initial)
	if err != nil {
		return err
	}
	prog, pkgs := ssautil.AllPackages(initial, 0)
	p.RootPkgPath = GetCommonPkgPath(pkgs)
	prog.Build()
	mainPkgs := ssautil.MainPackages(pkgs)
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
	p.Graph = ToGraph(g)
	p.removeMeaningLessNode()
	return nil
}

func ToGraph(g *callgraph.Graph) *Graph {
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
	return graph
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

func CheckErrors(pkgs []*packages.Package) error {
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
