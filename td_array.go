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

type ArrayEntries map[int]interface{}

func Array(model interface{}, entries ArrayEntries) TestDeep {
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
		a.populateExpectedEntries(entries)
		return &a
	}

	panic("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES)")
}

func Slice(model interface{}, entries ArrayEntries) TestDeep {
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
		a.populateExpectedEntries(entries)
		return &a
	}

	panic("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES)")
}

func (a *tdArray) populateExpectedEntries(entries ArrayEntries) {
	var maxLength, numEntries int

	maxIndex := -1
	for index := range entries {
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
	for index, expectedValue := range entries {
		vexpectedValue := reflect.ValueOf(expectedValue)

		if _, ok := expectedValue.(TestDeep); !ok {
			if !vexpectedValue.Type().AssignableTo(elemType) {
				panic(fmt.Sprintf(
					"type %s of #%d expected value differs from %s content (%s)",
					vexpectedValue.Type(),
					index,
					ternStr(maxLength < 0, "slice", "array"),
					elemType))
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
		if _, ok := entries[index]; ok {
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
		if _, ok := entries[index]; !ok {
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

		err = deepValueEqual(got.Index(index), expectedValue, curCtx)
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
			fmt.Fprintf(buf, "  %d: %s\n", index, toString(expectedValue))
		}

		buf.WriteString("})")
	}
	return buf.String()
}

func (a *tdArray) expectedTypeStr() string {
	if a.isPtr {
		return "*" + a.expectedModel.Type().String()
	}
	return a.expectedModel.Type().String()
}
