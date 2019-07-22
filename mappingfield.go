package automapper

import (
	"reflect"
)

type IStructConverter interface {
	Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error
	GetType() int
	GetFromField() *structField
	GetToField() *structField
}

type Array2ArrayMappingField struct {
	BaseMappingField
	FromFieldType  reflect.Type
	ToFieldType    reflect.Type
	Length int
	ChildMapping *MappingInfo
}

func (array2ArrayMappingField *Array2ArrayMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	for i:=0;i<array2ArrayMappingField.Length;i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), array2ArrayMappingField.ToFieldType)
		if err != nil {
			return err
		}
		destFieldValue.Set(reflect.ValueOf(childVal))
	}

	return nil
}

type ChildrenMappingField struct {
	BaseMappingField
	ChildMapping *MappingInfo
}

func (cildrenMappingField *ChildrenMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {

	childVal, err := cildrenMappingField.ChildMapping.mapper(sourceFieldValue.Interface())
	if err != nil {
		return err
	}
	if destFieldValue.Kind() == reflect.Ptr {
		destFieldValue.Set(childVal.Addr())
	}else{
		destFieldValue.Set(childVal)
	}
	return nil
}

type NoneMappingField struct {
	BaseMappingField
}

func (noneMappingField *NoneMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	panic("never come here")
}

type SameTypeMappingField struct {
	BaseMappingField
}

func (sameTypeMappingField *SameTypeMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

type AnyMappingField struct {
	BaseMappingField
}

func (anyMappingField *AnyMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

type BaseMappingField struct {
	Type      int
	FromField *structField
	ToField   *structField
}

func (mappingField *BaseMappingField) GetType() int {
	return mappingField.Type
}

func (mappingField *BaseMappingField) GetFromField() *structField {
	return mappingField.FromField
}

func (mappingField *BaseMappingField) GetToField() *structField {
	return mappingField.ToField
}

func (baseMappingField *BaseMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	panic("never come here")
}