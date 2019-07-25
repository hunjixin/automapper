package automapper

import (
	"reflect"
)

// IStructConverter define all conversion types
type IStructConverter interface {
	Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error
}

// MapToMapMapping map to map
type MapToMapMapping struct {
	SoureValueType reflect.Type
	DestValueType  reflect.Type
}

// Convert match the map value by key. the result of recursively mapping the value as a new map
func (mapToMapMapping *MapToMapMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	mappingInfo, _ := ensureMapping(mapToMapMapping.SoureValueType.Elem(), mapToMapMapping.DestValueType.Elem())
	sourceMapIter := sourceFieldValue.MapRange()
	destFieldValue.Set(reflect.MakeMap(mapToMapMapping.DestValueType))
	for sourceMapIter.Next() {
		val, err := mappingInfo.mapper(sourceMapIter.Value().Interface())
		if err != nil {
			return err
		}
		if val.IsValid() {
			if mapToMapMapping.DestValueType.Elem().Kind() == reflect.Ptr {
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

// MapToStructMapping map to struct
type MapToStructMapping struct {
	ToFields []*structField
}

// Convert match value by map key and struct field name, the result of recursively mapping the value as struct field value
func (mapToStructMapping *MapToStructMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	values := deepValue(destFieldValue)
	if sourceFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(nil))
	}
	sourceMap := reflect.Indirect(sourceFieldValue)
	mapIter := sourceMap.MapRange()
	for mapIter.Next() {
		for index, valueField := range values {
			if mapToStructMapping.ToFields[index].Name() == mapIter.Key().Interface().(string) {
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

// StructToMapMapping struct to map
type StructToMapMapping struct {
	FromFields []*structField
}

// Convert match value by struct field name and map key, the result of recursively mapping the value as new map value
func (structToMapMapping *StructToMapMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	values := deepValue(sourceFieldValue)
	if destFieldValue.IsNil() {
		destFieldValue.Set(reflect.ValueOf(map[string]interface{}{}))
	}

	for _, field := range structToMapMapping.FromFields {
		key := reflect.ValueOf(field.Name())
		val := values[field.FiledIndex]
		if val.IsValid() {
			destFieldValue.SetMapIndex(key, val)
		}
	}
	return nil
}

// Array2ArrayMapping array to array
type Array2ArrayMapping struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	Length        int
	ChildMapping  *MappingInfo
}

//Convert clone array element one by one, is ele type is different recursively element value to the new array element value
func (array2ArrayMapping *Array2ArrayMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	for i := 0; i < array2ArrayMapping.Length; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), array2ArrayMapping.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}
	return nil
}

// Array2ArrayMapping slice to array
type Slice2ArrayMapping struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ArrayLen      int
	ChildMapping  *MappingInfo
}

//Convert only clone minlength of 2, element copy like Array2ArrayMapping
func (slice2ArrayMapping *Slice2ArrayMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	if destFieldValue.Len() != slice2ArrayMapping.ArrayLen {
		return ErrLengthNotMatch
	}
	for i := 0; i < slice2ArrayMapping.ArrayLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), slice2ArrayMapping.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}

	return nil
}

// Slice2SliceMapping slice to slice
type Slice2SliceMapping struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ChildMapping  *MappingInfo
}

//Convert only clone minlength of 2, element copy like Array2ArrayMapping
func (slice2SliceMapping *Slice2SliceMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	sourceLen := sourceFieldValue.Len()
	destFieldValue.Set(reflect.MakeSlice(destFieldValue.Type(), sourceLen, sourceLen))
	for i := 0; i < sourceLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), slice2SliceMapping.ToFieldType)
		if err != nil {
			return err
		}
		setValue(destFieldValue.Index(i), childVal)
	}
	return nil
}

// Array2SliceMapping array to slice
type Array2SliceMapping struct {
	FromFieldType reflect.Type
	ToFieldType   reflect.Type
	ArrayLen      int
	ChildMapping  *MappingInfo
}

//Convert only clone minlength of 2, element copy like Array2ArrayMapping
func (array2SliceMapping *Array2SliceMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	if sourceFieldValue.Len() != array2SliceMapping.ArrayLen {
		return ErrLengthNotMatch
	}
	newSlices := reflect.MakeSlice(destFieldValue.Type(), array2SliceMapping.ArrayLen, array2SliceMapping.ArrayLen)
	for i := 0; i < array2SliceMapping.ArrayLen; i++ {
		childVal, err := mapper(sourceFieldValue.Index(i).Interface(), array2SliceMapping.ToFieldType)
		if err != nil {
			return err
		}
		setValue(newSlices.Index(i), childVal)
	}
	destFieldValue.Set(newSlices)
	return nil
}

// SameTypeMapping used to assign values between the same types
type SameTypeMapping struct {
	SourceType reflect.Type
	DestType   reflect.Type
}

//Convert set directly
func (sameTypeMapping *SameTypeMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

// AnyMapping for interace{} dest value
type AnyMapping struct {
}

//Convert interace{} can receive any value , set directly
func (anyMapping *AnyMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	destFieldValue.Set(reflect.ValueOf(sourceFieldValue.Interface()))
	return nil
}

// StructFieldMapping each field of the structure corresponds to a mapping relationship
type StructFieldMapping struct {
	FromField    *structField
	ToField      *structField
	ChildMapping *MappingInfo
}

// Convert invoke field mappig and use reslut as field value
func (structFieldMap *StructFieldMapping) Convert(sourceFieldValue reflect.Value, destFieldValue reflect.Value) error {
	val, err := structFieldMap.ChildMapping.mapper(sourceFieldValue.Interface())
	if err != nil {
		return err
	}
	setValue(destFieldValue, val)
	return nil
}
