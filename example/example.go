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

type Aaaaa struct {
	En
	EnB
	A string
}

type Bbbbb struct {
	En
	A string
}

func main() {
	automapper.CreateMapper(reflect.TypeOf((*Aaaaa)(nil)), reflect.TypeOf((*Bbbbb)(nil)))

	a := Aaaaa{En{"xxxxxxxxx", "VVVVVVVVVVVVV"}, EnB{}, "XXXXXX"}
	result := automapper.MustMapper(a, reflect.TypeOf((*Bbbbb)(nil)))
	fmt.Println(reflect.TypeOf(result).String())

	result2 := automapper.MustMapper(a, reflect.TypeOf(Bbbbb{}))
	fmt.Println(reflect.TypeOf(result2).String())
}
