// Copyright (c) 2018-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/location"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

func TestError(t *testing.T) {
	defer color.SaveState()()

	checkWithoutColors := func(err *ctxerr.Error) {
		t.Helper()
		test.EqualStr(t, err.ErrorWithoutColors(), err.Error())
		defer color.SaveState(true)()
		test.IsTrue(t, err.ErrorWithoutColors() != err.Error())
	}

	err := ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message:  "error message",
		Got:      1,
		Expected: 2,
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: error message
	     got: 1
	expected: 2`)
	checkWithoutColors(&err)
	test.EqualStr(t, err.GotString(), "1")
	test.EqualStr(t, err.ExpectedString(), "2")
	test.EqualStr(t, err.SummaryString(), "")

	err.Message = "Value of %% differ"
	test.EqualStr(t, err.Error(),
		`Value of DATA[12].Field differ
	     got: 1
	expected: 2`)
	checkWithoutColors(&err)

	err.Message = "Path at end: %%"
	test.EqualStr(t, err.Error(),
		`Path at end: DATA[12].Field
	     got: 1
	expected: 2`)
	checkWithoutColors(&err)

	err.Message = "%% <- the path!"
	test.EqualStr(t, err.Error(),
		`DATA[12].Field <- the path!
	     got: 1
	expected: 2`)
	checkWithoutColors(&err)

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message:  "error message",
		Got:      1,
		Expected: 2,
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: error message
	     got: 1
	expected: 2
[under operator Operator at file.go:23]`)
	checkWithoutColors(&err)

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message: "error message",
		Summary: ctxerr.NewSummary("666"),
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{
				Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field").AddCustomLevel("<All#1/2>"),
			},
			Message: "origin error message",
			Summary: ctxerr.NewSummary("42"),
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: origin error message
		42
	[under operator SubOperator at file2.go:236]
[under operator Operator at file.go:23]`)
	checkWithoutColors(&err)
	test.EqualStr(t, err.GotString(), "")
	test.EqualStr(t, err.ExpectedString(), "")
	test.EqualStr(t, err.SummaryString(), "666")

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message: "error message",
		Summary: ctxerr.NewSummary("666"),
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{
				Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field").AddCustomLevel("<All#1/2>"),
			},
			Message: "origin error message",
			Summary: ctxerr.NewSummary("42"),
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
		// Next error at same location
		Next: &ctxerr.Error{
			Context: ctxerr.Context{
				Path: ctxerr.NewPath("DATA").AddArrayIndex(13).AddField("Field"),
			},
			Message: "error message",
			Summary: ctxerr.NewSummary("888"),
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 23,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: origin error message
		42
	[under operator SubOperator at file2.go:236]
DATA[13].Field: error message
	888
[under operator Operator at file.go:23]`)
	checkWithoutColors(&err)

	err = ctxerr.Error{
		Context: ctxerr.Context{Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field")},
		Message: "error message",
		Summary: ctxerr.NewSummary("666"),
		Location: location.Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &ctxerr.Error{
			Context: ctxerr.Context{
				Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field").AddCustomLevel("<All#1/2>"),
			},
			Message: "origin error message",
			Summary: ctxerr.NewSummary("42"),
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
		// Next error at different location
		Next: &ctxerr.Error{
			Context: ctxerr.Context{
				Path: ctxerr.NewPath("DATA").AddArrayIndex(13).AddField("Field"),
			},
			Message: "error message",
			Summary: ctxerr.NewSummary("888"),
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 24,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: origin error message
		42
	[under operator SubOperator at file2.go:236]
[under operator Operator at file.go:23]
DATA[13].Field: error message
	888
[under operator Operator at file.go:24]`)
	checkWithoutColors(&err)

	//
	// ErrTooManyErrors
	test.EqualStr(t, ctxerr.ErrTooManyErrors.Error(),
		`Too many errors (use TESTDEEP_MAX_ERRORS=-1 to see all)`)
}

func TestTypeMismatch(t *testing.T) {
	rErr := ctxerr.TypeMismatch(reflect.TypeOf(0), reflect.TypeOf(""))
	test.EqualStr(t, rErr.Message, "type mismatch")
	test.EqualStr(t, string(rErr.Got.(types.RawString)), `int`)
	test.EqualStr(t, string(rErr.Expected.(types.RawString)), `string`)

	// It is the caller responsibility to check that both types
	// differ. To ease testing we can pass twice the same type, it is
	// the same as passing 2 different types but with the same short
	// name (a/too.Type vs b/foo.Type), util.TypeFullName() is called
	// for both types.
	rErr = ctxerr.TypeMismatch(reflect.TypeOf(0), reflect.TypeOf(0))
	test.EqualStr(t, rErr.Message, "type mismatch")
	test.EqualStr(t, string(rErr.Got.(types.RawString)), `int`)
	test.EqualStr(t, string(rErr.Expected.(types.RawString)), `int`)
}

func TestBooleanError(t *testing.T) {
	if ctxerr.BooleanError.Error() != "" {
		t.Errorf("BooleanError should stringify to empty string, not `%s'",
			ctxerr.BooleanError.Error())
	}
}
