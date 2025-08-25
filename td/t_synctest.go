// Copyright (c) 2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build !go1.25
// +build !go1.25

package td

// SyncTest is only supported from go1.25.0.
func (t *T) SyncTest(func(t *T)) {
	t.Fatal("SyncTest is only supported from go1.25.0")
}

// SyncTestAssertRequire is only supported from go1.25.0.
func (t *T) SyncTestAssertRequire(func(assert, require *T)) {
	t.Fatal("SyncTestAssertRequire is only supported from go1.25.0")
}
