// Copyright (c) 2020-2022, Maxime Soulé
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

func checkJSON(t *testing.T, gotJSON, expectedJSON string) {
	t.Helper()

	var expected any
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

			var expected any
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
		for _, tc := range []struct{ got, expected string }{
			// "," is accepted just before non-empty "}" or "]"
			{`{"foo": "bar", }`, `{"foo":"bar"}`},
			{`{"foo":"bar",}`, `{"foo":"bar"}`},
			{`[ 1, 2, 3, ]`, `[1,2,3]`},
			{`[ 1,2,3,]`, `[1,2,3]`},

			// No need to escape \n, \r & \t
			{"\"\n\r\t\"", `"\n\r\t"`},

			// Extend to golang accepted numbers
			// as int64
			{`+42`, `42`},

			{`0600`, `384`},
			{`-0600`, `-384`},
			{`+0600`, `384`},

			{`0xBadFace`, `195951310`},
			{`-0xBadFace`, `-195951310`},
			{`+0xBadFace`, `195951310`},

			// as float64
			{`0600.123`, `600.123`}, // float64 can not be an octal number
			{`0600.`, `600`},        // float64 can not be an octal number
			{`.25`, `0.25`},
			{`+123.`, `123`},

			// Extend to golang 1.13 accepted numbers
			// as int64
			{`4_2`, `42`},
			{`+4_2`, `42`},
			{`-4_2`, `-42`},

			{`0b101010`, `42`},
			{`-0b101010`, `-42`},
			{`+0b101010`, `42`},

			{`0b10_1010`, `42`},
			{`-0b_10_1010`, `-42`},
			{`+0b10_10_10`, `42`},

			{`0B101010`, `42`},
			{`-0B101010`, `-42`},
			{`+0B101010`, `42`},

			{`0B10_1010`, `42`},
			{`-0B_10_1010`, `-42`},
			{`+0B10_10_10`, `42`},

			{`0_600`, `384`},
			{`-0_600`, `-384`},
			{`+0_600`, `384`},

			{`0o600`, `384`},
			{`0o_600`, `384`},
			{`-0o600`, `-384`},
			{`-0o6_00`, `-384`},
			{`+0o600`, `384`},
			{`+0o60_0`, `384`},

			{`0O600`, `384`},
			{`0O_600`, `384`},
			{`-0O600`, `-384`},
			{`-0O6_00`, `-384`},
			{`+0O600`, `384`},
			{`+0O60_0`, `384`},

			{`0xBad_Face`, `195951310`},
			{`-0x_Bad_Face`, `-195951310`},
			{`+0xBad_Face`, `195951310`},

			{`0XBad_Face`, `195951310`},
			{`-0X_Bad_Face`, `-195951310`},
			{`+0XBad_Face`, `195951310`},

			// as float64
			{`0_600.123`, `600.123`}, // float64 can not be an octal number
			{`1_5.`, `15`},
			{`0.15e+0_2`, `15`},
			{`0x1p-2`, `0.25`},
			{`0x2.p10`, `2048`},
			{`0x1.Fp+0`, `1.9375`},
			{`0X.8p-0`, `0.5`},
			{`0X_1FFFP-16`, `0.1249847412109375`},

			// Raw strings
			{`r"pipo"`, `"pipo"`},
			{`r "pipo"`, `"pipo"`},
			{"r\n'pipo'", `"pipo"`},
			{`r%pipo%`, `"pipo"`},
			{`r·pipo·`, `"pipo"`},
			{"r`pipo`", `"pipo"`},
			{`r/pipo/`, `"pipo"`},
			{"r //comment\n`pipo`", `"pipo"`}, // comments accepted bw r and string
			{"r//comment\n`pipo`", `"pipo"`},
			{"r/*comment\n*/|pipo|", `"pipo"`},
			{"r(p\ni\rp\to)", `"p\ni\rp\to"`}, // accepted raw whitespaces
			{`r@pi\po\@`, `"pi\\po\\"`},       // backslash has no meaning
			// balanced delimiters
			{`r(p(i(hey)p)o)`, `"p(i(hey)p)o"`},
			{`r{p{i{hey}p}o}`, `"p{i{hey}p}o"`},
			{`r[p[i[hey]p]o]`, `"p[i[hey]p]o"`},
			{`r<p<i<hey>p>o>`, `"p<i<hey>p>o"`},
			{`r(pipo)`, `"pipo"`},
			{"r \t\n(pipo)", `"pipo"`},
			{`r{pipo}`, `"pipo"`},
			{`r[pipo]`, `"pipo"`},
			{`r<pipo>`, `"pipo"`},
			// Not balanced
			{`r)pipo)`, `"pipo"`},
			{`r}pipo}`, `"pipo"`},
			{`r]pipo]`, `"pipo"`},
			{`r>pipo>`, `"pipo"`},
		} {
			t.Run(tc.got, func(t *testing.T) {
				checkJSON(t, tc.got, tc.expected)
			})
		}
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
				Placeholders: []any{"foo", "bar"},
				PlaceholdersByName: map[string]any{
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
				strings.ReplaceAll(js, "\n", "\r"),
				strings.ReplaceAll(js, "\n", "\r\n"),
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

	t.Run("OK", func(t *testing.T) {
		opts := json.ParseOpts{
			OpFn: func(op json.Operator, pos json.Position) (any, error) {
				if op.Name == "KnownOp" {
					return "OK", nil
				}
				return nil, fmt.Errorf("hmm weird operator %q", op.Name)
			},
		}
		for _, js := range []string{
			`[ KnownOp ]`,
			`[ KnownOp() ]`,
			`[ $^KnownOp() ]`,
			`[ $^KnownOp ]`,
			`[ KnownOp($^KnownOp) ]`,
			`[ KnownOp( $^KnownOp() ) ]`,
			`[ $^KnownOp(KnownOp) ]`,
		} {
			_, err := json.Parse([]byte(js), opts)
			test.NoError(t, err, "json.Parse OK", js)
		}
	})

	t.Run("Reentrant parser", func(t *testing.T) {
		opts := json.ParseOpts{
			OpFn: func(op json.Operator, pos json.Position) (any, error) {
				if op.Name == "KnownOp" {
					return "OK", nil
				}
				return nil, fmt.Errorf("hmm weird operator %q", op.Name)
			},
		}
		for _, js := range []string{
			`[ "$^KnownOp(1, 2, 3)" ]`,
			`[ "$^KnownOp(1, 2, 3)  " ]`,
			`[ "$^KnownOp(r<$^KnownOp(11, 12)>, 2, KnownOp(31, 32))" ]`,
		} {
			_, err := json.Parse([]byte(js), opts)
			test.NoError(t, err, "json.Parse OK", js)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for i, tst := range []struct{ nam, js, err string }{
			// comment
			{
				nam: "unterminated comment",
				js:  " \n   /* unterminated",
				err: "multi-lines comment not terminated at line 2:3 (pos 5)",
			},
			{
				nam: "/ at EOF",
				js:  " \n /",
				err: "syntax error: unexpected '/' at line 2:1 (pos 3)",
			},
			{
				nam: "/toto",
				js:  " \n /toto",
				err: "syntax error: unexpected '/' at line 2:1 (pos 3)",
			},
			// string
			{
				nam: "unterminated string+multi lines",
				js:  "/* multi\nline\ncomment */ \"...",
				err: "unterminated string at line 3:11 (pos 25)",
			},
			{
				nam: "unterminated string",
				js:  `  "unterminated\`,
				err: "unterminated string at line 1:2 (pos 2)",
			},
			{
				nam: "bad escape",
				js:  `"bad escape \a"`,
				err: "invalid escape sequence at line 1:13 (pos 13)",
			},
			{
				nam: `bad escape \u`,
				js:  `"bad échappe \u123t"`,
				err: "invalid escape sequence at line 1:14 (pos 14)",
			},
			{
				nam: "bad rune",
				js:  "\"bad rune \007\"",
				err: "invalid character in string at line 1:10 (pos 10)",
			},
			// number
			{
				nam: "bad number",
				js:  "  \n 123.345.45",
				err: "invalid number at line 2:1 (pos 4)",
			},
			// dollar token
			{
				nam: "dollar at EOF",
				js:  "  $",
				err: "syntax error: unexpected '$' at line 1:2 (pos 2)",
			},
			{
				nam: "dollar alone",
				js:  "  $  ",
				err: "syntax error: unexpected '$' at line 1:2 (pos 2)",
			},
			{
				nam: "multi lines+dollar at EOF",
				js:  "  \n 123.345$",
				err: "syntax error: unexpected '$' at line 2:8 (pos 11)",
			},
			{
				nam: "bad num placeholder",
				js:  `  $123a `,
				err: "invalid numeric placeholder at line 1:2 (pos 2)",
			},
			{
				nam: "bad num placeholder in string",
				js:  `  "$123a" `,
				err: "invalid numeric placeholder at line 1:3 (pos 3)",
			},
			{
				nam: "bad 0 placeholder",
				js:  `  $00 `,
				err: `invalid numeric placeholder "$00", it should start at "$1" at line 1:2 (pos 2)`,
			},
			{
				nam: "bad 0 placeholder in string",
				js:  `  "$00" `,
				err: `invalid numeric placeholder "$00", it should start at "$1" at line 1:3 (pos 3)`,
			},
			{
				nam: "placeholder/params mismatch",
				js:  `  $1 `,
				err: `numeric placeholder "$1", but no params given at line 1:2 (pos 2)`,
			},
			{
				nam: "placeholder in string/params mismatch",
				js:  `[ "$1", 1, 2 ] `,
				err: `numeric placeholder "$1", but no params given at line 1:3 (pos 3)`,
			},
			{
				nam: "invalid operator in string",
				js:  ` "$^UnknownAndBad>" `,
				err: `invalid operator name "UnknownAndBad>" at line 1:4 (pos 4)`,
			},
			{
				nam: "unknown operator close paren",
				js:  ` UnknownAndBad)`,
				err: `unknown operator "UnknownAndBad" at line 1:1 (pos 1)`,
			},
			{
				nam: "unknown operator close paren in string",
				js:  ` "$^UnknownAndBad)" `,
				err: `unknown operator "UnknownAndBad" at line 1:4 (pos 4)`,
			},
			{
				nam: "op and syntax error",
				js:  ` KnownOp)`,
				err: `syntax error: unexpected ')' at line 1:8 (pos 8)`,
			},
			{
				nam: "op in string and syntax error",
				js:  ` "$^KnownOp)" `,
				err: `syntax error: unexpected ')' at line 1:11 (pos 11)`,
			},
			{
				nam: "op paren in string and syntax error",
				js:  ` "$^KnownOp())" `,
				err: `syntax error: unexpected ')' at line 1:13 (pos 13)`,
			},
			{
				nam: "invalid $^",
				js:  `  $^. `,
				err: `$^ must be followed by an operator name at line 1:2 (pos 2)`,
			},
			{
				nam: "invalid $^ in string",
				js:  ` "$^."`,
				err: `$^ must be followed by an operator name at line 1:2 (pos 2)`,
			},
			{
				nam: "invalid $^ at EOF",
				js:  `  $^`,
				err: `$^ must be followed by an operator name at line 1:2 (pos 2)`,
			},
			{
				nam: "invalid $^ in string at EOF",
				js:  ` "$^"`,
				err: `$^ must be followed by an operator name at line 1:2 (pos 2)`,
			},
			{
				nam: "bad placeholder",
				js:  `  $tag%`,
				err: `bad placeholder "$tag%" at line 1:2 (pos 2)`,
			},
			{
				nam: "bad placeholder in string",
				js:  `  "$tag%"`,
				err: `bad placeholder "$tag%" at line 1:3 (pos 3)`,
			},
			{
				nam: "unknown placeholder",
				js:  `  $tag`,
				err: `unknown placeholder "$tag" at line 1:2 (pos 2)`,
			},
			{
				nam: "unknown placeholder in string",
				js:  `  "$tag"`,
				err: `unknown placeholder "$tag" at line 1:3 (pos 3)`,
			},
			// operator
			{
				nam: "invalid operator",
				js:  "  AnyOpé",
				err: `invalid operator name "AnyOp\xc3" at line 1:2 (pos 2)`,
			},
			{
				nam: "invalid $^operator",
				js:  "  $^AnyOpé",
				err: `invalid operator name "AnyOp\xc3" at line 1:4 (pos 4)`,
			},
			{
				nam: "invalid $^operator in string",
				js:  `  "$^AnyOpé"`,
				err: `invalid operator name "AnyOp\xc3" at line 1:5 (pos 5)`,
			},
			{
				nam: "unknown operator",
				js:  "  AnyOp",
				err: `unknown operator "AnyOp" at line 1:2 (pos 2)`,
			},
			{
				nam: "unknown operator paren",
				js:  "  AnyOp()",
				err: `unknown operator "AnyOp" at line 1:2 (pos 2)`,
			},
			{
				nam: "unknown $^operator",
				js:  "$^AnyOp",
				err: `unknown operator "AnyOp" at line 1:2 (pos 2)`,
			},
			{
				nam: "unknown $^operator paren",
				js:  "$^AnyOp()",
				err: `unknown operator "AnyOp" at line 1:2 (pos 2)`,
			},
			{
				nam: "unknown $^operator in string",
				js:  `"$^AnyOp"`,
				err: `unknown operator "AnyOp" at line 1:3 (pos 3)`,
			},
			{
				nam: "unknown $^operator paren in string",
				js:  `"$^AnyOp()"`,
				err: `unknown operator "AnyOp" at line 1:3 (pos 3)`,
			},
			{
				nam: "unknown $^operator in rawstring",
				js:  `r<$^AnyOp>`,
				err: `unknown operator "AnyOp" at line 1:4 (pos 4)`,
			},
			{
				nam: "unknown $^operator paren in rawstring",
				js:  `r<$^AnyOp()>`,
				err: `unknown operator "AnyOp" at line 1:4 (pos 4)`,
			},
			// syntax error
			{
				nam: "syntax error num+bool",
				js:  "  \n 123.345true",
				err: "syntax error: unexpected TRUE at line 2:8 (pos 11)",
			},
			{
				nam: "syntax error num+%",
				js:  "  \n 123.345%",
				err: "syntax error: unexpected '%' at line 2:8 (pos 11)",
			},
			{
				nam: "syntax error num+ESC",
				js:  "  \n 123.345\x1f",
				err: `syntax error: unexpected '\u001f' at line 2:8 (pos 11)`,
			},
			{
				nam: "syntax error num+unicode",
				js:  "  \n 123.345\U0002f500",
				err: `syntax error: unexpected '\U0002f500' at line 2:8 (pos 11)`,
			},
			// multiple errors
			{
				nam: "multi errors placeholders",
				js:  "[$1,$2,",
				err: `numeric placeholder "$1", but no params given at line 1:1 (pos 1)
numeric placeholder "$2", but no params given at line 1:4 (pos 4)
syntax error: unexpected EOF at line 1:6 (pos 6)`,
			},
			{
				nam: "multi errors placeholder+operator",
				js:  `[$1,"$^Unknown1()","$^Unknown2()"]`,
				err: `numeric placeholder "$1", but no params given at line 1:1 (pos 1)
invalid operator name "Unknown1" at line 1:7 (pos 7)
invalid operator name "Unknown2" at line 1:22 (pos 22)`,
			},
			// raw strings
			{
				nam: "rawstring start delimiter",
				js:  "  \n   r   ",
				err: `cannot find r start delimiter at line 2:7 (pos 10)`,
			},
			{
				nam: "rawstring start delimiter EOF",
				js:  "  \n   r",
				err: `cannot find r start delimiter at line 2:4 (pos 7)`,
			},
			{
				nam: "rawstring bad delimiter",
				js:  `  rxpipox`,
				err: `invalid r delimiter 'x', should be either a punctuation or a symbol rune, excluding '_' at line 1:3 (pos 3)`,
			},
			{
				nam: "rawstring bad underscore delimiter",
				js:  `  r_pipo_`,
				err: `invalid r delimiter '_', should be either a punctuation or a symbol rune, excluding '_' at line 1:3 (pos 3)`,
			},
			{
				nam: "rawstring bad rune",
				js:  "  r:bad rune \007:",
				err: `invalid character in raw string at line 1:13 (pos 13)`,
			},
			{
				nam: "unterminated rawstring",
				js:  `  r!pipo...`,
				err: `unterminated raw string at line 1:3 (pos 3)`,
			},
		} {
			t.Run(tst.nam, func(t *testing.T) {
				opts := json.ParseOpts{
					OpFn: func(op json.Operator, pos json.Position) (any, error) {
						if op.Name == "KnownOp" {
							return "OK", nil
						}
						return nil, fmt.Errorf("unknown operator %q", op.Name)
					},
				}
				_, err := json.Parse([]byte(tst.js), opts)
				if test.Error(t, err, `#%d \n, json.Parse fails`, i) {
					test.EqualStr(t, err.Error(), tst.err, `#%d \n, err OK`, i)
				}

				_, err = json.Parse([]byte(strings.ReplaceAll(tst.js, "\n", "\r")), opts)
				if test.Error(t, err, `#%d \r, json.Parse fails`, i) {
					test.EqualStr(t, err.Error(), tst.err, `#%d \r, err OK`, i)
				}

				_, err = json.Parse([]byte(strings.ReplaceAll(tst.js, "\n", "\r\n")), opts)
				if test.Error(t, err, `#%d \r\n, json.Parse fails`, i) {
					test.EqualStr(t, err.Error(), tst.err, `#%d \r\n, err OK`, i)
				}
			})
		}

		_, err := json.Parse(
			[]byte(`[$2]`),
			json.ParseOpts{Placeholders: []any{1}},
		)
		if test.Error(t, err) {
			test.EqualStr(t, err.Error(),
				`numeric placeholder "$2", but only one param given at line 1:1 (pos 1)`)
		}

		_, err = json.Parse(
			[]byte(`[$3]`),
			json.ParseOpts{Placeholders: []any{1, 2}},
		)
		if test.Error(t, err) {
			test.EqualStr(t, err.Error(),
				`numeric placeholder "$3", but only 2 params given at line 1:1 (pos 1)`)
		}

		for _, js := range []string{
			`     KnownOp(  AnyOp()  )`,
			`     KnownOp(  AnyOp  )`,
			`    KnownOp("$^AnyOp()" )`,
			`    KnownOp("$^AnyOp" )`,
			`    KnownOp( $^AnyOp() )`,
			`  $^KnownOp(   AnyOp )`,
			` "$^KnownOp(   AnyOp )"`,
			` "$^KnownOp(   AnyOp() )"`,
			` "$^KnownOp( $^AnyOp() )"`,
			`"$^KnownOp(r'$^AnyOp()')"`,
		} {
			t.Run(js, func(t *testing.T) {
				var anyOpPos json.Position
				_, err = json.Parse([]byte(js), json.ParseOpts{
					OpFn: func(op json.Operator, pos json.Position) (any, error) {
						if op.Name == "KnownOp" {
							return "OK", nil
						}
						anyOpPos = pos
						return nil, fmt.Errorf("hmm weird operator %q", op.Name)
					},
				})
				if test.Error(t, err, "json.Parse fails") {
					test.EqualInt(t, anyOpPos.Pos, 15)
					test.EqualInt(t, anyOpPos.Line, 1)
					test.EqualInt(t, anyOpPos.Col, 15)
					test.EqualStr(t, err.Error(),
						`hmm weird operator "AnyOp" at line 1:15 (pos 15)`)
				}
			})
		}
	})

	t.Run("no operators", func(t *testing.T) {
		_, err := json.Parse([]byte("  Operator"))
		if test.Error(t, err, "json.Parse fails") {
			test.EqualStr(t, err.Error(),
				`unknown operator "Operator" at line 1:2 (pos 2)`)
		}
	})
}
