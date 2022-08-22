// Copyright (c) 2022, Maxime Soulé
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

func TestGrep(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := [...]int{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Grep(td.Gt(0), []int{1, 2, 3}))
				checkOK(t, tc.got, td.Grep(td.Not(td.Between(-2, 2)), []int{-3, 3}))

				checkOK(t, tc.got, td.Grep(
					func(x int) bool { return (x & 1) != 0 },
					[]int{-3, -1, 1, 3}))

				checkOK(t, tc.got, td.Grep(
					func(x int64) bool { return (x & 1) != 0 },
					[]int{-3, -1, 1, 3}),
					"int64 filter vs int items")

				checkOK(t, tc.got, td.Grep(
					func(x any) bool { return (x.(int) & 1) != 0 },
					[]int{-3, -1, 1, 3}),
					"any filter vs int items")
			})
		}
	})

	t.Run("struct", func(t *testing.T) {
		type person struct {
			ID   int64
			Name string
		}
		got := [...]person{
			{ID: 1, Name: "Joe"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Alice"},
			{ID: 4, Name: "Brian"},
			{ID: 5, Name: "Britt"},
		}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Grep(
					td.JSONPointer("/Name", td.HasPrefix("Br")),
					[]person{{ID: 4, Name: "Brian"}, {ID: 5, Name: "Britt"}}))

				checkOK(t, tc.got, td.Grep(
					func(p person) bool { return p.ID < 3 },
					[]person{{ID: 1, Name: "Joe"}, {ID: 2, Name: "Bob"}}))
			})
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		got := [...]any{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Grep(td.Gt(0), []any{1, 2, 3}))
				checkOK(t, tc.got, td.Grep(td.Not(td.Between(-2, 2)), []any{-3, 3}))

				checkOK(t, tc.got, td.Grep(
					func(x int) bool { return (x & 1) != 0 },
					[]any{-3, -1, 1, 3}))

				checkOK(t, tc.got, td.Grep(
					func(x int64) bool { return (x & 1) != 0 },
					[]any{-3, -1, 1, 3}),
					"int64 filter vs any/int items")

				checkOK(t, tc.got, td.Grep(
					func(x any) bool { return (x.(int) & 1) != 0 },
					[]any{-3, -1, 1, 3}),
					"any filter vs any/int items")
			})
		}
	})

	t.Run("interfaces error", func(t *testing.T) {
		got := [...]any{123, "foo"}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkError(t, tc.got,
					td.Grep(func(x int) bool { return true }, []string{"never reached"}),
					expectedError{
						Message:  mustBe("incompatible parameter type"),
						Path:     mustBe("DATA[1]"),
						Got:      mustBe("string"),
						Expected: mustBe("int"),
					})
			})
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var got []int
		testCases := []struct {
			name string
			got  any
		}{
			{"slice", got},
			{"*slice", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Grep(td.Gt(666), ([]int)(nil)))
			})
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		checkError(t, (*[]int)(nil), td.Grep(td.Ignore(), []int{33}),
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *slice (*[]int type)"),
				Expected: mustBe("non-nil *slice OR *array"),
			})
	})

	t.Run("JSON", func(t *testing.T) {
		got := map[string]any{
			"values": []int{1, 2, 3, 4},
		}
		checkOK(t, got, td.JSON(`{"values": Grep(Gt(2), [3, 4])}`))
	})

	t.Run("errors", func(t *testing.T) {
		for _, filter := range []any{nil, 33} {
			checkError(t, "never tested",
				td.Grep(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Grep operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Grep(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must be a function or FILTER_TESTDEEP_OPERATOR a TestDeep operator"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func() bool { return true },
			func(a, b int) bool { return true },
			func(a ...int) bool { return true },
		} {
			checkError(t, "never tested",
				td.Grep(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Grep operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Grep(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must take only one non-variadic argument"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func(a int) {},
			func(a int) int { return 0 },
			func(a int) (bool, bool) { return true, true },
		} {
			checkError(t, "never tested",
				td.Grep(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Grep operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Grep(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must return bool"),
				},
				"filter:", filter)
		}

		checkError(t, "never tested", td.Grep(td.Ignore(), 42),
			expectedError{
				Message: mustBe("bad usage of Grep operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Grep(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), EXPECTED_VALUE must be a slice not a int"),
			})

		checkError(t, &struct{}{}, td.Grep(td.Ignore(), []int{33}),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*struct (*struct {} type)"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})

		checkError(t, nil, td.Grep(td.Ignore(), []int{33}),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})
	})
}

func TestGrepTypeBehind(t *testing.T) {
	equalTypes(t, td.Grep(func(n int) bool { return true }, []int{33}), []int{})
	equalTypes(t, td.Grep(td.Gt("0"), []string{"33"}), []string{})

	// Erroneous op
	equalTypes(t, td.Grep(42, 33), nil)
}

func TestGrepString(t *testing.T) {
	test.EqualStr(t,
		td.Grep(func(n int) bool { return true }, []int{}).String(),
		"Grep(func(int) bool)")

	test.EqualStr(t, td.Grep(td.Gt(0), []int{}).String(), "Grep(> 0)")

	// Erroneous op
	test.EqualStr(t, td.Grep(42, []int{}).String(), "Grep(<ERROR>)")
}

func TestFirst(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := [...]int{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.First(td.Gt(0), 1))
				checkOK(t, tc.got, td.First(td.Not(td.Between(-3, 2)), 3))

				checkOK(t, tc.got, td.First(
					func(x int) bool { return (x & 1) == 0 },
					-2))

				checkOK(t, tc.got, td.First(
					func(x int64) bool { return (x & 1) != 0 },
					-3),
					"int64 filter vs int items")

				checkOK(t, tc.got, td.First(
					func(x any) bool { return (x.(int) & 1) == 0 },
					-2),
					"any filter vs int items")

				checkError(t, tc.got,
					td.First(td.Gt(666), "never reached"),
					expectedError{
						Message:  mustBe("item not found"),
						Path:     mustBe("DATA"),
						Got:      mustContain(`]int) (len=7 `),
						Expected: mustBe("First(> 666)"),
					})
			})
		}
	})

	t.Run("struct", func(t *testing.T) {
		type person struct {
			ID   int64
			Name string
		}
		got := [...]person{
			{ID: 1, Name: "Joe"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Alice"},
			{ID: 4, Name: "Brian"},
			{ID: 5, Name: "Britt"},
		}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.First(
					td.JSONPointer("/Name", td.HasPrefix("Br")),
					person{ID: 4, Name: "Brian"}))

				checkOK(t, tc.got, td.First(
					func(p person) bool { return p.ID < 3 },
					person{ID: 1, Name: "Joe"}))
			})
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		got := [...]any{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.First(td.Gt(0), 1))
				checkOK(t, tc.got, td.First(td.Not(td.Between(-3, 2)), 3))

				checkOK(t, tc.got, td.First(
					func(x int) bool { return (x & 1) == 0 },
					-2))

				checkOK(t, tc.got, td.First(
					func(x int64) bool { return (x & 1) != 0 },
					-3),
					"int64 filter vs any/int items")

				checkOK(t, tc.got, td.First(
					func(x any) bool { return (x.(int) & 1) == 0 },
					-2),
					"any filter vs any/int items")
			})
		}
	})

	t.Run("interfaces error", func(t *testing.T) {
		got := [...]any{123, "foo"}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkError(t, tc.got,
					td.First(func(x int) bool { return false }, "never reached"),
					expectedError{
						Message:  mustBe("incompatible parameter type"),
						Path:     mustBe("DATA[1]"),
						Got:      mustBe("string"),
						Expected: mustBe("int"),
					})
			})
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var got []int
		testCases := []struct {
			name string
			got  any
		}{
			{"slice", got},
			{"*slice", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkError(t, tc.got,
					td.First(td.Gt(666), "never reached"),
					expectedError{
						Message:  mustBe("item not found"),
						Path:     mustBe("DATA"),
						Got:      mustBe("([]int) <nil>"),
						Expected: mustBe("First(> 666)"),
					})
			})
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		checkError(t, (*[]int)(nil), td.First(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *slice (*[]int type)"),
				Expected: mustBe("non-nil *slice OR *array"),
			})
	})

	t.Run("JSON", func(t *testing.T) {
		got := map[string]any{
			"values": []int{1, 2, 3, 4},
		}
		checkOK(t, got, td.JSON(`{"values": First(Gt(2), 3)}`))
	})

	t.Run("errors", func(t *testing.T) {
		for _, filter := range []any{nil, 33} {
			checkError(t, "never tested",
				td.First(filter, 42),
				expectedError{
					Message: mustBe("bad usage of First operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: First(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must be a function or FILTER_TESTDEEP_OPERATOR a TestDeep operator"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func() bool { return true },
			func(a, b int) bool { return true },
			func(a ...int) bool { return true },
		} {
			checkError(t, "never tested",
				td.First(filter, 42),
				expectedError{
					Message: mustBe("bad usage of First operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: First(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must take only one non-variadic argument"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func(a int) {},
			func(a int) int { return 0 },
			func(a int) (bool, bool) { return true, true },
		} {
			checkError(t, "never tested",
				td.First(filter, 42),
				expectedError{
					Message: mustBe("bad usage of First operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: First(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must return bool"),
				},
				"filter:", filter)
		}

		checkError(t, &struct{}{}, td.First(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*struct (*struct {} type)"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})

		checkError(t, nil, td.First(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})
	})
}

func TestFirstString(t *testing.T) {
	test.EqualStr(t,
		td.First(func(n int) bool { return true }, 33).String(),
		"First(func(int) bool)")

	test.EqualStr(t, td.First(td.Gt(0), 33).String(), "First(> 0)")

	// Erroneous op
	test.EqualStr(t, td.First(42, 33).String(), "First(<ERROR>)")
}

func TestFirstTypeBehind(t *testing.T) {
	equalTypes(t, td.First(func(n int) bool { return true }, 33), []int{})
	equalTypes(t, td.First(td.Gt("x"), "x"), []string{})

	// Erroneous op
	equalTypes(t, td.First(42, 33), nil)
}

func TestLast(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := [...]int{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Last(td.Lt(0), -1))
				checkOK(t, tc.got, td.Last(td.Not(td.Between(1, 3)), 0))

				checkOK(t, tc.got, td.Last(
					func(x int) bool { return (x & 1) == 0 },
					2))

				checkOK(t, tc.got, td.Last(
					func(x int64) bool { return (x & 1) != 0 },
					3),
					"int64 filter vs int items")

				checkOK(t, tc.got, td.Last(
					func(x any) bool { return (x.(int) & 1) == 0 },
					2),
					"any filter vs int items")

				checkError(t, tc.got,
					td.Last(td.Gt(666), "never reached"),
					expectedError{
						Message:  mustBe("item not found"),
						Path:     mustBe("DATA"),
						Got:      mustContain(`]int) (len=7 `),
						Expected: mustBe("Last(> 666)"),
					})
			})
		}
	})

	t.Run("struct", func(t *testing.T) {
		type person struct {
			ID   int64
			Name string
		}
		got := [...]person{
			{ID: 1, Name: "Joe"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Alice"},
			{ID: 4, Name: "Brian"},
			{ID: 5, Name: "Britt"},
		}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Last(
					td.JSONPointer("/Name", td.HasPrefix("Br")),
					person{ID: 5, Name: "Britt"}))

				checkOK(t, tc.got, td.Last(
					func(p person) bool { return p.ID < 3 },
					person{ID: 2, Name: "Bob"}))
			})
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		got := [...]any{-3, -2, -1, 0, 1, 2, 3}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkOK(t, tc.got, td.Last(td.Lt(0), -1))
				checkOK(t, tc.got, td.Last(td.Not(td.Between(1, 3)), 0))

				checkOK(t, tc.got, td.Last(
					func(x int) bool { return (x & 1) == 0 },
					2))

				checkOK(t, tc.got, td.Last(
					func(x int64) bool { return (x & 1) != 0 },
					3),
					"int64 filter vs any/int items")

				checkOK(t, tc.got, td.Last(
					func(x any) bool { return (x.(int) & 1) == 0 },
					2),
					"any filter vs any/int items")
			})
		}
	})

	t.Run("interfaces error", func(t *testing.T) {
		got := [...]any{123, "foo", 456}
		sgot := got[:]

		testCases := []struct {
			name string
			got  any
		}{
			{"slice", sgot},
			{"array", got},
			{"*slice", &sgot},
			{"*array", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkError(t, tc.got,
					td.Last(func(x int) bool { return false }, "never reached"),
					expectedError{
						Message:  mustBe("incompatible parameter type"),
						Path:     mustBe("DATA[1]"),
						Got:      mustBe("string"),
						Expected: mustBe("int"),
					})
			})
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var got []int
		testCases := []struct {
			name string
			got  any
		}{
			{"slice", got},
			{"*slice", &got},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				checkError(t, tc.got,
					td.Last(td.Gt(666), "never reached"),
					expectedError{
						Message:  mustBe("item not found"),
						Path:     mustBe("DATA"),
						Got:      mustBe("([]int) <nil>"),
						Expected: mustBe("Last(> 666)"),
					})
			})
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		checkError(t, (*[]int)(nil), td.Last(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *slice (*[]int type)"),
				Expected: mustBe("non-nil *slice OR *array"),
			})
	})

	t.Run("JSON", func(t *testing.T) {
		got := map[string]any{
			"values": []int{1, 2, 3, 4},
		}
		checkOK(t, got, td.JSON(`{"values": Last(Lt(3), 2)}`))
	})

	t.Run("errors", func(t *testing.T) {
		for _, filter := range []any{nil, 33} {
			checkError(t, "never tested",
				td.Last(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Last operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Last(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must be a function or FILTER_TESTDEEP_OPERATOR a TestDeep operator"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func() bool { return true },
			func(a, b int) bool { return true },
			func(a ...int) bool { return true },
		} {
			checkError(t, "never tested",
				td.Last(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Last operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Last(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must take only one non-variadic argument"),
				},
				"filter:", filter)
		}

		for _, filter := range []any{
			func(a int) {},
			func(a int) int { return 0 },
			func(a int) (bool, bool) { return true, true },
		} {
			checkError(t, "never tested",
				td.Last(filter, 42),
				expectedError{
					Message: mustBe("bad usage of Last operator"),
					Path:    mustBe("DATA"),
					Summary: mustBe("usage: Last(FILTER_FUNC|FILTER_TESTDEEP_OPERATOR, TESTDEEP_OPERATOR|EXPECTED_VALUE), FILTER_FUNC must return bool"),
				},
				"filter:", filter)
		}

		checkError(t, &struct{}{}, td.Last(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*struct (*struct {} type)"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})

		checkError(t, nil, td.Last(td.Ignore(), 33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("slice OR array OR *slice OR *array"),
			})
	})
}

func TestLastString(t *testing.T) {
	test.EqualStr(t,
		td.Last(func(n int) bool { return true }, 33).String(),
		"Last(func(int) bool)")

	test.EqualStr(t, td.Last(td.Gt(0), 33).String(), "Last(> 0)")

	// Erroneous op
	test.EqualStr(t, td.Last(42, 33).String(), "Last(<ERROR>)")
}

func TestLastTypeBehind(t *testing.T) {
	equalTypes(t, td.Last(func(n int) bool { return true }, 33), []int{})
	equalTypes(t, td.Last(td.Gt("x"), "x"), []string{})

	// Erroneous op
	equalTypes(t, td.Last(42, 33), nil)
}
