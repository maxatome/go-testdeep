// Copyright (c) 2019-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type errReader struct{}

// Read implements io.Reader.
func (r errReader) Read(p []byte) (int, error) {
	return 0, errors.New("an error occurred")
}

const (
	insideOpJSON = " inside operator JSON at td_json_test.go:"
	underOpJSON  = "under operator JSON at td_json_test.go:"
)

func TestJSON(t *testing.T) {
	type MyStruct struct {
		Name   string `json:"name"`
		Age    uint   `json:"age"`
		Gender string `json:"gender"`
	}

	//
	// nil
	checkOK(t, nil, td.JSON(`null`))
	checkOK(t, (*int)(nil), td.JSON(`null`))

	//
	// Basic types
	checkOK(t, 123, td.JSON(`  123  `))
	checkOK(t, true, td.JSON(`  true  `))
	checkOK(t, false, td.JSON(`  false  `))
	checkOK(t, "foobar", td.JSON(`  "foobar"  `))

	//
	// struct
	//
	got := MyStruct{Name: "Bob", Age: 42, Gender: "male"}

	// No placeholder
	checkOK(t, got,
		td.JSON(`{"name":"Bob","age":42,"gender":"male"}`))

	checkOK(t, got, td.JSON(`$1`, got)) // json.Marshal() got for $1

	// Numeric placeholders
	checkOK(t, got,
		td.JSON(`{"name":"$1","age":$2,"gender":$3}`,
			"Bob", 42, "male")) // raw values

	checkOK(t, got,
		td.JSON(`{"name":"$1","age":$2,"gender":"$3"}`,
			td.Re(`^Bob`),
			td.Between(40, 45),
			td.NotEmpty()))

	// Same using Flatten
	checkOK(t, got,
		td.JSON(`{"name":"$1","age":$2,"gender":"$3"}`,
			td.Re(`^Bob`),
			td.Flatten([]td.TestDeep{td.Between(40, 45), td.NotEmpty()}),
		))

	// Operators are not JSON marshallable
	checkOK(t, got,
		td.JSON(`$1`, map[string]any{
			"name":   td.Re(`^Bob`),
			"age":    42,
			"gender": td.NotEmpty(),
		}))

	// Placeholder + unmarshal before comparison
	checkOK(t, json.RawMessage(`[1,2,3]`), td.JSON(`$1`, []int{1, 2, 3}))
	checkOK(t, json.RawMessage(`{"foo":[1,2,3]}`),
		td.JSON(`{"foo":$1}`, []int{1, 2, 3}))
	checkOK(t, json.RawMessage(`[1,2,3]`),
		td.JSON(`$1`, []any{1, td.Between(1, 3), 3}))

	// Tag placeholders
	checkOK(t, got,
		td.JSON(`{"name":"$name","age":$age,"gender":$gender}`,
			// raw values
			td.Tag("name", "Bob"), td.Tag("age", 42), td.Tag("gender", "male")))

	checkOK(t, got,
		td.JSON(`{"name":"$name","age":$age,"gender":"$gender"}`,
			td.Tag("name", td.Re(`^Bo`)),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.NotEmpty())))

	// Tag placeholders + numeric placeholders
	checkOK(t, []MyStruct{got, got},
		td.JSON(`[
				{"name":"$1","age":$age,"gender":"$3"},
				{"name":"$1","age":$2,"gender":"$3"}
			]`,
			td.Re(`^Bo`),                      // $1
			td.Tag("age", td.Between(40, 45)), // $2
			"male"))                           // $3

	// Tag placeholders + operators are not JSON marshallable
	checkOK(t, got,
		td.JSON(`$all`, td.Tag("all", map[string]any{
			"name":   td.Re(`^Bob`),
			"age":    42,
			"gender": td.NotEmpty(),
		})))

	checkError(t, got,
		td.JSON(`{"name":$1, "age":$1, "gender":$1}`,
			td.Tag("!!", td.Ignore())),
		expectedError{
			Message: mustBe("bad usage of Tag operator"),
			Summary: mustBe("Invalid tag, should match (Letter|_)(Letter|_|Number)*"),
			Under:   mustContain("under operator Tag"),
		})

	// Tag placeholders + nil
	checkOK(t, nil, td.JSON(`$all`, td.Tag("all", nil)))

	// Mixed placeholders + operator
	for _, op := range []string{
		"NotEmpty",
		"NotEmpty()",
		"$^NotEmpty",
		"$^NotEmpty()",
		`"$^NotEmpty"`,
		`"$^NotEmpty()"`,
		`r<$^NotEmpty>`,
		`r<$^NotEmpty()>`,
	} {
		checkOK(t, got,
			td.JSON(`{"name":"$name","age":$1,"gender":`+op+`}`,
				td.Tag("age", td.Between(40, 45)),
				td.Tag("name", td.Re(`^Bob`))),
			"using operator %s", op)
	}

	checkOK(t, got,
		td.JSON(`{"name":Re("^Bo\\w"),"age":Between(40,45),"gender":NotEmpty()}`))
	checkOK(t, got,
		td.JSON(`
{
  "name":   All(Re("^Bo\\w"), HasPrefix("Bo"), HasSuffix("ob")),
  "age":    Between(40,45),
  "gender": NotEmpty()
}`))
	checkOK(t, got,
		td.JSON(`
{
  "name":   All(Re("^Bo\\w"), HasPrefix("Bo"), HasSuffix("ob")),
  "age":    Between(40,45),
  "gender": NotEmpty
}`))

	// Same but operators in strings using "$^"
	checkOK(t, got,
		td.JSON(`{"name":Re("^Bo\\w"),"age":"$^Between(40,45)","gender":"$^NotEmpty()"}`))
	checkOK(t, got, // using classic "" string, so each \ has to be escaped
		td.JSON(`
{
  "name":   "$^All(Re(\"^Bo\\\\w\"), HasPrefix(\"Bo\"), HasSuffix(\"ob\"))",
  "age":    "$^Between(40,45)",
  "gender": "$^NotEmpty()",
}`))
	checkOK(t, got, // using raw strings, no escape needed
		td.JSON(`
{
  "name":   "$^All(Re(r(^Bo\\w)), HasPrefix(r{Bo}), HasSuffix(r'ob'))",
  "age":    "$^Between(40,45)",
  "gender": "$^NotEmpty()",
}`))

	// …with comments…
	checkOK(t, got,
		td.JSON(`
// This should be the JSON representation of MyStruct struct
{
  // A person:
  "name":   "$name",   // The name of this person
  "age":    $1,        /* The age of this person:
                          - placeholder unquoted, but could be without
                            any change
                          - to demonstrate a multi-lines comment */
  "gender": $^NotEmpty // Operator NotEmpty
}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	before := time.Now()
	timeGot := map[string]time.Time{"created_at": time.Now()}
	checkOK(t, timeGot,
		td.JSON(`{"created_at": Between($1, $2)}`, before, time.Now()))

	checkOK(t, timeGot,
		td.JSON(`{"created_at": $1}`, td.Between(before, time.Now())))

	// Len
	checkOK(t, []int{1, 2, 3}, td.JSON(`Len(3)`))

	//
	// []byte
	checkOK(t, got,
		td.JSON([]byte(`{"name":"$name","age":$1,"gender":"male"}`),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// json.RawMessage
	checkOK(t, got,
		td.JSON(json.RawMessage(`{"name":"$name","age":$1,"gender":"male"}`),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// nil++
	checkOK(t, nil, td.JSON(`$1`, nil))
	checkOK(t, (*int)(nil), td.JSON(`$1`, td.Nil()))

	checkOK(t, nil, td.JSON(`$x`, td.Tag("x", nil)))
	checkOK(t, (*int)(nil), td.JSON(`$x`, td.Tag("x", nil)))

	checkOK(t, json.RawMessage(`{"foo": null}`), td.JSON(`{"foo": null}`))

	checkOK(t,
		json.RawMessage(`{"foo": null}`),
		td.JSON(`{"foo": $1}`, nil))

	checkOK(t,
		json.RawMessage(`{"foo": null}`),
		td.JSON(`{"foo": $1}`, td.Nil()))

	checkOK(t,
		json.RawMessage(`{"foo": null}`),
		td.JSON(`{"foo": $x}`, td.Tag("x", nil)))

	checkOK(t,
		json.RawMessage(`{"foo": null}`),
		td.JSON(`{"foo": $x}`, td.Tag("x", td.Nil())))

	//
	// Loading a file
	tmpDir := t.TempDir()

	filename := tmpDir + "/test.json"
	err := os.WriteFile(
		filename, []byte(`{"name":$name,"age":$1,"gender":$^NotEmpty}`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	checkOK(t, got,
		td.JSON(filename,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// Reading (a file)
	tmpfile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	checkOK(t, got,
		td.JSON(tmpfile,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))
	tmpfile.Close()

	//
	// Escaping $ in strings
	checkOK(t, "$test", td.JSON(`"$$test"`))

	//
	// Errors
	checkError(t, func() {}, td.JSON(`null`),
		expectedError{
			Message: mustBe("json.Marshal failed"),
			Summary: mustContain("json: unsupported type"),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, map[string]string{"zip": "pipo"},
		td.All(td.JSON(`SuperMapOf({"zip":$1})`, "bingo")),
		expectedError{
			Path:    mustBe(`DATA`),
			Message: mustBe("compared (part 1 of 1)"),
			Got: mustBe(`(map[string]string) (len=1) {
 (string) (len=3) "zip": (string) (len=4) "pipo"
}`),
			Expected: mustBe(`JSON(SuperMapOf(map[string]interface {}{
       "zip": "bingo",
     }))`,
			),
			Under: mustContain("under operator All at "),
			Origin: &expectedError{
				Path:     mustBe(`DATA<All#1/1>["zip"]`),
				Message:  mustBe(`values differ`),
				Got:      mustBe(`"pipo"`),
				Expected: mustBe(`"bingo"`),
				Under:    mustContain("under operator SuperMapOf at line 1:0 (pos 0)" + insideOpJSON),
			},
		})

	checkError(t, map[string]string{"zip": "pipo"},
		td.JSON(`SuperMapOf({"zip":$1})`, "bingo"),
		expectedError{
			Path:     mustBe(`DATA["zip"]`),
			Message:  mustBe("values differ"),
			Got:      mustBe(`"pipo"`),
			Expected: mustBe(`"bingo"`),
			Under:    mustContain("under operator SuperMapOf at line 1:0 (pos 0)" + insideOpJSON),
		})

	checkError(t, json.RawMessage(`"pipo:bingo"`),
		td.JSON(`Re(r;^pipo:(\w+);, ["bad"])`),
		expectedError{
			Path:     mustBe(`(DATA =~ ^pipo:(\w+))[0]`),
			Got:      mustBe(`"bingo"`),
			Expected: mustBe(`"bad"`),
			Under:    mustContain("under operator Re at line 1:0 (pos 0)" + insideOpJSON),
		})

	checkError(t, json.RawMessage(`"pipo:bingo"`),
		td.JSON(`Re(r;^pipo:(\w+);, [$1])`, "bad"),
		expectedError{
			Path:     mustBe(`(DATA =~ ^pipo:(\w+))[0]`),
			Got:      mustBe(`"bingo"`),
			Expected: mustBe(`"bad"`),
			Under:    mustContain("under operator Re at line 1:0 (pos 0)" + insideOpJSON),
		})

	checkError(t, json.RawMessage(`"pipo:bingo"`),
		td.JSON(`Re(r;^pipo:(\w+);, [$param])`, td.Tag("param", "bad")),
		expectedError{
			Path:     mustBe(`(DATA =~ ^pipo:(\w+))[0]`),
			Got:      mustBe(`"bingo"`),
			Expected: mustBe(`"bad"`),
			Under:    mustContain("under operator Re at line 1:0 (pos 0)" + insideOpJSON),
		})

	checkError(t, json.RawMessage(`"pipo:bingo"`),
		td.JSON(`Re(r;^pipo:(\w+);, Bag($1))`, "bad"),
		expectedError{
			Path:    mustBe(`(DATA =~ ^pipo:(\w+))`),
			Summary: mustBe(`Missing item: ("bad")` + "\n" + `  Extra item: ("bingo")`),
			Under:   mustContain("under operator Bag at line 1:19 (pos 19)" + insideOpJSON),
		})

	checkError(t, json.RawMessage(`"pipo:bingo"`),
		td.JSON(`Re(r;^pipo:(\w+);, Bag($param))`, td.Tag("param", "bad")),
		expectedError{
			Path:    mustBe(`(DATA =~ ^pipo:(\w+))`),
			Summary: mustBe(`Missing item: ("bad")` + "\n" + `  Extra item: ("bingo")`),
			Under:   mustContain("under operator Bag at line 1:19 (pos 19)" + insideOpJSON),
		})

	//
	// Fatal errors
	checkError(t, "never tested",
		td.JSON("uNkNoWnFiLe.json"),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("JSON file uNkNoWnFiLe.json cannot be read: "),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(42),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: JSON(STRING_JSON|STRING_FILENAME|[]byte|json.RawMessage|io.Reader, ...), but received int as 1st parameter"),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(errReader{}),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("JSON read error: an error occurred"),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`pipo`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustContain("JSON unmarshal error: "),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[$foo]`,
			td.Tag("foo", td.Ignore()),
			td.Tag("foo", td.Ignore())),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`2 params have the same tag "foo"`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, []int{42},
		td.JSON(`[$1]`, func() {}),
		expectedError{
			Message: mustBe("an error occurred while unmarshalling JSON into func()"),
			Path:    mustBe("DATA[0]"),
			Summary: mustBe("json: cannot unmarshal number into Go value of type func()"),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, []int{42},
		td.JSON(`[$foo]`, td.Tag("foo", func() {})),
		expectedError{
			Message: mustBe("an error occurred while unmarshalling JSON into func()"),
			Path:    mustBe("DATA[0]"),
			Summary: mustBe("json: cannot unmarshal number into Go value of type func()"),
			Under:   mustContain(underOpJSON),
		})

	// numeric placeholders
	checkError(t, "never tested",
		td.JSON(`[1, "$123bad"]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder at line 1:5 (pos 5)`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[1, $000]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder "$000", it should start at "$1" at line 1:4 (pos 4)`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[1, $1]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$1", but no params given at line 1:4 (pos 4)`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[1, 2, $3]`, td.Ignore()),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$3", but only one param given at line 1:7 (pos 7)`),
			Under:   mustContain(underOpJSON),
		})

	// $^Operator
	checkError(t, "never tested",
		td.JSON(`[1, $^bad%]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:4 (pos 4)`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[1, "$^bad%"]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:5 (pos 5)`),
			Under:   mustContain(underOpJSON),
		})

	// named placeholders
	checkError(t, "never tested",
		td.JSON(`[
  1,
  "$bad%"
]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: bad placeholder "$bad%" at line 3:3 (pos 10)`),
			Under:   mustContain(underOpJSON),
		})

	checkError(t, "never tested",
		td.JSON(`[1, $unknown]`),
		expectedError{
			Message: mustBe("bad usage of JSON operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: unknown placeholder "$unknown" at line 1:4 (pos 4)`),
			Under:   mustContain(underOpJSON),
		})

	//
	// Stringification
	test.EqualStr(t, td.JSON(`1`).String(),
		`JSON(1)`)

	test.EqualStr(t, td.JSON(`[ 1, 2, 3 ]`).String(),
		`
JSON([
       1,
       2,
       3
     ])`[1:])

	test.EqualStr(t, td.JSON(` null `).String(), `JSON(null)`)

	test.EqualStr(t,
		td.JSON(`[ $1, $name, $2, Nil(), $nil, 26, Between(5, 6), Len(34), Len(Between(5, 6)), 28 ]`,
			td.Between(12, 20),
			"test",
			td.Tag("name", td.Code(func(s string) bool { return len(s) > 0 })),
			td.Tag("nil", nil),
			14,
		).String(),
		`
JSON([
       "$1" /* 12 ≤ got ≤ 20 */,
       "$name" /* Code(func(string) bool) */,
       "test",
       nil,
       null,
       26,
       5.0 ≤ got ≤ 6.0,
       len=34,
       len: 5.0 ≤ got ≤ 6.0,
       28
     ])`[1:])

	test.EqualStr(t,
		td.JSON(`[ $1, $name, $2, $^Nil, $nil ]`,
			td.Between(12, 20),
			"test",
			td.Tag("name", td.Code(func(s string) bool { return len(s) > 0 })),
			td.Tag("nil", nil),
		).String(),
		`
JSON([
       "$1" /* 12 ≤ got ≤ 20 */,
       "$name" /* Code(func(string) bool) */,
       "test",
       nil,
       null
     ])`[1:])

	test.EqualStr(t,
		td.JSON(`{"label": $value, "zip": $^NotZero}`,
			td.Tag("value", td.Bag(
				td.JSON(`{"name": $1,"age":$2,"surname":$3}`,
					td.HasPrefix("Bob"),
					td.Between(12, 24),
					"captain",
				),
				td.JSON(`{"name": $1}`, td.HasPrefix("Alice")),
			)),
		).String(),
		`
JSON({
       "label": "$value" /* Bag(JSON({
                                       "age": "$2" /* 12 ≤ got ≤ 24 */,
                                       "name": "$1" /* HasPrefix("Bob") */,
                                       "surname": "captain"
                                     }),
                                JSON({
                                       "name": "$1" /* HasPrefix("Alice") */
                                     })) */,
       "zip": NotZero()
     })`[1:])

	test.EqualStr(t,
		td.JSON(`
{
  "label": {"name": HasPrefix("Bob"), "age": Between(12,24)},
  "zip":   NotZero()
}`).String(),
		`
JSON({
       "label": {
                  "age": 12.0 ≤ got ≤ 24.0,
                  "name": HasPrefix("Bob")
                },
       "zip": NotZero()
     })`[1:])

	// Erroneous op
	test.EqualStr(t, td.JSON(`[`).String(), "JSON(<ERROR>)")
}

