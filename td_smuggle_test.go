// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
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
	// Convertible types
	checkOK(t, 123,
		testdeep.Smuggle(func(n float64) int { return int(n) }, 123))

	type xInt int
	checkOK(t, xInt(123),
		testdeep.Smuggle(func(n int) int64 { return int64(n) }, int64(123)))
	checkOK(t, xInt(123),
		testdeep.Smuggle(func(n uint32) int64 { return int64(n) }, int64(123)))

	type tVal struct{ Val interface{} }
	checkOK(t, tVal{Val: int32(123)},
		testdeep.Struct(tVal{}, testdeep.StructFields{
			"Val": testdeep.Smuggle(func(n int64) int { return int(n) }, 123),
		}))

	//
	// Errors
	checkError(t, "123",
		testdeep.Smuggle(func(n float64) int { return int(n) }, 123),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("float64"),
		})

	checkError(t, tVal{},
		testdeep.Struct(tVal{}, testdeep.StructFields{
			"Val": testdeep.Smuggle(func(n int64) int { return int(n) }, 123),
		}),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA.Val"),
			Got:      mustBe("interface {}"),
			Expected: mustBe("int64"),
		})

	checkError(t, tVal{Val: "str"},
		testdeep.Struct(tVal{}, testdeep.StructFields{
			"Val": testdeep.Smuggle(func(n int64) int { return int(n) }, 123),
		}),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA.Val"),
			Got:      mustBe("string"),
			Expected: mustBe("int64"),
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
		testdeep.Smuggle(func(n int) (int, error) {
			return n, errors.New("very custom error")
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

	checkError(t, gotStruct,
		testdeep.Smuggle("MyStructMid.MyStructBase.ValBool", false),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.MyStructMid.MyStructBase.ValBool"),
			Got:      mustBe("(bool) true"),
			Expected: mustBe("(bool) false"),
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
	test.CheckPanic(t, func() { testdeep.Smuggle(123, 12) }, "usage: Smuggle")
	test.CheckPanic(t, func() { testdeep.Smuggle("foo.9bingo", 12) },
		"bad field name `9bingo' in FIELDS_PATH")

	// Bad number of args
	test.CheckPanic(t, func() {
		testdeep.Smuggle(func() int { return 0 }, 12)
	}, "FUNC must take only one argument")

	test.CheckPanic(t, func() {
		testdeep.Smuggle(func(a int, b string) int { return 0 }, 12)
	}, "FUNC must take only one argument")

	// Bad number of returned values
	const errMesg = "FUNC must return value or (value, bool) or (value, bool, string) or (value, error)"

	test.CheckPanic(t, func() {
		testdeep.Smuggle(func(a int) {}, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		testdeep.Smuggle(
			func(a int) (int, bool, string, int) { return 0, false, "", 23 },
			12)
	}, errMesg)

	// Bad returned types
	test.CheckPanic(t, func() {
		testdeep.Smuggle(func(a int) (int, int) { return 0, 0 }, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		testdeep.Smuggle(func(a int) (int, bool, int) { return 0, false, 23 }, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		testdeep.Smuggle(func(a int) (int, error, string) { return 0, nil, "" }, 12)
	}, errMesg)

	//
	// String
	test.EqualStr(t,
		testdeep.Smuggle(func(n int) int { return 0 }, 12).String(),
		"Smuggle(func(int) int)")

	test.EqualStr(t,
		testdeep.Smuggle(func(n int) (int, bool) { return 23, false }, 12).String(),
		"Smuggle(func(int) (int, bool))")

	test.EqualStr(t,
		testdeep.Smuggle(func(n int) (int, error) { return 23, nil }, 12).String(),
		"Smuggle(func(int) (int, error))")

	test.EqualStr(t,
		testdeep.Smuggle(func(n int) (int, MyBool, MyString) { return 23, false, "" }, 12).
			String(),
		"Smuggle(func(int) (int, testdeep_test.MyBool, testdeep_test.MyString))")
}

func TestSmuggleFieldPath(t *testing.T) {
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

	type A struct {
		Num int
		Str string
	}
	type C struct {
		A      A
		PA1    *A
		PA2    *A
		Iface1 interface{}
		Iface2 interface{}
		Iface3 interface{}
		Iface4 interface{}
	}
	type B struct {
		A      A
		PA     *A
		PppA   ***A
		Iface  interface{}
		Iface2 interface{}
		Iface3 interface{}
		C      *C
	}
	pa := &A{Num: 3, Str: "three"}
	ppa := &pa
	b := B{
		A:      A{Num: 1, Str: "one"},
		PA:     &A{Num: 2, Str: "two"},
		PppA:   &ppa,
		Iface:  A{Num: 4, Str: "four"},
		Iface2: &ppa,
		Iface3: nil,
		C: &C{
			A:      A{Num: 5, Str: "five"},
			PA1:    &A{Num: 6, Str: "six"},
			PA2:    nil, // explicit to be clear
			Iface1: A{Num: 7, Str: "seven"},
			Iface2: &A{Num: 8, Str: "eight"},
			Iface3: nil, // explicit to be clear
			Iface4: (*A)(nil),
		},
	}

	//
	// OK
	checkOK(t, gotStruct, testdeep.Smuggle("ValInt", 123))
	checkOK(t, gotStruct,
		testdeep.Smuggle("MyStructMid.ValStr", testdeep.Contains("oob")))
	checkOK(t, gotStruct,
		testdeep.Smuggle("MyStructMid.MyStructBase.ValBool", true))
	checkOK(t, gotStruct, testdeep.Smuggle("ValBool", true)) // thanks to composition

	// OK across pointers
	checkOK(t, b, testdeep.Smuggle("PA.Num", 2))
	checkOK(t, b, testdeep.Smuggle("PppA.Num", 3))

	// OK with interface{}
	checkOK(t, b, testdeep.Smuggle("Iface.Num", 4))
	checkOK(t, b, testdeep.Smuggle("Iface2.Num", 3))
	checkOK(t, b, testdeep.Smuggle("C.Iface1.Num", 7))
	checkOK(t, b, testdeep.Smuggle("C.Iface2.Num", 8))

	// Errors
	checkError(t, 12, testdeep.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: (int) 12\nit failed coz: it is not a struct and should be"),
		})
	checkError(t, gotStruct, testdeep.Smuggle("ValInt.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `ValInt' is not a struct and should be"),
		})
	checkError(t, gotStruct, testdeep.Smuggle("MyStructMid.ValStr.foobar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `MyStructMid.ValStr' is not a struct and should be"),
		})

	checkError(t, gotStruct, testdeep.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `foo' not found"),
		})

	checkError(t, b, testdeep.Smuggle("C.PA2.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.PA2' is nil"),
		})
	checkError(t, b, testdeep.Smuggle("C.Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.Iface3' is nil"),
		})
	checkError(t, b, testdeep.Smuggle("C.Iface4.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.Iface4' is nil"),
		})
	checkError(t, b, testdeep.Smuggle("Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `Iface3' is nil"),
		})
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
