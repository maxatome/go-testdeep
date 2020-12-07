// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
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
		checkOK(t, got, td.Bag(5, 4, 1, 4, 3), testName)

		checkError(t, got, td.Bag(5, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (4)"),
			},
			testName)

		checkError(t, got, td.Bag(5, 4, 1, 4, 3, 66, 42),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				// items are sorted
				Summary: mustBe(`Missing 2 items: (42,
                  66)`),
			},
			testName)

		checkError(t, got, td.Bag(5, 66, 4, 1, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		checkError(t, got, td.Bag(5, 66, 4, 1, 4, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing 2 items: (66,\n                  66)"),
			},
			testName)

		checkError(t, got, td.Bag(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Bag"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)\n  Extra item: (4)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.Bag(float64(5), 4, 1, 4, 3)), testName)

		//
		// SubBagOf
		checkOK(t, got, td.SubBagOf(5, 4, 1, 4, 3), testName)
		checkOK(t, got, td.SubBagOf(5, 66, 4, 1, 4, 3), testName)

		checkError(t, got, td.SubBagOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SubBagOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (4)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.SubBagOf(float64(5), 4, 1, 4, 3)), testName)

		//
		// SuperBagOf
		checkOK(t, got, td.SuperBagOf(5, 4, 1, 4, 3), testName)
		checkOK(t, got, td.SuperBagOf(5, 4, 3), testName)

		checkError(t, got, td.SuperBagOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SuperBagOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.SuperBagOf(float64(5), 4, 1, 4, 3)), testName)
	}

	checkOK(t, []interface{}{123, "foo", nil, "bar", nil},
		td.Bag("foo", "bar", 123, nil, nil))

	var nilSlice MySlice
	for idx, got := range []interface{}{([]int)(nil), &nilSlice} {
		testName := fmt.Sprintf("Test #%d", idx)

		checkOK(t, got, td.Bag(), testName)
		checkOK(t, got, td.SubBagOf(), testName)
		checkOK(t, got, td.SubBagOf(1, 2), testName)
		checkOK(t, got, td.SuperBagOf(), testName)
	}

	for idx, bag := range []td.TestDeep{
		td.Bag(123),
		td.SubBagOf(123),
		td.SuperBagOf(123),
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
				Got:      mustBe("nil *td_test.MySlice"),
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
	test.EqualStr(t, td.Bag(1).String(), "Bag(1)")
	test.EqualStr(t, td.Bag(1, 2).String(), "Bag(1,\n    2)")

	test.EqualStr(t, td.SubBagOf(1).String(), "SubBagOf(1)")
	test.EqualStr(t, td.SubBagOf(1, 2).String(), "SubBagOf(1,\n         2)")

	test.EqualStr(t, td.SuperBagOf(1).String(), "SuperBagOf(1)")
	test.EqualStr(t, td.SuperBagOf(1, 2).String(),
		"SuperBagOf(1,\n           2)")
}

func TestBagTypeBehind(t *testing.T) {
	equalTypes(t, td.Bag(6), nil)
	equalTypes(t, td.SubBagOf(6), nil)
	equalTypes(t, td.SuperBagOf(6), nil)
}
