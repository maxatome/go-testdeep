// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
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

	fs = td.Flatten(map[int]int{1: 2, 3: 4})
	if s, ok := fs.Slice.(map[int]int); test.IsTrue(t, ok) {
		test.EqualInt(t, len(s), 2)
	}

	dark.CheckFatalizerBarrierErr(t, func() { td.Flatten(nil) },
		"usage: Flatten(SLICE|ARRAY|MAP), but received nil as 1st parameter")
	dark.CheckFatalizerBarrierErr(t, func() { td.Flatten(42) },
		"usage: Flatten(SLICE|ARRAY|MAP), but received int as 1st parameter")
}
