// Copyright (c) 2024-2025, Maxime Soulé
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

	"github.com/maxatome/go-testdeep/internal/test"
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
		name        string
		how         any
		got         any
		expected    any
		expectedErr expectedError
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
			name:     "asc float",
			how:      42.3,
			got:      []int{1, -2, -3, 0, -1, 3, 2},
			expected: []int{-3, -2, -1, 0, 1, 2, 3},
		},
		{
			name:     "no items",
			how:      1,
			got:      []int{},
			expected: []int{},
		},
		{
			name:     "one item",
			how:      1,
			got:      []int{23},
			expected: []int{23},
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
			name: "evenHigher",
			how: func(a, b int) bool {
				if (a%2 == 0) != (b%2 == 0) {
					return a%2 != 0
				}
				return a < b
			},
			got:      []int{-1, 1, 2, -3, 3, -2, 0},
			expected: []int{-3, -1, 1, 3, -2, 0, 2},
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
		{
			name:     "invalid fields-path",
			how:      "",
			got:      []int{1, 2},
			expected: []int{42, 42},
			expectedErr: expectedError{
				Message: mustBe("cannot sort items"),
				Path:    mustBe("DATA"),
				Summary: mustBe("FIELDS_PATH cannot be empty"),
				Under:   mustContain("under operator Sort at "),
			},
		},
		{
			name:     "bad usage",
			how:      1,
			got:      []int{1, 2},
			expected: 42,
			expectedErr: expectedError{
				Message: mustBe("bad usage of Sort operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Sort(SORT_FUNC|int|string|[]string, TESTDEEP_OPERATOR|EXPECTED_VALUE), EXPECTED_VALUE must be a slice or an array not a int"),
				Under:   mustContain("under operator Sort at "),
			},
		},
		{
			name:     "grepResolvePtr error",
			got:      (*int)(nil),
			expected: td.Ignore(),
			expectedErr: expectedError{
				Message:  mustBe("nil pointer"),
				Got:      mustBe("nil *int"),
				Expected: mustBe("non-nil *slice OR *array"),
				Under:    mustContain("under operator Sort at "),
			},
		},
		{
			name:     "grepBadKind",
			got:      123,
			expected: td.Ignore(),
			expectedErr: expectedError{
				Message:  mustBe("bad kind"),
				Got:      mustBe("int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
				Under:    mustContain("under operator Sort at "),
			},
		},
		{
			name:     "erroneous operator expected",
			how:      1,
			got:      []int{1, 2},
			expected: td.JSON("{"),
			expectedErr: expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustMatch(`DATA<sorted>`), // always without .Iface
				Summary: mustContain("JSON unmarshal error"),
				Under:   mustContain("under operator JSON at "),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedErr == (expectedError{}) {
				checkOK(t, tc.got, td.Sort(tc.how, tc.expected))
			} else {
				checkError(t, tc.got, td.Sort(tc.how, tc.expected), tc.expectedErr)
			}
		})
	}

	t.Run("JSON", func(t *testing.T) {
		checkOK(t,
			json.RawMessage(`["c","a","b"]`),
			td.JSON(`Sort(1, ["a","b","c"])`))

		checkOK(t,
			json.RawMessage(`{"x": ["c","a","b"]}`),
			td.JSON(`{"x": Sort(-1, ["c","b","a"])}`))

		checkOK(t,
			map[string][]string{"labels": {"c", "a", "b"}},
			td.JSON(`{"labels": Sort(1, ["a", "b", "c"])}`))

		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		got := struct {
			People []Person `json:"people"`
		}{
			People: []Person{
				{"Brian", 22},
				{"Bob", 19},
				{"Stephen", 19},
				{"Alice", 20},
				{"Marcel", 25},
			},
		}
		checkOK(t, got, td.JSON(`{
			"people": Sort("name", [ // sort by name ascending
				{"name": "Alice",   "age": 20},
				{"name": "Bob",     "age": 19},
				{"name": "Brian",   "age": 22},
				{"name": "Marcel",  "age": 25},
				{"name": "Stephen", "age": 19},
			])
		}`))
		checkOK(t, got, td.JSON(`{
			"people": Sort([ "-age", "name" ], [ // sort by age desc, then by name asc
				{"name": "Marcel",  "age": 25},
				{"name": "Brian",   "age": 22},
				{"name": "Alice",   "age": 20},
				{"name": "Bob",     "age": 19},
				{"name": "Stephen", "age": 19},
			])
		}`))
	})
}

func TestSortTypeBehind(t *testing.T) {
	equalTypes(t, td.Sort(1, []int{}), nil)

	// Erroneous op
	equalTypes(t, td.Sort(func() {}, []int{}), nil)
}

