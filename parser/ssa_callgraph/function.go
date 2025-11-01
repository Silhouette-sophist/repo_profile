package ssa_callgraph

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

// 获取函数所属包的名称（短名称，如 "main"）
func getPackageName(fn *ssa.Function) string {
	if fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return "builtin" // 内置函数（如 println）无包信息
	}
	return fn.Pkg.Pkg.Name() // 包的短名称（如 "json"）
}

// 获取函数所属包的完整路径（如 "encoding/json"）
func getPackagePath(fn *ssa.Function) string {
	if fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return "builtin"
	}
	return fn.Pkg.Pkg.Path() // 包的完整导入路径
}

// 获取函数定义所在的文件路径
func getFunctionFile(fn *ssa.Function, pkg *packages.Package) string {
	if fn == nil || pkg == nil || pkg.Fset == nil {
		return ""
	}
	// 获取函数在语法树中的定义对象
	obj := fn.Object()
	if obj == nil {
		return ""
	}
	// 获取定义位置的 token.Pos
	pos := obj.Pos()
	if !pos.IsValid() {
		return ""
	}
	// 通过文件集（FileSet）将 pos 转换为文件路径
	position := pkg.Fset.Position(pos)
	return position.Filename
}

// 获取函数的完整名称（包含接收者类型，如 "T.Add"）
func getFunctionName(fn *ssa.Function) string {
	// 处理方法（有接收者）
	if recv := fn.Signature.Recv(); recv != nil {
		// 接收者类型（如 *T、T）
		recvType := recv.Type()
		// 简化类型显示（去掉包路径，只保留类型名）
		if named, ok := recvType.(*types.Named); ok {
			return fmt.Sprintf("%s.%s", named.Obj().Name(), fn.Name())
		}
		return fmt.Sprintf("%s.%s", recvType.String(), fn.Name())
	}
	// 普通函数（无接收者）
	return fn.Name()
}

// 获取函数的全局唯一名称（包路径+函数名，如 "example.com/mypkg.T.Add"）
func getFullFunctionName(fn *ssa.Function) string {
	pkgPath := getPackagePath(fn)
	funcName := getFunctionName(fn)
	return fmt.Sprintf("%s.%s", pkgPath, funcName)
}
