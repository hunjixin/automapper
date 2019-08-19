package automapper

import (
	"github.com/hraban/lrucache"
	"reflect"
	"sync"
)

var (
	mapperStore = map[reflect.Type]map[reflect.Type]*MappingInfo{}
	lock        sync.Mutex
	cache       = lrucache.New(1025)

	ByteType    = reflect.TypeOf(byte(0))
)

// MustCreateMapper similar to CreateMapper but ignore err
func MustCreateMapper(sourceType, destType interface{}) *MappingInfo {
	mappingInfo, _ := CreateMapper(sourceType, destType)
	return mappingInfo
}

// CreateMapper judging the mapping relationship between two types,
// there are currently array slice mapped to each other, map and structure are mapped to each other.
// mapping between structures is more complicated
// if name is 1 to 1 :use name to match
// if name is 1 to many or many to 1: use key path to match
// if name is many to many :  use key path to match. may exist match more than one
func CreateMapper(sourceType, destType interface{}) (*MappingInfo, error) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	return createMapper(reflect.TypeOf(sourceType), reflect.TypeOf(destType))
}

func createMapper(sourceType, destType reflect.Type) (*MappingInfo, error) {
	newMappingInfo := &MappingInfo{
		Key:        sourceType.String() + "=>" + destType.String(),
		SourceType: sourceType,
		DestType:   destType,
		MapFileds:  []IStructConverter{},
		MapFunc:    []func(reflect.Value, interface{}){},
	}

	mappingInfosMap, ok := mapperStore[sourceType]
	if ok {
		oldmappingInfo, ok2 := mappingInfosMap[sourceType]
		if ok2 {
			return oldmappingInfo, nil
		} else {
			mappingInfosMap[destType] = newMappingInfo
		}
	} else {
		mapperStore[sourceType] = map[reflect.Type]*MappingInfo{destType: newMappingInfo}
	}

	if sourceType == destType {
		newMappingInfo.Type = SameType
		newMappingInfo.AddField(&SameTypeMapping{
			sourceType,
			destType,
		})
		goto End
	}

	//PTR=>PTR
	if sourceType.Kind() == reflect.Ptr && destType.Kind() == reflect.Ptr {
		newMappingInfo.Type = PtrToPtr
		childMapping, _ := ensureMapping(sourceType.Elem(), destType.Elem())
		newMappingInfo.AddField(&PtrToPtrMapping{
			childMapping,
		})
		goto End
	}
	//PTR
	if sourceType.Kind() == reflect.Ptr {
		newMappingInfo.Type = Ptr
		childMapping, _ := ensureMapping(sourceType.Elem(), destType)
		newMappingInfo.AddField(&PtrMapping{
			childMapping,
			true,
		})
		goto End
	}
	if destType.Kind() == reflect.Ptr {
		newMappingInfo.Type = Ptr
		childMapping, _ := ensureMapping(sourceType, destType.Elem())
		newMappingInfo.AddField(&PtrMapping{
			childMapping,
			false,
		})
		goto End
	}

	// map => struct
	if isString2InterfaceMap(sourceType) && destType.Kind() == reflect.Struct {
		newMappingInfo.Type = MapToStruct
		newMappingInfo.AddField(&MapToStructMapping{
			deepFields(destType),
		})
		goto End
	}

	//struct => map
	if isString2InterfaceMap(destType) && sourceType.Kind() == reflect.Struct {
		newMappingInfo.Type = StructToMap
		newMappingInfo.AddField(&StructToMapMapping{
			deepFields(sourceType),
		})
		goto End
	}

	//Slice => Array
	if sourceType.Kind() == reflect.Slice && destType.Kind() == reflect.Array {
		childMapping, _ := ensureMapping(sourceType.Elem(), destType.Elem())
		newMappingInfo.Type = SliceToArray
		newMappingInfo.AddField(&Slice2ArrayMapping{
			sourceType.Elem(),
			destType.Elem(),
			sourceType.Len(),
			childMapping,
		})
		goto End
	}
	//Array => Slice
	if sourceType.Kind() == reflect.Array && destType.Kind() == reflect.Slice {
		newMappingInfo.Type = ArrayToSlice
		childMapping, _ := ensureMapping(sourceType.Elem(), destType.Elem())
		newMappingInfo.Type = ArrayToArray
		newMappingInfo.AddField(&Array2SliceMapping{
			sourceType.Elem(),
			destType.Elem(),
			sourceType.Len(),
			childMapping,
		})
		goto End
	}

	// type below SameType can set directly
	//Array => Array
	if sourceType.Kind() == reflect.Array && destType.Kind() == reflect.Array {
		childMapping, _ := ensureMapping(sourceType.Elem(), destType.Elem())
		newMappingInfo.Type = ArrayToArray
		newMappingInfo.AddField(&Array2ArrayMapping{
			sourceType.Elem(),
			destType.Elem(),
			sourceType.Len(),
			childMapping,
		})
		goto End
	}

	//Slice => Slice
	if sourceType.Kind() == reflect.Slice && destType.Kind() == reflect.Slice {
		newMappingInfo.Type = ArrayToSlice
		childMapping, _ := ensureMapping(sourceType.Elem(), destType.Elem())
		newMappingInfo.Type = SliceToSlice
		newMappingInfo.AddField(&Slice2SliceMapping{
			sourceType.Elem(),
			destType.Elem(),
			childMapping,
		})
		goto End
	}

	//map => map
	if destType.Kind() == reflect.Map && sourceType.Kind() == reflect.Map {
		newMappingInfo.Type = MapToMap
		newMappingInfo.AddField(&MapToMapMapping{
			sourceType,
			destType,
		})
		goto End
	}

	//struct => struct
	if sourceType.Kind() == reflect.Struct && destType.Kind() == reflect.Struct {
		// get deep field and group fields by name
		allSourceFileds := deepFields(sourceType)
		sourceGroupFields := groupFiled(allSourceFileds)
		destFileds := deepFields(destType)
		destGroupFields := groupFiled(destFileds)
		newMappingInfo.FromFields = allSourceFileds
		newMappingInfo.ToField = destFileds

		for name, oneSourceGroupField := range sourceGroupFields {
			oneDestGoupField, exist := destGroupFields[name]
			if !exist {
				continue
			}
			if len(oneSourceGroupField) == 1 {
				if len(oneDestGoupField) == 1 {
					// 1 to 1
					sourceField := oneSourceGroupField[0]
					destField := oneDestGoupField[0]
					newMappingInfo.Type = StructToStrucgt
					childMapping, _ := ensureMapping(sourceField.Type, destField.Type)
					newMappingInfo.AddField(&StructFieldMapping{sourceField, destField, childMapping})
				} else {
					//1 to many
					sourceField := oneSourceGroupField[0]
					for j := 0; j < len(oneDestGoupField); j++ {
						destField := oneDestGoupField[j]
						if sourceField.Path == destField.Path {
							newMappingInfo.Type = StructToStrucgt
							childMapping, _ := ensureMapping(sourceField.Type, destField.Type)
							newMappingInfo.AddField(&StructFieldMapping{sourceField, destField, childMapping})
						}
					}
				}
			} else {
				if len(oneDestGoupField) == 1 {
					// many to 1
					destField := oneDestGoupField[0]
					for j := 0; j < len(oneSourceGroupField); j++ {
						sourceField := oneSourceGroupField[j]
						if sourceField.Path == destField.Path {
							newMappingInfo.Type = StructToStrucgt
							childMapping, _ := ensureMapping(sourceField.Type, destField.Type)
							newMappingInfo.AddField(&StructFieldMapping{sourceField, destField, childMapping})
							break
						}
					}
				} else {
					//many to many
					for i := 0; i < len(oneSourceGroupField); i++ {
						sourceField := oneSourceGroupField[i]
						for j := 0; j < len(oneDestGoupField); j++ {
							destField := oneDestGoupField[j]
							if sourceField.Path == destField.Path {
								newMappingInfo.Type = StructToStrucgt
								childMapping, _ := ensureMapping(sourceField.Type, destField.Type)
								newMappingInfo.AddField(&StructFieldMapping{sourceField, destField, childMapping})
							}
						}
					}
				}
			}
		}
		goto End
	}
End:
	//add simply type convert mapping such as int to string

	return newMappingInfo, nil
}

