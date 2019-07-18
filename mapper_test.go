package automapper

import (
	"reflect"
	"testing"
)


func TestOneToOneCreateMapper(t *testing.T) {
	type TestA struct {
		A string
		B string
	}
	type TestB struct {
		B string
		A string
	}
	sourceType := reflect.TypeOf((*TestA)(nil))
	destType := reflect.TypeOf((*TestB)(nil))
	err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, err := GetMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 2 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds), )
	}
}

func TestOneToManyCreateMapper(t *testing.T) {
	type Embed struct {
		A string
		B string
	}
	type TestA struct {
		Embed
	}
	type TestB struct {
		Embed
		B string
		A string
	}
	sourceType := reflect.TypeOf((*TestA)(nil))
	destType := reflect.TypeOf((*TestB)(nil))
	err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, err := GetMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 2 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds), )
	}
	for _, mapField := range mapping.MapFileds {
		if mapField.Name == "A" {
			if mapField.Key != ".Embed.A" {
				t.Errorf("Map field path error  %s but got %s", ".Embed.A", mapping.Key)
			}
		}
		if mapField.Name == "B" {
			if mapField.Key != ".Embed.B" {
				t.Errorf("Map field path error  %s but got %s", ".Embed.B", mapping.Key)
			}
		}
	}
}

func TestManyToManyCreateMapper(t *testing.T) {
	type Embed struct {
		A string
		B string
	}
	type TestA struct {
		Embed
		B string
		A string
	}
	type TestB struct {
		Embed
		B string
		A string
	}
	sourceType := reflect.TypeOf((*TestA)(nil))
	destType := reflect.TypeOf((*TestB)(nil))
	err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, err := GetMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 4 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds), )
	}
}