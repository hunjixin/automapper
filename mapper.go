package automapper

import (
	"github.com/hraban/lrucache"
	"reflect"
	"sync"
)

var (
	mapperStore = map[reflect.Type]map[reflect.Type]*MappingInfo{}
	lock        sync.Mutex
	cache             = lrucache.New(1025)
)

// MustCreateMapper similar CreateMapper just ignore error
func MustCreateMapper(sourceType, destType reflect.Type) *MappingInfo {
	mappingInfo, _ := CreateMapper(sourceType, destType)
	return mappingInfo
}

func CreateMapper(sourceType, destType reflect.Type) (*MappingInfo, error) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	return createMapper(sourceType, destType)
}

// CreateMapper build field mapping between sourceType and destType
// if name is 1 to 1 and map derect
// if name is 1 to many or many to 1: use key path to match
// if name is many to many :  use key path to match. may exist match more than one
func createMapper(sourceType, destType reflect.Type) (*MappingInfo, error) {
	sourceType = indirectType(sourceType)
	destType = indirectType(destType)
	newMappingInfo := &MappingInfo{
		Key: sourceType.String() +"=>" + destType.String(),
		SourceType: sourceType,
		DestType:   destType,
		MapFileds:  []IStructConverter{},
		MapFunc:	[]func (interface{}, interface{}){},
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
	//map => map
	if isString2InterfaceMap(destType) && isString2InterfaceMap(sourceType) {
		newMappingInfo.Type = MapToMap
		newMappingInfo.AddField(&MapToMapMappingField{
			sourceType.Elem(),
			destType.Elem(),
		})
		goto End
	}

	// map => struct
	if isString2InterfaceMap(sourceType) && destType.Kind() == reflect.Struct {
		newMappingInfo.Type = MapToStruct
		newMappingInfo.AddField(&MapToStructMappingField{
		})
		goto End
	}

	//struct => map
	if isString2InterfaceMap(destType) && sourceType.Kind() == reflect.Struct {
		newMappingInfo.Type = StructToMap
		newMappingInfo.AddField(&MapToStructMappingField{
		})
		goto End
	}

	//Array => Array
	if sourceType.Kind() == reflect.Array && destType.Kind() == reflect.Array {
		childMapping, _ := ensureMapping(sourceType, destType)
		newMappingInfo.Type = ArrayToArray
		newMappingInfo.AddField(&Array2ArrayMappingField{
			sourceType.Elem(),
			destType.Elem(),
			sourceType.Len(),
			childMapping,
		})
		goto End
	}
	//Slice => Array
	if sourceType.Kind() == reflect.Slice && destType.Kind() == reflect.Array {
		newMappingInfo.Type = SliceToArray
		newMappingInfo.AddField(&Slice2ArrayMappingField{
		})
		goto End
	}
	//Array => Slice
	if sourceType.Kind() == reflect.Array && destType.Kind() == reflect.Slice {
		newMappingInfo.Type = ArrayToSlice
		newMappingInfo.AddField(&Array2SliceMappingField{
		})
		goto End
	}
	//Slice => Slice
	if sourceType.Kind() == reflect.Slice && destType.Kind() == reflect.Slice {
		newMappingInfo.Type = SliceToSlice
		newMappingInfo.AddField(&Slice2SliceMappingField{
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
					newMappingInfo.tryAddNameFieldMapping(sourceField, destField)
				} else {
					//1 to many
					sourceFiled := oneSourceGroupField[0]
					for j := 0; j < len(oneDestGoupField); j++ {
						destField := oneDestGoupField[j]
						if sourceFiled.Path == destField.Path {
							newMappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
						}
					}
				}
			} else {
				if len(oneDestGoupField) == 1 {
					// many to 1
					destField := oneDestGoupField[0]
					for j := 0; j < len(oneSourceGroupField); j++ {
						sourceFiled := oneSourceGroupField[j]
						if sourceFiled.Path == destField.Path {
							newMappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
							break
						}
					}
				} else {
					//many to many
					for i := 0; i < len(oneSourceGroupField); i++ {
						sourceFiled := oneSourceGroupField[i]
						for j := 0; j < len(oneDestGoupField); j++ {
							destField := oneDestGoupField[j]
							if sourceFiled.Path == destField.Path {
								newMappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
							}
						}
					}
				}
			}
		}
		goto End
	}

	if sourceType == destType {
		newMappingInfo.Type = SameType
		newMappingInfo.AddField(&SameTypeMappingField{
			sourceType,
			destType,
		})
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
	val, _ := mapper(source, destType)
	return val
}


//Mapper similar to mapper but thread safe
func Mapper(source interface{}, destType reflect.Type) (interface{}, error) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	return mapper(source, destType)
}

// mapper convert source to destType
func mapper(source interface{}, destType reflect.Type) (interface{}, error) {
	mapping, _ := ensureMapping(reflect.TypeOf(source), destType)
	destValue, err := mapping.mapper(source)
	if err != nil {
		return nil, err
	}
	if destType.Kind() == reflect.Ptr {
		return destValue.Addr().Interface(), nil
	}else{
		return destValue.Interface(), nil
	}
}

// Mapper get a instance by given source value and dest type
func (mappingInfo *MappingInfo) mapper(source interface{}) (reflect.Value, error) {
	sourceValue   	:= reflect.Indirect(reflect.ValueOf(source))
	destValue    	:= reflect.New(indirectType(mappingInfo.DestType)).Elem()
	switch mappingInfo.Type {
		case 	None   :
		case	AnyType :
			//mappingInfo.
			//TODO
		case	SameType :
			err := mappingInfo.MapFileds[0].Convert(sourceValue, destValue)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			//TODO
		case	ArrayMap :
			//TODO
		case	MapToMap :
			//TODO
		case	StructToMap :
			//TODO
		case	MapToStruct :
			//TODO
		case	ArrayToArray :
			//TODO
		case	ArrayToSlice :
			//TODO
		case	SliceToArray :
			//TODO
		case	SliceToSlice :
			//TODO
		case StructToStrucgt:
			if !destValue.IsNil() {
				destFieldValues := deepValue(destValue)
				sourceFields 	:= deepValue(sourceValue)
				for _, mappingField := range mappingInfo.MapFileds {
					structFieldMapping :=  mappingField.(*StructFieldField)
					sourceFieldValue := sourceFields[structFieldMapping.FromField.FiledIndex]
					destFieldValue := destFieldValues[structFieldMapping.ToField.FiledIndex]
					err := structFieldMapping.Convert(sourceFieldValue, destFieldValue)
					if err != nil {
						return reflect.ValueOf(nil), err
					}
				}
			}
	}
	
	


	for _, mapFunc := range mappingInfo.MapFunc {
		mapFunc(destValue.Addr().Interface(), sourceValue.Addr().Interface())
	}
	return destValue, nil
}

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// EnsureMapping similar to ensureMapping but thread safe
func EnsureMapping(sourceType, destType reflect.Type) *MappingInfo {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	mappingInfo, _ := ensureMapping(sourceType, destType)
	return mappingInfo
}

// ensureMapping get mapping by source and dest type if not exist auto create
func ensureMapping(sourceType, destType reflect.Type) (*MappingInfo, bool) {
	if sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
	}
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	mappingInfosMaps, ok := mapperStore[sourceType]
	if !ok {
		mapping, _:=  createMapper(sourceType, destType)
		return mapping, true
	}
	mappingInfosMap, ok := mappingInfosMaps[destType]
	if !ok {
		mapping, _:=  createMapper(sourceType, destType)
		return mapping, true
	}
	return mappingInfosMap, false
}
