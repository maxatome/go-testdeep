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

type tdCode struct {
	Base
	function reflect.Value
	argType  reflect.Type
}

var _ TestDeep = &tdCode{}

// Code operator allows to check data using a custom function. So
// "fn" is a function that must take one parameter whose type must be
// the same as the type of the compared value.
//
// "fn" can return a single bool kind value, telling that yes or no
// the custom test is successful:
//   Code(func (date time.Time) bool {
//       return date.Year() == 2018
//     })
//
// or two values (bool, string) kinds. The bool value has the same
// meaning as above, and the string value is used to describe the
// test when it fails:
//   Code(func (date time.Time) (bool, string) {
//       if date.Year() == 2018 {
//         return true, ""
//       }
//       return false, "year must be 2018"
//     })
//
// This operator allows to handle any specific comparison not handled
// by standard operators.
//
// It is not recommended to call CmpDeeply (or any other Cmp*
// functions or *T methods) inside the body of "fn", because of
// confusion produced by output in case of failure. When the data
// needs to be transformed before being compared again, Smuggle
// operator should be used instead.
func Code(fn interface{}) TestDeep {
	vfn := reflect.ValueOf(fn)

	if vfn.Kind() != reflect.Func {
		panic("usage: Code(FUNC)")
	}

	fnType := vfn.Type()
	if fnType.NumIn() != 1 {
		panic("Code(FUNC): FUNC must take only one argument")
	}

	switch fnType.NumOut() {
	case 2:
		if fnType.Out(1).Kind() != reflect.String {
			break
		}
		fallthrough

	case 1:
		if fnType.Out(0).Kind() == reflect.Bool {
			return &tdCode{
				Base:     NewBase(3),
				function: vfn,
				argType:  fnType.In(0),
			}
		}
	}

	panic("Code(FUNC): FUNC must return bool or (bool, string)")
}

func (c *tdCode) Match(ctx Context, got reflect.Value) *Error {
	if !got.Type().AssignableTo(c.argType) {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "incompatible parameter type",
			Got:      rawString(got.Type().String()),
			Expected: rawString(c.argType.String()),
		})
	}

	// Refuse to override unexported fields access in this case. It is a
	// choice, as we think it is better to use Code() on surrounding
	// struct instead.
	if !got.CanInterface() {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message: "cannot compare unexported field",
			Summary: rawString("use Code() on surrounding struct instead"),
		})
	}

	ret := c.function.Call([]reflect.Value{got})
	if ret[0].Bool() {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}

	err := Error{
		Message: "ran code with %% as argument",
	}

	if len(ret) > 1 {
		err.Summary = tdCodeResult{
			Value:  got,
			Reason: ret[1].String(),
		}
	} else {
		err.Summary = tdCodeResult{
			Value: got,
		}
	}

	return ctx.CollectError(&err)
}

func (c *tdCode) String() string {
	return "Code(" + c.function.Type().String() + ")"
}

type tdCodeResult struct {
	Value  reflect.Value
	Reason string
}

var _ testDeepStringer = tdCodeResult{}

func (r tdCodeResult) _TestDeep() {}

func (r tdCodeResult) String() string {
	if r.Reason == "" {
		return fmt.Sprintf("  value: %s\nit failed but didn't say why",
			toString(r.Value))
	}
	return fmt.Sprintf("        value: %s\nit failed coz: %s",
		toString(r.Value), r.Reason)
}
