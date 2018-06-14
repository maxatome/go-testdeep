// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestGetInterface(t *testing.T) {
	type Private struct {
		private *Private // nolint: megacheck
	}

	// Cases not tested by TestEqualOthers()
	s := Private{
		private: &Private{},
	}

	_, ok := getInterface(reflect.ValueOf(s).Field(0), false)
	if ok {
		t.Error("getInterface() should return false for private field")
	}

	_, ok = getInterface(reflect.ValueOf(s).Field(0), true)
	if UnsafeDisabled {
		if ok {
			t.Error("unsafe package is disabled, getInterface should fail")
		}
	} else if !ok {
		t.Error("unsafe package is available, getInterface should succeed")
	}

	//var (
	//	panicked   bool
	//	panicParam interface{}
	//)
	//
	//func() {
	//	defer func() { panicParam = recover() }()
	//	panicked = true
	//	mustGetInterface(reflect.ValueOf(s).Field(0))
	//	panicked = false
	//}()
	//
	//if panicked {
	//	panicStr, ok := panicParam.(string)
	//	if ok {
	//		const expectedPanic = "getInterface() does not handle map kind"
	//		if panicStr != expectedPanic {
	//			t.Errorf("panic() string `%s' ≠ `%s'", panicStr, expectedPanic)
	//		}
	//	} else {
	//		t.Errorf("panic() occurred but recover()d %T type instead of string",
	//			panicParam)
	//	}
	//} else {
	//	t.Error("panic() did not occur")
	//}
}

type TestTestingT struct {
	LastMessage string
}

func (t *TestTestingT) Error(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
}

func (t *TestTestingT) Helper() {
	// Do nothing
}

func TestFormatError(t *testing.T) {
	ttt := &TestTestingT{}

	err := &Error{
		Context: NewContext(),
		Message: "test error message",
		Summary: rawString("test error summary"),
	}

	nonStringName := bytes.NewBufferString("zip!")

	//
	// Without args
	formatError(ttt, err)
	equalStr(t, ttt.LastMessage, `Failed test
DATA: test error message
	test error summary`)

	//
	// With one arg
	formatError(ttt, err, "foo bar!")
	equalStr(t, ttt.LastMessage, `Failed test 'foo bar!'
DATA: test error message
	test error summary`)

	formatError(ttt, err, nonStringName)
	equalStr(t, ttt.LastMessage, `Failed test 'zip!'
DATA: test error message
	test error summary`)

	//
	// With several args & Printf format
	formatError(ttt, err, "hello %d!", 123)
	equalStr(t, ttt.LastMessage, `Failed test 'hello 123!'
DATA: test error message
	test error summary`)

	//
	// With several args without Printf format
	formatError(ttt, err, "hello ", "world! ", 123)
	equalStr(t, ttt.LastMessage, `Failed test 'hello world! 123'
DATA: test error message
	test error summary`)

	formatError(ttt, err, nonStringName, "hello ", "world! ", 123)
	equalStr(t, ttt.LastMessage, `Failed test 'zip!hello world! 123'
DATA: test error message
	test error summary`)
}
