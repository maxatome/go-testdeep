package testdeep

import (
	"reflect"
)

type tdPtr struct {
	tdSmuggler
}

var _ TestDeep = &tdPtr{}

func Ptr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPtr{
			tdSmuggler: newSmuggler(val),
		}

		if !p.isTestDeeper {
			p.expectedValue = reflect.New(vval.Type())
			p.expectedValue.Elem().Set(vval)
		}
		return &p
	}
	panic("usage: Ptr(NON_NIL_VALUE)")
}

func (p *tdPtr) Match(ctx Context, got reflect.Value) (err *Error) {
	if got.Kind() != reflect.Ptr {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "pointer type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(p.String()),
			Location: p.GetLocation(),
		}
	}

	if p.isTestDeeper {
		err = deepValueEqual(ctx.AddPtr(1), got.Elem(), p.expectedValue)
	} else {
		err = deepValueEqual(ctx, got, p.expectedValue)
	}
	return err.SetLocationIfMissing(p)
}

func (p *tdPtr) String() string {
	if p.isTestDeeper {
		return "*<something>"
	}
	return p.expectedValue.Type().String()
}

type tdPPtr struct {
	tdSmuggler
}

var _ TestDeep = &tdPPtr{}

func PPtr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPPtr{
			tdSmuggler: newSmuggler(val),
		}

		if !p.isTestDeeper {
			pVval := reflect.New(vval.Type())
			pVval.Elem().Set(vval)

			p.expectedValue = reflect.New(pVval.Type())
			p.expectedValue.Elem().Set(pVval)
		}
		return &p
	}
	panic("usage: PPtr(NON_NIL_VALUE)")
}

func (p *tdPPtr) Match(ctx Context, got reflect.Value) (err *Error) {
	if got.Kind() != reflect.Ptr || got.Elem().Kind() != reflect.Ptr {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "pointer type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(p.String()),
			Location: p.GetLocation(),
		}
	}

	if p.isTestDeeper {
		err = deepValueEqual(ctx.AddPtr(2), got.Elem().Elem(), p.expectedValue)
	} else {
		err = deepValueEqual(ctx, got, p.expectedValue)
	}
	return err.SetLocationIfMissing(p)
}

func (p *tdPPtr) String() string {
	if p.isTestDeeper {
		return "**<something>"
	}
	return p.expectedValue.Type().String()
}
