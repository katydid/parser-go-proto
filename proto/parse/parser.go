//  Copyright 2026 Walter Schulze
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

package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/katydid/parser-go-proto/proto/desc"
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/jsonschema"
	"github.com/katydid/parser-go/parse"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// Parser represents a protocol buffer parser.
type Parser interface {
	parse.Parser
	//Init initialises the parser with a marshaled protocol buffer.
	Init([]byte)
	//Reset resets the parser to go back to the beginnig.
	Reset()
	//Message returns the current message's descriptor.
	Message() *descriptor.DescriptorProto
	//Field returns the current field's descriptor.
	Field() *descriptor.FieldDescriptorProto

	jsonschema.JSONSchemaAble
}

// NewParser returns a new protocol buffer parser the specific root message.
// When the value of a field name is requested this parser will return the field name using the String method.
func NewParser(rootPackage, rootMessage string, opts ...Option) (Parser, error) {
	return newProtoParser(rootPackage, rootMessage, opts...)
}

func newProtoParser(srcPackage, srcMessage string, opts ...Option) (*parser, error) {
	o := newOptions(opts...)
	descMap, err := desc.NewDescriptorMap(srcPackage, srcMessage, o.desc)
	if err != nil {
		return nil, err
	}
	root := descMap.GetRoot()
	fieldsMap := descMap.LookupFields(root)
	return &parser{
		alloc: o.alloc,

		root:          root,
		rootDescMap:   descMap,
		rootFieldsMap: fieldsMap,

		stack: make([]state, 0, 10),
	}, nil
}

type parser struct {
	alloc func(size int) []byte

	root          *descriptor.DescriptorProto
	rootDescMap   desc.DescMap
	rootFieldsMap map[uint64]*descriptor.FieldDescriptorProto

	buf []byte

	keyVarint  uint64
	decodedKey bool

	tokenized   bool
	tokenVarint uint64

	state
	stack []state
}

type state struct {
	// required
	parent    *descriptor.DescriptorProto
	fieldsMap map[uint64]*descriptor.FieldDescriptorProto
	kind      stateKind
	offset    int
	endOffset int
	// optional
	field       *descriptor.FieldDescriptorProto
	wireType    int
	fieldNumber int32
}

func (p *parser) Reset() {
	p.Init(p.buf)
}

func (p *parser) Init(buf []byte) {
	p.stack = p.stack[:0]
	p.buf = buf
	p.state = state{
		parent:    p.root,
		fieldsMap: p.rootFieldsMap,
		offset:    0,
		endOffset: len(buf),
	}
	return
}

type stateKind byte

const startState = stateKind(0)
const inMessageState = stateKind('m')
const atFieldState = stateKind('f')
const firstRepeatedValueState = stateKind('0')
const inRepeatedFieldState = stateKind('[')
const inPackedField = stateKind('+')
const isLeafState = stateKind('l')
const endState = stateKind('$')

func (p *parser) Message() *descriptor.DescriptorProto {
	return p.parent
}

func (p *parser) Field() *descriptor.FieldDescriptorProto {
	return p.field
}

func (p *parser) JSONSchemaType() jsonschema.JSONSchemaType {
	switch p.state.kind {
	case startState:
		return jsonschema.JSONSchemaTypeUnknown
	case inMessageState:
		return jsonschema.JSONSchemaTypeObject
	case atFieldState:
		return jsonschema.JSONSchemaTypeUnknown
	case isLeafState:
		return jsonschema.JSONSchemaTypeUnknown
	case firstRepeatedValueState:
		return jsonschema.JSONSchemaTypeArray
	case inRepeatedFieldState:
		return jsonschema.JSONSchemaTypeUnknown
	case inPackedField:
		return jsonschema.JSONSchemaTypeArray
	case endState:
		return jsonschema.JSONSchemaTypeObject
	}
	panic(fmt.Sprintf("unreachable kind %c", p.state.kind))
}

func (p *parser) Next() (parse.Hint, error) {
	p.tokenized = false
	switch p.state.kind {
	case startState:
		return p.nextStart()
	case inMessageState:
		return p.nextInMessage()
	case atFieldState:
		return p.nextAtField()
	case isLeafState:
		return p.nextIsLeaf()
	case firstRepeatedValueState:
		return p.nextFirstRepeatedValueState()
	case inRepeatedFieldState:
		return p.nextInRepeatedField()
	case inPackedField:
		return p.nextInPackedField()
	case endState:
		return p.nextEnd()
	}
	panic(fmt.Sprintf("unreachable kind %c", p.state.kind))
}

func (p *parser) nextStart() (parse.Hint, error) {
	p.state.kind = endState
	p.down(state{
		parent:    p.parent,
		fieldsMap: p.fieldsMap,
		kind:      inMessageState,
		offset:    p.offset,
		endOffset: p.endOffset,
	})
	return parse.EnterHint, nil
}