func TestSortString(t *testing.T) {
	test.EqualStr(t, td.Sort(nil, []int{}).String(), "Sort(<nil>, ([]int) {\n})")

	test.EqualStr(t, td.Sort(1, []int{}).String(), "Sort(1, ([]int) {\n})")

	test.EqualStr(t,
		td.Sort(func(a, b int) bool { return true }, []int{}).String(),
		"Sort(func(int, int) bool, ([]int) {\n})")

	// Erroneous op
	test.EqualStr(t, td.Sort(func() {}, []int{}).String(), "Sort(<ERROR>)")
}

// nolint: unused
func TestSorted(t *testing.T) {
	firstBecomesLast := func(x any) (any, int) {
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
		name        string
		got         any
		sorted      []any
		expectedErr expectedError
	}{
		{
			name: "slice",
			got:  []int{0, 1, 2, 2},
		},
		{
			name: "*slice",
			got:  &[]int{0, 1, 2, 2},
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
		{
			name:   "invalid fields-path",
			got:    []int{1, 2},
			sorted: []any{""},
			expectedErr: expectedError{
				Message: mustBe("cannot sort items"),
				Path:    mustBe("DATA"),
				Summary: mustBe("FIELDS_PATH cannot be empty"),
				Under:   mustContain("under operator Sorted at "),
			},
		},
		{
			name:   "bad usage",
			got:    []int{1, 2},
			sorted: []any{42, 23},
			expectedErr: expectedError{
				Message: mustBe("bad usage of Sorted operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Sorted(SORT_FUNC|int|[]string|string...), string... expected but received int as 1st parameter"),
				Under:   mustContain("under operator Sorted at "),
			},
		},
		{
			name:   "grepResolvePtr error",
			got:    (*int)(nil),
			sorted: []any{1},
			expectedErr: expectedError{
				Message:  mustBe("nil pointer"),
				Got:      mustBe("nil *int"),
				Expected: mustBe("non-nil *slice OR *array"),
				Under:    mustContain("under operator Sorted at "),
			},
		},
		{
			name:   "grepBadKind",
			got:    123,
			sorted: []any{1},
			expectedErr: expectedError{
				Message:  mustBe("bad kind"),
				Got:      mustBe("int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
				Under:    mustContain("under operator Sorted at "),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedErr == (expectedError{}) {
				checkOK(t, tc.got, td.Sorted(tc.sorted...))

				// expect an error after moving first item at last position
				if got, l := firstBecomesLast(tc.got); got != nil {
					checkError(t, got, td.Sorted(tc.sorted...),
						expectedError{
							Message: mustBe(
								fmt.Sprintf("not sorted, item #%d value is before #%d one while it should not",
									l-1, l-2)),
							Path: mustBe("DATA"),
							Summary: mustMatch(
								fmt.Sprintf(`(?s)^item #%d: .+\nitem #%d: `, l-2, l-1)),
						})
				}
				return
			}

			checkError(t, tc.got, td.Sorted(tc.sorted...), tc.expectedErr)
		})
	}

	t.Run("JSON", func(t *testing.T) {
		checkOK(t,
			json.RawMessage(`["a","b","c"]`),
			td.JSON(`Sorted`))

		checkOK(t,
			json.RawMessage(`{"x": ["c","b","a"]}`),
			td.JSON(`{ "x": Sorted(-1) }`))

		checkOK(t,
			map[string][]string{"labels": {"a", "b", "c"}},
			td.JSON(`{ "labels": Sorted }`))

		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		got := struct {
			People []Person `json:"people"`
		}{
			People: []Person{
				{"Alice", 20},
				{"Bob", 19},
				{"Brian", 22},
				{"Marcel", 25},
				{"Stephen", 19},
			},
		}
		// is sorted by name ascending
		checkOK(t, got, td.JSON(`{ "people": Sorted("name") }`))

		// is sorted by age desc, then by name asc
		checkError(t, got, td.JSON(`{ "people": Sorted("-age", "name") }`),
			expectedError{
				Message: mustBe("not sorted, item #2 value is before #1 one while it should not"),
				Path:    mustBe(`DATA["people"]`),
				Summary: mustMatch(`(?s)^item #1: .+"Bob".*\nitem #2: .+"Brian"`),
				Under:   mustContain("under operator Sorted at "),
			})
	})
}

func TestSortedTypeBehind(t *testing.T) {
	equalTypes(t, td.Sorted(), nil)
	equalTypes(t, td.Sorted(-1), nil)

	// Erroneous op
	equalTypes(t, td.Sorted(func() {}), nil)
}

func TestSortedString(t *testing.T) {
	test.EqualStr(t, td.Sorted(nil).String(), "Sorted(<nil>)")
	test.EqualStr(t, td.Sorted(1).String(), "Sorted(1)")
	test.EqualStr(t, td.Sorted("a.b", "-d.e").String(), `Sorted(a.b, -d.e)`)

	test.EqualStr(t,
		td.Sorted(func(a, b int) bool { return true }).String(),
		"Sorted(func(int, int) bool)")

	// Erroneous op
	test.EqualStr(t, td.Sorted(func() {}).String(), "Sorted(<ERROR>)")
}
