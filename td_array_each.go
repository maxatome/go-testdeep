// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
	"strings"
)

type tdArrayEach struct {
	BaseOKNil
	expected reflect.Value
}

var _ TestDeep = &tdArrayEach{}

// ArrayEach operator has to be applied on arrays or slices or on
// pointers on array/slice. It compares each item of data array/slice
// against expected value. During a match, all items have to match to
// succeed.
func ArrayEach(expectedValue interface{}) TestDeep {
	return &tdArrayEach{
		BaseOKNil: NewBaseOKNil(3),
		expected:  reflect.ValueOf(expectedValue),
	}
}

func (a *tdArrayEach) Match(ctx Context, got reflect.Value) (err *Error) {
	if !got.IsValid() {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "nil value",
			Got:      rawString("nil"),
			Expected: rawString("Slice OR Array OR *Slice OR *Array"),
			Location: a.GetLocation(),
		}
	}

	switch got.Kind() {
	case reflect.Ptr:
		gotElem := got.Elem()
		if !gotElem.IsValid() {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil pointer",
				Got:      rawString("nil " + got.Type().String()),
				Expected: rawString("Slice OR Array OR *Slice OR *Array"),
				Location: a.GetLocation(),
			}
		}

		if gotElem.Kind() != reflect.Array && gotElem.Kind() != reflect.Slice {
			break
		}
		got = gotElem
		fallthrough

	case reflect.Array, reflect.Slice:
		gotLen := got.Len()

		for idx := 0; idx < gotLen; idx++ {
			err = deepValueEqual(ctx.AddArrayIndex(idx), got.Index(idx), a.expected)
			if err != nil {
				return err.SetLocationIfMissing(a)
			}
		}
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "bad type",
		Got:      rawString(got.Type().String()),
		Expected: rawString("Slice OR Array OR *Slice OR *Array"),
		Location: a.GetLocation(),
	}
}

func (a *tdArrayEach) String() string {
	const prefix = "ArrayEach("

	content := toString(a.expected)
	if strings.Contains(content, "\n") {
		return prefix + indentString(content, "          ") + ")"
	}
	return prefix + content + ")"
}
