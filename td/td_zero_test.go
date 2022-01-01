// Copyright (c) 2018-2022, Maxime Soulé
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

func TestZero(t *testing.T) {
	checkOK(t, 0, td.Zero())
	checkOK(t, int64(0), td.Zero())
	checkOK(t, float64(0), td.Zero())
	checkOK(t, nil, td.Zero())
	checkOK(t, (map[string]int)(nil), td.Zero())
	checkOK(t, ([]int)(nil), td.Zero())
	checkOK(t, [3]int{}, td.Zero())
	checkOK(t, MyStruct{}, td.Zero())
	checkOK(t, (*MyStruct)(nil), td.Zero())
	checkOK(t, &MyStruct{}, td.Ptr(td.Zero()))
	checkOK(t, (chan int)(nil), td.Zero())
	checkOK(t, (func())(nil), td.Zero())
	checkOK(t, false, td.Zero())

	checkError(t, 12, td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("0"),
		})
	checkError(t, int64(12), td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 0"),
		})
	checkError(t, float64(12), td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12.0"),
			Expected: mustBe("0.0"),
		})
	checkError(t, map[string]int{}, td.Zero(),
		expectedError{
			Message:  mustBe("nil map"),
			Path:     mustBe("DATA"),
			Got:      mustBe("not nil"),
			Expected: mustBe("nil"),
		})
	checkError(t, []int{}, td.Zero(),
		expectedError{
			Message:  mustBe("nil slice"),
			Path:     mustBe("DATA"),
			Got:      mustBe("not nil"),
			Expected: mustBe("nil"),
		})
	checkError(t, [3]int{0, 12}, td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("12"),
			Expected: mustBe("0"),
		})
	checkError(t, MyStruct{ValInt: 12}, td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValInt"),
			Got:      mustBe("12"),
			Expected: mustBe("0"),
		})
	checkError(t, &MyStruct{}, td.Zero(),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("*DATA"),
			// in fact, pointer on 0'ed struct contents
			Got:      mustContain(`ValInt: (int) 0`),
			Expected: mustBe("nil"),
		})
	checkError(t, true, td.Zero(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("true"),
			Expected: mustBe("false"),
		})

	//
	// String
	test.EqualStr(t, td.Zero().String(), "Zero()")
}

func TestNotZero(t *testing.T) {
	checkOK(t, 12, td.NotZero())
	checkOK(t, int64(12), td.NotZero())
	checkOK(t, float64(12), td.NotZero())
	checkOK(t, map[string]int{}, td.NotZero())
	checkOK(t, []int{}, td.NotZero())
	checkOK(t, [3]int{1}, td.NotZero())
	checkOK(t, MyStruct{ValInt: 1}, td.NotZero())
	checkOK(t, &MyStruct{}, td.NotZero())
	checkOK(t, make(chan int), td.NotZero())
	checkOK(t, func() {}, td.NotZero())
	checkOK(t, true, td.NotZero())

	checkError(t, nil, td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, 0, td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("0"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, int64(0), td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(int64) 0"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, float64(0), td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("0.0"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, (map[string]int)(nil), td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(map[string]int) <nil>"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, ([]int)(nil), td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("([]int) <nil>"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, [3]int{}, td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0"),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, MyStruct{}, td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`ValInt: (int) 0`),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, &MyStruct{}, td.Ptr(td.NotZero()),
		expectedError{
			Message: mustBe("zero value"),
			Path:    mustBe("*DATA"),
			// in fact, pointer on 0'ed struct contents
			Got:      mustContain(`ValInt: (int) 0`),
			Expected: mustBe("NotZero()"),
		})
	checkError(t, false, td.NotZero(),
		expectedError{
			Message:  mustBe("zero value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("false"),
			Expected: mustBe("NotZero()"),
		})

	//
	// String
	test.EqualStr(t, td.NotZero().String(), "NotZero()")
}

func TestZeroTypeBehind(t *testing.T) {
	equalTypes(t, td.Zero(), nil)
	equalTypes(t, td.NotZero(), nil)
}
