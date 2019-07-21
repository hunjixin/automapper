# automapper
Package automapper provides data mapping between different struct

### Features

1. complex type mapping include embed field
2. child struct field mapping
3. support tag to redefine field name
4. func to customize field mapping content
5. automatic registration when no special requirement


## Example

```shell
    go get github.com/hunjixin/automapper
```

## Example

### embed mapping

```go
    package main
    
    import (
        "fmt"
        "github.com/hunjixin/automapper"
        "reflect"
    )
    
    type En struct {
        B string
        D string
    }
    
    type EnB struct {
        B string
        D string
    }
    
    type ExampleStructA struct {
        EnB
        En
        A string
    }
    
    type ExampleStructB struct {
        En
        A string
    }
    
    func main() {
        //automapper.MustCreateMapper(reflect.TypeOf((*ExampleStructA)(nil)), reflect.TypeOf((*ExampleStructB)(nil)))
    
        a := ExampleStructA{ EnB{}, En{"Sh", "Bj"},"XXXXXX"}
        result := automapper.MustMapper(a, reflect.TypeOf((*ExampleStructB)(nil)))
        fmt.Println(reflect.TypeOf(result).String())
    
        result2 := automapper.MustMapper(a, reflect.TypeOf(ExampleStructB{}))
        fmt.Println(reflect.TypeOf(result2).String())
    }
```

### tag mapping

```go
    package main
    
    import (
        "fmt"
        "github.com/hunjixin/automapper"
        "reflect"
    )
    
    type UserDto struct {
        Nick string
        Name string
    }
    
    type User struct {
        Name string  `mapping:"Nick"`
        Nick string  `mapping:"Name"`
    }
    
    func main() {
        user := &User{"NAME", "NICK"}
        result := automapper.MustMapper(user, reflect.TypeOf((*UserDto)(nil)))
        fmt.Println(reflect.TypeOf(result).String())
    }
```

### func mapping

```go
   package main
   
   import (
   	"fmt"
   	"github.com/hunjixin/automapper"
   	"reflect"
   	"time"
   )
   
   type UserDto struct {
   	Name string
   	Addr string
   	Age  int
   }
   
   type User struct {
   	Name string
   	Nick string
   	Addr string
   	Birth time.Time
   }
   
   func init() {
   	automapper.MustCreateMapper(reflect.TypeOf((*User)(nil)), reflect.TypeOf((*UserDto)(nil))).
   	Mapping(func(destVal interface{}, sourceVal interface{}) {
   		destVal.(*UserDto).Name = sourceVal.(User).Name + "|"+ sourceVal.(User).Nick
   	}).
   	Mapping(func(destVal interface{}, sourceVal interface{}) {
   		destVal.(*UserDto).Age = time.Now().Year() - sourceVal.(User).Birth.Year()
   	})
   }
   
   func main() {
   	user := &User{"NAME", "NICK", "B·J", time.Date(1992, 10,3,1,0,0,0,time.UTC)}
   	result := automapper.MustMapper(user, reflect.TypeOf((*UserDto)(nil)))
   	fmt.Println(reflect.TypeOf(result).String())
   }
```

## RoadMap
 1. array&&slice&map mapping