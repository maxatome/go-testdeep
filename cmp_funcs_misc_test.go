// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func ExampleCmpTrue() {
	t := &testing.T{}

	got := true
	ok := CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	got = false
	ok = CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpFalse() {
	t := &testing.T{}

	got := false
	ok := CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	got = true
	ok = CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := CmpError(t, got, "An error occurred")
	fmt.Println(ok)

	got = nil
	ok = CmpError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpNoError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := CmpNoError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	got = nil
	ok = CmpNoError(t, got, "An error occurred")
	fmt.Println(ok)

	// Output:
	// false
	// true
}
