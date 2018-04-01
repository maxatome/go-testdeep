package testdeep

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func isNilStr(isNil bool) rawString {
	if isNil {
		return "nil"
	}
	return "not nil"
}

func deepValueEqual(got, expected reflect.Value, ctx Context) (err *Error) {
	if !got.IsValid() || !expected.IsValid() {
		if got.IsValid() == expected.IsValid() {
			return
		}

		err = &Error{}

		if expected.IsValid() { // here: !got.IsValid()
			if expected.Type().Implements(testDeeper) {
				td := expected.Interface().(TestDeep)
				if td.HandleInvalid() {
					return td.Match(ctx, got)
				}
				if ctx.booleanError {
					return booleanError
				}

				// Special case if "expected" is a TestDeep operator which
				// does not handle invalid values: the operator is not called,
				// but for the user the error comes from it
				err.Location = td.GetLocation()
			} else if ctx.booleanError {
				return booleanError
			}

			err.Expected = expected
		} else { // here: !expected.IsValid() && got.IsValid()
			// Special case: "got" is a nil interface, so consider as equal
			// to "expected" nil.
			if got.Kind() == reflect.Interface && got.IsNil() {
				return nil
			}

			if ctx.booleanError {
				return booleanError
			}

			err.Got = got
		}

		err.Context = ctx
		err.Message = "values differ"
		return
	}

	// Special case, "got" implements testDeeper: only if allowed
	if got.Type().Implements(testDeeper) {
		panic("Found a TestDeep operator in got param, " +
			"can only use it in expected one!")
	}

	if got.Type() != expected.Type() {
		if expected.Type().Implements(testDeeper) {
			return expected.Interface().(TestDeep).Match(ctx, got)
		}

		// "expected" is not a TestDeep operator

		// If "got" is an interface, try to see what is behind before failing
		// Used by Set/Bag Match method in such cases:
		//     []interface{}{123, "foo"}  →  Bag("foo", 123)
		//    Interface kind -^-----^   but String-^ and ^- Int kinds
		if got.Kind() == reflect.Interface {
			return deepValueEqual(got.Elem(), expected, ctx)
		}

		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(expected.Type().String()),
		}
	}

	// if ctx.Depth > 10 { panic("deepValueEqual") }	// for debugging

	// We want to avoid putting more in the visited map than we need to.
	// For any possible reference cycle that might be encountered,
	// hard(t) needs to return true for at least one of the types in the cycle.
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if got.CanAddr() && expected.CanAddr() && hard(got.Kind()) {
		addr1 := unsafe.Pointer(got.UnsafeAddr())
		addr2 := unsafe.Pointer(expected.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		v := visit{
			a1:  addr1,
			a2:  addr2,
			typ: got.Type(),
		}
		if ctx.visited[v] {
			return
		}

		// Remember for later.
		ctx.visited[v] = true
	}

	switch got.Kind() {
	case reflect.Array:
		for i := 0; i < got.Len(); i++ {
			err = deepValueEqual(got.Index(i), expected.Index(i),
				ctx.AddArrayIndex(i))
			if err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		if got.IsNil() != expected.IsNil() {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil slice",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			}
		}
		if got.Len() != expected.Len() {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "slice len",
				Got:      rawInt(got.Len()),
				Expected: rawInt(expected.Len()),
			}
		}
		if got.Pointer() == expected.Pointer() {
			return
		}
		for i := 0; i < got.Len(); i++ {
			err = deepValueEqual(got.Index(i), expected.Index(i),
				ctx.AddArrayIndex(i))
			if err != nil {
				return
			}
		}
		return

	case reflect.Interface:
		if got.IsNil() || expected.IsNil() {
			if got.IsNil() == expected.IsNil() {
				return
			}
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil interface",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			}
		}
		return deepValueEqual(got.Elem(), expected.Elem(), ctx)

	case reflect.Ptr:
		if got.Pointer() == expected.Pointer() {
			return
		}
		return deepValueEqual(got.Elem(), expected.Elem(), ctx.AddPtr(1))

	case reflect.Struct:
		sType := got.Type()
		for i, n := 0, got.NumField(); i < n; i++ {
			err = deepValueEqual(got.Field(i), expected.Field(i),
				ctx.AddDepth("."+sType.Field(i).Name))
			if err != nil {
				return
			}
		}
		return

	case reflect.Map:
		if got.IsNil() != expected.IsNil() {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil map",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			}
		}

		if got.Pointer() == expected.Pointer() {
			return
		}

		var notFoundKeys []reflect.Value
		foundKeys := map[interface{}]bool{}

		for _, vkey := range expected.MapKeys() {
			gotValue := got.MapIndex(vkey)
			if !gotValue.IsValid() {
				notFoundKeys = append(notFoundKeys, vkey)
				continue
			}

			err = deepValueEqual(gotValue, expected.MapIndex(vkey),
				ctx.AddDepth("["+toString(vkey)+"]"))
			if err != nil {
				return
			}
			foundKeys[mustGetInterface(vkey)] = true
		}

		if got.Len() == len(foundKeys) {
			if len(notFoundKeys) == 0 {
				return
			}

			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context: ctx,
				Message: "comparing map",
				Summary: tdSetResult{
					Kind:    keysSetResult,
					Missing: notFoundKeys,
				},
			}
		}

		if ctx.booleanError {
			return booleanError
		}

		// Retrieve extra keys
		res := tdSetResult{
			Kind:    keysSetResult,
			Missing: notFoundKeys,
			Extra:   make([]reflect.Value, 0, got.Len()-len(foundKeys)),
		}

		for _, vkey := range got.MapKeys() {
			if !foundKeys[vkey.Interface()] {
				res.Extra = append(res.Extra, vkey)
			}
		}

		return &Error{
			Context: ctx,
			Message: "comparing map",
			Summary: res,
		}

	case reflect.Func:
		if got.IsNil() && expected.IsNil() {
			return
		}
		if ctx.booleanError {
			return booleanError
		}
		// Can't do better than this:
		return &Error{
			Context: ctx,
			Message: "functions mismatch",
			Summary: rawString("<can not be compared>"),
		}

	default:
		// Normal equality suffices
		if mustGetInterface(got) == mustGetInterface(expected) {
			return
		}
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "values differ",
			Got:      got,
			Expected: expected,
		}
	}
}

