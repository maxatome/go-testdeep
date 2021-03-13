// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestLen(t *testing.T) {
	checkOK(t, "abcd", td.Len(4))
	checkOK(t, "abcd", td.Len(td.Between(4, 6)))
	checkOK(t, "abcd", td.Len(td.Between(6, 4)))

	checkOK(t, []byte("abcd"), td.Len(4))
	checkOK(t, []byte("abcd"), td.Len(td.Between(4, 6)))

	checkOK(t, [5]int{}, td.Len(5))
	checkOK(t, [5]int{}, td.Len(td.Between(4, 6)))

	checkOK(t, map[int]bool{1: true, 2: false}, td.Len(2))
	checkOK(t, map[int]bool{1: true, 2: false},
		td.Len(td.Between(1, 6)))

	checkOK(t, make(chan int, 3), td.Len(0))

	checkError(t, [5]int{}, td.Len(4),
		expectedError{
			Message:  mustBe("bad length"),
			Path:     mustBe("DATA"),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, [5]int{}, td.Len(td.Lt(4)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("len(DATA)"),
			Got:      mustBe("5"),
			Expected: mustBe("< 4"),
		})

	checkError(t, 123, td.Len(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Array, Chan, Map, Slice or string"),
		})

	//
	// String
	test.EqualStr(t, td.Len(3).String(), "len=3")
	test.EqualStr(t,
		td.Len(td.Between(3, 8)).String(), "len: 3 ≤ got ≤ 8")
	test.EqualStr(t, td.Len(td.Gt(8)).String(), "len: > 8")

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Len(int64(12)) }, "usage: Len(")
}

func TestCap(t *testing.T) {
	checkOK(t, make([]byte, 0, 4), td.Cap(4))
	checkOK(t, make([]byte, 0, 4), td.Cap(td.Between(4, 6)))

	checkOK(t, [5]int{}, td.Cap(5))
	checkOK(t, [5]int{}, td.Cap(td.Between(4, 6)))

	checkOK(t, make(chan int, 3), td.Cap(3))

	checkError(t, [5]int{}, td.Cap(4),
		expectedError{
			Message:  mustBe("bad capacity"),
			Path:     mustBe("DATA"),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, [5]int{}, td.Cap(td.Between(2, 4)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("cap(DATA)"),
			Got:      mustBe("5"),
			Expected: mustBe("2 ≤ got ≤ 4"),
		})

	checkError(t, map[int]int{1: 2}, td.Cap(1),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[int]int"),
			Expected: mustBe("Array, Chan or Slice"),
		})

	//
	// String
	test.EqualStr(t, td.Cap(3).String(), "cap=3")
	test.EqualStr(t,
		td.Cap(td.Between(3, 8)).String(), "cap: 3 ≤ got ≤ 8")
	test.EqualStr(t, td.Cap(td.Gt(8)).String(), "cap: > 8")

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Cap(int64(12)) }, "usage: Cap(")
}

func TestLenCapTypeBehind(t *testing.T) {
	equalTypes(t, td.Cap(3), nil)
	equalTypes(t, td.Len(3), nil)
}
