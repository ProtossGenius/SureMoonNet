package main

import (
	"flag"
	"fmt"

	"com.suremoon.net/smn/analysis/smn_rpc_itf"
)

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc insterface;'s net accepter, from proto.Message call interface.")
	flag.Parse()
	_, err := smn_rpc_itf.GetItfListFromDir(*o)
	fmt.Println(err)
	fmt.Println(*i)
}
