## parser-go-proto

Parser for Protocol Buffers in Go.

The parser dynamically parses the serialized protocol buffer bytes and inspect fields without unmarshaling into a Go `struct` first.
This parser is fast, since it does not allocate any memory.

This can be used to dynamically inspect serialized protobufs or create a filter for protobufs stored on disk.

# Usage

Given a protocol buffer:

```proto
syntax = "proto3";
package mypackage;

...

message mymessage {
  string myfield = 1;
  int64 otherfield = 2;
}
```

We can marshal it and store it on disk or pass it over a TCP connection to another process:

```go
import "google.golang.org/protobuf/proto"

func NewMarshaledMyMessage() ([]byte, error) {
	msg := &Mymessage{Myfield: "myvalue"}
	return proto.Marshal(msg)
}
```

The new process can construct the parser for the marshaled bytes:

```go
import (
    "github.com/katydid/parser-go/parser"
    "github.com/katydid/parser-go-proto/proto"
)

func NewMyMessageParser(marshaledMyMessage []byte) (parser.Interface, error) {
	mymessageParser, err := proto.NewParser("mypackage", "mymessage")
	if err != nil {
		return nil, err
	}
	if err := mymessageParser.Init(marshaledMyMessage); err != nil {
		return nil, err
	}
	return mymessageParser, nil
}
```

We can then use the parser to decode only `myfield` and skip over other fields and return `"myvalue"`:

```go
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
```

## Known Issues

This is a online parser, which does not allocate memory, which means it cannot support certain protobuf features:

  * The parser does not return defaults or proto3 zero values. It will simply skip fields if they are not present in the serialized data.
  * Parsing of merged fields will result in those fields being returned twice, instead of once. You can use `NoLatentAppendingOrMerging` to check that the serialized data does not contain merged fields.

