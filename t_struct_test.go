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

func ExampleT_True() {
	t := NewT(&testing.T{})

	got := true
	ok := t.True(got, "check that got is true!")
	fmt.Println(ok)

	got = false
	ok = t.True(got, "check that got is true!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_False() {
	t := NewT(&testing.T{})

	got := false
	ok := t.False(got, "check that got is false!")
	fmt.Println(ok)

	got = true
	ok = t.False(got, "check that got is false!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}
