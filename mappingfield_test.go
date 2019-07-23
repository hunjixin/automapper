package automapper

import (
	"reflect"
	"testing"
)

func TestMapToMapMappingField(t *testing.T) {
	fromMap := map[string]string{
		"xx": "ggg",
	}
	mappingField := &MapToMapMappingField{
		reflect.TypeOf(fromMap),
		reflect.TypeOf(fromMap),
	}
	val := reflect.New(reflect.TypeOf(fromMap)).Elem()
	err := mappingField.Convert(reflect.ValueOf(fromMap), val)
	if err != nil {
		t.Error(err)
	}
	mapVal := val.Interface().(map[string]string)
	if mapVal["xx"] != fromMap["xx"] {
		t.Errorf("Expected value  %s is not equal to the original value %s", fromMap["xx"], mapVal["xx"])
	}
}

func TestStructToMapMappingField(t *testing.T) {
	type B struct {
		Aa string
	}
	fromFields := deepFields(reflect.TypeOf(B{}))
	mappingField := &StructToMapMappingField{
		fromFields,
	}

	sourceValue := reflect.ValueOf(B{"hello bob!"})
	destValue := reflect.ValueOf(map[string]interface{}{})
	mappingField.Convert(sourceValue, destValue)
	valMap := destValue.Interface().(map[string]interface{})
	if valMap["Aa"] != "hello bob!" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestMapToStructMappingField(t *testing.T) {
	type B struct {
		Aa string
	}
	type A struct {
		M  string
		N  int
		BB *B
	}

	toFields := deepFields(reflect.TypeOf(A{}))
	mappingField := &MapToStructMappingField{
		toFields,
	}

	destValue := reflect.New(reflect.TypeOf(A{})).Elem()
	sourceValue := reflect.ValueOf(map[string]interface{}{
		"M":  "XXXXXX",
		"N":  123,
		"BB": B{"XX"},
	})
	mappingField.Convert(sourceValue, destValue)
	convertValue := destValue.Interface().(A)
	if convertValue.M != "XXXXXX" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestArrayToArrayMappingField(t *testing.T) {
	type A struct {
		B string
	}
	arr1 := [5]A{}
	arr1[0].B = "xxxx"
	arr1[4].B = "xxxx"
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(A{}), reflect.TypeOf(A{}))
	arrMap := &Array2ArrayMappingField{
		reflect.TypeOf(A{}),
		reflect.TypeOf(A{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf(arr1)).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]A)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSliceToArrayMappingField(t *testing.T) {
	type A struct {
		B string
	}
	arr1 := make([]*A, 5)
	arr1[0] = &A{"xxxxxx"}
	arr1[4] = &A{"xxxxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(A{}), reflect.TypeOf(A{}))
	arrMap := &Slice2ArrayMappingField{
		reflect.TypeOf(&A{}),
		reflect.TypeOf(A{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([5]A{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]A)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSliceToSliceMappingField(t *testing.T) {
	type A struct {
		B string
	}
	arr1 := make([]*A, 5)
	arr1[0] = &A{"xxxxxx"}
	arr1[4] = &A{"xxxxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(&A{}), reflect.TypeOf(A{}))
	arrMap := &Slice2SliceMappingField{
		reflect.TypeOf(&A{}),
		reflect.TypeOf(A{}),
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]A{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]A)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestArrayToSliceMappingField(t *testing.T) {
	type A struct {
		B string
	}
	arr1 := [5]*A{}
	arr1[0] = &A{"xxxx"}
	arr1[4] = &A{"xxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(A{}), reflect.TypeOf(A{}))
	arrMap := &Array2SliceMappingField{
		reflect.TypeOf(&A{}),
		reflect.TypeOf(A{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]A{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]A)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}
