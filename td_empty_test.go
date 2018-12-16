// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestEmpty(t *testing.T) {
	checkOK(t, nil, testdeep.Empty())
	checkOK(t, "", testdeep.Empty())
	checkOK(t, ([]int)(nil), testdeep.Empty())
	checkOK(t, []int{}, testdeep.Empty())
	checkOK(t, (map[string]bool)(nil), testdeep.Empty())
	checkOK(t, map[string]bool{}, testdeep.Empty())
	checkOK(t, (chan int)(nil), testdeep.Empty())
	checkOK(t, make(chan int), testdeep.Empty())
	checkOK(t, [0]int{}, testdeep.Empty())

	type MySlice []int
	checkOK(t, MySlice{}, testdeep.Empty())
	checkOK(t, &MySlice{}, testdeep.Empty())

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, testdeep.Empty())

	l1 = nil
	checkOK(t, &l3, testdeep.Empty())

	checkError(t, 12, testdeep.Empty(),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
		})

	num := 12
	n1 := &num
	n2 := &n1
	n3 := &n2
	checkError(t, &n3, testdeep.Empty(),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("****int"),
			Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
		})

	n1 = nil
	checkError(t, &n3, testdeep.Empty(),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("****int"),
			Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
		})

	checkError(t, "foobar", testdeep.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foobar"`),
			Expected: mustBe("empty"),
		})
	checkError(t, []int{1}, testdeep.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("1"),
			Expected: mustBe("empty"),
		})
	checkError(t, map[string]bool{"foo": true}, testdeep.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo": (bool) true`),
			Expected: mustBe("empty"),
		})

	ch := make(chan int, 1)
	ch <- 42
	checkError(t, ch, testdeep.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int)"),
			Expected: mustBe("empty"),
		})

	checkError(t, [3]int{}, testdeep.Empty(),
		expectedError{
			Message:  mustBe("not empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0"),
			Expected: mustBe("empty"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Empty().String(), "Empty()")
}

func TestNotEmpty(t *testing.T) {
	checkOK(t, "foobar", testdeep.NotEmpty())
	checkOK(t, []int{1}, testdeep.NotEmpty())
	checkOK(t, map[string]bool{"foo": true}, testdeep.NotEmpty())
	checkOK(t, [3]int{}, testdeep.NotEmpty())

	ch := make(chan int, 1)
	ch <- 42
	checkOK(t, ch, testdeep.NotEmpty())

	type MySlice []int
	checkOK(t, MySlice{1}, testdeep.NotEmpty())
	checkOK(t, &MySlice{1}, testdeep.NotEmpty())

	l1 := &MySlice{1}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, testdeep.NotEmpty())

	checkError(t, 12, testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
		})

	checkError(t, nil, testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not empty"),
		})
	checkError(t, "", testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`""`),
			Expected: mustBe("not empty"),
		})
	checkError(t, ([]int)(nil), testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustBe("([]int) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, []int{}, testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("([]int)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, (map[string]bool)(nil), testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(map[string]bool) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, map[string]bool{}, testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(map[string]bool)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, (chan int)(nil), testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int) <nil>"),
			Expected: mustBe("not empty"),
		})
	checkError(t, make(chan int), testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("(chan int)"),
			Expected: mustBe("not empty"),
		})
	checkError(t, [0]int{}, testdeep.NotEmpty(),
		expectedError{
			Message:  mustBe("empty"),
			Path:     mustBe("DATA"),
			Got:      mustContain("([0]int)"),
			Expected: mustBe("not empty"),
		})

	//
	// String
	test.EqualStr(t, testdeep.NotEmpty().String(), "NotEmpty()")
}

func TestEmptyTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Empty(), nil)
	equalTypes(t, testdeep.NotEmpty(), nil)
}
