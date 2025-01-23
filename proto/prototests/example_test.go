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
	"io"
	"testing"

	protoparser "github.com/katydid/parser-go-proto/proto"
	"github.com/katydid/parser-go/parser"
	"google.golang.org/protobuf/proto"
)

func NewMarshaledMyMessage() ([]byte, error) {
	msg := &Mymessage{Myfield: "myvalue"}
	return proto.Marshal(msg)
}

func NewMyMessageParser(marshaledMyMessage []byte) (parser.Interface, error) {
	protoParser, err := protoparser.NewParser("prototests", "mymessage")
	if err != nil {
		return nil, err
	}
	if err := protoParser.Init(marshaledMyMessage); err != nil {
		return nil, err
	}
	return protoParser, nil
}

func GetMyField(p parser.Interface) (string, error) {
	for {
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		fieldName, err := p.String()
		if err != nil {
			return "", err
		}
		if fieldName != "myfield" {
			continue
		}
		p.Down()
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		return p.String()
	}
	return "", nil
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
