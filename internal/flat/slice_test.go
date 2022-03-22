// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package flat_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestLen(t *testing.T) {
	num, flattened := flat.Len([]any{1, 2, 3, 4})
	test.EqualInt(t, num, 4)
	test.IsTrue(t, flattened)

	num, flattened = flat.Len([]any{
		1, 2,
		flat.Slice{Slice: []int{3, 4, 5, 6}},
		flat.Slice{Slice: map[int]int{-1: -2, -3: -4}},
		7,
		flat.Slice{
			Slice: []any{
				flat.Slice{Slice: []int{8, 9}},
				flat.Slice{Slice: []int{10, 11}},
				flat.Slice{Slice: map[int]any{
					-5: -6,
					-7: flat.Slice{Slice: []int{-8, -9, -10}},
				}},
			},
		},
		12,
		flat.Slice{Slice: map[any]any{
			-11: flat.Slice{Slice: []int{-12, -13}},
		}},
	})
	test.EqualInt(t, num, 12+13)
	test.IsFalse(t, flattened)
}

func TestValues(t *testing.T) {
	sv := flat.Values(nil)
	test.IsTrue(t, sv != nil)
	test.EqualInt(t, len(sv), 0)

	sv = flat.Values([]any{1, 2})
	if test.EqualInt(t, len(sv), 2) {
		test.EqualInt(t, int(sv[0].Int()), 1)
		test.EqualInt(t, int(sv[1].Int()), 2)
	}

	sv = flat.Values([]any{
		1, 2,
		flat.Slice{Slice: []int{3, 4, 5, 6}},
		7,
		flat.Slice{
			Slice: []any{
				flat.Slice{Slice: []int{8, 9}},
				flat.Slice{Slice: []any{10, 11}},
				12,
				13,
			},
		},
		14,
		flat.Slice{
			Slice: map[int]any{
				15: flat.Slice{Slice: map[int]int{16: 17}},
			},
		},
	})
	if test.EqualInt(t, len(sv), 17) {
		for i, v := range sv {
			test.EqualInt(t, int(v.Int()), i+1)
		}
	}
}

func TestInterfaces(t *testing.T) {
	si := flat.Interfaces()
	test.IsTrue(t, si == nil)

	si = flat.Interfaces(1, 2)
	if test.EqualInt(t, len(si), 2) {
		test.EqualInt(t, si[0].(int), 1)
		test.EqualInt(t, si[1].(int), 2)
	}

	si = flat.Interfaces(
		1, 2,
		flat.Slice{Slice: []int{3, 4, 5, 6}},
		7,
		flat.Slice{
			Slice: []any{
				flat.Slice{Slice: []int{8, 9}},
				flat.Slice{Slice: []any{10, 11}},
				12,
				13,
			},
		},
		14,
		flat.Slice{
			Slice: map[int]any{
				15: flat.Slice{Slice: map[int]int{16: 17}},
			},
		},
	)
	if test.EqualInt(t, len(si), 17) {
		for i, iface := range si {
			test.EqualInt(t, iface.(int), i+1)
		}
	}
}
