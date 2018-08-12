// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestString(t *testing.T) {
	checkOK(t, "foobar", testdeep.String("foobar"))

	type MyString string
	checkOK(t, MyString("foobar"), testdeep.String("foobar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), testdeep.String("pipo bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, testdeep.String("pipo bingo"))

	checkError(t, "foo bar test", testdeep.String("pipo"), expectedError{
		Message:  mustBe("does not match"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustContain(`"pipo"`),
	})

	checkError(t, 12, testdeep.String("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestHasPrefix(t *testing.T) {
	checkOK(t, "foobar", testdeep.HasPrefix("foo"))

	type MyString string
	checkOK(t, MyString("foobar"), testdeep.HasPrefix("foo"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), testdeep.HasPrefix("pipo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, testdeep.HasPrefix("pipo"))

	checkError(t, "foo bar test", testdeep.HasPrefix("pipo"), expectedError{
		Message:  mustBe("has not prefix"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustMatch(`^HasPrefix\(.*"pipo"`),
	})

	checkError(t, 12, testdeep.HasPrefix("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestHasSuffix(t *testing.T) {
	checkOK(t, "foobar", testdeep.HasSuffix("bar"))

	type MyString string
	checkOK(t, MyString("foobar"), testdeep.HasSuffix("bar"))

	// error interface
	checkOK(t, errors.New("pipo bingo"), testdeep.HasSuffix("bingo"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, testdeep.HasSuffix("bingo"))

	checkError(t, "foo bar test", testdeep.HasSuffix("pipo"), expectedError{
		Message:  mustBe("has not suffix"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustMatch(`^HasSuffix\(.*"pipo"`),
	})

	checkError(t, 12, testdeep.HasSuffix("bar"), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
	})
}

func TestStringTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.String("x"), nil)
	equalTypes(t, testdeep.HasPrefix("x"), nil)
	equalTypes(t, testdeep.HasSuffix("x"), nil)
}
