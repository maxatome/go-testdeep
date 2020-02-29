// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"bytes"
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestFormatError(t *testing.T) {
	defer ctxerr.SaveColorState()()

	ttt := test.NewTestingT()

	err := &ctxerr.Error{
		Context: newContext(),
		Message: "test error message",
		Summary: ctxerr.NewSummary("test error summary"),
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

func TestCmp(t *testing.T) {
	tt := test.NewTestingFT(t.Name())
	test.IsTrue(t, Cmp(tt, 1, 1))
	test.IsFalse(t, tt.Failed())

	tt = test.NewTestingFT(t.Name())
	test.IsFalse(t, Cmp(tt, 1, 2))
	test.IsTrue(t, tt.Failed())
}

func TestCmpDeeply(t *testing.T) {
	tt := test.NewTestingFT(t.Name())
	test.IsTrue(t, CmpDeeply(tt, 1, 1))
	test.IsFalse(t, tt.Failed())

	tt = test.NewTestingFT(t.Name())
	test.IsFalse(t, CmpDeeply(tt, 1, 2))
	test.IsTrue(t, tt.Failed())
}
