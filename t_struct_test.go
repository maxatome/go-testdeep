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

func ExampleT_CmpError() {
	t := NewT(&testing.T{})

	got := fmt.Errorf("Error #%d", 42)
	ok := t.CmpError(got, "An error occurred")
	fmt.Println(ok)

	got = nil
	ok = t.CmpError(got, "An error occurred") // fails
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_CmpNoError() {
	t := NewT(&testing.T{})

	got := fmt.Errorf("Error #%d", 42)
	ok := t.CmpNoError(got, "An error occurred") // fails
	fmt.Println(ok)

	got = nil
	ok = t.CmpNoError(got, "An error occurred")
	fmt.Println(ok)

	// Output:
	// false
	// true
}

func TestT(tt *testing.T) {
	t := NewT(tt)
	CmpDeeply(tt, t.Config, DefaultContextConfig)

	t = NewT(tt, ContextConfig{})
	CmpDeeply(tt, t.Config, DefaultContextConfig)

	conf := ContextConfig{
		RootName:  "TEST",
		MaxErrors: 33,
	}
	t = NewT(tt, conf)
	CmpDeeply(tt, t.Config, conf)

	t2 := t.RootName("T2")
	CmpDeeply(tt, t.Config, conf)
	CmpDeeply(tt, t2.Config, ContextConfig{
		RootName:  "T2",
		MaxErrors: 33,
	})

	//
	// Bad usage
	checkPanic(tt,
		func() { NewT(tt, ContextConfig{}, ContextConfig{}) },
		"usage: NewT")
}

func TestRun(tt *testing.T) {
	t := NewT(tt)

	runPassed := false

	ok := t.Run("Test level1",
		func(t *T) {
			ok := t.Run("Test level2",
				func(t *T) {
					runPassed = t.True(true) // test succeeds!
				})

			t.True(ok)
		})

	t.True(ok)
	t.True(runPassed)
}
