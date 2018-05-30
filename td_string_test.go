// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestString(t *testing.T) {
	checkOK(t, "foobar", String("foobar"))

	type MyString string
	checkOK(t, MyString("foobar"), String("foobar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), String("pipo bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, String("pipo bingo"))

	checkError(t, "foo bar test", String("pipo"), expectedError{
		Message:  mustBe("does not match"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustContain(`"pipo"`),
	})

	checkError(t, 12, String("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestHasPrefix(t *testing.T) {
	checkOK(t, "foobar", HasPrefix("foo"))

	type MyString string
	checkOK(t, MyString("foobar"), HasPrefix("foo"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), HasPrefix("pipo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, HasPrefix("pipo"))

	checkError(t, "foo bar test", HasPrefix("pipo"), expectedError{
		Message:  mustBe("has not prefix"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustMatch(`^HasPrefix\(.*"pipo"`),
	})

	checkError(t, 12, HasPrefix("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestHasSuffix(t *testing.T) {
	checkOK(t, "foobar", HasSuffix("bar"))

	type MyString string
	checkOK(t, MyString("foobar"), HasSuffix("bar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), HasSuffix("bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, HasSuffix("bingo"))

	checkError(t, "foo bar test", HasSuffix("pipo"), expectedError{
		Message:  mustBe("has not suffix"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustMatch(`^HasSuffix\(.*"pipo"`),
	})

	checkError(t, 12, HasSuffix("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestContains(t *testing.T) {
	checkOK(t, "foobar", Contains("ooba"))

	type MyString string
	checkOK(t, MyString("foobar"), Contains("ooba"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), Contains("po bi"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, Contains("po bi"))

	checkError(t, "foo bar test", Contains("pipo"), expectedError{
		Message:  mustBe("does not contain"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustMatch(`^Contains\(.*"pipo"`),
	})

	checkError(t, 12, Contains("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestStringTypeOf(t *testing.T) {
	equalTypes(t, String("x"), nil)
	equalTypes(t, HasPrefix("x"), nil)
	equalTypes(t, HasSuffix("x"), nil)
	equalTypes(t, Contains("x"), nil)
}
