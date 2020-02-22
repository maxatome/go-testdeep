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

func TestAll(t *testing.T) {
	checkOK(t, 6, td.All(6, 6, 6))
	checkOK(t, nil, td.All(nil, nil, nil))

	checkError(t, 6, td.All(6, 5, 6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("5"),
		})

	checkError(t, 6, td.All(6, nil, 6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("nil"),
		})

	checkError(t, nil, td.All(nil, 5, nil),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("5"),
		})

	checkError(t,
		6,
		td.All(
			6,
			td.All(td.Between(3, 8), td.Between(4, 5)),
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
	test.EqualStr(t, td.All(6).String(), "All(6)")
	test.EqualStr(t, td.All(6, 7).String(), "All(6,\n    7)")
}

func TestAllTypeBehind(t *testing.T) {
	equalTypes(t, td.All(6, nil), nil)
	equalTypes(t, td.All(6, "toto"), nil)

	equalTypes(t, td.All(6, td.Zero(), 7, 8), 26)

	// Always the same non-interface type (even if we encounter several
	// interface types)
	equalTypes(t,
		td.All(
			td.Empty(),
			5,
			td.Isa((*error)(nil)), // interface type (in fact pointer to ...)
			td.All(6, 7),
			td.Isa((*fmt.Stringer)(nil)), // interface type
			8),
		42)

	// Only one interface type
	equalTypes(t,
		td.All(
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
			td.Isa((*error)(nil)),
		),
		(*error)(nil))

	// Several interface types, cannot be sure
	equalTypes(t,
		td.All(
			td.Isa((*error)(nil)),
			td.Isa((*fmt.Stringer)(nil)),
		),
		nil)
}
