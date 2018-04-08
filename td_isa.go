package testdeep

import (
	"reflect"
)

type tdIsa struct {
	Base
	expectedType reflect.Type
}

var _ TestDeep = &tdIsa{}

func Isa(model interface{}) TestDeep {
	return &tdIsa{
		Base:         NewBase(3),
		expectedType: reflect.ValueOf(model).Type(),
	}
}

func (i *tdIsa) Match(ctx Context, got reflect.Value) (err *Error) {
	if got.Type() == i.expectedType {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "type mismatch",
		Got:      rawString(got.Type().String()),
		Expected: rawString(i.expectedType.String()),
		Location: i.GetLocation(),
	}
}

func (i *tdIsa) String() string {
	return i.expectedType.String()
}
