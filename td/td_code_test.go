// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestCode(t *testing.T) {
	checkOK(t, 12, td.Code(func(n int) bool { return n >= 10 && n < 20 }))

	checkOK(t, 12, td.Code(func(val any) bool {
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
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Code(FUNC), but received nil as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Code((func(string) bool)(nil)),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC cannot be a nil function"),
		})

	checkError(t, "never tested",
		td.Code("test"),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Code(FUNC), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Code(func(x ...int) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested",
		td.Code(func() bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested",
		td.Code(func(a, b, c, d string) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested",
		td.Code(func(a int, b string) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested",
		td.Code(func(t *td.T, a int, b string) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested", // because it is certainly an error
		td.Code(func(assert, require *td.T) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must take only one non-variadic argument or (*td.T, arg) or (*td.T, *td.T, arg)"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (bool, int) { return true, 0 }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (error, string) { return nil, "" }), //nolint: staticcheck
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (int, string) { return 0, "" }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (string, bool) { return "", true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) (bool, string, int) { return true, "", 0 }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) {}),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(n int) int { return 0 }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return bool or (bool, string) or error"),
		})

	checkError(t, "never tested",
		td.Code(func(t *td.T, a int) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return nothing"),
		})

	checkError(t, "never tested",
		td.Code(func(assert, require *td.T, a int) bool { return true }),
		expectedError{
			Message: mustBe("bad usage of Code operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Code(FUNC): FUNC must return nothing"),
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

func TestCodeCustom(t *testing.T) {
	// Specific _checkOK func as td.Code(FUNC) with FUNC(t,arg) or
	// FUNC(assert,require,arg) works in non-boolean context but cannot
	// work in boolean context as there is no initial testing.TB instance
	_customCheckOK := func(t *testing.T, got, expected any, args ...any) bool {
		t.Helper()
		if !td.Cmp(t, got, expected, args...) {
			return false
		}
		// Should always fail in boolean context as no original testing.TB available
		err := td.EqDeeplyError(got, expected)
		if err == nil {
			t.Error(`Boolean context succeeded and it shouldn't`)
			return false
		}
		expErr := expectedError{
			Message: mustBe("cannot build *td.T instance"),
			Path:    mustBe("DATA"),
			Summary: mustBe("original testing.TB instance is missing"),
		}
		if !strings.HasPrefix(expected.(fmt.Stringer).String(), "Code") {
			expErr = ifaceExpectedError(t, expErr)
		}
		if !matchError(t, err.(*ctxerr.Error), expErr, true, args...) {
			return false
		}
		if td.EqDeeply(got, expected) {
			t.Error(`Boolean context succeeded and it shouldn't`)
			return false
		}
		return true
	}

	customCheckOK(t, _customCheckOK, 123, td.Code(func(t *td.T, n int) {
		t.Cmp(t.Config.FailureIsFatal, false)
		t.Cmp(n, 123)
	}))

	customCheckOK(t, _customCheckOK, 123, td.Code(func(assert, require *td.T, n int) {
		assert.Cmp(assert.Config.FailureIsFatal, false)
		assert.Cmp(require.Config.FailureIsFatal, true)
		assert.Cmp(n, 123)
		require.Cmp(n, 123)
	}))

	got := map[string]int{"foo": 123}

	t.Run("Simple success", func(t *testing.T) {
		mockT := test.NewTestingTB("TestCodeCustom")
		td.Cmp(mockT, got, td.Map(map[string]int{}, td.MapEntries{
			"foo": td.Code(func(t *td.T, n int) {
				t.Cmp(n, 123)
			}),
		}))
		test.EqualInt(t, len(mockT.Messages), 0)
	})

	t.Run("Simple failure", func(t *testing.T) {
		mockT := test.NewTestingTB("TestCodeCustom")
		td.NewT(mockT).
			RootName("PIPO").
			Cmp(got, td.Map(map[string]int{}, td.MapEntries{
				"foo": td.Code(func(t *td.T, n int) {
					t.Cmp(n, 124)                                   // inherit only RootName
					t.RootName(t.Config.OriginalPath()).Cmp(n, 125) // recover current path
					t.RootName("").Cmp(n, 126)                      // undo RootName inheritance
				}),
			}))
		test.IsTrue(t, mockT.HasFailed)
		test.IsFalse(t, mockT.IsFatal)
		missing := mockT.ContainsMessages(
			`PIPO: values differ`,
			`     got: 123`,
			`expected: 124`,
			`PIPO["foo"]: values differ`,
			`     got: 123`,
			`expected: 125`,
			`DATA: values differ`,
			`     got: 123`,
			`expected: 126`,
		)
		if len(missing) != 0 {
			t.Error("Following expected messages are not found:\n-", strings.Join(missing, "\n- "))
			t.Error("================================ in:")
			t.Error(strings.Join(mockT.Messages, "\n"))
			t.Error("====================================")
		}
	})

	t.Run("AssertRequire success", func(t *testing.T) {
		mockT := test.NewTestingTB("TestCodeCustom")
		td.Cmp(mockT, got, td.Map(map[string]int{}, td.MapEntries{
			"foo": td.Code(func(assert, require *td.T, n int) {
				assert.Cmp(n, 123)
				require.Cmp(n, 123)
			}),
		}))
		test.EqualInt(t, len(mockT.Messages), 0)
	})

	t.Run("AssertRequire failure", func(t *testing.T) {
		mockT := test.NewTestingTB("TestCodeCustom")
		td.NewT(mockT).
			RootName("PIPO").
			Cmp(got, td.Map(map[string]int{}, td.MapEntries{
				"foo": td.Code(func(assert, require *td.T, n int) {
					assert.Cmp(n, 124)                                         // inherit only RootName
					assert.RootName(assert.Config.OriginalPath()).Cmp(n, 125)  // recover current path
					assert.RootName(require.Config.OriginalPath()).Cmp(n, 126) // recover current path
					assert.RootName("").Cmp(n, 127)                            // undo RootName inheritance
				}),
			}))
		test.IsTrue(t, mockT.HasFailed)
		test.IsFalse(t, mockT.IsFatal)
		missing := mockT.ContainsMessages(
			`PIPO: values differ`,
			`     got: 123`,
			`expected: 124`,
			`PIPO["foo"]: values differ`,
			`     got: 123`,
			`expected: 125`,
			`PIPO["foo"]: values differ`,
			`     got: 123`,
			`expected: 126`,
			`DATA: values differ`,
			`     got: 123`,
			`expected: 127`,
		)
		if len(missing) != 0 {
			t.Error("Following expected messages are not found:\n-", strings.Join(missing, "\n- "))
			t.Error("================================ in:")
			t.Error(strings.Join(mockT.Messages, "\n"))
			t.Error("====================================")
		}
	})

	t.Run("AssertRequire fatalfailure", func(t *testing.T) {
		mockT := test.NewTestingTB("TestCodeCustom")
		td.NewT(mockT).
			RootName("PIPO").
			Cmp(got, td.Map(map[string]int{}, td.MapEntries{
				"foo": td.Code(func(assert, require *td.T, n int) {
					mockT.CatchFatal(func() {
						assert.RootName("FIRST").Cmp(n, 124)
						require.RootName("SECOND").Cmp(n, 125)
						assert.RootName("THIRD").Cmp(n, 126)
					})
				}),
			}))
		test.IsTrue(t, mockT.HasFailed)
		test.IsTrue(t, mockT.IsFatal)
		missing := mockT.ContainsMessages(
			`FIRST: values differ`,
			`     got: 123`,
			`expected: 124`,
			`SECOND: values differ`,
			`     got: 123`,
			`expected: 125`,
		)
		mesgs := strings.Join(mockT.Messages, "\n")
		if len(missing) != 0 {
			t.Error("Following expected messages are not found:\n-", strings.Join(missing, "\n- "))
			t.Error("================================ in:")
			t.Error(mesgs)
			t.Error("====================================")
		}

		if strings.Contains(mesgs, "THIRD") {
			t.Error("THIRD test found, but shouldn't, in:")
			t.Error(mesgs)
			t.Error("====================================")
		}
	})
}

func TestCodeTypeBehind(t *testing.T) {
	// Type behind is the code function parameter one

	equalTypes(t, td.Code(func(n int) bool { return n != 0 }), 23)
	equalTypes(t, td.Code(func(_ *td.T, n int) {}), 23)
	equalTypes(t, td.Code(func(_, _ *td.T, n int) {}), 23)

	type MyTime time.Time

	equalTypes(t,
		td.Code(func(t MyTime) bool { return time.Time(t).IsZero() }),
		MyTime{})
	equalTypes(t, td.Code(func(_ *td.T, t MyTime) {}), MyTime{})
	equalTypes(t, td.Code(func(_, _ *td.T, t MyTime) {}), MyTime{})

	// Erroneous op
	equalTypes(t, td.Code(nil), nil)
}
