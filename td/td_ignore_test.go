// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestIgnore(t *testing.T) {
	checkOK(t, "any value!", td.Ignore())

	checkOK(t, nil, td.Ignore())
	checkOK(t, (*int)(nil), td.Ignore())

	//
	// String
	test.EqualStr(t, td.Ignore().String(), "Ignore()")
}

func TestIgnoreTypeBehind(t *testing.T) {
	equalTypes(t, td.Ignore(), nil)
}
