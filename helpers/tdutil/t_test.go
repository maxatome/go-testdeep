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
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestT(t *testing.T) {
	mockT := tdutil.NewT("hey!")

	if n := mockT.Name(); n != "hey!" {
		t.Errorf(`Test name is not correct, got: %s, expected: hey!`, n)
	}

	mockT.Log("Hey this is a log message!")

	buf := mockT.LogBuf()
	if !strings.HasSuffix(buf, "Hey this is a log message!\n") {
		t.Errorf(`LogBuf does not work as expected: "%s"`, buf)
	}
}

func TestFailNow(t *testing.T) {
	mockT := tdutil.NewT("hey!")

	test.IsFalse(t, mockT.CatchFailNow(func() {}))

	test.IsTrue(t, mockT.CatchFailNow(func() { mockT.FailNow() }))
	test.IsTrue(t, mockT.CatchFailNow(func() { mockT.Fatal("Ouch!") }))
	test.IsTrue(t, mockT.CatchFailNow(func() { mockT.Fatalf("Ouch!") }))

	// No FailNow() but panic()
	var (
		panicked, failNowOccurred bool
		panicParam                any
	)
	func() {
		defer func() { panicParam = recover() }()

		panicked = true
		failNowOccurred = mockT.CatchFailNow(func() { panic("Boom!") })
		panicked = false
	}()

	test.IsFalse(t, failNowOccurred)
	if test.IsTrue(t, panicked) {
		panicStr, ok := panicParam.(string)
		if test.IsTrue(t, ok) {
			test.EqualStr(t, panicStr, "Boom!")
		}
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
