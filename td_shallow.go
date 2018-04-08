package testdeep

import (
	"fmt"
	"reflect"
)

type tdShallow struct {
	Base
	expectedKind    reflect.Kind
	expectedPointer uintptr
}

var _ TestDeep = &tdShallow{}

func Shallow(ptr interface{}) TestDeep {
	vptr := reflect.ValueOf(ptr)

	shallow := tdShallow{
		Base:         NewBase(3),
		expectedKind: vptr.Kind(),
	}

	// Note form reflect documentation:
	// If v's Kind is Func, the returned pointer is an underlying code
	// pointer, but not necessarily enough to identify a single function
	// uniquely. The only guarantee is that the result is zero if and
	// only if v is a nil func Value.

	switch shallow.expectedKind {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice,
		reflect.UnsafePointer:
		shallow.expectedPointer = vptr.Pointer()
		return &shallow

	default:
		panic("usage: Shallow(CHANNEL|FUNC|MAP|PTR|SLICE|UNSAFE_PTR)")
	}
}

func (s *tdShallow) Match(ctx Context, got reflect.Value) *Error {
	if got.Kind() != s.expectedKind {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad kind",
			Got:      rawString(got.Kind().String()),
			Expected: rawString(s.expectedKind.String()),
			Location: s.GetLocation(),
		}
	}

	if got.Pointer() != s.expectedPointer {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  fmt.Sprintf("%s pointer mismatch", s.expectedKind),
			Got:      rawString(fmt.Sprintf("0x%x", got.Pointer())),
			Expected: rawString(fmt.Sprintf("0x%x", s.expectedPointer)),
			Location: s.GetLocation(),
		}
	}
	return nil
}

func (s *tdShallow) String() string {
	return fmt.Sprintf("(%s) 0x%x", s.expectedKind, s.expectedPointer)
}
