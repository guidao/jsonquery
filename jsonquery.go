package jsonquery

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type JsonQuery struct {
	data interface{}
}

type Value interface {
	Bool() (bool, error)
	BoolOr(bool) bool

	String() (string, error)
	StringOr(string) string

	Float64() (float64, error)
	Float64Or(float64) float64

	Interface() (interface{}, error)
	InterfaceOr(interface{}) interface{}

	ForeachMap(func(string, Value)) error
	ForeachArray(func(int, Value)) error

	Error() error
}

type Lens struct {
	lens []func(interface{}) (interface{}, error)
	data interface{}
	err  error
}

func NewLens() *Lens {
	return &Lens{}
}

func (this *Lens) Key(k string) *Lens {
	this.lens = append(this.lens, func(input interface{}) (interface{}, error) {
		if input == nil {
			return nil, errors.Errorf("not found key:%v", k)
		}
		m, ok := input.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("not a object:%v", input)
		}
		if v, ok := m[k]; ok {
			return v, nil
		}
		return nil, errors.Errorf("not found key:%v", k)
	})
	return this
}

func (this *Lens) Index(i int) *Lens {
	this.lens = append(this.lens, func(input interface{}) (interface{}, error) {
		if input == nil {
			return nil, errors.Errorf("not found index:%v", i)
		}
		m, ok := input.([]interface{})
		if !ok {
			return nil, errors.Errorf("not a array:%v", input)
		}
		if len(m) <= i {
			return nil, errors.Errorf("array len less index:%v, len:%v", i, len(m))
		}
		return m[i], nil
	})
	return this
}

func (this *Lens) GetWithJson(doc string) Value {
	var v interface{}
	err := json.Unmarshal([]byte(doc), &v)
	if err != nil {
		return &jsonValue{err: err}
	}
	return this.GetWithValue(v)
}

func (this *Lens) GetWithValue(data interface{}) Value {
	if this.err != nil {
		return &jsonValue{err: this.err}
	}
	v := data
	for _, f := range this.lens {
		var err error
		v, err = f(v)
		if err != nil {
			return &jsonValue{err: err}
		}
	}
	return &jsonValue{value: v}
}

type jsonValue struct {
	value interface{}
	err   error
}

func (this *jsonValue) Interface() (interface{}, error) {
	if this.err != nil {
		return nil, this.err
	}
	return this.value, nil
}

func (this *jsonValue) InterfaceOr(def interface{}) interface{} {
	if v, err := this.Interface(); err == nil {
		return v
	}
	return def
}

func (this *jsonValue) Error() error {
	return this.err
}

func (this *jsonValue) StringOr(def string) string {
	if v, err := this.String(); err == nil {
		return v
	}
	return def
}

func (this *jsonValue) String() (string, error) {
	if this.err != nil {
		return "", this.err
	}
	if s, ok := this.value.(string); ok {
		return s, nil
	}
	return "", errors.New("value not string")
}

func (this *jsonValue) Float64Or(def float64) float64 {
	if v, err := this.Float64(); err == nil {
		return v
	}
	return def
}

func (this *jsonValue) Float64() (float64, error) {
	if this.err != nil {
		return 0, this.err
	}
	if n, ok := this.value.(float64); ok {
		return n, nil
	}
	return 0, errors.New("value not float64")
}

func (this *jsonValue) BoolOr(def bool) bool {
	if v, err := this.Bool(); err == nil {
		return v
	}
	return def
}
func (this *jsonValue) Bool() (bool, error) {
	if this.err != nil {
		return false, this.err
	}
	if b, ok := this.value.(bool); ok {
		return b, nil
	}
	return false, errors.New("value not bool")
}

func (this *jsonValue) Unmarshal(v interface{}) error {
	if this.err != nil {
		return this.err
	}
	data, err := json.Marshal(this.value)
	if err != nil {
		return errors.Cause(err)
	}
	if err = json.Unmarshal(data, v); err != nil {
		return errors.Cause(err)
	}
	return nil
}

func (this *jsonValue) ForeachMap(fn func(k string, v Value)) error {
	if this.err != nil {
		return errors.Cause(this.err)
	}
	if m, ok := this.value.(map[string]interface{}); ok {
		for key, value := range m {
			fn(key, &jsonValue{value: value})
		}
	}
	return errors.New("not object")
}

func (this *jsonValue) ForeachArray(fn func(i int, v Value)) error {
	if this.err != nil {
		return errors.Cause(this.err)
	}
	if m, ok := this.value.([]interface{}); ok {
		for i, value := range m {
			fn(i, &jsonValue{value: value})
		}
	}
	return errors.New("not object")
}
