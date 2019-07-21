// +build main4
package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
)

type UserDto struct {
	Nick string
	Name string
}

type User struct {
	Name string  `mapping:"Nick"`
	Nick string  `mapping:"Name"`
}

func init() {
	automapper.MustCreateMapper(reflect.TypeOf((*User)(nil)), reflect.TypeOf((*UserDto)(nil)))
}

func main() {
	user := &User{"NAME", "NICK"}
	result := automapper.MustMapper(user, reflect.TypeOf((*UserDto)(nil)))
	fmt.Println(reflect.TypeOf(result).String())
}
