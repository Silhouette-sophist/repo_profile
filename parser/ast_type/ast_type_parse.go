package ast_type

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	vs "github.com/Silhouette-sophist/repo_profile/parser/visitor"
	"github.com/Silhouette-sophist/repo_profile/service"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

type AstTypeAnalyzer struct {
	RepoPath    string
	RootPkg     string
	PkgLoadMap  []*packages.Package   // ast load的类型
	ModuleInfos []*service.ModuleInfo // ast识别的模块
}

func (r *AstTypeAnalyzer) AnalyzeRepo(ctx context.Context) {
	fmt.Println("AstTypeAnalyzer.AnalyzeRepo")
	// 1.先加载仓库类型
	loadPackages, err := service.LoadPackages(ctx, &service.LoadConfig{
		RepoPath: r.RepoPath,
		PkgPath:  r.RootPkg,
		LoadEnum: service.LoadAllPkg,
	})
	if err != nil {
		fmt.Printf("LoadAllPackages err: %v", err)
		return
	}
	r.PkgLoadMap = loadPackages
	// 2.分析所有模块
	repo, err := r.ParseRepo()
	if err != nil {
		fmt.Printf("ParseRepo err: %v", err)
		return
	}
	fmt.Printf("AstTypeAnalyzer.AnalyzeRepo %v\n", repo)
}

// ParseRepo 匹配仓库信息
func (r *AstTypeAnalyzer) ParseRepo() ([]*service.ModuleInfo, error) {
	modules, err := r.FindAllModules()
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// ParseModule 解析单个go.mod文件
func (r *AstTypeAnalyzer) ParseModule(dir string) (*service.ModuleInfo, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("ParseModule dir:%s cost: %v\n", dir, time.Since(start))
	}()
	info := &service.ModuleInfo{
		Dir:          dir,
		PkgFuncMap:   make(map[string][]*vs.FuncInfo),
		PkgVarMap:    make(map[string][]*vs.VarInfo),
		PkgStructMap: make(map[string][]*vs.StructInfo),
	}
	// 读取go.mod文件
	modPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return info, fmt.Errorf("读取go.mod失败: %v", err)
	}
	// 解析go.mod文件
	modeFile, err := modfile.Parse(modPath, data, nil)
	if err != nil {
		return info, fmt.Errorf("解析go.mod失败: %v", err)
	}
	info.Path = modeFile.Module.Mod.Path
	if modeFile.Go != nil {
		info.GoVersion = modeFile.Go.Version
	}
	// 匹配mod文件中内容
	r.AppendModuleInfo(info)
	// 解析依赖
	for _, req := range modeFile.Require {
		info.Requires = append(info.Requires, service.Dependency{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}
	// 解析替换规则
	for _, replace := range modeFile.Replace {
		info.Replaces = append(info.Replaces, service.ReplaceRule{
			OldPath:    replace.Old.Path,
			OldVersion: replace.Old.Version,
			NewPath:    replace.New.Path,
			NewVersion: replace.New.Version,
		})
	}
	// 解析目录中所有.go文件的导入
	imports, err := r.ParseImportsFromDir(dir)
	if err != nil {
		log.Printf("警告: 解析目录 %s 中的导入失败: %v", dir, err)
	}
	info.Imports = imports
	return info, nil
}

// AppendModuleInfo 解析模块中的所有.go文件
func (r *AstTypeAnalyzer) AppendModuleInfo(modInfo *service.ModuleInfo) {
	filepath.Walk(modInfo.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		rFilePath, err := filepath.Rel(modInfo.Dir, path)
		if err != nil {
			return err
		}
		relDir, err := r.DeductRelativeDir(modInfo.Dir, path)
		if err != nil {
			return err
		}
		curPkg := modInfo.Path
		if relDir != "" {
			curPkg = modInfo.Path + "/" + relDir
		}
		var curPkgPackage *packages.Package
		for _, p := range r.PkgLoadMap {
			if p.ID == curPkg {
				curPkgPackage = p
				break
			}
		}
		fileFuncVisitor, err := service.ParseSingleFileWithPackageTypes(curPkg, rFilePath, path, curPkgPackage)
		if err != nil {
			return err
		}
		modInfo.PkgFuncMap[curPkg] = append(modInfo.PkgFuncMap[curPkg], fileFuncVisitor.FileFuncInfos...)
		modInfo.PkgVarMap[curPkg] = append(modInfo.PkgVarMap[curPkg], fileFuncVisitor.FilePkgVars...)
		modInfo.PkgStructMap[curPkg] = append(modInfo.PkgStructMap[curPkg], fileFuncVisitor.FileStructs...)
		return nil
	})
}

// DeductRelativeDir 计算子文件相对于父目录的目录路径（排除文件名）
func (r *AstTypeAnalyzer) DeductRelativeDir(parentDir, childPath string) (string, error) {
	// 计算相对路径
	relPath, err := filepath.Rel(parentDir, childPath)
	if err != nil {
		return "", err
	}
	// 排除文件名，只保留目录部分
	dir := filepath.Dir(relPath)
	// 如果结果是 "."，说明就在父目录下，返回空字符串
	if dir == "." {
		return "", nil
	}
	return dir, nil
}

// ParseImportsFromDir 从目录中的所有.go文件解析导入的包
func (r *AstTypeAnalyzer) ParseImportsFromDir(dir string) ([]service.ImportInfo, error) {
	var imports []service.ImportInfo
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过非.go文件和测试文件
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		// 解析文件
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("解析文件 %s 失败: %v", path, err)
		}
		// 提取导入
		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			var importName *string
			if imp.Name != nil {
				importName = &imp.Name.Name
			}
			imports = append(imports, service.ImportInfo{
				Path: importPath,
				Name: importName,
			})
		}
		return nil
	})
	return imports, err
}

// FindAllModules 递归查找目录中的所有模块
func (r *AstTypeAnalyzer) FindAllModules() ([]*service.ModuleInfo, error) {
	modules := make([]*service.ModuleInfo, 0)
	err := filepath.WalkDir(r.RepoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过隐藏目录
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}
		// 检查是否为go.mod文件
		if !d.IsDir() && d.Name() == "go.mod" {
			moduleDir := filepath.Dir(path)
			module, err := r.ParseModule(moduleDir)
			if err != nil {
				module.Error = err
			}
			modules = append(modules, module)
			// 跳过子目录（避免处理嵌套模块，除非需要）
			return fs.SkipDir
		}
		return nil
	})
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Path < modules[j].Path
	})
	return modules, err
}
