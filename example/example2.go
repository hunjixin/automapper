package main

import (
	"fmt"
	"github.com/hunjixin/automapper"
	"reflect"
	"time"
)

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
}

type PersonDto struct {
	Name       string
	Age        int
	Address    string
	Sons       *PersonDto
	CreateDate time.Time
	DeleteDate time.Time
	IsDel      bool
}

func init() {
	automapper.CreateMapper(reflect.TypeOf((*PersonModel)(nil)), reflect.TypeOf((*PersonDto)(nil)))
	automapper.CreateMapper(reflect.TypeOf((*Son)(nil)), reflect.TypeOf((*PersonDto)(nil)))
}

type A struct {
}

func main() {
	children := &PersonModel{}
	children.Name = "bruth"
	children.Birth = time.Date(1993, 3, 4, 1, 2, 3, 4, time.UTC)
	children.Address = "S·H"
	children.CreateDate = time.Now()
	children.IsDel = true

	father := &PersonModel{}
	father.Name = "Jimmy"
	children.Birth = time.Date(1973, 3, 4, 1, 2, 3, 4, time.UTC)
	father.Address = "S·H"
	father.CreateDate = time.Now()
	father.IsDel = true
	father.Sons = &Son{*children}
	result := automapper.MustMapper(father, reflect.TypeOf((*PersonDto)(nil)))
	fmt.Println(reflect.TypeOf(result).String())
}
