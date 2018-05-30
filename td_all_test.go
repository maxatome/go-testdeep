// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestAll(t *testing.T) {
	checkOK(t, 6, All(6, 6, 6))
	checkOK(t, nil, All(nil, nil, nil))

	checkError(t, 6, All(6, 5, 6), expectedError{
		Message:  mustBe("compared (part 2 of 3)"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("(int) 5"),
	})

	checkError(t, 6, All(6, nil, 6), expectedError{
		Message:  mustBe("compared (part 2 of 3)"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("nil"),
	})

	checkError(t, nil, All(nil, 5, nil), expectedError{
		Message:  mustBe("compared (part 2 of 3)"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("(int) 5"),
	})

	checkError(t, 6, All(6, All(Between(3, 8), Between(4, 5)), 6),
		expectedError{
			Message:  mustBe("compared (part 2 of 3)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(int) 6"),
			Expected: mustBe("All(3 ≤ got ≤ 8,\n    4 ≤ got ≤ 5)"),
			Origin: &expectedError{
				Message:  mustBe("compared (part 2 of 2)"),
				Path:     mustBe("DATA<All#2/3>"),
				Got:      mustBe("(int) 6"),
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
	equalStr(t, All(6).String(), "All((int) 6)")
	equalStr(t, All(6, 7).String(), "All((int) 6,\n    (int) 7)")
}

func TestAllTypeOf(t *testing.T) {
	equalTypes(t, All(6), nil)
}
