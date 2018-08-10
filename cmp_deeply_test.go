// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestFormatError(t *testing.T) {
	ttt := &TestTestingT{}

	err := &Error{
		Context: NewContext(),
		Message: "test error message",
		Summary: rawString("test error summary"),
	}

	nonStringName := bytes.NewBufferString("zip!")

	for _, fatal := range []bool{false, true} {
		//
		// Without args
		formatError(ttt, fatal, err)
		test.EqualStr(t, ttt.LastMessage, `Failed test
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With one arg
		formatError(ttt, fatal, err, "foo bar!")
		test.EqualStr(t, ttt.LastMessage, `Failed test 'foo bar!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		formatError(ttt, fatal, err, nonStringName)
		test.EqualStr(t, ttt.LastMessage, `Failed test 'zip!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args & Printf format
		formatError(ttt, fatal, err, "hello %d!", 123)
		test.EqualStr(t, ttt.LastMessage, `Failed test 'hello 123!'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		//
		// With several args without Printf format
		formatError(ttt, fatal, err, "hello ", "world! ", 123)
		test.EqualStr(t, ttt.LastMessage, `Failed test 'hello world! 123'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)

		formatError(ttt, fatal, err, nonStringName, "hello ", "world! ", 123)
		test.EqualStr(t, ttt.LastMessage, `Failed test 'zip!hello world! 123'
DATA: test error message
	test error summary`)
		test.EqualBool(t, ttt.IsFatal, fatal)
	}
}
