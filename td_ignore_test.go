// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestIgnore(t *testing.T) {
	checkOK(t, "any value!", Ignore())

	checkOK(t, nil, Ignore())
	checkOK(t, (*int)(nil), Ignore())

	//
	// String
	equalStr(t, Ignore().String(), "Ignore()")
}
