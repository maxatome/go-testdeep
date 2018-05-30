// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdZero struct {
	BaseOKNil
}

var _ TestDeep = &tdZero{}

// Zero operator checks that data is zero regarding its type.
//
//   nil is the zero value of pointers, maps, slices, channels and functions;
//   0 is the zero value of numbers;
//   false is the zero value of booleans;
//   zero value of structs is the struct with no fields initialized.
func Zero() TestDeep {
	return &tdZero{
		BaseOKNil: NewBaseOKNil(3),
	}
}

func (z *tdZero) Match(ctx Context, got reflect.Value) (err *Error) {
	// nil case
	if !got.IsValid() {
		return nil
	}
	return deepValueEqual(ctx, got, reflect.New(got.Type()).Elem()).
		SetLocationIfMissing(z)
}

func (z *tdZero) String() string {
	return "Zero()"
}
