// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.13
// +build go1.13

package json_test

import (
	"testing"
)

func TestJSON_go113(t *testing.T) {
	// Extend to golang 1.13 accepted numbers

	// as int64
	checkJSON(t, `4_2`, `42`)
	checkJSON(t, `+4_2`, `42`)
	checkJSON(t, `-4_2`, `-42`)

	checkJSON(t, `0b101010`, `42`)
	checkJSON(t, `-0b101010`, `-42`)
	checkJSON(t, `+0b101010`, `42`)

	checkJSON(t, `0b10_1010`, `42`)
	checkJSON(t, `-0b_10_1010`, `-42`)
	checkJSON(t, `+0b10_10_10`, `42`)

	checkJSON(t, `0B101010`, `42`)
	checkJSON(t, `-0B101010`, `-42`)
	checkJSON(t, `+0B101010`, `42`)

	checkJSON(t, `0B10_1010`, `42`)
	checkJSON(t, `-0B_10_1010`, `-42`)
	checkJSON(t, `+0B10_10_10`, `42`)

	checkJSON(t, `0_600`, `384`)
	checkJSON(t, `-0_600`, `-384`)
	checkJSON(t, `+0_600`, `384`)

	checkJSON(t, `0o600`, `384`)
	checkJSON(t, `0o_600`, `384`)
	checkJSON(t, `-0o600`, `-384`)
	checkJSON(t, `-0o6_00`, `-384`)
	checkJSON(t, `+0o600`, `384`)
	checkJSON(t, `+0o60_0`, `384`)

	checkJSON(t, `0O600`, `384`)
	checkJSON(t, `0O_600`, `384`)
	checkJSON(t, `-0O600`, `-384`)
	checkJSON(t, `-0O6_00`, `-384`)
	checkJSON(t, `+0O600`, `384`)
	checkJSON(t, `+0O60_0`, `384`)

	checkJSON(t, `0xBad_Face`, `195951310`)
	checkJSON(t, `-0x_Bad_Face`, `-195951310`)
	checkJSON(t, `+0xBad_Face`, `195951310`)

	checkJSON(t, `0XBad_Face`, `195951310`)
	checkJSON(t, `-0X_Bad_Face`, `-195951310`)
	checkJSON(t, `+0XBad_Face`, `195951310`)

	// as float64
	checkJSON(t, `0_600.123`, `600.123`) // float64 can not be an octal number
	checkJSON(t, `1_5.`, `15`)
	checkJSON(t, `0.15e+0_2`, `15`)
	checkJSON(t, `0x1p-2`, `0.25`)
	checkJSON(t, `0x2.p10`, `2048`)
	checkJSON(t, `0x1.Fp+0`, `1.9375`)
	checkJSON(t, `0X.8p-0`, `0.5`)
	checkJSON(t, `0X_1FFFP-16`, `0.1249847412109375`)
}
