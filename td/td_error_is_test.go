// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestErrorIs(t *testing.T) {
	insideErr1 := fmt.Errorf("failure1")
	insideErr2 := fmt.Errorf("failure2: %w", insideErr1)
	insideErr3 := fmt.Errorf("failure3: %w", insideErr2)
	err := fmt.Errorf("failure4: %w", insideErr3)

	checkOK(t, err, td.ErrorIs(err))
	checkOK(t, err, td.ErrorIs(insideErr3))
	checkOK(t, err, td.ErrorIs(insideErr2))
	checkOK(t, err, td.ErrorIs(insideErr1))
	checkOK(t, nil, td.ErrorIs(nil))

	var errNil error
	checkOK(t, &errNil, td.Ptr(td.ErrorIs(nil)))

	checkError(t, nil, td.ErrorIs(insideErr1),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, 45, td.ErrorIs(insideErr1),
		expectedError{
			Message:  mustBe("int does not implement error interface"),
			Path:     mustBe("DATA"),
			Got:      mustBe("45"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, 45, td.ErrorIs(fmt.Errorf("another")),
		expectedError{
			Message:  mustBe("int does not implement error interface"),
			Path:     mustBe("DATA"),
			Got:      mustBe("45"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, err, td.ErrorIs(fmt.Errorf("another")),
		expectedError{
			Message:  mustBe("is not the error"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`(*fmt.wrapError) "failure4: failure3: failure2: failure1"`),
			Expected: mustBe(`(*errors.errorString) "another"`),
		})

	checkError(t, err, td.ErrorIs(nil),
		expectedError{
			Message:  mustBe("is not the error"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`(*fmt.wrapError) "failure4: failure3: failure2: failure1"`),
			Expected: mustBe(`nil`),
		})

	type private struct{ err error }
	got := private{err: err}
	for _, expErr := range []error{err, insideErr3} {
		expected := td.Struct(private{}, td.StructFields{"err": td.ErrorIs(expErr)})
		if dark.UnsafeDisabled {
			checkError(t, got, expected,
				expectedError{
					Message: mustBe("cannot compare"),
					Path:    mustBe("DATA.err"),
					Summary: mustBe("unexported field that cannot be overridden"),
				})
		} else {
			checkOK(t, got, expected)
		}
	}

	if !dark.UnsafeDisabled {
		got = private{}
		checkOK(t, got, td.Struct(private{}, td.StructFields{"err": td.ErrorIs(nil)}))
	}

	//
	// String
	test.EqualStr(t, td.ErrorIs(insideErr1).String(), "ErrorIs(failure1)")
	test.EqualStr(t, td.ErrorIs(nil).String(), "ErrorIs(nil)")
}

func TestErrorIsTypeBehind(t *testing.T) {
	equalTypes(t, td.ErrorIs(fmt.Errorf("another")), nil)
}
