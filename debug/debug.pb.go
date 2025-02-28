// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: debug.proto

package debug

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

type Debug struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	A             *int64                 `protobuf:"varint,1,opt,name=A" json:"A,omitempty"`
	B             []string               `protobuf:"bytes,2,rep,name=B" json:"B,omitempty"`
	C             *Debug                 `protobuf:"bytes,3,opt,name=C" json:"C,omitempty"`
	D             *int32                 `protobuf:"varint,4,opt,name=D" json:"D,omitempty"`
	E             []*Debug               `protobuf:"bytes,5,rep,name=E" json:"E,omitempty"`
	F             []uint32               `protobuf:"fixed32,6,rep,name=F" json:"F,omitempty"`
	G             *float64               `protobuf:"fixed64,7,opt,name=G" json:"G,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Debug) Reset() {
	*x = Debug{}
	mi := &file_debug_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Debug) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Debug) ProtoMessage() {}

func (x *Debug) ProtoReflect() protoreflect.Message {
	mi := &file_debug_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Debug.ProtoReflect.Descriptor instead.
func (*Debug) Descriptor() ([]byte, []int) {
	return file_debug_proto_rawDescGZIP(), []int{0}
}

func (x *Debug) GetA() int64 {
	if x != nil && x.A != nil {
		return *x.A
	}
	return 0
}

func (x *Debug) GetB() []string {
	if x != nil {
		return x.B
	}
	return nil
}

func (x *Debug) GetC() *Debug {
	if x != nil {
		return x.C
	}
	return nil
}

func (x *Debug) GetD() int32 {
	if x != nil && x.D != nil {
		return *x.D
	}
	return 0
}

func (x *Debug) GetE() []*Debug {
	if x != nil {
		return x.E
	}
	return nil
}

func (x *Debug) GetF() []uint32 {
	if x != nil {
		return x.F
	}
	return nil
}

func (x *Debug) GetG() float64 {
	if x != nil && x.G != nil {
		return *x.G
	}
	return 0
}

var File_debug_proto protoreflect.FileDescriptor

var file_debug_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x64, 0x65, 0x62, 0x75, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x64,
	0x65, 0x62, 0x75, 0x67, 0x22, 0x85, 0x01, 0x0a, 0x05, 0x44, 0x65, 0x62, 0x75, 0x67, 0x12, 0x0c,
	0x0a, 0x01, 0x41, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x41, 0x12, 0x0c, 0x0a, 0x01,
	0x42, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x01, 0x42, 0x12, 0x1a, 0x0a, 0x01, 0x43, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x62, 0x75, 0x67, 0x2e, 0x44, 0x65,
	0x62, 0x75, 0x67, 0x52, 0x01, 0x43, 0x12, 0x0c, 0x0a, 0x01, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x01, 0x44, 0x12, 0x1a, 0x0a, 0x01, 0x45, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x64, 0x65, 0x62, 0x75, 0x67, 0x2e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x52, 0x01, 0x45,
	0x12, 0x0c, 0x0a, 0x01, 0x46, 0x18, 0x06, 0x20, 0x03, 0x28, 0x07, 0x52, 0x01, 0x46, 0x12, 0x0c,
	0x0a, 0x01, 0x47, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x47, 0x42, 0x2a, 0x5a, 0x28,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x74, 0x79, 0x64,
	0x69, 0x64, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2d, 0x67, 0x6f, 0x2d, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x64, 0x65, 0x62, 0x75, 0x67,
}

var (
	file_debug_proto_rawDescOnce sync.Once
	file_debug_proto_rawDescData = file_debug_proto_rawDesc
)

func file_debug_proto_rawDescGZIP() []byte {
	file_debug_proto_rawDescOnce.Do(func() {
		file_debug_proto_rawDescData = protoimpl.X.CompressGZIP(file_debug_proto_rawDescData)
	})
	return file_debug_proto_rawDescData
}

var file_debug_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_debug_proto_goTypes = []any{
	(*Debug)(nil), // 0: debug.Debug
}
var file_debug_proto_depIdxs = []int32{
	0, // 0: debug.Debug.C:type_name -> debug.Debug
	0, // 1: debug.Debug.E:type_name -> debug.Debug
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_debug_proto_init() }
func file_debug_proto_init() {
	if File_debug_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_debug_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_debug_proto_goTypes,
		DependencyIndexes: file_debug_proto_depIdxs,
		MessageInfos:      file_debug_proto_msgTypes,
	}.Build()
	File_debug_proto = out.File
	file_debug_proto_rawDesc = nil
	file_debug_proto_goTypes = nil
	file_debug_proto_depIdxs = nil
}
