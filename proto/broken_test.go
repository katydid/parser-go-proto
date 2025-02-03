//  Copyright 2015 Walter Schulze
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
	"testing"

	"github.com/katydid/parser-go-proto/proto/prototests"
	"github.com/katydid/parser-go/parser/debug"
	"google.golang.org/protobuf/proto"
)

func TestBrokenLengthValue(t *testing.T) {
	msg := &prototests.Mymessage{Myfield: "myvalue"}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mymessage: %v\n", data)
	// change string length of first field to be incorrect.
	data[1] = 12
	t.Logf("mymessage: %v\n", data)
	protoParser, err := NewParser("prototests", "mymessage")
	if err != nil {
		t.Fatal(err)
	}
	if err := protoParser.Init(data); err != nil {
		t.Fatal(err)
	}
	// make sure the parser doesn't panic and only returns an error.
	if _, err := debug.Parse(protoParser); err == nil {
		t.Fatal("expected error, because of wrong length")
	}
}

func TestBrokenLengthMessage(t *testing.T) {
	msg := &prototests.BigMsg{
		Msg: &prototests.SmallMsg{
			ScarBusStop: proto.String("abc"),
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("bigmsg: %v\n", data)
	// change string length of first field to be incorrect.
	data[1] = 12
	t.Logf("bigmsg: %v\n", data)
	protoParser, err := NewParser("prototests", "BigMsg")
	if err != nil {
		t.Fatal(err)
	}
	if err := protoParser.Init(data); err != nil {
		t.Fatal(err)
	}
	// make sure the parser doesn't panic and only returns an error.
	if _, err := debug.Parse(protoParser); err == nil {
		t.Fatal("expected error, because of wrong length")
	}
}
