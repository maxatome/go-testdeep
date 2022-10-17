// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.20
// +build go1.20

package types_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

// go1.17 allows to convert []T to *[n]T.
// go1.20 allows to convert []T to [n]T.
func TestIsTypeOrConvertible_go117(t *testing.T) {
	type ArrP *[5]int
	type Arr [5]int

	// 1.17
	ok, convertible := types.IsTypeOrConvertible(
		reflect.ValueOf([]int{1, 2, 3, 4, 5}),
		reflect.TypeOf((ArrP)(nil)))
	test.IsTrue(t, ok)
	test.IsTrue(t, convertible)

	// 1.20
	ok, convertible = types.IsTypeOrConvertible(
		reflect.ValueOf([]int{1, 2, 3, 4, 5}),
		reflect.TypeOf([5]int{}))
	test.IsTrue(t, ok)
	test.IsTrue(t, convertible)

	// 1.20
	ok, convertible = types.IsTypeOrConvertible(
		reflect.ValueOf([]int{1, 2, 3, 4, 5}),
		reflect.TypeOf(Arr{}))
	test.IsTrue(t, ok)
	test.IsTrue(t, convertible)

	ok, convertible = types.IsTypeOrConvertible(
		reflect.ValueOf([]int{1, 2, 3, 4}), // not enough items
		reflect.TypeOf((ArrP)(nil)))
	test.IsFalse(t, ok)
	test.IsFalse(t, convertible)

	ok, convertible = types.IsTypeOrConvertible(
		reflect.ValueOf([]int{1, 2, 3, 4, 5}),
		reflect.TypeOf(&struct{}{}))
	test.IsFalse(t, ok)
	test.IsFalse(t, convertible)
}
