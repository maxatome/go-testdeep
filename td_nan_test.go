// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"math"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestNaN(t *testing.T) {
	checkOK(t, math.NaN(), NaN())
	checkOK(t, float32(math.NaN()), NaN())

	checkError(t, float32(12), NaN(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float32) 12"),
		Expected: mustBe("NaN"),
	})
	checkError(t, float64(12), NaN(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 12"),
		Expected: mustBe("NaN"),
	})

	checkError(t, 12, NaN(), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("float32 OR float64"),
	})
}

func TestNotNaN(t *testing.T) {
	checkOK(t, float64(12), NotNaN())
	checkOK(t, float32(12), NotNaN())

	checkError(t, float32(math.NaN()), NotNaN(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float32) NaN"),
		Expected: mustBe("not NaN"),
	})
	checkError(t, math.NaN(), NotNaN(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) NaN"),
		Expected: mustBe("not NaN"),
	})

	checkError(t, 12, NotNaN(), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("int"),
		Expected: mustBe("float32 OR float64"),
	})
}

func TestNaNTypeBehind(t *testing.T) {
	equalTypes(t, NaN(), nil)
	equalTypes(t, NotNaN(), nil)
}
