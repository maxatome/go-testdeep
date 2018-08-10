// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	. "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestIsa(t *testing.T) {
	var gotStruct = MyStruct{
		MyStructMid: MyStructMid{
			MyStructBase: MyStructBase{
				ValBool: true,
			},
			ValStr: "foobar",
		},
		ValInt: 123,
	}

	checkOK(t, &gotStruct, Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), Isa((*MyStruct)(nil)))
	checkOK(t, gotStruct, Isa(MyStruct{}))

	checkOK(t, bytes.NewBufferString("foobar"), Isa((*fmt.Stringer)(nil)),
		"checks bytes.NewBufferString() implements fmt.Stringer")

	var ifstr fmt.Stringer = regexp.MustCompile("aa")
	checkOK(t, bytes.NewBufferString("foobar"), Isa(&ifstr))

	checkError(t, &gotStruct, Isa(&MyStructBase{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStructBase"),
	})

	checkError(t, (*MyStruct)(nil), Isa(&MyStructBase{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStructBase"),
	})

	checkError(t, gotStruct, Isa(&MyStruct{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStruct"),
	})

	checkError(t, &gotStruct, Isa(MyStruct{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("testdeep_test.MyStruct"),
	})

	gotSlice := []int{1, 2, 3}
	checkOK(t, gotSlice, Isa([]int{}))
	checkOK(t, &gotSlice, Isa(((*[]int)(nil))))

	checkError(t, &gotSlice, Isa([]int{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*[]int"),
		Expected: mustContain("[]int"),
	})

	checkError(t, gotSlice, Isa((*[]int)(nil)), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("[]int"),
		Expected: mustContain("*[]int"),
	})

	checkError(t, gotSlice, Isa([1]int{2}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("[]int"),
		Expected: mustContain("[1]int"),
	})

	//
	// String
	test.EqualStr(t, Isa((*MyStruct)(nil)).String(), "*testdeep_test.MyStruct")
}

func TestIsaTypeBehind(t *testing.T) {
	equalTypes(t, Isa(([]int)(nil)), []int{})
}
