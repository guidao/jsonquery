package jsonquery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var s = `{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`

func TestKey(t *testing.T) {
	v := `{"key": [1,2,3]}`
	value, err := NewLens().Key("key").Index(1).GetWithJson(v).Float64()
	assert.NoError(t, err)
	assert.Equal(t, float64(2), value)
}

func TestSet(t *testing.T) {
	value := NewLens().GetWithJson(s)
	title := NewLens().Key("store").Key("book").Index(1).Key("title")
	assert.Equal(t, "Sword of Honour", title.GetWithJson(s).StringOr(""))
	value = value.Set(title, "hello world")
	assert.Equal(t, "hello world", title.GetWithValue(value.InterfaceOr(nil)).StringOr(""))

	//如果值不存在，可以更新
	name := NewLens().Key("store").Key("bicycle").Key("name")
	value = value.Set(name, "hello")
	assert.NoError(t, value.Error())
	assert.NoError(t, name.GetWithValue(value.InterfaceOr(nil)).Error())
	assert.Equal(t, "hello", name.GetWithValue(value.InterfaceOr(nil)).StringOr(""))
}
