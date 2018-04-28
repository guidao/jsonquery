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

	o := jsonquery.NewLens().Key("array").GetWithJson(`{"array": ["hello", "world"]}`)
	o = o.Set(jsonquery.NewLens().Index(0), "HELLO")
	fmt.Println(o.InterfaceOr(nil)) //["HELLO", "world"]
	fmt.Println(err)                //nil
}