func (p *parser) nextInMessage() (parse.Hint, error) {
	if p.offset >= p.endOffset {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.LeaveHint, nil
	}
	v, n, err := uvarint(p.buf[p.offset:])
	if err != nil {
		return parse.UnknownHint, err
	}
	p.offset += n
	p.wireType = int(v & 0x7)
	p.fieldNumber = int32(v >> 3)
	var ok bool
	p.field, ok = p.fieldsMap[v]
	if !ok {
		// skip unknown field
		if err := p.skipFieldValue(); err != nil {
			return parse.UnknownHint, err
		}
		return p.Next()
	}
	p.state.kind = atFieldState
	return parse.FieldHint, nil
}

func (p *parser) skipFieldValue() error {
	length, err := p.decodeLength(p.wireType)
	if err != nil {
		return err
	}
	p.offset += length
	return nil
}

func (p *parser) nextAtField() (parse.Hint, error) {
	if desc.IsRepeated(p.field) {
		if desc.IsScalar(p.field) && p.wireType == 2 { // isPacked
			length, err := p.decodeLength(p.wireType)
			if err != nil {
				return parse.UnknownHint, err
			}
			offset := p.offset
			p.offset += length
			p.state.kind = inMessageState
			p.down(state{
				parent:    p.parent,
				fieldsMap: p.fieldsMap,
				kind:      inPackedField,
				offset:    offset,
				endOffset: p.offset,

				fieldNumber: p.fieldNumber,
				wireType:    desc.WireType(p.field),
				field:       p.field,
			})
			return parse.EnterHint, nil
		} else {
			p.state.kind = inMessageState
			p.down(state{
				parent:    p.parent,
				fieldsMap: p.fieldsMap,
				kind:      firstRepeatedValueState,
				offset:    p.offset,
				// We cannot guess the end offset of all the repeated fields,
				// without decoding all repeated fields,
				// so we use the old end offset.
				endOffset: p.endOffset,

				fieldNumber: p.fieldNumber,
				wireType:    p.wireType,
				field:       p.field,
			})
			return parse.EnterHint, nil
		}
	} else if desc.IsMessage(p.field) {
		length, err := p.decodeLength(p.wireType)
		if err != nil {
			return parse.UnknownHint, err
		}
		offset := p.offset
		p.offset += length
		p.state.kind = inMessageState
		newParent := p.rootDescMap.LookupMessage(p.field)
		newFieldsMap := p.rootDescMap.LookupFields(newParent)
		p.down(state{
			parent:    newParent,
			fieldsMap: newFieldsMap,
			kind:      inMessageState,
			offset:    offset,
			endOffset: p.offset,
		})
		return parse.EnterHint, nil
	} else {
		length, err := p.decodeLength(p.wireType)
		if err != nil {
			return parse.UnknownHint, err
		}
		offset := p.offset
		p.offset += length
		p.state.kind = inMessageState
		p.down(state{
			parent:    p.parent,
			fieldsMap: p.fieldsMap,
			kind:      isLeafState,
			offset:    offset,
			endOffset: p.offset,

			fieldNumber: p.fieldNumber,
			wireType:    p.wireType,
			field:       p.field,
		})
		return parse.ValueHint, nil
	}
}

func (p *parser) nextIsLeaf() (parse.Hint, error) {
	if err := p.up(); err != nil {
		return parse.UnknownHint, err
	}
	return p.Next()
}

func (p *parser) nextFirstRepeatedValueState() (parse.Hint, error) {
	length, err := p.decodeLength(p.wireType)
	if err != nil {
		return parse.UnknownHint, err
	}
	offset := p.offset
	p.offset += length
	if desc.IsMessage(p.field) {
		p.state.kind = inRepeatedFieldState
		newParent := p.rootDescMap.LookupMessage(p.field)
		newFieldsMap := p.rootDescMap.LookupFields(newParent)
		p.down(state{
			parent:    newParent,
			fieldsMap: newFieldsMap,
			kind:      inMessageState,
			offset:    offset,
			endOffset: p.offset,
		})
		return parse.EnterHint, nil
	}
	p.state.kind = inRepeatedFieldState
	p.down(state{
		parent:    p.parent,
		fieldsMap: p.fieldsMap,
		kind:      isLeafState,
		offset:    offset,
		endOffset: p.offset,

		fieldNumber: p.fieldNumber,
		wireType:    p.wireType,
		field:       p.field,
	})

	return parse.ValueHint, nil
}

