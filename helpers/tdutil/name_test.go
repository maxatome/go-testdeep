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
		params   []interface{}
		expected string
	}{
		{
			params:   []interface{}{},
			expected: "",
		},
		{
			params:   []interface{}{"foobar"},
			expected: "foobar",
		},
		{
			params:   []interface{}{"foo", "bar"},
			expected: "foobar",
		},
		{
			params:   []interface{}{123, "zip"},
			expected: "123zip",
		},
		{
			params:   []interface{}{123, 456},
			expected: "123 456",
		},
		{
			params:   []interface{}{"foo(%d) bar(%s)", 123, "zip"},
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
