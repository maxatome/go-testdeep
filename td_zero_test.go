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

func TestZero(t *testing.T) {
	checkOK(t, 0, Zero())
	checkOK(t, int64(0), Zero())
	checkOK(t, float64(0), Zero())
	checkOK(t, nil, Zero())
	checkOK(t, (map[string]int)(nil), Zero())
	checkOK(t, ([]int)(nil), Zero())
	checkOK(t, [3]int{}, Zero())
	checkOK(t, MyStruct{}, Zero())
	checkOK(t, (*MyStruct)(nil), Zero())
	checkOK(t, (chan int)(nil), Zero())
	checkOK(t, (func())(nil), Zero())
	checkOK(t, false, Zero())

	checkError(t, 12, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, int64(12), Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int64) 12"),
		Expected: mustBe("(int64) 0"),
	})
	checkError(t, float64(12), Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 12"),
		Expected: mustBe("(float64) 0"),
	})
	checkError(t, map[string]int{}, Zero(), expectedError{
		Message:  mustBe("nil map"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, []int{}, Zero(), expectedError{
		Message:  mustBe("nil slice"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, [3]int{0, 12}, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA[1]"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, MyStruct{ValInt: 12}, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA.ValInt"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, &MyStruct{}, Zero(), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("*DATA"),
		// in fact, pointer on 0'ed struct contents
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("nil"),
	})
	checkError(t, true, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(bool) true"),
		Expected: mustBe("(bool) false"),
	})

	//
	// String
	equalStr(t, Zero().String(), "Zero()")
}
