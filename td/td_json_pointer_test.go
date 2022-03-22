// Copyright (c) 2020-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type jsonPtrTest int

func (j jsonPtrTest) UnmarshalJSON(b []byte) error {
	return errors.New("jsonPtrTest unmarshal custom error")
}

type jsonPtrMap map[string]any

func (j *jsonPtrMap) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*map[string]any)(j))
}

var _ = []json.Unmarshaler{jsonPtrTest(0), &jsonPtrMap{}}

func TestJSONPointer(t *testing.T) {
	//
	// nil
	t.Run("nil", func(t *testing.T) {
		checkOK(t, nil, td.JSONPointer("", nil))
		checkOK(t, (*int)(nil), td.JSONPointer("", nil))

		// Yes encoding/json succeeds to unmarshal nil into an int
		checkOK(t, nil, td.JSONPointer("", 0))
		checkOK(t, (*int)(nil), td.JSONPointer("", 0))

		checkError(t, map[string]int{"foo": 42}, td.JSONPointer("/foo", nil),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA.JSONPointer</foo>"),
				Got:      mustBe(`42.0`),
				Expected: mustBe(`nil`),
			})

		// As encoding/json succeeds to unmarshal nil into an int
		checkError(t, map[string]any{"foo": nil}, td.JSONPointer("/foo", 1),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA.JSONPointer</foo>"),
				Got:      mustBe(`0`), // as an int is expected, nil becomes 0
				Expected: mustBe(`1`),
			})
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
			map[string]any{
				"slice": []any{"bar", "baz"},
				"map": map[string]any{
					"test": map[string]any{
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
		checkOK(t, got, td.JSONPointer("/map/test", map[string]any{
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
			got := jsonPtrMap{"x": 123}
			checkOK(t, got, td.JSONPointer("", jsonPtrMap{"x": float64(123)}))
			checkOK(t, got, td.JSONPointer("", &jsonPtrMap{"x": float64(123)}))
			checkOK(t, got, td.JSONPointer("", got))
			checkOK(t, got, td.JSONPointer("", &got))
		})

		t.Run("struct", func(t *testing.T) {
			type jpStruct struct {
				Num any
			}
			got := jpStruct{Num: 123}
			checkOK(t, got, td.JSONPointer("", jpStruct{Num: float64(123)}))
			checkOK(t, jpStruct{Num: got}, td.JSONPointer("/Num", jpStruct{Num: float64(123)}))
			checkOK(t, got, td.JSONPointer("", got))
			checkOK(t, got, td.JSONPointer("", &got))

			expected := int8(123)
			checkOK(t, got, td.JSONPointer("/Num", expected))
			checkOK(t, got, td.JSONPointer("/Num", &expected))
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
	// Bad usage
	checkError(t, "never tested",
		td.JSONPointer("x", 1234),
		expectedError{
			Message: mustBe("bad usage of JSONPointer operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`bad JSON pointer "x"`),
		})

	//
	// String
	test.EqualStr(t, td.JSONPointer("/x", td.Gt(2)).String(),
		"JSONPointer(/x, > 2)")
	test.EqualStr(t, td.JSONPointer("/x", 2).String(),
		"JSONPointer(/x, 2)")
	test.EqualStr(t, td.JSONPointer("/x", nil).String(),
		"JSONPointer(/x, nil)")

	// Erroneous op
	test.EqualStr(t, td.JSONPointer("x", 1234).String(), "JSONPointer(<ERROR>)")
}

func TestJSONPointerTypeBehind(t *testing.T) {
	equalTypes(t, td.JSONPointer("", 42), nil)

	// Erroneous op
	equalTypes(t, td.JSONPointer("x", 1234), nil)
}
