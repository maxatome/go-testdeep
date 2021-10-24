// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.15
// +build go1.15

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestSmuggleFieldsPath_go115(t *testing.T) {
	type C struct {
		Iface interface{}
	}

	got := C{
		Iface: []interface{}{
			map[complex64]interface{}{complex(42, 0): []string{"pipo"}},
			map[complex128]interface{}{complex(42, 0): []string{"pipo"}},
		},
	}

	for i := 0; i < 2; i++ {
		checkOK(t, got, td.Smuggle(fmt.Sprintf("Iface[%d][42][0]", i), "pipo"))
		checkOK(t, got, td.Smuggle(fmt.Sprintf("Iface[%d][42][0]", i-2), "pipo"))
	}
}
