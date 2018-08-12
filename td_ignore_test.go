// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestIgnore(t *testing.T) {
	checkOK(t, "any value!", testdeep.Ignore())

	checkOK(t, nil, testdeep.Ignore())
	checkOK(t, (*int)(nil), testdeep.Ignore())

	//
	// String
	test.EqualStr(t, testdeep.Ignore().String(), "Ignore()")
}

func TestIgnoreTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Ignore(), nil)
}
