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

package debug

import (
	proto "google.golang.org/protobuf/proto"

	"katydid.org.za/go/parser-go/hedge"
)

// Input is a sample instance of the Debug struct.
var Input = &Debug{
	A: proto.Int64(1),
	B: []string{"b2", "b3"},
	C: &Debug{
		A: proto.Int64(2),
		D: proto.Int32(3),
		E: []*Debug{
			{
				B: []string{"b4"},
			},
			{
				B: []string{"b5"},
			},
		},
	},
	D: proto.Int32(4),
	F: []uint32{5},
}

// Output is a sample instance of Nodes that repesents the Input variable after it has been parsed by Walk.
var Output = hedge.Hedge{
	{Label: hedge.NewStringToken("A"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(1), Children: nil}}},
	{Label: hedge.NewStringToken("B"), Children: hedge.Hedge{
		{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b2"), Children: nil}}},
		{Label: hedge.NewInt64Token(1), Children: hedge.Hedge{{Label: hedge.NewStringToken("b3"), Children: nil}}},
	}},
	{Label: hedge.NewStringToken("C"), Children: hedge.Hedge{
		{Label: hedge.NewStringToken("A"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(2), Children: nil}}},
		{Label: hedge.NewStringToken("D"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(3), Children: nil}}},
		{Label: hedge.NewStringToken("E"), Children: hedge.Hedge{
			{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{
				{Label: hedge.NewStringToken("B"), Children: hedge.Hedge{
					{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b4"), Children: nil}}},
				}},
			}},
			{Label: hedge.NewInt64Token(1), Children: hedge.Hedge{
				{Label: hedge.NewStringToken("B"), Children: hedge.Hedge{
					{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b5"), Children: nil}}},
				}},
			}},
		}},
	}},
	{Label: hedge.NewStringToken("D"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(4), Children: nil}}},
	{Label: hedge.NewStringToken("F"), Children: hedge.Hedge{
		{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewInt64Token(5), Children: nil}}},
	}},
}
