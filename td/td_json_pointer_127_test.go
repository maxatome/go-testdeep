// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td_test

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestJSONPointerv2(t *testing.T) {
	checkOK(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": 42}}},
		td.JSONPointer("/a/b/c", 42, json.DefaultOptionsV2()))

	checkOK(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": "42"}}},
		td.JSONPointer("/a/b/c", int8(42),
			json.DefaultOptionsV2(),
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *int8) error {
				_, err := fmt.Sscanf(string(b), `"%d"`, v)
				return err
			})),
		))
}

func ExampleJSONPointer_jsonv2() {
	t := &testing.T{}

	got := jsonv1.RawMessage(`{"a": {"b": {"c": "42"}}}`)

	ok := td.Cmp(t, got, td.JSONPointer("/a/b/c", "42", jsonv1.DefaultOptionsV1()))
	fmt.Println("using json/v1:", ok, "(same as default)")

	ok = td.Cmp(t, got, td.JSONPointer("/a/b/c", "42", json.DefaultOptionsV2()))
	fmt.Println("using json/v2:", ok)

	ok = td.Cmp(t, got,
		td.JSONPointer("/a/b/c", int8(42),
			json.DefaultOptionsV2(),
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *int8) error {
				_, err := fmt.Sscanf(string(b), `"%d"`, v)
				return err
			})),
		))
	fmt.Println("multiple json/v2.Options + custom unmarshaler:", ok)

	// Output:
	// using json/v1: true (same as default)
	// using json/v2: true
	// multiple json/v2.Options + custom unmarshaler: true
}
