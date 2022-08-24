// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build !js && !appengine && !safe && !disableunsafe
// +build !js,!appengine,!safe,!disableunsafe

package td_test

import (
	"testing"
)

// Map, unsafe access is mandatory here.
func TestEqualMapUnsafe(t *testing.T) {
	type key struct{ k string }
	type A struct{ x map[key]struct{} }

	checkError(t, A{x: map[key]struct{}{{k: "z"}: {}}},
		A{x: map[key]struct{}{{k: "x"}: {}}},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA.x"),
			Summary: mustBe(`Missing key: ((td_test.key) {
               k: (string) (len=1) "x"
              })
  Extra key: ((td_test.key) {
               k: (string) (len=1) "z"
              })`),
		})
}
