package proto_read_lang

type MsgReader func(protoPath, pkgHead, goPath, ext, output string) error