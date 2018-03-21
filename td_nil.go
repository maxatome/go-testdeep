package testdeep

import (
	"reflect"
)

type tdNil struct {
	TestDeepBaseOKNil
}

var _ TestDeep = &tdNil{}

func Nil() TestDeep {
	return &tdNil{
		TestDeepBaseOKNil: NewTestDeepBaseOKNil(3),
	}
}

func (n *tdNil) Match(ctx Context, got reflect.Value) *Error {
	if !got.IsValid() {
		return nil
	}

	switch got.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		if got.IsNil() {
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "non-nil",
		Got:      got,
		Expected: n,
		Location: n.GetLocation(),
	}
}

func (n *tdNil) String() string {
	return "nil"
}

type tdNotNil struct {
	TestDeepBaseOKNil
}

var _ TestDeep = &tdNotNil{}

func NotNil() TestDeep {
	return &tdNotNil{
		TestDeepBaseOKNil: NewTestDeepBaseOKNil(3),
	}
}

func (n *tdNotNil) Match(ctx Context, got reflect.Value) *Error {
	if got.IsValid() {
		switch got.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			if !got.IsNil() {
				return nil
			}

			// All other kinds are non-nil by nature
		default:
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "nil value",
		Got:      got,
		Expected: n,
		Location: n.GetLocation(),
	}
}

func (n *tdNotNil) String() string {
	return "not nil"
}
