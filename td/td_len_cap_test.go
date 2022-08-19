// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"math"
	"testing"

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
	checkOK(t, [5]int{}, td.Len(int8(5)))
	checkOK(t, [5]int{}, td.Len(int16(5)))
	checkOK(t, [5]int{}, td.Len(int32(5)))
	checkOK(t, [5]int{}, td.Len(int64(5)))
	checkOK(t, [5]int{}, td.Len(uint(5)))
	checkOK(t, [5]int{}, td.Len(uint8(5)))
	checkOK(t, [5]int{}, td.Len(uint16(5)))
	checkOK(t, [5]int{}, td.Len(uint32(5)))
	checkOK(t, [5]int{}, td.Len(uint64(5)))
	checkOK(t, [5]int{}, td.Len(float32(5)))
	checkOK(t, [5]int{}, td.Len(float64(5)))
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
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("array OR chan OR map OR slice OR string"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Len(nil),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received nil as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Len("12"),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received string as 1st parameter"),
		})

	// out of bounds
	checkError(t, "never tested",
		td.Len(uint64(math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (18446744073709551615), should be in int range"),
		})

	checkError(t, "never tested",
		td.Len(float64(math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (1.8446744073709552e+19), should be in int range"),
		})

	checkError(t, "never tested",
		td.Len(float64(-math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (-1.8446744073709552e+19), should be in int range"),
		})

	checkError(t, "never tested",
		td.Len(3.1),
		expectedError{
			Message: mustBe("bad usage of Len operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Len(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (3.1), should be in int range"),
		})

	//
	// String
	test.EqualStr(t, td.Len(3).String(), "len=3")
	test.EqualStr(t,
		td.Len(td.Between(3, 8)).String(), "len: 3 ≤ got ≤ 8")
	test.EqualStr(t, td.Len(td.Gt(8)).String(), "len: > 8")

	// Erroneous
	test.EqualStr(t, td.Len("12").String(), "Len(<ERROR>)")
}

func TestCap(t *testing.T) {
	checkOK(t, make([]byte, 0, 4), td.Cap(4))
	checkOK(t, make([]byte, 0, 4), td.Cap(td.Between(4, 6)))

	checkOK(t, [5]int{}, td.Cap(5))
	checkOK(t, [5]int{}, td.Cap(int8(5)))
	checkOK(t, [5]int{}, td.Cap(int16(5)))
	checkOK(t, [5]int{}, td.Cap(int32(5)))
	checkOK(t, [5]int{}, td.Cap(int64(5)))
	checkOK(t, [5]int{}, td.Cap(uint(5)))
	checkOK(t, [5]int{}, td.Cap(uint8(5)))
	checkOK(t, [5]int{}, td.Cap(uint16(5)))
	checkOK(t, [5]int{}, td.Cap(uint32(5)))
	checkOK(t, [5]int{}, td.Cap(uint64(5)))
	checkOK(t, [5]int{}, td.Cap(float32(5)))
	checkOK(t, [5]int{}, td.Cap(float64(5)))
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
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map (map[int]int type)"),
			Expected: mustBe("array OR chan OR slice"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Cap(nil),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received nil as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Cap("12"),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received string as 1st parameter"),
		})

	// out of bounds
	checkError(t, "never tested",
		td.Cap(uint64(math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (18446744073709551615), should be in int range"),
		})

	checkError(t, "never tested",
		td.Cap(float64(math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (1.8446744073709552e+19), should be in int range"),
		})

	checkError(t, "never tested",
		td.Cap(float64(-math.MaxUint64)),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (-1.8446744073709552e+19), should be in int range"),
		})

	checkError(t, "never tested",
		td.Cap(3.1),
		expectedError{
			Message: mustBe("bad usage of Cap operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Cap(TESTDEEP_OPERATOR|INT), but received an out of bounds or not integer 1st parameter (3.1), should be in int range"),
		})

	//
	// String
	test.EqualStr(t, td.Cap(3).String(), "cap=3")
	test.EqualStr(t,
		td.Cap(td.Between(3, 8)).String(), "cap: 3 ≤ got ≤ 8")
	test.EqualStr(t, td.Cap(td.Gt(8)).String(), "cap: > 8")

	// Erroneous op
	test.EqualStr(t, td.Cap("12").String(), "Cap(<ERROR>)")
}

func TestLenCapTypeBehind(t *testing.T) {
	equalTypes(t, td.Cap(3), nil)
	equalTypes(t, td.Len(3), nil)

	// Erroneous op
	equalTypes(t, td.Cap("12"), nil)
	equalTypes(t, td.Len("12"), nil)
}
