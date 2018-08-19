// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/location"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestError(t *testing.T) {
	err := ctxerr.Error{
		Context:  ctxerr.Context{Path: "DATA[12].Field"},
		Message:  "Error message",
		Got:      1,
		Expected: 2,
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	     got: (int) 1
	expected: (int) 2`)
	test.EqualStr(t, err.GotString(), "(int) 1")
	test.EqualStr(t, err.ExpectedString(), "(int) 2")
	test.EqualStr(t, err.SummaryString(), "")

	err.Message = "Value of %% differ"
	test.EqualStr(t, err.Error(),
		`Value of DATA[12].Field differ
	     got: (int) 1
	expected: (int) 2`)

	err.Message = "Path at end: %%"
	test.EqualStr(t, err.Error(),
		`Path at end: DATA[12].Field
	     got: (int) 1
	expected: (int) 2`)

	err.Message = "%% <- the path!"
	test.EqualStr(t, err.Error(),
		`DATA[12].Field <- the path!
	     got: (int) 1
	expected: (int) 2`)

	err = ctxerr.Error{
		Context:  ctxerr.Context{Path: "DATA[12].Field"},
		Message:  "Error message",
		Got:      1,
		Expected: 2,
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	     got: (int) 1
	expected: (int) 2
[under TestDeep operator Operator at file.go:23]`)

	err = ctxerr.Error{
		Context: ctxerr.Context{Path: "DATA[12].Field"},
		Message: "Error message",
		Summary: 666,
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{Path: "DATA[12].Field<All#1/2>"},
			Message: "Origin error message",
			Summary: 42,
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	(int) 666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		(int) 42
	[under TestDeep operator SubOperator at file2.go:236]
[under TestDeep operator Operator at file.go:23]`)
	test.EqualStr(t, err.GotString(), "")
	test.EqualStr(t, err.ExpectedString(), "")
	test.EqualStr(t, err.SummaryString(), "(int) 666")

	err = ctxerr.Error{
		Context: ctxerr.Context{Path: "DATA[12].Field"},
		Message: "Error message",
		Summary: 666,
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{Path: "DATA[12].Field<All#1/2>"},
			Message: "Origin error message",
			Summary: 42,
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
		// Next error at same location
		Next: &ctxerr.Error{
			Context: ctxerr.Context{Path: "DATA[13].Field"},
			Message: "Error message",
			Summary: 888,
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 23,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	(int) 666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		(int) 42
	[under TestDeep operator SubOperator at file2.go:236]
DATA[13].Field: Error message
	(int) 888
[under TestDeep operator Operator at file.go:23]`)

	err = ctxerr.Error{
		Context: ctxerr.Context{Path: "DATA[12].Field"},
		Message: "Error message",
		Summary: 666,
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{Path: "DATA[12].Field<All#1/2>"},
			Message: "Origin error message",
			Summary: 42,
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
		// Next error at different location
		Next: &ctxerr.Error{
			Context: ctxerr.Context{Path: "DATA[13].Field"},
			Message: "Error message",
			Summary: 888,
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 24,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	(int) 666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		(int) 42
	[under TestDeep operator SubOperator at file2.go:236]
[under TestDeep operator Operator at file.go:23]
DATA[13].Field: Error message
	(int) 888
[under TestDeep operator Operator at file.go:24]`)

	//
	// ErrTooManyErrors
	test.EqualStr(t, ctxerr.ErrTooManyErrors.Error(),
		`Too many errors (use TESTDEEP_MAX_ERRORS=-1 to see all)`)
}

func TestBooleanError(t *testing.T) {
	if ctxerr.BooleanError.Error() != "" {
		t.Errorf("BooleanError should stringify to empty string, not `%s'",
			ctxerr.BooleanError.Error())
	}
}
