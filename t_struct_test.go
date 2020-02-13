// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestT(tt *testing.T) {
	// We don't want to include "anchors" field in comparison
	cmp := func(tt *testing.T, got, expected testdeep.ContextConfig) {
		tt.Helper()
		testdeep.Cmp(tt, got,
			testdeep.SStruct(expected, testdeep.StructFields{
				"anchors": testdeep.Ignore(),
			}),
		)
	}

	tt.Run("without config", func(tt *testing.T) {
		t := testdeep.NewT(tt)
		cmp(tt, t.Config, testdeep.DefaultContextConfig)

		tDup := testdeep.NewT(t)
		cmp(tt, tDup.Config, testdeep.DefaultContextConfig)
	})

	tt.Run("explicit default config", func(tt *testing.T) {
		t := testdeep.NewT(tt, testdeep.ContextConfig{})
		cmp(tt, t.Config, testdeep.DefaultContextConfig)

		tDup := testdeep.NewT(t)
		cmp(tt, tDup.Config, testdeep.DefaultContextConfig)
	})

	tt.Run("specific config", func(tt *testing.T) {
		conf := testdeep.ContextConfig{
			RootName:  "TEST",
			MaxErrors: 33,
		}
		t := testdeep.NewT(tt, conf)
		cmp(tt, t.Config, conf)

		tDup := testdeep.NewT(t)
		cmp(tt, tDup.Config, conf)

		newConf := conf
		newConf.MaxErrors = 34
		tDup = testdeep.NewT(t, newConf)
		cmp(tt, tDup.Config, newConf)

		t2 := t.RootName("T2")
		cmp(tt, t.Config, conf)
		cmp(tt, t2.Config, testdeep.ContextConfig{
			RootName:  "T2",
			MaxErrors: 33,
		})

		t3 := t.RootName("")
		cmp(tt, t3.Config, testdeep.ContextConfig{
			RootName:  "DATA",
			MaxErrors: 33,
		})
	})

	//
	// Bad usages
	test.CheckPanic(tt,
		func() {
			testdeep.NewT(tt, testdeep.ContextConfig{}, testdeep.ContextConfig{})
		},
		"usage: NewT")

	test.CheckPanic(tt, func() { testdeep.NewT(nil) }, "usage: NewT")
}

func TestTCmp(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())
	t := testdeep.NewT(ttt)
	testdeep.CmpTrue(tt, t.Cmp(1, 1))
	testdeep.CmpFalse(tt, ttt.Failed())

	ttt = test.NewTestingFT(tt.Name())
	t = testdeep.NewT(ttt)
	testdeep.CmpFalse(tt, t.Cmp(1, 2))
	testdeep.CmpTrue(tt, ttt.Failed())
}

func TestTCmpDeeply(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())
	t := testdeep.NewT(ttt)
	testdeep.CmpTrue(tt, t.CmpDeeply(1, 1))
	testdeep.CmpFalse(tt, ttt.Failed())

	ttt = test.NewTestingFT(tt.Name())
	t = testdeep.NewT(ttt)
	testdeep.CmpFalse(tt, t.CmpDeeply(1, 2))
	testdeep.CmpTrue(tt, ttt.Failed())
}

func TestRunT(tt *testing.T) {
	t := testdeep.NewT(tt)

	runPassed := false

	ok := t.RunT("Test level1",
		func(t *testdeep.T) {
			ok := t.RunT("Test level2",
				func(t *testdeep.T) {
					runPassed = t.True(true) // test succeeds!
				})

			t.True(ok)
		})

	t.True(ok)
	t.True(runPassed)
}

func TestFailureIsFatal(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())

	// All t.True(false) tests of course fail

	// Using default config
	t := testdeep.NewT(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "by default it not fatal")

	// Using specific config
	t = testdeep.NewT(ttt, testdeep.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal()
	t = testdeep.NewT(ttt).FailureIsFatal()
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal(true)
	t = testdeep.NewT(ttt).FailureIsFatal(true)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require()
	t = testdeep.Require(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require() with specific config (cannot override FailureIsFatal)
	t = testdeep.Require(ttt, testdeep.ContextConfig{FailureIsFatal: false})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Canceling specific config
	t = testdeep.NewT(ttt, testdeep.ContextConfig{FailureIsFatal: false}).
		FailureIsFatal(false)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert()
	t = testdeep.Assert(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert() with specific config (cannot override FailureIsFatal)
	t = testdeep.Assert(ttt, testdeep.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / assert
	t, _ = testdeep.AssertRequire(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using AssertRequire() / assert with specific config (cannot
	// override FailureIsFatal)
	t, _ = testdeep.AssertRequire(ttt, testdeep.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / require
	_, t = testdeep.AssertRequire(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using AssertRequire() / require with specific config (cannot
	// override FailureIsFatal)
	_, t = testdeep.AssertRequire(ttt, testdeep.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")
}

func TestUseEqual(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())

	var time1, time2 time.Time
	for {
		time1 = time.Now()
		time2 = time1.Truncate(0)
		if !time1.Equal(time2) {
			tt.Fatal("time.Equal() does not work as expected")
		}
		if time1 != time2 { // to avoid the bad luck case where time1.wall=0
			break
		}
	}

	// Using default config
	t := testdeep.NewT(ttt)
	test.IsFalse(tt, t.Cmp(time1, time2))

	// UseEqual
	t = testdeep.NewT(ttt).UseEqual()
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = testdeep.NewT(ttt).UseEqual(true)
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = testdeep.NewT(ttt).UseEqual(false)
	test.IsFalse(tt, t.Cmp(time1, time2))
}

func TestBeLax(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())

	// Using default config
	t := testdeep.NewT(ttt)
	test.IsFalse(tt, t.Cmp(int64(123), 123))

	// BeLax
	t = testdeep.NewT(ttt).BeLax()
	test.IsTrue(tt, t.Cmp(int64(123), 123))

	t = testdeep.NewT(ttt).BeLax(true)
	test.IsTrue(tt, t.Cmp(int64(123), 123))

	t = testdeep.NewT(ttt).BeLax(false)
	test.IsFalse(tt, t.Cmp(int64(123), 123))
}
