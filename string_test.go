// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestToString(t *testing.T) {
	for _, curTest := range []struct {
		Got      interface{}
		Expected string
	}{
		{Got: "foobar", Expected: `"foobar"`},
		{Got: "foo\rbar", Expected: `(string) (len=7) "foo\rbar"`},
		{Got: "foo\u2028bar", Expected: `(string) (len=9) "foo\u2028bar"`},
		{Got: reflect.ValueOf("foobar"), Expected: `"foobar"`},
		{Got: rawString("test"), Expected: "test"},
		{Got: rawInt(42), Expected: "42"},
		{Got: Nil(), Expected: "nil"},
		{Got: 42, Expected: "(int) 42"},
	} {
		test.EqualStr(t, toString(curTest.Got), curTest.Expected)
	}
}
