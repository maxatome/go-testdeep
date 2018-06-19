// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import "runtime"

// CmpTrue is a shortcut for:
//
//   CmpDeeply(t, got, true, args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpTrue(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, true, args...)
}

// CmpFalse is a shortcut for:
//
//   CmpDeeply(t, got, false, args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpFalse(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, false, args...)
}

func cmpError(ctx Context, t TestingT, got error, args ...interface{}) bool {
	if got != nil {
		return true
	}

	t.Helper()
	formatError(t,
		ctx.failureIsFatal,
		&Error{
			Context:  ctx,
			Message:  "should be an error",
			Got:      rawString("nil"),
			Expected: rawString("non-nil error"),
		},
		args...)

	return false
}

func cmpNoError(ctx Context, t TestingT, got error, args ...interface{}) bool {
	if got == nil {
		return true
	}

	t.Helper()
	formatError(t,
		ctx.failureIsFatal,
		&Error{
			Context:  ctx,
			Message:  "should NOT be an error",
			Got:      got,
			Expected: rawString("nil"),
		},
		args...)

	return false
}

// CmpError checks that "got" is non-nil error.
//
//   _, err := MyFunction(1, 2, 3)
//   CmpError(t, err, "MyFunction(1, 2, 3) should return an error")
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpError(t TestingT, got error, args ...interface{}) bool {
	t.Helper()
	return cmpError(NewContext(), t, got, args...)
}

// CmpNoError checks that "got" is nil error.
//
//   value, err := MyFunction(1, 2, 3)
//   if CmpNoError(t, err) {
//     // one can now check value...
//   }
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNoError(t TestingT, got error, args ...interface{}) bool {
	t.Helper()
	return cmpNoError(NewContext(), t, got, args...)
}

func cmpPanic(ctx Context, t TestingT, fn func(), expected interface{}, args ...interface{}) bool {
	t.Helper()

	if ctx.path == contextDefaultRootName {
		ctx.path = contextPanicRootName
	}

	var (
		panicked   bool
		panicParam interface{}
	)

	func() {
		defer func() { panicParam = recover() }()
		panicked = true
		fn()
		panicked = false
	}()

	if !panicked {
		formatError(t,
			ctx.failureIsFatal,
			&Error{
				Context: ctx,
				Message: "should have panicked",
				Summary: rawString("did not panic"),
			},
			args...)
		return false
	}

	return cmpDeeply(ctx.AddDepth("→panic()"), t, panicParam, expected, args...)
}

func cmpNotPanic(ctx Context, t TestingT, fn func(), args ...interface{}) bool {
	var (
		panicked   bool
		stackTrace rawString
	)

	func() {
		defer func() {
			panicParam := recover()
			if panicked {
				buf := make([]byte, 8192)
				n := runtime.Stack(buf, false)
				for ; n > 0; n-- {
					if buf[n-1] != '\n' {
						break
					}
				}
				stackTrace = rawString("panic: " + toString(panicParam) + "\n\n" +
					string(buf[:n]))
			}
		}()
		panicked = true
		fn()
		panicked = false
	}()

	if !panicked {
		return true
	}

	t.Helper()

	if ctx.path == contextDefaultRootName {
		ctx.path = contextPanicRootName
	}

	formatError(t,
		ctx.failureIsFatal,
		&Error{
			Context:  ctx,
			Message:  "should NOT have panicked",
			Got:      stackTrace,
			Expected: rawString("not panicking at all"),
		})
	return false
}

// CmpPanic calls "fn" and checks a panic() occurred with the
// "expectedPanic" parameter. It returns true only if both conditions
// are fulfilled.
//
// Note that calling panic(nil) in "fn" body is detected as a panic
// (in this case "expectedPanic" has to be nil.)
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpPanic(t TestingT, fn func(), expectedPanic interface{},
	args ...interface{}) bool {
	t.Helper()
	return cmpPanic(NewContext(), t, fn, expectedPanic, args...)
}

// CmpNotPanic calls "fn" and checks no panic() occurred. If a panic()
// occurred false is returned then the panic() parameter and the stack
// trace appear in the test report.
//
// Note that calling panic(nil) in "fn" body is detected as a panic.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNotPanic(t TestingT, fn func(), args ...interface{}) bool {
	t.Helper()
	return cmpNotPanic(NewContext(), t, fn, args...)
}
