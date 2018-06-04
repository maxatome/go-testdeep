// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"
)

type tdLen struct {
	tdSmuggler
}

var _ TestDeep = &tdLen{}

// Len is a smuggler operator. It takes data, applies len() function
// on it and compares its result to "val". Of course, the compared
// value must be an array, a channel, a map, a slice or a string.
//
// "val" can be an int value:
//   Len(12)
// as well as an other operator:
//   Len(Between(3, 4))
func Len(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		l := tdLen{
			tdSmuggler: newSmuggler(val),
		}

		if l.isTestDeeper {
			return &l
		}

		// A len is always an int
		if vval.Type() == intType {
			l.expectedValue = vval
			return &l
		}
	}
	panic("usage: Len(TESTDEEP_OPERATOR|INT)")
}

func (l *tdLen) String() string {
	if l.isTestDeeper {
		return "len: " + l.expectedValue.Interface().(TestDeep).String()
	}
	return fmt.Sprintf("len=%d", l.expectedValue.Int())
}

func (l *tdLen) Match(ctx Context, got reflect.Value) *Error {
	switch got.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if l.isTestDeeper {
			vlen := reflect.New(intType)
			vlen.Elem().SetInt(int64(got.Len()))

			return deepValueEqual(ctx.AddFunctionCall("len"),
				vlen.Elem(), l.expectedValue)
		}

		if got.Len() == int(l.expectedValue.Int()) {
			return nil
		}
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "bad length",
			Got:      rawInt(got.Len()),
			Expected: rawInt(l.expectedValue.Int()),
		})

	default:
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: rawString("Array, Chan, Map, Slice or string"),
		})
	}
}

type tdCap struct {
	tdSmuggler
}

var _ TestDeep = &tdCap{}

// Cap is a smuggler operator. It takes data, applies cap() function
// on it and compares its result to "val". Of course, the compared
// value must be an array, a channel or a slice.
//
// "val" can be an int value:
//   Cap(12)
// as well as an other operator:
//   Cap(Between(3, 4))
func Cap(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		c := tdCap{
			tdSmuggler: newSmuggler(val),
		}

		if c.isTestDeeper {
			return &c
		}

		// A len is always an int
		if vval.Type() == intType {
			c.expectedValue = vval
			return &c
		}
	}
	panic("usage: Cap(TESTDEEP_OPERATOR|INT)")
}

func (c *tdCap) String() string {
	if c.isTestDeeper {
		return "cap: " + c.expectedValue.Interface().(TestDeep).String()
	}
	return fmt.Sprintf("cap=%d", c.expectedValue.Int())
}

func (c *tdCap) Match(ctx Context, got reflect.Value) *Error {
	switch got.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if c.isTestDeeper {
			vcap := reflect.New(intType)
			vcap.Elem().SetInt(int64(got.Cap()))

			return deepValueEqual(ctx.AddFunctionCall("cap"),
				vcap.Elem(), c.expectedValue)
		}

		if got.Cap() == int(c.expectedValue.Int()) {
			return nil
		}
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "bad capacity",
			Got:      rawInt(got.Cap()),
			Expected: rawInt(c.expectedValue.Int()),
		})

	default:
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: rawString("Array, Chan or Slice"),
		})
	}
}
