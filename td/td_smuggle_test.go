// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
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
		td.Smuggle(
			func(date time.Time) int {
				return date.Year()
			},
			td.Between(2010, 2020)))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) td.SmuggledGot {
				return td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}
			},
			td.Contains("oob")))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) *td.SmuggledGot {
				return &td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}
			},
			td.Contains("oob")))

	//
	// 2 returned values
	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (string, bool) {
				if s.ValStr == "" {
					return "", false
				}
				return s.ValStr, true
			},
			td.Contains("oob")))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (td.SmuggledGot, bool) {
				if s.ValStr == "" {
					return td.SmuggledGot{}, false
				}
				return td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true
			},
			td.Contains("oob")))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (*td.SmuggledGot, bool) {
				if s.ValStr == "" {
					return nil, false
				}
				return &td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true
			},
			td.Contains("oob")))

	//
	// 3 returned values
	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (string, bool, string) {
				if s.ValStr == "" {
					return "", false, "ValStr must not be empty"
				}
				return s.ValStr, true, ""
			},
			td.Contains("oob")))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (td.SmuggledGot, bool, string) {
				if s.ValStr == "" {
					return td.SmuggledGot{}, false, "ValStr must not be empty"
				}
				return td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true, ""
			},
			td.Contains("oob")))

	checkOK(t,
		gotStruct,
		td.Smuggle(
			func(s MyStruct) (*td.SmuggledGot, bool, string) {
				if s.ValStr == "" {
					return nil, false, "ValStr must not be empty"
				}
				return &td.SmuggledGot{
					Name: "ValStr",
					Got:  s.ValStr,
				}, true, ""
			},
			td.Contains("oob")))

	//
	// Convertible types
	checkOK(t, 123,
		td.Smuggle(func(n float64) int { return int(n) }, 123))

	type xInt int
	checkOK(t, xInt(123),
		td.Smuggle(func(n int) int64 { return int64(n) }, int64(123)))
	checkOK(t, xInt(123),
		td.Smuggle(func(n uint32) int64 { return int64(n) }, int64(123)))

	checkOK(t, int32(123),
		td.Smuggle(func(n int64) int { return int(n) }, 123))

	checkOK(t, gotTime,
		td.Smuggle(func(t fmt.Stringer) string { return t.String() },
			"2018-05-23 12:13:14 +0000 UTC"))

	//
	// Errors
	checkError(t, "123",
		td.Smuggle(func(n float64) int { return int(n) }, 123),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("float64"),
		})

	checkError(t, nil,
		td.Smuggle(func(n int64) int { return int(n) }, 123),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("int64"),
		})

	checkError(t, 12,
		td.Smuggle(func(n int) (int, bool) { return n, false }, 12),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("  value: 12\nit failed but didn't say why"),
		})

	type MyBool bool
	type MyString string
	checkError(t, 12,
		td.Smuggle(func(n int) (int, MyBool, MyString) {
			return n, false, "very custom error"
		}, 12),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: very custom error"),
		})

	checkError(t, 12,
		td.Smuggle(func(n int) (int, error) {
			return n, errors.New("very custom error")
		}, 12),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: very custom error"),
		})

	checkError(t, 12,
		td.Smuggle(func(n int) *td.SmuggledGot { return nil }, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("(int64) 13"),
		})

	// Internal use
	checkError(t, 12,
		td.Smuggle(func(n int) (int, error) {
			return n, &ctxerr.Error{
				Message: "my message",
				Summary: ctxerr.NewSummary("my summary"),
			}
		}, 13),
		expectedError{
			Message: mustBe("my message"),
			Path:    mustBe("DATA"),
			Summary: mustBe("my summary"),
		})

	//
	// Errors behind Smuggle()
	checkError(t, 12,
		td.Smuggle(func(n int) int64 { return int64(n) }, int64(13)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA<smuggled>"),
			Got:      mustBe("(int64) 12"),
			Expected: mustBe("(int64) 13"),
		})

	checkError(t, gotStruct,
		td.Smuggle("MyStructMid.MyStructBase.ValBool", false),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.MyStructMid.MyStructBase.ValBool"),
			Got:      mustBe("true"),
			Expected: mustBe("false"),
		})

	checkError(t, 12,
		td.Smuggle(func(n int) td.SmuggledGot {
			return td.SmuggledGot{
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
		td.Smuggle(func(n int) *td.SmuggledGot {
			return &td.SmuggledGot{
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
		td.Smuggle(func(n int) *td.SmuggledGot {
			return &td.SmuggledGot{
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
	test.CheckPanic(t, func() { td.Smuggle(123, 12) }, "usage: Smuggle")
	test.CheckPanic(t, func() { td.Smuggle("foo.9bingo", 12) },
		"bad field name `9bingo' in FIELDS_PATH")

	// Bad number of args
	test.CheckPanic(t, func() {
		td.Smuggle(func() int { return 0 }, 12)
	}, "FUNC must take only one argument")

	test.CheckPanic(t, func() {
		td.Smuggle(func(x ...int) int { return 0 }, 12)
	}, "FUNC must take only one argument")

	test.CheckPanic(t, func() {
		td.Smuggle(func(a int, b string) int { return 0 }, 12)
	}, "FUNC must take only one argument")

	// Bad number of returned values
	const errMesg = "FUNC must return value or (value, bool) or (value, bool, string) or (value, error)"

	test.CheckPanic(t, func() {
		td.Smuggle(func(a int) {}, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		td.Smuggle(
			func(a int) (int, bool, string, int) { return 0, false, "", 23 },
			12)
	}, errMesg)

	// Bad returned types
	test.CheckPanic(t, func() {
		td.Smuggle(func(a int) (int, int) { return 0, 0 }, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		td.Smuggle(func(a int) (int, bool, int) { return 0, false, 23 }, 12)
	}, errMesg)

	test.CheckPanic(t, func() {
		td.Smuggle(func(a int) (int, error, string) { return 0, nil, "" }, 12) // nolint: staticcheck
	}, errMesg)

	//
	// String
	test.EqualStr(t,
		td.Smuggle(func(n int) int { return 0 }, 12).String(),
		"Smuggle(func(int) int)")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, bool) { return 23, false }, 12).String(),
		"Smuggle(func(int) (int, bool))")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, error) { return 23, nil }, 12).String(),
		"Smuggle(func(int) (int, error))")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, MyBool, MyString) { return 23, false, "" }, 12).
			String(),
		"Smuggle(func(int) (int, td_test.MyBool, td_test.MyString))")
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
	checkOK(t, gotStruct, td.Smuggle("ValInt", 123))
	checkOK(t, gotStruct,
		td.Smuggle("MyStructMid.ValStr", td.Contains("oob")))
	checkOK(t, gotStruct,
		td.Smuggle("MyStructMid.MyStructBase.ValBool", true))
	checkOK(t, gotStruct, td.Smuggle("ValBool", true)) // thanks to composition

	// OK across pointers
	checkOK(t, b, td.Smuggle("PA.Num", 2))
	checkOK(t, b, td.Smuggle("PppA.Num", 3))

	// OK with interface{}
	checkOK(t, b, td.Smuggle("Iface.Num", 4))
	checkOK(t, b, td.Smuggle("Iface2.Num", 3))
	checkOK(t, b, td.Smuggle("C.Iface1.Num", 7))
	checkOK(t, b, td.Smuggle("C.Iface2.Num", 8))

	// Errors
	checkError(t, 12, td.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: it is not a struct and should be"),
		})
	checkError(t, gotStruct, td.Smuggle("ValInt.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `ValInt' is not a struct and should be"),
		})
	checkError(t, gotStruct, td.Smuggle("MyStructMid.ValStr.foobar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `MyStructMid.ValStr' is not a struct and should be"),
		})

	checkError(t, gotStruct, td.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `foo' not found"),
		})

	checkError(t, b, td.Smuggle("C.PA2.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.PA2' is nil"),
		})
	checkError(t, b, td.Smuggle("C.Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.Iface3' is nil"),
		})
	checkError(t, b, td.Smuggle("C.Iface4.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `C.Iface4' is nil"),
		})
	checkError(t, b, td.Smuggle("Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field `Iface3' is nil"),
		})
}

func TestSmuggleTypeBehind(t *testing.T) {
	// Type behind is the smuggle function parameter one

	equalTypes(t, td.Smuggle(func(n int) bool { return n != 0 }, true), 23)

	type MyTime time.Time

	equalTypes(t,
		td.Smuggle(
			func(t MyTime) time.Time { return time.Time(t) },
			time.Now()),
		MyTime{})

	equalTypes(t,
		td.Smuggle(func(from interface{}) interface{} { return from }, nil),
		reflect.TypeOf((*interface{})(nil)).Elem())
}
