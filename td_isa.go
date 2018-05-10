package testdeep

import (
	"reflect"
)

type tdIsa struct {
	Base
	expectedType   reflect.Type
	checkImplement bool
}

var _ TestDeep = &tdIsa{}

// Special case when pointer on an interface XXX
func Isa(model interface{}) TestDeep {
	modelType := reflect.ValueOf(model).Type()

	return &tdIsa{
		Base:         NewBase(3),
		expectedType: modelType,
		checkImplement: modelType.Kind() == reflect.Ptr &&
			modelType.Elem().Kind() == reflect.Interface,
	}
}

func (i *tdIsa) Match(ctx Context, got reflect.Value) (err *Error) {
	gotType := got.Type()

	if gotType == i.expectedType {
		return nil
	}

	if i.checkImplement {
		if gotType.Implements(i.expectedType.Elem()) {
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "type mismatch",
		Got:      rawString(gotType.String()),
		Expected: rawString(i.expectedType.String()),
		Location: i.GetLocation(),
	}
}

func (i *tdIsa) String() string {
	return i.expectedType.String()
}
