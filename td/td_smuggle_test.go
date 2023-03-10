// Copyright (c) 2018-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

// reArmReader is a bytes.Reader that re-arms when an error occurs,
// typically on EOF.
type reArmReader bytes.Reader

var _ io.Reader = (*reArmReader)(nil)

func newReArmReader(b []byte) *reArmReader {
	return (*reArmReader)(bytes.NewReader(b))
}

func (r *reArmReader) Read(b []byte) (n int, err error) {
	n, err = (*bytes.Reader)(r).Read(b)
	if err != nil {
		(*bytes.Reader)(r).Seek(0, io.SeekStart) //nolint: errcheck
	}
	return
}

func (r *reArmReader) String() string { return "<no string here>" }

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

	checkOK(t, []byte("{}"),
		td.Smuggle(
			func(x json.RawMessage) json.RawMessage { return x },
			td.JSON(`{}`)))

	//
	// bytes slice caster variations
	checkOK(t, []byte(`{"foo":1}`),
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":1}`)))

	checkOK(t, []byte(`{"foo":1}`),
		td.Smuggle(json.RawMessage(nil), td.JSON(`{"foo":1}`)))

	checkOK(t, []byte(`{"foo":1}`),
		td.Smuggle(reflect.TypeOf(json.RawMessage(nil)), td.JSON(`{"foo":1}`)))

	checkOK(t, `{"foo":1}`,
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":1}`)))

	checkOK(t, newReArmReader([]byte(`{"foo":1}`)), // io.Reader first
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":1}`)))

	checkError(t, nil,
		td.Smuggle(json.RawMessage{}, td.JSON(`{}`)),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("json.RawMessage or convertible or io.Reader"),
		})

	checkError(t, MyStruct{},
		td.Smuggle(json.RawMessage{}, td.JSON(`{}`)),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MyStruct"),
			Expected: mustBe("json.RawMessage or convertible or io.Reader"),
		})

	checkError(t, errReader{}, // erroneous io.Reader
		td.Smuggle(json.RawMessage{}, td.JSON(`{}`)),
		expectedError{
			Message: mustBe("an error occurred while reading from io.Reader"),
			Path:    mustBe("DATA"),
			Summary: mustBe("an error occurred"),
		})

	//
	// strings caster variations
	type myString string
	checkOK(t, `pipo bingo`,
		td.Smuggle("", td.HasSuffix("bingo")))

	checkOK(t, []byte(`pipo bingo`),
		td.Smuggle(myString(""), td.HasSuffix("bingo")))

	checkOK(t, []byte(`pipo bingo`),
		td.Smuggle(reflect.TypeOf(myString("")), td.HasSuffix("bingo")))

	checkOK(t, newReArmReader([]byte(`pipo bingo`)), // io.Reader first
		td.Smuggle(myString(""), td.HasSuffix("bingo")))

	checkError(t, nil,
		td.Smuggle("", "bingo"),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("string or convertible or io.Reader"),
		})

	checkError(t, MyStruct{},
		td.Smuggle(myString(""), "bingo"),
		expectedError{
			Message:  mustBe("incompatible parameter type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MyStruct"),
			Expected: mustBe("td_test.myString or convertible or io.Reader"),
		})

	checkError(t, errReader{}, // erroneous io.Reader
		td.Smuggle("", "bingo"),
		expectedError{
			Message: mustBe("an error occurred while reading from io.Reader"),
			Path:    mustBe("DATA"),
			Summary: mustBe("an error occurred"),
		})

	//
	// Any other caster variations
	checkOK(t, `pipo bingo`,
		td.Smuggle([]rune{}, td.Contains([]rune(`bing`))))
	checkOK(t, `pipo bingo`,
		td.Smuggle(([]rune)(nil), td.Contains([]rune(`bing`))))
	checkOK(t, `pipo bingo`,
		td.Smuggle(reflect.TypeOf([]rune{}), td.Contains([]rune(`bing`))))

	checkOK(t, 123.456, td.Smuggle(int64(0), int64(123)))
	checkOK(t, 123.456, td.Smuggle(reflect.TypeOf(int64(0)), int64(123)))

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
	const usage = "Smuggle(FUNC|FIELDS_PATH|ANY_TYPE, TESTDEEP_OPERATOR|EXPECTED_VALUE): "
	checkError(t, "never tested",
		td.Smuggle(nil, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: " + usage[:len(usage)-2] + ", ANY_TYPE cannot be nil nor Interface"),
		})

	checkError(t, nil,
		td.Smuggle(reflect.TypeOf((*fmt.Stringer)(nil)).Elem(), 1234),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: " + usage[:len(usage)-2] + ", ANY_TYPE reflect.Type cannot be Func nor Interface"),
		})

	checkError(t, nil,
		td.Smuggle(reflect.TypeOf(func() {}), 1234),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: " + usage[:len(usage)-2] + ", ANY_TYPE reflect.Type cannot be Func nor Interface"),
		})

	checkError(t, "never tested",
		td.Smuggle((func(string) int)(nil), 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Smuggle(FUNC): FUNC cannot be a nil function"),
		})

	checkError(t, "never tested",
		td.Smuggle("bad[path", 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(usage + `cannot find final ']' in FIELD_PATH "bad[path"`),
		})

	// Bad number of args
	checkError(t, "never tested",
		td.Smuggle(func() int { return 0 }, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(usage + "FUNC must take only one non-variadic argument"),
		})

	checkError(t, "never tested",
		td.Smuggle(func(x ...int) int { return 0 }, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(usage + "FUNC must take only one non-variadic argument"),
		})

	checkError(t, "never tested",
		td.Smuggle(func(a int, b string) int { return 0 }, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(usage + "FUNC must take only one non-variadic argument"),
		})

	// Bad number of returned values
	const errMesg = usage + "FUNC must return value or (value, bool) or (value, bool, string) or (value, error)"

	checkError(t, "never tested",
		td.Smuggle(func(a int) {}, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(errMesg),
		})

	checkError(t, "never tested",
		td.Smuggle(
			func(a int) (int, bool, string, int) { return 0, false, "", 23 },
			12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(errMesg),
		})

	// Bad returned types
	checkError(t, "never tested",
		td.Smuggle(func(a int) (int, int) { return 0, 0 }, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(errMesg),
		})

	checkError(t, "never tested",
		td.Smuggle(func(a int) (int, bool, int) { return 0, false, 23 }, 12),
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(errMesg),
		})

	checkError(t, "never tested",
		td.Smuggle(func(a int) (int, error, string) { return 0, nil, "" }, 12), //nolint: staticcheck
		expectedError{
			Message: mustBe("bad usage of Smuggle operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(errMesg),
		})

	//
	// String
	test.EqualStr(t,
		td.Smuggle(func(n int) int { return 0 }, 12).String(),
		"Smuggle(func(int) int, 12)")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, bool) { return 23, false }, 12).String(),
		"Smuggle(func(int) (int, bool), 12)")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, error) { return 23, nil }, 12).String(),
		"Smuggle(func(int) (int, error), 12)")

	test.EqualStr(t,
		td.Smuggle(func(n int) (int, MyBool, MyString) { return 23, false, "" }, 12).
			String(),
		"Smuggle(func(int) (int, td_test.MyBool, td_test.MyString), 12)")

	test.EqualStr(t,
		td.Smuggle(reflect.TypeOf(42), 23).String(),
		"Smuggle(type:int, 23)")

	test.EqualStr(t,
		td.Smuggle(666, 23).String(),
		"Smuggle(type:int, 23)")

	test.EqualStr(t,
		td.Smuggle("", 23).String(),
		"Smuggle(type:string, 23)")

	test.EqualStr(t,
		td.Smuggle("name", "bob").String(),
		`Smuggle("name", "bob")`)

	// Erroneous op
	test.EqualStr(t,
		td.Smuggle((func(int) int)(nil), 12).String(),
		"Smuggle(<ERROR>)")
}

