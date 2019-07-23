package automapper

import (
	"reflect"
)

type IStructConverter interface {
	Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error
}

// StructToMapMappingField
type MapToMapMappingField struct {
	SoureValueType reflect.Type
	DestValueType  reflect.Type
}

// MapToStructMappingField deep child field and convert to map
func (mapToMapMappingField *MapToMapMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	mappingInfo, _ := ensureMapping(mapToMapMappingField.SoureValueType.Elem(), mapToMapMappingField.DestValueType.Elem())
	sourceMapIter := sourceFieldValue.MapRange()
	destFieldValue.Set(reflect.MakeMap(mapToMapMappingField.DestValueType))
	for sourceMapIter.Next() {
		val, err := mappingInfo.mapper(sourceMapIter.Value().Interface())
		if err != nil {
			return err
		}
		if val.IsValid() {
			if mapToMapMappingField.DestValueType.Kind() == reflect.Ptr {
				if val.CanAddr() {
					destFieldValue.SetMapIndex(sourceMapIter.Key(), val.Addr())
				}else{
					copyVal := reflect.New(val.Type()).Elem()
					copyVal.Set(val)
					destFieldValue.SetMapIndex(sourceMapIter.Key(), copyVal.Addr())
				}
			} else {
				destFieldValue.SetMapIndex(sourceMapIter.Key(), val)
			}
		}
	}
	return nil
}

// StructToMapMappingField
type MapToStructMappingField struct {
	ToFields []*structField
}

// MapToStructMappingField deep child field and convert to map
func (mapToStructMappingField *MapToStructMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	values := deepValue(destFieldValue)
	if sourceFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(nil))
	}
	sourceMap := reflect.Indirect(sourceFieldValue)
	mapIter := sourceMap.MapRange()
	for mapIter.Next() {
		for index, valueField := range values {
			if mapToStructMappingField.ToFields[index].Name() == mapIter.Key().Interface().(string) {
				setValue(valueField, mapIter.Value().Elem())
			}
		}
	}

	return nil
}

func setValue(destValue, sourceValue reflect.Value){
	if !sourceValue.IsValid() {
		return
	}
	if destValue.Kind() == reflect.Ptr {
		if sourceValue.CanAddr() {
			destValue.Set(sourceValue.Addr())
		}else{
			val := reflect.New(sourceValue.Type()).Elem()
			val.Set(sourceValue)
			destValue.Set(val.Addr())
		}
	} else {
		destValue.Set(sourceValue)
	}
}
// StructToMapMappingField
type StructToMapMappingField struct {
	FromFields []*structField
}

// StructToMapMappingField deep child field and convert to map
func (structToMapMappingField *StructToMapMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	values := deepValue(sourceFieldValue)
	if destFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(map[string]interface{}{}))
	}

	for _, field := range structToMapMappingField.FromFields {
		key := reflect.ValueOf(field.Name())
		val := values[field.FiledIndex]
		if val.IsValid() {
			if destFieldValue.Type().Kind() == reflect.Ptr {
				if val.CanAddr() {
					destFieldValue.SetMapIndex(key, val.Addr())
				}else{
					copyVal := reflect.New(val.Type()).Elem()
					copyVal.Set(val)
					destFieldValue.SetMapIndex(key, copyVal.Addr())
				}
			} else {
				destFieldValue.SetMapIndex(key, val)
			}
		}
	}
	return nil
}

type Array2ArrayMappingField struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	Length        int
	ChildMapping  *MappingInfo
}

func (array2ArrayMappingField *Array2ArrayMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	for i := 0; i < array2ArrayMappingField.Length; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), array2ArrayMappingField.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}
	return nil
}

type Slice2ArrayMappingField struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ArrayLen      int
	ChildMapping  *MappingInfo
}

func (slice2ArrayMappingField *Slice2ArrayMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	if destFieldValue.Len() != slice2ArrayMappingField.ArrayLen {
		return ErrLengthNotMatch
	}
	for i := 0; i < slice2ArrayMappingField.ArrayLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), slice2ArrayMappingField.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}

	return nil
}

type Slice2SliceMappingField struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ChildMapping  *MappingInfo
}

func (slice2SliceMappingField *Slice2SliceMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	sourceLen := sourceFieldValue.Len()
	destFieldValue.Set(reflect.MakeSlice(destFieldValue.Type(), sourceLen, sourceLen))
	for i := 0; i < sourceLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), slice2SliceMappingField.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}
	return nil
}

type Array2SliceMappingField struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ArrayLen      int
	ChildMapping  *MappingInfo
}

func (array2ArrayMappingField *Array2SliceMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	if sourceFieldValue.Len() != array2ArrayMappingField.ArrayLen {
		return ErrLengthNotMatch
	}
	newSlices := reflect.MakeSlice(destFieldValue.Type(), array2ArrayMappingField.ArrayLen, array2ArrayMappingField.ArrayLen)
	for i := 0; i < array2ArrayMappingField.ArrayLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), array2ArrayMappingField.ToFieldType)
		if err != nil {
			return err
		}
		setValue(newSlices.Index(i), childVal)
	}
	destFieldValue.Set(newSlices)
	return nil
}

type ChildrenMappingField struct {
	ChildMapping *MappingInfo
}

func (cildrenMappingField *ChildrenMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	childVal, err := cildrenMappingField.ChildMapping.mapper(sourceFieldValue.Interface())
	if err != nil {
		return err
	}
	setValue(destFieldValue, childVal)
	return nil
}

type NoneMappingField struct {
}

func (noneMappingField *NoneMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	panic("never come here")
}

type SameTypeMappingField struct {
	SourceType reflect.Type
	DestType   reflect.Type
}

func (sameTypeMappingField *SameTypeMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

type AnyMappingField struct {
}

func (anyMappingField *AnyMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

type StructFieldMap struct {
	FromField    *structField
	ToField      *structField
	ChildMapping *MappingInfo
}

func (structFieldMap *StructFieldMap) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	val, err := structFieldMap.ChildMapping.mapper(sourceFieldValue.Interface())
	if err != nil {
		return err
	}
	setValue(destFieldValue, val)
	return nil
}
