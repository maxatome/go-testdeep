package testdeep

import (
	"fmt"
	"reflect"
	"time"
)

type tdTruncTime struct {
	TestDeepBase
	expectedType reflect.Type
	expectedTime time.Time
	trunc        time.Duration
}

var _ TestDeep = &tdTruncTime{}

func TruncTime(val interface{}, trunc ...time.Duration) TestDeep {
	if len(trunc) <= 1 {
		t := tdTruncTime{
			TestDeepBase: NewTestDeepBase(3),
		}

		if len(trunc) == 1 {
			t.trunc = trunc[0]
		}

		vval := reflect.ValueOf(val)

		t.expectedType = vval.Type()
		if t.expectedType == timeType {
			t.expectedTime = val.(time.Time).Truncate(t.trunc)
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

	var (
		gotIf interface{}
		ok    bool
	)
	if got.Type() == timeType {
		gotIf, ok = getInterface(got, true)
	} else {
		gotIf, ok = getInterface(got.Convert(timeType), true)
	}
	if !ok {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context: ctx,
			Message: "cannot compare unexported field that cannot be overridden",
			Summary: "",
		}
	}

	gotTime := gotIf.(time.Time)
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
