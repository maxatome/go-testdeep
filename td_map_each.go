package testdeep

import (
	"reflect"
	"strings"
)

type tdMapEach struct {
	BaseOKNil
	expected reflect.Value
}

var _ TestDeep = &tdMapEach{}

func MapEach(item interface{}) TestDeep {
	return &tdMapEach{
		BaseOKNil: NewBaseOKNil(3),
		expected:  reflect.ValueOf(item),
	}
}

func (m *tdMapEach) Match(ctx Context, got reflect.Value) (err *Error) {
	if !got.IsValid() {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "nil value",
			Got:      rawString("nil"),
			Expected: rawString("Map OR *Map"),
			Location: m.GetLocation(),
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
				Expected: rawString("Map OR *Map"),
				Location: m.GetLocation(),
			}
		}

		if gotElem.Kind() != reflect.Map {
			break
		}
		got = gotElem
		fallthrough

	case reflect.Map:
		for _, key := range got.MapKeys() {
			err = deepValueEqual(ctx.AddDepth("["+toString(key)+"]"),
				got.MapIndex(key), m.expected)
			if err != nil {
				return err.SetLocationIfMissing(m)
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
		Expected: rawString("Map OR *Map"),
		Location: m.GetLocation(),
	}
}

func (m *tdMapEach) String() string {
	const prefix = "MapEach("

	content := toString(m.expected)
	if strings.Contains(content, "\n") {
		return prefix + indentString(content, "        ") + ")"
	}
	return prefix + content + ")"
}
