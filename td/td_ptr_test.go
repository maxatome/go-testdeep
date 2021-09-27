// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestPtr(t *testing.T) {
	//
	// Ptr
	num := 12
	str := "test"
	pNum := &num
	pStr := &str
	pStruct := &struct{}{}

	checkOK(t, &num, td.Ptr(12))
	checkOK(t, &str, td.Ptr("test"))
	checkOK(t, &struct{}{}, td.Ptr(struct{}{}))

	checkError(t, &num, td.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("13"),
		})
	checkError(t, nil, td.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("*int"),
		})
	checkError(t, (*int)(nil), td.Ptr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"), // should be DATA, but seems hard to be done
			Got:      mustBe("nil"),
			Expected: mustBe("13"),
		})
	checkError(t, (*int)(nil), td.Ptr((*int)(nil)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("**int"),
		})
	checkError(t, &num, td.Ptr("test"),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("*string"),
		})

	checkError(t, &num, td.Ptr(td.Any(11)),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("*DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("Any(11)"),
		})

	checkError(t, &str, td.Ptr("foobar"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"),
			Got:      mustContain(`"test"`),
			Expected: mustContain(`"foobar"`),
		})

	checkError(t, 13, td.Ptr(13),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*int"),
		})
	checkError(t, &str, td.Ptr(12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*string"),
			Expected: mustBe("*int"),
		})

	checkError(t, &pNum, td.Ptr(12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("**int"),
			Expected: mustBe("*int"),
		})
	checkError(t, &pStr, td.Ptr("test"),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("**string"),
			Expected: mustBe("*string"),
		})

	//
	// PPtr
	checkOK(t, &pNum, td.PPtr(12))
	checkOK(t, &pStr, td.PPtr("test"))
	checkOK(t, &pStruct, td.PPtr(struct{}{}))

	checkError(t, &pNum, td.PPtr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("**DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("13"),
		})
	checkError(t, nil, td.PPtr(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("**int"),
		})
	checkError(t, &num, td.PPtr(13),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("**int"),
		})

	checkError(t, &pStr, td.PPtr("foobar"),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("**DATA"),
		})
	checkError(t, &str, td.PPtr("foobar"),
		expectedError{
			Message:  mustBe("pointer type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*string"),
			Expected: mustBe("**string"),
		})

	checkError(t, &pNum, td.PPtr(td.Any(11)),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("**DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("Any(11)"),
		})

	pStruct = nil
	checkError(t, &pStruct, td.PPtr(struct{}{}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("**DATA"), // should be *DATA, but seems hard to be done
			Got:      mustBe("nil"),
			Expected: mustContain("struct"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Ptr(nil),
		expectedError{
			Message: mustBe("bad usage of Ptr operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("usage: Ptr("),
		})

	checkError(t, "never tested",
		td.Ptr(MyInterface(nil)),
		expectedError{
			Message: mustBe("bad usage of Ptr operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("usage: Ptr("),
		})

	checkError(t, "never tested",
		td.PPtr(nil),
		expectedError{
			Message: mustBe("bad usage of PPtr operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("usage: PPtr("),
		})

	checkError(t, "never tested",
		td.PPtr(MyInterface(nil)),
		expectedError{
			Message: mustBe("bad usage of PPtr operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("usage: PPtr("),
		})

	//
	// String
	test.EqualStr(t, td.Ptr(13).String(), "*int")
	test.EqualStr(t, td.PPtr(13).String(), "**int")
	test.EqualStr(t, td.Ptr(td.Ptr(13)).String(), "*<something>")
	test.EqualStr(t, td.PPtr(td.Ptr(13)).String(), "**<something>")

	// Erroneous op
	test.EqualStr(t, td.Ptr(nil).String(), "Ptr(<ERROR>)")
	test.EqualStr(t, td.PPtr(nil).String(), "PPtr(<ERROR>)")
}

func TestPtrTypeBehind(t *testing.T) {
	var num int
	equalTypes(t, td.Ptr(6), &num)

	// Another TestDeep operator delegation
	var num64 int64
	equalTypes(t, td.Ptr(td.Between(int64(1), int64(2))), &num64)
	equalTypes(t, td.Ptr(td.Any(1, 1.2)), nil)

	// Erroneous op
	equalTypes(t, td.Ptr(nil), nil)
}

func TestPPtrTypeBehind(t *testing.T) {
	var pnum *int
	equalTypes(t, td.PPtr(6), &pnum)

	// Another TestDeep operator delegation
	var pnum64 *int64
	equalTypes(t, td.PPtr(td.Between(int64(1), int64(2))), &pnum64)
	equalTypes(t, td.PPtr(td.Any(1, 1.2)), nil)

	// Erroneous op
	equalTypes(t, td.PPtr(nil), nil)
}
