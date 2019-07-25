// +build main7

package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
	"strconv"
	"time"
)

func main(){
	automapper.MustCreateMapper(reflect.TypeOf(time.Time{}), reflect.TypeOf("")).
		Mapping(func(destVal reflect.Value, sourceVal interface{}) {
			str := sourceVal.(time.Time).String()
			destVal.Elem().SetString(str)
		})
	automapper.MustCreateMapper(reflect.TypeOf(0), reflect.TypeOf("")).
		Mapping(func(destVal reflect.Value, sourceVal interface{}) {
			intVal := sourceVal.(int)
			destVal.Elem().SetString(strconv.Itoa(intVal))
		})
	type A struct {
		M time.Time
		N int
	}

	type B struct {
		M string
		N string
	}

	str := automapper.MustMapper(A{time.Now(),123}, reflect.TypeOf(B{}))
	fmt.Println(str)

	mapping := automapper.EnsureMapping(reflect.TypeOf(A{}), reflect.TypeOf(B{}))
	mapping.Mapping(func(destVal reflect.Value, sourceVal interface{}) {
		str := sourceVal.(A).M.String()
		destVal.Interface().(*B).M = "北京时间："+str
		destVal.Interface().(*B).N = "到达次数："+ strconv.Itoa(sourceVal.(A).N)
	})
	str = automapper.MustMapper(A{time.Now(),456}, reflect.TypeOf(B{}))
	fmt.Println(str)
}
