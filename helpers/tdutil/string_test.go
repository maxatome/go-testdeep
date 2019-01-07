// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestFormatString(t *testing.T) {
	for _, curTest := range []struct {
		paramGot string
		expected string
	}{
		{paramGot: "foobar", expected: `"foobar"`},
		{paramGot: "foo\rbar", expected: `(string) (len=7) "foo\rbar"`},
		{paramGot: "foo\u2028bar", expected: `(string) (len=9) "foo\u2028bar"`},
		{paramGot: `foo"bar`, expected: "`foo\"bar`"},
		{paramGot: "foo\n\"bar", expected: "`foo\n\"bar`"},
		{paramGot: "foo`\"\nbar", expected: "(string) (len=9) \"foo`\\\"\\nbar\""},
		{paramGot: "foo`\n\"bar", expected: "(string) (len=9) \"foo`\\n\\\"bar\""},
		{paramGot: "foo\n`\"bar", expected: "(string) (len=9) \"foo\\n`\\\"bar\""},
		{paramGot: "foo\n\"`bar", expected: "(string) (len=9) \"foo\\n\\\"`bar\""},
	} {
		got := tdutil.FormatString(curTest.paramGot)
		if got != curTest.expected {
			t.Errorf(`got "%s" ≠ expected "%s"`, got, curTest.expected)
		}
	}
}
