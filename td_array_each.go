package testdeep

import (
	"reflect"
	"strings"
)

type tdArrayEach struct {
	BaseOKNil
	expected reflect.Value
}

var _ TestDeep = &tdArrayEach{}

func ArrayEach(item interface{}) TestDeep {
	return &tdArrayEach{
		BaseOKNil: NewBaseOKNil(3),
		expected:  reflect.ValueOf(item),
	}
}

func (a *tdArrayEach) Match(ctx Context, got reflect.Value) (err *Error) {
	if !got.IsValid() {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "nil value",
			Got:      rawString("nil"),
			Expected: rawString("Slice OR Array OR *Slice OR *Array"),
			Location: a.GetLocation(),
		}
	}

	switch got.Kind() {
	case reflect.Ptr:
		gotElem := got.Elem()
		if !gotElem.IsValid() {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil pointer",
				Got:      rawString("nil " + got.Type().String()),
				Expected: rawString("Slice OR Array OR *Slice OR *Array"),
				Location: a.GetLocation(),
			}
		}

		if gotElem.Kind() != reflect.Array && gotElem.Kind() != reflect.Slice {
			break
		}
		got = gotElem
		fallthrough

	case reflect.Array, reflect.Slice:
		gotLen := got.Len()

		for idx := 0; idx < gotLen; idx++ {
			err = deepValueEqual(got.Index(idx), a.expected, ctx.AddArrayIndex(idx))
			if err != nil {
				return err.SetLocationIfMissing(a)
			}
		}
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "bad type",
		Got:      rawString(got.Type().String()),
		Expected: rawString("Slice OR Array OR *Slice OR *Array"),
		Location: a.GetLocation(),
	}
}

func (a *tdArrayEach) String() string {
	const prefix = "ArrayEach("

	content := toString(a.expected)
	if strings.Contains(content, "\n") {
		return prefix + indentString(content, "          ") + ")"
	}
	return prefix + content + ")"
}
