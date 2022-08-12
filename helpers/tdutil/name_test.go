// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestBuildTestName(t *testing.T) {
	for i, curTest := range []struct {
		params   []any
		expected string
	}{
		{
			params:   []any{},
			expected: "",
		},
		{
			params:   []any{"foobar"},
			expected: "foobar",
		},
		{
			params:   []any{"foobar %d"},
			expected: "foobar %d",
		},
		{
			params:   []any{"foobar %", 12},
			expected: "foobar %12",
		},
		{
			params:   []any{"foo", "bar"},
			expected: "foobar",
		},
		{
			params:   []any{123, "zip"},
			expected: "123zip",
		},
		{
			params:   []any{123, 456},
			expected: "123 456",
		},
		{
			params:   []any{"foo(%d) bar(%s)", 123, "zip"},
			expected: "foo(123) bar(zip)",
		},
	} {
		name := tdutil.BuildTestName(curTest.params...)
		if name != curTest.expected {
			t.Errorf(`BuildTestName#%d == "%s" but ≠ "%s"`, i, name, curTest.expected)
		}
	}

	tdutil.FbuildTestName(nil)
}
