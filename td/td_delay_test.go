// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestDelay(t *testing.T) {
	called := 0
	op := td.Delay(func() td.TestDeep {
		called++
		return td.Lt(13)
	})
	test.EqualInt(t, called, 0)
	checkOK(t, 12, op)
	test.EqualInt(t, called, 1)
	checkOK(t, 12, op)
	test.EqualInt(t, called, 1)

	delayNil := td.Delay(td.Nil)
	checkOK(t, nil, delayNil)

	test.EqualStr(t, delayNil.String(), "nil")

	checkError(t, 8,
		td.Delay(
			func() td.TestDeep {
				return td.Gt(13)
			},
		),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("8"),
			Expected: mustBe("> 13"),
		})

	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Delay(nil) },
		"Delay(DELAYED): DELAYED must be non-nil")
}

func TestDelayTypeBehind(t *testing.T) {
	equalTypes(t, td.Delay(func() td.TestDeep { return td.String("x") }), nil)
	equalTypes(t, td.Delay(func() td.TestDeep { return td.Gt(16) }), 42)
}
