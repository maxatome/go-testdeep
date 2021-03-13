// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package trace_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/trace"
)

func TestStack(t *testing.T) {
	s := trace.Stack{
		{Package: "A", Func: "Aaa.func1"},
		{Package: "A", Func: "Aaa.func2"},
		{Package: "B", Func: "Bbb"},
		{Package: "C", Func: "Ccc"},
	}

	test.IsFalse(t, s.Match(100, "A"))
	test.IsFalse(t, s.Match(-100, "A"))

	test.IsFalse(t, s.Match(3, "B"))
	test.IsFalse(t, s.Match(-1, "B"))

	test.IsTrue(t, s.Match(3, "C"))
	test.IsTrue(t, s.Match(-1, "C"))

	test.IsFalse(t, s.Match(1, "A", "Aaa.func3", "Aaa.func1"))
	test.IsTrue(t, s.Match(1, "A", "Aaa.func3", "Aaa.func2"))
	test.IsTrue(t, s.Match(1, "A", "Aaa.func3", "Aaa.func*"))
}
