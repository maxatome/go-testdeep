// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/util"
)

func TestJSONPointer(t *testing.T) {
	var ref any
	err := json.Unmarshal([]byte(`
{
   "foo": ["bar", "baz"],
   "": 0,
   "a/b": 1,
   "c%d": 2,
   "e^f": 3,
   "g|h": 4,
   "i\\j": 5,
   "k\"l": 6,
   " ": 7,
   "m~n": 8
}`),
		&ref)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %s", err)
	}

	checkOK := func(pointer string, expected any) {
		t.Helper()

		got, err := util.JSONPointer(ref, pointer)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got: %v expected: %v", got, expected)
		}
		if err != nil {
			t.Errorf("error <%s> received instead of nil", err)
		}
	}

	checkErr := func(pointer, errExpected string) {
		t.Helper()

		got, err := util.JSONPointer(ref, pointer)
		if got != nil {
			t.Errorf("got: %v expected: nil", got)
		}
		if err == nil {
			t.Errorf("error nil received instead of <%s>", errExpected)
		} else if err.Error() != errExpected {
			t.Errorf("error <%s> received instead of <%s>", err, errExpected)
		}
	}

	checkOK(``, ref)
	checkOK(`/foo`, []any{"bar", "baz"})
	checkOK(`/foo/0`, "bar")
	checkOK(`/`, float64(0))
	checkOK(`/a~1b`, float64(1))
	checkOK(`/c%d`, float64(2))
	checkOK(`/e^f`, float64(3))
	checkOK(`/g|h`, float64(4))
	checkOK(`/i\j`, float64(5))
	checkOK(`/k"l`, float64(6))
	checkOK(`/ `, float64(7))
	checkOK(`/m~0n`, float64(8))

	checkErr("x", "invalid JSON pointer")
	checkErr("/8", "key not found @/8")
	checkErr("/foo/-1/pipo", "array but not an index in JSON pointer @/foo/-1")
	checkErr("/foo/bingo/pipo", "array but not an index in JSON pointer @/foo/bingo")
	checkErr("/foo/2/pipo", "out of array range @/foo/2")
	checkErr("/foo/1/pipo", "not a map nor an array @/foo/1/pipo")
}
