// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestEmpty(t *testing.T) {
	checkOK(t, nil, Empty())
	checkOK(t, "", Empty())
	checkOK(t, ([]int)(nil), Empty())
	checkOK(t, []int{}, Empty())
	checkOK(t, (map[string]bool)(nil), Empty())
	checkOK(t, map[string]bool{}, Empty())
	checkOK(t, (chan int)(nil), Empty())
	checkOK(t, make(chan int), Empty())
	checkOK(t, [0]int{}, Empty())

	type MySlice []int
	checkOK(t, MySlice{}, Empty())
	checkOK(t, &MySlice{}, Empty())

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, Empty())

	l1 = nil
	checkOK(t, &l3, Empty())

	checkError(t, 12, Empty(), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
	})

	num := 12
	n1 := &num
	n2 := &n1
	n3 := &n2
	checkError(t, &n3, Empty(), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("****int"),
		Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
	})

	n1 = nil
	checkError(t, &n3, Empty(), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("****int"),
		Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
	})

	checkError(t, "foobar", Empty(), expectedError{
		Message:  mustBe("not empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foobar"`),
		Expected: mustBe("empty"),
	})
	checkError(t, []int{1}, Empty(), expectedError{
		Message:  mustBe("not empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(int) 1"),
		Expected: mustBe("empty"),
	})
	checkError(t, map[string]bool{"foo": true}, Empty(), expectedError{
		Message:  mustBe("not empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo": (bool) true`),
		Expected: mustBe("empty"),
	})

	ch := make(chan int, 1)
	ch <- 42
	checkError(t, ch, Empty(), expectedError{
		Message:  mustBe("not empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(chan int)"),
		Expected: mustBe("empty"),
	})

	checkError(t, [3]int{}, Empty(), expectedError{
		Message:  mustBe("not empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(int) 0"),
		Expected: mustBe("empty"),
	})

	//
	// String
	test.EqualStr(t, Empty().String(), "Empty()")
}

func TestNotEmpty(t *testing.T) {
	checkOK(t, "foobar", NotEmpty())
	checkOK(t, []int{1}, NotEmpty())
	checkOK(t, map[string]bool{"foo": true}, NotEmpty())
	checkOK(t, [3]int{}, NotEmpty())

	ch := make(chan int, 1)
	ch <- 42
	checkOK(t, ch, NotEmpty())

	type MySlice []int
	checkOK(t, MySlice{1}, NotEmpty())
	checkOK(t, &MySlice{1}, NotEmpty())

	l1 := &MySlice{1}
	l2 := &l1
	l3 := &l2
	checkOK(t, &l3, NotEmpty())

	checkError(t, 12, NotEmpty(), expectedError{
		Message:  mustBe("bad type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("Array, Chan, Map, Slice, string or pointer(s) on them"),
	})

	checkError(t, nil, NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("nil"),
		Expected: mustBe("not empty"),
	})
	checkError(t, "", NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`""`),
		Expected: mustBe("not empty"),
	})
	checkError(t, ([]int)(nil), NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustBe("([]int) <nil>"),
		Expected: mustBe("not empty"),
	})
	checkError(t, []int{}, NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("([]int)"),
		Expected: mustBe("not empty"),
	})
	checkError(t, (map[string]bool)(nil), NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(map[string]bool) <nil>"),
		Expected: mustBe("not empty"),
	})
	checkError(t, map[string]bool{}, NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(map[string]bool)"),
		Expected: mustBe("not empty"),
	})
	checkError(t, (chan int)(nil), NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(chan int) <nil>"),
		Expected: mustBe("not empty"),
	})
	checkError(t, make(chan int), NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(chan int)"),
		Expected: mustBe("not empty"),
	})
	checkError(t, [0]int{}, NotEmpty(), expectedError{
		Message:  mustBe("empty"),
		Path:     mustBe("DATA"),
		Got:      mustContain("([0]int)"),
		Expected: mustBe("not empty"),
	})

	//
	// String
	test.EqualStr(t, NotEmpty().String(), "NotEmpty()")
}

func TestEmptyTypeBehind(t *testing.T) {
	equalTypes(t, Empty(), nil)
	equalTypes(t, NotEmpty(), nil)
}
