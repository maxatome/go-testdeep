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
)

func TestShallow(t *testing.T) {
	//
	// Slice
	gotSlice := []int{1, 2, 3}
	expectedSlice := []int{1, 2, 3}
	checkError(t, gotSlice, testdeep.Shallow(expectedSlice),
		expectedError{
			Message:  mustBe("slice pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedSlice = gotSlice
	checkOK(t, gotSlice, testdeep.Shallow(expectedSlice))
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
	checkPanic(t, func() { testdeep.Shallow("test") }, "usage: Shallow")

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
