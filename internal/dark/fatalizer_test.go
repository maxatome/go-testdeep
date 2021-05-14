// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package dark_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestFatalizer(t *testing.T) {
	test.IsTrue(t, dark.GetFatalizer() != nil)

	tt := test.NewTestingTB("TestFatalizer")
	test.IsFalse(t, dark.CheckFatalizerBarrierErr(tt, func() {}, "no"))
	test.EqualStr(t, tt.LastMessage(), "dark.FatalizerBarrier() did not return an error")
	test.IsFalse(t, tt.IsFatal)
	test.IsTrue(t, tt.HasFailed)

	tt = test.NewTestingTB("TestFatalizer")
	test.IsFalse(t, dark.CheckFatalizerBarrierErr(tt,
		func() {
			f := dark.GetFatalizer()
			f.Helper()
			dark.Fatal(f, "Ouch!")
		}, "Hcuo!"))
	test.EqualStr(t, tt.LastMessage(), `dark.FatalizerBarrier() error "Ouch!"
does not contain "Hcuo!"`)
	test.IsFalse(t, tt.IsFatal)
	test.IsTrue(t, tt.HasFailed)

	test.IsTrue(t, dark.CheckFatalizerBarrierErr(t,
		func() {
			f := dark.GetFatalizer()
			f.Helper()
			dark.Fatal(f, "Ouch!")
		}, "Ouch!"))

	// For the unsafe case, classic panic has to be rethrown
	test.CheckPanic(t,
		func() {
			dark.CheckFatalizerBarrierErr(t, func() { panic("PANIC!") }, "")
		},
		"PANIC!")
}

func TestFatalPanic(t *testing.T) {
	fpIn := dark.FatalPanic("")

	fpIn.Helper()

	test.EqualStr(t, fpIn.String(), "")

	var recovered interface{}
	func() {
		defer func() { recovered = recover() }()
		fpIn.Fatal("User error")
	}()

	if fpOut, ok := recovered.(dark.FatalPanic); test.IsTrue(t, ok) {
		test.EqualStr(t, fpOut.String(), "User error")
	}
}
