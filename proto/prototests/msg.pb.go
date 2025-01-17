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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: msg.proto

package prototests

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// BigMsg contains a field and a message field.
type BigMsg struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Field         *int64                 `protobuf:"varint,1,opt,name=Field" json:"Field,omitempty"`
	Msg           *SmallMsg              `protobuf:"bytes,3,opt,name=Msg" json:"Msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BigMsg) Reset() {
	*x = BigMsg{}
	mi := &file_msg_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BigMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BigMsg) ProtoMessage() {}

func (x *BigMsg) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BigMsg.ProtoReflect.Descriptor instead.
func (*BigMsg) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *BigMsg) GetField() int64 {
	if x != nil && x.Field != nil {
		return *x.Field
	}
	return 0
}

func (x *BigMsg) GetMsg() *SmallMsg {
	if x != nil {
		return x.Msg
	}
	return nil
}

// SmallMsg only contains some native fields.
type SmallMsg struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	ScarBusStop     *string                `protobuf:"bytes,1,opt,name=ScarBusStop" json:"ScarBusStop,omitempty"`
	FlightParachute []uint32               `protobuf:"fixed32,12,rep,name=FlightParachute" json:"FlightParachute,omitempty"`
	MapShark        *string                `protobuf:"bytes,18,opt,name=MapShark" json:"MapShark,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *SmallMsg) Reset() {
	*x = SmallMsg{}
	mi := &file_msg_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SmallMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SmallMsg) ProtoMessage() {}

func (x *SmallMsg) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SmallMsg.ProtoReflect.Descriptor instead.
func (*SmallMsg) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *SmallMsg) GetScarBusStop() string {
	if x != nil && x.ScarBusStop != nil {
		return *x.ScarBusStop
	}
	return ""
}

func (x *SmallMsg) GetFlightParachute() []uint32 {
	if x != nil {
		return x.FlightParachute
	}
	return nil
}

func (x *SmallMsg) GetMapShark() string {
	if x != nil && x.MapShark != nil {
		return *x.MapShark
	}
	return ""
}

// Packed contains some repeated packed fields.
type Packed struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ints          []int64                `protobuf:"varint,4,rep,packed,name=Ints" json:"Ints,omitempty"`
	Floats        []float64              `protobuf:"fixed64,5,rep,packed,name=Floats" json:"Floats,omitempty"`
	Uints         []uint32               `protobuf:"varint,6,rep,packed,name=Uints" json:"Uints,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Packed) Reset() {
	*x = Packed{}
	mi := &file_msg_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Packed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packed) ProtoMessage() {}

func (x *Packed) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packed.ProtoReflect.Descriptor instead.
func (*Packed) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

func (x *Packed) GetInts() []int64 {
	if x != nil {
		return x.Ints
	}
	return nil
}

func (x *Packed) GetFloats() []float64 {
	if x != nil {
		return x.Floats
	}
	return nil
}

func (x *Packed) GetUints() []uint32 {
	if x != nil {
		return x.Uints
	}
	return nil
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x22, 0x46, 0x0a, 0x06, 0x42, 0x69, 0x67, 0x4d, 0x73,
	0x67, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x26, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x65, 0x73, 0x74,
	0x73, 0x2e, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x4d, 0x73, 0x67, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x22,
	0x72, 0x0a, 0x08, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x4d, 0x73, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x53,
	0x63, 0x61, 0x72, 0x42, 0x75, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x53, 0x63, 0x61, 0x72, 0x42, 0x75, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x28, 0x0a,
	0x0f, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x50, 0x61, 0x72, 0x61, 0x63, 0x68, 0x75, 0x74, 0x65,
	0x18, 0x0c, 0x20, 0x03, 0x28, 0x07, 0x52, 0x0f, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x50, 0x61,
	0x72, 0x61, 0x63, 0x68, 0x75, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x61, 0x70, 0x53, 0x68,
	0x61, 0x72, 0x6b, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x61, 0x70, 0x53, 0x68,
	0x61, 0x72, 0x6b, 0x22, 0x56, 0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x16, 0x0a,
	0x04, 0x49, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x03, 0x42, 0x02, 0x10, 0x01, 0x52,
	0x04, 0x49, 0x6e, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x06, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x01, 0x42, 0x02, 0x10, 0x01, 0x52, 0x06, 0x46, 0x6c, 0x6f, 0x61, 0x74,
	0x73, 0x12, 0x18, 0x0a, 0x05, 0x55, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0d,
	0x42, 0x02, 0x10, 0x01, 0x52, 0x05, 0x55, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x74, 0x79, 0x64, 0x69,
	0x64, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2d, 0x67, 0x6f, 0x2d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x65, 0x73,
	0x74, 0x73,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_msg_proto_goTypes = []any{
	(*BigMsg)(nil),   // 0: prototests.BigMsg
	(*SmallMsg)(nil), // 1: prototests.SmallMsg
	(*Packed)(nil),   // 2: prototests.Packed
}
var file_msg_proto_depIdxs = []int32{
	1, // 0: prototests.BigMsg.Msg:type_name -> prototests.SmallMsg
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}
