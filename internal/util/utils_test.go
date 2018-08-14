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

func TestTern(t *testing.T) {
	test.EqualStr(t, util.TernStr(true, "A", "B"), "A")
	test.EqualStr(t, util.TernStr(false, "A", "B"), "B")

	test.EqualInt(t, int(util.TernRune(true, 'A', 'B')), int('A'))
	test.EqualInt(t, int(util.TernRune(false, 'A', 'B')), int('B'))
}
