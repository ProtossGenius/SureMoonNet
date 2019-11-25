package main

import (
	"fmt"
	"reflect"
)

type StructTest struct {
	Hello string `fd:"about desc"`
}

func main() {
	t := &StructTest{Hello: "hello world"}
	tt := reflect.TypeOf(*t)
	vt := reflect.ValueOf(t).Elem()
	//	tvt := vt.Type()
	for i := 0; i < tt.NumField(); i++ {
		ft := tt.Field(i)
		fmt.Printf("name %v tag %v type %v\n", ft.Name, ft.Tag, ft.Type)
		v := vt.Field(i).Interface()
		fmt.Printf("~~~~~~~~~~~~~ %v", v)
		v = "ccccccccccccc"
		fmt.Printf(t.Hello)
	}
}
