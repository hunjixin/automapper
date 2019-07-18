package automapper

import (
	"reflect"
	"sync"
)

var (
	mapper     = map[reflect.Type]map[reflect.Type]*MappingInfo{}
	mapperLock = sync.RWMutex{}
)

// StructField wrap reflect.StructField.  FiledIndex recored field index include enbed types, Key indirect the path to enbed field like .Enbed.XXX
type StructField struct {
	reflect.StructField
	FiledIndex int
	Key        string
}

// MappingInfo recored field mapping information
type MappingInfo struct {
	Key        string
	SourceType reflect.Type
	DestType   reflect.Type
	MapFileds  map[*StructField]*StructField
}

// CreateMapper build field mapping between sourceType and destType
// if name is 1 to 1 and map derect
// if name is 1 to many or many to 1: use key path to match
// if name is many to many :  use key path to match. may exist match more than one
func CreateMapper(sourceType, destType reflect.Type) error {
	if sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
	}
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if sourceType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	mappingInfo := &MappingInfo{}
	mappingInfo.SourceType = sourceType
	mappingInfo.DestType = destType
	mappingInfo.MapFileds = map[*StructField]*StructField{}

	// get deep field
	allSourceFileds := deepFields(sourceType)
	destFileds := deepFields(destType)
	//group fields by name
	sourceGroupFields := groupFiled(allSourceFileds)
	destGroupFields := groupFiled(destFileds)

	for name, oneSourceGroupField := range sourceGroupFields {
		oneDestGoupField, exist := destGroupFields[name]
		if !exist {
			continue
		}
		if len(oneSourceGroupField) == 1 {
			if len(oneDestGoupField) == 1 {
				// 1 to 1
				if oneSourceGroupField[0].Type == oneDestGoupField[0].Type {
					mappingInfo.MapFileds[oneSourceGroupField[0]] = oneDestGoupField[0]
				}
			} else {
				//1 to many
				sourceFiled := oneSourceGroupField[0]
				for j := 0; j < len(oneDestGoupField); j++ {
					destField := oneDestGoupField[j]
					if sourceFiled.Key == destField.Key && sourceFiled.Type == destField.Type {
						mappingInfo.MapFileds[sourceFiled] = destField
					}
				}
			}
		} else {
			if len(oneDestGoupField) == 1 {
				// many to 1
				destField := oneDestGoupField[0]
				for j := 0; j < len(oneSourceGroupField); j++ {
					sourceFiled := oneSourceGroupField[j]
					if sourceFiled.Key == destField.Key && sourceFiled.Type == destField.Type {
						mappingInfo.MapFileds[sourceFiled] = destField
						break
					}
				}
			} else {
				//many to many
				for i := 0; i < len(oneSourceGroupField); i++ {
					sourceFiled := oneSourceGroupField[i]
					for j := 0; j < len(oneDestGoupField); j++ {
						destField := oneDestGoupField[j]
						if sourceFiled.Key == destField.Key && sourceFiled.Type == destField.Type {
							mappingInfo.MapFileds[sourceFiled] = destField
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
		_, ok2 := mappingInfosMap[sourceType]
		if ok2 {
			return nil
		} else {
			mappingInfosMap[destType] = mappingInfo
		}
	} else {
		mapper[sourceType] = map[reflect.Type]*MappingInfo{destType: mappingInfo}
	}

	return nil
}

// groupFiled group field by name
func groupFiled(fileds []*StructField) map[string][]*StructField {
	groupFileds := map[string][]*StructField{}
	for _, field := range fileds {
		oneGroupFields, exist := groupFileds[field.Name]
		if exist {
			oneGroupFields = append(oneGroupFields, field)
			groupFileds[field.Name] = oneGroupFields
		} else {
			groupFileds[field.Name] = []*StructField{field}
		}
	}
	return groupFileds
}

func contain(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// MustMapper like Mapper just ignore error
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
	mapperLock.RLock()
	defer func() {
		mapperLock.RUnlock()
	}()
	sourceType := reflect.TypeOf(source)
	sourceValue := deepValue(reflect.ValueOf(source))
	mappingInfosMap, ok := mapper[sourceType]
	if !ok {
		return nil, ErrMapperNotDefine
	}

	mappingInfo, ok := mappingInfosMap[destType]
	if !ok {
		return nil, ErrMapperNotDefine
	}
	destValue := reflect.New(destType).Elem()
	destFieldValues := deepValue(destValue)
	for sourcerField, destFiled := range mappingInfo.MapFileds {
		sourceValue := sourceValue[sourcerField.FiledIndex]
		if sourceValue.CanAddr()&&sourceValue.IsNil() {
			continue
		}else{
			field := destFieldValues[destFiled.FiledIndex]
			field.Set(reflect.ValueOf(sourceValue.Interface()))
		}
	}

	if isDestPtr {
		return destValue.Addr().Interface(), nil
	} else {
		return destValue.Interface(), nil
	}
}

func deepValue(ifv reflect.Value) []reflect.Value {
	fields := make([]reflect.Value, 0)
	for i := 0; i < ifv.Type().NumField(); i++ {
		v := ifv.Field(i)
		t := ifv.Type().Field(i)
		if v.Kind() == reflect.Struct && t.Anonymous{
			fields = append(fields, deepValue(v)...)
		}else{
			fields = append(fields, v)
		}
	}

	return fields
}

func deepFields(ift reflect.Type) []*StructField {
	return internalDeepFields(ift, &intVal{0}, "")
}

func GetMapping(sourceType, destType reflect.Type) (*MappingInfo, error) {
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
	mappingInfosMap, ok :=mappingInfosMaps[destType]
	if !ok {
		return nil, ErrMapperNotDefine
	}
	return mappingInfosMap, nil
}

type intVal struct {
	I int
}

func (intval *intVal) plus(i int) {
	intval.I = intval.I + i
}

func internalDeepFields(ift reflect.Type, index *intVal, key string) []*StructField {
	fields := make([]*StructField, 0)
	for i := 0; i < ift.NumField(); i++ {
		f := ift.Field(i)
		newKey := key + "." + f.Name
		if f.Type.Kind() ==  reflect.Struct && f.Anonymous {
			fields = append(fields, internalDeepFields(f.Type, index, newKey)...)
		}else{
			fields = append(fields, &StructField{f, index.I, newKey})
			index.plus(1)
		}
	}
	return fields
}
