//  Copyright 2013 Walter Schulze
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

// Package proto contains an implementation of a protocol buffer parser.
//
// Defaults and proto3 zero values will not be returned. Fields that are not present in the serialized data, will not be returned.
//
// Merging of fields and splitting of arrays are not supported by this parser for optimization reasons.
// Use the NoLatentAppendingOrMerging function to check whether the marshaled buffer conforms to the limitations.
package proto
