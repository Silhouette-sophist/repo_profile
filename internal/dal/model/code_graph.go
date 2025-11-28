package model

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	astKeyReceiverRegex = regexp.MustCompile(`\(([^)]+)\)`)
)

// neo4j 节点类型
type (
	Repository struct {
		Repo       string
		GitRepo    string
		Branch     string
		CommitHash string
		ModuleInfo string
		UniqueId   string
	}

	AstStruct struct {
		Repo      string
		Pkg       string
		File      string
		Name      string
		Content   string
		StartLine int32
		EndLine   int32
		UniqueId  string
	}

	AstFunction struct {
		Repo      string
		Pkg       string
		File      string
		Name      string
		Content   string
		Type      string
		StartLine int32
		EndLine   int32
		Receiver  *string
		Params    []string
		Results   []string
		UniqueId  string
	}

	AstVariable struct {
		Repo      string
		Pkg       string
		File      string
		Name      string
		Content   string
		StartLine int32
		EndLine   int32
		UniqueId  string
	}
)

// neo4j 边关系关系
type (
	RelationType string

	BaseRelation struct {
		SourceID   string
		TargetID   string
		RelType    RelationType
		Properties map[string]interface{}
	}
)

const (
	InvokeRelationType    RelationType = "INVOKE"    // 函数调用关系
	ReferenceRelationType RelationType = "REFERENCE" // 函数引用包变量
	AssociateRelationType RelationType = "ASSOCIATE" // 函数关联结构体
	DependsOnRelationType RelationType = "DEPEND"    // 结构体依赖结构体
)

func (r RelationType) String() string {
	return string(r)
}

func GenerateCodeGraphFromFuncInfo(ctx context.Context, info *CodeFuncInfo) string {
	if info == nil {
		fmt.Printf("info is nil\n")
		return ""
	}
	if info.RecvType != nil {
		recvType := info.RecvType.BaseType
		if recvType == "" {
			recvType = info.RecvType.Type
		}
		lastIndex := strings.LastIndex(recvType, ".")
		var recvName string
		if lastIndex > -1 {
			recvName = recvType[lastIndex+1:]
		}
		return GenerateNodeID(info.PkgPath, "function", recvName, info.FuncName, info.BlockSpan.StartLine)
	} else {
		return GenerateNodeID(info.PkgPath, "function", "", info.FuncName, info.BlockSpan.StartLine)
	}
}

func GenerateCodeGraphFromStructInfo(ctx context.Context, info *StructInfo) string {
	if info == nil {
		fmt.Printf("info is nil\n")
		return ""
	}
	return GenerateNodeID(info.Pkg, "struct", "", info.Name, info.BlockSpan.StartLine)
}

func GenerateCodeGraphFromVariableInfo(ctx context.Context, info *FilePkgVar) string {
	if info == nil {
		fmt.Printf("info is nil\n")
		return ""
	}
	return GenerateNodeID(info.Pkg, "variable", "", info.Name, info.BlockSpan.StartLine)
}

func GenerateCodeGraphFromPkgWithName(repo, relPkg, nodeType, nodeName string) string {
	pkgPath := fmt.Sprintf("%s%s", repo, relPkg)
	return GenerateNodeID(pkgPath, nodeType, "", nodeName, 0)
}

func GenerateCodeGraphFromCallee(info *CalleeInfo) string {
	if info == nil || strings.Contains(info.Name, "$") {
		return ""
	}
	var receiver string
	if submatch := astKeyReceiverRegex.FindStringSubmatch(info.AstKey); len(submatch) >= 2 {
		receiver = submatch[1]
	}
	return GenerateNodeID(info.PkgPath, "function", receiver, info.Name, 0)
}

func GenerateNodeID(pkgPath, nodeType, receiver, nodeName string, startLine int32) string {
	// 清理特殊字符（避免ID解析错误，如冒号替换为下划线）
	cleanPkgPath := strings.ReplaceAll(pkgPath, ":", "_")
	cleanNodeType := strings.ReplaceAll(nodeType, ":", "_")
	cleanReceiver := strings.ReplaceAll(receiver, ":", "_")
	cleanNodeName := strings.ReplaceAll(nodeName, ":", "_")

	// 接收者为空时统一用"_"占位（避免歧义）
	if cleanReceiver == "" {
		cleanReceiver = "_"
	}

	if index := strings.Index(cleanNodeName, "$"); index > -1 {
		cleanNodeName = fmt.Sprintf("%s$%d", cleanNodeName[:index], startLine)
	}
	// 拼接ID（按规则顺序）
	return strings.Join(
		[]string{cleanPkgPath, cleanNodeType, cleanReceiver, cleanNodeName},
		":",
	)
}
