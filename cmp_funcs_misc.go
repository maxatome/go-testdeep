// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"testing"
)

// CmpTrue is a shortcut for:
//
//   CmpDeeply(t, got, true, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpTrue(t *testing.T, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, true, args...)
}

// CmpFalse is a shortcut for:
//
//   CmpDeeply(t, got, false, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpFalse(t *testing.T, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, false, args...)
}

func cmpError(ctx Context, t *testing.T, got error, args ...interface{}) bool {
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

func cmpNoError(ctx Context, t *testing.T, got error, args ...interface{}) bool {
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
func CmpError(t *testing.T, got error, args ...interface{}) bool {
	t.Helper()
	return cmpError(NewContext(), t, got, args...)
}

// CmpNoError checks that "got" is nil error.
//
//   value, err := MyFunction(1, 2, 3)
//   if CmpNoError(t, err) {
//     // one can now check value...
//   }
func CmpNoError(t *testing.T, got error, args ...interface{}) bool {
	t.Helper()
	return cmpNoError(NewContext(), t, got, args...)
}
