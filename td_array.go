// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"fmt"
	"reflect"
)

type tdArray struct {
	Base
	expectedModel   reflect.Value
	expectedEntries []reflect.Value
	isPtr           bool
}

var _ TestDeep = &tdArray{}

// ArrayEntries allows to pass array or slice entries to check in
// functions Array and Slice. It is a map whose each key is the item
// index and the corresponding value the expected item value (which
// can be a TestDeep operator as well as a zero value.)
type ArrayEntries map[int]interface{}

// Array operator compares the contents of an array or a pointer on an
// array against the non-zero values of "model" (if any) and the
// values of "expectedEntries".
//
// "model" must be the same type as compared data.
//
// "expectedEntries" can be nil, if no zero entries are expected and
// no TestDeep operator are involved.
//
// TypeBehind method returns the reflect.Type of "model".
func Array(model interface{}, expectedEntries ArrayEntries) TestDeep {
	vmodel := reflect.ValueOf(model)

	a := tdArray{
		Base: NewBase(3),
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Array {
			break
		}
		a.isPtr = true
		fallthrough

	case reflect.Array:
		a.expectedModel = vmodel
		a.populateExpectedEntries(expectedEntries)
		return &a
	}

	panic("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES)")
}

// Slice operator compares the contents of a slice or a pointer on a
// slice against the non-zero values of "model" (if any) and the
// values of "expectedEntries".
//
// "model" must be the same type as compared data.
//
// "expectedEntries" can be nil, if no zero entries are expected and
// no TestDeep operator are involved.
//
// TypeBehind method returns the reflect.Type of "model".
func Slice(model interface{}, expectedEntries ArrayEntries) TestDeep {
	vmodel := reflect.ValueOf(model)

	a := tdArray{
		Base: NewBase(3),
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Slice {
			break
		}
		a.isPtr = true
		fallthrough

	case reflect.Slice:
		a.expectedModel = vmodel
		a.populateExpectedEntries(expectedEntries)
		return &a
	}

	panic("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES)")
}

func (a *tdArray) populateExpectedEntries(expectedEntries ArrayEntries) {
	var maxLength, numEntries int

	maxIndex := -1
	for index := range expectedEntries {
		if index > maxIndex {
			maxIndex = index
		}
	}

	if a.expectedModel.Kind() == reflect.Array {
		maxLength = a.expectedModel.Len()

		if maxLength <= maxIndex {
			panic(fmt.Sprintf(
				"array length is %d, so cannot have #%d expected index",
				maxLength,
				maxIndex))
		}
		numEntries = maxLength
	} else {
		maxLength = -1

		numEntries = maxIndex + 1
		if numEntries < a.expectedModel.Len() {
			numEntries = a.expectedModel.Len()
		}
	}

	a.expectedEntries = make([]reflect.Value, numEntries)

	elemType := a.expectedModel.Type().Elem()
	var vexpectedValue reflect.Value
	for index, expectedValue := range expectedEntries {
		if expectedValue == nil {
			switch elemType.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
				reflect.Ptr, reflect.Slice:
				vexpectedValue = reflect.New(elemType).Elem() // change to a typed nil
			default:
				panic(fmt.Sprintf(
					"expected value of #%d cannot be nil as items type is %s",
					index,
					elemType))
			}
		} else {
			vexpectedValue = reflect.ValueOf(expectedValue)

			if _, ok := expectedValue.(TestDeep); !ok {
				if !vexpectedValue.Type().AssignableTo(elemType) {
					panic(fmt.Sprintf(
						"type %s of #%d expected value differs from %s contents (%s)",
						vexpectedValue.Type(),
						index,
						ternStr(maxLength < 0, "slice", "array"),
						elemType))
				}
			}
		}

		a.expectedEntries[index] = vexpectedValue
	}

	// Check initialized entries in model
	vzero := reflect.Zero(elemType)
	zero := vzero.Interface()
	for index := a.expectedModel.Len() - 1; index >= 0; index-- {
		ventry := a.expectedModel.Index(index)

		// Entry already expected
		if _, ok := expectedEntries[index]; ok {
			// If non-zero entry, consider it as an error (= 2 expected
			// values for the same item)
			if !reflect.DeepEqual(zero, ventry.Interface()) {
				panic(fmt.Sprintf(
					"non zero #%d entry in model already exists in expectedEntries",
					index))
			}
			continue
		}

		a.expectedEntries[index] = ventry
	}

	// Array case, all is OK
	if maxLength >= 0 {
		return
	}

	// Slice case, initialize missing expected items to zero
	for index := a.expectedModel.Len(); index < numEntries; index++ {
		if _, ok := expectedEntries[index]; !ok {
			a.expectedEntries[index] = vzero
		}
	}
}

func (a *tdArray) Match(ctx Context, got reflect.Value) (err *Error) {
	if a.isPtr {
		if got.Kind() != reflect.Ptr {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "type mismatch",
				Got:      rawString(got.Type().String()),
				Expected: rawString(a.expectedTypeStr()),
				Location: a.GetLocation(),
			}
		}
		got = got.Elem()
	}

	if got.Type() != a.expectedModel.Type() {
		if ctx.booleanError {
			return booleanError
		}
		var gotType rawString
		if a.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      gotType,
			Expected: rawString(a.expectedTypeStr()),
			Location: a.GetLocation(),
		}
	}

	gotLen := got.Len()
	for index, expectedValue := range a.expectedEntries {
		curCtx := ctx.AddArrayIndex(index)

		if index >= gotLen {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  curCtx,
				Message:  "expected value out of range",
				Got:      rawString("<non-existent value>"),
				Expected: expectedValue,
				Location: a.GetLocation(),
			}
		}

		err = deepValueEqual(curCtx, got.Index(index), expectedValue)
		if err != nil {
			return err.SetLocationIfMissing(a)
		}
	}

	if gotLen > len(a.expectedEntries) {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx.AddArrayIndex(len(a.expectedEntries)),
			Message:  "got value out of range",
			Got:      got.Index(len(a.expectedEntries)),
			Expected: rawString("<non-existent value>"),
			Location: a.GetLocation(),
		}
	}

	return nil
}

func (a *tdArray) String() string {
	buf := bytes.NewBufferString(ternStr(a.expectedModel.Kind() == reflect.Array,
		"Array(", "Slice("))

	buf.WriteString(a.expectedTypeStr())

	if len(a.expectedEntries) == 0 {
		buf.WriteString("{})")
	} else {
		buf.WriteString("{\n")

		for index, expectedValue := range a.expectedEntries {
			fmt.Fprintf(buf, "  %d: %s\n", // nolint: errcheck
				index, toString(expectedValue))
		}

		buf.WriteString("})")
	}
	return buf.String()
}

func (s *tdArray) TypeBehind() reflect.Type {
	if s.isPtr {
		return reflect.New(s.expectedModel.Type()).Type()
	}
	return s.expectedModel.Type()
}

func (a *tdArray) expectedTypeStr() string {
	if a.isPtr {
		return "*" + a.expectedModel.Type().String()
	}
	return a.expectedModel.Type().String()
}
