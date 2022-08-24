// Copyright (c) 2020-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

func TestIsStruct(t *testing.T) {
	s := struct{}{}
	ps := &s
	pps := &ps
	m := map[string]struct{}{}

	for i, test := range []struct {
		val any
		ok  bool
	}{
		{val: s, ok: true},
		{val: ps, ok: true},
		{val: pps, ok: true},
		{val: &pps, ok: true},
		{val: m, ok: false},
		{val: &m, ok: false},
	} {
		if types.IsStruct(reflect.TypeOf(test.val)) != test.ok {
			t.Errorf("#%d IsStruct() mismatch as ≠ %t", i, test.ok)
		}
	}
}

func TestIsTypeOrConvertible(t *testing.T) {
	type MyInt int

	ok, convertible := types.IsTypeOrConvertible(reflect.ValueOf(123), reflect.TypeOf(123))
	test.IsTrue(t, ok)
	test.IsFalse(t, convertible)

	ok, convertible = types.IsTypeOrConvertible(reflect.ValueOf(123), reflect.TypeOf(123.45))
	test.IsTrue(t, ok)
	test.IsTrue(t, convertible)

	ok, convertible = types.IsTypeOrConvertible(reflect.ValueOf(123), reflect.TypeOf(MyInt(123)))
	test.IsTrue(t, ok)
	test.IsTrue(t, convertible)

	ok, convertible = types.IsTypeOrConvertible(reflect.ValueOf("xx"), reflect.TypeOf(123))
	test.IsFalse(t, ok)
	test.IsFalse(t, convertible)
}

func TestKindType(t *testing.T) {
	for _, tc := range []struct {
		val      any
		expected string
	}{
		{nil, "nil"},
		{42, "int"},
		{(*int)(nil), "*int"},
		{(*[]int)(nil), "*slice (*[]int type)"},
		{(***int)(nil), "***int"},
	} {
		vval := reflect.ValueOf(tc.val)
		name := "nil"
		if tc.val != nil {
			name = vval.Type().String()
		}
		t.Run(name, func(t *testing.T) {
			test.EqualStr(t, types.KindType(vval), tc.expected)
		})
	}
}
