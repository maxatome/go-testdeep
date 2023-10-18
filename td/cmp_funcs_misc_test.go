// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func ExampleCmpTrue() {
	t := &testing.T{}

	got := true
	ok := td.CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	got = false
	ok = td.CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpFalse() {
	t := &testing.T{}

	got := false
	ok := td.CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	got = true
	ok = td.CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := td.CmpError(t, got, "An error occurred")
	fmt.Println(ok)

	got = nil
	ok = td.CmpError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpNoError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := td.CmpNoError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	got = nil
	ok = td.CmpNoError(t, got, "An error occurred")
	fmt.Println(ok)

	// Output:
	// false
	// true
}

func ExampleCmpPanic() {
	t := &testing.T{}

	ok := td.CmpPanic(t,
		func() { panic("I am panicking!") }, "I am panicking!",
		"Checks for panic")
	fmt.Println("checks exact panic() string:", ok)

	// Can use TestDeep operator too
	ok = td.CmpPanic(t,
		func() { panic("I am panicking!") }, td.Contains("panicking!"),
		"Checks for panic")
	fmt.Println("checks panic() sub-string:", ok)

	// Can detect panic(nil)
	// Before Go 1.21, programs that called panic(nil) observed recover
	// returning nil. Starting in Go 1.21, programs that call panic(nil)
	// observe recover returning a *PanicNilError. Programs can change
	// back to the old behavior by setting GODEBUG=panicnil=1.
	// See https://pkg.go.dev/runtime#PanicNilError
	ok = td.CmpPanic(t, func() { panic(nil) }, nil, "Checks for panic(nil)")
	fmt.Println("checks for panic(nil):", ok)

	// As well as structured data panic
	type PanicStruct struct {
		Error string
		Code  int
	}

	ok = td.CmpPanic(t,
		func() {
			panic(PanicStruct{Error: "Memory violation", Code: 11})
		},
		PanicStruct{
			Error: "Memory violation",
			Code:  11,
		})
	fmt.Println("checks exact panic() struct:", ok)

	// or combined with TestDeep operators too
	ok = td.CmpPanic(t,
		func() {
			panic(PanicStruct{Error: "Memory violation", Code: 11})
		},
		td.Struct(PanicStruct{}, td.StructFields{
			"Code": td.Between(10, 20),
		}))
	fmt.Println("checks panic() struct against TestDeep operators:", ok)

	// Of course, do not panic = test failure, even for expected nil
	// panic parameter
	ok = td.CmpPanic(t, func() {}, nil)
	fmt.Println("checks a panic occurred:", ok)

	// Output:
	// checks exact panic() string: true
	// checks panic() sub-string: true
	// checks for panic(nil): true
	// checks exact panic() struct: true
	// checks panic() struct against TestDeep operators: true
	// checks a panic occurred: false
}

func ExampleCmpNotPanic() {
	t := &testing.T{}

	ok := td.CmpNotPanic(t, func() {})
	fmt.Println("checks a panic DID NOT occur:", ok)

	// Classic panic
	ok = td.CmpNotPanic(t, func() { panic("I am panicking!") },
		"Hope it does not panic!")
	fmt.Println("still no panic?", ok)

	// Can detect panic(nil)
	ok = td.CmpNotPanic(t, func() { panic(nil) }, "Checks for panic(nil)")
	fmt.Println("last no panic?", ok)

	// Output:
	// checks a panic DID NOT occur: true
	// still no panic? false
	// last no panic? false
}
