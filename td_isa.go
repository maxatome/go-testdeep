package testdeep

import (
	"reflect"
)

type tdIsa struct {
	TestDeepBase
	expectedModel reflect.Value
}

var _ TestDeep = &tdIsa{}

func Isa(model interface{}) TestDeep {
	return &tdIsa{
		TestDeepBase:  NewTestDeepBase(3),
		expectedModel: reflect.ValueOf(model),
	}
}

func (i *tdIsa) Match(ctx Context, got reflect.Value) (err *Error) {
	if got.Type() == i.expectedModel.Type() {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "type mismatch",
		Got:      rawString(got.Type().String()),
		Expected: rawString(i.expectedModel.Type().String()),
		Location: i.GetLocation(),
	}
}

func (i *tdIsa) String() string {
	return i.expectedModel.Type().String()
}
