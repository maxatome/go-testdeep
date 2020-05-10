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

func TestBad(t *testing.T) {
	defer ctxerr.SaveColorState()()

	test.EqualStr(t, ctxerr.Bad("test"), "test")
	test.EqualStr(t, ctxerr.Bad("test %d", 123), "test 123")
}

func TestBadUsage(t *testing.T) {
	defer ctxerr.SaveColorState()()

	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", nil, 1, true),
		"usage: Zzz(STRING), but received nil as 1st parameter")

	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", 42, 1, true),
		"usage: Zzz(STRING), but received int as 1st parameter")

	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", []int{}, 1, true),
		"usage: Zzz(STRING), but received []int (slice) as 1st parameter")
	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", []int{}, 1, false),
		"usage: Zzz(STRING), but received []int as 1st parameter")

	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", nil, 1, true),
		"usage: Zzz(STRING), but received nil as 1st parameter")
	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", nil, 2, true),
		"usage: Zzz(STRING), but received nil as 2nd parameter")
	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", nil, 3, true),
		"usage: Zzz(STRING), but received nil as 3rd parameter")
	test.EqualStr(t,
		ctxerr.BadUsage("Zzz(STRING)", nil, 4, true),
		"usage: Zzz(STRING), but received nil as 4th parameter")
}

func TestTooManyParams(t *testing.T) {
	defer ctxerr.SaveColorState()()

	test.EqualStr(t, ctxerr.TooManyParams("Zzz(PARAM)"),
		"usage: Zzz(PARAM), too many parameters")
}

func TestError(t *testing.T) {
	defer ctxerr.SaveColorState()()

	err := ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message:  "Error message",
		Got:      1,
		Expected: 2,
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	     got: 1
	expected: 2`)
	test.EqualStr(t, err.GotString(), "1")
	test.EqualStr(t, err.ExpectedString(), "2")
	test.EqualStr(t, err.SummaryString(), "")

	err.Message = "Value of %% differ"
	test.EqualStr(t, err.Error(),
		`Value of DATA[12].Field differ
	     got: 1
	expected: 2`)

	err.Message = "Path at end: %%"
	test.EqualStr(t, err.Error(),
		`Path at end: DATA[12].Field
	     got: 1
	expected: 2`)

	err.Message = "%% <- the path!"
	test.EqualStr(t, err.Error(),
		`DATA[12].Field <- the path!
	     got: 1
	expected: 2`)

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
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
	     got: 1
	expected: 2
[under TestDeep operator Operator at file.go:23]`)

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message: "Error message",
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
			Message: "Origin error message",
			Summary: ctxerr.NewSummary("42"),
			Location: location.Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		42
	[under TestDeep operator SubOperator at file2.go:236]
[under TestDeep operator Operator at file.go:23]`)
	test.EqualStr(t, err.GotString(), "")
	test.EqualStr(t, err.ExpectedString(), "")
	test.EqualStr(t, err.SummaryString(), "666")

	err = ctxerr.Error{
		Context: ctxerr.Context{
			Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field"),
		},
		Message: "Error message",
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
			Message: "Origin error message",
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
			Message: "Error message",
			Summary: ctxerr.NewSummary("888"),
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 23,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		42
	[under TestDeep operator SubOperator at file2.go:236]
DATA[13].Field: Error message
	888
[under TestDeep operator Operator at file.go:23]`)

	err = ctxerr.Error{
		Context: ctxerr.Context{Path: ctxerr.NewPath("DATA").AddArrayIndex(12).AddField("Field")},
		Message: "Error message",
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
			Message: "Origin error message",
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
			Message: "Error message",
			Summary: ctxerr.NewSummary("888"),
			Location: location.Location{
				File: "file.go",
				Func: "Operator",
				Line: 24,
			},
		},
	}
	test.EqualStr(t, err.Error(),
		`DATA[12].Field: Error message
	666
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		42
	[under TestDeep operator SubOperator at file2.go:236]
[under TestDeep operator Operator at file.go:23]
DATA[13].Field: Error message
	888
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
