// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestContainsKey(t *testing.T) {
	type MyMap map[int]string

	for idx, got := range []any{
		map[int]string{12: "foo", 34: "bar", 28: "zip"},
		MyMap{12: "foo", 34: "bar", 28: "zip"},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.ContainsKey(34), testName)
		checkOK(t, got, td.ContainsKey(td.Between(30, 35)),
			testName)

		checkError(t, got, td.ContainsKey(35),
			expectedError{
				Message: mustBe("does not contain key"),
				Path:    mustBe("DATA"),
				Summary: mustMatch(`expected key: 35
 not in keys: \((12|28|34),
               (12|28|34),
               (12|28|34)\)`),
			}, testName)

		// Lax
		checkOK(t, got, td.Lax(td.ContainsKey(float64(34))), testName)
	}
}

// nil case.
func TestContainsKeyNil(t *testing.T) {
	type MyPtrMap map[*int]int

	num := 12345642
	for idx, got := range []any{
		map[*int]int{&num: 42, nil: 666},
		MyPtrMap{&num: 42, nil: 666},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, td.ContainsKey(nil), testName)
		checkOK(t, got, td.ContainsKey((*int)(nil)), testName)
		checkOK(t, got, td.ContainsKey(td.Nil()), testName)
		checkOK(t, got, td.ContainsKey(td.NotNil()), testName)

		checkError(t, got, td.ContainsKey((*uint8)(nil)),
			expectedError{
				Message: mustBe("does not contain key"),
				Path:    mustBe("DATA"),
				Summary: mustMatch(`expected key: \(\*uint8\)\(<nil>\)
 not in keys: \(\(\*int\)\((<nil>|.*12345642.*)\),
               \(\*int\)\((<nil>|.*12345642.*)\)\)`),
			}, testName)
	}

	checkError(t,
		map[string]int{"foo": 12, "bar": 34, "zip": 28}, // got
		td.ContainsKey(nil),
		expectedError{
			Message: mustBe("does not contain key"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`expected key: nil
 not in keys: \("(foo|bar|zip)",
               "(foo|bar|zip)",
               "(foo|bar|zip)"\)`),
		})

	checkError(t, "foobar", td.ContainsKey(nil),
		expectedError{
			Message:  mustBe("cannot check contains key"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("ContainsKey(nil)"),
		})

	checkError(t, "foobar", td.ContainsKey(123),
		expectedError{
			Message:  mustBe("cannot check contains key"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("int"),
		})

	// Caught by deepValueEqual, before Match() call
	checkError(t, nil, td.ContainsKey(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("ContainsKey(nil)"),
		})
}

func TestContainsKeyTypeBehind(t *testing.T) {
	equalTypes(t, td.ContainsKey("x"), nil)
}
