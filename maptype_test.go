package automapper

import (
	"fmt"
	"reflect"
	"testing"
)

type SimpleA struct {
	B string
}

type SimpleB struct {
	B string
}

type ComplexA struct {
	M  string
	N  int
	BB *SimpleB
}

func TestSimple_MapToMapMappingField(t *testing.T) {
	fromMap := map[string]string{
		"xx": "ggg",
	}
	mappingField := &MapToMapMapping{
		reflect.TypeOf(fromMap),
		reflect.TypeOf(fromMap),
	}
	val := reflect.New(mappingField.DestValueType).Elem()
	err := mappingField.Convert(reflect.ValueOf(fromMap), val)
	if err != nil {
		t.Error(err)
	}
	mapVal := val.Interface().(map[string]string)
	if mapVal["xx"] != fromMap["xx"] {
		t.Errorf("Expected value  %s is not equal to the original value %s", fromMap["xx"], mapVal["xx"])
	}
}

func TestStruct_MapToMapMappingField(t *testing.T) {
	fromMap := map[string]SimpleA{
		"xx": SimpleA{"xx"},
	}
	mappingField := &MapToMapMapping{
		reflect.TypeOf(fromMap),
		reflect.TypeOf(map[string]SimpleB{}),
	}
	val := reflect.New(mappingField.DestValueType).Elem()
	err := mappingField.Convert(reflect.ValueOf(fromMap), val)
	if err != nil {
		t.Error(err)
	}
	mapVal := val.Interface().(map[string]SimpleB)
	if mapVal["xx"].B != fromMap["xx"].B {
		t.Errorf("Expected value  %s is not equal to the original value %s", fromMap["xx"], mapVal["xx"])
	}
}

func TestPtr_MapToMapStructMappingField(t *testing.T) {
	fromMap := map[string]SimpleA{
		"xx": SimpleA{"xx"},
	}
	mappingField := &MapToMapMapping{
		reflect.TypeOf(fromMap),
		reflect.TypeOf(map[string]*SimpleB{}),
	}
	val := reflect.New(mappingField.DestValueType).Elem()
	err := mappingField.Convert(reflect.ValueOf(fromMap), val)
	if err != nil {
		t.Error(err)
	}
	mapVal := val.Interface().(map[string]*SimpleB)
	if mapVal["xx"].B != fromMap["xx"].B {
		t.Errorf("Expected value  %s is not equal to the original value %s", fromMap["xx"], mapVal["xx"])
	}
}

