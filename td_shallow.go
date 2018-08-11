// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
)

type tdShallow struct {
	Base
	expectedKind    reflect.Kind
	expectedPointer uintptr
}

var _ TestDeep = &tdShallow{}

// Shallow operator compares pointers only, not their contents. It
// applies on channels, functions (with some restrictions), maps,
// pointers and slices.
//
// During a match, the compared data must be the same as
// "expectedPointer" to succeed.
//
//   a, b := 123, 123
//   CmpDeeply(t, &a, Shallow(&a)) // succeeds
//   CmpDeeply(t, &a, Shallow(&b)) // fails even if a == b as &a != &b
func Shallow(expectedPtr interface{}) TestDeep {
	vptr := reflect.ValueOf(expectedPtr)

	shallow := tdShallow{
		Base:         NewBase(3),
		expectedKind: vptr.Kind(),
	}

	// Note from reflect documentation:
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

func (s *tdShallow) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if got.Kind() != s.expectedKind {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "bad kind",
			Got:      types.RawString(got.Kind().String()),
			Expected: types.RawString(s.expectedKind.String()),
		})
	}

	if got.Pointer() != s.expectedPointer {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  fmt.Sprintf("%s pointer mismatch", s.expectedKind),
			Got:      types.RawString(fmt.Sprintf("0x%x", got.Pointer())),
			Expected: types.RawString(fmt.Sprintf("0x%x", s.expectedPointer)),
		})
	}
	return nil
}

func (s *tdShallow) String() string {
	return fmt.Sprintf("(%s) 0x%x", s.expectedKind, s.expectedPointer)
}
