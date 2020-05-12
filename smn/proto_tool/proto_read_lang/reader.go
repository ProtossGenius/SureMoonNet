package proto_read_lang

import (
	"fmt"
)

//MsgReader proto message  read code.
type MsgReader func(protoPath, module string) error

//Readers lang->MsgReader.
var Readers = map[string]MsgReader{
	"go": GoMsgReader,
}

//Write write reader.
func Write(lang string, protoPath, module string) error {
	f, ok := Readers[lang]
	if !ok {
		return fmt.Errorf("Not found in Readers: %s", lang)
	}
	return f(protoPath, module)
}
