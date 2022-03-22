// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestAny(t *testing.T) {
	checkOK(t, 6, td.Any(nil, 5, 6, 7))
	checkOK(t, nil, td.Any(5, 6, 7, nil))

	checkError(t, 6, td.Any(5),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(5)"),
		})

	checkError(t, 6, td.Any(nil),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(nil)"),
		})

	checkError(t, nil, td.Any(6),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Any(6)"),
		})

	// Lax
	checkOK(t, float64(123), td.Lax(td.Any(122, 123, 124)))

	//
	// String
	test.EqualStr(t, td.Any(6).String(), "Any(6)")
	test.EqualStr(t, td.Any(6, 7).String(), "Any(6,\n    7)")
}

func TestAnyTypeBehind(t *testing.T) {
	equalTypes(t, td.Any(6, nil), nil)
	equalTypes(t, td.Any(6, "toto"), nil)

	equalTypes(t, td.Any(6, td.Zero(), 7, 8), 26)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		td.Any(
			td.Empty(),
			5,
			td.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			td.Any(6, 7),
			td.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		42)

	// Only one interface type
	equalTypes(t,
		td.Any(
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
		),
		(*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		td.Any(
			td.Isa((*error)(nil)),
			td.Isa((*fmt.Stringer)(nil)),
		),
		nil)

	equalTypes(t,
		td.Any(
			td.Code(func(x any) bool { return true }),
			td.Code(func(y int) bool { return true }),
		),
		12)
}
