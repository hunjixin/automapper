package automapper

import "reflect"

const (
	None = iota
	AnyType
	SameType
	ArrayMap
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
	MapFileds  []IStructConverter

	FromFields []*structField
	ToField    []*structField
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

// tryAddNameFieldMapping analysis mapping time and add it to MapFields
func (mappingInfo *MappingInfo) tryAddNameFieldMapping(sourceFiled, destFiled *structField) bool {
	mappingInfo.Type = StructToStrucgt
	childMapping, _ := ensureMapping(sourceFiled.Type, destFiled.Type)
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, &StructFieldMap{sourceFiled, destFiled, childMapping})
	return true
}

// Interface/Map/Ptr/Slice/String/Struct/UnsafePointer/Array
func isSimpleType(t reflect.Type) bool {
	if t.Kind() == reflect.Bool &&
		t.Kind() == reflect.Int &&
		t.Kind() == reflect.Int8 &&
		t.Kind() == reflect.Int16 &&
		t.Kind() == reflect.Int32 &&
		t.Kind() == reflect.Int64 &&
		t.Kind() == reflect.Uint &&
		t.Kind() == reflect.Uint8 &&
		t.Kind() == reflect.Uint16 &&
		t.Kind() == reflect.Uint32 &&
		t.Kind() == reflect.Uint64 &&
		t.Kind() == reflect.Uintptr &&
		t.Kind() == reflect.Float32 &&
		t.Kind() == reflect.Float64 &&
		t.Kind() == reflect.Complex64 &&
		t.Kind() == reflect.Complex128 &&
		t.Kind() == reflect.Chan &&
		t.Kind() == reflect.Func &&
		t.Kind() == reflect.Int &&
		t.Kind() == reflect.Int8 {
		return true
	}
	return false
}

func isString2InterfaceMap(t reflect.Type) bool {
	if t.Kind() == reflect.Map &&
		t.Key().Kind() == reflect.String &&
		t.Elem().Kind() == reflect.Interface {
		return true
	}
	return false
}
