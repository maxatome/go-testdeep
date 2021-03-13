// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type jsonPtrTest int

func (j jsonPtrTest) UnmarshalJSON(b []byte) error {
	return errors.New("jsonPtrTest unmarshal custom error")
}

type jsonPtrMap map[string]interface{}

func (j *jsonPtrMap) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*map[string]interface{})(j))
}

var _ = []json.Unmarshaler{jsonPtrTest(0), &jsonPtrMap{}}

func TestJSONPointer(t *testing.T) {
	//
	// nil
	t.Run("nil", func(t *testing.T) {
		checkOK(t, nil, td.JSONPointer("", nil))
		checkOK(t, (*int)(nil), td.JSONPointer("", nil))
	})

	//
	// Basic types
	t.Run("basic types", func(t *testing.T) {
		checkOK(t, 123, td.JSONPointer("", 123))
		checkOK(t, 123, td.JSONPointer("", td.Between(120, 130)))
		checkOK(t, true, td.JSONPointer("", true))
	})

	//
	// More complex type with encoding/json tags
	t.Run("complex type with json tags", func(t *testing.T) {
		type jpStruct struct {
			Slice []string             `json:"slice,omitempty"`
			Map   map[string]*jpStruct `json:"map,omitempty"`
			Num   int                  `json:"num"`
			Bool  bool                 `json:"bool"`
			Str   string               `json:"str,omitempty"`
		}

		got := jpStruct{
			Slice: []string{"bar", "baz"},
			Map: map[string]*jpStruct{
				"test": {
					Num: 2,
					Str: "level2",
				},
			},
			Num:  1,
			Bool: true,
			Str:  "level1",
		}

		// No filter, should match got or its map representation
		checkOK(t, got, td.JSONPointer("",
			map[string]interface{}{
				"slice": []interface{}{"bar", "baz"},
				"map": map[string]interface{}{
					"test": map[string]interface{}{
						"num":  2,
						"str":  "level2",
						"bool": false,
					},
				},
				"num":  int64(1), // should be OK as Lax is enabled
				"bool": true,
				"str":  "level1",
			}))
		checkOK(t, got, td.JSONPointer("", got))
		checkOK(t, got, td.JSONPointer("", &got))

		// A specific field
		checkOK(t, got, td.JSONPointer("/num", int64(1))) // Lax enabled
		checkOK(t, got, td.JSONPointer("/slice/1", "baz"))
		checkOK(t, got, td.JSONPointer("/map/test/num", 2))
		checkOK(t, got, td.JSONPointer("/map/test/str", td.Contains("vel2")))

		checkOK(t, got,
			td.JSONPointer("/map", td.JSONPointer("/test", td.JSONPointer("/num", 2))))

		checkError(t, got, td.JSONPointer("/zzz/pipo", 666),
			expectedError{
				Message: mustBe("cannot retrieve value via JSON pointer"),
				Path:    mustBe("DATA.JSONPointer</zzz>"),
				Summary: mustBe("key not found"),
			})

		checkError(t, got, td.JSONPointer("/num/pipo", 666),
			expectedError{
				Message: mustBe("cannot retrieve value via JSON pointer"),
				Path:    mustBe("DATA.JSONPointer</num/pipo>"),
				Summary: mustBe("not a map nor an array"),
			})

		checkError(t, got, td.JSONPointer("/slice/2", "zip"),
			expectedError{
				Message: mustBe("cannot retrieve value via JSON pointer"),
				Path:    mustBe("DATA.JSONPointer</slice/2>"),
				Summary: mustBe("out of array range"),
			})

		checkError(t, got, td.JSONPointer("/slice/xxx", "zip"),
			expectedError{
				Message: mustBe("cannot retrieve value via JSON pointer"),
				Path:    mustBe("DATA.JSONPointer</slice/xxx>"),
				Summary: mustBe("array but not an index in JSON pointer"),
			})

		checkError(t, got, td.JSONPointer("/slice/1", "zip"),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA.JSONPointer</slice/1>"),
				Got:      mustBe(`"baz"`),
				Expected: mustBe(`"zip"`),
			})

		// A struct behind a specific field
		checkOK(t, got, td.JSONPointer("/map/test", map[string]interface{}{
			"num":  2,
			"str":  "level2",
			"bool": false,
		}))
		checkOK(t, got, td.JSONPointer("/map/test", jpStruct{
			Num: 2,
			Str: "level2",
		}))
		checkOK(t, got, td.JSONPointer("/map/test", &jpStruct{
			Num: 2,
			Str: "level2",
		}))
		checkOK(t, got, td.JSONPointer("/map/test", td.Struct(&jpStruct{
			Num: 2,
			Str: "level2",
		}, nil)))
	})

	//
	// Complex type without encoding/json tags
	t.Run("complex type without json tags", func(t *testing.T) {
		type jpStruct struct {
			Slice []string
			Map   map[string]*jpStruct
			Num   int
			Bool  bool
			Str   string
		}

		got := jpStruct{
			Slice: []string{"bar", "baz"},
			Map: map[string]*jpStruct{
				"test": {
					Num: 2,
					Str: "level2",
				},
			},
			Num:  1,
			Bool: true,
			Str:  "level1",
		}

		checkOK(t, got, td.JSONPointer("/Num", 1))
		checkOK(t, got, td.JSONPointer("/Slice/1", "baz"))
		checkOK(t, got, td.JSONPointer("/Map/test/Num", 2))
		checkOK(t, got, td.JSONPointer("/Map/test/Str", td.Contains("vel2")))
	})

	//
	// Chained list
	t.Run("Chained list", func(t *testing.T) {
		type Item struct {
			Val  int   `json:"val"`
			Next *Item `json:"next"`
		}
		got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}
		checkOK(t, got, td.JSONPointer("/next/next", Item{Val: 3}))
		checkOK(t, got, td.JSONPointer("/next/next", &Item{Val: 3}))
		checkOK(t, got,
			td.JSONPointer("/next/next",
				td.Struct(Item{}, td.StructFields{"Val": td.Gte(3)})))

		checkOK(t, json.RawMessage(`{"foo":{"bar": {"zip": true}}}`),
			td.JSONPointer("/foo/bar", td.JSON(`{"zip": true}`)))
	})

	//
	// Lax cases
	t.Run("Lax", func(t *testing.T) {
		t.Run("json.Unmarshaler", func(t *testing.T) {
			// Lax not enabled when a json.Unmarshaler type is expected
			got := jsonPtrMap{"x": 123}
			checkOK(t, got, td.JSONPointer("", jsonPtrMap{"x": float64(123)}))
			checkOK(t, got, td.JSONPointer("", &jsonPtrMap{"x": float64(123)}))

			checkError(t, got, td.JSONPointer("", got),
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe(`DATA.JSONPointer<>["x"]`),
					Got:      mustBe(`float64`),
					Expected: mustBe(`int`),
				})
			checkError(t, got, td.JSONPointer("", &got),
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe(`(*DATA.JSONPointer<>)["x"]`),
					Got:      mustBe(`float64`),
					Expected: mustBe(`int`),
				})
		})

		t.Run("struct", func(t *testing.T) {
			// Lax not enabled when a struct or struct pointer is expected
			type jpStruct struct {
				Num interface{}
			}
			got := jpStruct{Num: 123}
			checkOK(t, got, td.JSONPointer("", jpStruct{Num: float64(123)}))
			checkOK(t, jpStruct{Num: got}, td.JSONPointer("/Num", jpStruct{Num: float64(123)}))

			checkError(t, got, td.JSONPointer("", got),
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe(`DATA.JSONPointer<>.Num`),
					Got:      mustBe(`float64`),
					Expected: mustBe(`int`),
				})
			checkError(t, got, td.JSONPointer("", &got),
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe(`DATA.JSONPointer<>.Num`),
					Got:      mustBe(`float64`),
					Expected: mustBe(`int`),
				})
		})
	})

	//
	// Errors
	t.Run("errors", func(t *testing.T) {
		checkError(t, func() {}, td.JSONPointer("", td.NotNil()),
			expectedError{
				Message: mustBe("json.Marshal failed"),
				Path:    mustBe("DATA"),
				Summary: mustContain("json: unsupported type"),
			})

		checkError(t,
			map[string]int{"zzz": 42},
			td.JSONPointer("/zzz", jsonPtrTest(56)),
			expectedError{
				Message: mustBe("an error occurred while unmarshalling JSON into td_test.jsonPtrTest"),
				Path:    mustBe("DATA.JSONPointer</zzz>"),
				Summary: mustBe("jsonPtrTest unmarshal custom error"),
			})
	})

	//
	// String
	test.EqualStr(t, td.JSONPointer("/x", td.Gt(2)).String(),
		"JSONPointer(/x, > 2)")
	test.EqualStr(t, td.JSONPointer("/x", 2).String(),
		"JSONPointer(/x, 2)")
	test.EqualStr(t, td.JSONPointer("/x", nil).String(),
		"JSONPointer(/x, nil)")

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.JSONPointer("x", 1234) },
		"JSONPointer(): bad JSON pointer x")
}

func TestJSONPointerTypeBehind(t *testing.T) {
	equalTypes(t, td.JSONPointer("", 42), nil)
}
