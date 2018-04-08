package testdeep

import (
	"reflect"
	"time"
)

func ternRune(cond bool, a, b rune) rune {
	if cond {
		return a
	}
	return b
}

func ternStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// getTime returns the time.Time that is inside got or that can be
// converted from got content.
func getTime(ctx Context, got reflect.Value, mustConvert bool) (time.Time, *Error) {
	var (
		gotIf interface{}
		ok    bool
	)
	if mustConvert {
		gotIf, ok = getInterface(got.Convert(timeType), true)
	} else {
		gotIf, ok = getInterface(got, true)
	}
	if !ok {
		if ctx.booleanError {
			return time.Time{}, booleanError
		}
		return time.Time{}, &Error{
			Context: ctx,
			Message: "cannot compare unexported field that cannot be overridden",
		}
	}
	return gotIf.(time.Time), nil
}
