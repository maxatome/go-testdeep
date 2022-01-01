// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestNone(t *testing.T) {
	checkOK(t, 6, td.None(7, 8, 9, nil))
	checkOK(t, nil, td.None(7, 8, 9))

	checkError(t, 6, td.None(6, 7),
		expectedError{
			Message:  mustBe("comparing with None (part 1 of 2 is OK)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("None(6,\n     7)"),
		})

	checkError(t, nil, td.None(7, nil),
		expectedError{
			Message:  mustBe("comparing with None (part 2 of 2 is OK)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("None(7,\n     nil)"),
		})

	// Lax
	checkError(t, float64(6), td.Lax(td.None(6, 7)),
		expectedError{
			Message:  mustBe("comparing with None (part 1 of 2 is OK)"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6.0"),
			Expected: mustBe("None(6,\n     7)"),
		})

	//
	// String
	test.EqualStr(t, td.None(6).String(), "None(6)")
	test.EqualStr(t, td.None(6, 7).String(), "None(6,\n     7)")
}

func TestNot(t *testing.T) {
	checkOK(t, 6, td.Not(7))
	checkOK(t, nil, td.Not(7))

	checkError(t, 6, td.Not(6),
		expectedError{
			Message:  mustBe("comparing with Not"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Not(6)"),
		})

	checkError(t, nil, td.Not(nil),
		expectedError{
			Message:  mustBe("comparing with Not"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Not(nil)"),
		})

	//
	// String
	test.EqualStr(t, td.Not(6).String(), "Not(6)")
}

func TestNoneTypeBehind(t *testing.T) {
	equalTypes(t, td.None(6), nil)
	equalTypes(t, td.Not(6), nil)
}
