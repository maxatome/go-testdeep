// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdNil struct {
	BaseOKNil
}

var _ TestDeep = &tdNil{}

// Nil operator checks that data is nil (or is a non-nil interface,
// but containing a nil pointer.)
func Nil() TestDeep {
	return &tdNil{
		BaseOKNil: NewBaseOKNil(3),
	}
}

func (n *tdNil) Match(ctx Context, got reflect.Value) *Error {
	if !got.IsValid() {
		return nil
	}

	switch got.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		if got.IsNil() {
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return ctx.CollectError(&Error{
		Message:  "non-nil",
		Got:      got,
		Expected: n,
	})
}

func (n *tdNil) String() string {
	return "nil"
}

type tdNotNil struct {
	BaseOKNil
}

var _ TestDeep = &tdNotNil{}

// NotNil operator checks that data is not nil (or is a non-nil
// interface, containing a non-nil pointer.)
func NotNil() TestDeep {
	return &tdNotNil{
		BaseOKNil: NewBaseOKNil(3),
	}
}

func (n *tdNotNil) Match(ctx Context, got reflect.Value) *Error {
	if got.IsValid() {
		switch got.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			if !got.IsNil() {
				return nil
			}

			// All other kinds are non-nil by nature
		default:
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return ctx.CollectError(&Error{
		Message:  "nil value",
		Got:      got,
		Expected: n,
	})
}

func (n *tdNotNil) String() string {
	return "not nil"
}
