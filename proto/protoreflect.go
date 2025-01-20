package proto

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// NewFileDescriptorSet is a helper function that converts multiple protoreflect.FileDescriptor into a FileDescriptorSet.
func NewFileDescriptorSet(reflectFileDescriptors ...protoreflect.FileDescriptor) *descriptor.FileDescriptorSet {
	fileDescriptors := make([]*descriptor.FileDescriptorProto, len(reflectFileDescriptors))
	for i, rfd := range reflectFileDescriptors {
		fileDescriptors[i] = protodesc.ToFileDescriptorProto(rfd)
	}
	return &descriptor.FileDescriptorSet{File: fileDescriptors}
}
