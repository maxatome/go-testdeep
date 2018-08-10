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

type tdMapEach struct {
	BaseOKNil
	expected reflect.Value
}

var _ TestDeep = &tdMapEach{}

// MapEach operator has to be applied on maps. It compares each value
// of data map against expected value. During a match, all values have
// to match to succeed.
func MapEach(expectedValue interface{}) TestDeep {
	return &tdMapEach{
		BaseOKNil: NewBaseOKNil(3),
		expected:  reflect.ValueOf(expectedValue),
	}
}

func (m *tdMapEach) Match(ctx Context, got reflect.Value) *Error {
	if !got.IsValid() {
		if ctx.BooleanError {
			return BooleanError
		}
		return ctx.CollectError(&Error{
			Message:  "nil value",
			Got:      rawString("nil"),
			Expected: rawString("Map OR *Map"),
		})
	}

	switch got.Kind() {
	case reflect.Ptr:
		gotElem := got.Elem()
		if !gotElem.IsValid() {
			if ctx.BooleanError {
				return BooleanError
			}
			return ctx.CollectError(&Error{
				Message:  "nil pointer",
				Got:      rawString("nil " + got.Type().String()),
				Expected: rawString("Map OR *Map"),
			})
		}

		if gotElem.Kind() != reflect.Map {
			break
		}
		got = gotElem
		fallthrough

	case reflect.Map:
		var err *Error
		for _, key := range got.MapKeys() {
			err = deepValueEqual(ctx.AddMapKey(key), got.MapIndex(key), m.expected)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if ctx.BooleanError {
		return BooleanError
	}
	return ctx.CollectError(&Error{
		Message:  "bad type",
		Got:      rawString(got.Type().String()),
		Expected: rawString("Map OR *Map"),
	})
}

func (m *tdMapEach) String() string {
	const prefix = "MapEach("

	content := toString(m.expected)
	if strings.Contains(content, "\n") {
		return prefix + indentString(content, "        ") + ")"
	}
	return prefix + content + ")"
}
