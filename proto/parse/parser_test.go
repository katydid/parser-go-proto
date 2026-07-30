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
	"math"
	"testing"

	protodebug "github.com/katydid/parser-go-proto/debug"
	"github.com/katydid/parser-go-proto/proto/desc"
	"github.com/katydid/parser-go-proto/proto/prototests"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
	"google.golang.org/protobuf/proto"
)

func TestSingleDebugAInt(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		A: proto.Int64(1),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSingleDebugBRepeated1String(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		B: []string{"a"},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSingleDebugBRepeated2String(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		B: []string{"a", "b"},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestDebugBRepeatedStringThenAnotherField(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		B: []string{"a", "b"},
		D: proto.Int32(123),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "D")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 123)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestDebugBRepeatedStringThenAnotherSkip(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		B: []string{"a", "b"},
		D: proto.Int32(123),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "D")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSingleDebugCAMessage(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		C: &protodebug.Debug{
			A: proto.Int64(123),
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "C")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 123)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSingleDebugERepeatedMessage(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		E: []*protodebug.Debug{
			{
				A: proto.Int64(123),
			},
			{
				A: proto.Int64(456),
			},
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "E")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 123)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 456)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestPacked1(t *testing.T) {
	var packedInput1 = &prototests.Packed{
		Ints: []int64{1, math.MaxInt64, math.MinInt64},
	}
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput1)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Ints")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, math.MaxInt64)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, math.MinInt64)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipPacked1(t *testing.T) {
	var packedInput1 = &prototests.Packed{
		Ints: []int64{1, math.MaxInt64, math.MinInt64},
	}
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput1)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Ints")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestDebugWithDesc(t *testing.T) {
	p, err := NewParser("debug", "Debug", WithFileDescriptorSet(desc.NewFileDescriptorSet()))
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		C: &protodebug.Debug{
			A: proto.Int64(123),
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "C")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 123)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestDebugWithSpecificDesc(t *testing.T) {
	p, err := NewParser("debug", "Debug", WithFileDescriptorSet(desc.NewFileDescriptorSet(protodebug.File_debug_proto)))
	if err != nil {
		t.Fatal(err)
	}
	msg := &protodebug.Debug{
		C: &protodebug.Debug{
			A: proto.Int64(123),
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "C")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 123)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
