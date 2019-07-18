# automapper
Package automapper provides data mapping between different struct

### Features

1. simple filed mapping 
2. support mapping enbeded fields


## Example

```shell
    go get github.com/hunjixin/automapper
```

## Example

```go
  package main
  
  import (
  	"github.com/hunjixin/automapper"
  	"time"
  	"reflect"
  	"fmt"
  )
  
  type PersonModel struct {
  	Name 		string
  	Age  		int
  	Address     string
  	CreateDate  time.Time
  	DeleteDate  time.Time
  	IsDel 		bool
  }
  
  type PersonDto struct {
  	Name 		string
  	Age  		int
  	Address     string
  	CreateDate  time.Time
  	DeleteDate  time.Time
  	IsDel 		bool
  }
  
  func init(){
  	automapper.CreateMapper(reflect.TypeOf((*PersonModel)(nil)), reflect.TypeOf((*PersonDto)(nil)))
  }
  
  func main() {
  	model := PersonModel{}
  	model.Name = "Jimmy"
  	model.Age = 12
  	model.Address = "S·H"
  	model.CreateDate = time.Now()
  	model.IsDel = true
  	result := automapper.MustMapper(model, reflect.TypeOf((*PersonDto)(nil)))
  	fmt.Println(reflect.TypeOf(result).String())
  }
```
## Roadmap

1. support tag to mapping
2. support injected func to mapping
3. add cache to faster type mapping 