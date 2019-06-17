package smn_pglang

type VarDef struct {
	Type    string `json:"type"`
	Var     string `json:"var"`
	ArrSize int    `json:"arr_size"` // not zero means is array; < 0 means don't know its size.
}

type StructDef struct {
	Name      string    `json:"name"`
	Variables []*VarDef `json:"variables"`
}

type FuncDef struct {
	Name    string    `json:"name"`
	Params  []*VarDef `json:"params"`
	Ret     bool      `json:"ret"`
	Returns []*VarDef `json:"returns"`
}

func (this *FuncDef) Parse() {
	if len(this.Returns) != 0 && this.Returns[0].Type != "void" {
		this.Ret = true
	}
}

func NewFuncDef() *FuncDef {
	return &FuncDef{Params: make([]*VarDef, 0), Returns: make([]*VarDef, 0)}
}

type ItfDef struct {
	Name      string     `json:"name"`
	Package   string     `json:"package"`
	Functions []*FuncDef `json:"functions"`
}

func NewItfDefine() *ItfDef {
	return &ItfDef{Functions: make([]*FuncDef, 0)}
}

type ClassDef struct {
	Name       string    `json:"name"`
	Package    string    `json:"package"`
	Father     []string  `json:"father"`
	Interfaces []string  `json:"interfaces"`
	Variables  []*VarDef `json:"var_def"`
}

type FuncMapping struct {
	SystemName string `json:"system_name"`
	FuncName   string `json:"func_name"`
	FuncRename string `json:"func_rename"`
}

type SystemDef struct {
	ItfDef
	SonSystemNames []*VarDef      `json:"son_systems"`
	SonModuleNames []*VarDef      `json:"modules"`
	FuncMappings   []*FuncMapping `json:"func_mappings"`
}
