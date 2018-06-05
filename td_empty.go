// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

const emptyBadType rawString = "Array, Chan, Map, Slice, string or pointer(s) on them"

type tdEmpty struct {
	BaseOKNil
}

var _ TestDeep = &tdEmpty{}

// Empty operator checks that an array, a channel, a map, a slice or a
// string is empty. As a special case (non-typed) nil, as well as nil
// channel, map or slice are considered empty.
//
// Note that the compared data can be a pointer (of pointer of pointer
// etc.) on an array, a channel, a map, a slice or a string.
func Empty() TestDeep {
	return &tdEmpty{
		BaseOKNil: NewBaseOKNil(3),
	}
}

// isEmpty returns (isEmpty, typeError) boolean values with only 3
// possible cases:
//  - true, false  → "got" is empty
//  - false, false → "got" is not empty
//  - false, true  → "got" type is not compatible with emptiness
func isEmpty(got reflect.Value) (bool, bool) {
	switch got.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return got.Len() == 0, false

	case reflect.Ptr:
		switch got.Type().Elem().Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice,
			reflect.String:
			if got.IsNil() {
				return true, false
			}
			fallthrough
		case reflect.Ptr:
			return isEmpty(got.Elem())
		default:
			return false, true // bad type
		}

	default:
		// nil case
		if !got.IsValid() {
			return true, false
		}
		return false, true // bad type
	}
}

func (e *tdEmpty) Match(ctx Context, got reflect.Value) (err *Error) {
	ok, badType := isEmpty(got)
	if ok {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}

	if badType {
		return ctx.CollectError(&Error{
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: emptyBadType,
		})
	}

	return ctx.CollectError(&Error{
		Message:  "not empty",
		Got:      got,
		Expected: rawString("empty"),
	})
}

func (e *tdEmpty) String() string {
	return "Empty()"
}

type tdNotEmpty struct {
	BaseOKNil
}

var _ TestDeep = &tdNotEmpty{}

// NotEmpty operator checks that an array, a channel, a map, a slice
// or a string is not empty. As a special case (non-typed) nil, as
// well as nil channel, map or slice are considered empty.
//
// Note that the compared data can be a pointer (of pointer of pointer
// etc.) on an array, a channel, a map, a slice or a string.
func NotEmpty() TestDeep {
	return &tdNotEmpty{
		BaseOKNil: NewBaseOKNil(3),
	}
}

func (e *tdNotEmpty) Match(ctx Context, got reflect.Value) (err *Error) {
	ok, badType := isEmpty(got)
	if ok {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "empty",
			Got:      got,
			Expected: rawString("not empty"),
		})
	}

	if badType {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: emptyBadType,
		})
	}
	return nil
}

func (e *tdNotEmpty) String() string {
	return "NotEmpty()"
}