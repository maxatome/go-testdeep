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

func TestArrayEach(t *testing.T) {
	type MyArray [3]int
	type MyEmptyArray [0]int
	type MySlice []int

	checkOKForEach(t,
		[]interface{}{
			[...]int{4, 4, 4},
			[]int{4, 4, 4},
			&[...]int{4, 4, 4},
			&[]int{4, 4, 4},
			MyArray{4, 4, 4},
			MySlice{4, 4, 4},
			&MyArray{4, 4, 4},
			&MySlice{4, 4, 4},
		},
		ArrayEach(4))

	checkOKForEach(t,
		[]interface{}{
			[0]int{},
			[]int{},
			&[0]int{},
			&[]int{},
			MyEmptyArray{},
			MySlice{},
			&MyEmptyArray{},
			&MySlice{},
		},
		ArrayEach(4))

	checkOK(t, ([]int)(nil), ArrayEach(4))
	checkOK(t, MySlice(nil), ArrayEach(4))
	checkError(t, (*MyArray)(nil), ArrayEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *testdeep_test.MyArray"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})
	checkError(t, (*MySlice)(nil), ArrayEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *testdeep_test.MySlice"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})

	checkOKForEach(t,
		[]interface{}{
			[...]int{20, 22, 29},
			[]int{20, 22, 29},
			MyArray{20, 22, 29},
			MySlice{20, 22, 29},
			&MyArray{20, 22, 29},
			&MySlice{20, 22, 29},
		},
		ArrayEach(Between(20, 30)))

	checkError(t, nil, ArrayEach(4), expectedError{
		Message:  mustBe("nil value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
	})

	checkErrorForEach(t,
		[]interface{}{
			[...]int{4, 5, 4},
			[]int{4, 5, 4},
			MyArray{4, 5, 4},
			MySlice{4, 5, 4},
			&MyArray{4, 5, 4},
			&MySlice{4, 5, 4},
		},
		ArrayEach(4),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 4"),
		})

	checkError(t, 666, ArrayEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})
	num := 666
	checkError(t, &num, ArrayEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})

	checkOK(t, []interface{}{nil, nil, nil}, ArrayEach(nil))
	checkError(t, []interface{}{nil, nil, nil, 66}, ArrayEach(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[3]"),
			Got:      mustBe("(int) 66"),
			Expected: mustBe("nil"),
		})

	//
	// String
	equalStr(t, ArrayEach(4).String(), "ArrayEach((int) 4)")
	equalStr(t, ArrayEach(All(1, 2)).String(),
		`ArrayEach(All((int) 1,
              (int) 2))`)
}

func TestArrayEachTypeBehind(t *testing.T) {
	equalTypes(t, ArrayEach(6), nil)
}
