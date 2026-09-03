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

package proto

import (
	"fmt"

	"katydid.org.za/go/parser-go/cp"
	"katydid.org.za/go/parser-go/parse"
)

// NoLatentAppendingOrMerging returns whether the byte buffer has some latent fields.
// Latent fields are those fields you have already seen on your walk, but then after seeing a different field you see this field again.
// This typically happens when the protocol buffer user created an object, marshaled it and then merged it with another value.
func NoLatentAppendingOrMerging(pkgName, msgName string, data []byte) error {
	parser, err := NewParser(pkgName, msgName)
	if err != nil {
		return err
	}
	parser.Init(data)
	return noLatentAppendingOrMerging(parser)
}

func noLatentAppendingOrMerging(parser parse.Parser) error {
	hint, err := parser.Next()
	if err != nil {
		return err
	}
	if hint == parse.ValueHint {
		return nil
	}
	if hint != parse.EnterHint {
		return nil
	}
	seen := make(map[string]bool)
	seeni := make(map[int64]bool)
	for {
		if hint, err := parser.Next(); err != nil || hint != parse.FieldHint {
			return err
		}
		kind, val, err := parser.Token()
		if err != nil {
			return err
		}
		switch kind {
		case parse.StringKind:
			fieldName := cp.ToString(val)
			if _, ok := seen[fieldName]; ok {
				return fmt.Errorf("%s requires merging", fieldName)
			}
			seen[fieldName] = true
		case parse.Int64Kind:
			index := cp.ToInt64(val)
			if _, ok := seeni[index]; ok {
				return fmt.Errorf("%d requires merging", index)
			}
			seeni[index] = true
		}
		if err := noLatentAppendingOrMerging(parser); err != nil {
			return err
		}
	}
}
