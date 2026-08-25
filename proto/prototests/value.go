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
	hedge.Field(`FieldA`, `0.123`),
	hedge.Nested(`FieldB`,
		hedge.Field(`SmallField`, `456`),
	),
	hedge.Nested(`FieldC`,
		hedge.Field(`BigField`, `789`),
	),
	hedge.Field(`Field1`, `123`),
}

// ABigContainer is a populated BigContainer instance.
var ABigContainer = &BigContainer{
	Field13: proto.Int64(987),
	M:       AContainer,
}

// ABigContainer is a populated BigContainer instance that has been parsed into hedge.Hedge.
var ABigContainerOutput = hedge.Hedge{
	hedge.Nested(`M`, AContainerOutput...),
	hedge.Field(`Field13`, `987`),
}
