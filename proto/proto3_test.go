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
	"testing"

	"google.golang.org/protobuf/proto"
	"katydid.org.za/go/parser-go-proto/proto/prototests"
	"katydid.org.za/go/parser-go/hedge"
	"katydid.org.za/go/parser-go/parse/debug"
	"katydid.org.za/go/parser-go/rand"
)

var proto3Input1 = &prototests.Proto3{
	Field: 97824789,
	Msg: &prototests.SmallMsg3{
		ScarBusStop:     "cde",
		FlightParachute: []uint32{1, 2, 3},
	},
	Ints: []int64{math.MinInt64},
}

var proto3Output1 = hedge.Hedge{
	{Label: hedge.NewStringToken("Field"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(97824789), Children: nil}}},
	{Label: hedge.NewStringToken("Msg"), Children: hedge.Hedge{
		{Label: hedge.NewStringToken("ScarBusStop"), Children: hedge.Hedge{{Label: hedge.NewStringToken("cde"), Children: nil}}},
		{Label: hedge.NewStringToken("FlightParachute"), Children: hedge.Hedge{
			{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewInt64Token(1), Children: nil}}},
			{Label: hedge.NewInt64Token(1), Children: hedge.Hedge{{Label: hedge.NewInt64Token(2), Children: nil}}},
			{Label: hedge.NewInt64Token(2), Children: hedge.Hedge{{Label: hedge.NewInt64Token(3), Children: nil}}},
		}},
	}},
	{Label: hedge.NewStringToken("Ints"), Children: hedge.Hedge{
		{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewInt64Token(math.MinInt64), Children: nil}}},
	}},
}

func TestProto31(t *testing.T) {
	p, err := NewParser("prototests", "Proto3")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(proto3Input1)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	m, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(proto3Output1) {
		t.Fatalf("expected %s but got %s", proto3Output1, m)
	}
}

func TestRandomProto31(t *testing.T) {
	p, err := NewParser("prototests", "Proto3")
	if err != nil {
		t.Fatal(err)
	}
	data, err := proto.Marshal(proto3Input1)
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
