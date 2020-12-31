// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package json_test

import (
	ejson "encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/maxatome/go-testdeep/internal/json"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestJSON(t *testing.T) {
	t.Run("Basics", func(t *testing.T) {
		for i, js := range []string{
			`true`,
			`   true   `,
			"\t\nfalse \n   ",
			`   null  `,
			`{}`,
			`[]`,
			`  123.456   `,
			`  123.456e4   `,
			`  123.456E-4   `,
			`  -123e-4   `,
			`0`,
			`""`,
			`"123.456$"`,
			` "foo bar \" \\ \/ \b \f \n\r \t \u2665 héhô" `,
			`[ 1, 2,3, 4 ]`,
			`{"foo":{"bar":true},"zip":1234}`,
		} {
			js := []byte(js)

			var expected interface{}
			err := ejson.Unmarshal(js, &expected)
			if err != nil {
				t.Fatalf("#%d, bad JSON: %s", i, err)
			}

			got, err := json.Parse(js)
			if !test.NoError(t, err, "#%d, json.Parse succeeds", i) {
				continue
			}

			if !reflect.DeepEqual(got, expected) {
				test.EqualErrorMessage(t,
					strings.TrimRight(spew.Sdump(got), "\n"),
					strings.TrimRight(spew.Sdump(expected), "\n"),
					"#%d is OK", i,
				)
			}
		}
	})

	t.Run("Comments", func(t *testing.T) {
		for i, js := range []string{
			"  // comment\ntrue",
			"  true // comment\n   ",
			"  true // comment",
			"  /* comment\nmulti\nline */true",
			"  true /* comment\nmulti\nline */",
			"  true /* comment\nmulti\nline */  \t",
			"  true /* comment\nmulti\nline */  // comment",
			"/**///\ntrue/**/",
		} {
			got, err := json.Parse([]byte(js))
			if !test.NoError(t, err, "#%d, json.Parse succeeds", i) {
				continue
			}

			if !reflect.DeepEqual(got, true) {
				test.EqualErrorMessage(t,
					got, true,
					"#%d is OK", i,
				)
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for i, tst := range []struct{ js, err string }{
			// comment
			{
				js:  " \n   /* unterminated",
				err: "multi-lines comment not terminated at line 2:3 (pos 5)",
			},
			{
				js:  " \n /",
				err: "syntax error: unexpected '/' at line 2:1 (pos 3)",
			},
			{
				js:  " \n /toto",
				err: "syntax error: unexpected '/' at line 2:1 (pos 3)",
			},
			// string
			{
				js:  "/* multi\nline\ncomment */ \"...",
				err: "unterminated string at line 3:11 (pos 25)",
			},
			{
				js:  `  "unterminated\`,
				err: "unterminated string at line 1:2 (pos 2)",
			},
			{
				js:  `"bad escape \a"`,
				err: "invalid escape sequence at line 1:13 (pos 13)",
			},
			{
				js:  `"bad échappe \u123t"`,
				err: "invalid escape sequence at line 1:14 (pos 14)",
			},
			{
				js:  "\"bad rune \007\"",
				err: "invalid character in string at line 1:10 (pos 10)",
			},
			// number
			{
				js:  "  \n 123.345.45",
				err: "invalid number at line 2:1 (pos 4)",
			},
			// dollar token
			{
				js:  "  \n 123.345$",
				err: "syntax error: unexpected '$' at line 2:8 (pos 11)",
			},
			{
				js:  `  $123a `,
				err: "invalid numeric placeholder at line 1:2 (pos 2)",
			},
			{
				js:  `  "$123a" `,
				err: "invalid numeric placeholder at line 1:3 (pos 3)",
			},
			{
				js:  `  $00 `,
				err: `invalid numeric placeholder "$00", it should start at "$1" at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$00" `,
				err: `invalid numeric placeholder "$00", it should start at "$1" at line 1:3 (pos 3)`,
			},
			{
				js:  `  $1 `,
				err: `numeric placeholder "$1", but only 0 param(s) given at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$1" `,
				err: `numeric placeholder "$1", but only 0 param(s) given at line 1:3 (pos 3)`,
			},
			{
				js:  `  $^AnyOp `,
				err: `bad operator shortcut "$^AnyOp" at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$^AnyOp" `,
				err: `bad operator shortcut "$^AnyOp" at line 1:3 (pos 3)`,
			},
			{
				js:  `  $tag%`,
				err: `bad placeholder "$tag%" at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$tag%"`,
				err: `bad placeholder "$tag%" at line 1:3 (pos 3)`,
			},
			{
				js:  `  $tag`,
				err: `unknown placeholder "$tag" at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$tag"`,
				err: `unknown placeholder "$tag" at line 1:3 (pos 3)`,
			},
			// operator
			{
				js:  "  AnyOpé",
				err: `invalid operator name "AnyOp\xc3" at line 1:2 (pos 2)`,
			},
			{
				js:  "  AnyOp()",
				err: `unknown operator "AnyOp" at line 1:2 (pos 2)`,
			},
			// syntax error
			{
				js:  "  \n 123.345true",
				err: "syntax error: unexpected TRUE at line 2:8 (pos 11)",
			},
			{
				js:  "  \n 123.345%",
				err: "syntax error: unexpected '%' at line 2:8 (pos 11)",
			},
			{
				js:  "  \n 123.345\x1f",
				err: `syntax error: unexpected '\u001f' at line 2:8 (pos 11)`,
			},
			{
				js:  "  \n 123.345\U0002f500",
				err: `syntax error: unexpected '\U0002f500' at line 2:8 (pos 11)`,
			},
		} {
			_, err := json.Parse([]byte(tst.js))
			if test.Error(t, err, "#%d, json.Parse fails", i) {
				test.EqualStr(t, err.Error(), tst.err, "#%d, err OK", i)
			}
		}

		_, err := json.Parse([]byte(`  KnownOp(  AnyOp()  )`),
			json.ParseOpts{
				OpFn: func(op json.Operator) (interface{}, error) {
					if op.Name == "KnownOp" {
						return "OK", nil
					}
					return nil, fmt.Errorf("unknown operator %q", op.Name)
				},
			})
		if test.Error(t, err, "json.Parse fails") {
			test.EqualStr(t, err.Error(),
				`unknown operator "AnyOp" at line 1:12 (pos 12)`)
		}
	})
}
