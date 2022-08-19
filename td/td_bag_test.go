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

	for idx, got := range []any{
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

	checkOK(t, []any{123, "foo", nil, "bar", nil},
		td.Bag("foo", "bar", 123, nil, nil))

	var nilSlice MySlice
	for idx, got := range []any{([]int)(nil), &nilSlice} {
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
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			},
			testName)

		num := 123
		checkError(t, &num, bag,
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			},
			testName)

		var list *MySlice
		checkError(t, list, bag,
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *slice (*td_test.MySlice type)"),
				Expected: mustBe("non-nil *slice OR *array"),
			},
			testName)

		checkError(t, nil, bag,
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
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
	equalTypes(t, td.Bag(6, 5), ([]int)(nil))
	equalTypes(t, td.Bag(6, "foo"), nil)

	equalTypes(t, td.SubBagOf(6, 5), ([]int)(nil))
	equalTypes(t, td.SubBagOf(6, "foo"), nil)

	equalTypes(t, td.SuperBagOf(6, 5), ([]int)(nil))
	equalTypes(t, td.SuperBagOf(6, "foo"), nil)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		td.Bag(
			td.Empty(),
			5,
			td.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			td.All(6, 7),
			td.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		([]int)(nil))

	// Only one interface type
	equalTypes(t,
		td.Bag(
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
		),
		([]*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		td.Bag(
			td.Isa((*error)(nil)),
			td.Isa((*fmt.Stringer)(nil)),
		),
		nil)
}
