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
	DestValueType reflect.Type
}

// MapToStructMappingField deep child field and convert to map
func (mapToMapMappingField *MapToMapMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	mappingInfo, _ :=ensureMapping(mapToMapMappingField.SoureValueType.Elem(), mapToMapMappingField.DestValueType.Elem())
	sourceMapIter := sourceFieldValue.MapRange()
	destFieldValue.Set(reflect.MakeMap(mapToMapMappingField.DestValueType))
	for ;sourceMapIter.Next(); {
		val, err := mappingInfo.mapper(sourceMapIter.Value().Interface())
		if err != nil {
			return err
		}
		if mapToMapMappingField.DestValueType.Kind() == reflect.Ptr {
			destFieldValue.SetMapIndex(sourceMapIter.Key(), val.Addr())
		}else{
			destFieldValue.SetMapIndex(sourceMapIter.Key(), val)
		}
	}
	return nil
}

// StructToMapMappingField
type MapToStructMappingField struct {

}

// MapToStructMappingField deep child field and convert to map
func (mapToStructMappingField *MapToStructMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	fields := deepFields(indirectType(destFieldValue.Type()))
	values := deepValue(destFieldValue)
	if sourceFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(nil))
	}
    sourceMap := reflect.Indirect(sourceFieldValue)
	mapIter := sourceMap.MapRange()
	for ;mapIter.Next(); {
		for index, valueField := range values {
			if fields[index].Name() == mapIter.Key().Interface().(string) {
				valueField.Set(mapIter.Value().Elem())
			}
		}
	}

	return nil
}


// StructToMapMappingField
type StructToMapMappingField struct {

}

// StructToMapMappingField deep child field and convert to map
func (structToMapMappingField *StructToMapMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	fields := deepFields(indirectType(sourceFieldValue.Type()))
	values := deepValue(sourceFieldValue)
	if destFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(map[string]interface{}{}))
	}

	for _, field := range fields {
		destFieldValue.SetMapIndex(reflect.ValueOf(field.Name()), values[field.FiledIndex])
	}
	return nil
}



type Array2ArrayMappingField struct {

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

type Slice2ArrayMappingField struct {

}

func (slice2ArrayMappingField *Slice2ArrayMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	//TODO
	return nil
}

type Slice2SliceMappingField struct {

}

func (slice2SliceMappingField *Slice2SliceMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	//TODO
	return nil
}

type Array2SliceMappingField struct {

}

func (array2ArrayMappingField *Array2SliceMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	//TODO
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
	if destFieldValue.Kind() == reflect.Ptr {
		destFieldValue.Set(childVal.Addr())
	}else{
		destFieldValue.Set(childVal)
	}
	return nil
}



type NoneMappingField struct {

}

func (noneMappingField *NoneMappingField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	panic("never come here")
}



type SameTypeMappingField struct {

	SourceType reflect.Type
	DestType reflect.Type
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

type StructFieldField struct {
	FromField *structField
	ToField   *structField
	ChildMapping *MappingInfo
}

func (structFieldField *StructFieldField) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	val, err :=structFieldField.ChildMapping.mapper(sourceFieldValue.Interface())
	if err != nil {
		return err
	}
	destFieldValue.Set(val)
	return nil
}

