package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestArray(t *testing.T) {
	type MyArray [5]int

	//
	// Simple array
	gotArray := [...]int{1, 2, 3, 4, 5}

	checkOK(t, gotArray, Array([5]int{}, nil))
	checkOK(t, gotArray, Array([5]int{0, 0, 0, 4}, nil))
	checkOK(t, gotArray, Array([5]int{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, gotArray, Array([5]int{0, 2}, ArrayEntries{2: 3, 0: 1}))

	checkError(t, gotArray, Array(MyArray{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[5]int"),
			Expected: mustBe("testdeep_test.MyArray"),
		})
	checkError(t, gotArray, Array([5]int{0, 0, 0, 0, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotArray, Array([5]int{}, ArrayEntries{4: 6}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})

	//
	// Array type
	gotTypedArray := MyArray{1, 2, 3, 4, 5}

	checkOK(t, gotTypedArray, Array(MyArray{}, nil))
	checkOK(t, gotTypedArray, Array(MyArray{0, 0, 0, 4}, nil))
	checkOK(t, gotTypedArray, Array(MyArray{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, gotTypedArray, Array(MyArray{0, 2}, ArrayEntries{2: 3, 0: 1}))

	checkOK(t, &gotTypedArray, Array(&MyArray{}, nil))
	checkOK(t, &gotTypedArray, Array(&MyArray{0, 0, 0, 4}, nil))
	checkOK(t, &gotTypedArray, Array(&MyArray{}, ArrayEntries{2: 3, 0: 1}))
	checkOK(t, &gotTypedArray, Array(&MyArray{0, 2}, ArrayEntries{2: 3, 0: 1}))

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
	checkError(t, gotTypedArray, Array(MyArray{0, 0, 0, 0, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, gotTypedArray, Array(MyArray{}, ArrayEntries{4: 6}),
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
	checkError(t, &gotTypedArray, Array(&MyArray{0, 0, 0, 0, 6}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[4]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 6"),
		})
	checkError(t, &gotTypedArray, Array(&MyArray{}, ArrayEntries{4: 6}),
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
	checkPanic(t, func() { Array([3]int{}, ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from array content (int)")
	checkPanic(t, func() { Array([1]int{12}, ArrayEntries{0: 21}) },
		"non zero #0 entry in model already exists in expectedEntries")

	//
	// String
	equalStr(t, Array(MyArray{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Array(testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	equalStr(t, Array(&MyArray{0, 0, 4}, ArrayEntries{1: 3, 0: 2}).String(),
		`Array(*testdeep_test.MyArray{
  0: (int) 2
  1: (int) 3
  2: (int) 4
})`)

	equalStr(t, Array(&MyArray{}, ArrayEntries{}).String(),
		`Array(*testdeep_test.MyArray{})`)
}

func TestSlice(t *testing.T) {
	type MySlice []int

	//
	// Simple slice
	gotSlice := []int{2, 3, 4}

	checkOK(t, gotSlice, Slice([]int{}, nil))
	checkOK(t, gotSlice, Slice([]int{0, 3}, nil))
	checkOK(t, gotSlice, Slice([]int{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, gotSlice, Slice([]int{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))

	checkError(t, gotSlice, Slice(MySlice{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("[]int"),
			Expected: mustBe("testdeep_test.MySlice"),
		})
	checkError(t, gotSlice, Slice([]int{0, 0, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotSlice, Slice([]int{}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})

	//
	// Slice type
	gotTypedSlice := MySlice{2, 3, 4}

	checkOK(t, gotTypedSlice, Slice(MySlice{}, nil))
	checkOK(t, gotTypedSlice, Slice(MySlice{0, 3}, nil))
	checkOK(t, gotTypedSlice, Slice(MySlice{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, gotTypedSlice, Slice(MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))

	checkOK(t, &gotTypedSlice, Slice(&MySlice{}, nil))
	checkOK(t, &gotTypedSlice, Slice(&MySlice{0, 3}, nil))
	checkOK(t, &gotTypedSlice, Slice(&MySlice{}, ArrayEntries{1: 3, 0: 2}))
	checkOK(t, &gotTypedSlice, Slice(&MySlice{0, 0, 4}, ArrayEntries{1: 3, 0: 2}))

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
	checkError(t, gotTypedSlice, Slice(MySlice{0, 0, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, gotTypedSlice, Slice(MySlice{}, ArrayEntries{65: 5}),
		expectedError{
			Message:  mustBe("expected value out of range"),
			Path:     mustBe("DATA[65]"),
			Got:      mustBe("<non-existent value>"),
			Expected: mustBe("(int) 5"),
		})

	checkError(t, &gotTypedSlice, Slice([]int{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MySlice"),
			Expected: mustBe("[]int"),
		})
	checkError(t, &gotTypedSlice, Slice(&MySlice{0, 0, 5}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})
	checkError(t, &gotTypedSlice, Slice(&MySlice{}, ArrayEntries{2: 5}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[2]"),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})

	//
	// Bad usage
	checkPanic(t, func() { Slice("test", nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice(&MyStruct{}, nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice([0]int{}, nil) }, "usage: Slice(")
	checkPanic(t, func() { Slice([]int{}, ArrayEntries{1: "bad"}) },
		"type string of #1 expected value differs from slice content (int)")
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
