// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.25
// +build go1.25

// Until go 1.23 in go.mod
//go:debug asynctimerchan=0

package td_test

import (
	"testing"
	"testing/synctest"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestSyncTest(t *testing.T) {
	t.Run("testing.T required", func(t *testing.T) {
		tt := tdutil.NewT(t.Name())

		run := false
		panicked := tt.CatchFailNow(func() {
			assert := td.Assert(tt)
			assert.SyncTest(func(assert *td.T) {
				run = true
			})
		})
		test.IsFalse(t, run)
		test.IsTrue(t, panicked)
		test.MatchStr(t, tt.LogBuf(),
			`^\s+t_synctest_test\.go:\d+: SyncTest only works if underlying T\.TB field is a \*testing\.T, so not a \*tdutil\.T\n\z`)
	})

	t.Run("OK1", func(t *testing.T) {
		run, belax, req := false, false, true
		assert := td.Assert(t).BeLax()
		assert.SyncTest(func(assert *td.T) {
			go func() {
				run = true
				belax = assert.Config.BeLax
				req = assert.Config.FailureIsFatal
			}()
			synctest.Wait()
		})
		test.IsTrue(t, run)
		test.IsTrue(t, belax)
		test.IsFalse(t, req)
	})

	t.Run("OK2", func(t *testing.T) {
		run, belax, req := false, true, false
		require := td.Require(t)
		require.SyncTest(func(require *td.T) {
			go func() {
				run = true
				belax = require.Config.BeLax
				req = require.Config.FailureIsFatal
			}()
			synctest.Wait()
		})
		test.IsTrue(t, run)
		test.IsFalse(t, belax)
		test.IsTrue(t, req)
	})
}

func TestSyncTestAssertRequire(t *testing.T) {
	t.Run("testing.T required", func(t *testing.T) {
		tt := tdutil.NewT(t.Name())

		run := false
		panicked := tt.CatchFailNow(func() {
			assert := td.Assert(tt)
			assert.SyncTestAssertRequire(func(assert, require *td.T) {
				run = true
			})
		})
		test.IsFalse(t, run)
		test.IsTrue(t, panicked)
		test.MatchStr(t, tt.LogBuf(),
			`^\s+t_synctest_test\.go:\d+: SyncTestAssertRequire only works if underlying T\.TB field is a \*testing\.T, so not a \*tdutil\.T\n\z`)
	})

	t.Run("OK1", func(t *testing.T) {
		var run, belaxA, belaxR, ass, req bool
		assert := td.Assert(t).BeLax()
		assert.SyncTestAssertRequire(func(assert, require *td.T) {
			go func() {
				run = true
				belaxA = assert.Config.BeLax
				ass = !assert.Config.FailureIsFatal
				belaxR = require.Config.BeLax
				req = require.Config.FailureIsFatal
			}()
			synctest.Wait()
		})
		test.IsTrue(t, run)
		test.IsTrue(t, belaxA)
		test.IsTrue(t, ass)
		test.IsTrue(t, belaxR)
		test.IsTrue(t, req)
	})

	t.Run("OK2", func(t *testing.T) {
		run, belaxA, belaxR, ass, req := false, true, true, false, false
		assert := td.Assert(t)
		assert.SyncTestAssertRequire(func(assert, require *td.T) {
			go func() {
				run = true
				belaxA = assert.Config.BeLax
				ass = !assert.Config.FailureIsFatal
				belaxR = require.Config.BeLax
				req = require.Config.FailureIsFatal
			}()
			synctest.Wait()
		})
		test.IsTrue(t, run)
		test.IsFalse(t, belaxA)
		test.IsTrue(t, ass)
		test.IsFalse(t, belaxR)
		test.IsTrue(t, req)
	})
}
