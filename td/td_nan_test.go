// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"math"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestNaN(t *testing.T) {
	checkOK(t, math.NaN(), td.NaN())
	checkOK(t, float32(math.NaN()), td.NaN())

	checkError(t, float32(12), td.NaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float32) 12"),
			Expected: mustBe("NaN"),
		})
	checkError(t, float64(12), td.NaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12.0"),
			Expected: mustBe("NaN"),
		})

	checkError(t, 12, td.NaN(),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float32 OR float64"),
		})
}

func TestNotNaN(t *testing.T) {
	checkOK(t, float64(12), td.NotNaN())
	checkOK(t, float32(12), td.NotNaN())

	checkError(t, float32(math.NaN()), td.NotNaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(float32) NaN"),
			Expected: mustBe("not NaN"),
		})
	checkError(t, math.NaN(), td.NotNaN(),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("NaN"),
			Expected: mustBe("not NaN"),
		})

	checkError(t, 12, td.NotNaN(),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float32 OR float64"),
		})
}

func TestNaNTypeBehind(t *testing.T) {
	equalTypes(t, td.NaN(), nil)
	equalTypes(t, td.NotNaN(), nil)
}
