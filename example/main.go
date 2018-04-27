package main

import (
	"fmt"
	"github.com/guidao/jsonquery"
)

func main() {
	v := jsonquery.NewLens().Key("inner").GetWithJson(`{}`)
	fmt.Println(v.StringOr("default"), v.Error())
}
