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
		td.ArrayEach(4))

	// Empty slice/array
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
			// nil cases
			([]int)(nil),
			MySlice(nil),
		},
		td.ArrayEach(4))

	checkError(t, (*MyArray)(nil), td.ArrayEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *td_test.MyArray"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})
	checkError(t, (*MySlice)(nil), td.ArrayEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *td_test.MySlice"),
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
		td.ArrayEach(td.Between(20, 30)))

	checkError(t, nil, td.ArrayEach(4),
		expectedError{
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
		td.ArrayEach(4),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, 666, td.ArrayEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})
	num := 666
	checkError(t, &num, td.ArrayEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
		})

	checkOK(t, []interface{}{nil, nil, nil}, td.ArrayEach(nil))
	checkError(t, []interface{}{nil, nil, nil, 66}, td.ArrayEach(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[3]"),
			Got:      mustBe("66"),
			Expected: mustBe("nil"),
		})

	//
	// String
	test.EqualStr(t, td.ArrayEach(4).String(), "ArrayEach(4)")
	test.EqualStr(t, td.ArrayEach(td.All(1, 2)).String(),
		`ArrayEach(All(1,
              2))`)
}

func TestArrayEachTypeBehind(t *testing.T) {
	equalTypes(t, td.ArrayEach(6), nil)
}
