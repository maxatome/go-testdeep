// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"math"
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestNaN(t *testing.T) {
	checkOK(t, math.NaN(), testdeep.NaN())
	checkOK(t, float32(math.NaN()), testdeep.NaN())

	checkError(t, float32(12), testdeep.NaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float32) 12"),
			Expected: mustBe("NaN"),
		})
	checkError(t, float64(12), testdeep.NaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float64) 12"),
			Expected: mustBe("NaN"),
		})

	checkError(t, 12, testdeep.NaN(),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float32 OR float64"),
		})
}

func TestNotNaN(t *testing.T) {
	checkOK(t, float64(12), testdeep.NotNaN())
	checkOK(t, float32(12), testdeep.NotNaN())

	checkError(t, float32(math.NaN()), testdeep.NotNaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float32) NaN"),
			Expected: mustBe("not NaN"),
		})
	checkError(t, math.NaN(), testdeep.NotNaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float64) NaN"),
			Expected: mustBe("not NaN"),
		})

	checkError(t, 12, testdeep.NotNaN(),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float32 OR float64"),
		})
}

func TestNaNTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.NaN(), nil)
	equalTypes(t, testdeep.NotNaN(), nil)
}
