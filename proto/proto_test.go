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
	"io"
	"testing"

	protodebug "github.com/katydid/parser-go-proto/debug"
	"github.com/katydid/parser-go-proto/proto/prototests"
	"github.com/katydid/parser-go/parser"
	"github.com/katydid/parser-go/parser/debug"
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
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	m, err := debug.Parse(parser)
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
		if err := p.Init(data); err != nil {
			t.Fatal(err)
		}
		//l := debug.NewLogger(p, debug.NewLineLogger())
		if err := debug.RandomWalk(p, debug.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
		//t.Logf("original %v vs random %v", debug.Output, m)
	}
}

func next(t *testing.T, parser parser.Interface) {
	if err := parser.Next(); err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}
}

func TestSkipRepeated1(t *testing.T) {
	p, err := NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	next(t, parser)
	next(t, parser)
	parser.Down()
	next(t, parser)
	parser.Down()
	parser.Up()
	next(t, parser)
	parser.Down()
	next(t, parser)
	next(t, parser)
	parser.Up()
	parser.Up()
	next(t, parser)
}

func TestSkipRepeated2(t *testing.T) {
	p, err := NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	next(t, parser)
	if _, err := parser.String(); err != nil {
		t.Fatal(err)
	}
	next(t, parser)
	if _, err := parser.String(); err != nil {
		t.Fatal(err)
	}
	parser.Down()
	next(t, parser)
	if _, err := parser.Int(); err != nil {
		t.Fatal(err)
	}
	parser.Up()
	next(t, parser)
}

func TestIndexIsNotAString(t *testing.T) {
	p, err := NewParser("debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	next(t, parser)
	if _, err := parser.String(); err != nil {
		t.Fatal(err)
	}
	next(t, parser)
	if _, err := parser.String(); err != nil {
		t.Fatal(err)
	}
	parser.Down()
	next(t, parser)
	if _, err := parser.String(); err == nil {
		t.Fatal("expected error, since an index is not a string")
	}
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
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	nodes, err := debug.Parse(p)
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
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	nodes, err := debug.Parse(p)
	if err != nil {
		t.Fatal(err)
	}
	if !nodes.Equal(prototests.ABigContainerOutput) {
		t.Fatalf("expected %v, but got %v", prototests.ABigContainerOutput, nodes)
	}
}

func TestDebugWithDesc(t *testing.T) {
	p, err := NewParserWithDesc("debug", "Debug", NewFileDescriptorSet())
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	m, err := debug.Parse(parser)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(protodebug.Output) {
		t.Fatalf("expected %s but got %s", protodebug.Output, m)
	}
}

func TestDebugWithSpecificDesc(t *testing.T) {
	p, err := NewParserWithDesc("debug", "Debug", NewFileDescriptorSet(protodebug.File_debug_proto))
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(protodebug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	parser := debug.NewLogger(p, debug.NewLineLogger())
	m, err := debug.Parse(parser)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(protodebug.Output) {
		t.Fatalf("expected %s but got %s", protodebug.Output, m)
	}
}
