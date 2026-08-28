//  Copyright 2016 Walter Schulze
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

package prototests

import (
	"google.golang.org/protobuf/proto"

	"katydid.org.za/go/parser-go/hedge"
)

// AContainer is a populated Container instance.
var AContainer = &Container{
	Field1: proto.Int64(123),
}

func init() {
	f := float64(0.123)
	proto.SetExtension(AContainer, E_FieldA, f)
	proto.SetExtension(AContainer, E_FieldB, &Small{SmallField: proto.Int64(456)})
	proto.SetExtension(AContainer, E_FieldC, &Big{BigField: proto.Int64(789)})
}

// AContainerOutput is a populated Container instance that has been parsed into hedge.Hedge.
var AContainerOutput = hedge.Hedge{
	{Label: hedge.NewStringToken("FieldA"), Children: hedge.Hedge{{Label: hedge.NewFloat64Token(0.123), Children: nil}}},
	{Label: hedge.NewStringToken("FieldB"), Children: hedge.Hedge{
		{Label: hedge.NewStringToken("SmallField"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(456), Children: nil}}},
	}},
	{Label: hedge.NewStringToken("FieldC"), Children: hedge.Hedge{
		{Label: hedge.NewStringToken("BigField"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(789), Children: nil}}},
	}},
	{Label: hedge.NewStringToken("Field1"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(123), Children: nil}}},
}

// ABigContainer is a populated BigContainer instance.
var ABigContainer = &BigContainer{
	Field13: proto.Int64(987),
	M:       AContainer,
}

// ABigContainer is a populated BigContainer instance that has been parsed into hedge.Hedge.
var ABigContainerOutput = hedge.Hedge{
	{Label: hedge.NewStringToken("M"), Children: AContainerOutput},
	{Label: hedge.NewStringToken("Field13"), Children: hedge.Hedge{{Label: hedge.NewInt64Token(987), Children: nil}}},
}