func (p *parser) nextInRepeatedField() (parse.Hint, error) {
	if p.offset == p.endOffset {
		offset := p.offset
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		// we could not guess the end of the repeated field,
		// without decoding all of the repeated field,
		// so we set it now.
		p.offset = offset
		return parse.LeaveHint, nil
	}
	v, n, err := uvarint(p.buf[p.offset:])
	if err != nil {
		return parse.UnknownHint, err
	}
	fieldNumber := int32(v >> 3)
	if fieldNumber != p.fieldNumber {
		// new wire type means we have reached the end of the repeated field
		offset := p.offset
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		// we could not guess the end of the repeated field,
		// without decoding all of the repeated field,
		// so we set it now.
		p.offset = offset
		return parse.LeaveHint, nil
	}
	p.offset += n
	length, err := p.decodeLength(p.wireType)
	if err != nil {
		return parse.UnknownHint, err
	}
	offset := p.offset
	p.offset += length
	if desc.IsMessage(p.field) {
		newParent := p.rootDescMap.LookupMessage(p.field)
		newFieldsMap := p.rootDescMap.LookupFields(newParent)
		p.down(state{
			parent:    newParent,
			fieldsMap: newFieldsMap,
			kind:      inMessageState,
			offset:    offset,
			endOffset: p.offset,
		})
		return parse.EnterHint, nil
	}
	p.down(state{
		parent:    p.parent,
		fieldsMap: p.fieldsMap,
		kind:      isLeafState,
		offset:    offset,
		endOffset: p.offset,

		wireType: p.wireType,
		field:    p.field,
	})
	return parse.ValueHint, nil
}

func (p *parser) skipRepeatedField() error {
	for {
		if p.offset == p.endOffset {
			offset := p.offset
			if err := p.up(); err != nil {
				return err
			}
			// we could not guess the end of the repeated field,
			// without decoding all of the repeated field,
			// so we set it now.
			p.offset = offset
			return nil
		}
		v, n, err := uvarint(p.buf[p.offset:])
		if err != nil {
			return err
		}
		fieldNumber := int32(v >> 3)
		if fieldNumber != p.fieldNumber {
			// new wire type means we have reached the end of the repeated field
			offset := p.offset
			if err := p.up(); err != nil {
				return err
			}
			// we could not guess the end of the repeated field,
			// without decoding all of the repeated field,
			// so we set it now.
			p.offset = offset
			return nil
		}
		p.offset += n
		if err := p.skipFieldValue(); err != nil {
			return err
		}
	}
}

func (p *parser) nextInPackedField() (parse.Hint, error) {
	if p.offset == p.endOffset {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.LeaveHint, nil
	}
	length, err := p.decodeLength(p.wireType)
	if err != nil {
		return parse.UnknownHint, err
	}
	offset := p.offset
	p.offset += length
	p.down(state{
		parent:    p.parent,
		fieldsMap: p.fieldsMap,
		kind:      isLeafState,
		offset:    offset,
		endOffset: p.offset,

		fieldNumber: p.fieldNumber,
		wireType:    p.wireType,
		field:       p.field,
	})
	return parse.ValueHint, nil
}

func (p *parser) nextEnd() (parse.Hint, error) {
	return parse.UnknownHint, io.EOF
}

func (p *parser) down(state state) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state
}

func (p *parser) up() error {
	if len(p.stack) == 0 {
		return io.ErrUnexpectedEOF
	}
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
	if len(p.stack) == 0 {
		p.state.kind = endState
	}
	return nil
}

func (p *parser) Skip() error {
	switch p.state.kind {
	case startState:
		_, err := p.Next()
		return err
	case inMessageState:
		// we just opened a message or finished parsing a field.
		// The end offset of the message is known,
		// so we can skip the whole message by going up.
		return p.up()
	case atFieldState:
		// we just parsed a field, which we can now skip.
		if err := p.skipFieldValue(); err != nil {
			return err
		}
		p.state.kind = inMessageState
		return nil
	case isLeafState:
		_, err := p.Next()
		return err
	case firstRepeatedValueState:
		if err := p.skipFieldValue(); err != nil {
			return err
		}
		return p.skipRepeatedField()
	case inRepeatedFieldState:
		return p.skipRepeatedField()
	case inPackedField:
		// we just opened a packed repeated field or finished parsing a item.
		// The end offset of the packed field is known,
		// so we can skip the whole repeated field by going up.
		return p.up()
	case endState:
		_, err := p.Next()
		return err
	}
	panic(fmt.Sprintf("unreachable kind %c", p.state.kind))
}

