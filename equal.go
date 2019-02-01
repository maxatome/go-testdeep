// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// deepValueEqual function is heavily based on reflect.deepValueEqual function
// licensed under the BSD-style license found in the LICENSE file in the
// golang repository: https://github.com/golang/go/blob/master/LICENSE

package testdeep

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/types"
)

func isNilStr(isNil bool) types.RawString {
	if isNil {
		return "nil"
	}
	return "not nil"
}

func deepValueEqualFinal(ctx ctxerr.Context, got, expected reflect.Value) (err *ctxerr.Error) {
	err = deepValueEqual(ctx, got, expected)
	if err == nil {
		// Try to merge pending errors
		errMerge := ctx.MergeErrors()
		if errMerge != nil {
			return errMerge
		}
	}
	return
}

// nilHandler is called when one of got or expected is nil (but never
// both, it is caller responsibility)
func nilHandler(ctx ctxerr.Context, got, expected reflect.Value) *ctxerr.Error {
	err := ctxerr.Error{}

	if expected.IsValid() { // here: !got.IsValid()
		if expected.Type().Implements(testDeeper) {
			curOperator := expected.Interface().(TestDeep)
			ctx.CurOperator = curOperator
			if curOperator.HandleInvalid() {
				return curOperator.Match(ctx, got)
			}
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}

			// Special case if "expected" is a TestDeep operator which
			// does not handle invalid values: the operator is not called,
			// but for the user the error comes from it
		} else if ctx.BooleanError {
			return ctxerr.BooleanError
		}

		err.Expected = expected
	} else { // here: !expected.IsValid() && got.IsValid()
		// Special case: "got" is a nil interface, so consider as equal
		// to "expected" nil.
		if got.Kind() == reflect.Interface && got.IsNil() {
			return nil
		}

		if ctx.BooleanError {
			return ctxerr.BooleanError
		}

		err.Got = got
	}

	err.Message = "values differ"
	return ctx.CollectError(&err)
}

