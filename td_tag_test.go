// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/util"
)

func TestTag(t *testing.T) {
	// expected value
	checkOK(t, 12, testdeep.Tag("number", 12))
	checkOK(t, nil, testdeep.Tag("number", nil))

	checkError(t, 8, testdeep.Tag("number", 9),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("8"),
			Expected: mustBe("9"),
		})

	// expected operator
	checkOK(t, 12, testdeep.Tag("number", testdeep.Between(9, 13)))

	checkError(t, 8, testdeep.Tag("number", testdeep.Between(9, 13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("8"),
			Expected: mustBe("9 ≤ got ≤ 13"),
		})

	//
	// Bad usage
	test.CheckPanic(t,
		func() { testdeep.Tag("1badTag", testdeep.Between(9, 13)) },
		util.ErrTagInvalid.Error())

	//
	// String
	test.EqualStr(t,
		testdeep.Tag("foo", testdeep.Gt(4)).String(),
		testdeep.Gt(4).String())
	test.EqualStr(t, testdeep.Tag("foo", 8).String(), "8")
	test.EqualStr(t, testdeep.Tag("foo", nil).String(), "nil")
}

func TestTagTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Tag("foo", 8), 0)
	equalTypes(t, testdeep.Tag("foo", testdeep.Gt(4)), 0)
	equalTypes(t, testdeep.Tag("foo", nil), nil)
}
