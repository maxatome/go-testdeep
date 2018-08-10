// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdExpectedType struct {
	Base
	expectedType reflect.Type
	isPtr        bool
}

func (t *tdExpectedType) errorTypeMismatch(gotType rawString) *Error {
	return &Error{
		Message:  "type mismatch",
		Got:      gotType,
		Expected: rawString(t.expectedTypeStr()),
	}
}

func (t *tdExpectedType) checkPtr(ctx Context, pGot *reflect.Value, nilAllowed bool) *Error {
	if t.isPtr {
		got := *pGot
		if got.Kind() != reflect.Ptr {
			if ctx.BooleanError {
				return BooleanError
			}
			return t.errorTypeMismatch(rawString(got.Type().String()))
		}

		if !nilAllowed && got.IsNil() {
			if ctx.BooleanError {
				return BooleanError
			}
			return &Error{
				Message:  "values differ",
				Got:      got,
				Expected: rawString("non-nil"),
			}
		}

		*pGot = got.Elem()
	}
	return nil
}

func (t *tdExpectedType) checkType(ctx Context, got reflect.Value) *Error {
	if got.Type() != t.expectedType {
		if ctx.BooleanError {
			return BooleanError
		}
		var gotType rawString
		if t.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return t.errorTypeMismatch(gotType)
	}
	return nil
}

func (t *tdExpectedType) TypeBehind() reflect.Type {
	if t.isPtr {
		return reflect.New(t.expectedType).Type()
	}
	return t.expectedType
}

func (t *tdExpectedType) expectedTypeStr() string {
	if t.isPtr {
		return "*" + t.expectedType.String()
	}
	return t.expectedType.String()
}
