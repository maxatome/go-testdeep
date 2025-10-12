// Copyright (c) 2020-2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestFlatten(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testCases := []struct {
			name            string
			sliceOrMapOrInt any
			fn              []any
			expectedType    reflect.Type
			expectedLen     int
		}{
			{
				name:            "slice",
				sliceOrMapOrInt: []int{1, 2, 3},
				expectedType:    reflect.TypeOf([]int{}),
				expectedLen:     3,
			},
			{
				name:            "array",
				sliceOrMapOrInt: [3]int{1, 2, 3},
				expectedType:    reflect.TypeOf([3]int{}),
				expectedLen:     3,
			},
			{
				name:            "map",
				sliceOrMapOrInt: map[int]int{1: 2, 3: 4},
				expectedType:    reflect.TypeOf(map[int]int{}),
				expectedLen:     4,
			},
			{
				name:            "slice+untyped nil fn",
				sliceOrMapOrInt: []int{1, 2, 3},
				fn:              []any{nil},
				expectedType:    reflect.TypeOf([]int{}),
				expectedLen:     3,
			},
			{
				name:            "int",
				sliceOrMapOrInt: 5,
				expectedType:    reflect.TypeOf(42),
				expectedLen:     5,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				s := td.Flatten(tc.sliceOrMapOrInt, tc.fn...)
				if reflect.TypeOf(s.Slice) != tc.expectedType {
					t.Errorf("types differ: got=%s, expected=%s",
						reflect.TypeOf(s.Slice), tc.expectedType)
					return
				}
				l, _ := flat.Len([]any{s})
				test.EqualInt(t, l, tc.expectedLen)
			})
		}
	})

	cmp := func(t *testing.T, got, expected []any) {
		t.Helper()

		if (got == nil) != (expected == nil) {
			t.Errorf("nil mismatch: got=%#v, expected=%#v", got, expected)
			return
		}

		lg, le := len(got), len(expected)
		l := lg
		if l > le {
			l = le
		}
		i := 0
		for ; i < l; i++ {
			if got[i] != expected[i] {
				t.Errorf("#%d item differ, got=%v, expected=%v", i, got[i], expected[i])
			}
		}
		for ; i < lg; i++ {
			t.Errorf("#%d item is extra, got=%v", i, got[i])
		}
		for ; i < le; i++ {
			t.Errorf("#%d item is missing, expected=%v", i, expected[i])
		}
	}

	t.Run("ok+func", func(t *testing.T) {
		testCases := []struct {
			name     string
			fn       any
			expected []any
		}{
			{
				name:     "func never called",
				fn:       func(s bool) bool { return true },
				expected: nil,
			},
			{
				name:     "double",
				fn:       func(a int) int { return a * 2 },
				expected: []any{0, 2, 4, 6, 8, 10, 12, 14, 16, 18},
			},
			{
				name:     "even",
				fn:       func(a int) (int, bool) { return a, a%2 == 0 },
				expected: []any{0, 2, 4, 6, 8},
			},
			{
				name:     "transform",
				fn:       func(a int) (string, bool) { return strconv.Itoa(a), a%2 == 0 },
				expected: []any{"0", "2", "4", "6", "8"},
			},
			{
				name:     "nil",
				fn:       func(a int) any { return nil },
				expected: []any{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			{
				name: "convertible",
				fn:   func(a int8) int8 { return a * 3 },
				expected: []any{
					int8(0), int8(3), int8(6), int8(9), int8(12),
					int8(15), int8(18), int8(21), int8(24), int8(27),
				},
			},
			{
				name: "any+variadic",
				fn: func(a any, opts ...bool) int {
					x, _ := a.(int)
					return x
				},
				expected: []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				s := td.Flatten([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, tc.fn)
				if sa, ok := s.Slice.([]any); test.IsTrue(t, ok) {
					cmp(t, sa, tc.expected)
				}
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    int
			fn       any
			expected []any
		}{
			{
				name:     "classic",
				input:    5,
				fn:       func(a int) int { return -a },
				expected: []any{0, -1, -2, -3, -4},
			},
			{
				name:     "convertible to int64",
				input:    5,
				fn:       func(a int64) int64 { return a * 1000 },
				expected: []any{int64(0), int64(1000), int64(2000), int64(3000), int64(4000)},
			},
			{
				name:  "empty",
				input: 0,
				fn:    func(a int) int { return a },
			},
			{
				name:  "empty coz neg",
				input: -42,
				fn:    func(a int) int { return a },
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				s := td.Flatten(tc.input, tc.fn)
				if sa, ok := s.Slice.([]any); test.IsTrue(t, ok) {
					cmp(t, sa, tc.expected)
				}
			})
		}
	})

	t.Run("complex", func(t *testing.T) {
		type person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		got := []person{{"alice", 22}, {"bob", 18}, {"brian", 34}, {"britt", 32}}

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"},
				func(name string) any { return td.Smuggle("Name", name) })))

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"}, "Smuggle:Name")))

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"},
				func(name string) any { return td.JSONPointer("/name", name) })))

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"}, "JSONPointer:/name")))

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"},
				func(name string) any { return td.SuperJSONOf(`{"name":$1}`, name) })))

		td.Cmp(t, got,
			td.Bag(td.Flatten(
				[]string{"alice", "britt", "brian", "bob"},
				func(name string) any { return td.Struct(person{Name: name}) })))
	})

	t.Run("errors", func(t *testing.T) {
		const (
			usage     = `usage: Flatten(SLICE|ARRAY|MAP|int[, FUNC])`
			usageFunc = usage + `, FUNC should be non-nil func(T) V or func(T) (V, bool) or a string "Smuggle:…" or "JSONPointer:…"`
		)
		testCases := []struct {
			name            string
			fn              []any
			sliceOrMapOrInt any
			expected        string
		}{
			{
				name:            "too many params",
				sliceOrMapOrInt: []int{},
				fn:              []any{1, 2},
				expected:        usage + ", too many parameters",
			},
			{
				name:     "nil sliceOrMapOrInt",
				expected: usage + ", but received nil as 1st parameter",
			},
			{
				name:            "bad sliceOrMapOrInt type",
				sliceOrMapOrInt: "42",
				expected:        usage + ", but received string as 1st parameter",
			},
			{
				name:            "not func",
				sliceOrMapOrInt: []int{},
				fn:              []any{42},
				expected:        usageFunc + ", but received int as 2nd parameter",
			},
			{
				name:            "func w/0 inputs",
				sliceOrMapOrInt: []int{},
				fn:              []any{func() int { return 0 }},
				expected:        usageFunc + ", but received func() int as 2nd parameter",
			},
			{
				name:            "func w/2 inputs",
				sliceOrMapOrInt: []int{},
				fn:              []any{func(a, b int) int { return 0 }},
				expected:        usageFunc + ", but received func(int, int) int as 2nd parameter",
			},
			{
				name:            "variadic func",
				sliceOrMapOrInt: []int{},
				fn:              []any{func(a ...int) int { return 0 }},
				expected:        usageFunc + ", but received func(...int) int as 2nd parameter",
			},
			{
				name:            "func w/0 output",
				sliceOrMapOrInt: []int{},
				fn:              []any{func(a int) {}},
				expected:        usageFunc + ", but received func(int) as 2nd parameter",
			},
			{
				name:            "func w/2 out without bool",
				sliceOrMapOrInt: []int{},
				fn:              []any{func(a int) (int, int) { return 0, 0 }},
				expected:        usageFunc + ", but received func(int) (int, int) as 2nd parameter",
			},
			{
				name:            "bad shortcut",
				sliceOrMapOrInt: []int{},
				fn:              []any{"Pipo"},
				expected:        usageFunc + `, but received "Pipo" as 2nd parameter`,
			},
			{
				name:            "typed nil func",
				sliceOrMapOrInt: []int{},
				fn:              []any{(func(a int) int)(nil)},
				expected:        usageFunc,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				test.CheckPanic(t, func() { td.Flatten(tc.sliceOrMapOrInt, tc.fn...) }, tc.expected)
			})
		}
	})
}
