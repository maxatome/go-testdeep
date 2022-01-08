// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"reflect"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestCanonStructField(t *testing.T) {
	for _, tst := range []struct{ got, expected string }{
		{"", ""},
		{"pipo", "pipo"},
		{">pipo", ">pipo"},
		{">    pipo  ", ">pipo"},
		{"123=~.*", "123=~.*"},
		{"  123  =~  .*   ", "123=~.*"},
		{"&badField", "&badField"},
	} {
		test.EqualStr(t, canonStructField(tst.got), tst.expected)
	}
}

func TestMergeStructFields(t *testing.T) {
	sfs := mergeStructFields()
	if sfs != nil {
		t.Errorf("not nil")
	}

	x := StructFields{}
	sfs = mergeStructFields(x)
	if reflect.ValueOf(sfs).Pointer() != reflect.ValueOf(x).Pointer() {
		t.Errorf("not x")
	}

	a := StructFields{"pipo": 1}
	b := StructFields{"pipo": 2}
	c := StructFields{"pipo": 3}
	sfs = mergeStructFields(a, b, c)
	if reflect.ValueOf(sfs).Pointer() == reflect.ValueOf(c).Pointer() {
		t.Errorf("is c")
	}
	test.EqualInt(t, len(sfs), 1)
	test.EqualInt(t, sfs["pipo"].(int), 3)

	a = StructFields{">pipo": 1}
	b = StructFields{">  pipo": 2}
	c = StructFields{">pipo  ": 3}
	sfs = mergeStructFields(a, b, c)
	if reflect.ValueOf(sfs).Pointer() == reflect.ValueOf(c).Pointer() {
		t.Errorf("is c")
	}
	test.EqualInt(t, len(sfs), 1)
	test.EqualInt(t, sfs[">pipo  "].(int), 3)

	a = StructFields{"1=~pipo": 1}
	b = StructFields{"  1  =~ pipo  ": 2}
	c = StructFields{"1\t=~\tpipo": 3}
	sfs = mergeStructFields(a, b, c)
	if reflect.ValueOf(sfs).Pointer() == reflect.ValueOf(c).Pointer() {
		t.Errorf("is c")
	}
	test.EqualInt(t, len(sfs), 1)
	test.EqualInt(t, sfs["1\t=~\tpipo"].(int), 3)
}

func TestFieldMatcher(t *testing.T) {
	_, err := newFieldMatcher("pipo", 123)
	if test.Error(t, err) {
		if err != errNotAMatcher {
			t.Errorf("got %q, but %q was expected", err, errNotAMatcher)
		}
	}

	for _, tst := range []struct {
		name  string
		order int
		match bool
	}{
		// Regexp
		{name: "=~.*", match: true},
		{name: "=~bc", match: true},
		{name: "=~3$", match: true},
		{name: "!~^b", match: false},
		{name: "134=~bc", match: true, order: 134},
		{name: "134 =~ bc", match: true, order: 134},
		{name: " 134 =~ bc", match: true, order: 134},
		// Shell pattern
		{name: "=*", match: true},
		{name: "=*bc*", match: true},
		{name: "=*3", match: true},
		{name: "!b*", match: false},
		{name: "134=*", match: true, order: 134},
		{name: "134 = *", match: true, order: 134},
		{name: " 134 = *", match: true, order: 134},
	} {
		fm, err := newFieldMatcher(tst.name, 123)
		test.NoError(t, err, tst.name)
		test.EqualStr(t, fm.name, tst.name, tst.name)
		test.EqualInt(t, fm.expected.(int), 123, tst.name)
		test.EqualInt(t, fm.order, tst.order, tst.name)
		test.EqualBool(t, fm.ok, strings.ContainsRune(tst.name, '='), tst.name)
		if test.IsTrue(t, fm.match != nil, tst.name) {
			ok, err := fm.match("abc123")
			test.NoError(t, err, tst.name)
			test.EqualBool(t, ok, tst.match)
		}
	}

	_, err = newFieldMatcher("=~bad(*", 123)
	if test.Error(t, err) {
		test.IsTrue(t, strings.HasPrefix(err.Error(), "bad regexp field `=~bad(*`: "))
	}
}
