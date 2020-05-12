package itf2rpc

import (
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
)

//FWriteRPC write RPC.
type FWriteRPC func(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error

//TargetMap from target to func.
var TargetMap = map[string]FWriteRPC{
	"go_s": GoSvr,
	"go_c": GoClient,
}
