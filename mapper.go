package automapper

import (
	"github.com/hraban/lrucache"
	"reflect"
	"sync"
)

var (
	mapper     = map[reflect.Type]map[reflect.Type]*MappingInfo{}
	mapperLock = sync.RWMutex{}
	cache      = lrucache.New(1025)
)

const (
	None     = iota
	SameNone = 1
	SameType = 2
	ChildMap = 3
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
	if sourceFiled.Type == destFiled.Type {
		mappingField := &SameTypeMappingField{
			BaseMappingField{
				Type:      SameType,
				FromField: sourceFiled,
				ToField:   destFiled,
			},
		}
		mappingInfo.MapFileds = append(mappingInfo.MapFileds, mappingField)
		return true
	} else {
		childMappingInfo, err := getMapping(sourceFiled.Type, destFiled.Type)
		if err == nil {
			mappingField := &ChildrenMappingField{
				BaseMappingField{
					Type:      ChildMap,
					FromField: sourceFiled,
					ToField:   destFiled,
				},
				childMappingInfo,
			}
			mappingInfo.MapFileds = append(mappingInfo.MapFileds, mappingField)
			return true
		}
	}
	mappingField := &NoneMappingField{
		BaseMappingField{
			Type:      SameNone,
			FromField: sourceFiled,
			ToField:   destFiled,
		},
	}
	mappingInfo.MapFileds = append(mappingInfo.MapFileds, mappingField)
	return false
}

// MustCreateMapper similar CreateMapper just ignore error
func MustCreateMapper(sourceType, destType reflect.Type) *MappingInfo {
	mappingInfo, _ := CreateMapper(sourceType, destType)
	return mappingInfo
}
// CreateMapper build field mapping between sourceType and destType
// if name is 1 to 1 and map derect
// if name is 1 to many or many to 1: use key path to match
// if name is many to many :  use key path to match. may exist match more than one
func CreateMapper(sourceType, destType reflect.Type) (*MappingInfo, error) {
	sourceType = toStructType(sourceType)
	destType = toStructType(destType)
	if sourceType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	mappingInfo := &MappingInfo{
		SourceType: sourceType,
		DestType:   destType,
		MapFileds:  []IStructConverter{},
		MapFunc:	[]func (interface{}, interface{}){},
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
				mappingInfo.tryAddNameFieldMapping(sourceField, destField)
			} else {
				//1 to many
				sourceFiled := oneSourceGroupField[0]
				for j := 0; j < len(oneDestGoupField); j++ {
					destField := oneDestGoupField[j]
					if sourceFiled.Path == destField.Path {
						mappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
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
						mappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
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
							mappingInfo.tryAddNameFieldMapping(sourceFiled, destField)
						}
					}
				}
			}
		}
	}
	mapperLock.Lock()
	defer func() {
		mapperLock.Unlock()
	}()
	mappingInfosMap, ok := mapper[sourceType]
	if ok {
		oldmappingInfo, ok2 := mappingInfosMap[sourceType]
		if ok2 {
			return oldmappingInfo, nil
		} else {
			mappingInfosMap[destType] = mappingInfo
		}
	} else {
		mapper[sourceType] = map[reflect.Type]*MappingInfo{destType: mappingInfo}
	}

	for _, mappingInfosMap := range mapper {
		for _, mappingInfo := range mappingInfosMap {
			for i := len(mappingInfo.MapFileds) - 1; i > -1; i-- {
				mapField := mappingInfo.MapFileds[i]
				if mapField.GetType() == SameNone &&
					toStructType(mapField.GetFromField().StructField.Type) == sourceType &&
					toStructType(mapField.GetToField().StructField.Type) == destType {
					field := mapField.(*NoneMappingField)
					childMapField := &ChildrenMappingField{
						BaseMappingField{
							Type:      ChildMap,
							FromField: field.GetFromField(),
							ToField:   field.GetToField(),
						},
						mappingInfo,
					}
					mappingInfo.MapFileds = append(mappingInfo.MapFileds, childMapField)
					mappingInfo.MapFileds = append(mappingInfo.MapFileds[:i], mappingInfo.MapFileds[i+1:]...)
				}
			}
		}
	}
	return mappingInfo, nil
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
	val, _ := Mapper(source, destType)
	return val
}

// Mapper get a instance by given source value and dest type
func Mapper(source interface{}, destType reflect.Type) (interface{}, error) {
	isDestPtr := false
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
		isDestPtr = true
	}
	sourceValue := reflect.Indirect(reflect.ValueOf(source))
	sourceType := sourceValue.Type()
	if destType.Kind() != reflect.Struct||sourceType.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	mappingInfo, err := getMapping(sourceType, destType)
	if err != nil {
		return nil, err
	}
	mapperLock.RLock()
	defer func() {
		mapperLock.RUnlock()
	}()

	destValue := reflect.New(destType).Elem()
	destFieldValues := deepValue(destValue)
	sourceFields := deepValue(sourceValue)

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
				return nil, err
			}
		}
	}
	for _, mapFunc := range mappingInfo.MapFunc {
		mapFunc(destValue.Addr().Interface(), sourceValue.Interface())
	}
	if isDestPtr {
		return destValue.Addr().Interface(), nil
	} else {
		return destValue.Interface(), nil
	}
}

func containMapping(sourceType, destType reflect.Type) bool {
	_, err := getMapping(sourceType, destType)
	if err != nil {
		return false
	} else {
		return true
	}
}

func toStructType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func getMapping(sourceType, destType reflect.Type) (*MappingInfo, error) {
	if sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
	}
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	mappingInfosMaps, ok := mapper[sourceType]
	if !ok {
		return nil, ErrMapperNotDefine
	}
	mappingInfosMap, ok := mappingInfosMaps[destType]
	if !ok {
		return nil, ErrMapperNotDefine
	}
	return mappingInfosMap, nil
}
