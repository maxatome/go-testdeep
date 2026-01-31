// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.18
// +build go1.18

package td_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestMust(t *testing.T) {
	fn := func(ok bool) (int, error) {
		if ok {
			return 42, nil
		}
		return 0, errors.New("error")
	}

	test.EqualInt(t, td.Must(fn(true)), 42)

	test.CheckPanic(t, func() { td.Must(fn(false)) }, "error")
}

func TestMust2(t *testing.T) {
	fn := func(ok bool) (int, string, error) {
		if ok {
			return 42, "pipo", nil
		}
		return 0, "", errors.New("error")
	}

	val1, val2 := td.Must2(fn(true))
	test.EqualInt(t, val1, 42)
	test.EqualStr(t, val2, "pipo")

	test.CheckPanic(t, func() { td.Must2(fn(false)) }, "error")
}

func TestMust3(t *testing.T) {
	fn := func(ok bool) (int, string, bool, error) {
		if ok {
			return 42, "pipo", true, nil
		}
		return 0, "", false, errors.New("error")
	}

	val1, val2, val3 := td.Must3(fn(true))
	test.EqualInt(t, val1, 42)
	test.EqualStr(t, val2, "pipo")
	test.EqualBool(t, val3, true)

	test.CheckPanic(t, func() { td.Must3(fn(false)) }, "error")
}