func deepValueEqualOK(got, expected reflect.Value) bool {
	return deepValueEqual(got, expected, NewBooleanContext()) == nil
}

func getInterface(val reflect.Value) (interface{}, bool) {
	if !val.IsValid() {
		return nil, true
	}

	if val.CanInterface() {
		return val.Interface(), true
	}

	// Kinds not treated specifically by deepValueEqual
	switch val.Kind() {
	case reflect.Int:
		return int(val.Int()), true
	case reflect.Int8:
		return int8(val.Int()), true
	case reflect.Int16:
		return int16(val.Int()), true
	case reflect.Int32:
		return int32(val.Int()), true
	case reflect.Int64:
		return val.Int(), true
	case reflect.Uint:
		return uint(val.Uint()), true
	case reflect.Uint8:
		return uint8(val.Uint()), true
	case reflect.Uint16:
		return uint16(val.Uint()), true
	case reflect.Uint32:
		return uint32(val.Uint()), true
	case reflect.Uint64, reflect.Uintptr:
		return val.Uint(), true
	case reflect.Float32:
		return float32(val.Float()), true
	case reflect.Float64:
		return val.Float(), true
	case reflect.Complex64:
		return complex64(val.Complex()), true
	case reflect.Complex128:
		return val.Complex(), true
	case reflect.String:
		return val.String(), true
	case reflect.Bool:
		return val.Bool(), true
	case reflect.Chan, reflect.UnsafePointer:
		return val.Pointer(), true
	default:
		return nil, false
	}
}

func mustGetInterface(val reflect.Value) interface{} {
	ret, ok := getInterface(val)
	if ok {
		return ret
	}
	panic("getInterface() does not handle " + val.Kind().String() + " kind")
}

func EqDeeply(got, expected interface{}) bool {
	return deepValueEqualOK(reflect.ValueOf(got), reflect.ValueOf(expected))
}

func EqDeeplyError(got, expected interface{}) *Error {
	return deepValueEqual(reflect.ValueOf(got), reflect.ValueOf(expected),
		NewContext("DATA"))
}

func CmpDeeply(t *testing.T, got, expected interface{},
	args ...interface{}) bool {
	err := EqDeeplyError(got, expected)
	if err == nil {
		return true
	}

	t.Helper()

	const failedTest = "Failed test"

	var label string
	switch len(args) {
	case 0:
		label = failedTest + "\n"
	case 1:
		label = failedTest + " '" + args[0].(string) + "'\n"
	default:
		label = fmt.Sprintf(failedTest+" '"+args[0].(string)+"'\n", args[1:]...)
	}

	t.Error(label + err.Error())
	return false
}
