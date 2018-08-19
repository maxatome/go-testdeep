// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package test

import (
	"fmt"
	"strings"
	"testing"
)

// EqualErrorMessage prints a test error message of the form:
//
//   Message
//   Failed test
//          got: got_value
//     expected: expected_value
func EqualErrorMessage(t *testing.T, got, expected interface{},
	args ...interface{}) {
	t.Helper()
	t.Errorf(`%sFailed test
	     got: %v
	expected: %v`,
		BuildTestName(args), got, expected)
}

// EqualStr checks that got equals expected.
func EqualStr(t *testing.T, got, expected string, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	EqualErrorMessage(t, got, expected, args...)
	return false
}

// EqualInt checks that got equals expected.
func EqualInt(t *testing.T, got, expected int, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	EqualErrorMessage(t, got, expected, args...)
	return false
}

// EqualBool checks that got equals expected.
func EqualBool(t *testing.T, got, expected bool, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	EqualErrorMessage(t, got, expected, args...)
	return false
}

// BuildTestName builds a string given args.
func BuildTestName(args []interface{}) string {
	switch len(args) {
	case 0:
		return ""

	case 1:
		return args[0].(string) + "\n"

	default:
		return fmt.Sprintf(args[0].(string)+"\n", args[1:]...)
	}
}

// IsTrue checks that got is true.
func IsTrue(t *testing.T, got bool, args ...interface{}) bool {
	if got {
		return true
	}

	t.Helper()
	EqualErrorMessage(t, false, true, args...)
	return false
}

// IsFalse checks that got is false.
func IsFalse(t *testing.T, got bool, args ...interface{}) bool {
	if !got {
		return true
	}

	t.Helper()
	EqualErrorMessage(t, true, false, args...)
	return false
}

// CheckPanic checks that fn() panics and that the panic() arg is a
// string that contains contains.
func CheckPanic(t *testing.T, fn func(), contains string) bool {
	t.Helper()

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
		t.Error("panic() did not occur")
		return false
	}

	panicStr, ok := panicParam.(string)
	if !ok {
		t.Errorf("panic() occurred but recover()d %T type instead of string",
			panicParam)
		return false
	}

	if !strings.Contains(panicStr, contains) {
		t.Errorf("panic() string `%s'\ndoes not contain `%s'", panicStr, contains)
		return false
	}
	return true
}
