// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util_test

import (
	"bytes"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/util"
)

type myTestDeepStringer struct {
	types.TestDeepStamp
}

func (m myTestDeepStringer) String() string {
	return "TesT!"
}

func TestToString(t *testing.T) {
	for _, curTest := range []struct {
		paramGot interface{}
		expected string
	}{
		{paramGot: nil, expected: "nil"},
		{paramGot: "foobar", expected: `"foobar"`},
		{paramGot: "foo\rbar", expected: `(string) (len=7) "foo\rbar"`},
		{paramGot: "foo\u2028bar", expected: `(string) (len=9) "foo\u2028bar"`},
		{paramGot: `foo"bar`, expected: "`foo\"bar`"},
		{paramGot: "foo\n\"bar", expected: "`foo\n\"bar`"},
		{paramGot: "foo`\"\nbar", expected: "(string) (len=9) \"foo`\\\"\\nbar\""},
		{paramGot: "foo`\n\"bar", expected: "(string) (len=9) \"foo`\\n\\\"bar\""},
		{paramGot: "foo\n`\"bar", expected: "(string) (len=9) \"foo\\n`\\\"bar\""},
		{paramGot: "foo\n\"`bar", expected: "(string) (len=9) \"foo\\n\\\"`bar\""},
		{paramGot: reflect.ValueOf("foobar"), expected: `"foobar"`},
		{paramGot: []reflect.Value{reflect.ValueOf("foo"), reflect.ValueOf("bar")},
			expected: `("foo",
 "bar")`},
		{paramGot: types.RawString("test"), expected: "test"},
		{paramGot: types.RawInt(42), expected: "42"},
		{paramGot: myTestDeepStringer{}, expected: "TesT!"},
		{paramGot: 42, expected: "42"},
		{paramGot: true, expected: "true"},
		{paramGot: false, expected: "false"},
		{paramGot: int64(42), expected: "(int64) 42"},
	} {
		test.EqualStr(t, util.ToString(curTest.paramGot), curTest.expected)
	}
}

func TestIndentString(t *testing.T) {
	for _, curTest := range []struct {
		ParamGot string
		Expected string
	}{
		{ParamGot: "", Expected: ""},
		{ParamGot: "pipo", Expected: "pipo"},
		{ParamGot: "pipo\nbingo\nzip", Expected: "pipo\n-bingo\n-zip"},
	} {
		test.EqualStr(t, util.IndentString(curTest.ParamGot, "-"), curTest.Expected)

		var buf bytes.Buffer
		util.IndentStringIn(&buf, curTest.ParamGot, "-")
		test.EqualStr(t, buf.String(), curTest.Expected)
	}
}

func TestSliceToBuffer(t *testing.T) {
	for _, curTest := range []struct {
		BufInit  string
		Items    []interface{}
		Expected string
	}{
		{BufInit: ">", Items: nil, Expected: ">()"},
		{BufInit: ">", Items: []interface{}{"pipo"}, Expected: `>("pipo")`},
		{
			BufInit: ">",
			Items:   []interface{}{"pipo", "bingo", "zip"},
			Expected: `>("pipo",
  "bingo",
  "zip")`,
		},
		{
			BufInit: "List\n  of\nitems:\n>",
			Items:   []interface{}{"pipo", "bingo", "zip"},
			Expected: `List
  of
items:
>("pipo",
  "bingo",
  "zip")`,
		},
	} {
		var items []reflect.Value
		if curTest.Items != nil {
			items = make([]reflect.Value, len(curTest.Items))
			for i, val := range curTest.Items {
				items[i] = reflect.ValueOf(val)
			}
		}

		buf := bytes.NewBufferString(curTest.BufInit)
		test.EqualStr(t, util.SliceToBuffer(buf, items).String(),
			curTest.Expected)
	}
}

func TestTypeFullName(t *testing.T) {
	// our full package name
	pc, _, _, _ := runtime.Caller(0)
	pkg := strings.TrimSuffix(runtime.FuncForPC(pc).Name(), ".TestTypeFullName")

	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(123)), "int")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf([]int{})), "[]int")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf([3]int{})), "[3]int")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf((**float64)(nil))), "**float64")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(map[int]float64{})), "map[int]float64")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(struct{}{})), "struct {}")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(struct {
		a int
		b bool
	}{})), "struct { a int; b bool }")
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(struct {
		s struct{ a []int }
		b bool
	}{})), "struct { s struct { a []int }; b bool }")

	type anon struct{ a []int } //nolint: structcheck,unused
	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(struct {
		anon
		b bool
	}{})), "struct { "+pkg+".anon; b bool }")

	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(func() {})), "func()")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func(a int) {})),
		"func(int)")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func(a int, b ...bool) rune { return 0 })),
		"func(int, ...bool) int32")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func() (int, bool, int) { return 0, true, 0 })),
		"func() (int, bool, int)")

	test.EqualStr(t, util.TypeFullName(reflect.TypeOf(func() {})), "func()")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func(a int) {})),
		"func(int)")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func(a int, b ...bool) rune { return 0 })),
		"func(int, ...bool) int32")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf(func() (int, bool, int) { return 0, true, 0 })),
		"func() (int, bool, int)")

	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf((<-chan []int)(nil))),
		"<-chan []int")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf((chan<- []int)(nil))),
		"chan<- []int")
	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf((chan []int)(nil))),
		"chan []int")

	test.EqualStr(t,
		util.TypeFullName(reflect.TypeOf((*interface{})(nil))),
		"*interface {}")
}
