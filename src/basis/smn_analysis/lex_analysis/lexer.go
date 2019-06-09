package lex_analysis

type LexType int

const (
	LEX_TYPE_STRING LexType = iota
)

type LexerInput struct {
	InputItf
	C rune
}

type LexicalUnit struct {
	ProductItf
	Type  LexType
	Value string
}

type LexicalCfg struct {
	StartCharset    string   `json:"start_charset"` //must in this charset.
	StartWith       []string `json:"start_with"`    //if len != 0, only when find those string can do.
	EndWith         []string `json:"end_with"`
	Contains        string   `json:"contains"`
	NotContains     string   `json:"not_contains"`
	FilterBackSlash bool     `json:"filter_back_slash"`
}
