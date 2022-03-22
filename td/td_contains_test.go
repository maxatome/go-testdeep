// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestContains(t *testing.T) {
	type (
		MySlice  []int
		MyArray  [3]int
		MyMap    map[string]int
		MyString string
	)

	for idx, got := range []any{
		[]int{12, 34, 28},
		MySlice{12, 34, 28},
		[...]int{12, 34, 28},
		MyArray{12, 34, 28},
		map[string]int{"foo": 12, "bar": 34, "zip": 28},
		MyMap{"foo": 12, "bar": 34, "zip": 28},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.Contains(34), testName)
		checkOK(t, got, td.Contains(td.Between(30, 35)), testName)

		checkError(t, got, td.Contains(35),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain("34"), // as well as other items in fact...
				Expected: mustBe("Contains(35)"),
			}, testName)

		// Lax
		checkOK(t, got, td.Lax(td.Contains(float64(34))), testName)
	}

	for idx, got := range []any{
		"foobar",
		MyString("foobar"),
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.Contains(td.Between('n', 'p')), testName)

		checkError(t, got, td.Contains(td.Between('y', 'z')),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`"foobar"`), // as well as other items in fact...
				Expected: mustBe(fmt.Sprintf("Contains((int32) %d ≤ got ≤ (int32) %d)", 'y', 'z')),
			}, testName)
	}
}

// nil case.
func TestContainsNil(t *testing.T) {
	type (
		MyPtrSlice []*int
		MyPtrArray [3]*int
		MyPtrMap   map[string]*int
	)

	num := 12345642
	for idx, got := range []any{
		[]*int{&num, nil},
		MyPtrSlice{&num, nil},
		[...]*int{&num, nil},
		MyPtrArray{&num},
		map[string]*int{"foo": &num, "bar": nil},
		MyPtrMap{"foo": &num, "bar": nil},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.Contains(nil), testName)
		checkOK(t, got, td.Contains((*int)(nil)), testName)
		checkOK(t, got, td.Contains(td.Nil()), testName)
		checkOK(t, got, td.Contains(td.NotNil()), testName)

		checkError(t, got, td.Contains((*uint8)(nil)),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain("12345642"),
				Expected: mustBe("Contains((*uint8)(<nil>))"),
			}, testName)
	}

	for idx, got := range []any{
		[]any{nil, 12345642},
		[]func(){nil, func() {}},
		[][]int{{}, nil},
		[...]any{nil, 12345642},
		[...]func(){nil, func() {}},
		[...][]int{{}, nil},
		map[bool]any{true: nil, false: 12345642},
		map[bool]func(){true: nil, false: func() {}},
		map[bool][]int{true: {}, false: nil},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.Contains(nil), testName)
		checkOK(t, got, td.Contains(td.Nil()), testName)
		checkOK(t, got, td.Contains(td.NotNil()), testName)
	}

	for idx, got := range []any{
		[]int{1, 2, 3},
		[...]int{1, 2, 3},
		map[string]int{"foo": 12, "bar": 34, "zip": 28},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkError(t, got, td.Contains(nil),
			expectedError{
				Message: mustBe("does not contain"),
				Path:    mustBe("DATA"),
				// Got
				Expected: mustBe("Contains(nil)"),
			}, testName)
	}

	checkError(t, "foobar", td.Contains(nil),
		expectedError{
			Message:  mustBe("cannot check contains"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("Contains(nil)"),
		})

	// Caught by deepValueEqual, before Match() call
	checkError(t, nil, td.Contains(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Contains(nil)"),
		})
}

