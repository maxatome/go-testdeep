// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.25
// +build go1.25

package td

import (
	"testing"
	"testing/synctest"
)

// SyncTest is a wrapper around [synctest.Test].
//
// The t param of f inherits the configuration of the self-reference.
//
// t.TB (set by [NewT], [Assert], [Require], …) must be a [*testing.T]
// instance.
func (t *T) SyncTest(f func(t *T)) {
	tt, ok := t.TB.(*testing.T)
	if !ok {
		t.Helper()
		t.Fatalf("SyncTest only works if underlying T.TB field is a *testing.T, so not a %T", t.TB)
	}
	conf := t.Config
	synctest.Test(tt, func(t *testing.T) {
		f(NewT(t, conf))
	})
}

// SyncTestAssertRequire is a wrapper around [synctest.Test].
//
// The assert and require params of f inherit the configuration
// of the self-reference, except that a failure is never fatal using
// assert and always fatal using require.
//
// t.TB (set by [NewT], [Assert], [Require], …) must be a [*testing.T]
// instance.
func (t *T) SyncTestAssertRequire(f func(assert, require *T)) {
	tt, ok := t.TB.(*testing.T)
	if !ok {
		t.Helper()
		t.Fatalf("SyncTestAssertRequire only works if underlying T.TB field is a *testing.T, so not a %T", t.TB)
	}
	conf := t.Config
	synctest.Test(tt, func(t *testing.T) {
		f(AssertRequire(t, conf))
	})
}
