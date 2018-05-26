// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"testing"
)

// Edge cases not tested elsewhere...

func TestBase(t *testing.T) {
	td := Base{}

	td.setLocation(200)
	if td.location.File != "???" && td.location.Line != 0 {
		t.Errorf("Location found! => %s", td.location)
	}
}

func TestTdSetResult(t *testing.T) {
	if tdSetResultKind(199).String() != "?" {
		t.Errorf("tdSetResultKind stringification failed => %s",
			tdSetResultKind(199))
	}
}
