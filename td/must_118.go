// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.18
// +build go1.18

package td

// Must panics if err is non-nil. Otherwise it returns ret.
//
//	fn := func() (int, error) { … }
//	value := td.Must(fn())
//
// It typically avoids to use 2 lines for the mostly same effect than:
//
//	value, err := fn()
//	td.Require(t).CmpNoError(err)
//
// except that in case of error a panic occurs instead of a clean
// [testing.T.Fatal] error like. Starting go1.27, you should prefer to
// use [T.Must] which does not panic but calls [testing.T.Fatal].
//
// See also [T.Must], [T.Must2], [T.Must3], [Must2], [Must3],
// [CmpNoError] and [T.CmpNoError].
func Must[X any](ret X, err error) X {
	if err != nil {
		panic("Must: " + err.Error())
	}
	return ret
}

// Must2 panics if err is non-nil. Otherwise it returns ret1, ret2.
//
//	fn := func() (int, string, error) { … }
//	value1, value2 := td.Must2(fn())
//
// It typically avoids to use 2 lines for the mostly same effect than:
//
//	value1, value2, err := fn()
//	td.Require(t).CmpNoError(err)
//
// except that in case of error a panic occurs instead of a clean
// [testing.T.Fatal] error. Starting go1.27, you should prefer to
// use [T.Must2] which does not panic but calls [testing.T.Fatal].
//
// See also [T.Must2], [T.Must], [T.Must3], [Must], [Must3],
// [CmpNoError] and [T.CmpNoError].
func Must2[X, Y any](ret1 X, ret2 Y, err error) (X, Y) {
	if err != nil {
		panic("Must2: " + err.Error())
	}
	return ret1, ret2
}

// Must3 panics if err is non-nil. Otherwise it returns ret1, ret2, ret3.
//
//	fn := func() (int, string, bool, error) { … }
//	value1, value2, value3 := td.Must3(fn())
//
// It typically avoids to use 2 lines for the mostly same effect than:
//
//	value1, value2, value3, err := fn()
//	td.Require(t).CmpNoError(err)
//
// except that in case of error a panic occurs instead of a clean
// [testing.T.Fatal] error. Starting go1.27, you should prefer to
// use [T.Must3] which does not panic but calls [testing.T.Fatal].
//
// See also [T.Must3], [T.Must], [T.Must2], [Must], [Must2],
// [CmpNoError] and [T.CmpNoError].
func Must3[X, Y, Z any](ret1 X, ret2 Y, ret3 Z, err error) (X, Y, Z) {
	if err != nil {
		panic("Must3: " + err.Error())
	}
	return ret1, ret2, ret3
}
