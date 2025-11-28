package ssa_callgraph

//
//import (
//	"context"
//	"fmt"
//	"go/token"
//	"go/types"
//	"strings"
//
//	"golang.org/x/tools/go/callgraph"
//	"golang.org/x/tools/go/callgraph/cha"
//	"golang.org/x/tools/go/callgraph/rta"
//	"golang.org/x/tools/go/callgraph/vta"
//	"golang.org/x/tools/go/packages"
//	"golang.org/x/tools/go/pointer"
//	"golang.org/x/tools/go/ssa"
//	"golang.org/x/tools/go/ssa/ssautil"
//)
//
//// SSAAnalysisResult SSA分析结果
//type SSAAnalysisResult struct {
//	RepoInfo    *RepoInfo
//	CallGraph   *CallGraph
//	PackageInfo map[string]*PackageSSAInfo
//}
//
//// RepoInfo 代码仓库信息
//type RepoInfo struct {
//	GitRepo    string
//	Branch     string
//	CommitHash string
//	ModuleInfo string
//}
//
//// CallGraph 调用图
//type CallGraph struct {
//	RootFunctions []string // 根函数列表（如main函数）
//	Functions     map[string]*FunctionSSAInfo
//	Edges         []*CallEdge
//}
//type Location struct {
//	StartLine int
//	EndLine   int
//}
//
//// FunctionSSAInfo 函数的SSA分析信息
//type FunctionSSAInfo struct {
//	ID              string
//	Name            string
//	Package         string
//	File            string
//	Location        Location
//	SSAInstructions []*SSAInstruction
//
//	// 调用关系
//	CallSites   []*CallSite // 调用点信息
//	CalledFuncs []string    // 被调用的函数ID列表
//	Callers     []string    // 调用者函数ID列表
//
//	// 数据流分析
//	UsedVars    []*VarUsage    // 使用的变量
//	UsedStructs []*StructUsage // 使用的结构体
//	DefVars     []*VarDef      // 定义的变量
//
//	// 控制流分析
//	BasicBlocks []*BasicBlock // 基本块
//	CFGEdges    []*CFGEdge    // 控制流边
//}
//
//// CallEdge 调用边
//type CallEdge struct {
//	ID         string
//	CallerID   string   // 调用者函数ID
//	CalleeID   string   // 被调用函数ID
//	CallSiteID string   // 调用点ID
//	Type       CallType // 调用类型
//}
//
//// CallSite 调用点
//type CallSite struct {
//	ID        string
//	Location  Location
//	CallExpr  string    // 调用表达式
//	CalleeID  string    // 被调用函数ID
//	InstrType InstrType // 指令类型
//}
//
//// VarUsage 变量使用信息
//type VarUsage struct {
//	VarID    string
//	VarName  string
//	VarType  string
//	Location Location
//	Usage    VarUsageType // 使用类型（读、写、地址取等）
//	InstrID  string       // 对应的SSA指令ID
//}
//
//// StructUsage 结构体使用信息
//type StructUsage struct {
//	StructID      string
//	StructName    string
//	Package       string
//	UsageType     StructUsageType // 使用类型
//	FieldAccesses []*FieldAccess  // 字段访问
//	MethodsCalled []*MethodCall   // 方法调用
//}
//
//// FieldAccess 字段访问
//type FieldAccess struct {
//	FieldName  string
//	Location   Location
//	AccessType FieldAccessType // 访问类型（读、写）
//}
//
//// MethodCall 方法调用
//type MethodCall struct {
//	MethodName string
//	Location   Location
//	ReceiverID string // 接收器变量ID
//}
//
//// BasicBlock 基本块
//type BasicBlock struct {
//	ID           string
//	Name         string
//	Location     Location
//	Instructions []string // 指令ID列表
//	Preds        []string // 前驱块ID列表
//	Succs        []string // 后继块ID列表
//}
//
//// CFGEdge 控制流边
//type CFGEdge struct {
//	FromBlockID string
//	ToBlockID   string
//	Type        CFGEdgeType // 边类型（条件真、条件假、无条件等）
//}
//
//// SSAInstruction SSA指令
//type SSAInstruction struct {
//	ID       string
//	OpCode   string // 操作码
//	Location Location
//	Operands []string  // 操作数ID列表
//	Result   string    // 结果变量ID（如果有）
//	Type     InstrType // 指令类型
//}
//
//// VarDef 变量定义
//type VarDef struct {
//	VarID    string
//	VarName  string
//	VarType  string
//	Location Location
//	InstrID  string // 定义该变量的指令ID
//}
//
//// PackageSSAInfo 包的SSA信息
//type PackageSSAInfo struct {
//	PackagePath string
//	Functions   []string         // 函数ID列表
//	Globals     []*GlobalVarInfo // 全局变量
//	Types       []*TypeInfo      // 类型信息
//	InitFunc    string           // init函数ID
//}
//
//// GlobalVarInfo 全局变量信息
//type GlobalVarInfo struct {
//	ID       string
//	Name     string
//	Type     string
//	Location Location
//	Uses     []string // 使用该全局变量的函数ID列表
//}
//
//// TypeInfo 类型信息
//type TypeInfo struct {
//	ID         string
//	Name       string
//	Package    string
//	Kind       TypeKind // 类型种类（结构体、接口等）
//	Methods    []*MethodInfo
//	Underlying string // 底层类型
//}
//
//// MethodInfo 方法信息
//type MethodInfo struct {
//	Name     string
//	Receiver string
//	FuncID   string // 对应的函数ID
//	Location Location
//}
//
//// 枚举类型定义
//type CallType int
//
//const (
//	CallStatic    CallType = iota // 静态调用
//	CallDynamic                   // 动态调用
//	CallInterface                 // 接口调用
//	CallClosure                   // 闭包调用
//)
//
//type InstrType int
//
//const (
//	InstrCall  InstrType = iota // 调用指令
//	InstrAlloc                  // 分配指令
//	InstrLoad                   // 加载指令
//	InstrStore                  // 存储指令
//	InstrBinOp                  // 二元操作
//	InstrUnOp                   // 一元操作
//	InstrPhi                    // Phi指令
//)
//
//type VarUsageType int
//
//const (
//	VarRead  VarUsageType = iota // 读
//	VarWrite                     // 写
//	VarAddr                      // 取地址
//)
//
//type StructUsageType int
//
//const (
//	StructCreate      StructUsageType = iota // 创建实例
//	StructFieldAccess                        // 字段访问
//	StructMethodCall                         // 方法调用
//	StructEmbedded                           // 嵌入结构体
//)
//
//type FieldAccessType int
//
//const (
//	FieldRead  FieldAccessType = iota // 字段读
//	FieldWrite                        // 字段写
//)
//
//type CFGEdgeType int
//
//const (
//	EdgeUnconditional CFGEdgeType = iota // 无条件跳转
//	EdgeTrue                             // 条件真
//	EdgeFalse                            // 条件假
//)
//
//type TypeKind int
//
//const (
//	TypeStruct    TypeKind = iota // 结构体类型
//	TypeInterface                 // 接口类型
//	TypeBasic                     // 基本类型
//	TypeSlice                     // 切片类型
//	TypeArray                     // 数组类型
//	TypeMap                       // 映射类型
//	TypePointer                   // 指针类型
//)
//
//// SSAAnalyzer SSA分析器
//type SSAAnalyzer struct {
//	prog      *ssa.Program
//	callGraph *callgraph.Graph
//	packages  []*ssa.Package
//	mainPkgs  []*ssa.Package
//	result    *SSAAnalysisResult
//	funcMap   map[string]*FunctionSSAInfo // 函数ID到信息的映射
//	varMap    map[string]*GlobalVarInfo   // 变量ID到信息的映射
//	typeMap   map[string]*TypeInfo        // 类型ID到信息的映射
//}
//
//// NewSSAAnalyzer 创建SSA分析器
//func NewSSAAnalyzer() *SSAAnalyzer {
//	return &SSAAnalyzer{
//		result:  &SSAAnalysisResult{},
//		funcMap: make(map[string]*FunctionSSAInfo),
//		varMap:  make(map[string]*GlobalVarInfo),
//		typeMap: make(map[string]*TypeInfo),
//	}
//}
//
//// Analyze 执行SSA分析
//func (a *SSAAnalyzer) Analyze(ctx context.Context, args SSAAnalyzerParam) (*SSAAnalysisResult, error) {
//	// 1. 加载包并构建SSA
//	if err := a.loadProgram(ctx, args); err != nil {
//		return nil, err
//	}
//
//	// 2. 构建调用图
//	if err := a.buildCallGraph(ctx, args); err != nil {
//		return nil, err
//	}
//
//	// 3. 收集包级别信息
//	a.collectPackageInfo()
//
//	// 4. 收集函数级别信息
//	a.collectFunctionInfo()
//
//	// 5. 构建调用关系
//	a.buildCallRelationships()
//
//	return a.result, nil
//}
//
//// loadProgram 加载程序并构建SSA
//func (a *SSAAnalyzer) loadProgram(ctx context.Context, args SSAAnalyzerParam) error {
//	cfg := &packages.Config{
//		Mode:  packages.LoadAllSyntax,
//		Tests: false,
//		Dir:   args.Path,
//	}
//
//	initialPkgs, err := packages.Load(cfg, "./...")
//	if err != nil {
//		return fmt.Errorf("failed to load packages: %v", err)
//	}
//
//	// 检查错误
//	if err := a.checkErrors(initialPkgs); err != nil {
//		return err
//	}
//
//	// 创建SSA程序
//	a.prog, a.packages = ssautil.AllPackages(initialPkgs, ssa.GlobalDebug)
//	a.prog.Build()
//
//	// 获取main包
//	a.mainPkgs = ssautil.MainPackages(a.packages)
//
//	return nil
//}
//
//type SSAAnalyzerParam struct {
//	GitRepo    string
//	Branch     string
//	CommitHash string
//	ModuleInfo string
//	Algorithm  string
//	Path       string
//}
//
//// buildCallGraph 构建调用图
//func (a *SSAAnalyzer) buildCallGraph(ctx context.Context, args SSAAnalyzerParam) error {
//	switch args.Algorithm {
//	case "cha":
//		a.callGraph = cha.CallGraph(a.prog)
//	case "rta":
//		roots := a.getRootFunctions()
//		res := rta.Analyze(roots, true)
//		a.callGraph = res.CallGraph
//	case "vta":
//		g := cha.CallGraph(a.prog)
//		a.callGraph = vta.CallGraph(ssautil.AllFunctions(a.prog), g)
//	case "pta":
//		result, err := pointer.Analyze(&pointer.Config{
//			Mains:          a.mainPkgs,
//			BuildCallGraph: true,
//		})
//		if err != nil {
//			return err
//		}
//		a.callGraph = result.CallGraph
//	default:
//		a.callGraph = cha.CallGraph(a.prog) // 默认使用CHA
//	}
//
//	// 初始化结果结构
//	a.result = &SSAAnalysisResult{
//		RepoInfo: &RepoInfo{
//			GitRepo:    args.GitRepo,
//			Branch:     args.Branch,
//			CommitHash: args.CommitHash,
//			ModuleInfo: args.ModuleInfo,
//		},
//		CallGraph: &CallGraph{
//			Functions: make(map[string]*FunctionSSAInfo),
//			Edges:     []*CallEdge{},
//		},
//		PackageInfo: make(map[string]*PackageSSAInfo),
//	}
//
//	return nil
//}
//
//// getRootFunctions 获取根函数
//func (a *SSAAnalyzer) getRootFunctions() []*ssa.Function {
//	var roots []*ssa.Function
//	for _, mainPkg := range a.mainPkgs {
//		if mainFunc := mainPkg.Func("main"); mainFunc != nil {
//			roots = append(roots, mainFunc)
//		}
//		if initFunc := mainPkg.Func("init"); initFunc != nil {
//			roots = append(roots, initFunc)
//		}
//	}
//	return roots
//}
//
//// collectPackageInfo 收集包级别信息
//func (a *SSAAnalyzer) collectPackageInfo() {
//	for _, pkg := range a.packages {
//		pkgInfo := &PackageSSAInfo{
//			PackagePath: pkg.Pkg.Path(),
//			Functions:   []string{},
//			Globals:     []*GlobalVarInfo{},
//			Types:       []*TypeInfo{},
//		}
//
//		// 收集全局变量
//		for _, member := range pkg.Members {
//			switch val := member.(type) {
//			case *ssa.Global:
//				globalInfo := a.collectGlobalVarInfo(val, pkg.Pkg.Path())
//				pkgInfo.Globals = append(pkgInfo.Globals, globalInfo)
//				a.varMap[globalInfo.ID] = globalInfo
//
//			case *ssa.Function:
//				funcID := a.getFunctionID(val)
//				pkgInfo.Functions = append(pkgInfo.Functions, funcID)
//
//				// 检查是否是init函数
//				if val.Name() == "init" {
//					pkgInfo.InitFunc = funcID
//				}
//			}
//		}
//
//		// 收集类型信息
//		a.collectTypeInfo(pkg, pkgInfo)
//
//		a.result.PackageInfo[pkg.Pkg.Path()] = pkgInfo
//	}
//}
//
//// collectGlobalVarInfo 收集全局变量信息
//func (a *SSAAnalyzer) collectGlobalVarInfo(global *ssa.Global, pkgPath string) *GlobalVarInfo {
//	return &GlobalVarInfo{
//		ID:   a.generateVarID(global),
//		Name: global.Name(),
//		Type: global.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(global.Pos()),
//			EndLine:   getLineNumber(global.Pos()),
//		},
//		Uses: []string{},
//	}
//}
//
//// collectTypeInfo 收集类型信息
//func (a *SSAAnalyzer) collectTypeInfo(pkg *ssa.Package, pkgInfo *PackageSSAInfo) {
//	scope := pkg.Pkg.Scope()
//	for _, name := range scope.Names() {
//		obj := scope.Lookup(name)
//		if typeName, ok := obj.(*types.TypeName); ok {
//			typeInfo := a.createTypeInfo(typeName, pkg.Pkg.Path())
//			if typeInfo != nil {
//				pkgInfo.Types = append(pkgInfo.Types, typeInfo)
//				a.typeMap[typeInfo.ID] = typeInfo
//			}
//		}
//	}
//}
//
//// createTypeInfo 创建类型信息
//func (a *SSAAnalyzer) createTypeInfo(typeName *types.TypeName, pkgPath string) *TypeInfo {
//	typeInfo := &TypeInfo{
//		ID:      a.generateTypeID(typeName),
//		Name:    typeName.Name(),
//		Package: pkgPath,
//		Methods: []*MethodInfo{},
//	}
//
//	// 确定类型种类
//	switch typeName.Type().Underlying().(type) {
//	case *types.Struct:
//		typeInfo.Kind = TypeStruct
//	case *types.Interface:
//		typeInfo.Kind = TypeInterface
//	case *types.Slice:
//		typeInfo.Kind = TypeSlice
//	case *types.Array:
//		typeInfo.Kind = TypeArray
//	case *types.Map:
//		typeInfo.Kind = TypeMap
//	case *types.Pointer:
//		typeInfo.Kind = TypePointer
//	default:
//		typeInfo.Kind = TypeBasic
//	}
//
//	typeInfo.Underlying = typeName.Type().Underlying().String()
//
//	// 收集方法信息
//	a.collectMethodInfo(typeName, typeInfo)
//
//	return typeInfo
//}
//
//// collectMethodInfo 收集方法信息
//func (a *SSAAnalyzer) collectMethodInfo(typeName *types.TypeName, typeInfo *TypeInfo) {
//	namedType := typeName.Type().(*types.Named)
//	for i := 0; i < namedType.NumMethods(); i++ {
//		method := namedType.Method(i)
//		if methodFunc := a.prog.FuncValue(method); methodFunc != nil {
//			methodInfo := &MethodInfo{
//				Name:     method.Name(),
//				Receiver: typeInfo.Name,
//				FuncID:   a.getFunctionID(methodFunc),
//				Location: Location{
//					StartLine: getLineNumber(method.Pos()),
//					EndLine:   getLineNumber(method.Pos()),
//				},
//			}
//			typeInfo.Methods = append(typeInfo.Methods, methodInfo)
//		}
//	}
//}
//
//// collectFunctionInfo 收集函数级别信息
//func (a *SSAAnalyzer) collectFunctionInfo() {
//	// 遍历调用图中的所有函数
//	a.callGraph.DeleteSyntheticNodes()
//
//	for _, node := range a.callGraph.Nodes {
//		if node.Func == nil {
//			continue
//		}
//
//		funcInfo := a.createFunctionInfo(node.Func)
//		a.result.CallGraph.Functions[funcInfo.ID] = funcInfo
//		a.funcMap[funcInfo.ID] = funcInfo
//
//		// 分析函数的SSA指令
//		a.analyzeFunctionInstructions(node.Func, funcInfo)
//	}
//}
//
//// createFunctionInfo 创建函数基本信息
//func (a *SSAAnalyzer) createFunctionInfo(fn *ssa.Function) *FunctionSSAInfo {
//	funcInfo := &FunctionSSAInfo{
//		ID:              a.getFunctionID(fn),
//		Name:            fn.Name(),
//		Package:         fn.Pkg.Pkg.Path(),
//		Location:        a.getFunctionLocation(fn),
//		SSAInstructions: []*SSAInstruction{},
//		CallSites:       []*CallSite{},
//		CalledFuncs:     []string{},
//		Callers:         []string{},
//		UsedVars:        []*VarUsage{},
//		UsedStructs:     []*StructUsage{},
//		DefVars:         []*VarDef{},
//		BasicBlocks:     []*BasicBlock{},
//		CFGEdges:        []*CFGEdge{},
//	}
//
//	// 设置文件路径
//	if fn.Pos() != 0 {
//		// 这里需要根据pos获取文件路径，具体实现取决于您的文件系统
//		funcInfo.File = a.getFilePath(fn.Pos())
//	}
//
//	return funcInfo
//}
//
//// analyzeFunctionInstructions 分析函数的SSA指令
//func (a *SSAAnalyzer) analyzeFunctionInstructions(fn *ssa.Function, funcInfo *FunctionSSAInfo) {
//	if fn.Blocks == nil {
//		return
//	}
//
//	// 分析基本块
//	for i, block := range fn.Blocks {
//		blockInfo := a.analyzeBasicBlock(block, i)
//		funcInfo.BasicBlocks = append(funcInfo.BasicBlocks, blockInfo)
//
//		// 分析控制流边
//		a.analyzeCFGEdges(block, funcInfo)
//	}
//
//	// 分析指令
//	for _, block := range fn.Blocks {
//		for _, instr := range block.Instrs {
//			instruction := a.analyzeInstruction(instr, funcInfo)
//			if instruction != nil {
//				funcInfo.SSAInstructions = append(funcInfo.SSAInstructions, instruction)
//			}
//		}
//	}
//}
//
//// analyzeBasicBlock 分析基本块
//func (a *SSAAnalyzer) analyzeBasicBlock(block *ssa.BasicBlock, index int) *BasicBlock {
//	blockInfo := &BasicBlock{
//		ID:           fmt.Sprintf("block_%d", index),
//		Name:         fmt.Sprintf("b%d", index),
//		Instructions: []string{},
//		Preds:        []string{},
//		Succs:        []string{},
//	}
//
//	// 计算位置信息
//	if len(block.Instrs) > 0 {
//		startInstr := block.Instrs[0]
//		endInstr := block.Instrs[len(block.Instrs)-1]
//		blockInfo.Location = Location{
//			StartLine: getLineNumber(startInstr.Pos()),
//			EndLine:   getLineNumber(endInstr.Pos()),
//		}
//	}
//
//	return blockInfo
//}
//
//// analyzeCFGEdges 分析控制流边
//func (a *SSAAnalyzer) analyzeCFGEdges(block *ssa.BasicBlock, funcInfo *FunctionSSAInfo) {
//	fromBlockID := fmt.Sprintf("block_%d", block.Index)
//
//	for _, succ := range block.Succs {
//		toBlockID := fmt.Sprintf("block_%d", succ.Index)
//		edge := &CFGEdge{
//			FromBlockID: fromBlockID,
//			ToBlockID:   toBlockID,
//			Type:        EdgeUnconditional, // 简化处理，实际应根据条件分支判断
//		}
//		funcInfo.CFGEdges = append(funcInfo.CFGEdges, edge)
//	}
//}
//
//// analyzeAllocInstruction 分析内存分配指令
//func (a *SSAAnalyzer) analyzeAllocInstruction(alloc *ssa.Alloc, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Alloc"
//	instruction.Type = InstrAlloc
//
//	// 记录变量定义
//	if alloc.Comment != "" {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(alloc),
//			VarName:  alloc.Comment,
//			VarType:  alloc.Type().String(),
//			Location: Location{StartLine: getLineNumber(alloc.Pos()), EndLine: getLineNumber(alloc.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//	}
//
//	// 分析分配的类型
//	if typeInfo := a.extractTypeInfo(alloc.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructCreate,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//}
//
//// analyzeFieldAccess 分析字段访问
//func (a *SSAAnalyzer) analyzeFieldAccess(fieldAddr *ssa.FieldAddr, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "FieldAddr"
//
//	// 分析结构体使用
//	if typeInfo := a.extractTypeInfo(fieldAddr.X.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//			FieldAccesses: []*FieldAccess{
//				{
//					FieldName:  fmt.Sprintf("Field%d", fieldAddr.Field),
//					Location:   Location{StartLine: getLineNumber(fieldAddr.Pos()), EndLine: getLineNumber(fieldAddr.Pos())},
//					AccessType: FieldRead, // 简化处理
//				},
//			},
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//}
//
//// 辅助函数
//func (a *SSAAnalyzer) getFunctionID(fn *ssa.Function) string {
//	return fmt.Sprintf("%s.%s", fn.Pkg.Pkg.Path(), fn.Name())
//}
//
//func (a *SSAAnalyzer) generateVarID(global *ssa.Global) string {
//	return fmt.Sprintf("var_%s_%s", global.Pkg.Pkg.Path(), global.Name())
//}
//
//func (a *SSAAnalyzer) generateTypeID(typeName *types.TypeName) string {
//	return fmt.Sprintf("type_%s_%s", typeName.Pkg().Path(), typeName.Name())
//}
//
//func (a *SSAAnalyzer) generateInstrID(instr ssa.Instruction) string {
//	return fmt.Sprintf("instr_%d", instr.Pos())
//}
//
//func (a *SSAAnalyzer) generateCallSiteID(call *ssa.Call) string {
//	return fmt.Sprintf("callsite_%d", call.Pos())
//}
//
//func (a *SSAAnalyzer) generateCallEdgeID(edge *callgraph.Edge) string {
//	if edge.Site != nil {
//		return fmt.Sprintf("edge_%d", edge.Site.Pos())
//	}
//	return fmt.Sprintf("edge_%p", edge)
//}
//
//func (a *SSAAnalyzer) generateVarIDFromValue(value ssa.Value) string {
//	return fmt.Sprintf("localvar_%d", value.Pos())
//}
//
//func (a *SSAAnalyzer) getFunctionLocation(fn *ssa.Function) Location {
//	return Location{
//		StartLine: getLineNumber(fn.Pos()),
//		EndLine:   getLineNumber(fn.Pos()), // 简化处理
//	}
//}
//
//func (a *SSAAnalyzer) getFilePath(pos token.Pos) string {
//	// 这里需要根据您的文件系统实现来获取文件路径
//	// 简化返回
//	return "unknown"
//}
//
//// checkErrors 检查包加载错误
//func (a *SSAAnalyzer) checkErrors(pkgs []*packages.Package) error {
//	for _, pkg := range pkgs {
//		if len(pkg.Errors) > 0 {
//			return fmt.Errorf("package %s has errors: %v", pkg.PkgPath, pkg.Errors)
//		}
//	}
//	return nil
//}
//
//// getLineNumber 从token.Pos获取行号（简化实现）
//func getLineNumber(pos token.Pos) int {
//	// 这里需要根据您的token.FileSet来获取准确的行号
//	// 简化返回
//	return 0
//}
//
//// analyzeBinaryOp 分析二元操作指令
//func (a *SSAAnalyzer) analyzeBinaryOp(binop *ssa.BinOp, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "BinOp"
//	instruction.Type = InstrBinOp
//	instruction.Operands = []string{
//		a.generateOperandID(binop.X),
//		a.generateOperandID(binop.Y),
//	}
//
//	// 分析操作数类型
//	operands := []ssa.Value{binop.X, binop.Y}
//	for _, operand := range operands {
//		if typeInfo := a.extractTypeInfo(operand.Type()); typeInfo != nil {
//			structUsage := &StructUsage{
//				StructID:   typeInfo.ID,
//				StructName: typeInfo.Name,
//				Package:    typeInfo.Package,
//				UsageType:  StructFieldAccess,
//			}
//			funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//		}
//
//		// 记录变量使用
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(operand),
//			VarName: a.getVarName(operand),
//			VarType: operand.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(binop.Pos()),
//				EndLine:   getLineNumber(binop.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	// 记录结果变量定义
//	if binop.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(binop),
//			VarName:  a.getVarName(binop),
//			VarType:  binop.Type().String(),
//			Location: Location{StartLine: getLineNumber(binop.Pos()), EndLine: getLineNumber(binop.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeStoreInstruction 分析存储指令
//func (a *SSAAnalyzer) analyzeStoreInstruction(store *ssa.Store, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Store"
//	instruction.Type = InstrStore
//	instruction.Operands = []string{
//		a.generateOperandID(store.Addr),
//		a.generateOperandID(store.Val),
//	}
//
//	// 分析地址类型（通常包含结构体信息）
//	if typeInfo := a.extractTypeInfo(store.Addr.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用（地址）
//	addrUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(store.Addr),
//		VarName: a.getVarName(store.Addr),
//		VarType: store.Addr.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(store.Pos()),
//			EndLine:   getLineNumber(store.Pos()),
//		},
//		Usage:   VarWrite,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, addrUsage)
//
//	// 记录变量使用（值）
//	valUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(store.Val),
//		VarName: a.getVarName(store.Val),
//		VarType: store.Val.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(store.Pos()),
//			EndLine:   getLineNumber(store.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, valUsage)
//
//	// 分析值的类型
//	if typeInfo := a.extractTypeInfo(store.Val.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//}
//
//// analyzeLoadInstruction 分析加载指令
//func (a *SSAAnalyzer) analyzeLoadInstruction(load *ssa.Load, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Load"
//	instruction.Type = InstrLoad
//	instruction.Operands = []string{a.generateOperandID(load.X)}
//
//	// 分析地址类型
//	if typeInfo := a.extractTypeInfo(load.X.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用（地址）
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(load.X),
//		VarName: a.getVarName(load.X),
//		VarType: load.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(load.Pos()),
//			EndLine:   getLineNumber(load.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if load.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(load),
//			VarName:  a.getVarName(load),
//			VarType:  load.Type().String(),
//			Location: Location{StartLine: getLineNumber(load.Pos()), EndLine: getLineNumber(load.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//
//	// 分析加载的类型
//	if typeInfo := a.extractTypeInfo(load.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//}
//
//// analyzePhiInstruction 分析Phi指令（SSA中的φ函数）
//func (a *SSAAnalyzer) analyzePhiInstruction(phi *ssa.Phi, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Phi"
//	instruction.Type = InstrPhi
//
//	// 收集所有操作数
//	for _, edge := range phi.Edges {
//		instruction.Operands = append(instruction.Operands, a.generateOperandID(edge))
//
//		// 记录变量使用
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(edge),
//			VarName: a.getVarName(edge),
//			VarType: edge.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(phi.Pos()),
//				EndLine:   getLineNumber(phi.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	// 记录结果变量定义
//	if phi.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(phi),
//			VarName:  a.getVarName(phi),
//			VarType:  phi.Type().String(),
//			Location: Location{StartLine: getLineNumber(phi.Pos()), EndLine: getLineNumber(phi.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeExtractInstruction 分析提取指令（从元组或结构中提取值）
//func (a *SSAAnalyzer) analyzeExtractInstruction(extract *ssa.Extract, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Extract"
//	instruction.Operands = []string{a.generateOperandID(extract.Tuple)}
//
//	// 分析元组类型
//	if typeInfo := a.extractTypeInfo(extract.Tuple.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//			FieldAccesses: []*FieldAccess{
//				{
//					FieldName:  fmt.Sprintf("Index%d", extract.Index),
//					Location:   Location{StartLine: getLineNumber(extract.Pos()), EndLine: getLineNumber(extract.Pos())},
//					AccessType: FieldRead,
//				},
//			},
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(extract.Tuple),
//		VarName: a.getVarName(extract.Tuple),
//		VarType: extract.Tuple.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(extract.Pos()),
//			EndLine:   getLineNumber(extract.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if extract.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(extract),
//			VarName:  a.getVarName(extract),
//			VarType:  extract.Type().String(),
//			Location: Location{StartLine: getLineNumber(extract.Pos()), EndLine: getLineNumber(extract.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeMakeInterfaceInstruction 分析创建接口指令
//func (a *SSAAnalyzer) analyzeMakeInterfaceInstruction(mi *ssa.MakeInterface, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "MakeInterface"
//	instruction.Operands = []string{a.generateOperandID(mi.X)}
//
//	// 分析具体类型和接口类型
//	if concreteType := a.extractTypeInfo(mi.X.Type()); concreteType != nil {
//		structUsage := &StructUsage{
//			StructID:   concreteType.ID,
//			StructName: concreteType.Name,
//			Package:    concreteType.Package,
//			UsageType:  StructEmbedded, // 接口实现可以看作是一种嵌入
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	if ifaceType := a.extractTypeInfo(mi.Type()); ifaceType != nil {
//		structUsage := &StructUsage{
//			StructID:   ifaceType.ID,
//			StructName: ifaceType.Name,
//			Package:    ifaceType.Package,
//			UsageType:  StructEmbedded,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(mi.X),
//		VarName: a.getVarName(mi.X),
//		VarType: mi.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(mi.Pos()),
//			EndLine:   getLineNumber(mi.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if mi.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(mi),
//			VarName:  a.getVarName(mi),
//			VarType:  mi.Type().String(),
//			Location: Location{StartLine: getLineNumber(mi.Pos()), EndLine: getLineNumber(mi.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeTypeAssertInstruction 分析类型断言指令
//func (a *SSAAnalyzer) analyzeTypeAssertInstruction(assert *ssa.TypeAssert, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "TypeAssert"
//	instruction.Operands = []string{a.generateOperandID(assert.X)}
//
//	// 分析断言的目标类型
//	if typeInfo := a.extractTypeInfo(assert.AssertedType); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(assert.X),
//		VarName: a.getVarName(assert.X),
//		VarType: assert.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(assert.Pos()),
//			EndLine:   getLineNumber(assert.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if assert.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(assert),
//			VarName:  a.getVarName(assert),
//			VarType:  assert.Type().String(),
//			Location: Location{StartLine: getLineNumber(assert.Pos()), EndLine: getLineNumber(assert.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeChangeInterfaceInstruction 分析接口转换指令
//func (a *SSAAnalyzer) analyzeChangeInterfaceInstruction(change *ssa.ChangeInterface, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "ChangeInterface"
//	instruction.Operands = []string{a.generateOperandID(change.X)}
//
//	// 分析源接口和目标接口类型
//	if srcType := a.extractTypeInfo(change.X.Type()); srcType != nil {
//		structUsage := &StructUsage{
//			StructID:   srcType.ID,
//			StructName: srcType.Name,
//			Package:    srcType.Package,
//			UsageType:  StructEmbedded,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	if dstType := a.extractTypeInfo(change.Type()); dstType != nil {
//		structUsage := &StructUsage{
//			StructID:   dstType.ID,
//			StructName: dstType.Name,
//			Package:    dstType.Package,
//			UsageType:  StructEmbedded,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(change.X),
//		VarName: a.getVarName(change.X),
//		VarType: change.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(change.Pos()),
//			EndLine:   getLineNumber(change.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if change.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(change),
//			VarName:  a.getVarName(change),
//			VarType:  change.Type().String(),
//			Location: Location{StartLine: getLineNumber(change.Pos()), EndLine: getLineNumber(change.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeChangeTypeInstruction 分析类型转换指令
//func (a *SSAAnalyzer) analyzeChangeTypeInstruction(change *ssa.ChangeType, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "ChangeType"
//	instruction.Operands = []string{a.generateOperandID(change.X)}
//
//	// 分析源类型和目标类型
//	if srcType := a.extractTypeInfo(change.X.Type()); srcType != nil {
//		structUsage := &StructUsage{
//			StructID:   srcType.ID,
//			StructName: srcType.Name,
//			Package:    srcType.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	if dstType := a.extractTypeInfo(change.Type()); dstType != nil {
//		structUsage := &StructUsage{
//			StructID:   dstType.ID,
//			StructName: dstType.Name,
//			Package:    dstType.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(change.X),
//		VarName: a.getVarName(change.X),
//		VarType: change.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(change.Pos()),
//			EndLine:   getLineNumber(change.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if change.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(change),
//			VarName:  a.getVarName(change),
//			VarType:  change.Type().String(),
//			Location: Location{StartLine: getLineNumber(change.Pos()), EndLine: getLineNumber(change.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// 辅助函数
//func (a *SSAAnalyzer) generateOperandID(value ssa.Value) string {
//	return fmt.Sprintf("operand_%d", value.Pos())
//}
//
//func (a *SSAAnalyzer) getVarName(value ssa.Value) string {
//	if value.Name() != "" {
//		return value.Name()
//	}
//
//	// 对于匿名变量，生成一个描述性名称
//	switch v := value.(type) {
//	case *ssa.Function:
//		return v.Name()
//	case *ssa.Global:
//		return v.Name()
//	case *ssa.Const:
//		return v.String()
//	default:
//		return fmt.Sprintf("var_%d", value.Pos())
//	}
//}
//
//// extractTypeInfo 提取类型信息（增强版）
//func (a *SSAAnalyzer) extractTypeInfo(typ types.Type) *TypeInfo {
//	// 处理指针类型
//	if ptr, ok := typ.(*types.Pointer); ok {
//		return a.extractTypeInfo(ptr.Elem())
//	}
//
//	// 处理命名类型
//	if named, ok := typ.(*types.Named); ok {
//		typeName := named.Obj()
//		typeID := a.generateTypeID(typeName)
//
//		if existing, exists := a.typeMap[typeID]; exists {
//			return existing
//		}
//
//		// 创建新的类型信息
//		typeInfo := &TypeInfo{
//			ID:         typeID,
//			Name:       typeName.Name(),
//			Package:    typeName.Pkg().Path(),
//			Underlying: named.Underlying().String(),
//		}
//
//		// 确定类型种类
//		switch named.Underlying().(type) {
//		case *types.Struct:
//			typeInfo.Kind = TypeStruct
//		case *types.Interface:
//			typeInfo.Kind = TypeInterface
//		case *types.Slice:
//			typeInfo.Kind = TypeSlice
//		case *types.Array:
//			typeInfo.Kind = TypeArray
//		case *types.Map:
//			typeInfo.Kind = TypeMap
//		default:
//			typeInfo.Kind = TypeBasic
//		}
//
//		a.typeMap[typeID] = typeInfo
//		return typeInfo
//	}
//
//	// 处理结构体类型（匿名结构体）
//	if structType, ok := typ.(*types.Struct); ok {
//		// 为匿名结构体生成ID
//		typeID := fmt.Sprintf("anon_struct_%p", structType)
//		typeInfo := &TypeInfo{
//			ID:         typeID,
//			Name:       "anonymous",
//			Package:    "unknown", // 匿名结构体没有包信息
//			Kind:       TypeStruct,
//			Underlying: structType.String(),
//		}
//		return typeInfo
//	}
//
//	return nil
//}
//
//// analyzeLoadOperation 分析加载操作（通过解引用实现）
//func (a *SSAAnalyzer) analyzeLoadOperation(unop *ssa.UnOp, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	// 检查是否是解引用操作（加载）
//	if unop.Op != token.MUL {
//		return // 不是加载操作
//	}
//
//	instruction.OpCode = "Load"
//	instruction.Type = InstrLoad
//	instruction.Operands = []string{a.generateOperandID(unop.X)}
//
//	// 分析地址类型
//	if typeInfo := a.extractTypeInfo(unop.X.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用（地址）
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(unop.X),
//		VarName: a.getVarName(unop.X),
//		VarType: unop.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(unop.Pos()),
//			EndLine:   getLineNumber(unop.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if unop.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(unop),
//			VarName:  a.getVarName(unop),
//			VarType:  unop.Type().String(),
//			Location: Location{StartLine: getLineNumber(unop.Pos()), EndLine: getLineNumber(unop.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//
//	// 分析加载的类型
//	if typeInfo := a.extractTypeInfo(unop.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//}
//
//// analyzeUnaryOp 分析一元操作指令（更新版本）
//func (a *SSAAnalyzer) analyzeUnaryOp(unop *ssa.UnOp, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	// 检查是否是加载操作
//	if unop.Op == token.MUL {
//		a.analyzeLoadOperation(unop, funcInfo, instruction)
//		return
//	}
//
//	// 其他一元操作
//	instruction.OpCode = "UnOp"
//	instruction.Type = InstrUnOp
//	instruction.Operands = []string{a.generateOperandID(unop.X)}
//
//	// 设置操作符
//	switch unop.Op {
//	case token.SUB:
//		instruction.OpCode = "Neg"
//	case token.NOT:
//		instruction.OpCode = "Not"
//	case token.XOR:
//		instruction.OpCode = "Xor"
//	case token.AND:
//		instruction.OpCode = "Addr" // 取地址操作
//	}
//
//	// 分析操作数类型
//	if typeInfo := a.extractTypeInfo(unop.X.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(unop.X),
//		VarName: a.getVarName(unop.X),
//		VarType: unop.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(unop.Pos()),
//			EndLine:   getLineNumber(unop.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 处理取地址操作的特殊情况
//	if unop.Op == token.AND {
//		addrUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(unop.X),
//			VarName: a.getVarName(unop.X),
//			VarType: unop.X.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(unop.Pos()),
//				EndLine:   getLineNumber(unop.Pos()),
//			},
//			Usage:   VarAddr, // 取地址操作
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, addrUsage)
//	}
//
//	// 如果有结果，记录变量定义
//	if unop.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(unop),
//			VarName:  a.getVarName(unop),
//			VarType:  unop.Type().String(),
//			Location: Location{StartLine: getLineNumber(unop.Pos()), EndLine: getLineNumber(unop.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeInstruction 分析SSA指令（修正版本）
//func (a *SSAAnalyzer) analyzeInstruction(instr ssa.Instruction, funcInfo *FunctionSSAInfo) *SSAInstruction {
//	instruction := &SSAInstruction{
//		ID:       a.generateInstrID(instr),
//		Location: Location{StartLine: getLineNumber(instr.Pos()), EndLine: getLineNumber(instr.Pos())},
//		Operands: []string{},
//	}
//
//	switch instr := instr.(type) {
//	case *ssa.Call:
//		a.analyzeCallInstruction(instr, funcInfo, instruction)
//	case *ssa.Alloc:
//		a.analyzeAllocInstruction(instr, funcInfo, instruction)
//	case *ssa.FieldAddr:
//		a.analyzeFieldAccess(instr, funcInfo, instruction)
//	case *ssa.UnOp:
//		a.analyzeUnaryOp(instr, funcInfo, instruction) // 这里会处理加载操作
//	case *ssa.BinOp:
//		a.analyzeBinaryOp(instr, funcInfo, instruction)
//	case *ssa.Store:
//		a.analyzeStoreInstruction(instr, funcInfo, instruction)
//	case *ssa.Phi:
//		a.analyzePhiInstruction(instr, funcInfo, instruction)
//	case *ssa.Extract:
//		a.analyzeExtractInstruction(instr, funcInfo, instruction)
//	case *ssa.MakeInterface:
//		a.analyzeMakeInterfaceInstruction(instr, funcInfo, instruction)
//	case *ssa.TypeAssert:
//		a.analyzeTypeAssertInstruction(instr, funcInfo, instruction)
//	case *ssa.ChangeInterface:
//		a.analyzeChangeInterfaceInstruction(instr, funcInfo, instruction)
//	case *ssa.ChangeType:
//		a.analyzeChangeTypeInstruction(instr, funcInfo, instruction)
//	case *ssa.MakeSlice:
//		a.analyzeMakeSliceInstruction(instr, funcInfo, instruction)
//	case *ssa.Slice:
//		a.analyzeSliceInstruction(instr, funcInfo, instruction)
//	case *ssa.IndexAddr:
//		a.analyzeIndexAddrInstruction(instr, funcInfo, instruction)
//	case *ssa.Index:
//		a.analyzeIndexInstruction(instr, funcInfo, instruction)
//	case *ssa.Lookup:
//		a.analyzeLookupInstruction(instr, funcInfo, instruction)
//	case *ssa.Range:
//		a.analyzeRangeInstruction(instr, funcInfo, instruction)
//	case *ssa.Next:
//		a.analyzeNextInstruction(instr, funcInfo, instruction)
//	case *ssa.Convert:
//		a.analyzeConvertInstruction(instr, funcInfo, instruction)
//	case *ssa.MakeMap:
//		a.analyzeMakeMapInstruction(instr, funcInfo, instruction)
//	case *ssa.MapUpdate:
//		a.analyzeMapUpdateInstruction(instr, funcInfo, instruction)
//	case *ssa.DebugRef:
//		// 调试引用，通常可以忽略
//		return nil
//	default:
//		// 处理其他未明确处理的指令类型
//		a.analyzeGenericInstruction(instr, funcInfo, instruction)
//	}
//
//	return instruction
//}
//
//// 新增的其他SSA指令分析函数
//
//// analyzeMakeSliceInstruction 分析创建切片指令
//func (a *SSAAnalyzer) analyzeMakeSliceInstruction(makeslice *ssa.MakeSlice, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "MakeSlice"
//	instruction.Operands = []string{
//		a.generateOperandID(makeslice.Len),
//		a.generateOperandID(makeslice.Cap),
//	}
//
//	// 分析切片元素类型
//	if typeInfo := a.extractTypeInfo(makeslice.Type().(*types.Slice).Elem()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录结果变量定义
//	if makeslice.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(makeslice),
//			VarName:  a.getVarName(makeslice),
//			VarType:  makeslice.Type().String(),
//			Location: Location{StartLine: getLineNumber(makeslice.Pos()), EndLine: getLineNumber(makeslice.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeSliceInstruction 分析切片操作指令
//func (a *SSAAnalyzer) analyzeSliceInstruction(slice *ssa.Slice, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Slice"
//	operands := []string{a.generateOperandID(slice.X)}
//
//	if slice.Low != nil {
//		operands = append(operands, a.generateOperandID(slice.Low))
//	}
//	if slice.High != nil {
//		operands = append(operands, a.generateOperandID(slice.High))
//	}
//	if slice.Max != nil {
//		operands = append(operands, a.generateOperandID(slice.Max))
//	}
//
//	instruction.Operands = operands
//
//	// 记录变量使用
//	varUsage := &VarUsage{
//		VarID:   a.generateVarIDFromValue(slice.X),
//		VarName: a.getVarName(slice.X),
//		VarType: slice.X.Type().String(),
//		Location: Location{
//			StartLine: getLineNumber(slice.Pos()),
//			EndLine:   getLineNumber(slice.Pos()),
//		},
//		Usage:   VarRead,
//		InstrID: instruction.ID,
//	}
//	funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//
//	// 记录结果变量定义
//	if slice.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(slice),
//			VarName:  a.getVarName(slice),
//			VarType:  slice.Type().String(),
//			Location: Location{StartLine: getLineNumber(slice.Pos()), EndLine: getLineNumber(slice.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeIndexAddrInstruction 分析索引地址指令
//func (a *SSAAnalyzer) analyzeIndexAddrInstruction(indexAddr *ssa.IndexAddr, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "IndexAddr"
//	instruction.Operands = []string{
//		a.generateOperandID(indexAddr.X),
//		a.generateOperandID(indexAddr.Index),
//	}
//
//	// 分析数组/切片元素类型
//	if typeInfo := a.extractTypeInfo(indexAddr.X.Type()); typeInfo != nil {
//		structUsage := &StructUsage{
//			StructID:   typeInfo.ID,
//			StructName: typeInfo.Name,
//			Package:    typeInfo.Package,
//			UsageType:  StructFieldAccess,
//		}
//		funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//	}
//
//	// 记录变量使用
//	for _, operand := range []ssa.Value{indexAddr.X, indexAddr.Index} {
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(operand),
//			VarName: a.getVarName(operand),
//			VarType: operand.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(indexAddr.Pos()),
//				EndLine:   getLineNumber(indexAddr.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	// 记录结果变量定义
//	if indexAddr.Referrers() != nil {
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(indexAddr),
//			VarName:  a.getVarName(indexAddr),
//			VarType:  indexAddr.Type().String(),
//			Location: Location{StartLine: getLineNumber(indexAddr.Pos()), EndLine: getLineNumber(indexAddr.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//}
//
//// analyzeGenericInstruction 分析通用指令
//func (a *SSAAnalyzer) analyzeGenericInstruction(instr ssa.Instruction, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	// 获取指令类型名称
//	instruction.OpCode = strings.TrimPrefix(fmt.Sprintf("%T", instr), "*ssa.")
//
//	// 尝试提取操作数
//	if value, ok := instr.(ssa.Value); ok {
//		// 对于有值的指令，记录变量定义
//		varDef := &VarDef{
//			VarID:    a.generateVarIDFromValue(value),
//			VarName:  a.getVarName(value),
//			VarType:  value.Type().String(),
//			Location: Location{StartLine: getLineNumber(instr.Pos()), EndLine: getLineNumber(instr.Pos())},
//			InstrID:  instruction.ID,
//		}
//		funcInfo.DefVars = append(funcInfo.DefVars, varDef)
//		instruction.Result = varDef.VarID
//	}
//
//	// 记录这是一个未特殊处理的指令类型
//	instruction.Type = InstrUnOp // 使用一个默认类型
//}
//
//// buildCallRelationships 构建调用关系（修正版本）
//func (a *SSAAnalyzer) buildCallRelationships() {
//	for _, node := range a.callGraph.Nodes {
//		if node.Func == nil {
//			continue
//		}
//
//		callerID := a.getFunctionID(node.Func)
//		callerInfo, exists := a.funcMap[callerID]
//		if !exists {
//			continue
//		}
//
//		// 处理出边（调用其他函数）
//		for _, edge := range node.Out {
//			if edge.Callee.Func == nil {
//				continue
//			}
//
//			calleeID := a.getFunctionID(edge.Callee.Func)
//
//			// 创建调用边
//			callEdge := &CallEdge{
//				ID:       a.generateCallEdgeID(edge),
//				CallerID: callerID,
//				CalleeID: calleeID,
//				Type:     a.determineCallType(edge),
//			}
//			a.result.CallGraph.Edges = append(a.result.CallGraph.Edges, callEdge)
//
//			// 更新调用者函数的被调用函数列表
//			if !a.containsString(callerInfo.CalledFuncs, calleeID) {
//				callerInfo.CalledFuncs = append(callerInfo.CalledFuncs, calleeID)
//			}
//
//			// 更新被调用函数的调用者信息
//			if calleeInfo, exists := a.funcMap[calleeID]; exists {
//				if !a.containsString(calleeInfo.Callers, callerID) {
//					calleeInfo.Callers = append(calleeInfo.Callers, callerID)
//				}
//			}
//		}
//
//		// 处理入边（被其他函数调用）
//		for _, edge := range node.In {
//			if edge.Caller.Func == nil {
//				continue
//			}
//
//			callerID := a.getFunctionID(edge.Caller.Func)
//
//			// 更新被调用函数的调用者信息（再次确认，确保完整性）
//			if !a.containsString(callerInfo.Callers, callerID) {
//				callerInfo.Callers = append(callerInfo.Callers, callerID)
//			}
//		}
//	}
//
//	// 设置根函数
//	a.setRootFunctions()
//}
//
//// setRootFunctions 设置根函数（如main函数）
//func (a *SSAAnalyzer) setRootFunctions() {
//	var rootFuncs []string
//
//	for _, mainPkg := range a.mainPkgs {
//		// 查找main函数
//		if mainFunc := mainPkg.Func("main"); mainFunc != nil {
//			mainFuncID := a.getFunctionID(mainFunc)
//			rootFuncs = append(rootFuncs, mainFuncID)
//
//			// 标记main函数为根函数
//			if mainFuncInfo, exists := a.funcMap[mainFuncID]; exists {
//				// 可以在这里添加根函数的特殊标记
//				mainFuncInfo.Name = "main" // 确保名称正确
//			}
//		}
//
//		// 查找init函数
//		if initFunc := mainPkg.Func("init"); initFunc != nil {
//			initFuncID := a.getFunctionID(initFunc)
//			rootFuncs = append(rootFuncs, initFuncID)
//		}
//	}
//
//	// 查找没有调用者的函数作为根函数
//	for funcID, funcInfo := range a.funcMap {
//		if len(funcInfo.Callers) == 0 && !a.containsString(rootFuncs, funcID) {
//			// 检查是否是包初始化函数或其他特殊函数
//			if strings.HasSuffix(funcInfo.Name, "init") ||
//				funcInfo.Package == "runtime" ||
//				strings.Contains(funcInfo.Name, "init.") {
//				rootFuncs = append(rootFuncs, funcID)
//			}
//		}
//	}
//
//	a.result.CallGraph.RootFunctions = rootFuncs
//}
//
//// containsString 检查字符串切片是否包含指定字符串
//func (a *SSAAnalyzer) containsString(slice []string, str string) bool {
//	for _, s := range slice {
//		if s == str {
//			return true
//		}
//	}
//	return false
//}
//
//// determineCallType 确定调用类型（增强版本）
//func (a *SSAAnalyzer) determineCallType(edge *callgraph.Edge) CallType {
//	if edge.Site == nil {
//		return CallStatic
//	}
//
//	switch site := edge.Site.(type) {
//	case *ssa.Call:
//		// 检查是否是接口调用
//		if site.Call.Value != nil {
//			if _, ok := site.Call.Value.Type().Underlying().(*types.Interface); ok {
//				return CallInterface
//			}
//
//			// 检查是否是函数值调用
//			if _, ok := site.Call.Value.Type().(*types.Signature); ok && site.Call.StaticCallee() == nil {
//				return CallDynamic
//			}
//		}
//
//		// 检查是否是闭包调用
//		if site.Call.StaticCallee() != nil &&
//			strings.Contains(site.Call.StaticCallee().Name(), "func") {
//			return CallClosure
//		}
//
//		return CallStatic
//
//	case *ssa.Go:
//		// Go语句调用
//		return CallDynamic
//
//	case *ssa.Defer:
//		// Defer语句调用
//		return CallDynamic
//
//	default:
//		return CallStatic
//	}
//}
//
//// analyzeCallInstruction 分析调用指令（增强版本）
//func (a *SSAAnalyzer) analyzeCallInstruction(call *ssa.Call, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Call"
//	instruction.Type = InstrCall
//
//	// 收集操作数
//	for _, arg := range call.Call.Args {
//		instruction.Operands = append(instruction.Operands, a.generateOperandID(arg))
//
//		// 记录变量使用（参数）
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(arg),
//			VarName: a.getVarName(arg),
//			VarType: arg.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(call.Pos()),
//				EndLine:   getLineNumber(call.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	if call.Call.StaticCallee() != nil {
//		calleeID := a.getFunctionID(call.Call.StaticCallee())
//
//		// 创建调用点信息
//		callSite := &CallSite{
//			ID:        a.generateCallSiteID(call),
//			Location:  Location{StartLine: getLineNumber(call.Pos()), EndLine: getLineNumber(call.Pos())},
//			CallExpr:  call.Call.String(),
//			CalleeID:  calleeID,
//			InstrType: InstrCall,
//		}
//		funcInfo.CallSites = append(funcInfo.CallSites, callSite)
//
//		// 记录被调用函数
//		if !a.containsString(funcInfo.CalledFuncs, calleeID) {
//			funcInfo.CalledFuncs = append(funcInfo.CalledFuncs, calleeID)
//		}
//	} else {
//		// 动态调用，无法确定具体被调用函数
//		callSite := &CallSite{
//			ID:        a.generateCallSiteID(call),
//			Location:  Location{StartLine: getLineNumber(call.Pos()), EndLine: getLineNumber(call.Pos())},
//			CallExpr:  call.Call.String(),
//			CalleeID:  "dynamic", // 标记为动态调用
//			InstrType: InstrCall,
//		}
//		funcInfo.CallSites = append(funcInfo.CallSites, callSite)
//	}
//
//	// 分析调用中的结构体使用
//	a.analyzeCallStructUsage(call, funcInfo)
//}
//
//// analyzeGoInstruction 分析Go语句调用
//func (a *SSAAnalyzer) analyzeGoInstruction(goInstr *ssa.Go, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Go"
//	instruction.Type = InstrCall
//
//	// 收集操作数
//	for _, arg := range goInstr.Call.Args {
//		instruction.Operands = append(instruction.Operands, a.generateOperandID(arg))
//
//		// 记录变量使用
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(arg),
//			VarName: a.getVarName(arg),
//			VarType: arg.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(goInstr.Pos()),
//				EndLine:   getLineNumber(goInstr.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	if goInstr.Call.StaticCallee() != nil {
//		calleeID := a.getFunctionID(goInstr.Call.StaticCallee())
//
//		// 创建调用点信息
//		callSite := &CallSite{
//			ID:        a.generateCallSiteID(goInstr),
//			Location:  Location{StartLine: getLineNumber(goInstr.Pos()), EndLine: getLineNumber(goInstr.Pos())},
//			CallExpr:  goInstr.Call.String(),
//			CalleeID:  calleeID,
//			InstrType: InstrCall,
//		}
//		funcInfo.CallSites = append(funcInfo.CallSites, callSite)
//
//		// 记录被调用函数
//		if !a.containsString(funcInfo.CalledFuncs, calleeID) {
//			funcInfo.CalledFuncs = append(funcInfo.CalledFuncs, calleeID)
//		}
//	}
//
//	// 分析调用中的结构体使用
//	a.analyzeCallStructUsage(goInstr.Call.Common(), funcInfo)
//}
//
//// analyzeDeferInstruction 分析Defer语句调用
//func (a *SSAAnalyzer) analyzeDeferInstruction(deferInstr *ssa.Defer, funcInfo *FunctionSSAInfo, instruction *SSAInstruction) {
//	instruction.OpCode = "Defer"
//	instruction.Type = InstrCall
//
//	// 收集操作数
//	for _, arg := range deferInstr.Call.Args {
//		instruction.Operands = append(instruction.Operands, a.generateOperandID(arg))
//
//		// 记录变量使用
//		varUsage := &VarUsage{
//			VarID:   a.generateVarIDFromValue(arg),
//			VarName: a.getVarName(arg),
//			VarType: arg.Type().String(),
//			Location: Location{
//				StartLine: getLineNumber(deferInstr.Pos()),
//				EndLine:   getLineNumber(deferInstr.Pos()),
//			},
//			Usage:   VarRead,
//			InstrID: instruction.ID,
//		}
//		funcInfo.UsedVars = append(funcInfo.UsedVars, varUsage)
//	}
//
//	if deferInstr.Call.StaticCallee() != nil {
//		calleeID := a.getFunctionID(deferInstr.Call.StaticCallee())
//
//		// 创建调用点信息
//		callSite := &CallSite{
//			ID:        a.generateCallSiteID(deferInstr),
//			Location:  Location{StartLine: getLineNumber(deferInstr.Pos()), EndLine: getLineNumber(deferInstr.Pos())},
//			CallExpr:  deferInstr.Call.String(),
//			CalleeID:  calleeID,
//			InstrType: InstrCall,
//		}
//		funcInfo.CallSites = append(funcInfo.CallSites, callSite)
//
//		// 记录被调用函数
//		if !a.containsString(funcInfo.CalledFuncs, calleeID) {
//			funcInfo.CalledFuncs = append(funcInfo.CalledFuncs, calleeID)
//		}
//	}
//
//	// 分析调用中的结构体使用
//	a.analyzeCallStructUsage(deferInstr.Call.Common(), funcInfo)
//}
//
//// analyzeCallStructUsage 分析调用中的结构体使用（增强版本）
//func (a *SSAAnalyzer) analyzeCallStructUsage(call *ssa.CallCommon, funcInfo *FunctionSSAInfo) {
//	for _, arg := range call.Args {
//		if typeInfo := a.extractTypeInfo(arg.Type()); typeInfo != nil {
//			// 检查是否已经记录过这个结构体使用
//			exists := false
//			for _, usage := range funcInfo.UsedStructs {
//				if usage.StructID == typeInfo.ID {
//					exists = true
//					break
//				}
//			}
//
//			if !exists {
//				structUsage := &StructUsage{
//					StructID:   typeInfo.ID,
//					StructName: typeInfo.Name,
//					Package:    typeInfo.Package,
//					UsageType:  StructMethodCall,
//				}
//				funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//			}
//		}
//	}
//
//	// 分析接收器（如果是方法调用）
//	if call.Value != nil {
//		if typeInfo := a.extractTypeInfo(call.Value.Type()); typeInfo != nil {
//			exists := false
//			for _, usage := range funcInfo.UsedStructs {
//				if usage.StructID == typeInfo.ID {
//					exists = true
//					break
//				}
//			}
//
//			if !exists {
//				structUsage := &StructUsage{
//					StructID:   typeInfo.ID,
//					StructName: typeInfo.Name,
//					Package:    typeInfo.Package,
//					UsageType:  StructMethodCall,
//				}
//				funcInfo.UsedStructs = append(funcInfo.UsedStructs, structUsage)
//			}
//		}
//	}
//}
