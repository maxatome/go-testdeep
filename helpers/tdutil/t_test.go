// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil_test

import (
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestT(t *testing.T) {
	mockT := &tdutil.T{}

	mockT.Log("Hey this is a log message!")

	buf := mockT.LogBuf()
	if !strings.HasSuffix(buf, "Hey this is a log message!\n") {
		t.Errorf(`LogBuf does not work as expected: "%s"`, buf)
	}
}

func TestRun(t *testing.T) {
	for i, curTest := range []struct {
		fn       func(*testing.T)
		expected bool
	}{
		{
			fn:       func(*testing.T) {},
			expected: true,
		},
		{
			fn: func(t *testing.T) {
				t.Error("An error occurred!")
			},
			expected: false,
		},
	} {
		mockT := &tdutil.T{}

		var called bool
		res := mockT.Run("testname", func(t *testing.T) {
			called = true
			curTest.fn(t)
		})
		if !called {
			t.Errorf("Run#%d func not called", i)
		}
		if res != curTest.expected {
			t.Errorf("Run#%d returned %v ≠ expected %v", i, res, curTest.expected)
		}
	}
}
