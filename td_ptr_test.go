// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
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

	checkOK(t, &num, Ptr(12))
	checkOK(t, &str, Ptr("test"))
	checkOK(t, &struct{}{}, Ptr(struct{}{}))

	checkError(t, &num, Ptr(13), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("*DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 13"),
	})
	checkError(t, nil, Ptr(13), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("*int"),
	})
	checkError(t, (*int)(nil), Ptr(13), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("*DATA"), // should be DATA, but seems hard to be done
		Got:      mustBe("nil"),
		Expected: mustBe("(int) 13"),
	})
	checkError(t, (*int)(nil), Ptr((*int)(nil)), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("*int"),
		Expected: mustBe("**int"),
	})
	checkError(t, &num, Ptr("test"), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("*int"),
		Expected: mustBe("*string"),
	})

	checkError(t, &num, Ptr(Any(11)), expectedError{
		Message:  mustBe("comparing with Any"),
		Path:     mustBe("*DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("Any((int) 11)"),
	})

	checkError(t, &str, Ptr("foobar"), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("*DATA"),
		Got:      mustContain(`"test"`),
		Expected: mustContain(`"foobar"`),
	})

	checkError(t, 13, Ptr(13), expectedError{
		Message:  mustBe("pointer type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("*int"),
	})
	checkError(t, &str, Ptr(12), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("*string"),
		Expected: mustBe("*int"),
	})

	checkError(t, &pNum, Ptr(12), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("**int"),
		Expected: mustBe("*int"),
	})
	checkError(t, &pStr, Ptr("test"), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("**string"),
		Expected: mustBe("*string"),
	})

	//
	// PPtr
	checkOK(t, &pNum, PPtr(12))
	checkOK(t, &pStr, PPtr("test"))
	checkOK(t, &pStruct, PPtr(struct{}{}))

	checkError(t, &pNum, PPtr(13), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("**DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 13"),
	})
	checkError(t, nil, PPtr(13), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("**int"),
	})
	checkError(t, &num, PPtr(13), expectedError{
		Message:  mustBe("pointer type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("*int"),
		Expected: mustBe("**int"),
	})

	checkError(t, &pStr, PPtr("foobar"), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("**DATA"),
	})
	checkError(t, &str, PPtr("foobar"), expectedError{
		Message:  mustBe("pointer type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("*string"),
		Expected: mustBe("**string"),
	})

	checkError(t, &pNum, PPtr(Any(11)), expectedError{
		Message:  mustBe("comparing with Any"),
		Path:     mustBe("**DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("Any((int) 11)"),
	})

	pStruct = nil
	checkError(t, &pStruct, PPtr(struct{}{}), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("**DATA"), // should be *DATA, but seems hard to be done
		Got:      mustBe("nil"),
		Expected: mustContain("struct"),
	})

	//
	// Bad usage
	checkPanic(t, func() { Ptr(nil) }, "usage: Ptr(")
	checkPanic(t, func() { Ptr(MyInterface(nil)) }, "usage: Ptr(")

	checkPanic(t, func() { PPtr(nil) }, "usage: PPtr(")
	checkPanic(t, func() { PPtr(MyInterface(nil)) }, "usage: PPtr(")

	//
	// String
	test.EqualStr(t, Ptr(13).String(), "*int")
	test.EqualStr(t, PPtr(13).String(), "**int")
	test.EqualStr(t, Ptr(Ptr(13)).String(), "*<something>")
	test.EqualStr(t, PPtr(Ptr(13)).String(), "**<something>")
}

func TestPtrTypeBehind(t *testing.T) {
	equalTypes(t, Ptr(6), nil)
	equalTypes(t, PPtr(6), nil)
}
