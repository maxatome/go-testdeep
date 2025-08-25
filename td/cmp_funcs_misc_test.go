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
