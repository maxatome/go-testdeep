// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/util"
)

func TestBadParam(t *testing.T) {
	test.EqualStr(t,
		util.BadParam(nil, 1, true),
		"but received nil as 1st parameter")

	test.EqualStr(t,
		util.BadParam(42, 1, true),
		"but received int as 1st parameter")

	test.EqualStr(t,
		util.BadParam([]int{}, 1, true),
		"but received []int (slice) as 1st parameter")
	test.EqualStr(t,
		util.BadParam([]int{}, 1, false),
		"but received []int as 1st parameter")

	test.EqualStr(t,
		util.BadParam(nil, 1, true),
		"but received nil as 1st parameter")
	test.EqualStr(t,
		util.BadParam(nil, 2, true),
		"but received nil as 2nd parameter")
	test.EqualStr(t,
		util.BadParam(nil, 3, true),
		"but received nil as 3rd parameter")
	test.EqualStr(t,
		util.BadParam(nil, 4, true),
		"but received nil as 4th parameter")
}

func TestTern(t *testing.T) {
	test.EqualInt(t, int(util.TernRune(true, 'A', 'B')), int('A'))
	test.EqualInt(t, int(util.TernRune(false, 'A', 'B')), int('B'))
}
