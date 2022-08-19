// Copyright (c) 2018-2022, Maxime Soulé
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

func TestSet(t *testing.T) {
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
		// Set
		checkOK(t, got, td.Set(5, 4, 1, 3), testName)
		checkOK(t, got,
			td.Set(5, 4, 1, 3, 3, 3, 3), testName) // duplicated fields
		checkOK(t, got,
			td.Set(
				td.Between(0, 5),
				td.Between(0, 5),
				td.Between(0, 5))) // dup too

		checkError(t, got, td.Set(5, 4),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				// items are sorted
				Summary: mustBe(`Extra 2 items: (1,
                3)`),
			},
			testName)

		checkError(t, got, td.Set(5, 4, 1, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		checkError(t, got, td.Set(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		checkError(t, got, td.Set(5, 67, 4, 1, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing 2 items: (66,\n                  67)"),
			},
			testName)

		checkError(t, got, td.Set(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)\n  Extra item: (1)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.Set(5, float64(4), 1, 3)), testName)

		//
		// SubSetOf
		checkOK(t, got, td.SubSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, td.SubSetOf(5, 4, 1, 3, 66), testName)

		checkError(t, got, td.SubSetOf(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a SubSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (1)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.SubSetOf(5, float64(4), 1, 3)), testName)

		//
		// SuperSetOf
		checkOK(t, got, td.SuperSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, td.SuperSetOf(5, 4), testName)

		checkError(t, got, td.SuperSetOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SuperSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (66)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.SuperSetOf(5, float64(4), 1, 3)), testName)

		//
		// NotAny
		checkOK(t, got, td.NotAny(10, 20, 30), testName)

		checkError(t, got, td.NotAny(3, 66),
			expectedError{
				Message: mustBe("comparing %% as a NotAny"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (3)"),
			},
			testName)

		// Lax
		checkOK(t, got, td.NotAny(float64(3)), testName)

		checkError(t, got, td.Lax(td.NotAny(float64(3))),
			expectedError{
				Message: mustBe("comparing %% as a NotAny"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra item: (3.0)"),
			},
			testName)
	}

	checkOK(t, []any{123, "foo", nil, "bar", nil},
		td.Set("foo", "bar", 123, nil))

	var nilSlice MySlice
	for idx, got := range []any{([]int)(nil), &nilSlice} {
		testName := fmt.Sprintf("Test #%d", idx)

		checkOK(t, got, td.Set(), testName)
		checkOK(t, got, td.SubSetOf(), testName)
		checkOK(t, got, td.SubSetOf(1, 2), testName)
		checkOK(t, got, td.SuperSetOf(), testName)
		checkOK(t, got, td.NotAny(), testName)
		checkOK(t, got, td.NotAny(1, 2), testName)
	}

	for idx, set := range []td.TestDeep{
		td.Set(123),
		td.SubSetOf(123),
		td.SuperSetOf(123),
		td.NotAny(123),
	} {
		testName := fmt.Sprintf("Test #%d → %s", idx, set)

		checkError(t, 123, set,
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			},
			testName)

		num := 123
		checkError(t, &num, set,
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*int"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			},
			testName)

		var list *MySlice
		checkError(t, list, set,
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *slice (*td_test.MySlice type)"),
				Expected: mustBe("non-nil *slice OR *array"),
			},
			testName)

		checkError(t, nil, set,
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
	test.EqualStr(t, td.Set(1).String(), "Set(1)")
	test.EqualStr(t, td.Set(1, 2).String(), "Set(1,\n    2)")

	test.EqualStr(t, td.SubSetOf(1).String(), "SubSetOf(1)")
	test.EqualStr(t, td.SubSetOf(1, 2).String(), "SubSetOf(1,\n         2)")

	test.EqualStr(t, td.SuperSetOf(1).String(), "SuperSetOf(1)")
	test.EqualStr(t, td.SuperSetOf(1, 2).String(),
		"SuperSetOf(1,\n           2)")

	test.EqualStr(t, td.NotAny(1).String(), "NotAny(1)")
	test.EqualStr(t, td.NotAny(1, 2).String(), "NotAny(1,\n       2)")
}

func TestSetTypeBehind(t *testing.T) {
	equalTypes(t, td.Set(6, 5), ([]int)(nil))
	equalTypes(t, td.Set(6, "foo"), nil)

	equalTypes(t, td.SubSetOf(6, 5), ([]int)(nil))
	equalTypes(t, td.SubSetOf(6, "foo"), nil)

	equalTypes(t, td.SuperSetOf(6, 5), ([]int)(nil))
	equalTypes(t, td.SuperSetOf(6, "foo"), nil)

	equalTypes(t, td.NotAny(6, 5), ([]int)(nil))
	equalTypes(t, td.NotAny(6, "foo"), nil)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		td.Set(
			td.Empty(),
			5,
			td.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			td.All(6, 7),
			td.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		([]int)(nil))

	// Only one interface type
	equalTypes(t,
		td.Set(
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
		),
		([]*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		td.Set(
			td.Isa((*error)(nil)),
			td.Isa((*fmt.Stringer)(nil)),
		),
		nil)
}
