// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestNone(t *testing.T) {
	checkOK(t, 6, testdeep.None(7, 8, 9, nil))
	checkOK(t, nil, testdeep.None(7, 8, 9))

	checkError(t, 6, testdeep.None(6, 7), expectedError{
		Message:  mustBe("comparing with None (part 1 of 2 is OK)"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("None((int) 6,\n     (int) 7)"),
	})

	checkError(t, nil, testdeep.None(7, nil), expectedError{
		Message:  mustBe("comparing with None (part 2 of 2 is OK)"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("None((int) 7,\n     nil)"),
	})

	//
	// String
	test.EqualStr(t, testdeep.None(6).String(), "None((int) 6)")
	test.EqualStr(t, testdeep.None(6, 7).String(), "None((int) 6,\n     (int) 7)")
}

func TestNot(t *testing.T) {
	checkOK(t, 6, testdeep.Not(7))
	checkOK(t, nil, testdeep.Not(7))

	checkError(t, 6, testdeep.Not(6), expectedError{
		Message:  mustBe("comparing with Not"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("Not((int) 6)"),
	})

	checkError(t, nil, testdeep.Not(nil), expectedError{
		Message:  mustBe("comparing with Not"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("Not(nil)"),
	})

	//
	// String
	test.EqualStr(t, testdeep.Not(6).String(), "Not((int) 6)")
}

func TestNoneTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.None(6), nil)
	equalTypes(t, testdeep.Not(6), nil)
}
