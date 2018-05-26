// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"
	"time"
)

type tdTruncTime struct {
	Base
	expectedType reflect.Type
	expectedTime time.Time
	trunc        time.Duration
}

var _ TestDeep = &tdTruncTime{}

// TruncTime operator compares time.Time (or assignable) values after
// truncating them to the optional "trunc" duration. See time.Truncate
// for details about the truncation.
//
// If "trunc" is missing, it defaults to 0.
//
// Whatever the "trunc" value is, the monotonic clock is stripped
// before the comparison against "expectedTime".
func TruncTime(expectedTime interface{}, trunc ...time.Duration) TestDeep {
	if len(trunc) <= 1 {
		t := tdTruncTime{
			Base: NewBase(3),
		}

		if len(trunc) == 1 {
			t.trunc = trunc[0]
		}

		vval := reflect.ValueOf(expectedTime)

		t.expectedType = vval.Type()
		if t.expectedType == timeType {
			t.expectedTime = expectedTime.(time.Time).Truncate(t.trunc)
			return &t
		}
		if t.expectedType.ConvertibleTo(timeType) {
			t.expectedTime = vval.Convert(timeType).
				Interface().(time.Time).Truncate(t.trunc)
			return &t
		}
	}
	panic("usage: TruncTime(time.Time[, time.Duration])")
}

func (t *tdTruncTime) Match(ctx Context, got reflect.Value) *Error {
	if got.Type() != t.expectedType {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(t.expectedType.String()),
			Location: t.GetLocation(),
		}
	}

	gotTime, err := getTime(ctx, got, got.Type() != timeType)
	if err != nil {
		return err
	}
	gotTimeTrunc := gotTime.Truncate(t.trunc)

	if gotTimeTrunc.Equal(t.expectedTime) {
		return nil
	}

	// Fail
	if ctx.booleanError {
		return booleanError
	}

	var gotRawStr, gotTruncStr string
	if t.expectedType != timeType &&
		t.expectedType.Implements(stringerInterface) {
		gotRawStr = got.Interface().(fmt.Stringer).String()
		gotTruncStr = reflect.ValueOf(gotTimeTrunc).Convert(t.expectedType).
			Interface().(fmt.Stringer).String()
	} else {
		gotRawStr = gotTime.String()
		gotTruncStr = gotTimeTrunc.String()
	}

	return &Error{
		Context:  ctx,
		Message:  "values differ",
		Got:      rawString(gotRawStr + "\ntruncated to:\n" + gotTruncStr),
		Expected: t,
		Location: t.GetLocation(),
	}
}

func (t *tdTruncTime) String() string {
	if t.expectedType.Implements(stringerInterface) {
		return reflect.ValueOf(t.expectedTime).Convert(t.expectedType).
			Interface().(fmt.Stringer).String()
	}
	return t.expectedTime.String()
}
