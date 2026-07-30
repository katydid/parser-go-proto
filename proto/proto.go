//  Copyright 2013 Walter Schulze
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

// Package proto contains the implementation of a Proto parser.
package proto

import (
	"github.com/katydid/parser-go-proto/proto/parse"
	goparse "github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/pool"
	"github.com/katydid/parser-go/tag"
)

type Parser interface {
	goparse.Parser
	Reset()
	// Init restarts the parser with a new byte buffer, without allocating a new parser.
	Init([]byte)
}

type parserWithReset interface {
	goparse.Parser
	Reset()
}

type parser struct {
	parserWithReset
	underlying Parser
	pool       pool.Pool
}

// NewParser returns a new Proto parser with indexes.
// Use this parser with other Katydid tools, such as the validator.
func NewParser(rootPackage, rootMessage string, opts ...Option) (Parser, error) {
	o := newOptions(opts...)
	p := pool.New()
	pos := []parse.Option{parse.WithAllocator(p.Alloc)}
	if o.desc != nil {
		pos = append(pos, parse.WithFileDescriptorSet(o.desc))
	}
	underlyingParser, err := parse.NewParser(rootPackage, rootMessage, pos...)
	if err != nil {
		return nil, err
	}
	tagged := tag.NewTagger(underlyingParser, tag.WithAllocator(p.Alloc), tag.WithIndexes())
	return &parser{parserWithReset: tagged, underlying: underlyingParser, pool: p}, nil
}

// NewJSONSchemaParser returns a new Proto parser that tags objects and arrays, so that the types can be checked by JSONSchema.
// The following json: `{"a": ["b", "c"]}`
// is parsed as: `{"object": {"a": {"array": {0: "b", 1: "c"}}}}`.
// The kind returned from the Token method for "object" and "array" will be parse.TagKind.
func NewJSONSchemaParser(rootPackage, rootMessage string, opts ...Option) (Parser, error) {
	o := newOptions(opts...)
	p := pool.New()
	pos := []parse.Option{parse.WithAllocator(p.Alloc)}
	if o.desc != nil {
		pos = append(pos, parse.WithFileDescriptorSet(o.desc))
	}
	underlyingParser, err := parse.NewParser(rootPackage, rootMessage, pos...)
	if err != nil {
		return nil, err
	}
	tagged := tag.NewTagger(underlyingParser, tag.WithAllocator(p.Alloc), tag.WithIndexes(), tag.WithTags())
	return &parser{parserWithReset: tagged, underlying: underlyingParser, pool: p}, nil
}

func (p *parser) Init(buf []byte) {
	// This Init really inits the underlying parser with the new buffer.
	p.parserWithReset.Reset()
	p.underlying.Init(buf)
	p.pool.FreeAll()
	return
}
