// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build !go1.27
// +build !go1.27

package td

import "testing"

// Only useful for full coverage on go1.26, can be removed once go1.27
// is used to compute coverage.
func TestJoinOptions(t *testing.T) {
	if joinOptions(nil, nil) != nil {
		t.Fatal("joinOptions returned a non-nil value")
	}
}
