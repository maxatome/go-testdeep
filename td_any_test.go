// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestAny(t *testing.T) {
	checkOK(t, 6, testdeep.Any(nil, 5, 6, 7))
	checkOK(t, nil, testdeep.Any(5, 6, 7, nil))

	checkError(t, 6, testdeep.Any(5),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(5)"),
		})

	checkError(t, 6, testdeep.Any(nil),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(nil)"),
		})

	checkError(t, nil, testdeep.Any(6),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Any(6)"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Any(6).String(), "Any(6)")
	test.EqualStr(t, testdeep.Any(6, 7).String(), "Any(6,\n    7)")
}

func TestAnyTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Any(6, nil), nil)
	equalTypes(t, testdeep.Any(6, "toto"), nil)

	equalTypes(t, testdeep.Any(6, testdeep.Zero(), 7, 8), 26)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		testdeep.Any(
			testdeep.Empty(),
			5,
			testdeep.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			testdeep.Any(6, 7),
			testdeep.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		42)

	// Only one interface type
	equalTypes(t,
		testdeep.Any(
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*error)(nil)),
		),
		(*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		testdeep.Any(
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*fmt.Stringer)(nil)),
		),
		nil)

	equalTypes(t,
		testdeep.Any(
			testdeep.Code(func(x interface{}) bool { return true }),
			testdeep.Code(func(y int) bool { return true }),
		),
		12)
}
