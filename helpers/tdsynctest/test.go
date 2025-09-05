// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.25
// +build go1.25

package tdsynctest

import (
	"testing"
	"testing/synctest"

	"github.com/maxatome/go-testdeep/td"
)

// Test is a wrapper around [synctest.Test].
//
// The t param of f inherits the configuration of t.
//
// t.TB (set by [td.NewT], [td.Assert], [td.Require], …) must be a
// [*testing.T] instance.
func Test(t *td.T, f func(t *td.T)) {
	tt, ok := t.TB.(*testing.T)
	if !ok {
		t.Helper()
		t.Fatalf("tdsynctest.Test only works if underlying T.TB field is a *testing.T, so not a %T", t.TB)
	}
	conf := t.Config
	synctest.Test(tt, func(t *testing.T) {
		f(td.NewT(t, conf))
	})
}

// TestAssertRequire is a wrapper around [synctest.Test].
//
// The assert and require params of f inherit the configuration of t,
// except that a failure is never fatal using assert and always fatal
// using require.
//
// t.TB (set by [td.NewT], [td.Assert], [td.Require], …) must be a
// [*testing.T] instance.
func TestAssertRequire(t *td.T, f func(assert, require *td.T)) {
	tt, ok := t.TB.(*testing.T)
	if !ok {
		t.Helper()
		t.Fatalf("tdsynctest.TestAssertRequire only works if underlying T.TB field is a *testing.T, so not a %T", t.TB)
	}
	conf := t.Config
	synctest.Test(tt, func(t *testing.T) {
		f(td.AssertRequire(t, conf))
	})
}

// Wait simply calls [synctest.Wait]. It only exists to avoid to
// import [testing/synctest].
func Wait() {
	synctest.Wait()
}
