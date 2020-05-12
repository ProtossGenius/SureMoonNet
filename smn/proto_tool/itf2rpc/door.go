package itf2rpc

import (
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
)

//FWriteRPC write RPC.
type FWriteRPC func(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error

//TargetMap from target to func.
var TargetMap = map[string]FWriteRPC{
	"go_s": GoSvr,
	"go_c": GoClient,
}

//Write itf to rpc.
func Write(target, path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	f, ok := TargetMap[target]
	if !ok {
		return fmt.Errorf("Can't Found Target %s", target)
	}
	return f(path, module, itfFullPkg, itf)
}
