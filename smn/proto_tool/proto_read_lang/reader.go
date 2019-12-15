package proto_read_lang

type MsgReader func(protoPath, pkgHead, goPath, output string) error
