// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestRe(t *testing.T) {
	//
	// string
	checkOK(t, "foo bar test", td.Re("bar"))
	checkOK(t, "foo bar test", td.Re(regexp.MustCompile("test$")))

	checkOK(t, "foo bar test",
		td.ReAll(`(\w+)`, td.Bag("bar", "test", "foo")))

	type MyString string
	checkOK(t, MyString("Ho zz hoho"),
		td.ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))
	checkOK(t, MyString("Ho zz hoho"),
		td.ReAll("(?i)(ho)", []any{"Ho", "ho", "ho"}))

	// error interface
	checkOK(t, errors.New("pipo bingo"), td.Re("bin"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, td.Re("bin"))

	checkError(t, 12, td.Re("bar"),
		expectedError{
			Message: mustBe("bad type"),
			Path:    mustBe("DATA"),
			Got:     mustBe("int"),
			Expected: mustBe(
				"string (convertible) OR fmt.Stringer OR error OR []uint8"),
		})

	checkError(t, "foo bar test", td.Re("pipo"),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("pipo"),
		})

	checkError(t, "foo bar test", td.Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("(pi)(po)"),
		})
	checkError(t, "foo bar test", td.Re("(pi)(po)", []any{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// bytes
	checkOK(t, []byte("foo bar test"), td.Re("bar"))

	checkOK(t, []byte("foo bar test"),
		td.ReAll(`(\w+)`, td.Bag("bar", "test", "foo")))

	type MySlice []byte
	checkOK(t, MySlice("Ho zz hoho"),
		td.ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))
	checkOK(t, MySlice("Ho zz hoho"),
		td.ReAll("(?i)(ho)", []any{"Ho", "ho", "ho"}))

	checkError(t, []int{12}, td.Re("bar"),
		expectedError{
			Message:  mustBe("bad slice type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("[]uint8"),
		})

	checkError(t, []byte("foo bar test"), td.Re("pipo"),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("pipo"),
		})

	checkError(t, []byte("foo bar test"),
		td.Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("(pi)(po)"),
		})
	checkError(t, []byte("foo bar test"),
		td.Re("(pi)(po)", []any{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// Bad usage
	const (
		ur = "(STRING|*regexp.Regexp[, NON_NIL_CAPTURE])"
		ua = "(STRING|*regexp.Regexp, NON_NIL_CAPTURE)"
	)

	checkError(t, "never tested",
		td.Re(123),
		expectedError{
			Message: mustBe("bad usage of Re operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Re" + ur + ", but received int as 1st parameter"),
		})

	checkError(t, "never tested",
		td.ReAll(123, nil),
		expectedError{
			Message: mustBe("bad usage of ReAll operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: ReAll" + ua + ", but received int as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Re("bar", []string{}, 1),
		expectedError{
			Message: mustBe("bad usage of Re operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Re" + ur + ", too many parameters"),
		})

	checkError(t, "never tested",
		td.ReAll(123, 456),
		expectedError{
			Message: mustBe("bad usage of ReAll operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: ReAll" + ua + ", but received int as 1st parameter"),
		})

	checkError(t, "never tested",
		td.ReAll(`12[3,4`, nil),
		expectedError{
			Message: mustBe("invalid regexp given to ReAll operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("error parsing regexp: "),
		})

	// Erroneous op
	test.EqualStr(t, td.Re(123).String(), "Re(<ERROR>)")
	test.EqualStr(t, td.ReAll(123, nil).String(), "ReAll(<ERROR>)")
}

func TestReTypeBehind(t *testing.T) {
	equalTypes(t, td.Re("x"), nil)
	equalTypes(t, td.ReAll("x", nil), nil)

	// Erroneous op
	equalTypes(t, td.Re(123), nil)
	equalTypes(t, td.ReAll(123, nil), nil)
}