func TestContainsString(t *testing.T) {
	type MyString string

	for idx, got := range []any{
		"pipo bingo",
		MyString("pipo bingo"),
		[]byte("pipo bingo"),
		errors.New("pipo bingo"), // error interface
		MyStringer{},             // fmt.Stringer interface
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.Contains("pipo"), testName)
		checkOK(t, got, td.Contains("po bi"), testName)
		checkOK(t, got, td.Contains("bingo"), testName)

		checkOK(t, got, td.Contains([]byte("pipo")), testName)
		checkOK(t, got, td.Contains([]byte("po bi")), testName)
		checkOK(t, got, td.Contains([]byte("bingo")), testName)

		checkOK(t, got, td.Contains('o'), testName)
		checkOK(t, got, td.Contains(byte('o')), testName)

		checkOK(t, got, td.Contains(""), testName)
		checkOK(t, got, td.Contains([]byte{}), testName)

		if _, ok := got.([]byte); ok {
			checkOK(t, got,
				td.Contains(td.Code(func(b byte) bool { return b == 'o' })),
				testName)
		} else {
			checkOK(t, got,
				td.Contains(td.Code(func(r rune) bool { return r == 'o' })),
				testName)
		}

		checkError(t, got, td.Contains("zip"),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`pipo bingo`),
				Expected: mustMatch(`^Contains\(.*"zip"`),
			})

		checkError(t, got, td.Contains([]byte("zip")),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`pipo bingo`),
				Expected: mustMatch(`^(?s)Contains\(.*zip`),
			})

		checkError(t, got, td.Contains('z'),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`pipo bingo`),
				Expected: mustBe(`Contains((int32) 122)`),
			})

		checkError(t, got, td.Contains(byte('z')),
			expectedError{
				Message:  mustBe("does not contain"),
				Path:     mustBe("DATA"),
				Got:      mustContain(`pipo bingo`),
				Expected: mustBe(`Contains((uint8) 122)`),
			})

		checkError(t, got, td.Contains(12),
			expectedError{
				Message:  mustBe("cannot check contains"),
				Path:     mustBe("DATA"),
				Got:      mustBe(reflect.TypeOf(got).String()),
				Expected: mustBe("int"),
			})

		checkError(t, got, td.Contains([]int{1, 2, 3}),
			expectedError{
				Message:  mustBe("cannot check contains"),
				Path:     mustBe("DATA"),
				Got:      mustBe(reflect.TypeOf(got).String()),
				Expected: mustBe("[]int"),
			})

		// Lax
		checkOK(t, got,
			td.Lax(td.Contains(td.Code(func(b int) bool { return b == 'o' }))),
			testName)
	}

	checkError(t, 12, td.Contains("bar"),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR []byte (convertible) OR fmt.Stringer OR error"),
		})

	checkError(t, "pipo", td.Contains(td.Code(func(x int) bool { return true })),
		expectedError{
			Message:  mustBe("Code operator has to match rune in string, but it does not"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("rune"),
		})
}

func TestContainsSlice(t *testing.T) {
	got := []int{1, 2, 3, 4, 5, 6}

	// Empty slice is always OK
	checkOK(t, got, td.Contains([]int{}))

	// Expected length > got length
	checkError(t, got, td.Contains([]int{1, 2, 3, 4, 5, 6, 7}),
		expectedError{
			Message:  mustBe("does not contain"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`([]int) (len=6 `),
			Expected: mustContain(`Contains(([]int) (len=7 `),
		})

	// Same length
	checkOK(t, got, td.Contains([]int{1, 2, 3, 4, 5, 6}))
	checkError(t, got, td.Contains([]int{8, 8, 8, 8, 8, 8}),
		expectedError{
			Message:  mustBe("does not contain"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`([]int) (len=6 `),
			Expected: mustContain(`Contains(([]int) (len=6 `),
		})

	checkOK(t, got, td.Contains([]int{1, 2, 3}))
	checkOK(t, got, td.Contains([]int{3, 4, 5}))
	checkOK(t, got, td.Contains([]int{4, 5, 6}))

	checkError(t, got, td.Contains([]int{8, 8, 8}),
		expectedError{
			Message:  mustBe("does not contain"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`([]int) (len=6 `),
			Expected: mustContain(`Contains(([]int) (len=3 `),
		})
}

func TestContainsTypeBehind(t *testing.T) {
	equalTypes(t, td.Contains("x"), nil)
}
