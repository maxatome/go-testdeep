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

func EqualStr(t *testing.T, got, expected string, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %s
	expected: %s`,
		BuildTestName(args), got, expected)
	return false
}

func EqualInt(t *testing.T, got, expected int, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %d
	expected: %d`,
		BuildTestName(args), got, expected)
	return false
}

func EqualBool(t *testing.T, got, expected bool, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %t
	expected: %t`,
		BuildTestName(args), got, expected)
	return false
}

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

func IsTrue(t *testing.T, got bool, args ...interface{}) bool {
	if got {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: false
	expected: true`, BuildTestName(args))
	return false
}

func IsFalse(t *testing.T, got bool, args ...interface{}) bool {
	if !got {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: true
	expected: false`, BuildTestName(args))
	return false
}

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
