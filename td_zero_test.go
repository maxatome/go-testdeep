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

func TestZero(t *testing.T) {
	checkOK(t, 0, testdeep.Zero())
	checkOK(t, int64(0), testdeep.Zero())
	checkOK(t, float64(0), testdeep.Zero())
	checkOK(t, nil, testdeep.Zero())
	checkOK(t, (map[string]int)(nil), testdeep.Zero())
	checkOK(t, ([]int)(nil), testdeep.Zero())
	checkOK(t, [3]int{}, testdeep.Zero())
	checkOK(t, MyStruct{}, testdeep.Zero())
	checkOK(t, (*MyStruct)(nil), testdeep.Zero())
	checkOK(t, &MyStruct{}, testdeep.Ptr(testdeep.Zero()))
	checkOK(t, (chan int)(nil), testdeep.Zero())
	checkOK(t, (func())(nil), testdeep.Zero())
	checkOK(t, false, testdeep.Zero())

	checkError(t, 12, testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, int64(12), testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int64) 12"),
		Expected: mustBe("(int64) 0"),
	})
	checkError(t, float64(12), testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 12"),
		Expected: mustBe("(float64) 0"),
	})
	checkError(t, map[string]int{}, testdeep.Zero(), expectedError{
		Message:  mustBe("nil map"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, []int{}, testdeep.Zero(), expectedError{
		Message:  mustBe("nil slice"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, [3]int{0, 12}, testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA[1]"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, MyStruct{ValInt: 12}, testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA.ValInt"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, &MyStruct{}, testdeep.Zero(), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("*DATA"),
		// in fact, pointer on 0'ed struct contents
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("nil"),
	})
	checkError(t, true, testdeep.Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(bool) true"),
		Expected: mustBe("(bool) false"),
	})

	//
	// String
	test.EqualStr(t, testdeep.Zero().String(), "Zero()")
}

func TestNotZero(t *testing.T) {
	checkOK(t, 12, testdeep.NotZero())
	checkOK(t, int64(12), testdeep.NotZero())
	checkOK(t, float64(12), testdeep.NotZero())
	checkOK(t, map[string]int{}, testdeep.NotZero())
	checkOK(t, []int{}, testdeep.NotZero())
	checkOK(t, [3]int{1}, testdeep.NotZero())
	checkOK(t, MyStruct{ValInt: 1}, testdeep.NotZero())
	checkOK(t, &MyStruct{}, testdeep.NotZero())
	checkOK(t, make(chan int), testdeep.NotZero())
	checkOK(t, func() {}, testdeep.NotZero())
	checkOK(t, true, testdeep.NotZero())

	checkError(t, nil, testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, 0, testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, int64(0), testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int64) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, float64(0), testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, (map[string]int)(nil), testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(map[string]int) <nil>"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, ([]int)(nil), testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("([]int) <nil>"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, [3]int{}, testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(int) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, MyStruct{}, testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, &MyStruct{}, testdeep.Ptr(testdeep.NotZero()), expectedError{
		Message: mustBe("zero value"),
		Path:    mustBe("*DATA"),
		// in fact, pointer on 0'ed struct contents
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, false, testdeep.NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(bool) false"),
		Expected: mustBe("NotZero()"),
	})

	//
	// String
	test.EqualStr(t, testdeep.NotZero().String(), "NotZero()")
}

func TestZeroTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Zero(), nil)
	equalTypes(t, testdeep.NotZero(), nil)
}
