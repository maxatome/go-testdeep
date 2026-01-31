// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestArray(t *testing.T) {
	type MyArray [5]int

	//
	// Simple array
	checkOK(t, [5]int{}, td.Array([5]int{}))
	checkOK(t, [5]int{}, td.Array([5]int{}, nil))
	checkOK(t, [5]int{0, 0, 0, 4}, td.Array([5]int{0, 0, 0, 4}))
	checkOK(t, [5]int{1, 0, 3},
		td.Array([5]int{}, td.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, [5]int{1, 2, 3},
		td.Array([5]int{0, 2}, td.ArrayEntries{2: 3, 0: 1}))

	checkOK(t, [5]any{1, 2, nil, 4, nil},
		td.Array([5]any{nil, 2, nil, 4}, td.ArrayEntries{0: 1, 2: nil}))

	checkOK(t, [5]any{1, 2, nil, 4, nil},
		td.Array([5]any{nil, 2, nil, 4},
			td.ArrayEntries{0: 3, 2: 28},
			td.ArrayEntries{0: 2, 2: "pipo"},
			td.ArrayEntries{0: 1, 2: nil},
		))

	zero, one, two := 0, 1, 2
	checkOK(t, [5]*int{nil, &zero, &one, &two},
		td.Array(
			[5]*int{}, td.ArrayEntries{1: &zero, 2: &one, 3: &two, 4: nil}))

	gotArray := [...]int{1, 2, 3, 4, 5}

	checkError(t, gotArray, td.Array(MyArray{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[5]int"),
			Expected: mustBe("td_test.MyArray"),
		})
	checkError(t, gotArray, td.Array([5]int{1, 2, 3, 4, 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})
	checkError(t, gotArray,
		td.Array([5]int{1, 2, 3, 4}, td.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})

	checkError(t, nil,
		td.Array([1]int{42}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustContain("Array("),
		})

	//
	// Array type
	checkOK(t, MyArray{}, td.Array(MyArray{}))
	checkOK(t, MyArray{}, td.Array(MyArray{}, nil))
	checkOK(t, MyArray{0, 0, 0, 4}, td.Array(MyArray{0, 0, 0, 4}))
	checkOK(t, MyArray{1, 0, 3},
		td.Array(MyArray{}, td.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, MyArray{1, 2, 3},
		td.Array(MyArray{0, 2}, td.ArrayEntries{2: 3, 0: 1}))

	checkOK(t, &MyArray{}, td.Array(&MyArray{}))
	checkOK(t, &MyArray{0, 0, 0, 4}, td.Array(&MyArray{0, 0, 0, 4}))
	checkOK(t, &MyArray{1, 0, 3},
		td.Array(&MyArray{}, td.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 0, 3},
		td.Array((*MyArray)(nil), td.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 2, 3},
		td.Array(&MyArray{0, 2}, td.ArrayEntries{2: 3, 0: 1}))

	gotTypedArray := MyArray{1, 2, 3, 4, 5}

	checkError(t, 123, td.Array(&MyArray{}, td.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*td_test.MyArray"),
		})

	checkError(t, &MyStruct{},
		td.Array(&MyArray{}, td.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MyStruct"),
			Expected: mustBe("*td_test.MyArray"),
		})

	checkError(t, gotTypedArray, td.Array([5]int{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, gotTypedArray, td.Array(MyArray{1, 2, 3, 4, 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})
	checkError(t, gotTypedArray,
		td.Array(MyArray{1, 2, 3, 4}, td.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})

	checkError(t, &gotTypedArray, td.Array([5]int{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, &gotTypedArray, td.Array(&MyArray{1, 2, 3, 4, 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})
	checkError(t, &gotTypedArray,
		td.Array(&MyArray{1, 2, 3, 4}, td.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("5"),
			Expected: mustBe("6"),
		})

	// Be lax...
	// Without Lax → error
	checkError(t, MyArray{}, td.Array([5]int{}),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkError(t, [5]int{}, td.Array(MyArray{}),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkOK(t, MyArray{}, td.Lax(td.Array([5]int{})))
	checkOK(t, [5]int{}, td.Lax(td.Array(MyArray{})))

	//
	// Bad usage
	checkError(t, "never tested",
		td.Array("test"),
		expectedError{
			Message: mustBe("bad usage of Array operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Array(&MyStruct{}),
		expectedError{
			Message: mustBe("bad usage of Array operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES), but received *td_test.MyStruct (ptr) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Array([]int{}),
		expectedError{
			Message: mustBe("bad usage of Array operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES), but received []int (slice) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Array([1]int{}, td.ArrayEntries{1: 34}),
		expectedError{
			Message: mustBe("bad usage of Array operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("array length is 1, so cannot have #1 expected index"),
		})

	checkError(t, [3]int{},
		td.Array([3]int{}, td.ArrayEntries{1: nil}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("0"),
			Expected: mustBe("nil"),
		})

	checkError(t, [3]int{},
		td.Array([3]int{}, td.ArrayEntries{1: "bad"}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("int"),
			Expected: mustBe("string"),
		})

	checkError(t, "never tested",
		td.Array([1]int{12}, td.ArrayEntries{0: 21}),
		expectedError{
			Message: mustBe("bad usage of Array operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("non zero #0 entry in model already exists in expectedEntries"),
		})

	//
	// String
	test.EqualStr(t,
		td.Array(MyArray{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}).String(),
		`Array(td_test.MyArray{
  0: 2
  1: 3
  2: 4
  3: 0
  4: 0
})`)

	test.EqualStr(t,
		td.Array(&MyArray{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}).String(),
		`Array(*td_test.MyArray{
  0: 2
  1: 3
  2: 4
  3: 0
  4: 0
})`)

	test.EqualStr(t, td.Array([0]int{}, td.ArrayEntries{}).String(),
		`Array([0]int{})`)

	// Erroneous op
	test.EqualStr(t, td.Array(12).String(), "Array(<ERROR>)")
}

func TestArrayTypeBehind(t *testing.T) {
	type MyArray [12]int

	equalTypes(t, td.Array([12]int{}), [12]int{})
	equalTypes(t, td.Array(MyArray{}), MyArray{})
	equalTypes(t, td.Array(&MyArray{}), &MyArray{})

	// Erroneous op
	equalTypes(t, td.Array(12), nil)
}

func TestSlice(t *testing.T) {
	type MySlice []int

	//
	// Simple slice
	checkOK(t, []int{}, td.Slice([]int{}))
	checkOK(t, []int{}, td.Slice([]int{}, nil))
	checkOK(t, []int{0, 3}, td.Slice([]int{0, 3}))
	checkOK(t, []int{2, 3},
		td.Slice([]int{}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3},
		td.Slice(([]int)(nil), td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4},
		td.Slice([]int{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4},
		td.Slice([]int{2, 3}, td.ArrayEntries{2: 4}))
	checkOK(t, []int{2, 3, 4, 0, 6},
		td.Slice([]int{2, 3}, td.ArrayEntries{2: 4, 4: 6}))

	gotSlice := []int{2, 3, 4}

	checkError(t, gotSlice, td.Slice(MySlice{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("td_test.MySlice"),
		})
	checkError(t, gotSlice, td.Slice([]int{2, 3, 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})
	checkError(t, gotSlice,
		td.Slice([]int{2, 3}, td.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})

	checkError(t, nil,
		td.Slice([]int{2, 3}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustContain("Slice("),
		})

	//
	// Slice type
	checkOK(t, MySlice{}, td.Slice(MySlice{}))
	checkOK(t, MySlice{}, td.Slice(MySlice{}, nil))
	checkOK(t, MySlice{0, 3}, td.Slice(MySlice{0, 3}))
	checkOK(t, MySlice{2, 3},
		td.Slice(MySlice{}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3},
		td.Slice((MySlice)(nil), td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4},
		td.Slice(MySlice{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4, 0, 6},
		td.Slice(MySlice{2, 3}, td.ArrayEntries{2: 4, 4: 6}))

	checkOK(t, &MySlice{}, td.Slice(&MySlice{}))
	checkOK(t, &MySlice{0, 3}, td.Slice(&MySlice{0, 3}))
	checkOK(t, &MySlice{2, 3},
		td.Slice(&MySlice{}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3},
		td.Slice((*MySlice)(nil), td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4},
		td.Slice(&MySlice{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4, 0, 6},
		td.Slice(&MySlice{2, 3}, td.ArrayEntries{2: 4, 4: 6}))

	gotTypedSlice := MySlice{2, 3, 4}

	checkError(t, 123, td.Slice(&MySlice{}, td.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*td_test.MySlice"),
		})

	checkError(t, &MyStruct{},
		td.Slice(&MySlice{}, td.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MyStruct"),
			Expected: mustBe("*td_test.MySlice"),
		})

	checkError(t, gotTypedSlice, td.Slice([]int{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, gotTypedSlice, td.Slice(MySlice{2, 3, 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})
	checkError(t, gotTypedSlice,
		td.Slice(MySlice{2, 3}, td.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})
	checkError(t, gotTypedSlice,
		td.Slice(MySlice{2, 3, 4}, td.ArrayEntries{3: 5}),
		expectedError{
			Message:  mustBe("expected value out of range"),
			Path:     mustBe("DATA[3]"),
			Got:      mustBe("<non-existent value>"),
			Expected: mustBe("5"),
		})
	checkError(t, gotTypedSlice, td.Slice(MySlice{2, 3}),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("<non-existent value>"),
		})

	checkError(t, &gotTypedSlice, td.Slice([]int{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, &gotTypedSlice, td.Slice(&MySlice{2, 3, 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})
	checkError(t, &gotTypedSlice,
		td.Slice(&MySlice{2, 3}, td.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})
	checkError(t, &gotTypedSlice, td.Slice(&MySlice{2, 3}),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("4"),
			Expected: mustBe("<non-existent value>"),
		})

	//
	// nil cases
	var (
		gotNilSlice      []int
		gotNilTypedSlice MySlice
	)

	checkOK(t, gotNilSlice, td.Slice([]int{}))
	checkOK(t, gotNilTypedSlice, td.Slice(MySlice{}))
	checkOK(t, &gotNilTypedSlice, td.Slice(&MySlice{}))

	// Be lax...
	// Without Lax → error
	checkError(t, MySlice{}, td.Slice([]int{}),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkError(t, []int{}, td.Slice(MySlice{}),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkOK(t, MySlice{}, td.Lax(td.Slice([]int{})))
	checkOK(t, []int{}, td.Lax(td.Slice(MySlice{})))

	//
	// Bad usage
	checkError(t, "never tested",
		td.Slice("test", nil),
		expectedError{
			Message: mustBe("bad usage of Slice operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Slice(&MyStruct{}),
		expectedError{
			Message: mustBe("bad usage of Slice operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES), but received *td_test.MyStruct (ptr) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Slice([0]int{}),
		expectedError{
			Message: mustBe("bad usage of Slice operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES), but received [0]int (array) as 1st parameter"),
		})

	checkError(t, []int{0, 5},
		td.Slice([]int{}, td.ArrayEntries{1: "bad"}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("int"),
			Expected: mustBe("string"),
		})

	checkError(t, "never tested",
		td.Slice([]int{12}, td.ArrayEntries{0: 21}),
		expectedError{
			Message: mustBe("bad usage of Slice operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("non zero #0 entry in model already exists in expectedEntries"),
		})

	//
	// String
	test.EqualStr(t,
		td.Slice(MySlice{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(td_test.MySlice{
  0: 2
  1: 3
  2: 4
})`)

	test.EqualStr(t,
		td.Slice(&MySlice{0, 0, 4}, td.ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(*td_test.MySlice{
  0: 2
  1: 3
  2: 4
})`)

	test.EqualStr(t, td.Slice(&MySlice{}, td.ArrayEntries{}).String(),
		`Slice(*td_test.MySlice{})`)

	// Erroneous op
	test.EqualStr(t, td.Slice(12).String(), "Slice(<ERROR>)")
}

func TestSliceTypeBehind(t *testing.T) {
	type MySlice []int

	equalTypes(t, td.Slice([]int{}), []int{})
	equalTypes(t, td.Slice(MySlice{}), MySlice{})
	equalTypes(t, td.Slice(&MySlice{}), &MySlice{})

	// Erroneous op
	equalTypes(t, td.Slice(12), nil)
}

func TestSuperSliceOf(t *testing.T) {
	t.Run("interface array", func(t *testing.T) {
		got := [5]any{"foo", "bar", nil, 666, 777}

		checkOK(t, got,
			td.SuperSliceOf([5]any{1: "bar"}, td.ArrayEntries{2: td.Nil()}))
		checkOK(t, got,
			td.SuperSliceOf([5]any{1: "bar"}, td.ArrayEntries{2: nil}))
		checkOK(t, got,
			td.SuperSliceOf([5]any{1: "bar"}, td.ArrayEntries{3: 666}))
		checkOK(t, got,
			td.SuperSliceOf([5]any{1: "bar"}, td.ArrayEntries{3: td.Between(665, 667)}))
		checkOK(t, &got,
			td.SuperSliceOf(&[5]any{1: "bar"}, td.ArrayEntries{3: td.Between(665, 667)}))

		checkOK(t, got,
			td.SuperSliceOf([5]any{1: "bar"},
				td.ArrayEntries{3: 42},
				td.ArrayEntries{3: "pipo"},
				td.ArrayEntries{3: 666},
			))

		checkError(t, got,
			td.SuperSliceOf([5]any{1: "foo"}, td.ArrayEntries{2: td.Nil()}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA[1]"),
				Got:      mustBe(`"bar"`),
				Expected: mustBe(`"foo"`),
			})
		checkError(t, got,
			td.SuperSliceOf([5]any{1: 666}, td.ArrayEntries{2: td.Nil()}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA[1]"),
				Got:      mustBe("string"),
				Expected: mustBe("int"),
			})
		checkError(t, &got,
			td.SuperSliceOf([5]any{1: 666}, td.ArrayEntries{2: td.Nil()}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*[5]interface {}"),
				Expected: mustBe("[5]interface {}"),
			})
		checkError(t, got,
			td.SuperSliceOf(&[5]any{1: 666}, td.ArrayEntries{2: td.Nil()}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe("[5]interface {}"),
				Expected: mustBe("*[5]interface {}"),
			})
	})

	t.Run("ints array", func(t *testing.T) {
		type MyArray [5]int

		checkOK(t, MyArray{}, td.SuperSliceOf(MyArray{}))

		got := MyArray{3: 4}
		checkOK(t, got, td.SuperSliceOf(MyArray{}))
		checkOK(t, got, td.SuperSliceOf(MyArray{3: 4}))
		checkOK(t, got, td.SuperSliceOf(MyArray{}, td.ArrayEntries{3: 4}))

		checkError(t, got,
			td.SuperSliceOf(MyArray{}, td.ArrayEntries{1: 666}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA[1]"),
				Got:      mustBe(`0`),
				Expected: mustBe(`666`),
			})

		// Be lax...
		// Without Lax → error
		checkError(t, got,
			td.SuperSliceOf([5]int{}, td.ArrayEntries{3: 4}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe(`td_test.MyArray`),
				Expected: mustBe(`[5]int`),
			})
		checkOK(t, got, td.Lax(td.SuperSliceOf([5]int{}, td.ArrayEntries{3: 4})))
		checkError(t, [5]int{3: 4},
			td.SuperSliceOf(MyArray{}, td.ArrayEntries{3: 4}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe(`[5]int`),
				Expected: mustBe(`td_test.MyArray`),
			})
		checkOK(t, [5]int{3: 4},
			td.Lax(td.SuperSliceOf(MyArray{}, td.ArrayEntries{3: 4})))

		checkError(t, "never tested",
			td.SuperSliceOf(MyArray{}, td.ArrayEntries{8: 34}),
			expectedError{
				Message: mustBe("bad usage of SuperSliceOf operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("array length is 5, so cannot have #8 expected index"),
			})
	})

	t.Run("ints slice", func(t *testing.T) {
		type MySlice []int

		checkOK(t, MySlice{}, td.SuperSliceOf(MySlice{}))
		checkOK(t, MySlice(nil), td.SuperSliceOf(MySlice{}))

		got := MySlice{3: 4}

		checkOK(t, got, td.SuperSliceOf(MySlice{}, td.ArrayEntries{3: td.N(5, 1)}))
		checkOK(t, got, td.SuperSliceOf(MySlice{3: 4}, td.ArrayEntries{2: 0}))

		checkError(t, got,
			td.SuperSliceOf(MySlice{}, td.ArrayEntries{1: 666}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA[1]"),
				Got:      mustBe(`0`),
				Expected: mustBe(`666`),
			})
		checkError(t, got,
			td.SuperSliceOf(MySlice{}, td.ArrayEntries{3: 0}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA[3]"),
				Got:      mustBe(`4`),
				Expected: mustBe(`0`),
			})
		checkError(t, got,
			td.SuperSliceOf(MySlice{}, td.ArrayEntries{28: 666}),
			expectedError{
				Message:  mustBe("expected value out of range"),
				Path:     mustBe("DATA[28]"),
				Got:      mustBe(`<non-existent value>`),
				Expected: mustBe(`666`),
			})
		checkError(t, got,
			td.SuperSliceOf(MySlice{28: 666}),
			expectedError{
				Message:  mustBe("expected value out of range"),
				Path:     mustBe("DATA[28]"),
				Got:      mustBe(`<non-existent value>`),
				Expected: mustBe(`666`),
			})

		// Be lax...
		// Without Lax → error
		checkError(t, got,
			td.SuperSliceOf([]int{}, td.ArrayEntries{3: 4}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe(`td_test.MySlice`),
				Expected: mustBe(`[]int`),
			})
		checkOK(t, got, td.Lax(td.SuperSliceOf([]int{}, td.ArrayEntries{3: 4})))
		checkError(t, []int{3: 4},
			td.SuperSliceOf(MySlice{}, td.ArrayEntries{3: 4}),
			expectedError{
				Message:  mustBe("type mismatch"),
				Path:     mustBe("DATA"),
				Got:      mustBe(`[]int`),
				Expected: mustBe(`td_test.MySlice`),
			})
		checkOK(t, []int{3: 4},
			td.Lax(td.SuperSliceOf(MySlice{}, td.ArrayEntries{3: 4})))
	})

	//
	// Bad usage
	checkError(t, "never tested",
		td.SuperSliceOf("test", nil),
		expectedError{
			Message: mustBe("bad usage of SuperSliceOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: SuperSliceOf(ARRAY|&ARRAY|SLICE|&SLICE, EXPECTED_ENTRIES), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.SuperSliceOf(&MyStruct{}),
		expectedError{
			Message: mustBe("bad usage of SuperSliceOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: SuperSliceOf(ARRAY|&ARRAY|SLICE|&SLICE, EXPECTED_ENTRIES), but received *td_test.MyStruct (ptr) as 1st parameter"),
		})

	checkError(t, []int{0, 1},
		td.SuperSliceOf([]int{}, td.ArrayEntries{1: "bad"}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("int"),
			Expected: mustBe("string"),
		})

	checkError(t, "never tested",
		td.SuperSliceOf([]int{12}, td.ArrayEntries{0: 21}),
		expectedError{
			Message: mustBe("bad usage of SuperSliceOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("non zero #0 entry in model already exists in expectedEntries"),
		})

	// Erroneous op
	test.EqualStr(t, td.SuperSliceOf(12).String(), "SuperSliceOf(<ERROR>)")
}

func TestSuperSliceOfTypeBehind(t *testing.T) {
	type MySlice []int

	equalTypes(t, td.SuperSliceOf([]int{}), []int{})
	equalTypes(t, td.SuperSliceOf(MySlice{}), MySlice{})
	equalTypes(t, td.SuperSliceOf(&MySlice{}), &MySlice{})

	type MyArray [12]int

	equalTypes(t, td.SuperSliceOf([12]int{}), [12]int{})
	equalTypes(t, td.SuperSliceOf(MyArray{}), MyArray{})
	equalTypes(t, td.SuperSliceOf(&MyArray{}), &MyArray{})

	// Erroneous op
	equalTypes(t, td.SuperSliceOf(12), nil)
}
