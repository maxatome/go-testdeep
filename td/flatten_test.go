// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestFlatten(t *testing.T) {
	fs := td.Flatten([]int{1, 2, 3})
	if s, ok := fs.Slice.([]int); test.IsTrue(t, ok) {
		test.EqualInt(t, len(s), 3)
	}

	fs = td.Flatten([...]int{1, 2, 3})
	_, ok := fs.Slice.([3]int)
	test.IsTrue(t, ok)

	test.CheckPanic(t, func() { td.Flatten(nil) }, "usage: Flatten(SLICE|ARRAY)")
	test.CheckPanic(t, func() { td.Flatten(42) }, "usage: Flatten(SLICE|ARRAY)")
}
