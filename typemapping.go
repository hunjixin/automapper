package automapper

import "reflect"

const (
	None     = iota
	AnyType  = 1
	SameType = 2
	ArrayMap = 3
	ChildMap = 4
)

// MappingInfo recored field mapping information
type MappingInfo struct {
	Key        string
	SourceType reflect.Type
	DestType   reflect.Type
	MapFileds  []IStructConverter
	MapFunc    []func (interface{}, interface{})
}

// Mapping add customize field mapping
func (mappingInfo *MappingInfo) Mapping(mapFunc func (interface{}, interface{})) *MappingInfo{
	mappingInfo.MapFunc = append(mappingInfo.MapFunc, mapFunc)
	return mappingInfo
}

// tryAddNameFieldMapping analysis mapping time and add it to MapFields
func (mappingInfo *MappingInfo) tryAddNameFieldMapping(sourceFiled, destFiled *structField) bool {
	if sourceFiled.Type.Kind() == reflect.Interface {
		return mappingInfo.tryAddAnyMapping(sourceFiled, destFiled)
	}
	if sourceFiled.Type == destFiled.Type {
		return mappingInfo.tryAddSameTypeMapping(sourceFiled, destFiled)
	} else {
		indirectSourceFiled := indirectType(sourceFiled.Type)
		indirectDestFiled := indirectType(destFiled.Type)
		//deferent type mapping
		//struct -> struct
		if indirectSourceFiled.Kind() == reflect.Struct &&
			indirectDestFiled.Kind() == reflect.Struct {
			return mappingInfo.tryAddStructToStructMapping(sourceFiled, destFiled)
		}
		//struct -> map[string]interface{}
		if 	indirectSourceFiled.Kind() == reflect.Struct &&
			indirectDestFiled.Kind() == reflect.Map {
			if indirectSourceFiled.Key().Kind() == reflect.String && indirectDestFiled.Elem().Kind()  == reflect.Interface {
				return mappingInfo.tryAddStructToMapMapping(sourceFiled, destFiled)
			}
		}
		//map[string]interface{} -> struct
		if 	indirectSourceFiled.Kind() == reflect.Map &&
			indirectDestFiled.Kind() == reflect.Struct {
			if indirectSourceFiled.Key().Kind() == reflect.String && indirectSourceFiled.Elem().Kind()  == reflect.Interface {
				return mappingInfo.tryAddStructToMapMapping(sourceFiled, destFiled)
			}
		}


		//Array=>Array
		if indirectSourceFiled.Kind() == reflect.Array &&
			indirectDestFiled.Kind() == reflect.Array {
			return mappingInfo.tryAddArrayToArrayMapping(sourceFiled, destFiled)
		}
		//Slice=>Array
		if indirectSourceFiled.Kind() == reflect.Slice &&
			indirectDestFiled.Kind() == reflect.Array {
			return mappingInfo.tryAddSliceToArrayMapping(sourceFiled, destFiled)
		}
		//Array=>Slice
		if indirectSourceFiled.Kind() == reflect.Array &&
			indirectDestFiled.Kind() == reflect.Slice {
			return mappingInfo.tryAddArrayToSliceMapping(sourceFiled, destFiled)
		}
		//Slice=>Slice
		if indirectSourceFiled.Kind() == reflect.Slice &&
			indirectDestFiled.Kind() == reflect.Slice {
			return mappingInfo.tryAddSliceToSliceMapping(sourceFiled, destFiled)
		}
	}
	return false
}

func (mappingInfo *MappingInfo) tryAddAnyMapping(sourceFiled, destFiled *structField) bool {
	anyMapingField := &AnyMappingField{
		BaseMappingField{
			Type:      AnyType,
			FromField: sourceFiled,
			ToField:   destFiled,
		},
	}
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, anyMapingField)
	return true
}

func (mappingInfo *MappingInfo) tryAddSameTypeMapping(sourceFiled, destFiled *structField) bool {
	mappingField := &SameTypeMappingField{
		BaseMappingField{
			Type:      SameType,
			FromField: sourceFiled,
			ToField:   destFiled,
		},
	}
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, mappingField)
	return true
}

func (mappingInfo *MappingInfo) tryAddStructToStructMapping(sourceFiled, destFiled *structField) bool {
	mapping, _ := ensureMapping(sourceFiled.Type, destFiled.Type)
	mappingField := &ChildrenMappingField{
		BaseMappingField{
			Type:      ChildMap,
			FromField: sourceFiled,
			ToField:   destFiled,
		},
		mapping,
	}
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, mappingField)
	return true
}

func (mappingInfo *MappingInfo) tryAddStructToMapMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
}

func (mappingInfo *MappingInfo) tryAddArrayToArrayMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
}

func (mappingInfo *MappingInfo) tryAddArrayToSliceMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
}

func (mappingInfo *MappingInfo) tryAddSliceToArrayMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
}

func (mappingInfo *MappingInfo) tryAddSliceToSliceMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
}

func (mappingInfo *MappingInfo) tryAddMapToStructMapping(sourceFiled, destFiled *structField) bool{
	//TODO
	panic("TODO")
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
		t.Kind() == reflect.Int8{
		return true
	}
	return false
}





