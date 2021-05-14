// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/trace"
)

func TestStripTrace(t *testing.T) {
	check := func(got, expected trace.Stack) {
		got = stripTrace(got)
		if !reflect.DeepEqual(got, expected) {
			t.Helper()
			t.Errorf("\n     got: %#v\nexpected: %#v", got, expected)
		}
	}

	check(nil, nil)

	s := trace.Stack{
		{Package: "test", Func: "A"},
	}
	check(s, s)

	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "TestSimple"},
	}
	check(s, s)

	// inside testing.Cleanup() call
	s = trace.Stack{
		{Package: "test", Func: "TestCleanup.func2"},
		{Package: "testing", Func: "(*common).Cleanup.func1"},
		{Package: "testing", Func: "(*common).runCleanup"},
		{Package: "testing", Func: "tRunner.func2"},
	}
	check(s, s[:1])

	//
	// td
	//
	// td.(*T).Run() call
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "TestSubtestTd.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).Run.func1"},
	}
	check(s, s[:2])

	// td.(*T).RunAssertRequire() call
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "TestSubtestTd.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).RunAssertRequire.func1"},
	}
	check(s, s[:2])

	//
	// tdhttp
	//
	// tdhttp.(*TestAPI).Run() call
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "TestSubtestTd.func1"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdhttp", Func: "(*TestAPI).Run.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).Run.func1"},
	}
	check(s, s[:2])

	//
	// tdsuite
	//
	// tdsuite.Run() call → TestSuite(*td.T)
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.TestSuite"},
		{Package: "reflect", Func: "Value.call"},
		{Package: "reflect", Func: "Value.Call"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func2"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).Run.func1"},
	}
	check(s, s[:2])

	// tdsuite.Run() call → TestSuite(assert, require *td.T)
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.TestSuite"},
		{Package: "reflect", Func: "Value.call"},
		{Package: "reflect", Func: "Value.Call"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).RunAssertRequire.func1"},
	}
	check(s, s[:2])

	// tdsuite.Run() call → Suite.Setup()
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.Setup"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteSetup"},
	}
	check(s, append(s[:2:2], s[4]))

	// tdsuite.Run() call → Suite.PreTest()
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.PreTest"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func2"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).Run.func1"},
	}
	check(s, s[:2])

	// tdsuite.Run() call → Suite.PostTest()
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.PostTest"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func2.1"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func2"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).Run.func1"},
	}
	check(s, s[:2])

	// tdsuite.Run() call → Suite.BetweenTests()
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.BetweenTests"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteBetweenTests"},
	}
	check(s, append(s[:2:2], s[4]))

	// tdsuite.Run() call → Suite.Destroy()
	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.Destroy"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func1"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteDestroy"},
	}
	check(s, append(s[:2:2], s[5]))

	// Improbable cases
	s = trace.Stack{
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteDestroy"},
	}
	check(s, s[:1])

	s = trace.Stack{
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "y"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "x"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteDestroy"},
	}
	check(s, s[:1])

	s = trace.Stack{
		{Package: "test", Func: "Suite.TestXxx"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "Run"},
		{Package: "test", Func: "TestSuiteDestroy"},
	}
	check(s, append(s[:1:1], s[2]))

	s = trace.Stack{
		{Package: "reflect", Func: "Value.call"},
		{Package: "reflect", Func: "Value.Call"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).RunAssertRequire.func1"},
	}
	check(s, nil)

	s = trace.Stack{
		{Package: "test", Func: "A"},
		{Package: "test", Func: "Suite.TestSuite"},
		{Package: "github.com/maxatome/go-testdeep/helpers/tdsuite", Func: "run.func1"},
		{Package: "github.com/maxatome/go-testdeep/td", Func: "(*T).RunAssertRequire.func1"},
	}
	check(s, s[:2])
}

func TestFormatError(t *testing.T) {
	err := &ctxerr.Error{
		Context: newContext(),
		Message: "test error message",
		Summary: ctxerr.NewSummary("test error summary"),
	}

	nonStringName := bytes.NewBufferString("zip!")

	for _, fatal := range []bool{false, true} {
		//
		// Without args
		ttt := test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With one arg
		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, "foo bar!") })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'foo bar!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, nonStringName) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'zip!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args & Printf format
		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, "hello %d!", 123) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'hello 123!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args & Printf format + Flatten
		ttt = test.NewTestingT()
		ttt.CatchFatal(func() {
			formatError(ttt, fatal, err, "hello %s → %d/%d!", "bob", Flatten([]int{123, 125}))
		})
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'hello bob → 123/125!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args without Printf format
		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, "hello ", "world! ", 123) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'hello world! 123'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args without Printf format + Flatten
		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, "hello ", "world! ", Flatten([]int{123, 125})) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'hello world! 123 125'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		ttt = test.NewTestingT()
		ttt.CatchFatal(func() { formatError(ttt, fatal, err, nonStringName, "hello ", "world! ", 123) })
		test.EqualStr(t, ttt.LastMessage(), `Failed test 'zip!hello world! 123'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)
	}
}

func TestCmp(t *testing.T) {
	tt := test.NewTestingTB(t.Name())
	test.IsTrue(t, Cmp(tt, 1, 1))
	test.IsFalse(t, tt.Failed())

	tt = test.NewTestingTB(t.Name())
	test.IsFalse(t, Cmp(tt, 1, 2))
	test.IsTrue(t, tt.Failed())
}

func TestCmpDeeply(t *testing.T) {
	tt := test.NewTestingTB(t.Name())
	test.IsTrue(t, CmpDeeply(tt, 1, 1))
	test.IsFalse(t, tt.Failed())

	tt = test.NewTestingTB(t.Name())
	test.IsFalse(t, CmpDeeply(tt, 1, 2))
	test.IsTrue(t, tt.Failed())
}
