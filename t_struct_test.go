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

func TestT(tt *testing.T) {
	t := testdeep.NewT(tt)
	testdeep.Cmp(tt, t.Config, testdeep.DefaultContextConfig)

	t = testdeep.NewT(tt, testdeep.ContextConfig{})
	testdeep.Cmp(tt, t.Config, testdeep.DefaultContextConfig)

	conf := testdeep.ContextConfig{
		RootName:  "TEST",
		MaxErrors: 33,
	}
	t = testdeep.NewT(tt, conf)
	testdeep.Cmp(tt, t.Config, conf)

	t2 := t.RootName("T2")
	testdeep.Cmp(tt, t.Config, conf)
	testdeep.Cmp(tt, t2.Config, testdeep.ContextConfig{
		RootName:  "T2",
		MaxErrors: 33,
	})

	//
	// Bad usage
	test.CheckPanic(tt,
		func() {
			testdeep.NewT(tt, testdeep.ContextConfig{}, testdeep.ContextConfig{})
		},
		"usage: NewT")
}

func TestTCmp(tt *testing.T) {
	ttt := &test.TestingFT{}
	t := testdeep.NewT(ttt)
	testdeep.CmpTrue(tt, t.Cmp(1, 1))
	testdeep.CmpFalse(tt, ttt.Failed())

	ttt = &test.TestingFT{}
	t = testdeep.NewT(ttt)
	testdeep.CmpFalse(tt, t.Cmp(1, 2))
	testdeep.CmpTrue(tt, ttt.Failed())
}

func TestTCmpDeeply(tt *testing.T) {
	ttt := &test.TestingFT{}
	t := testdeep.NewT(ttt)
	testdeep.CmpTrue(tt, t.CmpDeeply(1, 1))
	testdeep.CmpFalse(tt, ttt.Failed())

	ttt = &test.TestingFT{}
	t = testdeep.NewT(ttt)
	testdeep.CmpFalse(tt, t.CmpDeeply(1, 2))
	testdeep.CmpTrue(tt, ttt.Failed())
}

func TestRun(tt *testing.T) {
	t := testdeep.NewT(tt)

	runPassed := false

	ok := t.Run("Test level1",
		func(t *testdeep.T) {
			ok := t.Run("Test level2",
				func(t *testdeep.T) {
					runPassed = t.True(true) // test succeeds!
				})

			t.True(ok)
		})

	t.True(ok)
	t.True(runPassed)
}

func TestFailureIsFatal(tt *testing.T) {
	ttt := &test.TestingFT{}

	// All t.True(false) tests of course fail

	// Using default config
	t := testdeep.NewT(ttt)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "by default it not fatal")

	// Using specific config
	t = testdeep.NewT(ttt, testdeep.ContextConfig{FailureIsFatal: true})
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal()
	t = testdeep.NewT(ttt).FailureIsFatal()
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Using FailureIsFatal(true)
	t = testdeep.NewT(ttt).FailureIsFatal(true)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpTrue(tt, ttt.IsFatal, "it must be fatal")

	// Canceling specific config
	t = testdeep.NewT(ttt, testdeep.ContextConfig{FailureIsFatal: false}).
		FailureIsFatal(false)
	t.True(false) // failure
	testdeep.CmpNotEmpty(tt, ttt.LastMessage)
	testdeep.CmpFalse(tt, ttt.IsFatal, "it must be not fatal")
}
