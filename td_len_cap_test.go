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

func TestLen(t *testing.T) {
	checkOK(t, "abcd", testdeep.Len(4))
	checkOK(t, "abcd", testdeep.Len(testdeep.Between(4, 6)))
	checkOK(t, "abcd", testdeep.Len(testdeep.Between(6, 4)))

	checkOK(t, []byte("abcd"), testdeep.Len(4))
	checkOK(t, []byte("abcd"), testdeep.Len(testdeep.Between(4, 6)))

	checkOK(t, [5]int{}, testdeep.Len(5))
	checkOK(t, [5]int{}, testdeep.Len(testdeep.Between(4, 6)))

	checkOK(t, map[int]bool{1: true, 2: false}, testdeep.Len(2))
	checkOK(t, map[int]bool{1: true, 2: false},
		testdeep.Len(testdeep.Between(1, 6)))

	checkOK(t, make(chan int, 3), testdeep.Len(0))

	checkError(t, [5]int{}, testdeep.Len(4),
		expectedError{
			Message:  mustBe("bad length"),
			Path:     mustBe("DATA"),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, [5]int{}, testdeep.Len(testdeep.Lt(4)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("len(DATA)"),
			Got:      mustBe("5"),
			Expected: mustBe("< 4"),
		})

	checkError(t, 123, testdeep.Len(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Array, Chan, Map, Slice or string"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Len(3).String(), "len=3")
	test.EqualStr(t,
		testdeep.Len(testdeep.Between(3, 8)).String(), "len: 3 ≤ got ≤ 8")
	test.EqualStr(t, testdeep.Len(testdeep.Gt(8)).String(), "len: > 8")

	//
	// Bad usage
	checkPanic(t, func() { testdeep.Len(int64(12)) }, "usage: Len(")
}

func TestCap(t *testing.T) {
	checkOK(t, make([]byte, 0, 4), testdeep.Cap(4))
	checkOK(t, make([]byte, 0, 4), testdeep.Cap(testdeep.Between(4, 6)))

	checkOK(t, [5]int{}, testdeep.Cap(5))
	checkOK(t, [5]int{}, testdeep.Cap(testdeep.Between(4, 6)))

	checkOK(t, make(chan int, 3), testdeep.Cap(3))

	checkError(t, [5]int{}, testdeep.Cap(4),
		expectedError{
			Message:  mustBe("bad capacity"),
			Path:     mustBe("DATA"),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, [5]int{}, testdeep.Cap(testdeep.Between(2, 4)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("cap(DATA)"),
			Got:      mustBe("5"),
			Expected: mustBe("2 ≤ got ≤ 4"),
		})

	checkError(t, map[int]int{1: 2}, testdeep.Cap(1),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[int]int"),
			Expected: mustBe("Array, Chan or Slice"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Cap(3).String(), "cap=3")
	test.EqualStr(t,
		testdeep.Cap(testdeep.Between(3, 8)).String(), "cap: 3 ≤ got ≤ 8")
	test.EqualStr(t, testdeep.Cap(testdeep.Gt(8)).String(), "cap: > 8")

	//
	// Bad usage
	checkPanic(t, func() { testdeep.Cap(int64(12)) }, "usage: Cap(")
}

func TestLenCapTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Cap(3), nil)
	equalTypes(t, testdeep.Len(3), nil)
}
