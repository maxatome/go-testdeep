// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestString(t *testing.T) {
	checkOK(t, "foobar", td.String("foobar"))
	checkOK(t, []byte("foobar"), td.String("foobar"))

	type MyBytes []byte
	checkOK(t, MyBytes("foobar"), td.String("foobar"))

	type MyString string
	checkOK(t, MyString("foobar"), td.String("foobar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), td.String("pipo bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, td.String("pipo bingo"))

	checkError(t, "foo bar test", td.String("pipo"),
		expectedError{
			Message:  mustBe("does not match"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustContain(`"pipo"`),
		})

	checkError(t, []int{1, 2}, td.String("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})

	checkError(t, 12, td.String("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})
}

func TestHasPrefix(t *testing.T) {
	checkOK(t, "foobar", td.HasPrefix("foo"))
	checkOK(t, []byte("foobar"), td.HasPrefix("foo"))

	type MyBytes []byte
	checkOK(t, MyBytes("foobar"), td.HasPrefix("foo"))

	type MyString string
	checkOK(t, MyString("foobar"), td.HasPrefix("foo"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), td.HasPrefix("pipo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, td.HasPrefix("pipo"))

	checkError(t, "foo bar test", td.HasPrefix("pipo"),
		expectedError{
			Message:  mustBe("has not prefix"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustMatch(`^HasPrefix\(.*"pipo"`),
		})

	checkError(t, []int{1, 2}, td.HasPrefix("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})

	checkError(t, 12, td.HasPrefix("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})
}

func TestHasSuffix(t *testing.T) {
	checkOK(t, "foobar", td.HasSuffix("bar"))
	checkOK(t, []byte("foobar"), td.HasSuffix("bar"))

	type MyBytes []byte
	checkOK(t, MyBytes("foobar"), td.HasSuffix("bar"))

	type MyString string
	checkOK(t, MyString("foobar"), td.HasSuffix("bar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), td.HasSuffix("bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, td.HasSuffix("bingo"))

	checkError(t, "foo bar test", td.HasSuffix("pipo"),
		expectedError{
			Message:  mustBe("has not suffix"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustMatch(`^HasSuffix\(.*"pipo"`),
		})

	checkError(t, []int{1, 2}, td.HasSuffix("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})

	checkError(t, 12, td.HasSuffix("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})
}

func TestStringTypeBehind(t *testing.T) {
	equalTypes(t, td.String("x"), nil)
	equalTypes(t, td.HasPrefix("x"), nil)
	equalTypes(t, td.HasSuffix("x"), nil)
}
