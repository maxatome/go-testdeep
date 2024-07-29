// Copyright (c) 2024-2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

type sortA struct{ a, b, c int }

type sortB struct {
	name string
	idx  int
}

func (b sortB) Compare(x sortB) int {
	bn, xn := strings.ToLower(b.name), strings.ToLower(x.name)
	if bn == xn {
		return b.idx - x.idx
	}
	if bn < xn {
		return -1
	}
	return 1
}

func TestInitSortBase(t *testing.T) {
	newSortA := func() []sortA {
		return []sortA{
			{2, 3, 2},
			{1, 2, 3},
			{3, 1, 2},
			{2, 4, 2},
			{1, 2, 4},
			{2, 3, 1},
		}
	}

	testCases := []struct {
		name     string
		how      any
		slice    any
		got2str  func(any) string
		expected string
	}{
		{
			name:     "mkSortAsc",
			how:      0,
			slice:    []int{3, 5, 2, 1, 4},
			expected: "[1 2 3 4 5]",
		},
		{
			name:     "mkSortDesc",
			how:      -1,
			slice:    []int{3, 5, 2, 1, 4},
			expected: "[5 4 3 2 1]",
		},
		{
			name:     "mkSortAsc-float64",
			how:      float64(1),
			slice:    []int{3, 5, 2, 1, 4},
			expected: "[1 2 3 4 5]",
		},
		{
			name:     "mkSortDesc-float64",
			how:      float64(-1),
			slice:    []int{3, 5, 2, 1, 4},
			expected: "[5 4 3 2 1]",
		},
		{
			name:     "mkSortAsc-Compare",
			how:      1,
			slice:    []sortB{{"Zb", 1}, {"za", 8}, {"a", 4}, {"A", 5}},
			expected: "[{a 4} {A 5} {za 8} {Zb 1}]",
		},
		{
			name:     "mkSortFieldsPaths-asc1-1field",
			how:      "b",
			slice:    newSortA(),
			expected: "[{3 1 2} {1 2 3} {1 2 4} {2 3 1} {2 3 2} {2 4 2}]",
		},
		{
			name:     "mkSortFieldsPaths-asc1-1field-plus",
			how:      "+b",
			slice:    newSortA(),
			expected: "[{3 1 2} {1 2 3} {1 2 4} {2 3 1} {2 3 2} {2 4 2}]",
		},
		{
			name:     "mkSortFieldsPaths-desc1",
			how:      []string{"-b"},
			slice:    newSortA(),
			expected: "[{2 4 2} {2 3 1} {2 3 2} {1 2 3} {1 2 4} {3 1 2}]",
		},
		{
			name:     "mkSortFieldsPaths-multi",
			how:      []string{"-a", "b", "-c"},
			slice:    newSortA(),
			expected: "[{3 1 2} {2 3 2} {2 3 1} {2 4 2} {1 2 4} {1 2 3}]",
		},
		{
			name: "mkSortFieldsPaths-deref",
			how:  []string{"-b"},
			slice: func() []any {
				sl := newSortA()
				var res []any
				for _, v := range sl {
					a := v
					b := &a
					res = append(res, &b)
				}
				var ps *sortA
				n := 18
				var pn *int
				return append(res, &ps, &pn, (**sortA)(nil), n, nil, &n)
			}(),
			got2str: func(got any) string {
				sl := got.([]any)
				sl2 := make([]any, len(sl))
				for i, e := range sl {
					switch se := e.(type) {
					case nil:
						sl2[i] = "nil"
					case int:
						sl2[i] = e
					case *int:
						if se == nil {
							sl2[i] = "(*int)(nil)"
						} else {
							sl2[i] = fmt.Sprintf("&%d", *se)
						}
					case **int:
						switch {
						case se == nil:
							sl2[i] = "(**int)(nil)"
						case *se == nil:
							sl2[i] = "(**int)(&nil)"
						default:
							sl2[i] = fmt.Sprintf("&&%d", **se)
						}
					case **sortA:
						switch {
						case se == nil:
							sl2[i] = "(**sortA:)(nil)"
						case *se == nil:
							sl2[i] = "(**sortA:)(&nil)"
						default:
							sl2[i] = **se
						}
					default:
						sl2[i] = fmt.Sprintf("%[1]T(%[1]v)", se)
					}
				}
				return fmt.Sprintf("%v", sl2)
			},
			expected: "[{2 4 2} {2 3 1} {2 3 2} {1 2 3} {1 2 4} {3 1 2} nil (**int)(&nil) (**sortA:)(nil) (**sortA:)(&nil) &18 18]",
		},
		{
			name:     "mkSortFieldsPaths-unknown",
			how:      []string{"a", "unknown"},
			slice:    newSortA(),
			expected: "[{1 2 3} {1 2 4} {2 3 1} {2 3 2} {2 4 2} {3 1 2}]",
		},
		{
			name:     "custom func",
			how:      func(a, b int) bool { return a > b },
			slice:    []int{3, 5, 2, 1, 4},
			expected: "[5 4 3 2 1]",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sb tdSortBase
			err := sb.initSortBase(tc.how)
			test.NoError(t, err)
			fn, err := sb.mkSortFn(reflect.TypeOf(tc.slice).Elem())
			test.NoError(t, err)
			sl := reflect.ValueOf(tc.slice)
			sort.SliceStable(tc.slice, func(i, j int) bool {
				return fn.Call([]reflect.Value{sl.Index(i), sl.Index(j)})[0].Bool()
			})
			var got string
			if tc.got2str != nil {
				got = tc.got2str(tc.slice)
			} else {
				got = fmt.Sprintf("%v", tc.slice)
			}
			test.EqualStr(t, got, tc.expected)
		})
	}

	t.Run("Error", func(t *testing.T) {
		testCases := []struct {
			name        string
			how         []any
			expectedErr string
		}{
			{
				name:        "not a string",
				how:         []any{"ok", false},
				expectedErr: "string... expected but received bool as 2nd parameter",
			},
			{
				name:        "any slice as used in JSON",
				how:         []any{[]any{"zzz", true}},
				expectedErr: "slice of strings expected as how, bool encountered at pos 1",
			},
			{
				name:        "unknown how",
				how:         []any{true},
				expectedErr: "but received bool as 1st parameter",
			},
			{
				name:        "bad func variadic",
				how:         []any{func(...int) bool { return true }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(...int) bool",
			},
			{
				name:        "bad func num in1",
				how:         []any{func(a int) bool { return true }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int) bool",
			},
			{
				name:        "bad func num in3",
				how:         []any{func(a, b, c int) bool { return true }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int, int) bool",
			},
			{
				name:        "bad func in types",
				how:         []any{func(a int, b bool) bool { return true }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, bool) bool",
			},
			{
				name:        "bad func num out0",
				how:         []any{func(a, b int) {}},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int)",
			},
			{
				name:        "bad func num out2",
				how:         []any{func(a, b int) (bool, bool) { return true, true }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int) (bool, bool)",
			},
			{
				name:        "bad func out type",
				how:         []any{func(a, b int) int { return 0 }},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int) int",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var sb tdSortBase
				err := sb.initSortBase(tc.how...)
				if sb.mkSortFn != nil {
					t.Error("sortFunc should return nil function")
				}
				if test.Error(t, err) {
					test.EqualStr(t, err.Error(), tc.expectedErr)
				}
			})
		}

		t.Run("custom func", func(t *testing.T) {
			var sb tdSortBase
			err := sb.initSortBase(func(a, b int) bool { return a > b })
			test.NoError(t, err)
			_, err = sb.mkSortFn(types.Bool)
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(), "bool is not assignable to int")
			}
		})
	})
}
