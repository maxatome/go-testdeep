// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package str

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
)

type myTestDeepStringer struct {
	types.TestDeepStamp
}

func (m myTestDeepStringer) String() string {
	return "TesT!"
}

func TestToString(t *testing.T) {
	for _, curTest := range []struct {
		Got      interface{}
		Expected string
	}{
		{Got: "foobar", Expected: `"foobar"`},
		{Got: "foo\rbar", Expected: `(string) (len=7) "foo\rbar"`},
		{Got: "foo\u2028bar", Expected: `(string) (len=9) "foo\u2028bar"`},
		{Got: reflect.ValueOf("foobar"), Expected: `"foobar"`},
		{Got: types.RawString("test"), Expected: "test"},
		{Got: types.RawInt(42), Expected: "42"},
		{Got: myTestDeepStringer{}, Expected: "TesT!"},
		{Got: 42, Expected: "(int) 42"},
	} {
		test.EqualStr(t, ToString(curTest.Got), curTest.Expected)
	}
}
