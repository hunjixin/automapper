// +build main3
package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
	"time"
)

type UserDto struct {
	Name string
	Addr string
	Age  int
}

type User struct {
	Name string
	Nick string
	Addr string
	Birth time.Time
}

func init() {
	automapper.MustCreateMapper(reflect.TypeOf((*User)(nil)), reflect.TypeOf((*UserDto)(nil))).
	Mapping(func(destVal interface{}, sourceVal interface{}) {
		destVal.(*UserDto).Name = sourceVal.(User).Name + "|"+ sourceVal.(User).Nick
	}).
	Mapping(func(destVal interface{}, sourceVal interface{}) {
		destVal.(*UserDto).Age = time.Now().Year() - sourceVal.(User).Birth.Year()
	})
}

func main() {
	user := &User{"NAME", "NICK", "B·J", time.Date(1992, 10,3,1,0,0,0,time.UTC)}
	result := automapper.MustMapper(user, reflect.TypeOf((*UserDto)(nil)))
	fmt.Println(reflect.TypeOf(result).String())
}
