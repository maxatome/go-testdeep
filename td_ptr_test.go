// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestPtr(t *testing.T) {
	//
	// Ptr
	num := 12
	str := "test"
	pNum := &num
	pStr := &str
	pStruct := &struct{}{}

	checkOK(t, &num, testdeep.Ptr(12))
	checkOK(t, &str, testdeep.Ptr("test"))
	checkOK(t, &struct{}{}, testdeep.Ptr(struct{}{}))

	checkError(t, &num, testdeep.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("13"),
		})
	checkError(t, nil, testdeep.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("*int"),
		})
	checkError(t, (*int)(nil), testdeep.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"), // should be DATA, but seems hard to be done
			Got:      mustBe("nil"),
			Expected: mustBe("13"),
		})
	checkError(t, (*int)(nil), testdeep.Ptr((*int)(nil)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("**int"),
		})
	checkError(t, &num, testdeep.Ptr("test"),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("*string"),
		})

	checkError(t, &num, testdeep.Ptr(testdeep.Any(11)),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("*DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("Any(11)"),
		})

	checkError(t, &str, testdeep.Ptr("foobar"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"),
			Got:      mustContain(`"test"`),
			Expected: mustContain(`"foobar"`),
		})

	checkError(t, 13, testdeep.Ptr(13),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*int"),
		})
	checkError(t, &str, testdeep.Ptr(12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*string"),
			Expected: mustBe("*int"),
		})

	checkError(t, &pNum, testdeep.Ptr(12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("**int"),
			Expected: mustBe("*int"),
		})
	checkError(t, &pStr, testdeep.Ptr("test"),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("**string"),
			Expected: mustBe("*string"),
		})

	//
	// PPtr
	checkOK(t, &pNum, testdeep.PPtr(12))
	checkOK(t, &pStr, testdeep.PPtr("test"))
	checkOK(t, &pStruct, testdeep.PPtr(struct{}{}))

	checkError(t, &pNum, testdeep.PPtr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("**DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("13"),
		})
	checkError(t, nil, testdeep.PPtr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("**int"),
		})
	checkError(t, &num, testdeep.PPtr(13),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("**int"),
		})

	checkError(t, &pStr, testdeep.PPtr("foobar"),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("**DATA"),
		})
	checkError(t, &str, testdeep.PPtr("foobar"),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*string"),
			Expected: mustBe("**string"),
		})

	checkError(t, &pNum, testdeep.PPtr(testdeep.Any(11)),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("**DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("Any(11)"),
		})

	pStruct = nil
	checkError(t, &pStruct, testdeep.PPtr(struct{}{}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("**DATA"), // should be *DATA, but seems hard to be done
			Got:      mustBe("nil"),
			Expected: mustContain("struct"),
		})

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Ptr(nil) }, "usage: Ptr(")
	test.CheckPanic(t, func() { testdeep.Ptr(MyInterface(nil)) }, "usage: Ptr(")

	test.CheckPanic(t, func() { testdeep.PPtr(nil) }, "usage: PPtr(")
	test.CheckPanic(t, func() { testdeep.PPtr(MyInterface(nil)) }, "usage: PPtr(")

	//
	// String
	test.EqualStr(t, testdeep.Ptr(13).String(), "*int")
	test.EqualStr(t, testdeep.PPtr(13).String(), "**int")
	test.EqualStr(t, testdeep.Ptr(testdeep.Ptr(13)).String(), "*<something>")
	test.EqualStr(t, testdeep.PPtr(testdeep.Ptr(13)).String(), "**<something>")
}

func TestPtrTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Ptr(6), nil)
	equalTypes(t, testdeep.PPtr(6), nil)
}
