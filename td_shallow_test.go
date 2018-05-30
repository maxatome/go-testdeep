// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"regexp"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestShallow(t *testing.T) {
	//
	// Slice
	gotSlice := []int{1, 2, 3}
	expectedSlice := []int{1, 2, 3}
	checkError(t, gotSlice, Shallow(expectedSlice), expectedError{
		Message:  mustBe("slice pointer mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("0x"),
		Expected: mustContain("0x"),
	})

	expectedSlice = gotSlice
	checkOK(t, gotSlice, Shallow(expectedSlice))
	checkOK(t, ([]byte)(nil), ([]byte)(nil))

	//
	// Map
	gotMap := map[string]bool{"a": true, "b": false}
	expectedMap := map[string]bool{"a": true, "b": false}
	checkError(t, gotMap, Shallow(expectedMap), expectedError{
		Message:  mustBe("map pointer mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("0x"),
		Expected: mustContain("0x"),
	})

	expectedMap = gotMap
	checkOK(t, gotMap, Shallow(expectedMap))
	checkOK(t, (map[string]bool)(nil), (map[string]bool)(nil))

	//
	// Ptr
	type MyStruct struct {
		val int
	}
	gotPtr := &MyStruct{val: 12}
	expectedPtr := &MyStruct{val: 12}
	checkError(t, gotPtr, Shallow(expectedPtr), expectedError{
		Message:  mustBe("ptr pointer mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("0x"),
		Expected: mustContain("0x"),
	})

	expectedPtr = gotPtr
	checkOK(t, gotPtr, Shallow(expectedPtr))
	checkOK(t, (*MyStruct)(nil), (*MyStruct)(nil))

	//
	// Func
	gotFunc := func(a int) int { return a * 2 }
	expectedFunc := func(a int) int { return a * 2 }
	checkError(t, gotFunc, Shallow(expectedFunc), expectedError{
		Message:  mustBe("func pointer mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("0x"),
		Expected: mustContain("0x"),
	})

	expectedFunc = gotFunc
	checkOK(t, gotFunc, Shallow(expectedFunc))
	checkOK(t, (func(a int) int)(nil), (func(a int) int)(nil))

	//
	// Chan
	gotChan := make(chan int)
	expectedChan := make(chan int)
	checkError(t, gotChan, Shallow(expectedChan), expectedError{
		Message:  mustBe("chan pointer mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("0x"),
		Expected: mustContain("0x"),
	})

	expectedChan = gotChan
	checkOK(t, gotChan, Shallow(expectedChan))
	checkOK(t, (chan int)(nil), (chan int)(nil))

	//
	// Erroneous mix
	checkError(t, gotMap, Shallow(expectedChan), expectedError{
		Message:  mustBe("bad kind"),
		Path:     mustBe("DATA"),
		Got:      mustContain("map"),
		Expected: mustContain("chan"),
	})

	//
	// Bad usage
	checkPanic(t, func() { Shallow("test") }, "usage: Shallow")

	//
	//
	reg := regexp.MustCompile(`^\(map\) 0x[a-f0-9]+\z`)
	if !reg.MatchString(Shallow(expectedMap).String()) {
		t.Errorf("Shallow().String() failed\n     got: %s\nexpected: %s",
			Shallow(expectedMap).String(), reg)
	}
}

func TestShallowTypeBehind(t *testing.T) {
	equalTypes(t, Shallow(t), nil)
}
