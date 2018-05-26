// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestError(t *testing.T) {
	err := Error{
		Context:  NewContext("DATA[12].Field"),
		Message:  "Error message",
		Got:      1,
		Expected: 2,
	}
	equalStr(t, err.Error(),
		`DATA[12].Field: Error message
	     got: (int) 1
	expected: (int) 2`)

	err.Message = "Value of %% differ"
	equalStr(t, err.Error(),
		`Value of DATA[12].Field differ
	     got: (int) 1
	expected: (int) 2`)

	err.Message = "Path at end: %%"
	equalStr(t, err.Error(),
		`Path at end: DATA[12].Field
	     got: (int) 1
	expected: (int) 2`)

	err.Message = "%% <- the path!"
	equalStr(t, err.Error(),
		`DATA[12].Field <- the path!
	     got: (int) 1
	expected: (int) 2`)

	err = Error{
		Context:  NewContext("DATA[12].Field"),
		Message:  "Error message",
		Got:      1,
		Expected: 2,
		Location: Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
	}
	equalStr(t, err.Error(),
		`DATA[12].Field: Error message
	     got: (int) 1
	expected: (int) 2
[under TestDeep operator Operator at file.go:23]`)

	err = Error{
		Context: NewContext("DATA[12].Field"),
		Message: "Error message",
		Summary: 666,
		Location: Location{
			File: "file.go",
			Func: "Operator",
			Line: 23,
		},
		Origin: &Error{
			Context: NewContext("DATA[12].Field<All#1/2>"),
			Message: "Origin error message",
			Summary: 42,
			Location: Location{
				File: "file2.go",
				Func: "SubOperator",
				Line: 236,
			},
		},
	}
	equalStr(t, err.Error(),
		`DATA[12].Field: Error message
	(int) 666
[under TestDeep operator Operator at file.go:23]
Originates from following error:
	DATA[12].Field<All#1/2>: Origin error message
		(int) 42
	[under TestDeep operator SubOperator at file2.go:236]`)
}
