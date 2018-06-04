// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdPtr struct {
	tdSmuggler
}

var _ TestDeep = &tdPtr{}

// Ptr is a smuggler operator. It takes the address of data and
// compares it to "val".
//
// "val" depends on data type. For example, if the compared data is an
// *int, one can have:
//   Ptr(12)
// as well as an other operator:
//   Ptr(Between(3, 4))
func Ptr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPtr{
			tdSmuggler: newSmuggler(val),
		}

		if !p.isTestDeeper {
			p.expectedValue = reflect.New(vval.Type())
			p.expectedValue.Elem().Set(vval)
		}
		return &p
	}
	panic("usage: Ptr(NON_NIL_VALUE)")
}

func (p *tdPtr) Match(ctx Context, got reflect.Value) *Error {
	if got.Kind() != reflect.Ptr {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "pointer type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(p.String()),
		})
	}

	if p.isTestDeeper {
		return deepValueEqual(ctx.AddPtr(1), got.Elem(), p.expectedValue)
	}
	return deepValueEqual(ctx, got, p.expectedValue)
}

func (p *tdPtr) String() string {
	if p.isTestDeeper {
		return "*<something>"
	}
	return p.expectedValue.Type().String()
}

type tdPPtr struct {
	tdSmuggler
}

var _ TestDeep = &tdPPtr{}

// PPtr is a smuggler operator. It takes the address of the address of
// data and compares it to "val".
//
// "val" depends on data type. For example, if the compared data is an
// **int, one can have:
//   PPtr(12)
// as well as an other operator:
//   PPtr(Between(3, 4))
//
// It is more efficient and shorter to write than:
//   Ptr(Ptr(val))
func PPtr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPPtr{
			tdSmuggler: newSmuggler(val),
		}

		if !p.isTestDeeper {
			pVval := reflect.New(vval.Type())
			pVval.Elem().Set(vval)

			p.expectedValue = reflect.New(pVval.Type())
			p.expectedValue.Elem().Set(pVval)
		}
		return &p
	}
	panic("usage: PPtr(NON_NIL_VALUE)")
}

func (p *tdPPtr) Match(ctx Context, got reflect.Value) *Error {
	if got.Kind() != reflect.Ptr || got.Elem().Kind() != reflect.Ptr {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "pointer type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(p.String()),
		})
	}

	if p.isTestDeeper {
		return deepValueEqual(ctx.AddPtr(2), got.Elem().Elem(), p.expectedValue)
	}
	return deepValueEqual(ctx, got, p.expectedValue)
}

func (p *tdPPtr) String() string {
	if p.isTestDeeper {
		return "**<something>"
	}
	return p.expectedValue.Type().String()
}
