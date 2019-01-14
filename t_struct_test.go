// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestT(tt *testing.T) {
	t := testdeep.NewT(tt)
	testdeep.CmpDeeply(tt, t.Config, testdeep.DefaultContextConfig)

	t = testdeep.NewT(tt, testdeep.ContextConfig{})
	testdeep.CmpDeeply(tt, t.Config, testdeep.DefaultContextConfig)

	conf := testdeep.ContextConfig{
		RootName:  "TEST",
		MaxErrors: 33,
	}
	t = testdeep.NewT(tt, conf)
	testdeep.CmpDeeply(tt, t.Config, conf)

	t2 := t.RootName("T2")
	testdeep.CmpDeeply(tt, t.Config, conf)
	testdeep.CmpDeeply(tt, t2.Config, testdeep.ContextConfig{
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

// Just to test the case where t is an interface and not a *testing.T
// See t.Helper() issue in all tested methods.
func TestStructWithInterfaceT(tt *testing.T) {
	ttt := &test.TestingFT{}

	t := testdeep.NewT(ttt)

	test.IsTrue(tt, t.False(false))
	test.IsFalse(tt, t.CmpError(nil))
	test.IsFalse(tt, t.CmpNoError(errors.New("error")))
	test.IsFalse(tt, t.CmpPanic(func() {}, "panic")) // no panic occurred
	test.IsTrue(tt, t.CmpNotPanic(func() {}))
	test.IsTrue(tt, t.Run("test", func(t *testdeep.T) {}))

	test.IsFalse(tt, t.All(0, []interface{}{12}))
	test.IsFalse(tt, t.Any(0, nil))
	test.IsFalse(tt, t.Array(0, [2]int{}, nil))
	test.IsFalse(tt, t.ArrayEach(0, nil))
	test.IsFalse(tt, t.Bag(0, nil))
	test.IsFalse(tt, t.Between(0, 1, 2, testdeep.BoundsInIn))
	test.IsFalse(tt, t.Cap(nil, 12))
	test.IsFalse(tt, t.Code(0, func(n int) bool { return false }))
	test.IsFalse(tt, t.Contains(0, nil))
	test.IsFalse(tt, t.ContainsKey(map[bool]int{}, true))
	test.IsFalse(tt, t.Empty(0))
	test.IsFalse(tt, t.Gt(0, 12))
	test.IsFalse(tt, t.Gte(0, 12))
	test.IsFalse(tt, t.HasPrefix(0, "pipo"))
	test.IsFalse(tt, t.HasSuffix(0, "pipo"))
	test.IsFalse(tt, t.Isa(0, "string"))
	test.IsFalse(tt, t.Len(nil, 12))
	test.IsFalse(tt, t.Lt(0, -12))
	test.IsFalse(tt, t.Lte(0, -12))
	test.IsFalse(tt, t.Map(0, map[int]bool{}, nil))
	test.IsFalse(tt, t.MapEach(0, nil))
	test.IsFalse(tt, t.N(0, 12, 0))
	test.IsFalse(tt, t.NaN(0, nil))
	test.IsFalse(tt, t.Nil(0))
	test.IsFalse(tt, t.None(0, []interface{}{0}))
	test.IsFalse(tt, t.Not(0, 0))
	test.IsFalse(tt, t.NotAny(0, nil))
	test.IsFalse(tt, t.NotEmpty(0, nil))
	test.IsFalse(tt, t.NotNaN(0, nil))
	test.IsFalse(tt, t.NotNil(nil))
	test.IsFalse(tt, t.NotZero(0, nil))
	test.IsFalse(tt, t.PPtr(0, 12))
	test.IsFalse(tt, t.Ptr(0, 12))
	test.IsFalse(tt, t.Re(0, "pipo", nil))
	test.IsFalse(tt, t.ReAll(0, "pipo", nil))
	test.IsFalse(tt, t.Set(0, nil))
	test.IsFalse(tt, t.Shallow(0, []int{}))
	test.IsFalse(tt, t.Slice(0, []int{}, nil))
	test.IsFalse(tt, t.Smuggle(0, func(n int) int { return 0 }, 12))
	test.IsFalse(tt, t.String(0, "pipo"))
	test.IsFalse(tt, t.Struct(0, struct{}{}, nil))
	test.IsFalse(tt, t.SubBagOf(0, nil))
	test.IsFalse(tt, t.SubMapOf(0, map[int]bool{}, nil))
	test.IsFalse(tt, t.SubSetOf(0, nil))
	test.IsFalse(tt, t.SuperBagOf(0, nil))
	test.IsFalse(tt, t.SuperMapOf(0, map[int]bool{}, nil))
	test.IsFalse(tt, t.SuperSetOf(0, nil))
	test.IsFalse(tt, t.TruncTime(0, time.Now(), time.Second))
	test.IsFalse(tt, t.Zero(12))
}
