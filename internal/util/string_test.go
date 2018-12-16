// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util_test

import (
	"bytes"
	"reflect"
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
		ParamGot interface{}
		Expected string
	}{
		{ParamGot: nil, Expected: "nil"},
		{ParamGot: "foobar", Expected: `"foobar"`},
		{ParamGot: "foo\rbar", Expected: `(string) (len=7) "foo\rbar"`},
		{ParamGot: "foo\u2028bar", Expected: `(string) (len=9) "foo\u2028bar"`},
		{ParamGot: reflect.ValueOf("foobar"), Expected: `"foobar"`},
		{ParamGot: types.RawString("test"), Expected: "test"},
		{ParamGot: types.RawInt(42), Expected: "42"},
		{ParamGot: myTestDeepStringer{}, Expected: "TesT!"},
		{ParamGot: 42, Expected: "42"},
	} {
		test.EqualStr(t, util.ToString(curTest.ParamGot), curTest.Expected)
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
