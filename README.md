# automapper [![godoc reference](https://godoc.org/github.com/hunjixin/automapper?status.png)](https://godoc.org/github.com/hunjixin/automapper)  [![Build Status](https://travis-ci.org/hunjixin/automapper.svg?branch=master)](https://travis-ci.org/hunjixin/automapper)
Package automapper provides data mapping between different struct

### Features

1. Complex type mapping include Struct, Embed Field, Array, Slice, Map and even Multiple pointers
2. Support tag to redefine field name
3. Func to customize struct mapping or global simple type conversion
4. Automatic registration 
5. Easy-to-use API


## Import

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
        a := ExampleStructA{EnB{}, En{"Sh", "Bj"}, "XXXXXX"}
        result2 := automapper.MustMapper(a, reflect.TypeOf(ExampleStructB{}))
        fmt.Println(result2)
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
        Name string `mapping:"Nick"`
        Nick string `mapping:"Name"`
    }
    
    func main() {
        user := &User{"NAME", "NICK"}
        result := automapper.MustMapper(user, reflect.TypeOf((*UserDto)(nil)))
        fmt.Println(result)
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
        Name  string
        Nick  string
        Addr  string
        Birth time.Time
     }
     
     func init() {
        automapper.MustCreateMapper((*User)(nil), (*UserDto)(nil)).
            Mapping(func(destVal reflect.Value, sourceVal interface{}) error {
                destVal.Interface().(*UserDto).Name = sourceVal.(*User).Name + "|" + sourceVal.(*User).Nick
                return nil
            }).
            Mapping(func(destVal reflect.Value, sourceVal interface{}) error {
                destVal.Interface().(*UserDto).Age = time.Now().Year() - sourceVal.(*User).Birth.Year()
                return nil
            })
     }
     
     func main() {
        user := &User{"NAME", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}
        result := automapper.MustMapper(user, reflect.TypeOf(UserDto{}))
        fmt.Println(result)
     }
```

### Array mapping
Array and Slice can map to each other
```go
    package main
    
    import (
        "github.com/hunjixin/automapper"
        "reflect"
        "fmt"
        "time"
    )
    
    type UserDto struct {
        Name string
        Addr string
        Age  int
    }
    
    type User struct {
        Name  string
        Nick  string
        Addr  string
        Birth time.Time
    }
    
    func main(){
        users := [3]*User{}
        users[0] = &User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}
        users[2] = &User{"Jack", "neo", "W·S", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}
        result2 := automapper.MustMapper(users, reflect.TypeOf([]*UserDto{}))
        fmt.Println(result2)
    }
```

### map mapping
map -> map
map[string]interface{} -> struct
struct -> map[string]interface{}
```go
    package main
    
    import (
        "fmt"
        "github.com/hunjixin/automapper"
        "time"
        "reflect"
    )
    type UserDto struct {
        Name string
        Addr string
        Age  int
    }
    
    type User struct {
        Name  string
        Nick  string
        Addr  string
        Birth time.Time
    }
    
    func main(){
        //map => map
        map1 := map[string]*User{
            "Hellen":&User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)},
            "Jack":&User{"Jack", "neo", "W·S", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)},
        }
    
        newVal := automapper.MustMapper(map1, reflect.TypeOf(map[string]*UserDto{}))
        for key, val := range newVal.(map[string]*UserDto) {
            fmt.Print(key)
            fmt.Println(val)
        }
    
        //map => struct
        map2 := map[string]interface{}{
            "Name":  "Hellen",
            "Nick":  "NICK",
            "Addr":  "B·J",
            "Birth": time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC),
        }
        newVal = automapper.MustMapper(map2, reflect.TypeOf(User{}))
        fmt.Println(newVal)
    
        //struct => map
        newVal = automapper.MustMapper(&User{"Hellen", "NICK", "B·J", time.Date(1992, 10, 3, 1, 0, 0, 0, time.UTC)}, reflect.TypeOf(map[string]interface{}{}))
        fmt.Println(newVal)
    }
```

### Define global types conversion and than overwrite in complex type mapping
```go
    package main
    
    import (
        "fmt"
        "github.com/hunjixin/automapper"
        "reflect"
        "strconv"
        "time"
    )
    
    func main(){
        automapper.MustCreateMapper(time.Time{}, "").
            Mapping(func(destVal reflect.Value, sourceVal interface{}) {
                str := sourceVal.(time.Time).String()
                destVal.Elem().SetString(str)
            })
        automapper.MustCreateMapper(0, "").
            Mapping(func(destVal reflect.Value, sourceVal interface{}) {
                intVal := sourceVal.(int)
                destVal.Elem().SetString(strconv.Itoa(intVal))
            })
        type A struct {
            M time.Time
            N int
        }
    
        type B struct {
            M string
            N string
        }
    
        str := automapper.MustMapper(A{time.Now(),123}, reflect.TypeOf(B{}))
        fmt.Println(str)
    
        mapping := automapper.EnsureMapping(A{}, B{})
        mapping.Mapping(func(destVal reflect.Value, sourceVal interface{}) {
            str := sourceVal.(A).M.String()
            destVal.Interface().(*B).M = "北京时间："+str
            destVal.Interface().(*B).N = "到达次数："+ strconv.Itoa(sourceVal.(A).N)
        })
        str = automapper.MustMapper(A{time.Now(),456}, reflect.TypeOf(B{}))
        fmt.Println(str)
    }

```
