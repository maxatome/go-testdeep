// Copyright (c) 2020, 2021, Maxime Soulé
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
			` "foo bar \" \\ \/ \b \f \n\r \t \u20ac \u10e6 \u10E6 héhô" `,
			`"\""`,
			`"\\"`,
			`"\/"`,
			`"\b"`,
			`"\f"`,
			`"\n"`,
			`"\r"`,
			`"\t"`,
			`"\u20ac"`,
			`"zz\""`,
			`"zz\\"`,
			`"zz\/"`,
			`"zz\b"`,
			`"zz\f"`,
			`"zz\n"`,
			`"zz\r"`,
			`"zz\t"`,
			`"zz\u20ac"`,
			`["74.99 \u20ac"]`,
			`{"text": "74.99 \u20ac"}`,
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

	t.Run("JSON spec infringements", func(t *testing.T) {
		check := func(gotJSON, expectedJSON string) {
			t.Helper()
			var expected interface{}
			err := ejson.Unmarshal([]byte(expectedJSON), &expected)
			if err != nil {
				t.Fatalf("bad JSON: %s", err)
			}

			got, err := json.Parse([]byte(gotJSON))
			if !test.NoError(t, err, "json.Parse succeeds") {
				return
			}
			if !reflect.DeepEqual(got, expected) {
				test.EqualErrorMessage(t,
					strings.TrimRight(spew.Sdump(got), "\n"),
					strings.TrimRight(spew.Sdump(expected), "\n"),
					"got matches expected",
				)
			}
		}
		// "," is accepted just before non-empty "}" or "]"
		check(`{"foo": "bar", }`, `{"foo":"bar"}`)
		check(`{"foo":"bar",}`, `{"foo":"bar"}`)
		check(`[ 1, 2, 3, ]`, `[1,2,3]`)
		check(`[ 1,2,3,]`, `[1,2,3]`)
	})

	t.Run("Special string cases", func(t *testing.T) {
		for i, tst := range []struct{ in, expected string }{
			{
				in:       `"$"`,
				expected: `$`,
			},
			{
				in:       `"$$"`,
				expected: `$`,
			},
			{
				in:       `"$$toto"`,
				expected: `$toto`,
			},
		} {
			got, err := json.Parse([]byte(tst.in))
			if !test.NoError(t, err, "#%d, json.Parse succeeds", i) {
				continue
			}

			if !reflect.DeepEqual(got, tst.expected) {
				test.EqualErrorMessage(t,
					strings.TrimRight(spew.Sdump(got), "\n"),
					strings.TrimRight(spew.Sdump(tst.expected), "\n"),
					"#%d is OK", i,
				)
			}
		}
	})

	t.Run("Placeholder cases", func(t *testing.T) {
		for i, js := range []string{
			`  $2  `,
			` "$2" `,
			`  $ph  `,
			` "$ph" `,
			`  $héhé  `,
			` "$héhé" `,
		} {
			got, err := json.Parse([]byte(js), json.ParseOpts{
				Placeholders: []interface{}{"foo", "bar"},
				PlaceholdersByName: map[string]interface{}{
					"ph":   "bar",
					"héhé": "bar",
				},
			})
			if !test.NoError(t, err, "#%d, json.Parse succeeds", i) {
				continue
			}

			if !reflect.DeepEqual(got, `bar`) {
				test.EqualErrorMessage(t,
					strings.TrimRight(spew.Sdump(got), "\n"),
					strings.TrimRight(spew.Sdump(`bar`), "\n"),
					"#%d is OK", i,
				)
			}
		}
	})

	t.Run("Comments", func(t *testing.T) {
		for i, js := range []string{
			"  // comment\ntrue",
			"  true // comment\n   ",
			"  true // comment\n",
			"  true // comment",
			"  /* comment\nmulti\nline */true",
			"  true /* comment\nmulti\nline */",
			"  true /* comment\nmulti\nline */  \t",
			"  true /* comment\nmulti\nline */  // comment",
			"/**///\ntrue/**/",
		} {
			for j, s := range []string{
				js,
				strings.Replace(js, "\n", "\r", -1),   //nolint: gocritic
				strings.Replace(js, "\n", "\r\n", -1), //nolint: gocritic
			} {
				got, err := json.Parse([]byte(s))
				if !test.NoError(t, err, "#%d/%d, json.Parse succeeds", i, j) {
					continue
				}

				if !reflect.DeepEqual(got, true) {
					test.EqualErrorMessage(t,
						got, true,
						"#%d/%d is OK", i, j,
					)
				}
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
				err: `numeric placeholder "$1", but no params given at line 1:2 (pos 2)`,
			},
			{
				js:  `  "$1" `,
				err: `numeric placeholder "$1", but no params given at line 1:3 (pos 3)`,
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
			// multiple errors
			{
				js: "[$1,$2,",
				err: `numeric placeholder "$1", but no params given at line 1:1 (pos 1)
numeric placeholder "$2", but no params given at line 1:4 (pos 4)
syntax error: unexpected EOF at line 1:6 (pos 6)`,
			},
		} {
			_, err := json.Parse([]byte(tst.js))
			if test.Error(t, err, `#%d \n, json.Parse fails`, i) {
				test.EqualStr(t, err.Error(), tst.err, `#%d \n, err OK`, i)
			}

			_, err = json.Parse([]byte(strings.Replace(tst.js, "\n", "\r", -1))) //nolint: gocritic
			if test.Error(t, err, `#%d \r, json.Parse fails`, i) {
				test.EqualStr(t, err.Error(), tst.err, `#%d \r, err OK`, i)
			}

			_, err = json.Parse([]byte(strings.Replace(tst.js, "\n", "\r\n", -1))) //nolint: gocritic
			if test.Error(t, err, `#%d \r\n, json.Parse fails`, i) {
				test.EqualStr(t, err.Error(), tst.err, `#%d \r\n, err OK`, i)
			}
		}

		_, err := json.Parse(
			[]byte(`[$2]`),
			json.ParseOpts{Placeholders: []interface{}{1}},
		)
		if test.Error(t, err) {
			test.EqualStr(t, err.Error(),
				`numeric placeholder "$2", but only one param given at line 1:1 (pos 1)`)
		}

		_, err = json.Parse(
			[]byte(`[$3]`),
			json.ParseOpts{Placeholders: []interface{}{1, 2}},
		)
		if test.Error(t, err) {
			test.EqualStr(t, err.Error(),
				`numeric placeholder "$3", but only 2 params given at line 1:1 (pos 1)`)
		}

		var anyOpPos json.Position
		_, err = json.Parse([]byte(`  KnownOp(  AnyOp()  )`),
			json.ParseOpts{
				OpFn: func(op json.Operator, pos json.Position) (interface{}, error) {
					if op.Name == "KnownOp" {
						return "OK", nil
					}
					anyOpPos = pos
					return nil, fmt.Errorf("hmm weird operator %q", op.Name)
				},
			})
		if test.Error(t, err, "json.Parse fails") {
			test.EqualInt(t, anyOpPos.Pos, 12)
			test.EqualInt(t, anyOpPos.Line, 1)
			test.EqualInt(t, anyOpPos.Col, 12)
			test.EqualStr(t, err.Error(),
				`hmm weird operator "AnyOp" at line 1:12 (pos 12)`)
		}

		for _, js := range []string{
			`  [ $^KnownOp,    $^AnyOp ]`,
			`  [ "$^KnownOp", "$^AnyOp" ]`,
		} {
			_, err := json.Parse([]byte(js),
				json.ParseOpts{
					OpShortcutFn: func(name string, pos json.Position) (interface{}, bool) {
						if name == "KnownOp" {
							return "OK", true
						}
						anyOpPos = pos
						return nil, false
					},
				})
			if test.Error(t, err, "json.Parse fails", js) {
				test.EqualInt(t, anyOpPos.Pos, 18)
				test.EqualInt(t, anyOpPos.Line, 1)
				test.EqualInt(t, anyOpPos.Col, 18)
				test.EqualStr(t, err.Error(),
					`bad operator shortcut "$^AnyOp" at line 1:18 (pos 18)`,
					js)
			}
		}
	})
}
