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

func TestAny(t *testing.T) {
	checkOK(t, 6, Any(nil, 5, 6, 7))
	checkOK(t, nil, Any(5, 6, 7, nil))

	checkError(t, 6, Any(5), expectedError{
		Message:  mustBe("comparing with Any"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("Any((int) 5)"),
	})

	checkError(t, 6, Any(nil), expectedError{
		Message:  mustBe("comparing with Any"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 6"),
		Expected: mustBe("Any(nil)"),
	})

	checkError(t, nil, Any(6), expectedError{
		Message:  mustBe("comparing with Any"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("Any((int) 6)"),
	})

	//
	// String
	equalStr(t, Any(6).String(), "Any((int) 6)")
	equalStr(t, Any(6, 7).String(), "Any((int) 6,\n    (int) 7)")
}
