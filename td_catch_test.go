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
	//	"github.com/maxatome/go-testdeep/internal/util"
)

func TestCatch(t *testing.T) {
	var num int
	checkOK(t, 12, testdeep.Catch(&num, 12))
	test.EqualInt(t, num, 12)

	var num64 int64
	checkError(t, 12, testdeep.Catch(&num64, 12),
		expectedError{
			Message:  mustBe("type mismatch"),
			Got:      mustBe("int"),
			Expected: mustBe("int64"),
		})

	checkOK(t, 12, testdeep.Lax(testdeep.Catch(&num64, 12)))
	test.EqualInt(t, int(num64), 12)

	// Lax not needed for interfaces
	var any interface{}
	if checkOK(t, 12, testdeep.Catch(&any, 12)) {
		if n, ok := any.(int); ok {
			test.EqualInt(t, n, 12)
		} else {
			t.Errorf("any is not an int but a %T", any)
		}
	}

	//
	// Bad usages
	test.CheckPanic(t, func() { testdeep.Catch(12, 28) },
		"usage: Catch(NON_NIL_PTR, EXPECTED_VALUE)")
	test.CheckPanic(t, func() { testdeep.Catch(nil, 28) },
		"usage: Catch(NON_NIL_PTR, EXPECTED_VALUE)")
	test.CheckPanic(t, func() { testdeep.Catch((*int)(nil), 28) },
		"usage: Catch(NON_NIL_PTR, EXPECTED_VALUE)")

	//
	// String
	test.EqualStr(t, testdeep.Catch(&num, 12).String(), "12")
	test.EqualStr(t,
		testdeep.Catch(&num, testdeep.Gt(4)).String(),
		testdeep.Gt(4).String())
	test.EqualStr(t, testdeep.Catch(&num, nil).String(), "nil")
}

func TestCatchTypeBehind(t *testing.T) {
	var num int
	equalTypes(t, testdeep.Catch(&num, 8), 0)
	equalTypes(t, testdeep.Catch(&num, testdeep.Gt(4)), 0)
	equalTypes(t, testdeep.Catch(&num, nil), nil)
}
