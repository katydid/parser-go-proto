//  Copyright 2013 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package reflect

import (
	"github.com/katydid/katydid/serialize"
	"io"
	"reflect"
)

type state struct {
	parent        reflect.Value
	typ           reflect.StructField
	value         reflect.Value
	field         int
	maxField      int
	sliceValue    reflect.Value
	sliceIndex    int
	sliceMaxIndex int
	inSlice       bool
	isFieldValue  bool
}

func (this state) Copy() state {
	return state{
		parent:        this.parent,
		typ:           this.typ,
		value:         this.value,
		field:         this.field,
		maxField:      this.maxField,
		sliceValue:    this.sliceValue,
		sliceIndex:    this.sliceIndex,
		sliceMaxIndex: this.sliceMaxIndex,
		inSlice:       this.inSlice,
		isFieldValue:  this.isFieldValue,
	}
}

type parser struct {
	state
	stack []state
}

func (this *parser) Copy() serialize.Parser {
	s := &parser{
		state: this.state.Copy(),
		stack: make([]state, len(this.stack)),
	}
	for i := range this.stack {
		s.stack[i] = this.stack[i].Copy()
	}
	return s
}

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func newState(parent reflect.Value) state {
	value := deref(parent)
	return state{
		parent:   value,
		maxField: value.NumField(),
	}
}

func newFieldState(typ reflect.StructField, val reflect.Value) state {
	return state{
		typ:          typ,
		value:        val,
		isFieldValue: true,
		maxField:     1,
	}
}

func NewReflectParser() *parser {
	return &parser{stack: make([]state, 0, 10)}
}

func (s *parser) Init(value reflect.Value) *parser {
	s.state = newState(value)
	return s
}

func isSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8
}

func (s *parser) Next() error {
	if s.field >= s.maxField {
		return io.EOF
	}
	if !s.isFieldValue {
		if s.inSlice {
			if s.sliceIndex >= s.sliceMaxIndex {
				s.inSlice = false
				s.field++
				return s.Next()
			}
			s.value = s.sliceValue.Index(s.sliceIndex)
			s.sliceIndex++
			return nil
		}
		s.typ = s.parent.Type().Field(s.field)
		s.value = s.parent.Field(s.field)
		if s.value.Kind() == reflect.Ptr || s.value.Kind() == reflect.Slice {
			if s.value.IsNil() {
				s.field++
				return s.Next()
			}
		}
		if isSlice(s.value) {
			s.sliceValue = s.value
			s.sliceMaxIndex = s.value.Len()
			s.sliceIndex = 0
			s.inSlice = true
			return s.Next()
		}
	}
	s.field++
	return nil
}

func (s *parser) IsLeaf() bool {
	return s.isFieldValue
}

func (s *parser) getValue() reflect.Value {
	kind := s.value.Kind()
	if kind == reflect.Slice {
		childValue := s.value.Elem()
		childKind := childValue.Kind()
		switch childKind {
		case reflect.Uint8, reflect.Int8:
			return s.value
		case reflect.Slice:
			switch childValue.Elem().Kind() {
			case reflect.Uint8, reflect.Int8:
				return childValue
			}
		default:
			return childValue
		}
	}
	if kind == reflect.Ptr {
		return s.value.Elem()
	}
	return s.value
}

func (s *parser) Double() (float64, error) {
	if s.isFieldValue {
		value := s.getValue()
		switch value.Kind() {
		case reflect.Float64, reflect.Float32:
			return value.Float(), nil
		}
	}
	return 0, serialize.ErrNotDouble
}

func (s *parser) Int() (int64, error) {
	if s.isFieldValue {
		value := s.getValue()
		switch value.Kind() {
		case reflect.Int64, reflect.Int32:
			return value.Int(), nil
		}
	}
	return 0, serialize.ErrNotInt
}

func (s *parser) Uint() (uint64, error) {
	if s.isFieldValue {
		value := s.getValue()
		switch value.Kind() {
		case reflect.Uint64, reflect.Uint32:
			return value.Uint(), nil
		}
	}
	return 0, serialize.ErrNotUint
}

func (s *parser) Bool() (bool, error) {
	if s.isFieldValue {
		value := s.getValue()
		switch value.Kind() {
		case reflect.Bool:
			return value.Bool(), nil
		}
	}
	return false, serialize.ErrNotBool
}

func (s *parser) String() (string, error) {
	if !s.isFieldValue {
		return s.typ.Name, nil
	}
	value := s.getValue()
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	}
	return "", serialize.ErrNotString
}

func (s *parser) Bytes() ([]byte, error) {
	if s.isFieldValue {
		value := s.getValue()
		switch value.Kind() {
		case reflect.Slice, reflect.Uint8, reflect.Int8:
			return value.Bytes(), nil
		}
	}
	return nil, serialize.ErrNotBytes
}

func (s *parser) Up() {
	top := len(s.stack) - 1
	s.state = s.stack[top]
	s.stack = s.stack[:top]
}

func (s *parser) canDown() bool {
	if s.typ.Type.Kind() == reflect.Struct {
		return true
	}
	if s.typ.Type.Kind() == reflect.Ptr {
		if s.typ.Type.Elem().Kind() == reflect.Struct {
			return true
		}
	}
	if s.typ.Type.Kind() == reflect.Slice {
		if s.typ.Type.Elem().Kind() == reflect.Struct {
			return true
		}
		if s.typ.Type.Elem().Kind() == reflect.Ptr {
			if s.typ.Type.Elem().Elem().Kind() == reflect.Struct {
				return true
			}
		}
	}
	return false
}

func (s *parser) Down() {
	s.stack = append(s.stack, s.state)
	if s.canDown() {
		s.state = newState(s.state.value)
	} else {
		s.state = newFieldState(s.typ, s.value)
	}
}
