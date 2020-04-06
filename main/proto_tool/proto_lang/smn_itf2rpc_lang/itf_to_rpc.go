package main

import "flag"

func main() {
	i := flag.String("i", "./src/rpc_itf/", "go rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc interface's net accepter, from proto.Message call interface.")
	s := flag.Bool("s", true, "is product server code")
	c := flag.Bool("c", true, "is product client code")
	flag.Parse()
}
