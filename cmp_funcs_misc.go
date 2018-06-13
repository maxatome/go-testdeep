// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

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
	t.Helper()

	if got != nil {
		return true
	}

	formatError(t,
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
	t.Helper()

	if got == nil {
		return true
	}

	formatError(t,
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
