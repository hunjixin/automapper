# automapper [![godoc reference](https://godoc.org/github.com/hunjixin/automapper?status.png)](https://godoc.org/github.com/hunjixin/automapper)  [![Build Status](https://travis-ci.org/hunjixin/automapper.svg?branch=master)](https://travis-ci.org/hunjixin/automapper)
Package automapper provides data mapping between different struct

### Features

1. complex type mapping include embed field,array, slice and map
2. support tag to redefine field name
3. func to customize field mapping content
4. automatic registration when no special requirement


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
        automapper.MustCreateMapper(reflect.TypeOf((*User)(nil)), reflect.TypeOf((*UserDto)(nil))).
            Mapping(func(destVal interface{}, sourceVal interface{}) {
                destVal.(*UserDto).Name = sourceVal.(*User).Name + "|" + sourceVal.(*User).Nick
            }).
            Mapping(func(destVal interface{}, sourceVal interface{}) {
                destVal.(*UserDto).Age = time.Now().Year() - sourceVal.(*User).Birth.Year()
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