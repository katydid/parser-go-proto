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
    "katydid.org.za/go/parser-go/parse"
    protoparser "katydid.org.za/go/parser-go-proto/proto"
)

func NewMyMessageParser(marshaledMyMessage []byte) (parse.Parser, error) {
	protoParser, err := protoparser.NewParser("prototests", "mymessage")
	if err != nil {
		return nil, err
	}
	protoParser.Init(marshaledMyMessage)
	return protoParser, nil
}
```

We can then use the parser to decode only `myfield` and skip over other fields and return `"myvalue"`:

```go
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
```

For more details on how to use the parser's methods like Next, Skip and Token see the [Parser Interface](https://git.katydid.org.za/parser-go).

## Known Issues

This is a online parser, which does not allocate memory, which means it cannot support certain protobuf features:

  * The parser does not return defaults or proto3 zero values. It will simply skip fields if they are not present in the serialized data.
  * Parsing of merged fields will result in those fields being returned twice, instead of once. You can use `NoLatentAppendingOrMerging` to check that the serialized data does not contain merged fields.

