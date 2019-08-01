package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
)

type A struct {
	B string
}

type C struct {
	B string
}

func main() {
	a1 := &A{"Multiple pointers"}
	a2 := &a1
	a3 := &a2

	b1 := C{}
	b2 := &b1
	b3 := &b2
	fmt.Println(b3)
	newVal := automapper.MustMapper(a3, reflect.TypeOf(b3))

	fmt.Println(**(newVal.(**C)))
}
