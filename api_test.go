package json_query

import (
	"testing"
)

func TestKey(t *testing.T) {
	v := `{"key": [1,2,3]}`
	value, err := NewLens().Key("key").Index(1).GetWithJson(v).Float64()
	if err != nil {
		t.Error(err)
	}
	if value != 2 {
		t.Error("expect 2 have ", value)
	}
	t.Log("vvvv:", value)

}
