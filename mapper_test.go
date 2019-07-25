package automapper

import (
	"reflect"
	"testing"
)

func TestOneToOneCreateMapper(t *testing.T) {
	type TestAA struct {
		A string
		B string
	}
	type TestBB struct {
		B string
		A string
	}
	sourceType := reflect.TypeOf((*TestAA)(nil))
	destType := reflect.TypeOf((*TestBB)(nil))
	_, err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, _ := ensureMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 2 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds))
	}
}

func TestOneToManyCreateMapper(t *testing.T) {
	type Embed struct {
		A string
		B string
	}
	type TestAAA struct {
		Embed
	}
	type TestBBB struct {
		Embed
		B string
		A string
	}
	sourceType := reflect.TypeOf((*TestAAA)(nil))
	destType := reflect.TypeOf((*TestBBB)(nil))
	_, err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, _ := ensureMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 2 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds))
	}
	for _, mapField := range mapping.MapFileds {
		struct2strcutMapField := mapField.(*StructFieldMapping)
		if struct2strcutMapField.FromField.Name() == "A" {
			if struct2strcutMapField.ToField.Path != "[Embed][A]" {
				t.Errorf("Map field path error  %s but got %s", ".Embed.A", mapping.Key)
			}
		}
		if struct2strcutMapField.FromField.Name() == "B" {
			if struct2strcutMapField.ToField.Path != "[Embed][B]" {
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
	_, err := CreateMapper(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	mapping, _ := ensureMapping(sourceType, destType)
	if err != nil {
		t.Error(err)
	}
	if len(mapping.MapFileds) != 4 {
		t.Errorf("Inconsistent number of mapped fields expect %d but got %d", 2, len(mapping.MapFileds))
	}
}

func TestMap2Map(t *testing.T) {
	type Map2MapA struct {
		 B string
	}
	type Map2MapB struct {
		B string
	}
	testData := []string {"A","VV","FR"}
	map1 := map[string]Map2MapA{}
	for _, val := range testData {
		map1[val] = Map2MapA{val}
	}
	map2Interface, err := Mapper(map1, reflect.TypeOf(map[string]Map2MapB{}))
	if err != nil {
		t.Error(err)
	}
	map2 := map2Interface.(map[string]Map2MapB)
	for _, val := range testData {
		newVal, ok := map2[val]
		if !ok {
			t.Errorf("value disappear")
		}
		if newVal.B != val {
			t.Errorf("value got but not correct")
		}
	}
}

func TestMap2Struct(t *testing.T) {
	type Map2StructB struct {
		B string
	}

	type Map2StructReceive struct {
		RecevieField *Map2StructB
	}
	map1 := map[string]interface{}{"RecevieField":Map2StructB{"xxxxx"}}
	structInterface, err := Mapper(map1, reflect.TypeOf(Map2StructReceive{}))
	if err != nil {
		t.Error(err)
	}
	map2 := structInterface.(Map2StructReceive)
	if map2.RecevieField.B != "xxxxx" {
		t.Errorf("value got but not correct")
	}
}

func TestStruct2Map(t *testing.T) {
	type Map2StructB struct {
		B string
	}

	structInterface, err := Mapper(Map2StructB{"xxxxx"}, reflect.TypeOf(map[string]interface{}{}))
	if err != nil {
		t.Error(err)
	}
	map2 := structInterface.(map[string]interface{})
	if map2["B"] != "xxxxx" {
		t.Errorf("value got but not correct")
	}
}