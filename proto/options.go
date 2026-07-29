//  Copyright 2026 Walter Schulze
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

import descriptor "google.golang.org/protobuf/types/descriptorpb"

type Option func(*options)

// WithFileDescriptorSet allows the user to specifiy the FileDescriptorSet manually instead of relying on the protoregistry.
func WithFileDescriptorSet(desc *descriptor.FileDescriptorSet) Option {
	return func(o *options) {
		o.desc = desc
	}
}

// WithTags tags
// 1. each object with an object key, for example `{"a": null}` is parsed as `{"object": {"a": null}}`.
// 2. each array with an array key, for example `{"a": []}` is parsed as `{"a": {"array": []}}`.
func WithTags() func(*options) {
	return func(t *options) {
		t.tag = true
	}
}

// WithIndexes tags each array item with an index:
// for example `["a", "b"]` is parsed as `[0: "a", 1: "b"]`.
// Requires WithTags to also be passed as an option.
func WithIndexes() func(*options) {
	return func(t *options) {
		t.index = true
	}
}

// WithAllocator replaces the default `func(size int) []byte { return make([]byte, size) }` allocator
// with a different allocator function.
// Usually an allocator that uses a pool.
func WithAllocator(alloc func(int) []byte) func(*options) {
	return func(o *options) {
		o.alloc = alloc
	}
}

type options struct {
	alloc func(size int) []byte
	desc  *descriptor.FileDescriptorSet
	index bool
	tag   bool
}

func newOptions(opts ...Option) *options {
	o := &options{
		alloc: func(size int) []byte { return make([]byte, size) },
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.desc == nil {
		o.desc = NewFileDescriptorSet()
	}
	return o
}
