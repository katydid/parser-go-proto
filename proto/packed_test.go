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

	"google.golang.org/protobuf/proto"
	"katydid.org.za/go/parser-go-proto/proto/prototests"
	"katydid.org.za/go/parser-go/expect"
	"katydid.org.za/go/parser-go/hedge"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/parser-go/parse/debug"
	"katydid.org.za/go/parser-go/rand"
)

var packedInput1 = &prototests.Packed{
	Ints: []int64{1, math.MaxInt64, math.MinInt64},
}

var packedOutput1 = hedge.Hedge{
	hedge.Nested(`Ints`,
		hedge.Field(`0`, `1`),
		hedge.Field(`1`, strconv.Itoa(math.MaxInt64)),
		hedge.Field(`2`, strconv.Itoa(math.MinInt64)),
	),
}

func TestPacked1ParseInto(t *testing.T) {
	p, err := NewParser("prototests", "Packed")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(packedInput1)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	m, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(packedOutput1) {
		t.Fatalf("expected %s but got %s", packedOutput1, m)
	}
}

func TestSkipPackedField(t *testing.T) {
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
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipPackedItemValues(t *testing.T) {
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
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 2)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipPackedItems(t *testing.T) {
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
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
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
		p.Init(data)
		if err := debug.RandomWalk(p, rand.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
	}
}

var packedInput2 = &prototests.Packed{
	Ints:   []int64{1, math.MaxInt64, math.MinInt64},
	Floats: []float64{0.1},
	Uints:  []uint32{3, 4},
}

var packedOutput2 = hedge.Hedge{
	hedge.Nested(`Ints`,
		hedge.Field(`0`, `1`),
		hedge.Field(`1`, strconv.Itoa(math.MaxInt64)),
		hedge.Field(`2`, strconv.Itoa(math.MinInt64)),
	),
	hedge.Nested(`Floats`,
		hedge.Field(`0`, `0.1`),
	),
	hedge.Nested(`Uints`,
		hedge.Field(`0`, `3`),
		hedge.Field(`1`, `4`),
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
	p.Init(data)
	m, err := hedge.ParseInto(p)
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
		p.Init(data)
		if err := debug.RandomWalk(p, rand.NewRand(), 10, 3); err != nil {
			t.Fatal(err)
		}
	}
}
