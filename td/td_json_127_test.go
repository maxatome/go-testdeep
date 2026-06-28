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
	"strconv"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestJSONv2(t *testing.T) {
	type MyStruct struct {
		Name string `json:"name"`
		Age  uint   `json:"age"`
	}

	got := MyStruct{Name: "Bob", Age: 42}

	checkOK(t, got,
		td.JSON(`{"name":$2,"age":"🕯42🕯"}`,
			json.WithMarshalers(json.MarshalFunc(func(n uint) ([]byte, error) {
				return fmt.Appendf(nil, `"🕯%d🕯"`, n), nil
			})),
			"Bob"))

	checkOK(t, got,
		td.JSON(`{"name":$1,"age":"🕯42🕯"}`,
			"Bob",
			json.WithMarshalers(json.MarshalFunc(func(n uint) ([]byte, error) {
				return fmt.Appendf(nil, `"🕯%d🕯"`, n), nil
			}))))

	checkOK(t, got,
		td.JSON(`{"name":$1,"age":$4}`,
			"Bob",
			json.WithMarshalers(json.MarshalFunc(func(n uint) ([]byte, error) {
				return fmt.Appendf(nil, `"%d"`, n), nil
			})),
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *int) error {
				_, err := fmt.Sscanf(string(b), `"%d"`, v)
				return err
			})),
			42))

	// Pass jsonv2.Options to JSONPointer operator
	checkOK(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": 42}}},
		td.JSON(`{"a": JSONPointer("/b/c", 42)}`,
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *float64) error {
				var err error
				if b[0] != '"' {
					// Used by JSON to convert original `42` to a float64
					*v, err = strconv.ParseFloat(string(b), 64)
				} else {
					// Used by JSONPointer to convert `"🕯42🕯"` to the expected float64
					_, err = fmt.Sscanf(string(b), `"🕯%f🕯"`, v)
				}
				return err
			})),
			json.WithMarshalers(json.MarshalFunc(func(n float64) ([]byte, error) {
				// Used by JSON to convert float64(42) to JSON string `"🕯42🕯"`
				return fmt.Appendf(nil, `"🕯%.0f🕯"`, n), nil
			})),
		))

	// With json/v2 unmarshaling can fail
	checkError(t, got,
		td.JSON(`{"name":"Bob","age":42}`,
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *float64) error {
				return fmt.Errorf("an error occurred")
			}))),
		expectedError{
			Message: mustBe("json.Unmarshal failed"),
			Summary: mustContain(": an error occurred"),
			Under:   mustContain("under operator JSON at td_json_127_test.go:"),
		})

	test.EqualStr(t,
		td.JSON(`{"x":123}`, json.DefaultOptionsV2()).String(),
		`JSON({
       "x": 123
     })`)
}

func ExampleJSON_jsonv2() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	// Same as default
	ok := td.Cmp(t, got, td.JSON(`{"age":42,"fullname":"Bob"}`, jsonv1.DefaultOptionsV1()))
	fmt.Println("using json/v1:", ok, "(same as default)")

	ok = td.Cmp(t, got, td.JSON(`{"age":42,"fullname":"Bob"}`, json.DefaultOptionsV2()))
	fmt.Println("using json/v2:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"fullname":$1,"age":"🕯42🕯"}`,
			"Bob",
			json.WithMarshalers(json.MarshalFunc(func(n int) ([]byte, error) {
				return fmt.Appendf(nil, `"🕯%d🕯"`, n), nil
			}))),
	)
	fmt.Println("with custom marshaler:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"fullname":$1,"age":$4}`,
			"Bob",
			json.WithMarshalers(json.MarshalFunc(func(n int) ([]byte, error) {
				return fmt.Appendf(nil, `"%d"`, n), nil
			})),
			json.WithUnmarshalers(json.UnmarshalFunc(func(b []byte, v *int) error {
				_, err := fmt.Sscanf(string(b), `"%d"`, v)
				return err
			})),
			42),
	)
	fmt.Println("with custom marshaler & unmarshaler:", ok)

	// Output:
	// using json/v1: true (same as default)
	// using json/v2: true
	// with custom marshaler: true
	// with custom marshaler & unmarshaler: true
}
