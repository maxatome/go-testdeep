// Copyright (c) 2020-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestFlatten(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testCases := []struct {
			name         string
			sliceOrMap   any
			fn           []any
			expectedType reflect.Type
			expectedLen  int
		}{
			{
				name:         "slice",
				sliceOrMap:   []int{1, 2, 3},
				expectedType: reflect.TypeOf([]int{}),
				expectedLen:  3,
			},
			{
				name:         "array",
				sliceOrMap:   [3]int{1, 2, 3},
				expectedType: reflect.TypeOf([3]int{}),
				expectedLen:  3,
			},
			{
				name:         "map",
				sliceOrMap:   map[int]int{1: 2, 3: 4},
				expectedType: reflect.TypeOf(map[int]int{}),
				expectedLen:  2,
			},
			{
				name:         "slice+untyped nil fn",
				sliceOrMap:   []int{1, 2, 3},
				fn:           []any{nil},
				expectedType: reflect.TypeOf([]int{}),
				expectedLen:  3,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				s := td.Flatten(tc.sliceOrMap, tc.fn...)
				if reflect.TypeOf(s.Slice) != tc.expectedType {
					t.Errorf("types differ: got=%s, expected=%s",
						reflect.TypeOf(s.Slice), tc.expectedType)
					return
				}
				test.EqualInt(t, reflect.ValueOf(s.Slice).Len(), tc.expectedLen)
			})
		}
	})

	t.Run("ok+func", func(t *testing.T) {
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
			usage     = `usage: Flatten(SLICE|ARRAY|MAP[, FUNC])`
			usageFunc = usage + `, FUNC should be non-nil func(T) V or func(T) (V, bool) or a string "Smuggle:…" or "JSONPointer:…"`
		)
		testCases := []struct {
			name       string
			fn         []any
			sliceOrMap any
			expected   string
		}{
			{
				name:       "too many params",
				sliceOrMap: []int{},
				fn:         []any{1, 2},
				expected:   usage + ", too many parameters",
			},
			{
				name:     "nil sliceOrMap",
				expected: usage + ", but received nil as 1st parameter",
			},
			{
				name:       "bad sliceOrMap type",
				sliceOrMap: 42,
				expected:   usage + ", but received int as 1st parameter",
			},
			{
				name:       "not func",
				sliceOrMap: []int{},
				fn:         []any{42},
				expected:   usageFunc + ", but received int as 2nd parameter",
			},
			{
				name:       "func w/0 inputs",
				sliceOrMap: []int{},
				fn:         []any{func() int { return 0 }},
				expected:   usageFunc + ", but received func() int as 2nd parameter",
			},
			{
				name:       "func w/2 inputs",
				sliceOrMap: []int{},
				fn:         []any{func(a, b int) int { return 0 }},
				expected:   usageFunc + ", but received func(int, int) int as 2nd parameter",
			},
			{
				name:       "variadic func",
				sliceOrMap: []int{},
				fn:         []any{func(a ...int) int { return 0 }},
				expected:   usageFunc + ", but received func(...int) int as 2nd parameter",
			},
			{
				name:       "func w/0 output",
				sliceOrMap: []int{},
				fn:         []any{func(a int) {}},
				expected:   usageFunc + ", but received func(int) as 2nd parameter",
			},
			{
				name:       "func w/2 out without bool",
				sliceOrMap: []int{},
				fn:         []any{func(a int) (int, int) { return 0, 0 }},
				expected:   usageFunc + ", but received func(int) (int, int) as 2nd parameter",
			},
			{
				name:       "bad shortcut",
				sliceOrMap: []int{},
				fn:         []any{"Pipo"},
				expected:   usageFunc + `, but received "Pipo" as 2nd parameter`,
			},
			{
				name:       "typed nil func",
				sliceOrMap: []int{},
				fn:         []any{(func(a int) int)(nil)},
				expected:   usageFunc,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				test.CheckPanic(t, func() { td.Flatten(tc.sliceOrMap, tc.fn...) }, tc.expected)
			})
		}
	})
}
