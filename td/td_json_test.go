// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type errReader struct{}

// Read implements io.Reader.
func (r errReader) Read(p []byte) (int, error) {
	return 0, errors.New("an error occurred")
}

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

	// Tag placeholders
	checkOK(t, got,
		td.JSON(`{"name":"$name","age":$age,"gender":"$gender"}`,
			td.Tag("name", td.Re(`^Bob`)),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.NotEmpty())))

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		td.JSON(`{"name":"$name","age":$1,"gender":$^NotEmpty}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

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
  "gender": $^NotEmpty // Shortcut to operator NotEmpty
}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// []byte
	checkOK(t, got,
		td.JSON([]byte(`{"name":"$name","age":$1,"gender":"male"}`),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// Loading a file
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	err = ioutil.WriteFile(
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
		})

	//
	// Panics
	test.CheckPanic(t, func() { td.JSON("uNkNoWnFiLe.json") },
		"JSON(): JSON file uNkNoWnFiLe.json cannot be read: ")

	test.CheckPanic(t, func() { td.JSON(42) },
		"usage: JSON(STRING_JSON|STRING_FILENAME|[]byte|io.Reader, ...), but received int as 1st parameter")

	test.CheckPanic(t, func() { td.JSON(errReader{}) },
		"JSON(): JSON read error: an error occurred")

	test.CheckPanic(t, func() { td.JSON(`pipo`) },
		"JSON(): JSON unmarshal error: ")

	test.CheckPanic(t,
		func() {
			td.JSON(`[$foo]`,
				td.Tag("foo", td.Ignore()),
				td.Tag("foo", td.Ignore()))
		},
		`2 params have the same tag "foo"`)

	// numeric placeholders
	test.CheckPanic(t, func() { td.JSON(`[1, "$123bad"]`) },
		`JSON(): JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { td.JSON(`[1, $000]`) },
		`JSON(): JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { td.JSON(`[1, $1]`) },
		`JSON(): JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { td.JSON(`[1, 2, $3]`, td.Ignore()) },
		`JSON(): JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { td.JSON(`[1, "$^bad%"]`) },
		`JSON(): JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { td.JSON(`[1, "$bad%"]`) },
		`JSON(): JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { td.JSON(`[1, $unknown]`) },
		`JSON(): JSON obj[1] contains a unknown placeholder "$unknown"`)

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
		td.JSON(`[ $1, $name, $2, $^Nil ]`,
			td.Between(12, 20),
			"test",
			td.Tag("name", td.Code(
				func(s string) bool { return len(s) > 0 })),
		).String(),
		`
JSON([
       "$1" /* 12 ≤ got ≤ 20 */,
       "$name" /* Code(func(string) bool) */,
       "test",
       "$^Nil"
     ])`[1:])

	test.EqualStr(t,
		td.JSON(`{"label": $value, "zip": $^NotZero}`,
			td.Tag("value", td.Bag(
				td.JSON(`{"name": $1,"age":$2}`,
					td.HasPrefix("Bob"),
					td.Between(12, 24),
				),
				td.JSON(`{"name": $1}`, td.HasPrefix("Alice")),
			)),
		).String(),
		`
JSON({
       "label": "$value" /* Bag(JSON({
                                       "age": "$2" /* 12 ≤ got ≤ 24 */,
                                       "name": "$1" /* HasPrefix("Bob") */
                                     }),
                                JSON({
                                       "name": "$1" /* HasPrefix("Alice") */
                                     })) */,
       "zip": "$^NotZero"
     })`[1:])

	// Improbable edge-case
	test.EqualStr(t,
		td.JSON(`"<testdeep:opOn>"`).String(),
		`JSON("<testdeep:opOn>")`)
}

func TestJSONTypeBehind(t *testing.T) {
	equalTypes(t, td.JSON(`false`), true)
	equalTypes(t, td.JSON(`"foo"`), "")
	equalTypes(t, td.JSON(`42`), float64(0))
	equalTypes(t, td.JSON(`[1,2,3]`), ([]interface{})(nil))
	equalTypes(t, td.JSON(`{"a":12}`), (map[string]interface{})(nil))

	nullType := td.JSON(`null`).TypeBehind()
	if nullType != reflect.TypeOf((*interface{})(nil)).Elem() {
		t.Errorf("Failed test: got %s intead of interface {}", nullType)
	}
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

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		td.SubJSONOf(
			`{"name":"$name","age":$1,"gender":$^NotEmpty,"details":{}}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

	//
	// Errors
	checkError(t, func() {}, td.SubJSONOf(`{}`),
		expectedError{
			Message: mustBe("json.Marshal failed"),
			Summary: mustContain("json: unsupported type"),
		})

	for i, n := range []interface{}{
		nil,
		(map[string]interface{})(nil),
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
	// Panics
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, "$123bad"]`) },
		`SubJSONOf(): JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, $000]`) },
		`SubJSONOf(): JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, $1]`) },
		`SubJSONOf(): JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, 2, $3]`, td.Ignore()) },
		`SubJSONOf(): JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, "$^bad%"]`) },
		`SubJSONOf(): JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, "$bad%"]`) },
		`SubJSONOf(): JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { td.SubJSONOf(`[1, $unknown]`) },
		`SubJSONOf(): JSON obj[1] contains a unknown placeholder "$unknown"`)

	test.CheckPanic(t, func() { td.SubJSONOf("null") },
		"SubJSONOf() only accepts JSON objects {…}")

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
            "zip": "$^NotZero"
          })`[1:])
}

func TestSubJSONOfTypeBehind(t *testing.T) {
	equalTypes(t, td.SubJSONOf(`{"a":12}`), (map[string]interface{})(nil))
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

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		td.SuperJSONOf(
			`{"name":"$name","age":$1,"gender":$^NotEmpty}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.Re(`^Bob`))))

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

	for i, n := range []interface{}{
		nil,
		(map[string]interface{})(nil),
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
	// Panics
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, "$123bad"]`) },
		`SuperJSONOf(): JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, $000]`) },
		`SuperJSONOf(): JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, $1]`) },
		`SuperJSONOf(): JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, 2, $3]`, td.Ignore()) },
		`SuperJSONOf(): JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, "$^bad%"]`) },
		`SuperJSONOf(): JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, "$bad%"]`) },
		`SuperJSONOf(): JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { td.SuperJSONOf(`[1, $unknown]`) },
		`SuperJSONOf(): JSON obj[1] contains a unknown placeholder "$unknown"`)

	test.CheckPanic(t, func() { td.SuperJSONOf("null") },
		"SuperJSONOf() only accepts JSON objects {…}")

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
              "zip": "$^NotZero"
            })`[1:])
}

func TestSuperJSONOfTypeBehind(t *testing.T) {
	equalTypes(t, td.SuperJSONOf(`{"a":12}`), (map[string]interface{})(nil))
}