func TestSmuggleFieldsPath(t *testing.T) {
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
		Iface1 any
		Iface2 any
		Iface3 any
		Iface4 any
	}
	type B struct {
		A      A
		PA     *A
		PppA   ***A
		Iface  any
		Iface2 any
		Iface3 any
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

	// OK with any
	checkOK(t, b, td.Smuggle("Iface.Num", 4))
	checkOK(t, b, td.Smuggle("Iface2.Num", 3))
	checkOK(t, b, td.Smuggle("C.Iface1.Num", 7))
	checkOK(t, b, td.Smuggle("C.Iface2.Num", 8))

	// Errors
	checkError(t, 12, td.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustBe("        value: 12\nit failed coz: it is a int and should be a struct"),
		})
	checkError(t, gotStruct, td.Smuggle("ValInt.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"ValInt\" is a int and should be a struct"),
		})
	checkError(t, gotStruct, td.Smuggle("MyStructMid.ValStr.foobar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"MyStructMid.ValStr\" is a string and should be a struct"),
		})

	checkError(t, gotStruct, td.Smuggle("foo.bar", 23),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"foo\" not found"),
		})

	checkError(t, b, td.Smuggle("C.PA2.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"C.PA2\" is nil"),
		})
	checkError(t, b, td.Smuggle("C.Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"C.Iface3\" is nil"),
		})
	checkError(t, b, td.Smuggle("C.Iface4.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"C.Iface4\" is nil"),
		})
	checkError(t, b, td.Smuggle("Iface3.Num", 456),
		expectedError{
			Message: mustBe("ran smuggle code with %% as argument"),
			Path:    mustBe("DATA"),
			Summary: mustContain("\nit failed coz: field \"Iface3\" is nil"),
		})

	// Referencing maps and array/slices
	x := B{
		Iface: map[string]any{
			"test": []int{2, 3, 4},
		},
		C: &C{
			Iface1: []any{
				map[int]any{42: []string{"pipo"}, 66: [2]string{"foo", "bar"}},
				map[int8]any{42: []string{"pipo"}},
				map[int16]any{42: []string{"pipo"}},
				map[int32]any{42: []string{"pipo"}},
				map[int64]any{42: []string{"pipo"}},
				map[uint]any{42: []string{"pipo"}},
				map[uint8]any{42: []string{"pipo"}},
				map[uint16]any{42: []string{"pipo"}},
				map[uint32]any{42: []string{"pipo"}},
				map[uint64]any{42: []string{"pipo"}},
				map[uintptr]any{42: []string{"pipo"}},
				map[float32]any{42: []string{"pipo"}},
				map[float64]any{42: []string{"pipo"}},
			},
		},
	}
	checkOK(t, x, td.Smuggle("Iface[test][1]", 3))
	checkOK(t, x, td.Smuggle("C.Iface1[0][66][1]", "bar"))
	for i := 0; i < 12; i++ {
		checkOK(t, x,
			td.Smuggle(fmt.Sprintf("C.Iface1[%d][42][0]", i), "pipo"))

		checkOK(t, x,
			td.Smuggle(fmt.Sprintf("C.Iface1[%d][42][-1]", i-12), "pipo"))
	}

	checkOK(t, x, td.Lax(td.Smuggle("PppA", nil)))
	checkOK(t, x, td.Smuggle("PppA", td.Nil()))

	//
	type D struct {
		Iface any
	}

	got := D{
		Iface: []any{
			map[complex64]any{complex(42, 0): []string{"pipo"}},
			map[complex128]any{complex(42, 0): []string{"pipo"}},
		},
	}

	for i := 0; i < 2; i++ {
		checkOK(t, got, td.Smuggle(fmt.Sprintf("Iface[%d][42][0]", i), "pipo"))
		checkOK(t, got, td.Smuggle(fmt.Sprintf("Iface[%d][42][0]", i-2), "pipo"))
	}
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
		td.Smuggle(func(from any) any { return from }, nil),
		reflect.TypeOf((*any)(nil)).Elem())

	equalTypes(t,
		td.Smuggle("foo.bar", nil),
		reflect.TypeOf((*any)(nil)).Elem())

	// Erroneous op
	equalTypes(t, td.Smuggle((func(int) int)(nil), 12), nil)
}
