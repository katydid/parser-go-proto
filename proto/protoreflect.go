//  Copyright 2025 Walter Schulze
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
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// NewFileDescriptorSet is a helper function that retrieves the FileDescriptorSet from the protoregistry.
// If protoreflect.FileDescriptor's are manually provided, then the protoregistry is not queried and
// the provided protoreflect.FileDescriptor's are simply converted to the FileDescriptorSet.
func NewFileDescriptorSet(reflectFileDescriptors ...protoreflect.FileDescriptor) *descriptor.FileDescriptorSet {
	if len(reflectFileDescriptors) == 0 {
		reflectFileDescriptors = make([]protoreflect.FileDescriptor, 0, protoregistry.GlobalFiles.NumFiles())
		protoregistry.GlobalFiles.RangeFiles(func(f protoreflect.FileDescriptor) bool {
			reflectFileDescriptors = append(reflectFileDescriptors, f)
			return true
		})
	}
	fileDescriptors := make([]*descriptor.FileDescriptorProto, len(reflectFileDescriptors))
	for i, rfd := range reflectFileDescriptors {
		fileDescriptors[i] = protodesc.ToFileDescriptorProto(rfd)
	}
	return &descriptor.FileDescriptorSet{File: fileDescriptors}
}
