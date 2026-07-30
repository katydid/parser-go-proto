//  Copyright 2015 Walter Schulze
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

package proto

import (
	"testing"

	protodebug "github.com/katydid/parser-go-proto/debug"
	"github.com/katydid/parser-go-proto/proto/prototests"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/hedge"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/parse/debug"
	"github.com/katydid/parser-go/rand"
	"google.golang.org/protobuf/proto"
)

func TestDebug(t *testing.T) {
	p, err := NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	p.(parse.ParserWithInit).Init(data)
	m, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(protodebug.Output) {
		t.Fatalf("expected %s but got %s", protodebug.Output, m)
	}
}

func TestRandomDebug(t *testing.T) {
	p, err := NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		p.Init(data)
		if err := debug.RandomWalk(p, rand.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSkipRepeated1(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "C")
}

func TestSkipRepeated2(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
}

func TestIndexIsNotAString(t *testing.T) {
	var p parse.ParserWithInit
	var err error
	p, err = NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
}

func TestExtensionsSmallContainer(t *testing.T) {
	p, err := NewParser("prototests", "Container")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(prototests.AContainer)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	nodes, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if !nodes.Equal(prototests.AContainerOutput) {
		t.Fatalf("expected %v, but got %v", prototests.AContainerOutput, nodes)
	}
}

func TestExtensionsBigContainer(t *testing.T) {
	p, err := NewParser("prototests", "BigContainer")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(prototests.ABigContainer)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	nodes, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if !nodes.Equal(prototests.ABigContainerOutput) {
		t.Fatalf("expected %v, but got %v", prototests.ABigContainerOutput, nodes)
	}
}
