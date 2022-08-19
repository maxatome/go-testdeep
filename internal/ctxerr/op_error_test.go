// Copyright (c) 2021-2022 Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

const prefix = ": bad usage of Zzz operator\n\t"

func TestOpBadUsage(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", nil, 1, true).Error(),
		prefix+"usage: Zzz(STRING), but received nil as 1st parameter")

	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", 42, 1, true).Error(),
		prefix+"usage: Zzz(STRING), but received int as 1st parameter")

	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", []int{}, 1, true).Error(),
		prefix+"usage: Zzz(STRING), but received []int (slice) as 1st parameter")
	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", []int{}, 1, false).Error(),
		prefix+"usage: Zzz(STRING), but received []int as 1st parameter")

	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", nil, 1, true).Error(),
		prefix+"usage: Zzz(STRING), but received nil as 1st parameter")
	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", nil, 2, true).Error(),
		prefix+"usage: Zzz(STRING), but received nil as 2nd parameter")
	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", nil, 3, true).Error(),
		prefix+"usage: Zzz(STRING), but received nil as 3rd parameter")
	test.EqualStr(t,
		ctxerr.OpBadUsage("Zzz", "(STRING)", nil, 4, true).Error(),
		prefix+"usage: Zzz(STRING), but received nil as 4th parameter")
}

func TestOpTooManyParams(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t, ctxerr.OpTooManyParams("Zzz", "(PARAM)").Error(),
		prefix+"usage: Zzz(PARAM), too many parameters")
}

func TestBad(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t,
		ctxerr.OpBad("Zzz", "test").Error(),
		prefix+"test")

	test.EqualStr(t,
		ctxerr.OpBad("Zzz", "test %d", 123).Error(),
		prefix+"test 123")
}

func TestBadKind(t *testing.T) {
	defer color.SaveState()()

	expected := func(got string) string {
		return ": bad kind\n\t     got: " + got + "\n\texpected: some kinds"
	}

	test.EqualStr(t,
		ctxerr.BadKind(reflect.ValueOf(42), "some kinds").Error(),
		expected("int"))

	test.EqualStr(t,
		ctxerr.BadKind(reflect.ValueOf(&[]int{}), "some kinds").Error(),
		expected("*slice (*[]int type)"))

	test.EqualStr(t,
		ctxerr.BadKind(reflect.ValueOf((***int)(nil)), "some kinds").Error(),
		expected("***int"))

	test.EqualStr(t,
		ctxerr.BadKind(reflect.ValueOf(nil), "some kinds").Error(),
		expected("nil"))
}

func TestNilPointer(t *testing.T) {
	defer color.SaveState()()

	expected := func(got string) string {
		return ": nil pointer\n\t     got: nil " + got + "\n\texpected: non-nil blah blah"
	}

	test.EqualStr(t,
		ctxerr.NilPointer(reflect.ValueOf((*int)(nil)), "non-nil blah blah").Error(),
		expected("*int"))

	test.EqualStr(t,
		ctxerr.NilPointer(reflect.ValueOf((*[]int)(nil)), "non-nil blah blah").Error(),
		expected("*slice (*[]int type)"))

	test.EqualStr(t,
		ctxerr.NilPointer(reflect.ValueOf((***int)(nil)), "non-nil blah blah").Error(),
		expected("***int"))

	test.EqualStr(t,
		ctxerr.NilPointer(reflect.ValueOf(nil), "non-nil blah blah").Error(),
		expected("nil"))
}
