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

type tdContains struct {
	tdSmugglerBase
}

var _ TestDeep = &tdContains{}

// Contains is a smuggler operator with a little convenient exception
// for strings. Contains has to be applied on arrays, slices, maps or
// strings. It compares each item of data array/slice/map/string (rune
// for strings) against "expectedValue".
//
//   list := []int{12, 34, 28}
//   CmpDeeply(t, list, Contains(34))              // succeeds
//   CmpDeeply(t, list, Contains(Between(30, 35))) // succeeds too
//   CmpDeeply(t, list, Contains(35))              // fails
//
//   hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
//   CmpDeeply(t, hash, Contains(34))              // succeeds
//   CmpDeeply(t, hash, Contains(Between(30, 35))) // succeeds too
//   CmpDeeply(t, hash, Contains(35))              // fails
//
//   got := "foo bar"
//   CmpDeeply(t, hash, Contains('o'))               // succeeds
//   CmpDeeply(t, hash, Contains(rune('o')))         // succeeds
//   CmpDeeply(t, hash, Contains(Between('n', 'p'))) // succeeds
//
// When Contains(nil) is used, nil is automatically converted to a
// typed nil on the fly to avoid confusion (if the array/slice/map
// item type allows it of course.) So all following CmpDeeply calls
// are equivalent (except the (*byte)(nil) one):
//
//   num := 123
//   list := []*int{&num, nil}
//   CmpDeeply(t, list, Contains(nil))         // succeeds → (*int)(nil)
//   CmpDeeply(t, list, Contains((*int)(nil))) // succeeds
//   CmpDeeply(t, list, Contains(Nil()))       // succeeds
//   // But...
//   CmpDeeply(t, list, Contains((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
//
// As well as these ones:
//
//   hash := map[string]*int{"foo": nil, "bar": &num}
//   CmpDeeply(t, hash, Contains(nil))         // succeeds → (*int)(nil)
//   CmpDeeply(t, hash, Contains((*int)(nil))) // succeeds
//   CmpDeeply(t, hash, Contains(Nil()))       // succeeds
//
// As a special case for string (or convertible), error or
// fmt.Stringer interface (error interface is tested before
// fmt.Stringer), "expectedValue" can be a string, a rune or a
// byte. In this case, it will test if the got string contains this
// expected string, rune or byte.
//
//   type Foobar string
//   CmpDeeply(t, Foobar("foobar"), Contains("ooba")) // succeeds
//
//   err := errors.New("error!")
//   CmpDeeply(t, err, Contains("ror")) // succeeds
//
//   bstr := bytes.NewBufferString("fmt.Stringer!")
//   CmpDeeply(t, bstr, Contains("String")) // succeeds
func Contains(expectedValue interface{}) TestDeep {
	c := tdContains{
		tdSmugglerBase: newSmugglerBase(expectedValue),
	}

	if !c.isTestDeeper {
		c.expectedValue = reflect.ValueOf(expectedValue)
	}
	return &c
}

func (c *tdContains) doesNotContains(ctx Context, got interface{}) *Error {
	if ctx.BooleanError {
		return booleanError
	}
	return ctx.CollectError(&Error{
		Message:  "does not contain",
		Got:      got,
		Expected: c,
	})
}

// getExpectedValue returns the expected value handling the
// Contains(nil) case: in this case it returns a typed nil (same type
// as the items of got).
// got is an array, a slice or a map (it's the caller responsibility to check)
func (c *tdContains) getExpectedValue(got reflect.Value) reflect.Value {
	// If the expectValue is non-typed nil
	if !c.expectedValue.IsValid() {
		// AND the kind of items in got is...
		switch got.Type().Elem().Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			// returns a typed nil
			return reflect.Zero(got.Type().Elem())
		}
	}
	return c.expectedValue
}

func (c *tdContains) Match(ctx Context, got reflect.Value) *Error {
	switch got.Kind() {
	case reflect.Array, reflect.Slice:
		expectedValue := c.getExpectedValue(got)
		for index := got.Len() - 1; index >= 0; index-- {
			if deepValueEqualOK(got.Index(index), expectedValue) {
				return nil
			}
		}
		return c.doesNotContains(ctx, got)

	case reflect.Map:
		expectedValue := c.getExpectedValue(got)
		for _, vkey := range got.MapKeys() {
			if deepValueEqualOK(got.MapIndex(vkey), expectedValue) {
				return nil
			}
		}
		return c.doesNotContains(ctx, got)

		// For String kind *AND* TestDeep operator, applies this operator on
		// each character of the string
	case reflect.String:
		if c.isTestDeeper {
			for _, chr := range got.String() {
				if deepValueEqualOK(reflect.ValueOf(chr), c.expectedValue) {
					return nil
				}
			}
			return c.doesNotContains(ctx, got)
		}
	}

	// A TestDeep operator can only be applied on arrays, slices, map
	// and as a special feature on strings (all handled in switch
	// above). For all other cases, it is an error.
	if !c.isTestDeeper {
		// If expectedValue is a string, a rune or a byte, we try to get a
		// string from got and check whether expectedValue is contained in
		// this string or not
		switch expectedKind := c.expectedValue.Kind(); expectedKind {
		case reflect.String, reflect.Int32, reflect.Uint8: // string, rune & byte
			str, err := getString(ctx, got)
			if err != nil {
				return err
			}

			switch expectedKind {
			case reflect.String:
				if strings.Contains(str, c.expectedValue.String()) {
					return nil
				}
			case reflect.Int32:
				if strings.ContainsRune(str, rune(c.expectedValue.Int())) {
					return nil
				}
			default: // = case reflect.Uint8:
				if strings.IndexByte(str, byte(c.expectedValue.Uint())) >= 0 {
					return nil
				}
			}
			return c.doesNotContains(ctx, str)
		}
	}

	if ctx.BooleanError {
		return booleanError
	}
	var expectedType interface{}
	if c.expectedValue.IsValid() {
		expectedType = rawString(c.expectedValue.Type().String())
	} else {
		expectedType = c
	}

	return ctx.CollectError(&Error{
		Message:  "cannot check contains",
		Got:      rawString(got.Type().String()),
		Expected: expectedType,
	})
}

func (c *tdContains) String() string {
	return "Contains(" + toString(c.expectedValue) + ")"
}
