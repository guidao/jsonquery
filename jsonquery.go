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

	Set(keys *Lens, value interface{}) Value

	Error() error
}

type Access interface {
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
}

type KeyType int

var (
	ObjectKey  KeyType = 0
	ArrayIndex KeyType = 1
)

type access struct {
	keyType KeyType
	key     interface{}
}

func (this *access) Get(o interface{}) (interface{}, error) {
	if o == nil {
		return nil, errors.New("object is nil")
	}
	switch this.keyType {
	case ObjectKey:
		m, ok := o.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("key is not a map:%v", this.key)
		}
		if v, ok := m[this.key.(string)]; ok {
			return v, nil
		}
		return nil, errors.Errorf("key not found:%v", this.key)
	case ArrayIndex:
		m, ok := o.([]interface{})
		if !ok {
			return nil, errors.Errorf("key is not a array:%v", this.key)
		}
		if len(m) <= this.key.(int) {
			return nil, errors.Errorf("key is out of array:%v", this.key)
		}
		return m[this.key.(int)], nil
	default:
		return nil, errors.Errorf("key type error:%v", this.key)
	}
}

func (this *access) Set(o interface{}, value interface{}) error {
	if o == nil {
		return errors.New("object is nil")
	}
	switch this.keyType {
	case ObjectKey:
		m, ok := o.(map[string]interface{})
		if !ok {
			return errors.Errorf("key is not a map:%v", this.key)
		}
		m[this.key.(string)] = value
		return nil
	case ArrayIndex:
		m, ok := o.([]interface{})
		if !ok {
			return errors.Errorf("key is not a array:%v", this.key)
		}
		m[this.key.(int)] = value
		return nil
	default:
		return errors.Errorf("key type error:%v", this.key)
	}
}

type Lens struct {
	lens []Access
	err  error
}

func NewLens() *Lens {
	return &Lens{}
}

func (this *Lens) Key(k string) *Lens {
	this.lens = append(this.lens, &access{keyType: ObjectKey, key: k})
	return this
}

func (this *Lens) Index(i int) *Lens {
	this.lens = append(this.lens, &access{keyType: ArrayIndex, key: i})
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
		v, err = f.Get(v)
		if err != nil {
			return &jsonValue{err: err}
		}
	}
	return &jsonValue{value: v}
}

func (this *Lens) set(object, value interface{}) Value {
	if this.err != nil {
		return &jsonValue{err: this.err}
	}
	length := len(this.lens)
	o := object
	for i, f := range this.lens {
		if i == length-1 {
			f.Set(o, value)
			return &jsonValue{value: object}
		}
		var err error
		o, err = f.Get(o)
		if err != nil {
			return &jsonValue{err: err}
		}
	}
	return &jsonValue{value: object}
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
		return nil
	}
	return errors.New("not object")
}

func (this *jsonValue) Set(keys *Lens, value interface{}) Value {
	if len(keys.lens) == 0 || keys == nil {
		this.value = value
		return this
	}
	if this.err != nil {
		return this
	}
	return keys.set(this.InterfaceOr(nil), value)
}
