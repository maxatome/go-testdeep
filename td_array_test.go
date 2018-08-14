// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestArray(t *testing.T) {
	type MyArray [5]int

	//
	// Simple array
	checkOK(t, [5]int{}, testdeep.Array([5]int{}, nil))
	checkOK(t, [5]int{0, 0, 0, 4}, testdeep.Array([5]int{0, 0, 0, 4}, nil))
	checkOK(t, [5]int{1, 0, 3},
		testdeep.Array([5]int{}, testdeep.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, [5]int{1, 2, 3},
		testdeep.Array([5]int{0, 2}, testdeep.ArrayEntries{2: 3, 0: 1}))

	zero, one, two := 0, 1, 2
	checkOK(t, [5]*int{nil, &zero, &one, &two},
		testdeep.Array(
			[5]*int{}, testdeep.ArrayEntries{1: &zero, 2: &one, 3: &two, 4: nil}))

	gotArray := [...]int{1, 2, 3, 4, 5}

	checkError(t, gotArray, testdeep.Array(MyArray{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[5]int"),
			Expected: mustBe("testdeep_test.MyArray"),
		})
	checkError(t, gotArray, testdeep.Array([5]int{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotArray,
		testdeep.Array([5]int{1, 2, 3, 4}, testdeep.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	//
	// Array type
	checkOK(t, MyArray{}, testdeep.Array(MyArray{}, nil))
	checkOK(t, MyArray{0, 0, 0, 4}, testdeep.Array(MyArray{0, 0, 0, 4}, nil))
	checkOK(t, MyArray{1, 0, 3},
		testdeep.Array(MyArray{}, testdeep.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, MyArray{1, 2, 3},
		testdeep.Array(MyArray{0, 2}, testdeep.ArrayEntries{2: 3, 0: 1}))

	checkOK(t, &MyArray{}, testdeep.Array(&MyArray{}, nil))
	checkOK(t, &MyArray{0, 0, 0, 4}, testdeep.Array(&MyArray{0, 0, 0, 4}, nil))
	checkOK(t, &MyArray{1, 0, 3},
		testdeep.Array(&MyArray{}, testdeep.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 0, 3},
		testdeep.Array((*MyArray)(nil), testdeep.ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 2, 3},
		testdeep.Array(&MyArray{0, 2}, testdeep.ArrayEntries{2: 3, 0: 1}))

	gotTypedArray := MyArray{1, 2, 3, 4, 5}

	checkError(t, 123, testdeep.Array(&MyArray{}, testdeep.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*testdeep_test.MyArray"),
		})

	checkError(t, &MyStruct{},
		testdeep.Array(&MyArray{}, testdeep.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyStruct"),
			Expected: mustBe("*testdeep_test.MyArray"),
		})

	checkError(t, gotTypedArray, testdeep.Array([5]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("testdeep_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, gotTypedArray, testdeep.Array(MyArray{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotTypedArray,
		testdeep.Array(MyArray{1, 2, 3, 4}, testdeep.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	checkError(t, &gotTypedArray, testdeep.Array([5]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, &gotTypedArray, testdeep.Array(&MyArray{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, &gotTypedArray,
		testdeep.Array(&MyArray{1, 2, 3, 4}, testdeep.ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Array("test", nil) }, "usage: Array(")
	test.CheckPanic(t,
		func() { testdeep.Array(&MyStruct{}, nil) },
		"usage: Array(")
	test.CheckPanic(t, func() { testdeep.Array([]int{}, nil) }, "usage: Array(")
	test.CheckPanic(t,
		func() { testdeep.Array([1]int{}, testdeep.ArrayEntries{1: 34}) },
		"array length is 1, so cannot have #1 expected index")
	test.CheckPanic(t,
		func() { testdeep.Array([3]int{}, testdeep.ArrayEntries{1: nil}) },
		"expected value of #1 cannot be nil as items type is int")
	test.CheckPanic(t,
		func() { testdeep.Array([3]int{}, testdeep.ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from array contents (int)")
	test.CheckPanic(t,
		func() { testdeep.Array([1]int{12}, testdeep.ArrayEntries{0: 21}) },
		"non zero #0 entry in model already exists in expectedEntries")

	//
	// String
	test.EqualStr(t,
		testdeep.Array(MyArray{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}).String(),
		`Array(testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
  3: (int) 0
  4: (int) 0
})`)

	test.EqualStr(t,
		testdeep.Array(&MyArray{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}).String(),
		`Array(*testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
  3: (int) 0
  4: (int) 0
})`)

	test.EqualStr(t, testdeep.Array([0]int{}, testdeep.ArrayEntries{}).String(),
		`Array([0]int{})`)
}

func TestArrayTypeBehind(t *testing.T) {
	type MyArray [12]int

	equalTypes(t, testdeep.Array([12]int{}, nil), [12]int{})
	equalTypes(t, testdeep.Array(MyArray{}, nil), MyArray{})
	equalTypes(t, testdeep.Array(&MyArray{}, nil), &MyArray{})
}

func TestSlice(t *testing.T) {
	type MySlice []int

	//
	// Simple slice
	checkOK(t, []int{}, testdeep.Slice([]int{}, nil))
	checkOK(t, []int{0, 3}, testdeep.Slice([]int{0, 3}, nil))
	checkOK(t, []int{2, 3},
		testdeep.Slice([]int{}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3},
		testdeep.Slice(([]int)(nil), testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4},
		testdeep.Slice([]int{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4},
		testdeep.Slice([]int{2, 3}, testdeep.ArrayEntries{2: 4}))
	checkOK(t, []int{2, 3, 4, 0, 6},
		testdeep.Slice([]int{2, 3}, testdeep.ArrayEntries{2: 4, 4: 6}))

	gotSlice := []int{2, 3, 4}

	checkError(t, gotSlice, testdeep.Slice(MySlice{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("testdeep_test.MySlice"),
		})
	checkError(t, gotSlice, testdeep.Slice([]int{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotSlice,
		testdeep.Slice([]int{2, 3}, testdeep.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})

	//
	// Slice type
	checkOK(t, MySlice{}, testdeep.Slice(MySlice{}, nil))
	checkOK(t, MySlice{0, 3}, testdeep.Slice(MySlice{0, 3}, nil))
	checkOK(t, MySlice{2, 3},
		testdeep.Slice(MySlice{}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3},
		testdeep.Slice((MySlice)(nil), testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4},
		testdeep.Slice(MySlice{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4, 0, 6},
		testdeep.Slice(MySlice{2, 3}, testdeep.ArrayEntries{2: 4, 4: 6}))

	checkOK(t, &MySlice{}, testdeep.Slice(&MySlice{}, nil))
	checkOK(t, &MySlice{0, 3}, testdeep.Slice(&MySlice{0, 3}, nil))
	checkOK(t, &MySlice{2, 3},
		testdeep.Slice(&MySlice{}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3},
		testdeep.Slice((*MySlice)(nil), testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4},
		testdeep.Slice(&MySlice{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4, 0, 6},
		testdeep.Slice(&MySlice{2, 3}, testdeep.ArrayEntries{2: 4, 4: 6}))

	gotTypedSlice := MySlice{2, 3, 4}

	checkError(t, 123, testdeep.Slice(&MySlice{}, testdeep.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*testdeep_test.MySlice"),
		})

	checkError(t, &MyStruct{},
		testdeep.Slice(&MySlice{}, testdeep.ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyStruct"),
			Expected: mustBe("*testdeep_test.MySlice"),
		})

	checkError(t, gotTypedSlice, testdeep.Slice([]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("testdeep_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, gotTypedSlice, testdeep.Slice(MySlice{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice,
		testdeep.Slice(MySlice{2, 3}, testdeep.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice,
		testdeep.Slice(MySlice{2, 3, 4}, testdeep.ArrayEntries{3: 5}),
		expectedError{
			Message:  mustBe("expected value out of range"),
			Path:     mustBe("DATA[3]"),
			Got:      mustBe("<non-existent value>"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, testdeep.Slice(MySlice{2, 3}, nil),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("<non-existent value>"),
		})

	checkError(t, &gotTypedSlice, testdeep.Slice([]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, &gotTypedSlice, testdeep.Slice(&MySlice{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, &gotTypedSlice,
		testdeep.Slice(&MySlice{2, 3}, testdeep.ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, &gotTypedSlice, testdeep.Slice(&MySlice{2, 3}, nil),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("<non-existent value>"),
		})

	//
	// nil cases
	var (
		gotNilSlice      []int
		gotNilTypedSlice MySlice
	)

	checkOK(t, gotNilSlice, testdeep.Slice([]int{}, nil))
	checkOK(t, gotNilTypedSlice, testdeep.Slice(MySlice{}, nil))
	checkOK(t, &gotNilTypedSlice, testdeep.Slice(&MySlice{}, nil))

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Slice("test", nil) }, "usage: Slice(")
	test.CheckPanic(t,
		func() { testdeep.Slice(&MyStruct{}, nil) },
		"usage: Slice(")
	test.CheckPanic(t, func() { testdeep.Slice([0]int{}, nil) }, "usage: Slice(")
	test.CheckPanic(t,
		func() { testdeep.Slice([]int{}, testdeep.ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from slice contents (int)")
	test.CheckPanic(t,
		func() { testdeep.Slice([]int{12}, testdeep.ArrayEntries{0: 21}) },
		"non zero #0 entry in model already exists in expectedEntries")

	//
	// String
	test.EqualStr(t,
		testdeep.Slice(MySlice{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(testdeep_test.MySlice{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	test.EqualStr(t,
		testdeep.Slice(&MySlice{0, 0, 4}, testdeep.ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(*testdeep_test.MySlice{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	test.EqualStr(t, testdeep.Slice(&MySlice{}, testdeep.ArrayEntries{}).String(),
		`Slice(*testdeep_test.MySlice{})`)
}

func TestSliceTypeBehind(t *testing.T) {
	type MySlice []int

	equalTypes(t, testdeep.Slice([]int{}, nil), []int{})
	equalTypes(t, testdeep.Slice(MySlice{}, nil), MySlice{})
	equalTypes(t, testdeep.Slice(&MySlice{}, nil), &MySlice{})
}
