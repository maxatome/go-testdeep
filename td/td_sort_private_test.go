// Copyright (c) 2024, Maxime Soulé
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

func TestSortFunc(t *testing.T) {
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
			name:     "mkSortStructFields-asc1-1field",
			how:      "b",
			slice:    newSortA(),
			expected: "[{3 1 2} {1 2 3} {1 2 4} {2 3 1} {2 3 2} {2 4 2}]",
		},
		{
			name:     "mkSortStructFields-desc1",
			how:      []string{"-b"},
			slice:    newSortA(),
			expected: "[{2 4 2} {2 3 1} {2 3 2} {1 2 3} {1 2 4} {3 1 2}]",
		},
		{
			name:     "mkSortStructFields-multi",
			how:      []string{"-a", "b", "-c"},
			slice:    newSortA(),
			expected: "[{3 1 2} {2 3 2} {2 3 1} {2 4 2} {1 2 4} {1 2 3}]",
		},
		{
			name: "mkSortStructFields-deref",
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
						if se == nil {
							sl2[i] = "(**int)(nil)"
						} else if *se == nil {
							sl2[i] = "(**int)(&nil)"
						} else {
							sl2[i] = fmt.Sprintf("&&%d", **se)
						}
					case **sortA:
						if se == nil {
							sl2[i] = "(**sortA:)(nil)"
						} else if *se == nil {
							sl2[i] = "(**sortA:)(&nil)"
						} else {
							sl2[i] = **se
						}
					default:
						sl2[i] = fmt.Sprintf("%[1]T(%[1]v)", se)
					}
				}
				return fmt.Sprintf("%v", sl2)
			},
			expected: "[(**sortA:)(nil) (**sortA:)(&nil) {2 4 2} {2 3 1} {2 3 2} {1 2 3} {1 2 4} {3 1 2} nil (**int)(&nil) &18 18]",
		},
		{
			name:     "mkSortStructFields-unknown",
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
			mkFn, err := sortFunc(tc.how)
			test.NoError(t, err)
			fn, err := mkFn(reflect.TypeOf(tc.slice).Elem())
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
			how         any
			expectedErr string
		}{
			{
				name:        "unknown how",
				how:         true,
				expectedErr: "SORT_FUNC must be an int, string, []string or func(T, T) bool",
			},
			{
				name:        "bad func variadic",
				how:         func(...int) bool { return true },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(...int) bool",
			},
			{
				name:        "bad func num in1",
				how:         func(a int) bool { return true },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int) bool",
			},
			{
				name:        "bad func num in3",
				how:         func(a, b, c int) bool { return true },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int, int) bool",
			},
			{
				name:        "bad func in types",
				how:         func(a int, b bool) bool { return true },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, bool) bool",
			},
			{
				name:        "bad func num out0",
				how:         func(a, b int) {},
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int)",
			},
			{
				name:        "bad func num out2",
				how:         func(a, b int) (bool, bool) { return true, true },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int) (bool, bool)",
			},
			{
				name:        "bad func out type",
				how:         func(a, b int) int { return 0 },
				expectedErr: "SORT_FUNC must match func(T, T) bool signature, not func(int, int) int",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mkFn, err := sortFunc(tc.how)
				if mkFn != nil {
					t.Error("sortFunc should return nil function")
				}
				if test.Error(t, err) {
					test.EqualStr(t, err.Error(), tc.expectedErr)
				}
			})
		}

		t.Run("custom func", func(t *testing.T) {
			mkFn, err := sortFunc(func(a, b int) bool { return a > b })
			test.NoError(t, err)
			_, err = mkFn(types.Bool)
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(), "bool is not assignable to int")
			}
		})
	})
}
