package automapper

import "reflect"

// StructField wrap reflect.StructField.  FiledIndex recored field index include enbed types, Path indirect the path to enbed field like .Enbed.XXX
type StructField struct {
	reflect.StructField
	TagName    string
	Ignore     bool
	FiledIndex int
	Path       string
}

func (structField *StructField) Name() string {
	if structField.TagName == "" {
		return structField.StructField.Name
	} else {
		return structField.TagName
	}
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
		if f.Type.Kind() == reflect.Chan ||
			f.Type.Kind() == reflect.Func ||
			f.Type.Kind() == reflect.UnsafePointer {
			continue
		}
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			fields = append(fields, internalDeepFields(f.Type, index, newKey)...)
		} else {
			fields = append(fields, &StructField{f, "", false, index.I, newKey})
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

func deepFields(ift reflect.Type) []*StructField {
	cacheKey := ift.String()
	if val, err := cache.Get(cacheKey); err == nil {
		//exist
		return val.([]*StructField)
	} else {
		// not exist parser and set cache
		mapFields := internalDeepFields(ift, &intVal{0}, "")
		cache.Set(cacheKey, mapFields)
		return mapFields
	}
}
