// Copyright (c) 2021 Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
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
