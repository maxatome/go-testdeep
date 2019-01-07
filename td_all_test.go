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

func TestAll(t *testing.T) {
	checkOK(t, 6, testdeep.All(6, 6, 6))
	checkOK(t, nil, testdeep.All(nil, nil, nil))

	checkError(t, 6, testdeep.All(6, 5, 6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("5"),
		})

	checkError(t, 6, testdeep.All(6, nil, 6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("nil"),
		})

	checkError(t, nil, testdeep.All(nil, 5, nil),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("5"),
		})

	checkError(t,
		6,
		testdeep.All(
			6,
			testdeep.All(testdeep.Between(3, 8), testdeep.Between(4, 5)),
			6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("All(3 ≤ got ≤ 8,\n    4 ≤ got ≤ 5)"),
			Origin: &expectedError{
				Message:  mustBe("compared (part 2 of 2)"),
				Path:     mustBe("DATA<All#2/3>"),
				Got:      mustBe("6"),
				Expected: mustBe("4 ≤ got ≤ 5"),
				Origin: &expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("DATA<All#2/3><All#2/2>"),
					Got:      mustBe("6"),
					Expected: mustBe("4 ≤ got ≤ 5"),
				},
			},
		})

	//
	// String
	test.EqualStr(t, testdeep.All(6).String(), "All(6)")
	test.EqualStr(t, testdeep.All(6, 7).String(), "All(6,\n    7)")
}

func TestAllTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.All(6, nil), nil)
	equalTypes(t, testdeep.All(6, "toto"), nil)

	equalTypes(t, testdeep.All(6, testdeep.Zero(), 7, 8), 26)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		testdeep.All(
			testdeep.Empty(),
			5,
			testdeep.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			testdeep.All(6, 7),
			testdeep.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		42)

	// Only one interface type
	equalTypes(t,
		testdeep.All(
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*error)(nil)),
		),
		(*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		testdeep.All(
			testdeep.Isa((*error)(nil)),
			testdeep.Isa((*fmt.Stringer)(nil)),
		),
		nil)
}
