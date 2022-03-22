// Copyright (c) 2019, Maxime Soulé
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

func TestCatch(t *testing.T) {
	var num int
	checkOK(t, 12, td.Catch(&num, 12))
	test.EqualInt(t, num, 12)

	var num64 int64
	checkError(t, 12, td.Catch(&num64, 12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Got:      mustBe("int"),
			Expected: mustBe("int64"),
		})

	checkOK(t, 12, td.Lax(td.Catch(&num64, 12)))
	test.EqualInt(t, int(num64), 12)

	// Lax not needed for interfaces
	var val any
	if checkOK(t, 12, td.Catch(&val, 12)) {
		if n, ok := val.(int); ok {
			test.EqualInt(t, n, 12)
		} else {
			t.Errorf("val is not an int but a %T", val)
		}
	}

	//
	// Bad usages
	checkError(t, "never tested",
		td.Catch(12, 28),
		expectedError{
			Message: mustBe("bad usage of Catch operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Catch(NON_NIL_PTR, EXPECTED_VALUE), but received int as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Catch(nil, 28),
		expectedError{
			Message: mustBe("bad usage of Catch operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Catch(NON_NIL_PTR, EXPECTED_VALUE), but received nil as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Catch((*int)(nil), 28),
		expectedError{
			Message: mustBe("bad usage of Catch operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Catch(NON_NIL_PTR, EXPECTED_VALUE), but received *int (ptr) as 1st parameter"),
		})

	//
	// String
	test.EqualStr(t, td.Catch(&num, 12).String(), "12")
	test.EqualStr(t,
		td.Catch(&num, td.Gt(4)).String(),
		td.Gt(4).String())
	test.EqualStr(t, td.Catch(&num, nil).String(), "nil")

	// Erroneous op
	test.EqualStr(t, td.Catch(nil, 28).String(), "Catch(<ERROR>)")
}

func TestCatchTypeBehind(t *testing.T) {
	var num int
	equalTypes(t, td.Catch(&num, 8), 0)
	equalTypes(t, td.Catch(&num, td.Gt(4)), 0)
	equalTypes(t, td.Catch(&num, td.Ignore()), 0) // fallback on *target
	equalTypes(t, td.Catch(&num, nil), nil)

	// Erroneous op
	equalTypes(t, td.Catch(nil, 28), nil)
}
