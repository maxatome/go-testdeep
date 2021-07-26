// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package nocolor_test

import (
	"os"
	"testing"

	_ "github.com/maxatome/go-testdeep/helpers/nocolor"
)

func TestNocolor(t *testing.T) {
	tdColor := os.Getenv("TESTDEEP_COLOR")
	if tdColor != "off" {
		t.Errorf(`TESTDEEP_COLOR expected to be "off" but is %q`, tdColor)
	}
}
