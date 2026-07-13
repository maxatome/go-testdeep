// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td

// Must fails immediately if err is non-nil. Otherwise it returns ret.
//
//	assert := td.Assert(t)
//	fn := func() (int, error) { … }
//	value := assert.Must2(fn())
//
// It typically avoids to use 2 lines for the same effect than:
//
//	value, err := fn()
//	td.Require(t).CmpNoError(err)
//
// See also [T.Must2], [T.Must3], [T.CmpNoError], [Must], [Must2],
// [Must3], and [CmpNoError].
func (t *T) Must[X any](ret X, err error) X {
	t.Helper()
	t.Require().CmpNoError(err)
	return ret
}

// Must2 fails immediately if err is non-nil. Otherwise it returns ret1, ret2.
//
//	assert := td.Assert(t)
//	fn := func() (int, string, error) { … }
//	value1, value2 := assert.Must2(fn())
//
// It typically avoids to use 2 lines for the same effect than:
//
//	value1, value2, err := fn()
//	td.Require(t).CmpNoError(err)
//
// See also [T.Must], [T.Must3], [T.CmpNoError], [Must], [Must2],
// [Must3] and [CmpNoError].
func (t *T) Must2[X, Y any](ret1 X, ret2 Y, err error) (X, Y) {
	t.Helper()
	t.Require().CmpNoError(err)
	return ret1, ret2
}

// Must3 fails immediately if err is non-nil. Otherwise it returns
// ret1, ret2, ret3.
//
//	assert := td.Assert(t)
//	fn := func() (int, string, bool, error) { … }
//	value1, value2, value3 := assert.Must3(fn())
//
// It typically avoids to use 2 lines for the same effect than:
//
//	value1, value2, value3, err := fn()
//	td.Require(t).CmpNoError(err)
//
// See also [T.Must], [T.Must2], [T.CmpNoError], [Must], [Must2],
// [Must3] and [CmpNoError].
func (t *T) Must3[X, Y, Z any](ret1 X, ret2 Y, ret3 Z, err error) (X, Y, Z) {
	t.Helper()
	t.Require().CmpNoError(err)
	return ret1, ret2, ret3
}
