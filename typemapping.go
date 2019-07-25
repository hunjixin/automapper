package automapper

import "reflect"

const (
	None = iota
	AnyType
	SameType
	MapToMap
	StructToMap
	MapToStruct
	ArrayToArray
	ArrayToSlice
	SliceToArray
	SliceToSlice
	StructToStrucgt
	StructField
)

// MappingInfo recored field mapping information
type MappingInfo struct {
	Key        string
	Type       int
	SourceType reflect.Type
	DestType   reflect.Type
	FromFields []*structField
	ToField    []*structField

	MapFileds  []IStructConverter
	MapFunc    []func(interface{}, interface{})
}

func (mappingInfo *MappingInfo) AddField(field IStructConverter) {
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, field)
}

// Mapping add customize field mapping
func (mappingInfo *MappingInfo) Mapping(mapFunc func(interface{}, interface{})) *MappingInfo {
	mappingInfo.MapFunc = append(mappingInfo.MapFunc, mapFunc)
	return mappingInfo
}

// Mapper get a instance by given source value and dest type
func (mappingInfo *MappingInfo) mapper(source interface{}) (reflect.Value, error) {
	originSourceValue := reflect.ValueOf(source)
	sourceValue := reflect.Indirect(originSourceValue)
	destValue := reflect.ValueOf(nil)
	if isNil(originSourceValue) {
		return destValue, nil
	}
	destValue = reflect.New(indirectType(mappingInfo.DestType)).Elem()
	switch mappingInfo.Type {
	case None:
	case AnyType:
		//mappingInfo.
		//TODO
	case SameType:
		err := mappingInfo.MapFileds[0].Convert(sourceValue, destValue)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	case ArrayToArray:
		fallthrough
	case ArrayToSlice:
		fallthrough
	case SliceToArray:
		fallthrough
	case SliceToSlice:
		fallthrough
	case MapToMap:
		fallthrough
	case StructToMap:
		fallthrough
	case MapToStruct:
		err := mappingInfo.MapFileds[0].Convert(sourceValue, destValue)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	case StructToStrucgt:
		destFieldValues := deepValue(destValue)
		sourceFields := deepValue(sourceValue)
		for _, mappingField := range mappingInfo.MapFileds {
			structFieldMapping := mappingField.(*StructFieldMapping)
			sourceFieldValue := sourceFields[structFieldMapping.FromField.FiledIndex]
			destFieldValue := destFieldValues[structFieldMapping.ToField.FiledIndex]
			err := structFieldMapping.Convert(sourceFieldValue, destFieldValue)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
		}
	}

	for _, mapFunc := range mappingInfo.MapFunc {
		mapFunc(destValue.Addr().Interface(), sourceValue.Addr().Interface())
	}
	return destValue, nil
}

// isString2InterfaceMap map in map2struct and struct2map must be string=> Interface{}
func isString2InterfaceMap(t reflect.Type) bool {
	if t.Kind() == reflect.Map &&
		t.Key().Kind() == reflect.String &&
		t.Elem().Kind() == reflect.Interface {
		return true
	}
	return false
}
