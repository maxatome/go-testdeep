// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestSet(t *testing.T) {
	type MyArray [5]int
	type MySlice []int

	for idx, got := range []interface{}{
		[]int{1, 3, 4, 4, 5},
		[...]int{1, 3, 4, 4, 5},
		MySlice{1, 3, 4, 4, 5},
		MyArray{1, 3, 4, 4, 5},
		&MySlice{1, 3, 4, 4, 5},
		&MyArray{1, 3, 4, 4, 5},
	} {
		testName := fmt.Sprintf("Test #%d → %v", idx, got)

		//
		// Set
		checkOK(t, got, testdeep.Set(5, 4, 1, 3), testName)
		checkOK(t, got,
			testdeep.Set(5, 4, 1, 3, 3, 3, 3), testName) // duplicated fields
		checkOK(t, got,
			testdeep.Set(
				testdeep.Between(0, 5),
				testdeep.Between(0, 5),
				testdeep.Between(0, 5))) // dup too

		checkError(t, got, testdeep.Set(5, 4),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				// items are sorted
				Summary: mustBe(`Extra 2 items: (1,
                3)`),
			},
			testName)

		checkError(t, got, testdeep.Set(5, 4, 1, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		checkError(t, got, testdeep.Set(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		checkError(t, got, testdeep.Set(5, 67, 4, 1, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing 2 items: (66,\n                  67)"),
			},
			testName)

		checkError(t, got, testdeep.Set(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)\n  Extra item: (1)"),
			},
			testName)

		//
		// SubSetOf
		checkOK(t, got, testdeep.SubSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, testdeep.SubSetOf(5, 4, 1, 3, 66), testName)

		checkError(t, got, testdeep.SubSetOf(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a SubSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (1)"),
			},
			testName)

		//
		// SuperSetOf
		checkOK(t, got, testdeep.SuperSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, testdeep.SuperSetOf(5, 4), testName)

		checkError(t, got, testdeep.SuperSetOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SuperSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		//
		// NotAny
		checkOK(t, got, testdeep.NotAny(10, 20, 30), testName)

		checkError(t, got, testdeep.NotAny(3, 66),
			expectedError{
				Message: mustBe("comparing %% as a NotAny"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (3)"),
			},
			testName)
	}

	checkOK(t, []interface{}{123, "foo", nil, "bar", nil},
		testdeep.Set("foo", "bar", 123, nil))

	var nilSlice MySlice
	for idx, got := range []interface{}{([]int)(nil), &nilSlice} {
		testName := fmt.Sprintf("Test #%d", idx)

		checkOK(t, got, testdeep.Set(), testName)
		checkOK(t, got, testdeep.SubSetOf(), testName)
		checkOK(t, got, testdeep.SubSetOf(1, 2), testName)
		checkOK(t, got, testdeep.SuperSetOf(), testName)
		checkOK(t, got, testdeep.NotAny(), testName)
		checkOK(t, got, testdeep.NotAny(1, 2), testName)
	}

	for idx, set := range []testdeep.TestDeep{
		testdeep.Set(123),
		testdeep.SubSetOf(123),
		testdeep.SuperSetOf(123),
		testdeep.NotAny(123),
	} {
		testName := fmt.Sprintf("Test #%d → %s", idx, set)

		checkError(t, 123, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		num := 123
		checkError(t, &num, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		var list *MySlice
		checkError(t, list, set,
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *testdeep_test.MySlice"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		checkError(t, nil, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)
	}

	//
	// String
	test.EqualStr(t, testdeep.Set(1).String(), "Set(1)")
	test.EqualStr(t, testdeep.Set(1, 2).String(), "Set(1,\n    2)")

	test.EqualStr(t, testdeep.SubSetOf(1).String(), "SubSetOf(1)")
	test.EqualStr(t, testdeep.SubSetOf(1, 2).String(), "SubSetOf(1,\n         2)")

	test.EqualStr(t, testdeep.SuperSetOf(1).String(), "SuperSetOf(1)")
	test.EqualStr(t, testdeep.SuperSetOf(1, 2).String(),
		"SuperSetOf(1,\n           2)")

	test.EqualStr(t, testdeep.NotAny(1).String(), "NotAny(1)")
	test.EqualStr(t, testdeep.NotAny(1, 2).String(), "NotAny(1,\n       2)")
}

func TestSetTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Set(6), nil)
	equalTypes(t, testdeep.SubSetOf(6), nil)
	equalTypes(t, testdeep.SuperSetOf(6), nil)
	equalTypes(t, testdeep.NotAny(6), nil)
}
