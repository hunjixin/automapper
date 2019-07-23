package automapper

import (
	"fmt"
	"reflect"
	"testing"
)
type InnerA struct {
	A string
}
type FromA struct {
	InnerA InnerA
}
type ToB struct {
	InnerA map[string]interface{}
}

func TestStructToMapMappingField(t *testing.T) {
	fromFields := deepFields(reflect.TypeOf(FromA{}))
	toFields := deepFields(reflect.TypeOf(ToB{}))
	mappingField := &StructToMapMappingField{
		BaseMappingField{
			Type:      StructToMap,
			FromField: fromFields[0],
			ToField:   toFields[0],
		},
	}

	sourceValue := reflect.ValueOf(InnerA{A:"hello bob!"})
	destValue := reflect.ValueOf(map[string]interface{}{})
	mappingField.Convert(sourceValue, destValue)
	fmt.Println(destValue.Interface().(map[string]interface{})["A"])
	if destValue.Interface().(map[string]interface{})["A"] != "hello bob!" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestMapToStructMappingField(t *testing.T) {
	type B struct {
		Aa string
	}
	type A struct {
		M string
		N int
		BB B
	}

	toFields := deepFields(reflect.TypeOf(A{}))
	mappingField := &MapToStructMappingField{
		BaseMappingField{
			Type:      MapToStruct,
			FromField:  nil,
			ToField:   toFields[0],
		},
	}

	destValue := reflect.New(reflect.TypeOf(A{})).Elem()
	sourceValue  := reflect.ValueOf(map[string]interface{}{
		"M":"XXXXXX",
		"N":123,
		"BB":B{"XX"},
	})
	mappingField.Convert(sourceValue, destValue)
	convertValue := destValue.Interface().(A)
	if convertValue.M != "XXXXXX" {
		t.Error("the converted value is not equal to the previous value")
	}
}