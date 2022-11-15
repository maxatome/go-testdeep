// Copyright (c) 2018-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/trace"
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

	test.CheckPanic(tt, func() { td.NewT(nil) }, "usage: NewT")
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

func TestParallel(t *testing.T) {
	t.Run("without Parallel", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt)
		t.Parallel()

		// has no effect
	})

	t.Run("with Parallel", func(tt *testing.T) {
		ttt := test.NewParallelTestingTB(tt.Name())

		t := td.NewT(ttt)
		t.Parallel()

		test.IsTrue(tt, ttt.IsParallel)
	})

	t.Run("Run with Parallel", func(tt *testing.T) {
		// This test verifies that subtests with t.Parallel() are run
		// in parallel. We use a WaitGroup to make both subtests block
		// until they're both ready. This test will block forever if
		// the tests are not run together.
		var ready sync.WaitGroup
		ready.Add(2)

		t := td.NewT(tt)

		t.Run("level 1", func(t *td.T) {
			t.Parallel()

			ready.Done() // I'm ready.
			ready.Wait() // Are you?
		})

		t.Run("level 2", func(t *td.T) {
			t.Parallel()

			ready.Done() // I'm ready.
			ready.Wait() // Are you?
		})
	})
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

		ok := t.RunT("Test level1", //nolint: staticcheck
			func(t *td.T) {
				ok := t.RunT("Test level2", //nolint: staticcheck
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

		ok := t.RunT("Test level1", //nolint: staticcheck
			func(t *td.T) {
				ok := t.RunT("Test level2", //nolint: staticcheck
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
	test.IsFalse(tt, ttt.IsFatal, "by default it is not fatal")

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

	// Using T.Assert()
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt, td.ContextConfig{FailureIsFatal: true}).Assert()
	t.True(false) // failure
	test.IsTrue(tt, ttt.LastMessage() != "")
	test.IsFalse(tt, ttt.IsFatal, "by default it is not fatal")

	// Using T.Require()
	ttt = test.NewTestingTB(tt.Name())
	t = td.NewT(ttt).Require()
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
	t = td.NewT(ttt).UseEqual() // enable globally
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = td.NewT(ttt).UseEqual(true) // enable globally
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = td.NewT(ttt).UseEqual(false) // disable globally
	test.IsFalse(tt, t.Cmp(time1, time2))

	t = td.NewT(ttt).UseEqual(time.Time{}) // enable only for time.Time
	test.IsTrue(tt, t.Cmp(time1, time2))

	t = t.UseEqual().UseEqual(false)     // enable then disable globally
	test.IsTrue(tt, t.Cmp(time1, time2)) // Equal() still used

	test.EqualStr(tt,
		ttt.CatchFatal(func() { td.NewT(ttt).UseEqual(42) }),
		"UseEqual expects type int owns an Equal method (@0)")
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

func TestIgnoreUnexported(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	type SType1 struct {
		Public  int
		private string
	}
	a1, b1 := SType1{Public: 42, private: "test"}, SType1{Public: 42}

	type SType2 struct {
		Public  int
		private string
	}
	a2, b2 := SType2{Public: 42, private: "test"}, SType2{Public: 42}

	// Using default config
	t := td.NewT(ttt)
	test.IsFalse(tt, t.Cmp(a1, b1))

	// IgnoreUnexported
	t = td.NewT(ttt).IgnoreUnexported() // ignore unexported globally
	test.IsTrue(tt, t.Cmp(a1, b1))
	test.IsTrue(tt, t.Cmp(a2, b2))

	t = td.NewT(ttt).IgnoreUnexported(true) // ignore unexported globally
	test.IsTrue(tt, t.Cmp(a1, b1))
	test.IsTrue(tt, t.Cmp(a2, b2))

	t = td.NewT(ttt).IgnoreUnexported(false) // handle unexported globally
	test.IsFalse(tt, t.Cmp(a1, b1))
	test.IsFalse(tt, t.Cmp(a2, b2))

	t = td.NewT(ttt).IgnoreUnexported(SType1{}) // ignore only for SType1
	test.IsTrue(tt, t.Cmp(a1, b1))
	test.IsFalse(tt, t.Cmp(a2, b2))

	t = t.UseEqual().UseEqual(false) // enable then disable globally
	test.IsTrue(tt, t.Cmp(a1, b1))
	test.IsFalse(tt, t.Cmp(a2, b2))

	t = td.NewT(ttt).IgnoreUnexported(SType1{}, SType2{}) // enable for both
	test.IsTrue(tt, t.Cmp(a1, b1))
	test.IsTrue(tt, t.Cmp(a2, b2))

	test.EqualStr(tt,
		ttt.CatchFatal(func() { td.NewT(ttt).IgnoreUnexported(42) }),
		"IgnoreUnexported expects type int be a struct, not a int (@0)")
}

func TestTestDeepInGotOK(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	var t *td.T
	cmp := func() bool { return t.Cmp(td.Ignore(), td.Ignore()) }

	// Using default config
	t = td.NewT(ttt)
	test.CheckPanic(tt, func() { cmp() },
		"Found a TestDeep operator in got param, can only use it in expected one!")

	t = td.NewT(ttt).TestDeepInGotOK()
	test.IsTrue(tt, cmp())

	t = t.TestDeepInGotOK(false)
	test.CheckPanic(tt, func() { cmp() },
		"Found a TestDeep operator in got param, can only use it in expected one!")

	t = t.TestDeepInGotOK(true)
	test.IsTrue(tt, cmp())
}

func TestLogTrace(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	t := td.NewT(ttt)

//line /t_struct_test.go:100
	t.LogTrace()
	test.EqualStr(tt, ttt.LastMessage(), `Stack trace:
	TestLogTrace() /t_struct_test.go:100`)
	test.IsFalse(tt, ttt.HasFailed)
	test.IsFalse(tt, ttt.IsFatal)
	ttt.ResetMessages()

//line /t_struct_test.go:110
	t.LogTrace("This is the %s:", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack:
	TestLogTrace() /t_struct_test.go:110`)
	ttt.ResetMessages()

//line /t_struct_test.go:120
	t.LogTrace("This is the %s:\n", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack:
	TestLogTrace() /t_struct_test.go:120`)
	ttt.ResetMessages()

//line /t_struct_test.go:130
	t.LogTrace("This is the ", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack
	TestLogTrace() /t_struct_test.go:130`)
	ttt.ResetMessages()

	trace.IgnorePackage()
	defer trace.UnignorePackage()
//line /t_struct_test.go:140
	t.LogTrace("Stack:\n")
	test.EqualStr(tt, ttt.LastMessage(), `Stack:
	Empty stack trace`)
}

func TestErrorTrace(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	t := td.NewT(ttt)

//line /t_struct_test.go:200
	t.ErrorTrace()
	test.EqualStr(tt, ttt.LastMessage(), `Stack trace:
	TestErrorTrace() /t_struct_test.go:200`)
	test.IsTrue(tt, ttt.HasFailed)
	test.IsFalse(tt, ttt.IsFatal)
	ttt.ResetMessages()

//line /t_struct_test.go:210
	t.ErrorTrace("This is the %s:", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack:
	TestErrorTrace() /t_struct_test.go:210`)
	ttt.ResetMessages()

//line /t_struct_test.go:220
	t.ErrorTrace("This is the %s:\n", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack:
	TestErrorTrace() /t_struct_test.go:220`)
	ttt.ResetMessages()

//line /t_struct_test.go:230
	t.ErrorTrace("This is the ", "stack")
	test.EqualStr(tt, ttt.LastMessage(), `This is the stack
	TestErrorTrace() /t_struct_test.go:230`)
	ttt.ResetMessages()

	trace.IgnorePackage()
	defer trace.UnignorePackage()
//line /t_struct_test.go:240
	t.ErrorTrace("Stack:\n")
	test.EqualStr(tt, ttt.LastMessage(), `Stack:
	Empty stack trace`)
}

func TestFatalTrace(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	t := td.NewT(ttt)

	match := func(got, expectedRe string) {
		tt.Helper()
		re := regexp.MustCompile(expectedRe)
		if !re.MatchString(got) {
			test.EqualErrorMessage(tt, got, expectedRe)
		}
	}

//line /t_struct_test.go:300
	match(ttt.CatchFatal(func() { t.FatalTrace() }), `Stack trace:
	TestFatalTrace\.func\d\(\)   /t_struct_test\.go:300
	\(\*TestingT\)\.CatchFatal\(\) internal/test/types\.go:\d+
	TestFatalTrace\(\)         /t_struct_test\.go:300`)
	test.IsTrue(tt, ttt.HasFailed)
	test.IsTrue(tt, ttt.IsFatal)
	ttt.ResetMessages()

//line /t_struct_test.go:310
	match(ttt.CatchFatal(func() { t.FatalTrace("This is the %s:", "stack") }),
		`This is the stack:
	TestFatalTrace\.func\d\(\)   /t_struct_test\.go:310
	\(\*TestingT\)\.CatchFatal\(\) internal/test/types\.go:\d+
	TestFatalTrace\(\)         /t_struct_test\.go:310`)
	ttt.ResetMessages()

//line /t_struct_test.go:320
	match(ttt.CatchFatal(func() { t.FatalTrace("This is the %s:\n", "stack") }),
		`This is the stack:
	TestFatalTrace\.func\d\(\)   /t_struct_test\.go:320
	\(\*TestingT\)\.CatchFatal\(\) internal/test/types\.go:\d+
	TestFatalTrace\(\)         /t_struct_test\.go:320`)
	ttt.ResetMessages()

//line /t_struct_test.go:330
	match(ttt.CatchFatal(func() { t.FatalTrace("This is the ", "stack") }),
		`This is the stack
	TestFatalTrace\.func\d\(\)   /t_struct_test\.go:330
	\(\*TestingT\)\.CatchFatal\(\) internal/test/types\.go:\d+
	TestFatalTrace\(\)         /t_struct_test\.go:330`)
	ttt.ResetMessages()

	trace.IgnorePackage()
	defer trace.UnignorePackage()
//line /t_struct_test.go:340
	test.EqualStr(tt, ttt.CatchFatal(func() { t.FatalTrace("Stack:\n") }),
		`Stack:
	Empty stack trace`)
}
