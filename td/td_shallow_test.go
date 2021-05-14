// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"regexp"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/td"
)

func TestShallow(t *testing.T) {
	checkOK(t, nil, nil)

	//
	// Slice
	back := [...]int{1, 2, 3, 1, 2, 3}
	as := back[:3]
	bs := back[3:]
	checkError(t, bs, td.Shallow(back[:]),
		expectedError{
			Message:  mustBe("slice pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	checkOK(t, as, td.Shallow(back[:]))
	checkOK(t, ([]byte)(nil), ([]byte)(nil))

	//
	// Map
	gotMap := map[string]bool{"a": true, "b": false}
	expectedMap := map[string]bool{"a": true, "b": false}
	checkError(t, gotMap, td.Shallow(expectedMap),
		expectedError{
			Message:  mustBe("map pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedMap = gotMap
	checkOK(t, gotMap, td.Shallow(expectedMap))
	checkOK(t, (map[string]bool)(nil), (map[string]bool)(nil))

	//
	// Ptr
	type MyStruct struct {
		val int
	}
	gotPtr := &MyStruct{val: 12}
	expectedPtr := &MyStruct{val: 12}
	checkError(t, gotPtr, td.Shallow(expectedPtr),
		expectedError{
			Message:  mustBe("ptr pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedPtr = gotPtr
	checkOK(t, gotPtr, td.Shallow(expectedPtr))
	checkOK(t, (*MyStruct)(nil), (*MyStruct)(nil))

	//
	// Func
	gotFunc := func(a int) int { return a * 2 }
	expectedFunc := func(a int) int { return a * 2 }
	checkError(t, gotFunc, td.Shallow(expectedFunc),
		expectedError{
			Message:  mustBe("func pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedFunc = gotFunc
	checkOK(t, gotFunc, td.Shallow(expectedFunc))
	checkOK(t, (func(a int) int)(nil), (func(a int) int)(nil))

	//
	// Chan
	gotChan := make(chan int)
	expectedChan := make(chan int)
	checkError(t, gotChan, td.Shallow(expectedChan),
		expectedError{
			Message:  mustBe("chan pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	expectedChan = gotChan
	checkOK(t, gotChan, td.Shallow(expectedChan))
	checkOK(t, (chan int)(nil), (chan int)(nil))

	//
	// String
	backStr := "foobarfoobar!"
	a := backStr[:6]
	b := backStr[6:12]
	checkOK(t, a, td.Shallow(backStr))
	checkOK(t, backStr, td.Shallow(a))
	checkOK(t, b, td.Shallow(backStr[6:7]))

	checkError(t, backStr, td.Shallow(b),
		expectedError{
			Message:  mustBe("string pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})
	checkError(t, b, td.Shallow(backStr),
		expectedError{
			Message:  mustBe("string pointer mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"),
			Expected: mustContain("0x"),
		})

	//
	// Erroneous mix
	checkError(t, gotMap, td.Shallow(expectedChan),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustContain("map"),
			Expected: mustContain("chan"),
		})

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Shallow(42) }, "usage: Shallow")

	//
	//
	reg := regexp.MustCompile(`^\(map\) 0x[a-f0-9]+\z`)
	if !reg.MatchString(td.Shallow(expectedMap).String()) {
		t.Errorf("Shallow().String() failed\n     got: %s\nexpected: %s",
			td.Shallow(expectedMap).String(), reg)
	}
}

func TestShallowTypeBehind(t *testing.T) {
	equalTypes(t, td.Shallow(t), nil)
}
