// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestList(t *testing.T) {
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

		checkOK(t, got, td.List(1, 3, 4, 4, 5), testName)

		typ := reflect.TypeOf(got)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		checkError(t, got, td.List(1, 3, 4),
			expectedError{
				Message: mustBe(fmt.Sprintf("comparing %s, from index #3", typ.Kind())),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra 2 items: (4,\n                5)"),
			},
			testName)

		checkError(t, got, td.List(1, 3, 4, 4, 5, 666),
			expectedError{
				Message: mustBe(fmt.Sprintf("comparing %s, from index #5", typ.Kind())),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing item: (666)"),
			},
			testName)

		checkError(t, got, td.List(1, 3, 666, 4, 5),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA[2]"),
				Got:      mustBe("4"),
				Expected: mustBe("666"),
			},
			testName)

		// Lax
		checkOK(t, got, td.Lax(td.Bag(float64(1), 3, 4, 4, uint8(5))), testName)
	}

	var nilSlice MySlice
	for idx, got := range []any{([]int)(nil), &nilSlice} {
		checkOK(t, got, td.List(), "Test #%d", idx)
	}

	checkError(t, 123, td.List(123),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("slice OR array OR *slice OR *array"),
		})

	num := 123
	checkError(t, &num, td.List(123),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("slice OR array OR *slice OR *array"),
		})

	var list *MySlice
	checkError(t, list, td.List(123),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *slice (*td_test.MySlice type)"),
			Expected: mustBe("non-nil *slice OR *array"),
		})

	checkError(t, nil, td.List(123),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("slice OR array OR *slice OR *array"),
		})

	//
	// String
	test.EqualStr(t, td.List(1).String(), "List(1)")
	test.EqualStr(t, td.List(1, 2).String(), "List(1,\n     2)")
}

func TestListTypeBehind(t *testing.T) {
	equalTypes(t, td.List(6, 5), ([]int)(nil))
	equalTypes(t, td.List(6, "foo"), nil)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		td.List(
			td.Empty(),
			5,
			td.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			td.All(6, 7),
			td.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		([]int)(nil))

	// Only one interface type
	equalTypes(t,
		td.List(
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
		),
		([]*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		td.List(
			td.Isa((*error)(nil)),
			td.Isa((*fmt.Stringer)(nil)),
		),
		nil)
}
