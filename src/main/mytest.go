package main

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

type TypeA interface {
}

type TypeB struct {
	TypeA
}

func main() {
	b := TypeB{}
	ch := make(chan TypeA)
	ch <- b
}
