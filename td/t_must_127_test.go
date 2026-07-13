// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestTMust(t *testing.T) {
	fn := func(ok bool) (int, error) {
		if ok {
			return 42, nil
		}
		return 0, errors.New("error")
	}

	ttt := test.NewTestingTB(t.Name())
	assert := td.NewT(ttt)

	test.EqualInt(t, assert.Must(fn(true)), 42)

	test.MatchStr(t, ttt.CatchFatal(func() { assert.Must(fn(false)) }),
		`Failed test
DATA: should NOT be an error
	     got: \S+\(error\)
	expected: nil
`)
}

func TestTMust2(t *testing.T) {
	fn := func(ok bool) (int, string, error) {
		if ok {
			return 42, "pipo", nil
		}
		return 0, "", errors.New("error")
	}

	ttt := test.NewTestingTB(t.Name())
	assert := td.NewT(ttt)

	val1, val2 := assert.Must2(fn(true))
	test.EqualInt(t, val1, 42)
	test.EqualStr(t, val2, "pipo")

	test.MatchStr(t, ttt.CatchFatal(func() { assert.Must2(fn(false)) }),
		`Failed test
DATA: should NOT be an error
	     got: \S+\(error\)
	expected: nil
`)
}

func TestTMust3(t *testing.T) {
	fn := func(ok bool) (int, string, bool, error) {
		if ok {
			return 42, "pipo", true, nil
		}
		return 0, "", false, errors.New("error")
	}

	ttt := test.NewTestingTB(t.Name())
	assert := td.NewT(ttt)

	val1, val2, val3 := assert.Must3(fn(true))
	test.EqualInt(t, val1, 42)
	test.EqualStr(t, val2, "pipo")
	test.IsTrue(t, val3)

	test.MatchStr(t, ttt.CatchFatal(func() { assert.Must3(fn(false)) }),
		`Failed test
DATA: should NOT be an error
	     got: \S+\(error\)
	expected: nil
`)
}
