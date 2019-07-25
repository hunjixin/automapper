// +build main2

package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
)

type En struct {
	B string
	D string
}

type EnB struct {
	B string
	D string
}

type ExampleStructA struct {
	EnB
	En
	A string
}

type ExampleStructB struct {
	En
	A string
}

func main() {
	a := ExampleStructA{EnB{}, En{"Sh", "Bj"}, "XXXXXX"}
	result2 := automapper.MustMapper(a, reflect.TypeOf(ExampleStructB{}))
	fmt.Println(result2)
}
