// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"testing"
)

// CmpTrue is a shortcut for:
//
//   CmpDeeply(t, got, true, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpTrue(t *testing.T, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, true, args...)
}

// CmpFalse is a shortcut for:
//
//   CmpDeeply(t, got, false, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpFalse(t *testing.T, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, false, args...)
}
