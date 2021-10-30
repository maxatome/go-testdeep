// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package dark

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestGetInterface(t *testing.T) {
	type Private struct {
		private *Private //nolint: megacheck
		privInt int
	}

	// Cases not tested by TestEqualOthers()
	s := Private{
		private: &Private{},
		privInt: 42,
	}

	//
	// GetInterface
	val, ok := GetInterface(reflect.ValueOf(nil), false)
	if val != nil {
		test.EqualErrorMessage(t, val, "nil")
	}
	test.IsTrue(t, ok)

	val, ok = GetInterface(reflect.ValueOf(123), false)
	if test.IsTrue(t, ok) {
		valInt, ok := val.(int)
		if test.IsTrue(t, ok) {
			test.EqualInt(t, valInt, 123)
		}
	}

	_, ok = GetInterface(reflect.ValueOf(s).Field(0), false)
	test.IsFalse(t, ok, "private field")

	val, ok = GetInterface(reflect.ValueOf(s).Field(1), false)
	if test.IsTrue(t, ok, "private field, BUT contents can be copied") {
		valInt, ok := val.(int)
		if test.IsTrue(t, ok) {
			test.EqualInt(t, valInt, 42)
		}
	}

	_, ok = GetInterface(reflect.ValueOf(s).Field(0), true)
	if UnsafeDisabled {
		test.IsFalse(t, ok, "unsafe package is disabled, GetInterface should fail")
	} else {
		test.IsTrue(t, ok,
			"unsafe package is available, GetInterface should succeed")
	}

	//
	// MustGetInterface
	val = MustGetInterface(reflect.ValueOf(123))
	valInt, ok := val.(int)
	if test.IsTrue(t, ok) {
		test.EqualInt(t, valInt, 123)
	}

	if UnsafeDisabled {
		test.CheckPanic(t,
			func() {
				MustGetInterface(reflect.ValueOf(s).Field(0))
			},
			"dark.GetInterface() does not handle private")
	} else {
		val = MustGetInterface(reflect.ValueOf(s).Field(0))
		if val == nil {
			test.EqualErrorMessage(t, val, "non-nil")
		}
	}

	// Private field BUT contents can be copied
	val = MustGetInterface(reflect.ValueOf(s).Field(1))
	valInt, ok = val.(int)
	if test.IsTrue(t, ok) {
		test.EqualInt(t, valInt, 42)
	}
}
