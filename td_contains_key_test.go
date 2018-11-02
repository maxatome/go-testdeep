// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestContainsKey(t *testing.T) {
	type MyMap map[int]string

	for idx, got := range []interface{}{
		map[int]string{12: "foo", 34: "bar", 28: "zip"},
		MyMap{12: "foo", 34: "bar", 28: "zip"},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.ContainsKey(34), testName)
		checkOK(t, got, testdeep.ContainsKey(testdeep.Between(30, 35)), testName)

		checkError(t, got, testdeep.ContainsKey(35),
			expectedError{
				Message: mustBe("does not contain key"),
				Path:    mustBe("DATA"),
				Summary: mustMatch(`expected key: \(int\) 35
 not in keys: \(\(int\) (12|28|34),
               \(int\) (12|28|34),
               \(int\) (12|28|34)\)`),
			}, testName)
	}
}

// nil case
func TestContainsKeyNil(t *testing.T) {
	type MyPtrMap map[*int]int

	num := 12345642
	for idx, got := range []interface{}{
		map[*int]int{&num: 42, nil: 666},
		MyPtrMap{&num: 42, nil: 666},
	} {
		testName := fmt.Sprintf("#%d: got=%v", idx, got)

		checkOK(t, got, testdeep.ContainsKey(nil), testName)
		checkOK(t, got, testdeep.ContainsKey((*int)(nil)), testName)
		checkOK(t, got, testdeep.ContainsKey(testdeep.Nil()), testName)
		checkOK(t, got, testdeep.ContainsKey(testdeep.NotNil()), testName)

		checkError(t, got, testdeep.ContainsKey((*uint8)(nil)),
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
		testdeep.ContainsKey(nil),
		expectedError{
			Message: mustBe("does not contain key"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`expected key: nil
 not in keys: \("(foo|bar|zip)",
               "(foo|bar|zip)",
               "(foo|bar|zip)"\)`),
		})

	checkError(t, "foobar", testdeep.ContainsKey(nil),
		expectedError{
			Message:  mustBe("cannot check contains key"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("ContainsKey(nil)"),
		})

	checkError(t, "foobar", testdeep.ContainsKey(123),
		expectedError{
			Message:  mustBe("cannot check contains key"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("int"),
		})

	// Caught by deepValueEqual, before Match() call
	checkError(t, nil, testdeep.ContainsKey(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("ContainsKey(nil)"),
		})
}

func TestContainsKeyTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.ContainsKey("x"), nil)
}
