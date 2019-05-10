package smn_data_file

import (
	"errors"
	"github.com/json-iterator/go"
	"strings"
)

type DataType int

const (
	_ DataType = iota
	DATA_TYPE_UNKNOW
	DATA_TYPE_JSON
	DATA_TYPE_XML
)

func dataType(data string) DataType {
	data = strings.TrimSpace(data)
	fc, ec := data[0], data[len(data)-1]
	switch fc {
	case '{':
		if ec == '}' {
			return DATA_TYPE_JSON
		}
	case '<':
		if ec == '{' {
			return DATA_TYPE_XML
		}
	}
	return DATA_TYPE_UNKNOW
}

func JsonToMap(data string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	err := jsoniter.Unmarshal([]byte(data), &res)
	return res, err
}

func GetDataMapFromStr(data string) (map[string]interface{}, error) {
	switch dataType(data) {
	case DATA_TYPE_JSON:
		return JsonToMap(data)
	default:
		return nil, errors.New(ERR_UNKNOW_TYPE)
	}
}