func (p *parser) Token() (parse.Kind, []byte, error) {
	switch p.kind {
	case startState, endState, inMessageState:
		return parse.UnknownKind, nil, nil
	case atFieldState:
		if p.field != nil && p.field.Name != nil {
			s := *p.field.Name
			token := cast.FromString(s, p.alloc)
			return parse.StringKind, token, nil
		}
		return parse.UnknownKind, nil, nil
	case isLeafState:
		typ := p.field.GetType()
		switch typ {
		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			token := p.slice()
			return parse.Float64Kind, token, nil
		case descriptor.FieldDescriptorProto_TYPE_FLOAT:
			f32 := cast.ToFloat32(p.slice())
			token := cast.FromFloat64(float64(f32), p.alloc)
			return parse.Float64Kind, token, nil
		case descriptor.FieldDescriptorProto_TYPE_INT64:
			i64, err := p.decodeInt64()
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
			i64, err := p.decodeSfixed64()
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_SINT64:
			i64, err := p.decodeSint64()
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_INT32:
			i32, err := p.decodeInt32()
			i64 := int64(i32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
			i32, err := p.decodeSfixed32()
			i64 := int64(i32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_SINT32:
			i32, err := p.decodeSint32()
			i64 := int64(i32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			i32, err := p.decodeInt32()
			i64 := int64(i32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_UINT64:
			u64, err := p.decodeUint64()
			if u64 <= math.MaxInt64 {
				i64 := int64(u64)
				token := cast.FromInt64(i64, p.alloc)
				return parse.Int64Kind, token, err
			}
			s := strconv.FormatUint(u64, 10)
			token := cast.FromString(s, p.alloc)
			return parse.DecimalKind, token, err
		case descriptor.FieldDescriptorProto_TYPE_FIXED64:
			u64, err := p.decodeFixed64()
			if u64 <= math.MaxInt64 {
				i64 := int64(u64)
				token := cast.FromInt64(i64, p.alloc)
				return parse.Int64Kind, token, err
			}
			s := strconv.FormatUint(u64, 10)
			token := cast.FromString(s, p.alloc)
			return parse.DecimalKind, token, err
		case descriptor.FieldDescriptorProto_TYPE_UINT32:
			u32, err := p.decodeUint32()
			i64 := int64(u32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_FIXED32:
			u32, err := p.decodeFixed32()
			i64 := int64(u32)
			token := cast.FromInt64(i64, p.alloc)
			return parse.Int64Kind, token, err
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			v, err := p.decodeBool()
			if err != nil {
				return parse.UnknownKind, nil, err
			}
			if v {
				return parse.TrueKind, nil, nil
			}
			return parse.FalseKind, nil, nil
		case descriptor.FieldDescriptorProto_TYPE_STRING:
			buf := p.slice()
			return parse.StringKind, buf, nil
		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			buf := p.slice()
			return parse.BytesKind, buf, nil
		}
		return parse.UnknownKind, nil, errUnknownFieldType
	}
	panic(fmt.Sprintf("unreachable %c", p.kind))
}

func (p *parser) decodeLength(wireType int) (int, error) {
	n, l, err := length(wireType, p.buf[p.offset:])
	if err != nil {
		return 0, err
	}
	p.offset += n
	if p.offset+l > p.endOffset {
		return 0, io.ErrShortBuffer
	}
	return l, nil
}

func (p *parser) decodeBool() (bool, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return v != 0, err
}

func (p *parser) decodeInt64() (int64, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return int64(v), err
}

func (p *parser) decodeUint64() (uint64, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return v, err
}

func (p *parser) decodeInt32() (int32, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return int32(v), err
}

func (p *parser) decodeFixed64() (uint64, error) {
	buf := p.slice()
	if len(buf) < 8 {
		return 0, io.ErrShortBuffer
	}
	return cast.ToUint64(buf), nil
}

func (p *parser) decodeFixed32() (uint32, error) {
	buf := p.slice()
	if len(buf) < 4 {
		return 0, io.ErrShortBuffer
	}
	return cast.ToUint32(buf), nil
}

func (p *parser) decodeUint32() (uint32, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return uint32(v), err
}

func (p *parser) decodeSfixed32() (int32, error) {
	buf := p.slice()
	if len(buf) < 4 {
		return 0, io.ErrShortBuffer
	}
	return cast.ToInt32(buf), nil
}

func (p *parser) decodeSfixed64() (int64, error) {
	buf := p.slice()
	if len(buf) < 8 {
		return 0, io.ErrShortBuffer
	}
	return cast.ToInt64(buf), nil
}

func (p *parser) decodeSint32() (int32, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31)), err
}

func (p *parser) decodeSint64() (int64, error) {
	buf := p.slice()
	v, err := p.tokenizeVarint(buf)
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63)), err
}

func (p *parser) slice() []byte {
	return p.buf[p.offset:p.endOffset]
}

func (p *parser) tokenizeVarint(bs []byte) (uint64, error) {
	if p.tokenized {
		return p.tokenVarint, nil
	}
	v, n := binary.Uvarint(bs)
	if n <= 0 {
		return 0, errDecodeVarint
	}
	p.tokenVarint = v
	p.tokenized = true
	p.endOffset = p.offset + n
	return v, nil
}
