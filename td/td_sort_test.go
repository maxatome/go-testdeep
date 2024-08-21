// Copyright (c) 2024, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestSort(t *testing.T) {
	type sortTest1 struct {
		s string
		a int
		b int
	}
	type sortTest2 struct{ a, b, c int }
	testCases := []struct {
		name     string
		how      any
		got      any
		expected any
	}{
		{
			name:     "slice",
			how:      1,
			got:      []int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "*slice",
			how:      1,
			got:      &[]int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "array",
			how:      1,
			got:      [...]int{1, -2, -3, 0, -1, 3, 2},
			expected: [...]int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "*array",
			how:      1,
			got:      &[...]int{1, -2, -3, 0, -1, 3, 2},
			expected: [...]int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "asc0",
			how:      0,
			got:      []int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "desc",
			how:      -1,
			got:      []int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{3, 2, 1, 0, -1, -2, -3},
		},
		{
			name: "func",
			how: func(a, b int) bool {
				if a == 0 || b == 0 {
					return b == 0
				}
				return a > b
			},
			got:      []int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{3, 2, 1, -1, -2, -3, 0},
		},
		{
			name:     "fields-path",
			how:      "s",
			got:      []sortTest1{{"c", 4, 2}, {"a", 8, 1}, {"b", 0, 3}},
			expected: []sortTest1{{"a", 8, 1}, {"b", 0, 3}, {"c", 4, 2}},
		},
		{
			name:     "multiple fields-paths",
			how:      []string{"a", "-b", "c"},
			got:      []sortTest2{{1, 9, 5}, {2, 0, 0}, {1, 9, 4}, {1, 8, 0}},
			expected: []sortTest2{{1, 9, 4}, {1, 9, 5}, {1, 8, 0}, {2, 0, 0}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkOK(t, tc.got, td.Sort(tc.how, tc.expected))
		})
	}

	t.Run("JSON", func(t *testing.T) {
		checkOK(t,
			json.RawMessage(`["c","a","b"]`),
			td.JSON(`Sort(1, ["a","b","c"])`))

		checkOK(t,
			json.RawMessage(`{"x": ["c","a","b"]}`),
			td.JSON(`{"x": Sort(-1, ["c","b","a"])}`))
	})
}

func TestSortTypeBehind(t *testing.T) {
	equalTypes(t, td.Sort(1, []int{}), nil)

	// Erroneous op
	equalTypes(t, td.Sort(func() {}, []int{}), nil)
}

// nolint: unused
func TestSorted(t *testing.T) {
	lastBecomesFirst := func(x any) (any, int) {
		vx := reflect.ValueOf(x)
		if vx.Kind() == reflect.Ptr {
			vx = vx.Elem()
		}
		l := vx.Len()
		if l <= 1 {
			return nil, 0
		}
		if vx.Kind() == reflect.Array {
			vx2 := reflect.New(vx.Type()).Elem()
			reflect.Copy(vx2, vx)
			vx = vx2
		}
		first := vx.Index(0).Interface()
		reflect.Copy(vx, vx.Slice(1, l))
		vx.Index(l - 1).Set(reflect.ValueOf(first))
		return vx.Interface(), l
	}

	type nested3 struct {
		val   int
		dummy int
	}
	type nested2 struct {
		val int
		n3  *nested3
	}
	type nested1 struct {
		val int
		n2  nested2
	}
	testCases := []struct {
		name   string
		got    any
		sorted []any
	}{
		{
			name: "slice",
			got:  []int{0, 1, 2, 2},
		},
		{
			name: "*slice",
			got:  ptr([]int{0, 1, 2, 2}),
		},
		{
			name: "array",
			got:  [...]int{0, 1, 2, 2},
		},
		{
			name: "*array",
			got:  ptr([...]int{0, 1, 2, 2}),
		},
		{
			name:   "asc",
			got:    []int{0, 1, 2, 2},
			sorted: []any{1},
		},
		{
			name:   "asc",
			got:    []int{4, 3, 2, 2},
			sorted: []any{-1},
		},
		{
			name:   "flatten struct field",
			got:    []struct{ name string }{{"a"}, {"b"}, {"c"}, {"c"}},
			sorted: []any{"name"},
		},
		{
			name:   "struct field in slice",
			got:    []struct{ name string }{{"a"}, {"b"}, {"c"}, {"c"}},
			sorted: []any{[]string{"name"}},
		},
		{
			name:   "flatten struct field desc",
			got:    []struct{ name string }{{"d"}, {"c"}, {"b"}, {"b"}},
			sorted: []any{"-name"},
		},
		{
			name:   "flatten multiple struct fields",
			got:    []struct{ a, b, c int }{{1, 9, 4}, {1, 9, 5}, {1, 8, 0}, {2, 0, 0}},
			sorted: []any{"a", "-b", "c"},
		},
		{
			name: "nested fields-path",
			got: []*nested1{
				{1, nested2{1, &nested3{1, 8}}},
				{1, nested2{1, &nested3{2, 7}}},
				{1, nested2{2, &nested3{2, 6}}},
				{2, nested2{2, &nested3{2, 5}}},
			},
			sorted: []any{"n2.n3.val", "n2.val", "val"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkOK(t, tc.got, td.Sorted(tc.sorted...))

			if got, l := lastBecomesFirst(tc.got); got != nil {
				checkError(t, got, td.Sorted(tc.sorted...),
					expectedError{
						Message: mustBe("not sorted"),
						Path:    mustBe("DATA"),
						Summary: mustBe(fmt.Sprintf(
							"item #%d value is lesser than #%d one while it should not",
							l-1, l-2)),
					})
			}
		})
	}
}

func TestSortedTypeBehind(t *testing.T) {
	equalTypes(t, td.Sorted(), nil)
	equalTypes(t, td.Sorted(-1), nil)

	// Erroneous op
	equalTypes(t, td.Sorted(func() {}), nil)
}