func deepValueEqual(ctx ctxerr.Context, got, expected reflect.Value) (err *ctxerr.Error) {
	if !got.IsValid() || !expected.IsValid() {
		if got.IsValid() == expected.IsValid() {
			return
		}
		return nilHandler(ctx, got, expected)
	}

	// Special case, "got" implements testDeeper: only if allowed
	if got.Type().Implements(testDeeper) {
		panic("Found a TestDeep operator in got param, " +
			"can only use it in expected one!")
	}

	if got.Type() != expected.Type() {
		if expected.Type().Implements(testDeeper) {
			curOperator := expected.Interface().(TestDeep)

			// Resolve interface
			if got.Kind() == reflect.Interface {
				got = got.Elem()

				if !got.IsValid() {
					return nilHandler(ctx, got, expected)
				}
			}

			ctx.CurOperator = curOperator
			return curOperator.Match(ctx, got)
		}

		// "expected" is not a TestDeep operator

		// If "got" is an interface, try to see what is behind before failing
		// Used by Set/Bag Match method in such cases:
		//     []interface{}{123, "foo"}  →  Bag("foo", 123)
		//    Interface kind -^-----^   but String-^ and ^- Int kinds
		if got.Kind() == reflect.Interface {
			return deepValueEqual(ctx, got.Elem(), expected)
		}

		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "type mismatch",
			Got:      types.RawString(got.Type().String()),
			Expected: types.RawString(expected.Type().String()),
		})
	}

	// if ctx.Depth > 10 { panic("deepValueEqual") }	// for debugging

	// We want to avoid putting more in the Visited map than we need to.
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
			// Canonicalize order to reduce number of entries in Visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		v := ctxerr.Visit{
			A1:  addr1,
			A2:  addr2,
			Typ: got.Type(),
		}
		if ctx.Visited[v] {
			return
		}

		// Remember for later.
		ctx.Visited[v] = true
	}

	switch got.Kind() {
	case reflect.Array:
		for i := 0; i < got.Len(); i++ {
			err = deepValueEqual(ctx.AddArrayIndex(i),
				got.Index(i), expected.Index(i))
			if err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		if got.IsNil() != expected.IsNil() {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(&ctxerr.Error{
				Message:  "nil slice",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			})
		}

		var (
			gotLen      = got.Len()
			expectedLen = expected.Len()
		)

		// Shortcut in boolean context
		if ctx.BooleanError && gotLen != expectedLen {
			return ctxerr.BooleanError
		}

		if got.Pointer() == expected.Pointer() {
			return
		}

		var maxLen int
		if gotLen >= expectedLen {
			maxLen = expectedLen
		} else {
			maxLen = gotLen
		}

		for i := 0; i < maxLen; i++ {
			err = deepValueEqual(ctx.AddArrayIndex(i),
				got.Index(i), expected.Index(i))
			if err != nil {
				return
			}
		}

		if gotLen != expectedLen {
			res := tdSetResult{
				Kind: itemsSetResult,
			}

			if gotLen > expectedLen {
				res.Extra = make([]reflect.Value, gotLen-expectedLen)
				for i := expectedLen; i < gotLen; i++ {
					res.Extra[i-expectedLen] = got.Index(i)
				}
			} else {
				res.Missing = make([]reflect.Value, expectedLen-gotLen)
				for i := gotLen; i < expectedLen; i++ {
					res.Missing[i-gotLen] = expected.Index(i)
				}
			}

			return ctx.CollectError(&ctxerr.Error{
				Message: fmt.Sprintf("comparing slices, from index #%d", maxLen),
				Summary: res.Summary(),
			})
		}
		return

	case reflect.Interface:
		if got.IsNil() || expected.IsNil() {
			if got.IsNil() == expected.IsNil() {
				return
			}
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(&ctxerr.Error{
				Message:  "nil interface",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			})
		}
		return deepValueEqual(ctx, got.Elem(), expected.Elem())

	case reflect.Ptr:
		if got.Pointer() == expected.Pointer() {
			return
		}
		return deepValueEqual(ctx.AddPtr(1), got.Elem(), expected.Elem())

	case reflect.Struct:
		sType := got.Type()
		for i, n := 0, got.NumField(); i < n; i++ {
			err = deepValueEqual(ctx.AddDepth("."+sType.Field(i).Name),
				got.Field(i), expected.Field(i))
			if err != nil {
				return
			}
		}
		return

	case reflect.Map:
		if got.IsNil() != expected.IsNil() {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(&ctxerr.Error{
				Message:  "nil map",
				Got:      isNilStr(got.IsNil()),
				Expected: isNilStr(expected.IsNil()),
			})
		}

		// Shortcut in boolean context
		if ctx.BooleanError && got.Len() != expected.Len() {
			return ctxerr.BooleanError
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

			err = deepValueEqual(ctx.AddMapKey(vkey),
				gotValue, expected.MapIndex(vkey))
			if err != nil {
				return
			}
			foundKeys[dark.MustGetInterface(vkey)] = true
		}

		if got.Len() == len(foundKeys) {
			if len(notFoundKeys) == 0 {
				return
			}
			return ctx.CollectError(&ctxerr.Error{
				Message: "comparing map",
				Summary: (tdSetResult{
					Kind:    keysSetResult,
					Missing: notFoundKeys,
				}).Summary(),
			})
		}

		if ctx.BooleanError {
			return ctxerr.BooleanError
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

		return ctx.CollectError(&ctxerr.Error{
			Message: "comparing map",
			Summary: res.Summary(),
		})

	case reflect.Func:
		if got.IsNil() && expected.IsNil() {
			return
		}
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		// Can't do better than this:
		return ctx.CollectError(&ctxerr.Error{
			Message: "functions mismatch",
			Summary: ctxerr.NewSummary("<can not be compared>"),
		})

	default:
		// Normal equality suffices
		if dark.MustGetInterface(got) == dark.MustGetInterface(expected) {
			return
		}
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "values differ",
			Got:      got,
			Expected: expected,
		})
	}
}

func deepValueEqualOK(got, expected reflect.Value) bool {
	return deepValueEqualFinal(newBooleanContext(), got, expected) == nil
}

// EqDeeply returns true if "got" matches "expected". "expected" can
// be the same type as "got" is, or contains some TestDeep operators.
func EqDeeply(got, expected interface{}) bool {
	return deepValueEqualOK(reflect.ValueOf(got), reflect.ValueOf(expected))
}

// EqDeeplyError returns nil if "got" matches "expected". "expected"
// can be the same type as got is, or contains some TestDeep
// operators. If "got" does not match "expected", the returned *ctxerr.Error
// contains the reason of the first mismatch detected.
func EqDeeplyError(got, expected interface{}) error {
	err := deepValueEqualFinal(newContext(),
		reflect.ValueOf(got), reflect.ValueOf(expected))
	if err == nil {
		return nil
	}
	return err
}
