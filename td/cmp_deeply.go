// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/trace"
)

func init() {
	trace.Init()
	trace.IgnorePackage()
}

// stripTrace removes go-testdeep useless calls in a trace returned by
// trace.Retrieve().
//
// Remove (*T).Run() call:
//    A()                   zz.go:20
//    TestSubtestTd.func1() zz_test.go:24
//   →(*T).Run.func1()      github.com/maxatome/go-testdeep@xxx/td/t_struct.go:554
//
// Remove (*T).RunAssertRequire() call:
//    A()                           zz.go:20
//    TestSubtestTd.func1()         zz_test.go:24
//   →(*T).RunAssertRequire.func1() github.com/maxatome/go-testdeep@xxx/td/t_struct.go:554
//
// Remove (*TestAPI).Run() call:
//    Subtest.func1() zz.go:20
//   →(*TestAPI).Run.func1() github.com/maxatome/go-testdeep@xxx/helpers/tdhttp/test_api.go:119
//   ⇒(*T).Run.func1()       github.com/maxatome/go-testdeep@xxx/td/t_struct.go:554
//
// Remove github.com/maxatome/go-testdeep/helpers/tdsuite calls:
//    A()               zz.go:20
//    Suite.TestSuite() zz_test.go:35
//   →Value.call()      $GOROOT/src/reflect/value.go:476
//   →Value.Call()      $GOROOT/src/reflect/value.go:337
//   →run.func2()       github.com/maxatome/go-testdeep@xxx/helpers/tdsuite/suite.go:301
//   ⇒(*T).Run.func1()  github.com/maxatome/go-testdeep@xxx/td/t_struct.go:554
func stripTrace(tce []trace.Level) []trace.Level {
	if len(tce) <= 1 {
		return tce
	}

	// Remove useless possible (*T).Run() or (*T).RunAssertRequire() first call
	first := tce[len(tce)-1]
	if first.Package != "github.com/maxatome/go-testdeep/td" ||
		(first.Func != "(*T).Run.func1" &&
			first.Func != "(*T).RunAssertRequire.func1") {
		return tce
	}

	tce = tce[:len(tce)-1]
	first = tce[len(tce)-1]

	// Remove useless tdhttp (*TestAPI).Run() call
	if first.Package == "github.com/maxatome/go-testdeep/helpers/tdhttp" &&
		first.Func == "(*TestAPI).Run.func1" {
		return tce[:len(tce)-1]
	}

	// Remove useless tdsuite calls
	if first.Package != "github.com/maxatome/go-testdeep/helpers/tdsuite" ||
		!strings.HasPrefix(first.Func, "run.func") {
		return tce
	}

	for i := len(tce) - 2; i >= 1; i-- {
		if tce[i].Package != "reflect" {
			return tce[:i+1]
		}
	}
	return nil
}

func formatError(t TestingT, isFatal bool, err *ctxerr.Error, args ...interface{}) {
	t.Helper()

	const failedTest = "Failed test"

	args = flat.Interfaces(args...)

	var buf bytes.Buffer
	color.AppendTestNameOn(&buf)
	if len(args) == 0 {
		buf.WriteString(failedTest + "\n")
	} else {
		buf.WriteString(failedTest + " '")
		tdutil.FbuildTestName(&buf, args...)
		buf.WriteString("'\n")
	}
	color.AppendTestNameOff(&buf)

	err.Append(&buf, "")

	// Stask trace
	if tce := stripTrace(trace.Retrieve(0, "testing.tRunner")); len(tce) > 1 {
		buf.WriteString("\nThis is how we got here:\n")

		fnMaxLen := 0
		for _, level := range tce {
			if len(level.Func) > fnMaxLen {
				fnMaxLen = len(level.Func)
			}
		}
		fnMaxLen += 2

		nl := ""
		for _, level := range tce {
			fmt.Fprintf(&buf, "%s\t%-*s %s", nl, fnMaxLen, level.Func+"()", level.FileLine)
			nl = "\n"
		}
	}

	if isFatal {
		t.Fatal(buf.String())
	} else {
		t.Error(buf.String())
	}
}

func cmpDeeply(ctx ctxerr.Context, t TestingT, got, expected interface{},
	args ...interface{}) bool {
	err := deepValueEqualFinal(ctx,
		reflect.ValueOf(got), reflect.ValueOf(expected))
	if err == nil {
		return true
	}

	t.Helper()
	formatError(t, ctx.FailureIsFatal, err, args...)
	return false
}

// Cmp returns true if "got" matches "expected". "expected" can
// be the same type as "got" is, or contains some TestDeep
// operators. If "got" does not match "expected", it returns false and
// the reason of failure is logged with the help of "t" Error()
// method.
//
//   got := "foobar"
//   td.Cmp(t, got, "foobar")            // succeeds
//   td.Cmp(t, got, td.HasPrefix("foo")) // succeeds
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func Cmp(t TestingT, got, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return cmpDeeply(newContext(), t, got, expected, args...)
}

// CmpDeeply works the same as Cmp and is still available for
// compatibility purpose. Use shorter Cmp in new code.
//
//   got := "foobar"
//   td.CmpDeeply(t, got, "foobar")            // succeeds
//   td.CmpDeeply(t, got, td.HasPrefix("foo")) // succeeds
func CmpDeeply(t TestingT, got, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return cmpDeeply(newContext(), t, got, expected, args...)
}
