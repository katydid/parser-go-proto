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

package proto_test

import (
	"encoding/binary"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/katydid/parser-go-proto/debug"
	protoparser "github.com/katydid/parser-go-proto/proto"
	"github.com/katydid/parser-go-proto/proto/prototests"
	"google.golang.org/protobuf/proto"
)

func noMerge(data []byte, pkgName, msgName string) error {
	parser, err := protoparser.NewParser(pkgName, msgName)
	if err != nil {
		return err
	}
	if err := parser.Init(data); err != nil {
		return err
	}
	return protoparser.NoLatentAppendingOrMerging(parser)
}

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestNoMergeNoMerge(t *testing.T) {
	m := debug.Input
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	err = noMerge(data, "debug", "Debug")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoMergeMerge(t *testing.T) {
	m := debug.Input
	m.G = proto.Float64(1.1)
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	key := byte(uint32(7)<<3 | uint32(1))
	data = append(data, key, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	err = noMerge(data, "debug", "Debug")
	if err == nil || !strings.Contains(err.Error(), "G requires merging") {
		t.Fatalf("G should require merging")
	}
}

func TestNoMergeLatent(t *testing.T) {
	m := debug.Input
	m.F = []uint32{5, 6, 7}
	m.G = proto.Float64(1.1)
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	key := byte(uint32(6)<<3 | uint32(5))
	data = append(data, key, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	err = noMerge(data, "debug", "Debug")
	if err == nil || !strings.Contains(err.Error(), "F") {
		t.Fatalf("F should have latent appending")
	}
}

func TestNoMergeNestedNoMerge(t *testing.T) {
	bigm := prototests.NewPopulatedBigMsg(r)
	data, err := proto.Marshal(bigm)
	if err != nil {
		t.Fatal(err)
	}
	err = noMerge(data, "prototests", "BigMsg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoMergeMessageMerge(t *testing.T) {
	bigm := prototests.NewPopulatedBigMsg(r)
	bigm.Msg = prototests.NewPopulatedSmallMsg(r)
	data, err := proto.Marshal(bigm)
	if err != nil {
		t.Fatal(err)
	}
	smallMsgfieldKey := byte(uint32(3)<<3 | uint32(2))         // 3 field number, 2 wire type
	flightParachuteFieldKey := byte(uint32(12)<<3 | uint32(5)) // 12 field number, 5 wire type
	data = append(data, smallMsgfieldKey, 5, flightParachuteFieldKey, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	err = noMerge(data, "prototests", "BigMsg")
	if err == nil || !strings.Contains(err.Error(), "Msg requires merging") {
		t.Fatalf("Msg should require merging, but got Error: <%v>", err)
	}
}

func TestNoMergeNestedMerge(t *testing.T) {

	m := prototests.NewPopulatedSmallMsg(r)
	if len(m.FlightParachute) == 0 {
		m.FlightParachute = []uint32{1}
	}
	m.MapShark = proto.String("a")
	mdata, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	flightParachuteFieldKey := byte(uint32(12)<<3 | uint32(5)) // 12 field number, 5 wire type
	mdata = append(mdata, flightParachuteFieldKey, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))

	bigm := &prototests.BigMsg{
		Field: proto.Int64(int64(r.Intn(256))),
	}
	bigdata, err := proto.Marshal(bigm)
	if err != nil {
		t.Fatal(err)
	}
	smallMsgfieldKey := byte(uint32(3)<<3 | uint32(2)) // 3 field number, 2 wire type
	bigdata = append(bigdata, smallMsgfieldKey, byte(len(mdata)))
	bigdata = append(bigdata, mdata...)
	err = noMerge(bigdata, "prototests", "BigMsg")
	if err == nil || !strings.Contains(err.Error(), "FlightParachute requires merging") {
		t.Fatalf("FlightParachute should require merging, but got Error: <%v>", err)
	}
}

func TestNoMergeExtensionNoMerge(t *testing.T) {
	bigm := prototests.AContainer
	data, err := proto.Marshal(bigm)
	if err != nil {
		t.Fatal(err)
	}
	err = noMerge(data, "prototests", "Container")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoMergeExtensionMerge(t *testing.T) {
	bigm := prototests.AContainer
	m := &prototests.Small{SmallField: proto.Int64(1)}
	data, err := proto.Marshal(bigm)
	if err != nil {
		t.Fatal(err)
	}
	mdata, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	key := uint32(101)<<3 | uint32(2)
	datakey := make([]byte, 10)
	n := binary.PutUvarint(datakey, uint64(key))
	datakey = datakey[:n]
	datalen := make([]byte, 10)
	n = binary.PutUvarint(datalen, uint64(len(mdata)))
	datalen = datalen[:n]
	data = append(data, append(datakey, append(datalen, mdata...)...)...)
	err = noMerge(data, "prototests", "Container")
	if err == nil || !strings.Contains(err.Error(), "FieldB requires merging") {
		t.Fatalf("FieldB should require merging, but error is %v", err)
	}
	t.Log(err)
}
