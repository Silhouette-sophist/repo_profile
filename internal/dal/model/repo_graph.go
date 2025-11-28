package model

type RepoAstCallGraph struct {
	GitRepo        string
	Branch         string
	CommitHash     string
	FileFuncMap    map[string]*CodeFuncInfo
	FileVarMap     map[string][]*FilePkgVar
	FileStructMap  map[string][]*StructInfo
	PathVersionMap map[string]string
	ModuleInfo     *string
}

type StructInfo struct {
	Pkg        string
	File       string
	Name       string
	Content    string
	BlockSpan  *BlockSpan
	DepsStruct map[string]map[string][]string
	UniqueId   string
}

type CodeFuncInfo struct {
	OffsetKey             string
	Content               string
	FilePath              string
	FuncName              string
	BlockSpan             *BlockSpan
	PkgPath               string
	RecvType              *ParamVar
	Type                  string
	Params                []*ParamVar
	Results               []*ParamVar
	CalleeInfos           []*CalleeInfo
	RelatedRepoPkgStructs map[string]map[string][]string
	RelatedRepoPkgVars    map[string]map[string][]string
	AstKey                string
	UniqueId              string
}

type BlockSpan struct {
	StartLine int32
	EndLine   int32
}

type ParamVar struct {
	Type     string
	Name     string
	BaseType string
}

type CalleeInfo struct {
	PkgPath   string
	Name      string
	Receiver  *string
	AstKey    string
	OffsetKey *string
	UniqueId  string
}

type FilePkgVar struct {
	Pkg       string
	File      string
	Type      string
	Name      string
	BlockSpan *BlockSpan
	Content   string
	UniqueId  string
}
