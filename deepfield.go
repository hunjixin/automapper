package automapper

import (
	"reflect"
	"strings"
)

// structField wrap reflect.structField.  FiledIndex recored field index include enbed types, Path indirect the path to enbed field like .Enbed.XXX
type structField struct {
	reflect.StructField
	MappingTag string
	Ignore     bool
	ForceName  string
	FiledIndex int
	Path       string
}

func newStructField(rField reflect.StructField, mappingTag string, fieldIndex int, path string) *structField {
	structField := &structField{
		MappingTag: mappingTag,
		FiledIndex: fieldIndex,
		Path:       path,
	}
	structField.StructField = rField
	if len(mappingTag) == 0 {
		return structField
	}
	arr := strings.Split(mappingTag, ",")
	for i := 0; i < len(arr); i++ {
		arr[i] = strings.Trim(arr[i], " \t")
	}
	except(arr, " ")
	if contain(arr, "ignore") {
		structField.Ignore = true
	}
	except(arr, "ignore")
	if len(arr) > 0 {
		structField.ForceName = arr[0]
	}
	return structField
}

func except(strs []string, s string) {
	length := len(strs)
	for i := length - 1; i > -1; i-- {
		if strs[i] == s {
			if i == length-1 {
				strs = strs[:i]
			} else {
				strs = append(strs[:i], strs[i+1:]...)
			}
		}
	}
}

func contain(strs []string, s string) bool {
	for i := 0; i < len(strs); i++ {
		if strs[i] == s {
			return true
		}
	}
	return false
}

func (structField *structField) Name() string {
	if len(structField.ForceName) > 0 { //tag
		return structField.ForceName
	}
	return structField.StructField.Name
}

type intVal struct {
	I int
}

func (intval *intVal) plus(i int) {
	intval.I = intval.I + i
}

func internalDeepFields(ift reflect.Type, index *intVal, key string) []*structField {
	fields := make([]*structField, 0)
	for i := 0; i < ift.NumField(); i++ {
		f := ift.Field(i)
		newKey := key + "[" + f.Name + "]"
		if f.Type.Kind() == reflect.Chan ||
			f.Type.Kind() == reflect.Func ||
			f.Type.Kind() == reflect.UnsafePointer {
			continue
		}
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			fields = append(fields, internalDeepFields(f.Type, index, newKey)...)
		} else {
			field := newStructField(f, f.Tag.Get("mapping"), index.I, newKey)
			fields = append(fields, field)
			index.plus(1)
		}
	}
	return fields
}

func deepValue(ifv reflect.Value) []reflect.Value {
	ifv = reflect.Indirect(ifv)
	fields := make([]reflect.Value, 0)
	for i := 0; i < ifv.Type().NumField(); i++ {
		v := ifv.Field(i)
		t := ifv.Type().Field(i)
		if t.Type.Kind() == reflect.Chan ||
			t.Type.Kind() == reflect.Func ||
			t.Type.Kind() == reflect.UnsafePointer {
			continue
		}
		if v.Kind() == reflect.Struct && t.Anonymous {
			fields = append(fields, deepValue(v)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

func deepFields(ift reflect.Type) []*structField {
	cacheKey := ift.String()
	if val, err := cache.Get(cacheKey); err == nil {
		//exist
		return val.([]*structField)
	} else {
		// not exist parser and set cache
		mapFields := internalDeepFields(ift, &intVal{0}, "")
		cache.Set(cacheKey, mapFields)
		return mapFields
	}
}
