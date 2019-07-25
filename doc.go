/*
Package automapper provides data mapping between different data struct include strcut-> struct, map->struct, struct->map. and you can also convert  between array and slice each other
Getting Started
To get started, you'll want to import the library:

	import log "github.com/hunjixin/automapper"

register type, this is not required, the library will automatically register the type

	automapper.CreateMapper(reflect.TypeOf(User(nil)), reflect.TypeOf(UserDto(nil)))

data mapping

	automapper.MustMapper(User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)},  reflect.TypeOf(UserDto(nil)))

struct mapping

The two structures match the name (case sensitive), and then determine the mapping relationship between the two structure corresponding field types. If the type is the same,
it belongs to SameType, otherwise it will generate the submap recursively.
The library supports the embed type. When traversing the field, the fields that traverse the embed are merged together. If the names are the same, the access path is matched.

	type A struct {
	   M string
	}

	type B struct {
	   A
	   M string
	}

	type C struct {
	   A
	   M string
	}

In this example, the conversion of B->C generates two rules, M-M and [A].M->[A].M., so that the field values can be correctly matched.

Map conversion

The mapping between map and map is similar to the structure, matching the same key, and then mapping the corresponding value.

When mapping between map and structure, map must be map[string]interface{} type, structure field type is different,  only using interface{} work

Slice and Array mapping

Arrays and slices can be mapped between two and two. Map the child elements one by one, and then set the result of the mapping to the new value.
When an array is mapped to a slice or between a slice and a slice, only the element of the minimum length is converted.
If the length of the array is 10 and the length of the slice is 5, only the first 5 elements are converted.

Func to customize field mapping content

sometimes we need a more free mapping scheme, such as combining two fields together under certain conditions.

	automapper.MustCreateMapper(reflect.TypeOf((*User)(nil)), reflect.TypeOf((*UserDto)(nil))).
	Mapping(func(destVal interface{}, sourceVal interface{}) {
		destVal.(*UserDto).Name = sourceVal.(*User).Name + "|" + sourceVal.(*User).Nick
	}).
	Mapping(func(destVal interface{}, sourceVal interface{}) {
		destVal.(*UserDto).Age = time.Now().Year() - sourceVal.(*User).Birth.Year()
	})
*/
package automapper