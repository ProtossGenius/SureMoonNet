package lex_analysis

import "github.com/ProtossGenius/pglang/snreader"

type LexType int

const (
	LEX_TYPE_STRING LexType = iota
)

type LexerInput struct {
	snreader.InputItf
	C rune
}

type LexicalUnit struct {
	snreader.ProductItf
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
