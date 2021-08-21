// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestCode(t *testing.T) {
	checkOK(t, 12, td.Code(func(n int) bool { return n >= 10 && n < 20 }))

	checkOK(t, 12, td.Code(func(val interface{}) bool {
		num, ok := val.(int)
		return ok && num == 12
	}))

	checkOK(t, errors.New("foobar"), td.Code(func(val error) bool {
		return val.Error() == "foobar"
	}))

	checkOK(t, json.RawMessage(`[42]`),
		td.Code(func(b json.RawMessage) error {
			var l []int
			err := json.Unmarshal(b, &l)
			if err != nil {
				return err
			}
			if len(l) != 1 || l[0] != 42 {
				return errors.New("42 not found")
			}
			return nil
		}))

	// Lax
	checkOK(t, 123, td.Lax(td.Code(func(n float64) bool { return n == 123 })))

	checkError(t, 123, td.Code(func(n float64) bool { return true }),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float64"),
		})

	type xInt int
	checkError(t, xInt(12),
		td.Code(func(n int) bool { return n >= 10 && n < 20 }),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.xInt"),
			Expected: mustBe("int"),
		})

	checkError(t, 12,
		td.Code(func(n int) (bool, string) { return false, "custom error" }),
		expectedError{
			Message: mustBe("ran code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: custom error"),
		})

	checkError(t, 12,
		td.Code(func(n int) bool { return false }),
		expectedError{
			Message: mustBe("ran code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("  value: 12\nit failed but didn't say why"),
		})

	type MyBool bool
	type MyString string
	checkError(t, 12,
		td.Code(func(n int) (MyBool, MyString) { return false, "very custom error" }),
		expectedError{
			Message: mustBe("ran code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: very custom error"),
		})

	checkError(t, 12,
		td.Code(func(i int) error {
			return errors.New("very custom error")
		}),
		expectedError{
			Message: mustBe("ran code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: very custom error"),
		})

	// Internal use
	checkError(t, 12,
		td.Code(func(i int) error {
			return &ctxerr.Error{
				Message: "my message",
				Summary: ctxerr.NewSummary("my summary"),
			}
		}),
		expectedError{
			Message: mustBe("my message"),
			Path:    mustBe("DATA"),
			Summary: mustBe("my summary"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Code(nil),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Code(FUNC), but received nil as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Code((func(string) bool)(nil)),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC cannot be a nil function"),
		})

	checkError(t, "never tested",
		td.Code("test"),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Code(FUNC), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Code(func() bool { return true }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument"),
		})

	checkError(t, "never tested",
		td.Code(func(x ...int) bool { return true }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument"),
		})

	checkError(t, "never tested",
		td.Code(func(a int, b string) bool { return true }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (bool, int) { return true, 0 }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (error, string) { return nil, "" }), // nolint: staticcheck
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (int, string) { return 0, "" }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (string, bool) { return "", true }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (bool, string, int) { return true, "", 0 }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) {}),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) int { return 0 }),
		expectedError{
			Message: mustBe("Bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	//
	// String
	test.EqualStr(t,
		td.Code(func(n int) bool { return false }).String(),
		"Code(func(int) bool)")
	test.EqualStr(t,
		td.Code(func(n int) (bool, string) { return false, "" }).String(),
		"Code(func(int) (bool, string))")
	test.EqualStr(t,
		td.Code(func(n int) error { return nil }).String(),
		"Code(func(int) error)")
	test.EqualStr(t,
		td.Code(func(n int) (MyBool, MyString) { return false, "" }).String(),
		"Code(func(int) (td_test.MyBool, td_test.MyString))")

	// Erroneous op
	test.EqualStr(t, td.Code(nil).String(), "Code(<ERROR>)")
}

func TestCodeTypeBehind(t *testing.T) {
	// Type behind is the code function parameter one

	equalTypes(t, td.Code(func(n int) bool { return n != 0 }), 23)

	type MyTime time.Time

	equalTypes(t,
		td.Code(func(t MyTime) bool { return time.Time(t).IsZero() }),
		MyTime{})

	// Erroneous op
	equalTypes(t, td.Code(nil), nil)
}
