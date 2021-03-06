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

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/trace"
)

func init() {
	trace.Init()
	trace.IgnorePackage(1)
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
	if trace := trace.Retrieve(2, "testing.tRunner"); len(trace) > 1 {
		buf.WriteString("\nThis is how we got here:\n")

		fnMaxLen := 0
		for _, level := range trace {
			if len(level.Func) > fnMaxLen {
				fnMaxLen = len(level.Func)
			}
		}
		fnMaxLen += 2

		nl := ""
		for _, level := range trace {
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
