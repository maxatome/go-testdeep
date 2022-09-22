// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/types"
)

type tdErrorIs struct {
	baseOKNil
	expected error
}

var _ TestDeep = &tdErrorIs{}

func errorToRawString(err error) types.RawString {
	if err == nil {
		return "nil"
	}
	return types.RawString(fmt.Sprintf("(%[1]T) %[1]q", err))
}

// summary(ErrorIs): checks the data is an error and matches a wrapped error
// input(ErrorIs): if(error)

// ErrorIs operator reports whether any error in an error's chain
// matches expected.
//
//	_, err := os.Open("/unknown/file")
//	td.Cmp(t, err, os.ErrNotExist)             // fails
//	td.Cmp(t, err, td.ErrorIs(os.ErrNotExist)) // succeeds
//
//	err1 := fmt.Errorf("failure1")
//	err2 := fmt.Errorf("failure2: %w", err1)
//	err3 := fmt.Errorf("failure3: %w", err2)
//	err := fmt.Errorf("failure4: %w", err3)
//	td.Cmp(t, err, td.ErrorIs(err))  // succeeds
//	td.Cmp(t, err, td.ErrorIs(err1)) // succeeds
//	td.Cmp(t, err1, td.ErrorIs(err)) // fails
//
// Behind the scene it uses [errors.Is] function.
//
// Note that like [errors.Is], expected can be nil: in this case the
// comparison succeeds when got is nil too.
//
// See also [CmpError] and [CmpNoError].
func ErrorIs(expected error) TestDeep {
	return &tdErrorIs{
		baseOKNil: newBaseOKNil(3),
		expected:  expected,
	}
}

func (e *tdErrorIs) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	// nil case
	if !got.IsValid() {
		// Special case
		if e.expected == nil {
			return nil
		}

		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "nil value",
			Got:      types.RawString("nil"),
			Expected: types.RawString("anything implementing error interface"),
		})
	}

	gotIf, ok := dark.GetInterface(got, true)
	if !ok {
		return ctx.CollectError(ctx.CannotCompareError())
	}

	gotErr, ok := gotIf.(error)
	if !ok {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  got.Type().String() + " does not implement error interface",
			Got:      gotIf,
			Expected: types.RawString("anything implementing error interface"),
		})
	}

	if errors.Is(gotErr, e.expected) {
		return nil
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	return ctx.CollectError(&ctxerr.Error{
		Message:  "is not the error",
		Got:      errorToRawString(gotErr),
		Expected: errorToRawString(e.expected),
	})
}

func (e *tdErrorIs) String() string {
	if e.expected == nil {
		return "ErrorIs(nil)"
	}
	return "ErrorIs(" + e.expected.Error() + ")"
}
