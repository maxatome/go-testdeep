// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"regexp"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestShallow(t *testing.T) {
	checkOK(t, nil, nil)

	//
	// Slice
	back := [...]int{1, 2, 3, 1, 2, 3}
	as := back[:3]
	bs := back[3:]
	checkError(t, bs, testdeep.Shallow(back[:]),
		expectedError{
			Message:  mustBe("slice pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	checkOK(t, as, testdeep.Shallow(back[:]))
	checkOK(t, ([]byte)(nil), ([]byte)(nil))

	//
	// Map
	gotMap := map[string]bool{"a": true, "b": false}
	expectedMap := map[string]bool{"a": true, "b": false}
	checkError(t, gotMap, testdeep.Shallow(expectedMap),
		expectedError{
			Message:  mustBe("map pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedMap = gotMap
	checkOK(t, gotMap, testdeep.Shallow(expectedMap))
	checkOK(t, (map[string]bool)(nil), (map[string]bool)(nil))

	//
	// Ptr
	type MyStruct struct {
		val int
	}
	gotPtr := &MyStruct{val: 12}
	expectedPtr := &MyStruct{val: 12}
	checkError(t, gotPtr, testdeep.Shallow(expectedPtr),
		expectedError{
			Message:  mustBe("ptr pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedPtr = gotPtr
	checkOK(t, gotPtr, testdeep.Shallow(expectedPtr))
	checkOK(t, (*MyStruct)(nil), (*MyStruct)(nil))

	//
	// Func
	gotFunc := func(a int) int { return a * 2 }
	expectedFunc := func(a int) int { return a * 2 }
	checkError(t, gotFunc, testdeep.Shallow(expectedFunc),
		expectedError{
			Message:  mustBe("func pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedFunc = gotFunc
	checkOK(t, gotFunc, testdeep.Shallow(expectedFunc))
	checkOK(t, (func(a int) int)(nil), (func(a int) int)(nil))

	//
	// Chan
	gotChan := make(chan int)
	expectedChan := make(chan int)
	checkError(t, gotChan, testdeep.Shallow(expectedChan),
		expectedError{
			Message:  mustBe("chan pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedChan = gotChan
	checkOK(t, gotChan, testdeep.Shallow(expectedChan))
	checkOK(t, (chan int)(nil), (chan int)(nil))

	//
	// String
	backStr := "foobarfoobar!"
	a := backStr[:6]
	b := backStr[6:12]
	checkOK(t, a, testdeep.Shallow(backStr))
	checkOK(t, backStr, testdeep.Shallow(a))
	checkOK(t, b, testdeep.Shallow(backStr[6:7]))

	checkError(t, backStr, testdeep.Shallow(b),
		expectedError{
			Message:  mustBe("string pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})
	checkError(t, b, testdeep.Shallow(backStr),
		expectedError{
			Message:  mustBe("string pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	//
	// Erroneous mix
	checkError(t, gotMap, testdeep.Shallow(expectedChan),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustContain("map"),
			Expected: mustContain("chan"),
		})

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Shallow(42) }, "usage: Shallow")

	//
	//
	reg := regexp.MustCompile(`^\(map\) 0x[a-f0-9]+\z`)
	if !reg.MatchString(testdeep.Shallow(expectedMap).String()) {
		t.Errorf("Shallow().String() failed\n     got: %s\nexpected: %s",
			testdeep.Shallow(expectedMap).String(), reg)
	}
}

func TestShallowTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Shallow(t), nil)
}
