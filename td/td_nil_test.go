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

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestNil(t *testing.T) {
	checkOK(t, (func())(nil), td.Nil())
	checkOK(t, ([]int)(nil), td.Nil())
	checkOK(t, (map[bool]bool)(nil), td.Nil())
	checkOK(t, (*int)(nil), td.Nil())
	checkOK(t, (chan int)(nil), td.Nil())
	checkOK(t, nil, td.Nil())
	checkOK(t,
		map[string]any{"foo": nil},
		map[string]any{"foo": td.Nil()},
	)

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkOK(t, got, td.Nil())

	checkError(t, 42, td.Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustBe("42"),
			Expected: mustBe("nil"),
		})

	num := 42
	checkError(t, &num, td.Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustMatch(`\(\*int\).*42`),
			Expected: mustBe("nil"),
		})

	//
	// String
	test.EqualStr(t, td.Nil().String(), "nil")
}

func TestNotNil(t *testing.T) {
	num := 42
	checkOK(t, func() {}, td.NotNil())
	checkOK(t, []int{}, td.NotNil())
	checkOK(t, map[bool]bool{}, td.NotNil())
	checkOK(t, &num, td.NotNil())
	checkOK(t, 42, td.NotNil())

	checkError(t, (func())(nil), td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, ([]int)(nil), td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (map[bool]bool)(nil), td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (*int)(nil), td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (chan int)(nil), td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, nil, td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	var got fmt.Stringer = (*bytes.Buffer)(nil)
	checkError(t, got, td.NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("<nil>"),
			Expected: mustBe("not nil"),
		})

	//
	// String
	test.EqualStr(t, td.NotNil().String(), "not nil")
}

func TestNilTypeBehind(t *testing.T) {
	equalTypes(t, td.Nil(), nil)
	equalTypes(t, td.NotNil(), nil)
}
