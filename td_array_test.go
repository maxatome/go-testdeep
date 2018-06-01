// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestArray(t *testing.T) {
	type MyArray [5]int

	//
	// Simple array
	checkOK(t, [5]int{}, Array([5]int{}, nil))
	checkOK(t, [5]int{0, 0, 0, 4}, Array([5]int{0, 0, 0, 4}, nil))
	checkOK(t, [5]int{1, 0, 3}, Array([5]int{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, [5]int{1, 2, 3}, Array([5]int{0, 2}, ArrayEntries{2: 3, 0: 1}))

	zero, one, two := 0, 1, 2
	checkOK(t, [5]*int{nil, &zero, &one, &two},
		Array([5]*int{}, ArrayEntries{1: &zero, 2: &one, 3: &two, 4: nil}))

	gotArray := [...]int{1, 2, 3, 4, 5}

	checkError(t, gotArray, Array(MyArray{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[5]int"),
			Expected: mustBe("testdeep_test.MyArray"),
		})
	checkError(t, gotArray, Array([5]int{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotArray, Array([5]int{1, 2, 3, 4}, ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	//
	// Array type
	checkOK(t, MyArray{}, Array(MyArray{}, nil))
	checkOK(t, MyArray{0, 0, 0, 4}, Array(MyArray{0, 0, 0, 4}, nil))
	checkOK(t, MyArray{1, 0, 3}, Array(MyArray{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, MyArray{1, 2, 3}, Array(MyArray{0, 2}, ArrayEntries{2: 3, 0: 1}))

	checkOK(t, &MyArray{}, Array(&MyArray{}, nil))
	checkOK(t, &MyArray{0, 0, 0, 4}, Array(&MyArray{0, 0, 0, 4}, nil))
	checkOK(t, &MyArray{1, 0, 3}, Array(&MyArray{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 0, 3}, Array((*MyArray)(nil), ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &MyArray{1, 2, 3}, Array(&MyArray{0, 2}, ArrayEntries{2: 3, 0: 1}))

	gotTypedArray := MyArray{1, 2, 3, 4, 5}

	checkError(t, 123, Array(&MyArray{}, ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*testdeep_test.MyArray"),
		})

	checkError(t, &MyStruct{}, Array(&MyArray{}, ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyStruct"),
			Expected: mustBe("*testdeep_test.MyArray"),
		})

	checkError(t, gotTypedArray, Array([5]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("testdeep_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, gotTypedArray, Array(MyArray{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotTypedArray, Array(MyArray{1, 2, 3, 4}, ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	checkError(t, &gotTypedArray, Array([5]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyArray"),
			Expected: mustBe("[5]int"),
		})
	checkError(t, &gotTypedArray, Array(&MyArray{1, 2, 3, 4, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, &gotTypedArray, Array(&MyArray{1, 2, 3, 4}, ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	//
	// Bad usage
	checkPanic(t, func() { Array("test", nil) }, "usage: Array(")
	checkPanic(t, func() { Array(&MyStruct{}, nil) }, "usage: Array(")
	checkPanic(t, func() { Array([]int{}, nil) }, "usage: Array(")
	checkPanic(t, func() { Array([1]int{}, ArrayEntries{1: 34}) },
		"array length is 1, so cannot have #1 expected index")
	checkPanic(t, func() { Array([3]int{}, ArrayEntries{1: nil}) },
		"expected value of #1 cannot be nil as items type is int")
	checkPanic(t, func() { Array([3]int{}, ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from array contents (int)")
	checkPanic(t, func() { Array([1]int{12}, ArrayEntries{0: 21}) },
		"non zero #0 entry in model already exists in expectedEntries")

	//
	// String
	equalStr(t, Array(MyArray{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Array(testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
  3: (int) 0
  4: (int) 0
})`)

	equalStr(t, Array(&MyArray{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Array(*testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
  3: (int) 0
  4: (int) 0
})`)

	equalStr(t, Array([0]int{}, ArrayEntries{}).String(),
		`Array([0]int{})`)
}

func TestArrayTypeBehind(t *testing.T) {
	type MyArray [12]int

	equalTypes(t, Array([12]int{}, nil), [12]int{})
	equalTypes(t, Array(MyArray{}, nil), MyArray{})
	equalTypes(t, Array(&MyArray{}, nil), &MyArray{})
}

func TestSlice(t *testing.T) {
	type MySlice []int

	//
	// Simple slice
	checkOK(t, []int{}, Slice([]int{}, nil))
	checkOK(t, []int{0, 3}, Slice([]int{0, 3}, nil))
	checkOK(t, []int{2, 3}, Slice([]int{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3}, Slice(([]int)(nil), ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4}, Slice([]int{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, []int{2, 3, 4}, Slice([]int{2, 3}, ArrayEntries{2: 4}))
	checkOK(t, []int{2, 3, 4, 0, 6}, Slice([]int{2, 3}, ArrayEntries{2: 4, 4: 6}))

	gotSlice := []int{2, 3, 4}

	checkError(t, gotSlice, Slice(MySlice{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("testdeep_test.MySlice"),
		})
	checkError(t, gotSlice, Slice([]int{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotSlice, Slice([]int{2, 3}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})

	//
	// Slice type
	checkOK(t, MySlice{}, Slice(MySlice{}, nil))
	checkOK(t, MySlice{0, 3}, Slice(MySlice{0, 3}, nil))
	checkOK(t, MySlice{2, 3}, Slice(MySlice{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3}, Slice((MySlice)(nil), ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4},
		Slice(MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, MySlice{2, 3, 4, 0, 6},
		Slice(MySlice{2, 3}, ArrayEntries{2: 4, 4: 6}))

	checkOK(t, &MySlice{}, Slice(&MySlice{}, nil))
	checkOK(t, &MySlice{0, 3}, Slice(&MySlice{0, 3}, nil))
	checkOK(t, &MySlice{2, 3}, Slice(&MySlice{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3}, Slice((*MySlice)(nil), ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4},
		Slice(&MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &MySlice{2, 3, 4, 0, 6},
		Slice(&MySlice{2, 3}, ArrayEntries{2: 4, 4: 6}))

	gotTypedSlice := MySlice{2, 3, 4}

	checkError(t, 123, Slice(&MySlice{}, ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("*testdeep_test.MySlice"),
		})

	checkError(t, &MyStruct{}, Slice(&MySlice{}, ArrayEntries{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyStruct"),
			Expected: mustBe("*testdeep_test.MySlice"),
		})

	checkError(t, gotTypedSlice, Slice([]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("testdeep_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{2, 3}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{2, 3, 4}, ArrayEntries{3: 5}),
		expectedError{
			Message:  mustBe("expected value out of range"),
			Path:     mustBe("DATA[3]"),
			Got:      mustBe("<non-existent value>"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{2, 3}, nil),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("<non-existent value>"),
		})

	checkError(t, &gotTypedSlice, Slice([]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, &gotTypedSlice, Slice(&MySlice{2, 3, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, &gotTypedSlice, Slice(&MySlice{2, 3}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, &gotTypedSlice, Slice(&MySlice{2, 3}, nil),
		expectedError{
			Message:  mustBe("got value out of range"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("<non-existent value>"),
		})

	//
	// Bad usage
	checkPanic(t, func() { Slice("test", nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice(&MyStruct{}, nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice([0]int{}, nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice([]int{}, ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from slice contents (int)")
	checkPanic(t, func() { Slice([]int{12}, ArrayEntries{0: 21}) },
		"non zero #0 entry in model already exists in expectedEntries")

	//
	// String
	equalStr(t, Slice(MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(testdeep_test.MySlice{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	equalStr(t, Slice(&MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Slice(*testdeep_test.MySlice{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	equalStr(t, Slice(&MySlice{}, ArrayEntries{}).String(),
		`Slice(*testdeep_test.MySlice{})`)
}

func TestSliceTypeBehind(t *testing.T) {
	type MySlice []int

	equalTypes(t, Slice([]int{}, nil), []int{})
	equalTypes(t, Slice(MySlice{}, nil), MySlice{})
	equalTypes(t, Slice(&MySlice{}, nil), &MySlice{})
}