func TestJSONInside(t *testing.T) {
	// Between
	t.Run("Between", func(t *testing.T) {
		got := map[string]int{"val1": 1, "val2": 2}

		checkOK(t, got,
			td.JSON(`{"val1": Between(0, 2), "val2": Between(2, 3, "[[")}`))
		checkOK(t, got,
			td.JSON(`{"val1": Between(0, 2), "val2": Between(2, 3, "BoundsInOut")}`))
		for _, bounds := range []string{"[[", "BoundsInOut"} {
			checkError(t, got,
				td.JSON(`{"val1": Between(0, 2), "val2": Between(1, 2, $1)}`, bounds),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe(`DATA["val2"]`),
					Got:      mustBe("2.0"),
					Expected: mustBe("1.0 ≤ got < 2.0"),
					Under:    mustContain("under operator Between at line 1:32 (pos 32)" + insideOpJSON),
				})
		}

		checkError(t, json.RawMessage(`123`),
			td.JSON(`Between(0, 2)`),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe(`DATA`),
				Got:      mustBe("123.0"),
				Expected: mustBe("0.0 ≤ got ≤ 2.0"),
				Under:    mustContain("under operator Between at line 1:0 (pos 0)" + insideOpJSON),
			})

		checkOK(t, got,
			td.JSON(`{"val1": Between(1, 1), "val2": Between(2, 2, "[]")}`))
		checkOK(t, got,
			td.JSON(`{"val1": Between(1, 1), "val2": Between(2, 2, "BoundsInIn")}`))

		checkOK(t, got,
			td.JSON(`{"val1": Between(0, 1, "]]"), "val2": Between(1, 3, "][")}`))
		checkOK(t, got,
			td.JSON(`{"val1": Between(0, 1, "BoundsOutIn"), "val2": Between(1, 3, "BoundsOutOut")}`))
		for _, bounds := range []string{"]]", "BoundsOutIn"} {
			checkError(t, got,
				td.JSON(`{"val1": 1, "val2": Between(2, 3, $1)}`, bounds),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe(`DATA["val2"]`),
					Got:      mustBe("2.0"),
					Expected: mustBe("2.0 < got ≤ 3.0"),
					Under:    mustContain("under operator Between at line 1:20 (pos 20)" + insideOpJSON),
				})
		}
		for _, bounds := range []string{"][", "BoundsOutOut"} {
			checkError(t, got,
				td.JSON(`{"val1": 1, "val2": Between(2, 3, $1)}`, bounds),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe(`DATA["val2"]`),
					Got:      mustBe("2.0"),
					Expected: mustBe("2.0 < got < 3.0"),
					Under:    mustContain("under operator Between at line 1:20 (pos 20)" + insideOpJSON),
				},
				"using bounds %q", bounds)
			checkError(t, got,
				td.JSON(`{"val1": 1, "val2": Between(1, 2, $1)}`, bounds),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe(`DATA["val2"]`),
					Got:      mustBe("2.0"),
					Expected: mustBe("1.0 < got < 2.0"),
					Under:    mustContain("under operator Between at line 1:20 (pos 20)" + insideOpJSON),
				},
				"using bounds %q", bounds)
		}

		// Bad 3rd parameter
		checkError(t, "never tested",
			td.JSON(`{
  "val2": Between(1, 2, "<>")
}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Between() bad 3rd parameter, use "[]", "[[", "]]" or "][" at line 2:10 (pos 12)`),
				Under:   mustContain(underOpJSON),
			})
		checkError(t, "never tested",
			td.JSON(`{
  "val2": Between(1, 2, 125)
}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Between() bad 3rd parameter, use "[]", "[[", "]]" or "][" at line 2:10 (pos 12)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`{"val2": Between(1)}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Between() requires 2 or 3 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`{"val2": Between(1,2,3,4)}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Between() requires 2 or 3 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})
	})

	// N
	t.Run("N", func(t *testing.T) {
		got := map[string]float32{"val": 2.1}

		checkOK(t, got, td.JSON(`{"val": N(2.1)}`))
		checkOK(t, got, td.JSON(`{"val": N(2, 0.1)}`))

		checkError(t, "never tested",
			td.JSON(`{"val2": N()}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: N() requires 1 or 2 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`{"val2": N(1,2,3)}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: N() requires 1 or 2 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})
	})

	// Re
	t.Run("Re", func(t *testing.T) {
		got := map[string]string{"val": "Foo bar"}

		checkOK(t, got, td.JSON(`{"val": Re("^Foo")}`))
		checkOK(t, got, td.JSON(`{"val": Re("^(\\w+)", ["Foo"])}`))
		checkOK(t, got, td.JSON(`{"val": Re("^(\\w+)", Bag("Foo"))}`))

		checkError(t, "never tested",
			td.JSON(`{"val2": Re()}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Re() requires 1 or 2 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`{"val2": Re(1,2,3)}`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Re() requires 1 or 2 parameters at line 1:9 (pos 9)`),
				Under:   mustContain(underOpJSON),
			})
	})

	// SubMapOf
	t.Run("SubMapOf", func(t *testing.T) {
		got := []map[string]int{{"val1": 1, "val2": 2}}

		checkOK(t, got, td.JSON(`[ SubMapOf({"val1":1, "val2":2, "xxx": "yyy"}) ]`))

		checkError(t, "never tested",
			td.JSON(`[ SubMapOf() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: SubMapOf() requires only one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ SubMapOf(1, 2) ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: SubMapOf() requires only one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})
	})

	// SuperMapOf
	t.Run("SuperMapOf", func(t *testing.T) {
		got := []map[string]int{{"val1": 1, "val2": 2}}

		checkOK(t, got, td.JSON(`[ SuperMapOf({"val1":1}) ]`))

		checkError(t, "never tested",
			td.JSON(`[ SuperMapOf() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: SuperMapOf() requires only one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ SuperMapOf(1, 2) ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: SuperMapOf() requires only one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})
	})

	// errors
	t.Run("Errors", func(t *testing.T) {
		checkError(t, "never tested",
			td.JSON(`[ UnknownOp() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: unknown operator UnknownOp() at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ Catch() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Catch() is not usable in JSON() at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ JSON() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: JSON() is not usable in JSON(), use literal JSON instead at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ All() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: All() requires at least one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ Empty(12) ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: Empty() requires no parameters at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ HasPrefix() ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: HasPrefix() requires only one parameter at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ JSONPointer(1, 2, 3) ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: JSONPointer() requires 2 parameters at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		checkError(t, "never tested",
			td.JSON(`[ JSONPointer(1, 2) ]`),
			expectedError{
				Message: mustBe("bad usage of JSON operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`JSON unmarshal error: JSONPointer() bad #1 parameter type: string required but float64 received at line 1:2 (pos 2)`),
				Under:   mustContain(underOpJSON),
			})

		// This one is not caught by JSON, but by Re itself, as the number
		// of parameters is correct
		checkError(t, json.RawMessage(`"never tested"`),
			td.JSON(`Re(1)`),
			expectedError{
				Message: mustBe("bad usage of Re operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe(`usage: Re(STRING|*regexp.Regexp[, NON_NIL_CAPTURE]), but received float64 as 1st parameter`),
				Under:   mustContain("under operator Re at line 1:0 (pos 0)" + insideOpJSON),
			})
	})
}

func TestJSONTypeBehind(t *testing.T) {
	equalTypes(t, td.JSON(`false`), true)
	equalTypes(t, td.JSON(`"foo"`), "")
	equalTypes(t, td.JSON(`42`), float64(0))
	equalTypes(t, td.JSON(`[1,2,3]`), ([]any)(nil))
	equalTypes(t, td.JSON(`{"a":12}`), (map[string]any)(nil))

	// operator at the root → delegate it TypeBehind() call
	equalTypes(t, td.JSON(`$1`, td.SuperMapOf(map[string]any{"x": 1}, nil)), (map[string]any)(nil))
	equalTypes(t, td.JSON(`SuperMapOf({"x":1})`), (map[string]any)(nil))

	equalTypes(t, td.JSON(`$1`, 123), 42)

	nullType := td.JSON(`null`).TypeBehind()
	if nullType != reflect.TypeOf((*any)(nil)).Elem() {
		t.Errorf("Failed test: got %s intead of interface {}", nullType)
	}

	// Erroneous op
	equalTypes(t, td.JSON(`[`), nil)
}

func TestSubJSONOf(t *testing.T) {
	type MyStruct struct {
		Name   string `json:"name"`
		Age    uint   `json:"age"`
		Gender string `json:"gender"`
	}

	//
	// struct
	//
	got := MyStruct{Name: "Bob", Age: 42, Gender: "male"}

	// No placeholder
	checkOK(t, got,
		td.SubJSONOf(`
{
  "name":    "Bob",
  "age":     42,
  "gender":  "male",
  "details": {  // ← we don't want to test this field
    "city": "Test City",
    "zip":  666
  }
}`))

	// Numeric placeholders
	checkOK(t, got,
		td.SubJSONOf(`{"name":"$1","age":$2,"gender":$3,"details":{}}`,
			"Bob", 42, "male")) // raw values

	checkOK(t, got,
		td.SubJSONOf(`{"name":"$1","age":$2,"gender":$3,"details":{}}`,
			td.Re(`^Bob`),
			td.Between(40, 45),
			td.NotEmpty()))

	// Same using Flatten
	checkOK(t, got,
		td.SubJSONOf(`{"name":"$1","age":$2,"gender":$3,"details":{}}`,
			td.Re(`^Bob`),
			td.Flatten([]td.TestDeep{td.Between(40, 45), td.NotEmpty()}),
		))

	// Tag placeholders
	checkOK(t, got,
		td.SubJSONOf(
			`{"name":"$name","age":$age,"gender":"$gender","details":{}}`,
			td.Tag("name", td.Re(`^Bob`)),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.NotEmpty())))

	// Mixed placeholders + operator
	for _, op := range []string{
		"NotEmpty",
		"NotEmpty()",
		"$^NotEmpty",
		"$^NotEmpty()",
		`"$^NotEmpty"`,
		`"$^NotEmpty()"`,
		`r<$^NotEmpty>`,
		`r<$^NotEmpty()>`,
	} {
		checkOK(t, got,
			td.SubJSONOf(
				`{"name":"$name","age":$1,"gender":`+op+`,"details":{}}`,
				td.Tag("age", td.Between(40, 45)),
				td.Tag("name", td.Re(`^Bob`))),
			"using operator %s", op)
	}

	//
	// Errors
	checkError(t, func() {}, td.SubJSONOf(`{}`),
		expectedError{
			Message: mustBe("json.Marshal failed"),
			Summary: mustContain("json: unsupported type"),
		})

	for i, n := range []any{
		nil,
		(map[string]any)(nil),
		(map[string]bool)(nil),
		([]int)(nil),
	} {
		checkError(t, n, td.SubJSONOf(`{}`),
			expectedError{
				Message:  mustBe("values differ"),
				Got:      mustBe("null"),
				Expected: mustBe("non-null"),
			},
			"nil test #%d", i)
	}

	//
	// Fatal errors
	checkError(t, "never tested",
		td.SubJSONOf(`[1, "$123bad"]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder at line 1:5 (pos 5)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf(`[1, $000]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder "$000", it should start at "$1" at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf(`[1, $1]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$1", but no params given at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf(`[1, 2, $3]`, td.Ignore()),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$3", but only one param given at line 1:7 (pos 7)`),
		})

	// $^Operator
	checkError(t, "never tested",
		td.SubJSONOf(`[1, $^bad%]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf(`[1, "$^bad%"]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:5 (pos 5)`),
		})

	// named placeholders
	checkError(t, "never tested",
		td.SubJSONOf(`[1, "$bad%"]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: bad placeholder "$bad%" at line 1:5 (pos 5)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf(`[1, $unknown]`),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: unknown placeholder "$unknown" at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SubJSONOf("null"),
		expectedError{
			Message: mustBe("bad usage of SubJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("SubJSONOf() only accepts JSON objects {…}"),
		})

	//
	// Stringification
	test.EqualStr(t, td.SubJSONOf(`{}`).String(), `SubJSONOf({})`)

	test.EqualStr(t, td.SubJSONOf(`{"foo":1, "bar":2}`).String(),
		`
SubJSONOf({
            "bar": 2,
            "foo": 1
          })`[1:])

	test.EqualStr(t,
		td.SubJSONOf(`{"label": $value, "zip": $^NotZero}`,
			td.Tag("value", td.Bag(
				td.SubJSONOf(`{"name": $1,"age":$2}`,
					td.HasPrefix("Bob"),
					td.Between(12, 24),
				),
				td.SubJSONOf(`{"name": $1}`, td.HasPrefix("Alice")),
			)),
		).String(),
		`
SubJSONOf({
            "label": "$value" /* Bag(SubJSONOf({
                                                 "age": "$2" /* 12 ≤ got ≤ 24 */,
                                                 "name": "$1" /* HasPrefix("Bob") */
                                               }),
                                     SubJSONOf({
                                                 "name": "$1" /* HasPrefix("Alice") */
                                               })) */,
            "zip": NotZero()
          })`[1:])

	// Erroneous op
	test.EqualStr(t, td.SubJSONOf(`123`).String(), "SubJSONOf(<ERROR>)")
}

func TestSubJSONOfTypeBehind(t *testing.T) {
	equalTypes(t, td.SubJSONOf(`{"a":12}`), (map[string]any)(nil))

	// Erroneous op
	equalTypes(t, td.SubJSONOf(`123`), nil)
}

func TestSuperJSONOf(t *testing.T) {
	type MyStruct struct {
		Name    string `json:"name"`
		Age     uint   `json:"age"`
		Gender  string `json:"gender"`
		Details string `json:"details"`
	}

	//
	// struct
	//
	got := MyStruct{Name: "Bob", Age: 42, Gender: "male", Details: "Nice"}

	// No placeholder
	checkOK(t, got, td.SuperJSONOf(`{"name": "Bob"}`))

	// Numeric placeholders
	checkOK(t, got,
		td.SuperJSONOf(`{"name":"$1","age":$2}`,
			"Bob", 42)) // raw values

	checkOK(t, got,
		td.SuperJSONOf(`{"name":"$1","age":$2}`,
			td.Re(`^Bob`),
			td.Between(40, 45)))

	// Same using Flatten
	checkOK(t, got,
		td.SuperJSONOf(`{"name":"$1","age":$2}`,
			td.Flatten([]td.TestDeep{td.Re(`^Bob`), td.Between(40, 45)}),
		))

	// Tag placeholders
	checkOK(t, got,
		td.SuperJSONOf(`{"name":"$name","gender":"$gender"}`,
			td.Tag("name", td.Re(`^Bob`)),
			td.Tag("gender", td.NotEmpty())))

	// Mixed placeholders + operator
	for _, op := range []string{
		"NotEmpty",
		"NotEmpty()",
		"$^NotEmpty",
		"$^NotEmpty()",
		`"$^NotEmpty"`,
		`"$^NotEmpty()"`,
		`r<$^NotEmpty>`,
		`r<$^NotEmpty()>`,
	} {
		checkOK(t, got,
			td.SuperJSONOf(
				`{"name":"$name","age":$1,"gender":`+op+`}`,
				td.Tag("age", td.Between(40, 45)),
				td.Tag("name", td.Re(`^Bob`))),
			"using operator %s", op)
	}

	// …with comments…
	checkOK(t, got,
		td.SuperJSONOf(`
// This should be the JSON representation of MyStruct struct
{
  // A person:
  "name":   "$name",   // The name of this person
  "age":    $1,        /* The age of this person:
                          - placeholder unquoted, but could be without
                            any change
                          - to demonstrate a multi-lines comment */
  "gender": $^NotEmpty // Shortcut to operator NotEmpty
}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// Errors
	checkError(t, func() {}, td.SuperJSONOf(`{}`),
		expectedError{
			Message: mustBe("json.Marshal failed"),
			Summary: mustContain("json: unsupported type"),
		})

	for i, n := range []any{
		nil,
		(map[string]any)(nil),
		(map[string]bool)(nil),
		([]int)(nil),
	} {
		checkError(t, n, td.SuperJSONOf(`{}`),
			expectedError{
				Message:  mustBe("values differ"),
				Got:      mustBe("null"),
				Expected: mustBe("non-null"),
			},
			"nil test #%d", i)
	}

	//
	// Fatal errors
	checkError(t, "never tested",
		td.SuperJSONOf(`[1, "$123bad"]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder at line 1:5 (pos 5)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf(`[1, $000]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: invalid numeric placeholder "$000", it should start at "$1" at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf(`[1, $1]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$1", but no params given at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf(`[1, 2, $3]`, td.Ignore()),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: numeric placeholder "$3", but only one param given at line 1:7 (pos 7)`),
		})

	// $^Operator
	checkError(t, "never tested",
		td.SuperJSONOf(`[1, $^bad%]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf(`[1, "$^bad%"]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: $^ must be followed by an operator name at line 1:5 (pos 5)`),
		})

	// named placeholders
	checkError(t, "never tested",
		td.SuperJSONOf(`[1, "$bad%"]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: bad placeholder "$bad%" at line 1:5 (pos 5)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf(`[1, $unknown]`),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`JSON unmarshal error: unknown placeholder "$unknown" at line 1:4 (pos 4)`),
		})

	checkError(t, "never tested",
		td.SuperJSONOf("null"),
		expectedError{
			Message: mustBe("bad usage of SuperJSONOf operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("SuperJSONOf() only accepts JSON objects {…}"),
		})

	//
	// Stringification
	test.EqualStr(t, td.SuperJSONOf(`{}`).String(), `SuperJSONOf({})`)

	test.EqualStr(t, td.SuperJSONOf(`{"foo":1, "bar":2}`).String(),
		`
SuperJSONOf({
              "bar": 2,
              "foo": 1
            })`[1:])

	test.EqualStr(t,
		td.SuperJSONOf(`{"label": $value, "zip": $^NotZero}`,
			td.Tag("value", td.Bag(
				td.SuperJSONOf(`{"name": $1,"age":$2}`,
					td.HasPrefix("Bob"),
					td.Between(12, 24),
				),
				td.SuperJSONOf(`{"name": $1}`, td.HasPrefix("Alice")),
			)),
		).String(),
		`
SuperJSONOf({
              "label": "$value" /* Bag(SuperJSONOf({
                                                     "age": "$2" /* 12 ≤ got ≤ 24 */,
                                                     "name": "$1" /* HasPrefix("Bob") */
                                                   }),
                                       SuperJSONOf({
                                                     "name": "$1" /* HasPrefix("Alice") */
                                                   })) */,
              "zip": NotZero()
            })`[1:])

	// Erroneous op
	test.EqualStr(t, td.SuperJSONOf(`123`).String(), "SuperJSONOf(<ERROR>)")
}

func TestSuperJSONOfTypeBehind(t *testing.T) {
	equalTypes(t, td.SuperJSONOf(`{"a":12}`), (map[string]any)(nil))

	// Erroneous op
	equalTypes(t, td.SuperJSONOf(`123`), nil)
}