// groupFiled group field by name
func groupFiled(fileds []*structField) map[string][]*structField {
	groupFileds := map[string][]*structField{}
	for _, field := range fileds {
		oneGroupFields, exist := groupFileds[field.Name()]
		if exist {
			oneGroupFields = append(oneGroupFields, field)
			groupFileds[field.Name()] = oneGroupFields
		} else {
			groupFileds[field.Name()] = []*structField{field}
		}
	}
	return groupFileds
}

// MustMapper similar to Mapper just ignore error
func MustMapper(source interface{}, destType reflect.Type) interface{} {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	reflectValue, _ := mapper(source, destType)
	return reflectValue.Interface()
}

// mapper convert source to destType
func Mapper(source interface{}, destType reflect.Type) (interface{}, error) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	reflectValue, err := mapper(source, destType)
	if err != nil {
		return nil, err
	}
	return reflectValue.Interface(), err
}

func mapper(source interface{}, destType reflect.Type) (reflect.Value, error) {
	mapping, _ := ensureMapping(reflect.TypeOf(source), destType)
	return mapping.mapper(reflect.ValueOf(source))
}

func isNil(val reflect.Value) bool {
	if val.Kind() == reflect.Invalid {
		return true
	}
	if (val.Kind() == reflect.Chan ||
		val.Kind() == reflect.Func ||
		val.Kind() == reflect.Chan ||
		val.Kind() == reflect.Map ||
		val.Kind() == reflect.Ptr ||
		val.Kind() == reflect.UnsafePointer ||
		val.Kind() == reflect.Interface ||
		val.Kind() == reflect.Slice) && val.IsNil() {
		return true
	}
	return false
}

// EnsureMapping get mapping relationship by source and dest type if not exist auto create
func EnsureMapping(sourceType, destType interface{}) *MappingInfo {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	mappingInfo, _ := ensureMapping(reflect.TypeOf(sourceType), reflect.TypeOf(destType))
	return mappingInfo
}

// ensureMapping get mapping by source and dest type if not exist auto create
func ensureMapping(sourceType, destType reflect.Type) (*MappingInfo, bool) {
	mappingInfosMaps, ok := mapperStore[sourceType]
	if !ok {
		mapping, _ := createMapper(sourceType, destType)
		return mapping, true
	}
	mappingInfosMap, ok := mappingInfosMaps[destType]
	if !ok {
		mapping, _ := createMapper(sourceType, destType)
		return mapping, true
	}
	return mappingInfosMap, false
}
