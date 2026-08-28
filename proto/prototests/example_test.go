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

package prototests

import (
	"testing"

	"google.golang.org/protobuf/proto"
	protoparser "katydid.org.za/go/parser-go-proto/proto"
	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/parse"
)

func NewMarshaledMyMessage() ([]byte, error) {
	msg := &Mymessage{Myfield: "myvalue"}
	return proto.Marshal(msg)
}

func NewMyMessageParser(marshaledMyMessage []byte) (parse.Parser, error) {
	protoParser, err := protoparser.NewParser("prototests", "mymessage")
	if err != nil {
		return nil, err
	}
	protoParser.Init(marshaledMyMessage)
	return protoParser, nil
}

func GetMyField(p parse.Parser) (string, error) {
	if hint, err := p.Next(); err != nil || hint != parse.EnterHint {
		return "", err
	}
	for {
		if hint, err := p.Next(); err != nil || hint != parse.FieldHint {
			return "", err
		}
		kind, fieldNameBytes, err := p.Token()
		if err != nil || kind != parse.StringKind {
			return "", err
		}
		var fieldName string
		cast.ToStringPtr(fieldNameBytes, &fieldName) // cast, do not copy
		if fieldName != "myfield" {
			if err := p.Skip(); err != nil {
				return "", err
			}
			continue
		}
		if hint, err := p.Next(); err != nil || hint != parse.ValueHint {
			return "", err
		}
		kind, val, err := p.Token()
		if err != nil || kind != parse.StringKind {
			return "", err
		}
		return string(val), nil
	}
}

func TestExample(t *testing.T) {
	data, err := NewMarshaledMyMessage()
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewMyMessageParser(data)
	if err != nil {
		t.Fatal(err)
	}
	myvalue, err := GetMyField(parser)
	if err != nil {
		t.Fatal(err)
	}
	if myvalue != "myvalue" {
		t.Fatalf("want %v got %v", "myvalue", myvalue)
	}
}
