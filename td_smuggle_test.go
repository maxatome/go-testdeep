// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestSmuggle(t *testing.T) {
	num := 42
	gotStruct := MyStruct{
		MyStructMid: MyStructMid{
			MyStructBase: MyStructBase{
				ValBool: true,
			},
			ValStr: "foobar",
		},
		ValInt: 123,
		Ptr:    &num,
	}

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	//
	// One returned value
	checkOK(t,
		gotTime,
		testdeep.Smuggle(
			func(date time.Time) int {
				return date.Year()
			},
			testdeep.Between(2010, 2020)))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) testdeep.SmuggledGot {
				return testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}
			},
			testdeep.Contains("oob")))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) *testdeep.SmuggledGot {
				return &testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}
			},
			testdeep.Contains("oob")))

	//
	// 2 returned values
	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (string, bool) {
				if s.ValStr == "" {
					return "", false
				}
				return s.ValStr, true
			},
			testdeep.Contains("oob")))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (testdeep.SmuggledGot, bool) {
				if s.ValStr == "" {
					return testdeep.SmuggledGot{}, false
				}
				return testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true
			},
			testdeep.Contains("oob")))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (*testdeep.SmuggledGot, bool) {
				if s.ValStr == "" {
					return nil, false
				}
				return &testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true
			},
			testdeep.Contains("oob")))

	//
	// 3 returned values
	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (string, bool, string) {
				if s.ValStr == "" {
					return "", false, "ValStr must not be empty"
				}
				return s.ValStr, true, ""
			},
			testdeep.Contains("oob")))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (testdeep.SmuggledGot, bool, string) {
				if s.ValStr == "" {
					return testdeep.SmuggledGot{}, false, "ValStr must not be empty"
				}
				return testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true, ""
			},
			testdeep.Contains("oob")))

	checkOK(t,
		gotStruct,
		testdeep.Smuggle(
			func(s MyStruct) (*testdeep.SmuggledGot, bool, string) {
				if s.ValStr == "" {
					return nil, false, "ValStr must not be empty"
				}
				return &testdeep.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true, ""
			},
			testdeep.Contains("oob")))

	//
	// Errors
	checkError(t, 123,
		testdeep.Smuggle(func(n float64) int { return int(n) }, 123),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("float64"),
		})

	type xInt int
	checkError(t, xInt(12),
		testdeep.Smuggle(func(n int) int64 { return int64(n) }, 12),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("testdeep_test.xInt"),
			Expected: mustBe("int"),
		})

	checkError(t, 12,
		testdeep.Smuggle(func(n int) (int, bool) { return n, false }, 12),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("  value: (int) 12\nit failed but didn't say why"),
		})

	type MyBool bool
	type MyString string
	checkError(t, 12,
		testdeep.Smuggle(func(n int) (int, MyBool, MyString) {
			return n, false, "very custom error"
		}, 12),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: (int) 12\nit failed coz: very custom error"),
		})

	checkError(t, 12,
		testdeep.Smuggle(func(n int) *testdeep.SmuggledGot { return nil }, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("(int64) 13"),
		})

	//
	// Errors behind Smuggle()
	checkError(t, 12,
		testdeep.Smuggle(func(n int) int64 { return int64(n) }, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA<smuggled>"),
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 13"),
		})

	checkError(t, 12,
		testdeep.Smuggle(func(n int) testdeep.SmuggledGot {
			return testdeep.SmuggledGot{
				// With Name = ""
				Got: int64(n),
			}
		}, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA<smuggled>"),
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 13"),
		})

	checkError(t, 12,
		testdeep.Smuggle(func(n int) *testdeep.SmuggledGot {
			return &testdeep.SmuggledGot{
				Name: "<int64>",
				Got:  int64(n),
			}
		}, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA<int64>"), // no dot added between DATA and <int64>
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 13"),
		})

	checkError(t, 12,
		testdeep.Smuggle(func(n int) *testdeep.SmuggledGot {
			return &testdeep.SmuggledGot{
				Name: "Int64",
				Got:  int64(n),
			}
		}, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.Int64"), // dot added between DATA and Int64
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 13"),
		})

	//
	// Bad usage
	checkPanic(t, func() { testdeep.Smuggle("test", 12) }, "usage: Smuggle")

	// Bad number of args
	checkPanic(t, func() {
		testdeep.Smuggle(func() int { return 0 }, 12)
	}, "FUNC must take only one argument")

	checkPanic(t, func() {
		testdeep.Smuggle(func(a int, b string) int { return 0 }, 12)
	}, "FUNC must take only one argument")

	// Bad number of returned values
	checkPanic(t, func() {
		testdeep.Smuggle(func(a int) {}, 12)
	}, "FUNC must return value or (value, bool) or (value, bool, string)")

	checkPanic(t, func() {
		testdeep.Smuggle(
			func(a int) (int, bool, string, int) { return 0, false, "", 23 },
			12)
	}, "FUNC must return value or (value, bool) or (value, bool, string)")

	// Bad returned types
	checkPanic(t, func() {
		testdeep.Smuggle(func(a int) (int, int) { return 0, 0 }, 12)
	}, "FUNC must return value or (value, bool) or (value, bool, string)")

	checkPanic(t, func() {
		testdeep.Smuggle(func(a int) (int, bool, int) { return 0, false, 23 }, 12)
	}, "FUNC must return value or (value, bool) or (value, bool, string)")

	//
	// String
	test.EqualStr(t,
		testdeep.Smuggle(func(n int) int { return 0 }, 12).String(),
		"Smuggle(func(int) int)")

	test.EqualStr(t,
		testdeep.Smuggle(func(n int) (int, bool) { return 23, false }, 12).String(),
		"Smuggle(func(int) (int, bool))")

	test.EqualStr(t,
		testdeep.Smuggle(func(n int) (int, MyBool, MyString) { return 23, false, "" }, 12).
			String(),
		"Smuggle(func(int) (int, testdeep_test.MyBool, testdeep_test.MyString))")
}

func TestSmuggleTypeBehind(t *testing.T) {
	// Type behind is the smuggle function parameter one

	equalTypes(t, testdeep.Smuggle(func(n int) bool { return n != 0 }, true), 23)

	type MyTime time.Time

	equalTypes(t,
		testdeep.Smuggle(
			func(t MyTime) time.Time { return time.Time(t) },
			time.Now()),
		MyTime{})
}
