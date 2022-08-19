// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestEmpty(t *testing.T) {
	checkOK(t, nil, td.Empty())
	checkOK(t, "", td.Empty())
	checkOK(t, ([]int)(nil), td.Empty())
	checkOK(t, []int{}, td.Empty())
	checkOK(t, (map[string]bool)(nil), td.Empty())
	checkOK(t, map[string]bool{}, td.Empty())
	checkOK(t, (chan int)(nil), td.Empty())
	checkOK(t, make(chan int), td.Empty())
	checkOK(t, [0]int{}, td.Empty())

	type MySlice []int
	checkOK(t, MySlice{}, td.Empty())
	checkOK(t, &MySlice{}, td.Empty())

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, td.Empty())

	l1 = nil
	checkOK(t, &l3, td.Empty())

	checkError(t, 12, td.Empty(),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("array OR chan OR map OR slice OR string OR pointer(s) on them"),
		})

	num := 12
	n1 := &num
	n2 := &n1
	n3 := &n2
	checkError(t, &n3, td.Empty(),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("****int"),
			Expected: mustBe("array OR chan OR map OR slice OR string OR pointer(s) on them"),
		})

	n1 = nil
	checkError(t, &n3, td.Empty(),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("****int"),
			Expected: mustBe("array OR chan OR map OR slice OR string OR pointer(s) on them"),
		})

	checkError(t, "foobar", td.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foobar"`),
			Expected: mustBe("empty"),
		})
	checkError(t, []int{1}, td.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("1"),
			Expected: mustBe("empty"),
		})
	checkError(t, map[string]bool{"foo": true}, td.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo": (bool) true`),
			Expected: mustBe("empty"),
		})

	ch := make(chan int, 1)
	ch <- 42
	checkError(t, ch, td.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int)"),
			Expected: mustBe("empty"),
		})

	checkError(t, [3]int{}, td.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0"),
			Expected: mustBe("empty"),
		})

	//
	// String
	test.EqualStr(t, td.Empty().String(), "Empty()")
}

func TestNotEmpty(t *testing.T) {
	checkOK(t, "foobar", td.NotEmpty())
	checkOK(t, []int{1}, td.NotEmpty())
	checkOK(t, map[string]bool{"foo": true}, td.NotEmpty())
	checkOK(t, [3]int{}, td.NotEmpty())

	ch := make(chan int, 1)
	ch <- 42
	checkOK(t, ch, td.NotEmpty())

	type MySlice []int
	checkOK(t, MySlice{1}, td.NotEmpty())
	checkOK(t, &MySlice{1}, td.NotEmpty())

	l1 := &MySlice{1}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, td.NotEmpty())

	checkError(t, 12, td.NotEmpty(),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("array OR chan OR map OR slice OR string OR pointer(s) on them"),
		})

	checkError(t, nil, td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not empty"),
		})
	checkError(t, "", td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`""`),
			Expected: mustBe("not empty"),
		})
	checkError(t, ([]int)(nil), td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustBe("([]int) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, []int{}, td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("([]int)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, (map[string]bool)(nil), td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(map[string]bool) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, map[string]bool{}, td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(map[string]bool)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, (chan int)(nil), td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, make(chan int), td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, [0]int{}, td.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("([0]int)"),
			Expected: mustBe("not empty"),
		})

	//
	// String
	test.EqualStr(t, td.NotEmpty().String(), "NotEmpty()")
}

func TestEmptyTypeBehind(t *testing.T) {
	equalTypes(t, td.Empty(), nil)
	equalTypes(t, td.NotEmpty(), nil)
}
