//  Copyright 2017 Walter Schulze
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
	"math"
	"strconv"
	"testing"

	"github.com/katydid/parser-go-proto/proto/prototests"
	"github.com/katydid/parser-go/parser/debug"
	"google.golang.org/protobuf/proto"
)

var packedInput1 = &prototests.Packed{
	Ints: []int64{1, math.MaxInt64, math.MinInt64},
}

var packedOutput1 = debug.Nodes{
	debug.Nested(`Ints`,
		debug.Field(`0`, `1`),
		debug.Field(`1`, strconv.Itoa(math.MaxInt64)),
		debug.Field(`2`, strconv.Itoa(math.MinInt64)),
	),
}

func TestPacked1(t *testing.T) {
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput1)
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
	if !m.Equal(packedOutput1) {
		t.Fatalf("expected %s but got %s", packedOutput1, m)
	}
}

func TestRandomPacked1(t *testing.T) {
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput1)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if err := p.Init(data); err != nil {
			t.Fatal(err)
		}
		l := debug.NewLogger(p, debug.NewLineLogger())
		if err := debug.RandomWalk(l, debug.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
	}
}

var packedInput2 = &prototests.Packed{
	Ints:   []int64{1, math.MaxInt64, math.MinInt64},
	Floats: []float64{0.1},
	Uints:  []uint32{3, 4},
}

var packedOutput2 = debug.Nodes{
	debug.Nested(`Ints`,
		debug.Field(`0`, `1`),
		debug.Field(`1`, strconv.Itoa(math.MaxInt64)),
		debug.Field(`2`, strconv.Itoa(math.MinInt64)),
	),
	debug.Nested(`Floats`,
		debug.Field(`0`, `0.1`),
	),
	debug.Nested(`Uints`,
		debug.Field(`0`, `3`),
		debug.Field(`1`, `4`),
	),
}

func TestPacked2(t *testing.T) {
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput2)
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
	if !m.Equal(packedOutput2) {
		t.Fatalf("expected %s but got %s", packedOutput2, m)
	}
}

func TestRandomPacked2(t *testing.T) {
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput2)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if err := p.Init(data); err != nil {
			t.Fatal(err)
		}
		l := debug.NewLogger(p, debug.NewLineLogger())
		if err := debug.RandomWalk(l, debug.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
	}
}
