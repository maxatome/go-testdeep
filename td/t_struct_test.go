// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestT(tt *testing.T) {
	// We don't want to include "anchors" field in comparison
	cmp := func(tt *testing.T, got, expected td.ContextConfig) {
		tt.Helper()
		td.Cmp(tt, got,
			td.SStruct(expected, td.StructFields{
				"anchors": td.Ignore(),
			}),
		)
	}

	tt.Run("without config", func(tt *testing.T) {
		t := td.NewT(tt)
		cmp(tt, t.Config, td.DefaultContextConfig)

		tDup := td.NewT(t)
		cmp(tt, tDup.Config, td.DefaultContextConfig)
	})

	tt.Run("explicit default config", func(tt *testing.T) {
		t := td.NewT(tt, td.ContextConfig{})
		cmp(tt, t.Config, td.DefaultContextConfig)

		tDup := td.NewT(t)
		cmp(tt, tDup.Config, td.DefaultContextConfig)
	})

	tt.Run("specific config", func(tt *testing.T) {
		conf := td.ContextConfig{
			RootName:  "TEST",
			MaxErrors: 33,
		}
		t := td.NewT(tt, conf)
		cmp(tt, t.Config, conf)

		tDup := td.NewT(t)
		cmp(tt, tDup.Config, conf)

		newConf := conf
		newConf.MaxErrors = 34
		tDup = td.NewT(t, newConf)
		cmp(tt, tDup.Config, newConf)

		t2 := t.RootName("T2")
		cmp(tt, t.Config, conf)
		cmp(tt, t2.Config, td.ContextConfig{
			RootName:  "T2",
			MaxErrors: 33,
		})

		t3 := t.RootName("")
		cmp(tt, t3.Config, td.ContextConfig{
			RootName:  "DATA",
			MaxErrors: 33,
		})
	})

	//
	// Bad usages
	test.CheckPanic(tt,
		func() {
			td.NewT(tt, td.ContextConfig{}, td.ContextConfig{})
		},
		"usage: NewT")

	test.CheckPanic(tt, func() { td.NewT(nil) }, "usage: NewT")
}

func TestTCmp(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())
	t := td.NewT(ttt)
	td.CmpTrue(tt, t.Cmp(1, 1))
	td.CmpFalse(tt, ttt.Failed())

	ttt = test.NewTestingFT(tt.Name())
	t = td.NewT(ttt)
	td.CmpFalse(tt, t.Cmp(1, 2))
	td.CmpTrue(tt, ttt.Failed())
}

func TestTCmpDeeply(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())
	t := td.NewT(ttt)
	td.CmpTrue(tt, t.CmpDeeply(1, 1))
	td.CmpFalse(tt, ttt.Failed())

	ttt = test.NewTestingFT(tt.Name())
	t = td.NewT(ttt)
	td.CmpFalse(tt, t.CmpDeeply(1, 2))
	td.CmpTrue(tt, ttt.Failed())
}

func TestRunT(tt *testing.T) {
	t := td.NewT(tt)

	runPassed := false

	ok := t.RunT("Test level1",
		func(t *td.T) {
			ok := t.RunT("Test level2",
				func(t *td.T) {
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
	t := td.NewT(ttt)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "by default it not fatal")

	// Using specific config
	t = td.NewT(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal()
	t = td.NewT(ttt).FailureIsFatal()
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal(true)
	t = td.NewT(ttt).FailureIsFatal(true)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require()
	t = td.Require(ttt)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require() with specific config (cannot override FailureIsFatal)
	t = td.Require(ttt, td.ContextConfig{FailureIsFatal: false})
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Canceling specific config
	t = td.NewT(ttt, td.ContextConfig{FailureIsFatal: false}).
		FailureIsFatal(false)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert()
	t = td.Assert(ttt)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert() with specific config (cannot override FailureIsFatal)
	t = td.Assert(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / assert
	t, _ = td.AssertRequire(ttt)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using AssertRequire() / assert with specific config (cannot
	// override FailureIsFatal)
	t, _ = td.AssertRequire(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / require
	_, t = td.AssertRequire(ttt)
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using AssertRequire() / require with specific config (cannot
	// override FailureIsFatal)
	_, t = td.AssertRequire(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	td.CmpNotEmpty(tt, ttt.LastMessage)
	td.CmpTrue(tt, ttt.IsFatal, "it must be fatal")
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
	t := td.NewT(ttt)
	test.IsFalse(tt, t.Cmp(time1, time2))

	// UseEqual
	t = td.NewT(ttt).UseEqual()
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = td.NewT(ttt).UseEqual(true)
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = td.NewT(ttt).UseEqual(false)
	test.IsFalse(tt, t.Cmp(time1, time2))
}

func TestBeLax(tt *testing.T) {
	ttt := test.NewTestingFT(tt.Name())

	// Using default config
	t := td.NewT(ttt)
	test.IsFalse(tt, t.Cmp(int64(123), 123))

	// BeLax
	t = td.NewT(ttt).BeLax()
	test.IsTrue(tt, t.Cmp(int64(123), 123))

	t = td.NewT(ttt).BeLax(true)
	test.IsTrue(tt, t.Cmp(int64(123), 123))

	t = td.NewT(ttt).BeLax(false)
	test.IsFalse(tt, t.Cmp(int64(123), 123))
}
