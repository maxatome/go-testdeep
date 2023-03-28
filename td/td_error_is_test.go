// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type errorIsSimpleErr string

func (e errorIsSimpleErr) Error() string {
	return string(e)
}

type errorIsWrappedErr struct {
	s   string
	err error
}

func (e errorIsWrappedErr) Error() string {
	if e.err != nil {
		return e.s + ": " + e.err.Error()
	}
	return e.s + ": nil"
}

func (e errorIsWrappedErr) Unwrap() error {
	return e.err
}

var _ = []error{errorIsSimpleErr(""), errorIsWrappedErr{}}

func TestErrorIs(t *testing.T) {
	insideErr1 := errorIsSimpleErr("failure1")
	insideErr2 := errorIsWrappedErr{"failure2", insideErr1}
	insideErr3 := errorIsWrappedErr{"failure3", insideErr2}
	err := errorIsWrappedErr{"failure4", insideErr3}

	checkOK(t, err, td.ErrorIs(err))
	checkOK(t, err, td.ErrorIs(insideErr3))
	checkOK(t, err, td.ErrorIs(insideErr2))
	checkOK(t, err, td.ErrorIs(insideErr1))
	checkOK(t, nil, td.ErrorIs(nil))

	checkOK(t, err, td.ErrorIs(td.All(
		td.Isa(errorIsSimpleErr("")),
		td.String("failure1"),
	)))

	// many errorIsWrappedErr in the err's tree, so only the first
	// encountered matches
	checkOK(t, err, td.ErrorIs(td.All(
		td.Isa(errorIsWrappedErr{}),
		td.HasPrefix("failure4"),
	)))

	// HasPrefix().TypeBehind() always returns nil
	// so errors.As() is called with &any, so the toplevel error matches
	checkOK(t, err, td.ErrorIs(td.HasPrefix("failure4")))

	var errNil error
	checkOK(t, &errNil, td.Ptr(td.ErrorIs(nil)))

	var inside errorIsSimpleErr
	checkOK(t, err, td.ErrorIs(td.Catch(&inside, td.String("failure1"))))
	test.EqualStr(t, string(inside), "failure1")

	checkError(t, nil, td.ErrorIs(insideErr1),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("nil value"),
			Got:      mustBe("nil"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, 45, td.ErrorIs(insideErr1),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("int does not implement error interface"),
			Got:      mustBe("45"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, 45, td.ErrorIs(fmt.Errorf("another")),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("int does not implement error interface"),
			Got:      mustBe("45"),
			Expected: mustBe("anything implementing error interface"),
		})

	checkError(t, err, td.ErrorIs(fmt.Errorf("another")),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("is not found in err's tree"),
			Got:      mustBe(`(td_test.errorIsWrappedErr) "failure4: failure3: failure2: failure1"`),
			Expected: mustBe(`(*errors.errorString) "another"`),
		})

	checkError(t, err, td.ErrorIs(td.String("nonono")),
		expectedError{
			Path:     mustBe("DATA.ErrorIs(interface {})"),
			Message:  mustBe("does not match"),
			Got:      mustBe(`"failure4: failure3: failure2: failure1"`),
			Expected: mustBe(`"nonono"`),
		})

	checkError(t, err, td.ErrorIs(td.Isa(fmt.Errorf("another"))),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("type is not found in err's tree"),
			Got:      mustBe(`(td_test.errorIsWrappedErr) failure4: failure3: failure2: failure1`),
			Expected: mustBe(`*errors.errorString`),
		})

	checkError(t, err, td.ErrorIs(td.Smuggle(io.ReadAll, td.String("xx"))),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("type is not found in err's tree"),
			Got:      mustBe(`(td_test.errorIsWrappedErr) failure4: failure3: failure2: failure1`),
			Expected: mustBe(`io.Reader`),
		})

	checkError(t, err, td.ErrorIs(nil),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("is not nil"),
			Got:      mustBe(`(td_test.errorIsWrappedErr) "failure4: failure3: failure2: failure1"`),
			Expected: mustBe(`nil`),
		})

	// As errors.Is, it does not match
	checkError(t, errorIsWrappedErr{"failure", nil}, td.ErrorIs(nil),
		expectedError{
			Path:     mustBe("DATA"),
			Message:  mustBe("is not nil"),
			Got:      mustBe(`(td_test.errorIsWrappedErr) "failure: nil"`),
			Expected: mustBe(`nil`),
		})

	checkError(t, err, td.ErrorIs(td.Gt(0)),
		expectedError{
			Path:    mustBe("DATA"),
			Message: mustBe("bad usage of ErrorIs operator"),
			Summary: mustBe(`ErrorIs(Gt): type int behind Gt operator is not an interface or does not implement error`),
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
	test.EqualStr(t, td.ErrorIs(td.HasPrefix("pipo")).String(),
		`ErrorIs(HasPrefix("pipo"))`)
	test.EqualStr(t, td.ErrorIs(12).String(), "ErrorIs(<ERROR>)")
}

func TestErrorIsTypeBehind(t *testing.T) {
	equalTypes(t, td.ErrorIs(fmt.Errorf("another")), nil)
}
