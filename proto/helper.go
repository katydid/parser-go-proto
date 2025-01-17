// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

import (
	"strings"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func GetMapFields(msg *descriptor.DescriptorProto) (*descriptor.FieldDescriptorProto, *descriptor.FieldDescriptorProto) {
	if !msg.GetOptions().GetMapEntry() {
		return nil, nil
	}
	return msg.GetField()[0], msg.GetField()[1]
}

func WireType(field *descriptor.FieldDescriptorProto) (wire int) {
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return 1
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return 5
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return 1
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return 5
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return 2
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return 2
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return 2
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return 2
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return 5
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return 1
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return 0
	}
	panic("unreachable")
}

func GetKeyUint64(field *descriptor.FieldDescriptorProto) (x uint64) {
	packed := IsPacked(field)
	wireType := WireType(field)
	fieldNumber := field.GetNumber()
	if packed {
		wireType = 2
	}
	x = uint64(uint32(fieldNumber)<<3 | uint32(wireType))
	return x
}

func GetKey3Uint64(field *descriptor.FieldDescriptorProto) (x uint64) {
	packed := IsPacked(field)
	wireType := WireType(field)
	fieldNumber := field.GetNumber()
	if packed {
		wireType = 2
	}
	x = uint64(uint32(fieldNumber)<<3 | uint32(wireType))
	return x
}

func GetKey(field *descriptor.FieldDescriptorProto) []byte {
	x := GetKeyUint64(field)
	i := 0
	keybuf := make([]byte, 0)
	for i = 0; x > 127; i++ {
		keybuf = append(keybuf, 0x80|uint8(x&0x7F))
		x >>= 7
	}
	keybuf = append(keybuf, uint8(x))
	return keybuf
}

func GetKey3(field *descriptor.FieldDescriptorProto) []byte {
	x := GetKey3Uint64(field)
	i := 0
	keybuf := make([]byte, 0)
	for i = 0; x > 127; i++ {
		keybuf = append(keybuf, 0x80|uint8(x&0x7F))
		x >>= 7
	}
	keybuf = append(keybuf, uint8(x))
	return keybuf
}

func GetField(desc *descriptor.FileDescriptorSet, packageName, messageName, fieldName string) *descriptor.FieldDescriptorProto {
	msg := GetMessage(desc, packageName, messageName)
	if msg == nil {
		return nil
	}
	for _, field := range msg.GetField() {
		if field.GetName() == fieldName {
			return field
		}
	}
	return nil
}

func GetNestedMessage(file *descriptor.FileDescriptorProto, msg *descriptor.DescriptorProto, typeName string) *descriptor.DescriptorProto {
	for _, nes := range msg.GetNestedType() {
		if nes.GetName() == typeName {
			return nes
		}
		res := GetNestedMessage(file, nes, strings.TrimPrefix(typeName, nes.GetName()+"."))
		if res != nil {
			return res
		}
	}
	return nil
}

func GetMessage(desc *descriptor.FileDescriptorSet, packageName string, typeName string) *descriptor.DescriptorProto {
	for _, file := range desc.GetFile() {
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) {
			continue
		}
		for _, msg := range file.GetMessageType() {
			if msg.GetName() == typeName {
				return msg
			}
		}
		for _, msg := range file.GetMessageType() {
			for _, nes := range msg.GetNestedType() {
				if nes.GetName() == typeName {
					return nes
				}
				if msg.GetName()+"."+nes.GetName() == typeName {
					return nes
				}
			}
		}
	}
	return nil
}

func IsProto3(desc *descriptor.FileDescriptorSet, packageName string, typeName string) bool {
	for _, file := range desc.GetFile() {
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) {
			continue
		}
		for _, msg := range file.GetMessageType() {
			if msg.GetName() == typeName {
				return file.GetSyntax() == "proto3"
			}
		}
		for _, msg := range file.GetMessageType() {
			for _, nes := range msg.GetNestedType() {
				if nes.GetName() == typeName {
					return file.GetSyntax() == "proto3"
				}
				if msg.GetName()+"."+nes.GetName() == typeName {
					return file.GetSyntax() == "proto3"
				}
			}
		}
	}
	return false
}

func IsExtendable(msg *descriptor.DescriptorProto) bool {
	return len(msg.GetExtensionRange()) > 0
}

func FindExtension(desc *descriptor.FileDescriptorSet, packageName string, typeName string, fieldName string) (extPackageName string, field *descriptor.FieldDescriptorProto) {
	parent := GetMessage(desc, packageName, typeName)
	if parent == nil {
		return "", nil
	}
	if !IsExtendable(parent) {
		return "", nil
	}
	extendee := "." + packageName + "." + typeName
	for _, file := range desc.GetFile() {
		for _, ext := range file.GetExtension() {
			if strings.Map(dotToUnderscore, file.GetPackage()) == strings.Map(dotToUnderscore, packageName) {
				if !(ext.GetExtendee() == typeName || ext.GetExtendee() == extendee) {
					continue
				}
			} else {
				if ext.GetExtendee() != extendee {
					continue
				}
			}
			if ext.GetName() == fieldName {
				return file.GetPackage(), ext
			}
		}
	}
	return "", nil
}

