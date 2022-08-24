// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package dark_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
)

func checkFieldValueOK(t *testing.T,
	s reflect.Value, fieldName string, value any,
) {
	t.Helper()

	testName := "field " + fieldName

	fieldOrig := s.FieldByName(fieldName)
	test.IsFalse(t, fieldOrig.CanInterface(), testName+" + fieldOrig.CanInterface()")

	fieldCopy, ok := dark.CopyValue(fieldOrig)
	if test.IsTrue(t, ok, "Can copy "+testName) {
		if test.IsTrue(t, fieldCopy.CanInterface(), testName+" + fieldCopy.CanInterface()") {
			test.IsTrue(t, reflect.DeepEqual(fieldCopy.Interface(), value),
				testName+" + fieldCopy contents")
		}
	}
}

func checkFieldValueNOK(t *testing.T, s reflect.Value, fieldName string) {
	t.Helper()

	testName := "field " + fieldName

	fieldOrig := s.FieldByName(fieldName)
	test.IsFalse(t, fieldOrig.CanInterface(), testName+" + fieldOrig.CanInterface()")

	_, ok := dark.CopyValue(fieldOrig)
	test.IsFalse(t, ok, "Could not copy "+testName)
}

func TestCopyValue(t *testing.T) {
	// Note that even if all the fields are public, a Struct cannot be copied
	type SubPublic struct {
		Public int
	}

	type SubPrivate struct {
		private int //nolint: unused,megacheck,staticcheck
	}

	type Private struct {
		boolean  bool
		integer  int
		uinteger uint
		cplx     complex128
		flt      float64
		str      string
		array    [3]any
		slice    []any
		hash     map[any]any
		pint     *int
		iface    any
		fn       func()
	}

	//
	// Copy OK
	num := 123
	private := Private{
		boolean: true,
		integer: 42,
		cplx:    complex(2, -2),
		flt:     1.234,
		str:     "foobar",
		array:   [3]any{1, 2, SubPublic{Public: 3}},
		slice:   append(make([]any, 0, 10), 4, 5, SubPublic{Public: 6}),
		hash: map[any]any{
			"foo":                 &SubPublic{Public: 34},
			SubPublic{Public: 78}: 42,
		},
		pint:  &num,
		iface: &num,
	}
	privateStruct := reflect.ValueOf(private)

	checkFieldValueOK(t, privateStruct, "boolean", private.boolean)
	checkFieldValueOK(t, privateStruct, "integer", private.integer)
	checkFieldValueOK(t, privateStruct, "uinteger", private.uinteger)
	checkFieldValueOK(t, privateStruct, "cplx", private.cplx)
	checkFieldValueOK(t, privateStruct, "flt", private.flt)
	checkFieldValueOK(t, privateStruct, "str", private.str)
	checkFieldValueOK(t, privateStruct, "array", private.array)
	checkFieldValueOK(t, privateStruct, "slice", private.slice)
	checkFieldValueOK(t, privateStruct, "hash", private.hash)
	checkFieldValueOK(t, privateStruct, "pint", private.pint)
	checkFieldValueOK(t, privateStruct, "iface", private.iface)

	//
	// Not able to copy...
	private = Private{
		array: [3]any{1, 2, SubPrivate{}},
		slice: append(make([]any, 0, 10), &SubPrivate{}, &SubPrivate{}),
		hash:  map[any]any{"foo": &SubPrivate{}},
		iface: &SubPrivate{},
		fn:    func() {},
	}
	privateStruct = reflect.ValueOf(private)

	checkFieldValueNOK(t, privateStruct, "array")
	checkFieldValueNOK(t, privateStruct, "slice")
	checkFieldValueNOK(t, privateStruct, "hash")
	checkFieldValueNOK(t, privateStruct, "iface")
	checkFieldValueNOK(t, privateStruct, "fn")

	private.hash = map[any]any{SubPrivate{}: 123}
	privateStruct = reflect.ValueOf(private)
	checkFieldValueNOK(t, privateStruct, "hash")

	//
	// nil cases
	private = Private{}
	privateStruct = reflect.ValueOf(private)
	checkFieldValueOK(t, privateStruct, "slice", private.slice)
	checkFieldValueOK(t, privateStruct, "hash", private.hash)
	checkFieldValueOK(t, privateStruct, "pint", private.pint)
	checkFieldValueOK(t, privateStruct, "iface", private.iface)
}
