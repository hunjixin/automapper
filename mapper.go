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
	if sourceType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
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
	// get deep field and group fields by name
	allSourceFileds := deepFields(sourceType)
	sourceGroupFields := groupFiled(allSourceFileds)
	destFileds := deepFields(destType)
	destGroupFields := groupFiled(destFileds)

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

	/*
	for _, mappingInfosMap := range mapperStore {
		for _, oldMappingInfo := range mappingInfosMap {
			for i := len(oldMappingInfo.MapFileds) - 1; i > -1; i-- {
				mapField := oldMappingInfo.MapFileds[i]
				if mapField.GetType() == None &&
					indirectType(mapField.GetFromField().StructField.Type) == sourceType &&
					indirectType(mapField.GetToField().StructField.Type) == destType {
					field := mapField.(*NoneMappingField)
					childMapField := &ChildrenMappingField{
						BaseMappingField{
							Type:      ChildMap,
							FromField: field.GetFromField(),
							ToField:   field.GetToField(),
						},
						newMappingInfo,
					}
					oldMappingInfo.MapFileds = append(oldMappingInfo.MapFileds, childMapField)
					oldMappingInfo.MapFileds = append(oldMappingInfo.MapFileds[:i], oldMappingInfo.MapFileds[i+1:]...)
				}
			}
		}
	}
	*/
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

	destFieldValues := deepValue(destValue)
	sourceFields 	:= deepValue(sourceValue)

	for _, mappingField := range mappingInfo.MapFileds {
		sourceFieldValue := sourceFields[mappingField.GetFromField().FiledIndex]
		if (sourceFieldValue.Kind() == reflect.Ptr ||
			sourceFieldValue.Kind() == reflect.Map ||
			sourceFieldValue.Kind() == reflect.Slice) &&
			sourceFieldValue.IsNil() {
			continue
		} else {
			destFieldValue := destFieldValues[mappingField.GetToField().FiledIndex]
			err := mappingField.Convert(sourceFieldValue, destFieldValue)
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
