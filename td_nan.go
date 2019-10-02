// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"math"
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
)

type tdNaN struct {
	base
}

var _ TestDeep = &tdNaN{}

// summary(NaN): checks a floating number is [`math.NaN`]

// NaN operator checks that data is a float and is not-a-number.
func NaN() TestDeep {
	return &tdNaN{
		base: newBase(3),
	}
}

func (n *tdNaN) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	switch got.Kind() {
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(got.Float()) {
			return nil
		}

		return ctx.CollectError(&ctxerr.Error{
			Message:  "values differ",
			Got:      got,
			Expected: n,
		})
	}

	return ctx.CollectError(&ctxerr.Error{
		Message:  "type mismatch",
		Got:      types.RawString(got.Type().String()),
		Expected: types.RawString("float32 OR float64"),
	})
}

func (n *tdNaN) String() string {
	return "NaN"
}

type tdNotNaN struct {
	base
}

var _ TestDeep = &tdNotNaN{}

// summary(NotNaN): checks a floating number is not [`math.NaN`]

// NotNaN operator checks that data is a float and is not not-a-number.
func NotNaN() TestDeep {
	return &tdNotNaN{
		base: newBase(3),
	}
}

func (n *tdNotNaN) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	switch got.Kind() {
	case reflect.Float32, reflect.Float64:
		if !math.IsNaN(got.Float()) {
			return nil
		}

		return ctx.CollectError(&ctxerr.Error{
			Message:  "values differ",
			Got:      got,
			Expected: n,
		})
	}

	return ctx.CollectError(&ctxerr.Error{
		Message:  "type mismatch",
		Got:      types.RawString(got.Type().String()),
		Expected: types.RawString("float32 OR float64"),
	})
}

func (n *tdNotNaN) String() string {
	return "not NaN"
}
