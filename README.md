

# jsonquery
一个方便的json值获取库, go的json解析方式要么需要预先定义结构，要么需要处理map[string]interface{}。第一种取值方便，解析麻烦。第二种解析时方便，取值时麻烦。而这个库综合了这两种方式的优点。


```go
package main

import (
	"fmt"
	"github.com/guidao/jsonquery"
)

func main() {
	//error
	v := jsonquery.NewLens().Key("inner").GetWithJson(`{}`)
	fmt.Println(v.StringOr("default"), v.Error())

	//object
	hello := jsonquery.NewLens().Key("hello").GetWithJson(`{"hello": "world"}`).StringOr("")
	fmt.Println(hello) // world

	//array
	f23, err := jsonquery.NewLens().Key("array").Index(2).GetWithJson(`{"array":["hello", "world", 23]}`).Float64()
	fmt.Println(f23, err) //23, nil

	//foreach
	err = jsonquery.NewLens().Key("array").GetWithJson(`{"array":["hello", "world", 23]}`).ForeachArray(func(i int, v jsonquery.Value) {
		fmt.Printf("i:%v, v:%v\n", i, v.InterfaceOr(nil))
		// i:0, v:hello
		// i:1, v:world
		// i:2, v:23
	})
	fmt.Println(err) //nil
}

```

更多例子参考test文件




