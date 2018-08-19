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

func TestBag(t *testing.T) {
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
		// Bag
		checkOK(t, got, testdeep.Bag(5, 4, 1, 4, 3), testName)

		checkError(t, got, testdeep.Bag(5, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra items: ((int) 4)"),
			},
			testName)

		checkError(t, got, testdeep.Bag(5, 4, 1, 4, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)

		checkError(t, got, testdeep.Bag(5, 66, 4, 1, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)

		checkError(t, got, testdeep.Bag(5, 66, 4, 1, 4, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66,\n                (int) 66)"),
			},
			testName)

		checkError(t, got, testdeep.Bag(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)\n  Extra items: ((int) 4)"),
			},
			testName)

		//
		// SubBagOf
		checkOK(t, got, testdeep.SubBagOf(5, 4, 1, 4, 3), testName)
		checkOK(t, got, testdeep.SubBagOf(5, 66, 4, 1, 4, 3), testName)

		checkError(t, got, testdeep.SubBagOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SubBagOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra items: ((int) 4)"),
			},
			testName)

		//
		// SuperBagOf
		checkOK(t, got, testdeep.SuperBagOf(5, 4, 1, 4, 3), testName)
		checkOK(t, got, testdeep.SuperBagOf(5, 4, 3), testName)

		checkError(t, got, testdeep.SuperBagOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SuperBagOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)
	}

	checkOK(t, []interface{}{123, "foo", nil, "bar", nil},
		testdeep.Bag("foo", "bar", 123, nil, nil))

	var nilSlice MySlice
	for idx, got := range []interface{}{([]int)(nil), &nilSlice} {
		testName := fmt.Sprintf("Test #%d", idx)

		checkOK(t, got, testdeep.Bag(), testName)
		checkOK(t, got, testdeep.SubBagOf(), testName)
		checkOK(t, got, testdeep.SubBagOf(1, 2), testName)
		checkOK(t, got, testdeep.SuperBagOf(), testName)
	}

	for idx, bag := range []testdeep.TestDeep{
		testdeep.Bag(123),
		testdeep.SubBagOf(123),
		testdeep.SuperBagOf(123),
	} {
		testName := fmt.Sprintf("Test #%d → %s", idx, bag)

		checkError(t, 123, bag,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		num := 123
		checkError(t, &num, bag,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		var list *MySlice
		checkError(t, list, bag,
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *testdeep_test.MySlice"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		checkError(t, nil, bag,
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
	test.EqualStr(t, testdeep.Bag(1).String(), "Bag((int) 1)")
	test.EqualStr(t, testdeep.Bag(1, 2).String(), "Bag((int) 1,\n    (int) 2)")

	test.EqualStr(t, testdeep.SubBagOf(1).String(), "SubBagOf((int) 1)")
	test.EqualStr(t, testdeep.SubBagOf(1, 2).String(), "SubBagOf((int) 1,\n         (int) 2)")

	test.EqualStr(t, testdeep.SuperBagOf(1).String(), "SuperBagOf((int) 1)")
	test.EqualStr(t, testdeep.SuperBagOf(1, 2).String(),
		"SuperBagOf((int) 1,\n           (int) 2)")
}

func TestBagTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Bag(6), nil)
	equalTypes(t, testdeep.SubBagOf(6), nil)
	equalTypes(t, testdeep.SuperBagOf(6), nil)
}
