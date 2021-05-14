// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/dark"
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
				"hooks":   td.Ignore(),
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
	ttb := test.NewTestingTB("usage params")
	ttb.CatchFatal(func() {
		td.NewT(ttb, td.ContextConfig{}, td.ContextConfig{})
	})
	test.IsTrue(tt, ttb.IsFatal)
	test.IsTrue(tt, strings.Contains(ttb.Messages[0], "usage: NewT("))

	dark.CheckFatalizerBarrierErr(tt, func() { td.NewT(nil) }, "usage: NewT")
}

func TestTCmp(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())
	t := td.NewT(ttt)
	test.IsTrue(tt, t.Cmp(1, 1))
	test.IsFalse(tt, ttt.Failed())

	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt)
	test.IsFalse(tt, t.Cmp(1, 2))
	test.IsTrue(tt, ttt.Failed())
}

func TestTCmpDeeply(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())
	t := td.NewT(ttt)
	test.IsTrue(tt, t.CmpDeeply(1, 1))
	test.IsFalse(tt, ttt.Failed())

	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt)
	test.IsFalse(tt, t.CmpDeeply(1, 2))
	test.IsTrue(tt, ttt.Failed())
}

func TestRun(t *testing.T) {
	t.Run("test.TB with Run", func(tt *testing.T) {
		t := td.NewT(tt)

		runPassed := false
		nestedFailureIsFatal := false

		ok := t.Run("Test level1",
			func(t *td.T) {
				ok := t.FailureIsFatal().Run("Test level2",
					func(t *td.T) {
						runPassed = t.True(true) // test succeeds!

						// Check we inherit config from caller
						nestedFailureIsFatal = t.Config.FailureIsFatal
					})

				t.True(ok)
			})

		test.IsTrue(tt, ok)
		test.IsTrue(tt, runPassed)
		test.IsTrue(tt, nestedFailureIsFatal)
	})

	t.Run("test.TB without Run", func(tt *testing.T) {
		t := td.NewT(test.NewTestingTB("gg"))

		runPassed := false

		ok := t.Run("Test level1",
			func(t *td.T) {
				ok := t.Run("Test level2",
					func(t *td.T) {
						runPassed = t.True(true) // test succeeds!
					})

				t.True(ok)
			})

		t.True(ok)
		t.True(runPassed)
	})
}

func TestRunAssertRequire(t *testing.T) {
	t.Run("test.TB with Run", func(tt *testing.T) {
		t := td.NewT(tt)

		runPassed := false
		assertIsFatal := true
		requireIsFatal := false

		ok := t.RunAssertRequire("Test level1",
			func(assert, require *td.T) {
				assertIsFatal = assert.Config.FailureIsFatal
				requireIsFatal = require.Config.FailureIsFatal

				ok := assert.RunAssertRequire("Test level2",
					func(assert, require *td.T) {
						runPassed = assert.True(true)               // test succeeds!
						runPassed = runPassed && require.True(true) // test succeeds!

						assertIsFatal = assertIsFatal || assert.Config.FailureIsFatal
						requireIsFatal = requireIsFatal && require.Config.FailureIsFatal
					})
				assert.True(ok)
				require.True(ok)

				ok = require.RunAssertRequire("Test level2",
					func(assert, require *td.T) {
						runPassed = runPassed && assert.True(true)  // test succeeds!
						runPassed = runPassed && require.True(true) // test succeeds!

						assertIsFatal = assertIsFatal || assert.Config.FailureIsFatal
						requireIsFatal = requireIsFatal && require.Config.FailureIsFatal
					})
				assert.True(ok)
				require.True(ok)
			})

		test.IsTrue(tt, ok)
		test.IsTrue(tt, runPassed)
		test.IsFalse(tt, assertIsFatal)
		test.IsTrue(tt, requireIsFatal)
	})

	t.Run("test.TB without Run", func(tt *testing.T) {
		t := td.NewT(test.NewTestingTB("gg"))

		runPassed := false
		assertIsFatal := true
		requireIsFatal := false

		ok := t.RunAssertRequire("Test level1",
			func(assert, require *td.T) {
				assertIsFatal = assert.Config.FailureIsFatal
				requireIsFatal = require.Config.FailureIsFatal

				ok := assert.RunAssertRequire("Test level2",
					func(assert, require *td.T) {
						runPassed = assert.True(true)               // test succeeds!
						runPassed = runPassed && require.True(true) // test succeeds!

						assertIsFatal = assertIsFatal || assert.Config.FailureIsFatal
						requireIsFatal = requireIsFatal && require.Config.FailureIsFatal
					})
				assert.True(ok)
				require.True(ok)

				ok = require.RunAssertRequire("Test level2",
					func(assert, require *td.T) {
						runPassed = runPassed && assert.True(true)  // test succeeds!
						runPassed = runPassed && require.True(true) // test succeeds!

						assertIsFatal = assertIsFatal || assert.Config.FailureIsFatal
						requireIsFatal = requireIsFatal && require.Config.FailureIsFatal
					})
				assert.True(ok)
				require.True(ok)
			})

		test.IsTrue(tt, ok)
		test.IsTrue(tt, runPassed)
		test.IsFalse(tt, assertIsFatal)
		test.IsTrue(tt, requireIsFatal)
	})
}

