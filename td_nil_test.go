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

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestNil(t *testing.T) {
	checkOK(t, (func())(nil), testdeep.Nil())
	checkOK(t, ([]int)(nil), testdeep.Nil())
	checkOK(t, (map[bool]bool)(nil), testdeep.Nil())
	checkOK(t, (*int)(nil), testdeep.Nil())
	checkOK(t, (chan int)(nil), testdeep.Nil())
	checkOK(t, nil, testdeep.Nil())

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkOK(t, got, testdeep.Nil())

	checkError(t, 42, testdeep.Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustBe("42"),
			Expected: mustBe("nil"),
		})

	num := 42
	checkError(t, &num, testdeep.Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustMatch(`\(\*int\).*42`),
			Expected: mustBe("nil"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Nil().String(), "nil")
}

func TestNotNil(t *testing.T) {
	num := 42
	checkOK(t, func() {}, testdeep.NotNil())
	checkOK(t, []int{}, testdeep.NotNil())
	checkOK(t, map[bool]bool{}, testdeep.NotNil())
	checkOK(t, &num, testdeep.NotNil())
	checkOK(t, 42, testdeep.NotNil())

	checkError(t, (func())(nil), testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, ([]int)(nil), testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (map[bool]bool)(nil), testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (*int)(nil), testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (chan int)(nil), testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, nil, testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkError(t, got, testdeep.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("<nil>"),
			Expected: mustBe("not nil"),
		})

	//
	// String
	test.EqualStr(t, testdeep.NotNil().String(), "not nil")
}

func TestNilTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Nil(), nil)
	equalTypes(t, testdeep.NotNil(), nil)
}
