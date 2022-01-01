// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

type compareType int

func (i compareType) Compare(j compareType) int {
	if i < j {
		return -1
	}
	if i > j {
		return 1
	}
	return 0
}

type lessType int

func (i lessType) Less(j lessType) bool {
	return i < j
}

type badType1 int

func (i badType1) Compare(j ...badType1) int { return 0 }     // IsVariadic()
func (i badType1) Less(j, k badType1) bool   { return false } // NumIn() == 3

type badType2 int

func (i badType2) Compare() int    { return 0 } // NumIn() == 1
func (i badType2) Less(j badType2) {}           // NumOut() == 0

type badType3 int

func (i badType3) Compare(j badType3) (int, int) { return 0, 0 }  // NumOut() == 2
func (i badType3) Less(j int) bool               { return false } // In(1) ≠ in

type badType4 int

func (i badType4) Compare(j badType4) bool { return false } // Out(0) ≠ out
func (i badType4) Less(j badType4) int     { return 0 }     // Out(0) ≠ out

func TestOrder(t *testing.T) {
	if types.NewOrder(reflect.TypeOf(0)) != nil {
		t.Error("types.NewOrder(int) returned non-nil func")
	}

	fn := types.NewOrder(reflect.TypeOf(compareType(0)))
	if fn == nil {
		t.Error("types.NewOrder(compareType) returned nil func")
	} else {
		a, b := reflect.ValueOf(compareType(1)), reflect.ValueOf(compareType(2))
		test.EqualInt(t, fn(a, b), -1)
		test.EqualInt(t, fn(b, a), 1)
		test.EqualInt(t, fn(a, a), 0)
	}

	fn = types.NewOrder(reflect.TypeOf(lessType(0)))
	if fn == nil {
		t.Error("types.NewOrder(lessType) returned nil func")
	} else {
		a, b := reflect.ValueOf(lessType(1)), reflect.ValueOf(lessType(2))
		test.EqualInt(t, fn(a, b), -1)
		test.EqualInt(t, fn(b, a), 1)
		test.EqualInt(t, fn(a, a), 0)
	}

	if types.NewOrder(reflect.TypeOf(badType1(0))) != nil {
		t.Error("types.NewOrder(badType1) returned non-nil func")
	}
	if types.NewOrder(reflect.TypeOf(badType2(0))) != nil {
		t.Error("types.NewOrder(badType2) returned non-nil func")
	}
	if types.NewOrder(reflect.TypeOf(badType3(0))) != nil {
		t.Error("types.NewOrder(badType3) returned non-nil func")
	}
	if types.NewOrder(reflect.TypeOf(badType4(0))) != nil {
		t.Error("types.NewOrder(badType4) returned non-nil func")
	}
}
