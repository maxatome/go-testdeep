// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/util"
)

type tdContainsKey struct {
	tdSmugglerBase
}

var _ TestDeep = &tdContainsKey{}

// ContainsKey is a smuggler operator and works on maps only. It
// compares each key of map against "expectedValue".
//
//   hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
//   CmpDeeply(t, hash, ContainsKey("foo"))          // succeeds
//   CmpDeeply(t, hash, ContainsKey(HasPrefix("z"))) // succeeds
//   CmpDeeply(t, hash, ContainsKey(HasPrefix("x"))  // fails
//
//   hnum := map[int]string{1: "foo", 42: "bar"}
//   CmpDeeply(t, hash, ContainsKey(42))             // succeeds
//   CmpDeeply(t, hash, ContainsKey(Between(40, 45)) // succeeds
//
// When ContainsKey(nil) is used, nil is automatically converted to a
// typed nil on the fly to avoid confusion (if the map key type allows
// it of course.) So all following CmpDeeply calls are equivalent
// (except the (*byte)(nil) one):
//
//   num := 123
//   hnum := map[*int]bool{&num: true, nil: true}
//   CmpDeeply(t, hnum, ContainsKey(nil))          // succeeds → (*int)(nil)
//   CmpDeeply(t, hnum, ContainsKey((*int)(nil)))  // succeeds
//   CmpDeeply(t, hnum, ContainsKey(Nil()))        // succeeds
//   // But...
//   CmpDeeply(t, hnum, ContainsKey((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
func ContainsKey(expectedValue interface{}) TestDeep {
	c := tdContainsKey{
		tdSmugglerBase: newSmugglerBase(expectedValue),
	}

	if !c.isTestDeeper {
		c.expectedValue = reflect.ValueOf(expectedValue)
	}
	return &c
}

func (c *tdContainsKey) doesNotContainKey(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if ctx.BooleanError {
		return ctxerr.BooleanError
	}

	keys := append(make([]reflect.Value, 0, got.Len()), got.MapKeys()...)

	return ctx.CollectError(&ctxerr.Error{
		Message: "does not contain key",
		Summary: ctxerr.ErrorSummaryItems{
			{
				Label: "expected key",
				Value: util.ToString(c.expectedValue),
			},
			{
				Label: "not in keys",
				Value: util.ToString(keys),
			},
		},
	})
}

// getExpectedValue returns the expected value handling the
// Contains(nil) case: in this case it returns a typed nil (same type
// as the keys of got).
// got is a map (it's the caller responsibility to check)
func (c *tdContainsKey) getExpectedValue(got reflect.Value) reflect.Value {
	// If the expectValue is non-typed nil
	if !c.expectedValue.IsValid() {
		// AND the kind of items in got is...
		switch got.Type().Key().Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			// returns a typed nil
			return reflect.Zero(got.Type().Key())
		}
	}
	return c.expectedValue
}

func (c *tdContainsKey) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if got.Kind() == reflect.Map {
		expectedValue := c.getExpectedValue(got)

		for _, vkey := range got.MapKeys() {
			if deepValueEqualOK(vkey, expectedValue) {
				return nil
			}
		}
		return c.doesNotContainKey(ctx, got)
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	var expectedType interface{}
	if c.expectedValue.IsValid() {
		expectedType = types.RawString(c.expectedValue.Type().String())
	} else {
		expectedType = c
	}

	return ctx.CollectError(&ctxerr.Error{
		Message:  "cannot check contains key",
		Got:      types.RawString(got.Type().String()),
		Expected: expectedType,
	})
}

func (c *tdContainsKey) String() string {
	return "ContainsKey(" + util.ToString(c.expectedValue) + ")"
}
