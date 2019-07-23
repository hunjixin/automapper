
package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
	"time"
)
type A struct {
	Name       string
}

type B struct {
	Name       string
}

type Son struct {
	PersonModel
}

type PersonModel struct {
	Name       string
	Birth      time.Time
	Address    string
	Sons       *Son
	CreateDate time.Time
	DeleteDate time.Time
	IsDel      bool
	XX A
}

type PersonDto struct {
	Name       string
	Age        int
	Address    string
	Sons       *PersonDto
	CreateDate time.Time
	DeleteDate time.Time
	IsDel      bool
	XX B
}

func init() {
	//automapper.MustCreateMapper(reflect.TypeOf((*PersonModel)(nil)), reflect.TypeOf((*PersonDto)(nil)))
	//automapper.MustCreateMapper(reflect.TypeOf((*Son)(nil)), reflect.TypeOf((*PersonDto)(nil)))
}

func main() {
	children := &PersonModel{}
	children.Name = "bruth"
	children.Birth = time.Date(1993, 3, 4, 1, 2, 3, 4, time.UTC)
	children.Address = "S·H"
	children.CreateDate = time.Now()
	children.IsDel = true
	children.XX = A{ "children"}

	father := &PersonModel{}
	father.Name = "Jimmy"
	children.Birth = time.Date(1973, 3, 4, 1, 2, 3, 4, time.UTC)
	father.Address = "S·H"
	father.CreateDate = time.Now()
	father.IsDel = true
	father.XX = A{"father"}
	father.Sons = &Son{*children}
	result := automapper.MustMapper(father, reflect.TypeOf((*PersonDto)(nil)))
	fmt.Println(result)
}
