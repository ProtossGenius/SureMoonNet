package smn_pglang

type VarDef struct {
	Type    string `json:"type"`
	Var     string `json:"var"`
	ArrSize int    `json:"is_arr"`
}

type Function struct {
	Name   string    `json:"name"`
	Params []*VarDef `json:"params"`
	Return string    `json:"return"`
}

type Interface struct {
	Name      string      `json:"name"`
	Package   string      `json:"package"`
	Functions []*Function `json:"functions"`
}

type Class struct {
	Father     string `json:"father"`
	Interfaces string `json:"interfaces"`
	VarDef     string `json:"var_def"`
}

type FuncMapping struct {
	SystemName string `json:"system_name"`
	FuncName   string `json:"func_name"`
	FuncRename string `json:"func_rename"`
}

type System struct {
	Interface
	SonSystemNames []*VarDef      `json:"son_systems"`
	SonModuleNames []*VarDef      `json:"modules"`
	FuncMappings   []*FuncMapping `json:"func_mappings"`
}