// Deprecated RunT.
func TestRunT(t *testing.T) {
	t.Run("test.TB with Run", func(tt *testing.T) {
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

		test.IsTrue(tt, ok)
		test.IsTrue(tt, runPassed)
	})

	t.Run("test.TB without Run", func(tt *testing.T) {
		t := td.NewT(test.NewTestingTB("gg"))

		runPassed := false

		ok := t.RunT("Test level1",
			func(t *td.T) {
				ok := t.RunT("Test level2",
					func(t *td.T) {
						runPassed = t.True(true) // test succeeds!
					})

				t.True(ok)
			})

		test.IsTrue(tt, ok)
		test.IsTrue(tt, runPassed)
	})
}

func TestFailureIsFatal(tt *testing.T) {
	// All t.True(false) tests of course fail

	// Using default config
	ttt := test.NewTestingTB(tt.Name())
	t := td.NewT(ttt)
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "by default it not fatal")

	// Using specific config
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt, td.ContextConfig{FailureIsFatal: true})
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal()
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt).FailureIsFatal()
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal(true)
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt).FailureIsFatal(true)
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require()
	ttt = test.NewTestingTB(tt.Name())
	t = td.Require(ttt)
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using Require() with specific config (cannot override FailureIsFatal)
	ttt = test.NewTestingTB(tt.Name())
	t = td.Require(ttt, td.ContextConfig{FailureIsFatal: false})
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Canceling specific config
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt, td.ContextConfig{FailureIsFatal: false}).
		FailureIsFatal(false)
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert()
	ttt = test.NewTestingTB(tt.Name())
	t = td.Assert(ttt)
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using Assert() with specific config (cannot override FailureIsFatal)
	ttt = test.NewTestingTB(tt.Name())
	t = td.Assert(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / assert
	ttt = test.NewTestingTB(tt.Name())
	t, _ = td.AssertRequire(ttt)
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "it must be not fatal")

	// Using AssertRequire() / assert with specific config (cannot
	// override FailureIsFatal)
	ttt = test.NewTestingTB(tt.Name())
	t, _ = td.AssertRequire(ttt, td.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "it must be not fatal")

	// AssertRequire() / require
	ttt = test.NewTestingTB(tt.Name())
	_, t = td.AssertRequire(ttt)
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using AssertRequire() / require with specific config (cannot
	// override FailureIsFatal)
	ttt = test.NewTestingTB(tt.Name())
	_, t = td.AssertRequire(ttt, td.ContextConfig{FailureIsFatal: true})
	ttt.CatchFatal(func() { t.True(false) }) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsTrue(tt, ttt.IsFatal, "it must be fatal")
}

func TestUseEqual(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

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
	ttt := test.NewTestingTB(tt.Name())

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
