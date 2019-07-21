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

type ChildrenMappingField struct {
	BaseMappingField
	ChildMapping *MappingInfo
}

func (mappingField *ChildrenMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	childVal, err := Mapper(sourceFieldValue.Interface(), destFieldValue.Type())
	if err != nil {
		return err
	}
	destFieldValue.Set(reflect.ValueOf(childVal))
	return nil
}

type NoneMappingField struct {
	BaseMappingField
}

func (mappingField *NoneMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	panic("never come here")
}

type SameTypeMappingField struct {
	BaseMappingField
}

func (mappingField *BaseMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
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
