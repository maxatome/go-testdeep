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

func TestKeysValues(t *testing.T) {
	var m map[string]int

	//
	t.Run("nil map", func(t *testing.T) {
		checkOK(t, m, td.Keys([]string{}))
		checkOK(t, m, td.Values([]int{}))

		checkOK(t, m, td.Keys(td.Empty()))
		checkOK(t, m, td.Values(td.Empty()))

		checkError(t, m, td.Keys(td.NotEmpty()),
			expectedError{
				Message:  mustBe("empty"),
				Path:     mustBe("keys(DATA)"),
				Expected: mustBe("not empty"),
			})
		checkError(t, m, td.Values(td.NotEmpty()),
			expectedError{
				Message:  mustBe("empty"),
				Path:     mustBe("values(DATA)"),
				Expected: mustBe("not empty"),
			})
	})

	//
	t.Run("non-nil but empty map", func(t *testing.T) {
		m = map[string]int{}
		checkOK(t, m, td.Keys([]string{}))
		checkOK(t, m, td.Values([]int{}))

		checkOK(t, m, td.Keys(td.Empty()))
		checkOK(t, m, td.Values(td.Empty()))

		checkError(t, m, td.Keys(td.NotEmpty()),
			expectedError{
				Message:  mustBe("empty"),
				Path:     mustBe("keys(DATA)"),
				Expected: mustBe("not empty"),
			})
		checkError(t, m, td.Values(td.NotEmpty()),
			expectedError{
				Message:  mustBe("empty"),
				Path:     mustBe("values(DATA)"),
				Expected: mustBe("not empty"),
			})
	})

	//
	t.Run("Filled map", func(t *testing.T) {
		m = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6}
		checkOK(t, m, td.Keys([]string{"a", "b", "c", "d", "e", "f"}))
		checkOK(t, m, td.Values([]int{1, 2, 3, 4, 5, 6}))

		checkOK(t, m, td.Keys(td.Bag("a", "b", "c", "d", "e", "f")))
		checkOK(t, m, td.Values(td.Bag(1, 2, 3, 4, 5, 6)))

		checkOK(t, m, td.Keys(td.ArrayEach(td.Between("a", "f"))))
		checkOK(t, m, td.Values(td.ArrayEach(td.Between(1, 6))))

		checkError(t, m, td.Keys(td.Empty()),
			expectedError{
				Message:  mustBe("not empty"),
				Path:     mustBe("keys(DATA)"),
				Expected: mustBe("empty"),
			})
		checkError(t, m, td.Values(td.Empty()),
			expectedError{
				Message:  mustBe("not empty"),
				Path:     mustBe("values(DATA)"),
				Expected: mustBe("empty"),
			})
	})

	//
	t.Run("Errors", func(t *testing.T) {
		checkError(t, nil, td.Keys([]int{1, 2, 3}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustContain("keys=([]int)"),
			})
		checkError(t, nil, td.Values([]int{1, 2, 3}),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustContain("values=([]int)"),
			})

		checkError(t, nil, td.Keys(td.Empty()),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("keys: Empty()"),
			})
		checkError(t, nil, td.Values(td.Empty()),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("values: Empty()"),
			})

		checkError(t, 123, td.Keys(td.Empty()),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("map"),
			})
		checkError(t, 123, td.Values(td.Empty()),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("map"),
			})
	})

	//
	t.Run("Bad usage", func(t *testing.T) {
		checkError(t, "never tested",
			td.Keys(12),
			expectedError{
				Message: mustBe("Bad usage of Keys operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Keys(TESTDEEP_OPERATOR|SLICE), but received int as 1st parameter"),
			})

		checkError(t, "never tested",
			td.Values(12),
			expectedError{
				Message: mustBe("Bad usage of Values operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Values(TESTDEEP_OPERATOR|SLICE), but received int as 1st parameter"),
			})
	})

	// Erroneous op
	test.EqualStr(t, td.Keys(12).String(), "Keys(<ERROR>)")
	test.EqualStr(t, td.Values(12).String(), "Values(<ERROR>)")
}

func TestKeysValuesTypeBehind(t *testing.T) {
	equalTypes(t, td.Keys([]string{}), nil)
	equalTypes(t, td.Values([]string{}), nil)

	// Erroneous op
	equalTypes(t, td.Keys(12), nil)
	equalTypes(t, td.Values(12), nil)
}
