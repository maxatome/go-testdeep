// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdAny struct {
	tdList
}

var _ TestDeep = &tdAny{}

// Any operator compares data against several expected values. During
// a match, at least one of them has to match to succeed.
func Any(expectedValues ...interface{}) TestDeep {
	return &tdAny{
		tdList: newList(expectedValues...),
	}
}

func (a *tdAny) Match(ctx Context, got reflect.Value) *Error {
	for _, item := range a.items {
		if deepValueEqualOK(got, item) {
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "comparing with Any",
		Got:      got,
		Expected: a,
		Location: a.GetLocation(),
	}
}
