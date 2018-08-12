// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestRe(t *testing.T) {
	//
	// string
	checkOK(t, "foo bar test", testdeep.Re("bar"))
	checkOK(t, "foo bar test", testdeep.Re(regexp.MustCompile("test$")))

	checkOK(t, "foo bar test",
		testdeep.ReAll(`(\w+)`, testdeep.Bag("bar", "test", "foo")))

	type MyString string
	checkOK(t, MyString("Ho zz hoho"),
		testdeep.ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))

	// error interface
	checkOK(t, errors.New("pipo bingo"), testdeep.Re("bin"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, testdeep.Re("bin"))

	checkError(t, 12, testdeep.Re("bar"),
		expectedError{
			Message: mustBe("bad type"),
			Path:    mustBe("DATA"),
			Got:     mustBe("int"),
			Expected: mustBe(
				"string (convertible) OR fmt.Stringer OR error OR []uint8"),
		})

	checkError(t, "foo bar test", testdeep.Re("pipo"),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("pipo"),
		})

	checkError(t, "foo bar test", testdeep.Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// bytes
	checkOK(t, []byte("foo bar test"), testdeep.Re("bar"))

	checkOK(t, []byte("foo bar test"),
		testdeep.ReAll(`(\w+)`, testdeep.Bag("bar", "test", "foo")))

	type MySlice []byte
	checkOK(t, MySlice("Ho zz hoho"),
		testdeep.ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))

	checkError(t, []int{12}, testdeep.Re("bar"),
		expectedError{
			Message:  mustBe("bad slice type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("[]uint8"),
		})

	checkError(t, []byte("foo bar test"), testdeep.Re("pipo"),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("pipo"),
		})

	checkError(t, []byte("foo bar test"),
		testdeep.Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// Bad usage
	const reUsage = "usage: Re("
	test.CheckPanic(t, func() { testdeep.Re(123) }, reUsage)
	test.CheckPanic(t, func() { testdeep.Re("bar", []string{}, 1) }, reUsage)

	const reAllUsage = "usage: ReAll("
	test.CheckPanic(t, func() { testdeep.ReAll(123, 456) }, reAllUsage)
}

func TestReTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Re("x"), nil)
}
