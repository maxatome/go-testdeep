// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestNil(t *testing.T) {
	checkOK(t, (func())(nil), Nil())
	checkOK(t, ([]int)(nil), Nil())
	checkOK(t, (map[bool]bool)(nil), Nil())
	checkOK(t, (*int)(nil), Nil())
	checkOK(t, (chan int)(nil), Nil())
	checkOK(t, nil, Nil())

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkOK(t, got, Nil())

	checkError(t, 42, Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(int) 42"),
			Expected: mustBe("nil"),
		})

	num := 42
	checkError(t, &num, Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustMatch(`\(\*int\).*42`),
			Expected: mustBe("nil"),
		})

	//
	// String
	equalStr(t, Nil().String(), "nil")
}

func TestNotNil(t *testing.T) {
	num := 42
	checkOK(t, func() {}, NotNil())
	checkOK(t, []int{}, NotNil())
	checkOK(t, map[bool]bool{}, NotNil())
	checkOK(t, &num, NotNil())
	checkOK(t, 42, NotNil())

	checkError(t, (func())(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, ([]int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (map[bool]bool)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (*int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (chan int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, nil, NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkError(t, got, NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("<nil>"),
			Expected: mustBe("not nil"),
		})

	//
	// String
	equalStr(t, NotNil().String(), "not nil")
}

func TestNilTypeBehind(t *testing.T) {
	equalTypes(t, Nil(), nil)
	equalTypes(t, NotNil(), nil)
}
