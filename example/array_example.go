// +build main1

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
	Name  string
	Nick  string
	Addr  string
	Birth time.Time
}

func main() {
	automapper.MustMapper(User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}, reflect.TypeOf(UserDto{}))
	users := [3]*User{}
	users[0] = &User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}
	users[2] = &User{"Jack", "neo", "W·S", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}
	result2 := automapper.MustMapper(users, reflect.TypeOf([]*UserDto{}))
	fmt.Println(result2)
}