func FindExtensionByFieldNumber(desc *descriptor.FileDescriptorSet, packageName string, typeName string, fieldNum int32) (extPackageName string, field *descriptor.FieldDescriptorProto) {
	parent := GetMessage(desc, packageName, typeName)
	if parent == nil {
		return "", nil
	}
	if !IsExtendable(parent) {
		return "", nil
	}
	extendee := "." + packageName + "." + typeName
	for _, file := range desc.GetFile() {
		for _, ext := range file.GetExtension() {
			if strings.Map(dotToUnderscore, file.GetPackage()) == strings.Map(dotToUnderscore, packageName) {
				if !(ext.GetExtendee() == typeName || ext.GetExtendee() == extendee) {
					continue
				}
			} else {
				if ext.GetExtendee() != extendee {
					continue
				}
			}
			if ext.GetNumber() == fieldNum {
				return file.GetPackage(), ext
			}
		}
	}
	return "", nil
}

func FindMessage(desc *descriptor.FileDescriptorSet, packageName string, typeName string, fieldName string) (msgPackageName string, msgName string) {
	parent := GetMessage(desc, packageName, typeName)
	if parent == nil {
		return "", ""
	}
	field := GetFieldDescriptor(parent, fieldName)
	if field == nil {
		var extPackageName string
		extPackageName, field = FindExtension(desc, packageName, typeName, fieldName)
		if field == nil {
			return "", ""
		}
		packageName = extPackageName
	}
	typeNames := strings.Split(field.GetTypeName(), ".")
	if len(typeNames) == 1 {
		msg := GetMessage(desc, packageName, typeName)
		if msg == nil {
			return "", ""
		}
		return packageName, msg.GetName()
	}
	if len(typeNames) > 2 {
		for i := 1; i < len(typeNames)-1; i++ {
			packageName = strings.Join(typeNames[1:len(typeNames)-i], ".")
			typeName = strings.Join(typeNames[len(typeNames)-i:], ".")
			msg := GetMessage(desc, packageName, typeName)
			if msg != nil {
				typeNames := strings.Split(msg.GetName(), ".")
				if len(typeNames) == 1 {
					return packageName, msg.GetName()
				}
				return strings.Join(typeNames[1:len(typeNames)-1], "."), typeNames[len(typeNames)-1]
			}
		}
	}
	return "", ""
}

func GetFieldDescriptor(msg *descriptor.DescriptorProto, fieldName string) *descriptor.FieldDescriptorProto {
	for _, field := range msg.GetField() {
		if field.GetName() == fieldName {
			return field
		}
	}
	return nil
}

func GetEnum(desc *descriptor.FileDescriptorSet, packageName string, typeName string) *descriptor.EnumDescriptorProto {
	for _, file := range desc.GetFile() {
		if strings.Map(dotToUnderscore, file.GetPackage()) != strings.Map(dotToUnderscore, packageName) {
			continue
		}
		for _, enum := range file.GetEnumType() {
			if enum.GetName() == typeName {
				return enum
			}
		}
	}
	return nil
}

func IsEnum(f *descriptor.FieldDescriptorProto) bool {
	return *f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM
}

func IsMessage(f *descriptor.FieldDescriptorProto) bool {
	return *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

func IsBytes(f *descriptor.FieldDescriptorProto) bool {
	return *f.Type == descriptor.FieldDescriptorProto_TYPE_BYTES
}

func IsRepeated(f *descriptor.FieldDescriptorProto) bool {
	return f.Label != nil && *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

func IsString(f *descriptor.FieldDescriptorProto) bool {
	return *f.Type == descriptor.FieldDescriptorProto_TYPE_STRING
}

func IsBool(f *descriptor.FieldDescriptorProto) bool {
	return *f.Type == descriptor.FieldDescriptorProto_TYPE_BOOL
}

func IsRequired(f *descriptor.FieldDescriptorProto) bool {
	return f.Label != nil && *f.Label == descriptor.FieldDescriptorProto_LABEL_REQUIRED
}

func IsPacked(f *descriptor.FieldDescriptorProto) bool {
	return f.Options != nil && f.GetOptions().GetPacked()
}

func IsPacked3(f *descriptor.FieldDescriptorProto) bool {
	if IsRepeated(f) && IsScalar(f) {
		if f.Options == nil || f.GetOptions().Packed == nil {
			return true
		}
		return f.Options != nil && f.GetOptions().GetPacked()
	}
	return false
}

// Is this field a scalar numeric type?
func IsScalar(field *descriptor.FieldDescriptorProto) bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return true
	default:
		return false
	}
}

func HasExtension(m *descriptor.DescriptorProto) bool {
	return len(m.ExtensionRange) > 0
}
