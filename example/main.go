package main

import (
	"fmt"
	jq "github.com/guidao/jsonquery"
)

func main() {
	//error
	v := jq.NewLens().Key("inner").GetWithJson(`{}`)
	fmt.Println(v.StringOr("default"), v.Error())

	//object
	hello := jq.NewLens().Key("hello").GetWithJson(`{"hello": "world"}`).StringOr("")
	fmt.Println(hello) // world

	//array
	f23, err := jq.NewLens().Key("array").Index(2).GetWithJson(`{"array":["hello", "world", 23]}`).Float64()
	fmt.Println(f23, err) //23, nil

	//foreach
	err = jq.NewLens().Key("array").GetWithJson(`{"array":["hello", "world", 23]}`).ForeachArray(func(i int, v jq.Value) {
		fmt.Printf("i:%v, v:%v\n", i, v.InterfaceOr(nil))
		// i:0, v:hello
		// i:1, v:world
		// i:2, v:23
	})
	fmt.Println(err) //nil

	o := jq.NewLens().Key("array").GetWithJson(`{"array": ["hello", "world"]}`)
	o = o.Set(jq.NewLens().Index(0), "HELLO")
	fmt.Println(o.InterfaceOr(nil)) //["HELLO", "world"]
}
