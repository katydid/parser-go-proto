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

package auto_test

import (
	"github.com/katydid/katydid/parser"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/auto"
	"testing"
)

type reset interface {
	parser.Interface
	Reset() error
}

func bench(b *testing.B, grammar *ast.Grammar, gen func() parser.Interface) {
	num := 1000
	parsers := make([]reset, num)
	for i := 0; i < num; i++ {
		parsers[i] = gen().(reset)
	}
	a, err := auto.Compile(grammar)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := parsers[i%num].Reset(); err != nil {
			b.Fatal(err)
		}
		a.Validate(parsers[i%num])
	}
	b.ReportAllocs()
}
