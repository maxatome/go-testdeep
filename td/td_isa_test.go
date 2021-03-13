// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
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

	checkOK(t, &gotStruct, td.Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), td.Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), td.Isa((*MyStruct)(nil)))
	checkOK(t, gotStruct, td.Isa(MyStruct{}))

	checkOK(t, bytes.NewBufferString("foobar"),
		td.Isa((*fmt.Stringer)(nil)),
		"checks bytes.NewBufferString() implements fmt.Stringer")

	// does bytes.NewBufferString("foobar") implements fmt.Stringer?
	checkOK(t, bytes.NewBufferString("foobar"), td.Isa((*fmt.Stringer)(nil)))

	checkError(t, &gotStruct, td.Isa(&MyStructBase{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*td_test.MyStruct"),
			Expected: mustContain("*td_test.MyStructBase"),
		})

	checkError(t, (*MyStruct)(nil), td.Isa(&MyStructBase{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*td_test.MyStruct"),
			Expected: mustContain("*td_test.MyStructBase"),
		})

	checkError(t, gotStruct, td.Isa(&MyStruct{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("td_test.MyStruct"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, &gotStruct, td.Isa(MyStruct{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*td_test.MyStruct"),
			Expected: mustContain("td_test.MyStruct"),
		})

	gotSlice := []int{1, 2, 3}
	checkOK(t, gotSlice, td.Isa([]int{}))
	checkOK(t, &gotSlice, td.Isa(((*[]int)(nil))))

	checkError(t, &gotSlice, td.Isa([]int{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*[]int"),
			Expected: mustContain("[]int"),
		})

	checkError(t, gotSlice, td.Isa((*[]int)(nil)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("[]int"),
			Expected: mustContain("*[]int"),
		})

	checkError(t, gotSlice, td.Isa([1]int{2}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("[]int"),
			Expected: mustContain("[1]int"),
		})

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Isa(nil) },
		"Isa(nil) is not allowed. To check an interface, try Isa((*fmt.Stringer)(nil)), for fmt.Stringer for example")

	//
	// String
	test.EqualStr(t, td.Isa((*MyStruct)(nil)).String(),
		"*td_test.MyStruct")
}

func TestIsaTypeBehind(t *testing.T) {
	equalTypes(t, td.Isa(([]int)(nil)), []int{})

	equalTypes(t, td.Isa((*fmt.Stringer)(nil)), (*fmt.Stringer)(nil))
}