func TestSimple_StructToMapMappingField(t *testing.T) {
	type B struct {
		Aa string
	}
	fromFields := deepFields(reflect.TypeOf(B{}))
	mappingField := &StructToMapMapping{
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

func TestStruct_StructToMapStructMappingField(t *testing.T) {
	type AAA struct {
		AA SimpleA
	}
	fromFields := deepFields(reflect.TypeOf(AAA{}))
	mappingField := &StructToMapMapping{
		fromFields,
	}

	sourceValue := reflect.ValueOf(AAA{SimpleA{"hello bob!"}})
	destValue := reflect.ValueOf(map[string]interface{}{})
	mappingField.Convert(sourceValue, destValue)
	valMap := destValue.Interface().(map[string]interface{})
	if valMap["AA"].(SimpleA).B != "hello bob!" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtr_StructToMapMappingField(t *testing.T) {
	type AAA struct {
		AA *SimpleA
	}
	fromFields := deepFields(reflect.TypeOf(AAA{}))
	mappingField := &StructToMapMapping{
		fromFields,
	}

	sourceValue := reflect.ValueOf(AAA{&SimpleA{"hello bob!"}})
	destValue := reflect.ValueOf(map[string]interface{}{})
	mappingField.Convert(sourceValue, destValue)
	valMap := destValue.Interface().(map[string]interface{})
	if valMap["AA"].(*SimpleA).B != "hello bob!" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSimple_MapToStructMappingField(t *testing.T) {
	toFields := deepFields(reflect.TypeOf(ComplexA{}))
	mappingField := &MapToStructMapping{
		toFields,
	}
	type ComplexAPtr struct {
		M  string
		N  int
		BB *SimpleB
	}
	destValue := reflect.New(reflect.TypeOf(ComplexAPtr{})).Elem()
	sourceValue := reflect.ValueOf(map[string]interface{}{
		"M":  "XXXXXX",
		"N":  123,
		"BB": SimpleB{"XX"},
	})
	mappingField.Convert(sourceValue, destValue)
	convertValue := destValue.Interface().(ComplexAPtr)
	if convertValue.M != "XXXXXX" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtr_MapToStructMappingField(t *testing.T) {
	toFields := deepFields(reflect.TypeOf(ComplexA{}))
	mappingField := &MapToStructMapping{
		toFields,
	}
	type ComplexAPtr struct {
		M  string
		N  int
		BB *SimpleB
	}
	destValue := reflect.New(reflect.TypeOf(ComplexAPtr{})).Elem()
	sourceValue := reflect.ValueOf(map[string]interface{}{
		"M":  "XXXXXX",
		"N":  123,
		"BB": SimpleB{"XX"},
	})
	mappingField.Convert(sourceValue, destValue)
	convertValue := destValue.Interface().(ComplexAPtr)
	if convertValue.M != "XXXXXX" {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSimple_ArrayToArrayMappingField(t *testing.T) {
	arr1 := [5]SimpleB{}
	arr1[0].B = "xxxx"
	arr1[4].B = "xxxx"
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(SimpleB{}), reflect.TypeOf(SimpleB{}))
	arrMap := &Array2ArrayMapping{
		reflect.TypeOf(SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf(arr1)).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtr_ArrayToArrayMappingField(t *testing.T) {
	arr1 := [5]*SimpleB{}
	arr1[0] = &SimpleB{"xxxx"}
	arr1[4] = &SimpleB{"xxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(&SimpleB{}), reflect.TypeOf(&SimpleB{}))
	arrMap := &Array2ArrayMapping{
		reflect.TypeOf(&SimpleB{}),
		reflect.TypeOf(&SimpleB{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf(arr1)).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]*SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSimple_SliceToArrayMappingField(t *testing.T) {
	arr1 := make([]SimpleB, 5)
	arr1[0] = SimpleB{"xxxxxx"}
	arr1[4] = SimpleB{"xxxxxx"}

	ptrMapping, _ := ensureMapping(reflect.TypeOf(SimpleB{}), reflect.TypeOf(SimpleB{}))
	arrMap := &Slice2ArrayMapping{
		reflect.TypeOf(SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		5,
		ptrMapping,
	}
	destValue := reflect.New(reflect.TypeOf([5]SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtrSliceToArrayMappingField(t *testing.T) {
	arr1 := make([]*SimpleB, 5)
	arr1[0] = &SimpleB{"xxxxxx"}
	arr1[4] = &SimpleB{"xxxxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(&SimpleB{}), reflect.TypeOf(&SimpleB{}))
	arrMap := &Slice2ArrayMapping{
		reflect.TypeOf(&SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([5]*SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([5]*SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestStruct_SliceToSliceMappingField(t *testing.T) {
	arr1 := make([]SimpleB, 5)
	arr1[0] = SimpleB{"xxxxxx"}
	arr1[4] = SimpleB{"xxxxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(SimpleB{}), reflect.TypeOf(SimpleB{}))
	arrMap := &Slice2SliceMapping{
		reflect.TypeOf(SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtr_SliceToSliceMappingField(t *testing.T) {
	arr1 := make([]*SimpleB, 5)
	arr1[0] = &SimpleB{"xxxxxx"}
	arr1[4] = &SimpleB{"xxxxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(&SimpleB{}), reflect.TypeOf(&SimpleB{}))
	arrMap := &Slice2SliceMapping{
		reflect.TypeOf(&SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]*SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]*SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestSimple_ArrayToSliceMappingField(t *testing.T) {
	arr1 := [5]SimpleB{}
	arr1[0] = SimpleB{"xxxx"}
	arr1[4] = SimpleB{"xxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(SimpleB{}), reflect.TypeOf(SimpleB{}))
	arrMap := &Array2SliceMapping{
		reflect.TypeOf(SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}

func TestPtr_ArrayToSliceMappingField(t *testing.T) {
	arr1 := [5]*SimpleB{}
	arr1[0] = &SimpleB{"xxxx"}
	arr1[4] = &SimpleB{"xxxx"}
	sameTypeMapping, _ := ensureMapping(reflect.TypeOf(&SimpleB{}), reflect.TypeOf(&SimpleB{}))
	arrMap := &Array2SliceMapping{
		reflect.TypeOf(&SimpleB{}),
		reflect.TypeOf(SimpleB{}),
		5,
		sameTypeMapping,
	}
	destValue := reflect.New(reflect.TypeOf([]*SimpleB{})).Elem()
	err := arrMap.Convert(reflect.ValueOf(arr1), destValue)
	if err != nil {
		t.Error(err)
	}
	newArr := destValue.Interface().([]*SimpleB)
	if newArr[0].B != arr1[0].B {
		t.Error("the converted value is not equal to the previous value")
	}
}


func TestPtrToPtrMapping(t *testing.T) {
	type PtrToPtrMappingA struct {
		M string
	}
	type PtrToPtrMappingB struct {
		M string
	}

	mapping := EnsureMapping(&PtrToPtrMappingA{}, &PtrToPtrMappingB{})
	realMapping :=  reflect.TypeOf(mapping.MapFileds[0])
	if realMapping != reflect.TypeOf(&PtrToPtrMapping{}) {
		t.Errorf("expect PtrToPtrMapping but got %s", realMapping.String())
	}

	sourceValue := reflect.ValueOf(&PtrToPtrMappingA{"Data"})
	t1 := reflect.TypeOf(&PtrToPtrMappingB{})
	destValue := reflect.New(t1).Elem()
	err := mapping.MapFileds[0].Convert(sourceValue, destValue)
	if err != nil {
		t.Error(err)
	}
	mapingValue := destValue.Elem().Field(0).Interface().(string)
	if mapingValue != "Data" {
		t.Errorf("expect 'Data' but got %s", mapingValue)
	}
}
func TestPtrToPtrMdddapping(t *testing.T) {
	type A struct {

	}
	type B struct {

	}
	t1 := reflect.TypeOf([]A{})
	t2 := reflect.TypeOf([]B{})
	fmt.Println(t1)
	fmt.Println(t2)
	if t1 == t2 {
		fmt.Println("xx")
	}
}
