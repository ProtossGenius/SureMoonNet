package lex_analysis

import "basis/smn_analysis"

type LexType int

const (
	LEX_TYPE_STRING LexType = iota
)

type LexerInput struct {
	smn_analysis.InputItf
	C rune
}

type LexicalUnit struct {
	smn_analysis.ProductItf
	Type  LexType
	Value string
}

func (this *LexicalUnit) ProductType() int {
	return int(this.Type)
}

type LexicalCfg struct {
	StartCharset    string   `json:"start_charset"` //must in this charset.
	StartWith       []string `json:"start_with"`    //if len != 0, only when find those string can do.
	EndWith         []string `json:"end_with"`
	Contains        string   `json:"contains"`
	NotContains     string   `json:"not_contains"`
	FilterBackSlash bool     `json:"filter_back_slash"`
}
