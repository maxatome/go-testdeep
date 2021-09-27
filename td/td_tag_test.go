// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/util"
	"github.com/maxatome/go-testdeep/td"
)

func TestTag(t *testing.T) {
	// expected value
	checkOK(t, 12, td.Tag("number", 12))
	checkOK(t, nil, td.Tag("number", nil))

	checkError(t, 8, td.Tag("number", 9),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("8"),
			Expected: mustBe("9"),
		})

	// expected operator
	checkOK(t, 12, td.Tag("number", td.Between(9, 13)))

	checkError(t, 8, td.Tag("number", td.Between(9, 13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("8"),
			Expected: mustBe("9 ≤ got ≤ 13"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Tag("1badTag", td.Between(9, 13)),
		expectedError{
			Message: mustBe("bad usage of Tag operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(util.ErrTagInvalid.Error()),
		})

	//
	// String
	test.EqualStr(t,
		td.Tag("foo", td.Gt(4)).String(),
		td.Gt(4).String())
	test.EqualStr(t, td.Tag("foo", 8).String(), "8")
	test.EqualStr(t, td.Tag("foo", nil).String(), "nil")

	// Erroneous op
	test.EqualStr(t, td.Tag("1badTag", 12).String(), "Tag(<ERROR>)")
}

func TestTagTypeBehind(t *testing.T) {
	equalTypes(t, td.Tag("foo", 8), 0)
	equalTypes(t, td.Tag("foo", td.Gt(4)), 0)
	equalTypes(t, td.Tag("foo", nil), nil)

	// Erroneous op
	equalTypes(t, td.Tag("1badTag", 12), nil)
}
