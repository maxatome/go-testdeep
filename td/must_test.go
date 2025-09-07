// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.18
// +build go1.18

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestMust(t *testing.T) {
	fn := func(isErr bool) (int64, error) {
		if isErr {
			return 0, fmt.Errorf("ERROR")
		}
		return 1234, nil
	}

	checkOK := func(t testing.TB, fn func(t testing.TB, isErr bool) int64) {
		ttt := test.NewTestingTB(t.Name())
		t.Helper()
		res := fn(ttt, false)
		td.CmpFalse(t, ttt.Failed())
		td.Cmp(t, res, int64(1234))
	}

	checkFailure := func(t testing.TB, fn func(t testing.TB, isErr bool) int64) {
		ttt := test.NewTestingTB(t.Name())
		t.Helper()
		td.CmpNotEmpty(t, ttt.CatchFatal(func() { fn(ttt, true) }))
		td.CmpTrue(t, ttt.Failed())
	}

	testCases := []struct {
		name string
		doit func(t testing.TB, isErr bool) int64
	}{
		{
			name: "old way",
			doit: func(t testing.TB, isErr bool) int64 {
				integer, err := fn(isErr)
				td.Require(t).CmpNoError(err)
				return integer
			},
		},
		{
			name: "Must1",
			doit: func(t testing.TB, isErr bool) int64 {
				return td.Must1[int64](t)(fn(isErr))
			},
		},
		{
			name: "Must2",
			doit: func(t testing.TB, isErr bool) int64 {
				return td.Must2(fn(isErr))(t)
			},
		},
		{
			name: "Must3+Call",
			doit: func(t testing.TB, isErr bool) int64 {
				return td.Must3(t, td.Call(fn(isErr)))
			},
		},
		{
			name: "Call+Must",
			doit: func(t testing.TB, isErr bool) int64 {
				return td.Call(fn(isErr)).Must(t)
			},
		},
		{
			name: "td.T.Must",
			doit: func(t testing.TB, isErr bool) int64 {
				return td.Assert(t).Must(fn(isErr)).(int64)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("OK", func(t *testing.T) {
				checkOK(t, tc.doit)
			})
			t.Run("Failure", func(t *testing.T) {
				checkFailure(t, tc.doit)
			})
		})
	}
}
