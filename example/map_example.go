// +build main4

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
	//map => map
	map1 := map[string]*User{
		"Hellen": &User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)},
		"Jack":   &User{"Jack", "neo", "W·S", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)},
	}

	newVal := automapper.MustMapper(map1, reflect.TypeOf(map[string]*UserDto{}))
	for key, val := range newVal.(map[string]*UserDto) {
		fmt.Print(key)
		fmt.Println(val)
	}

	//map => struct
	map2 := map[string]interface{}{
		"Name":  "Hellen",
		"Nick":  "NICK",
		"Addr":  "B·J",
		"Birth": time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC),
	}
	newVal = automapper.MustMapper(map2, reflect.TypeOf(User{}))
	fmt.Println(newVal)

	//struct => map
	newVal = automapper.MustMapper(&User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}, reflect.TypeOf(map[string]interface{}{}))
	fmt.Println(newVal)
}
