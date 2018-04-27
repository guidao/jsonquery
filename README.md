

# jsonquery
一个方便的json值获取库, go的json解析方式要么需要预先定义结构，要么需要处理map[string]interface{}。第一种取值方便，解析麻烦。第二种解析时方便，取值时麻烦。而这个库综合了这两种方式的优点。


```go
func main(){
    v := `{"key": [1,2,3]}`
    value, err := NewLens().Key("key").Index(1).GetWithJson(v).Float64()
    fmt.Println(v) // 2
}
```

更多例子参考test文件




