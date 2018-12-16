// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestContains(t *testing.T) {
	type (
		MySlice  []int
		MyArray  [3]int
		MyMap    map[string]int
		MyString string
	)

	for idx, got := range []interface{}{
		[]int{12, 34, 28},
		MySlice{12, 34, 28},
		[...]int{12, 34, 28},
		MyArray{12, 34, 28},
		map[string]int{"foo": 12, "bar": 34, "zip": 28},
		MyMap{"foo": 12, "bar": 34, "zip": 28},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.Contains(34), testName)
		checkOK(t, got, testdeep.Contains(testdeep.Between(30, 35)), testName)

		checkError(t, got, testdeep.Contains(35),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain("34"), // as well as other items in fact...
				Expected: mustBe("Contains(35)"),
			}, testName)
	}

	for idx, got := range []interface{}{
		"foobar",
		MyString("foobar"),
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.Contains(testdeep.Between('n', 'p')), testName)

		checkError(t, got, testdeep.Contains(testdeep.Between('y', 'z')),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`"foobar"`), // as well as other items in fact...
				Expected: mustBe(fmt.Sprintf("Contains(%d ≤ got ≤ %d)", 'y', 'z')),
			}, testName)
	}
}

// nil case
func TestContainsNil(t *testing.T) {
	type (
		MyPtrSlice []*int
		MyPtrArray [3]*int
		MyPtrMap   map[string]*int
	)

	num := 12345642
	for idx, got := range []interface{}{
		[]*int{&num, nil},
		MyPtrSlice{&num, nil},
		[...]*int{&num, nil},
		MyPtrArray{&num},
		map[string]*int{"foo": &num, "bar": nil},
		MyPtrMap{"foo": &num, "bar": nil},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.Contains(nil), testName)
		checkOK(t, got, testdeep.Contains((*int)(nil)), testName)
		checkOK(t, got, testdeep.Contains(testdeep.Nil()), testName)
		checkOK(t, got, testdeep.Contains(testdeep.NotNil()), testName)

		checkError(t, got, testdeep.Contains((*uint8)(nil)),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain("12345642"),
				Expected: mustBe("Contains((*uint8)(<nil>))"),
			}, testName)
	}

	for idx, got := range []interface{}{
		[]interface{}{nil, 12345642},
		[]func(){nil, func() {}},
		[][]int{{}, nil},
		[...]interface{}{nil, 12345642},
		[...]func(){nil, func() {}},
		[...][]int{{}, nil},
		map[bool]interface{}{true: nil, false: 12345642},
		map[bool]func(){true: nil, false: func() {}},
		map[bool][]int{true: {}, false: nil},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.Contains(nil), testName)
		checkOK(t, got, testdeep.Contains(testdeep.Nil()), testName)
		checkOK(t, got, testdeep.Contains(testdeep.NotNil()), testName)
	}

	for idx, got := range []interface{}{
		[]int{1, 2, 3},
		[...]int{1, 2, 3},
		map[string]int{"foo": 12, "bar": 34, "zip": 28},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkError(t, got, testdeep.Contains(nil),
			expectedError{
				Message: mustBe("does not contain"),
				Path:    mustBe("DATA"),
				// Got
				Expected: mustBe("Contains(nil)"),
			}, testName)
	}

	checkError(t, "foobar", testdeep.Contains(nil),
		expectedError{
			Message:  mustBe("cannot check contains"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("Contains(nil)"),
		})

	// Caught by deepValueEqual, before Match() call
	checkError(t, nil, testdeep.Contains(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Contains(nil)"),
		})
}

func TestContainsString(t *testing.T) {
	type MyString string

	for idx, got := range []interface{}{
		"pipo bingo",
		MyString("pipo bingo"),
		errors.New("pipo bingo"), // error interface
		MyStringer{},             // fmt.Stringer interface
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.Contains("po bi"), testName)
		checkOK(t, got, testdeep.Contains('o'), testName)
		checkOK(t, got, testdeep.Contains(byte('o')), testName)

		checkError(t, got, testdeep.Contains("zip"),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`"pipo bingo"`),
				Expected: mustMatch(`^Contains\(.*"zip"`),
			})

		checkError(t, got, testdeep.Contains(12),
			expectedError{
				Message:  mustBe("cannot check contains"),
				Path:     mustBe("DATA"),
				Got:      mustBe(reflect.TypeOf(got).String()),
				Expected: mustBe("int"),
			})
	}

	checkError(t, 12, testdeep.Contains("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR fmt.Stringer OR error"),
		})
}

func TestContainsTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Contains("x"), nil)
}
