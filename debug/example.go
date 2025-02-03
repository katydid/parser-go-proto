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

	"github.com/katydid/parser-go/parser/debug"
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
var Output = debug.Nodes{
	debug.Field(`A`, `1`),
	debug.Nested(`B`,
		debug.Field(`0`, `b2`),
		debug.Field(`1`, `b3`),
	),
	debug.Nested(`C`,
		debug.Field(`A`, `2`),
		debug.Field(`D`, `3`),
		debug.Nested(`E`,
			debug.Nested(`0`,
				debug.Nested(`B`,
					debug.Field(`0`, `b4`),
				),
			),
			debug.Nested(`1`,
				debug.Nested(`B`,
					debug.Field(`0`, `b5`),
				),
			),
		),
	),
	debug.Field(`D`, `4`),
	debug.Nested(`F`,
		debug.Field(`0`, `5`),
	),
}
