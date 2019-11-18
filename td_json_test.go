// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
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
	checkOK(t, nil, testdeep.JSON(`null`))
	checkOK(t, (*int)(nil), testdeep.JSON(`null`))

	//
	// Basic types
	checkOK(t, 123, testdeep.JSON(`  123  `))
	checkOK(t, true, testdeep.JSON(`  true  `))
	checkOK(t, false, testdeep.JSON(`  false  `))
	checkOK(t, "foobar", testdeep.JSON(`  "foobar"  `))

	//
	// struct
	//
	got := MyStruct{Name: "Bob", Age: 42, Gender: "male"}

	// No placeholder
	checkOK(t, got,
		testdeep.JSON(`{"name":"Bob","age":42,"gender":"male"}`))

	// Numeric placeholders
	checkOK(t, got,
		testdeep.JSON(`{"name":"$1","age":$2,"gender":$3}`,
			"Bob", 42, "male")) // raw values

	checkOK(t, got,
		testdeep.JSON(`{"name":"$1","age":$2,"gender":"$3"}`,
			testdeep.Re(`^Bob`),
			testdeep.Between(40, 45),
			testdeep.NotEmpty()))

	// Tag placeholders
	checkOK(t, got,
		testdeep.JSON(`{"name":"$name","age":$age,"gender":"$gender"}`,
			testdeep.Tag("name", testdeep.Re(`^Bob`)),
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("gender", testdeep.NotEmpty())))

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		testdeep.JSON(`{"name":"$name","age":$1,"gender":$^NotEmpty}`,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	// …with comments…
	checkOK(t, got,
		testdeep.JSON(`
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
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// []byte
	checkOK(t, got,
		testdeep.JSON([]byte(`{"name":"$name","age":$1,"gender":"male"}`),
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

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
		testdeep.JSON(filename,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// Reading (a file)
	tmpfile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	checkOK(t, got,
		testdeep.JSON(tmpfile,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))
	tmpfile.Close()

	//
	// Escaping $ in strings
	checkOK(t, "$test", testdeep.JSON(`"$$test"`))

	//
	// Errors
	checkError(t, func() {}, testdeep.JSON(`null`),
		expectedError{
			Message: mustBe("json.Marshal failed"),
			Summary: mustContain("json: unsupported type"),
		})

	//
	// Panics
	test.CheckPanic(t, func() { testdeep.JSON("uNkNoWnFiLe.json") },
		"JSON file uNkNoWnFiLe.json cannot be read: ")

	test.CheckPanic(t, func() { testdeep.JSON(42) },
		"usage: JSON(STRING_JSON|STRING_FILENAME|[]byte|io.Reader, ...)")

	test.CheckPanic(t, func() { testdeep.JSON(errReader{}) },
		"JSON read error: an error occurred")

	test.CheckPanic(t, func() { testdeep.JSON(`pipo`) },
		"JSON unmarshal error: ")

	test.CheckPanic(t,
		func() {
			testdeep.JSON(`[$foo]`,
				testdeep.Tag("foo", testdeep.Ignore()),
				testdeep.Tag("foo", testdeep.Ignore()))
		},
		`2 params have the same tag "foo"`)

	// numeric placeholders
	test.CheckPanic(t, func() { testdeep.JSON(`[1, "$123bad"]`) },
		`JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { testdeep.JSON(`[1, $000]`) },
		`JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { testdeep.JSON(`[1, $1]`) },
		`JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { testdeep.JSON(`[1, 2, $3]`, testdeep.Ignore()) },
		`JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { testdeep.JSON(`[1, "$^bad%"]`) },
		`JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { testdeep.JSON(`[1, "$bad%"]`) },
		`JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { testdeep.JSON(`[1, $unknown]`) },
		`JSON obj[1] contains a unknown placeholder "$unknown"`)

	//
	// Stringification
	test.EqualStr(t, testdeep.JSON(`1`).String(),
		`JSON(1)`)

	test.EqualStr(t, testdeep.JSON(`[ 1, 2, 3 ]`).String(),
		`
JSON([
       1,
       2,
       3
     ])`[1:])

	test.EqualStr(t, testdeep.JSON(` null `).String(), `JSON(null)`)

	test.EqualStr(t,
		testdeep.JSON(`[ $1, $name, $2, $^Nil ]`,
			testdeep.Between(12, 20),
			"test",
			testdeep.Tag("name", testdeep.Code(
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
		testdeep.JSON(`{"label": $value, "zip": $^NotZero}`,
			testdeep.Tag("value", testdeep.Bag(
				testdeep.JSON(`{"name": $1,"age":$2}`,
					testdeep.HasPrefix("Bob"),
					testdeep.Between(12, 24),
				),
				testdeep.JSON(`{"name": $1}`, testdeep.HasPrefix("Alice")),
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
		testdeep.JSON(`"<testdeep:opOn>"`).String(),
		`JSON("<testdeep:opOn>")`)
}

func TestJSONTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.JSON(`false`), true)
	equalTypes(t, testdeep.JSON(`"foo"`), "")
	equalTypes(t, testdeep.JSON(`42`), float64(0))
	equalTypes(t, testdeep.JSON(`[1,2,3]`), ([]interface{})(nil))
	equalTypes(t, testdeep.JSON(`{"a":12}`), (map[string]interface{})(nil))

	nullType := testdeep.JSON(`null`).TypeBehind()
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
		testdeep.SubJSONOf(`
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
		testdeep.SubJSONOf(`{"name":"$1","age":$2,"gender":$3,"details":{}}`,
			"Bob", 42, "male")) // raw values

	// Tag placeholders
	checkOK(t, got,
		testdeep.SubJSONOf(
			`{"name":"$name","age":$age,"gender":"$gender","details":{}}`,
			testdeep.Tag("name", testdeep.Re(`^Bob`)),
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("gender", testdeep.NotEmpty())))

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		testdeep.SubJSONOf(
			`{"name":"$name","age":$1,"gender":$^NotEmpty,"details":{}}`,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// Errors
	checkError(t, func() {}, testdeep.SubJSONOf(`{}`),
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
		checkError(t, n, testdeep.SubJSONOf(`{}`),
			expectedError{
				Message:  mustBe("values differ"),
				Got:      mustBe("null"),
				Expected: mustBe("non-null"),
			},
			"nil test #%d", i)
	}

	//
	// Panics
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, "$123bad"]`) },
		`JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, $000]`) },
		`JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, $1]`) },
		`JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, 2, $3]`, testdeep.Ignore()) },
		`JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, "$^bad%"]`) },
		`JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, "$bad%"]`) },
		`JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { testdeep.SubJSONOf(`[1, $unknown]`) },
		`JSON obj[1] contains a unknown placeholder "$unknown"`)

	test.CheckPanic(t, func() { testdeep.SubJSONOf("null") },
		"SubJSONOf only accepts JSON objects {…}")

	//
	// Stringification
	test.EqualStr(t, testdeep.SubJSONOf(`{}`).String(), `SubJSONOf({})`)

	test.EqualStr(t, testdeep.SubJSONOf(`{"foo":1, "bar":2}`).String(),
		`
SubJSONOf({
            "bar": 2,
            "foo": 1
          })`[1:])

	test.EqualStr(t,
		testdeep.SubJSONOf(`{"label": $value, "zip": $^NotZero}`,
			testdeep.Tag("value", testdeep.Bag(
				testdeep.SubJSONOf(`{"name": $1,"age":$2}`,
					testdeep.HasPrefix("Bob"),
					testdeep.Between(12, 24),
				),
				testdeep.SubJSONOf(`{"name": $1}`, testdeep.HasPrefix("Alice")),
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
	equalTypes(t, testdeep.SubJSONOf(`{"a":12}`), (map[string]interface{})(nil))
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
	checkOK(t, got, testdeep.SuperJSONOf(`{"name": "Bob"}`))

	// Numeric placeholders
	checkOK(t, got,
		testdeep.SuperJSONOf(`{"name":"$1","age":$2}`,
			"Bob", 42)) // raw values

	// Tag placeholders
	checkOK(t, got,
		testdeep.SuperJSONOf(`{"name":"$name","gender":"$gender"}`,
			testdeep.Tag("name", testdeep.Re(`^Bob`)),
			testdeep.Tag("gender", testdeep.NotEmpty())))

	// Mixed placeholders + operator shortcut
	checkOK(t, got,
		testdeep.SuperJSONOf(
			`{"name":"$name","age":$1,"gender":$^NotEmpty}`,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	// …with comments…
	checkOK(t, got,
		testdeep.SuperJSONOf(`
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
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// Errors
	checkError(t, func() {}, testdeep.SuperJSONOf(`{}`),
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
		checkError(t, n, testdeep.SuperJSONOf(`{}`),
			expectedError{
				Message:  mustBe("values differ"),
				Got:      mustBe("null"),
				Expected: mustBe("non-null"),
			},
			"nil test #%d", i)
	}

	//
	// Panics
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, "$123bad"]`) },
		`JSON obj[1] contains a bad numeric placeholder "$123bad"`)
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, $000]`) },
		`JSON obj[1] contains invalid numeric placeholder "$000", it should start at "$1"`)
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, $1]`) },
		`JSON obj[1] contains numeric placeholder "$1", but only 0 params given`)
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, 2, $3]`, testdeep.Ignore()) },
		`JSON obj[2] contains numeric placeholder "$3", but only 1 params given`)

	// operator shortcut
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, "$^bad%"]`) },
		`JSON obj[1] contains a bad operator shortcut "$^bad%"`)
	// named placeholders
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, "$bad%"]`) },
		`JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { testdeep.SuperJSONOf(`[1, $unknown]`) },
		`JSON obj[1] contains a unknown placeholder "$unknown"`)

	test.CheckPanic(t, func() { testdeep.SuperJSONOf("null") },
		"SuperJSONOf only accepts JSON objects {…}")

	//
	// Stringification
	test.EqualStr(t, testdeep.SuperJSONOf(`{}`).String(), `SuperJSONOf({})`)

	test.EqualStr(t, testdeep.SuperJSONOf(`{"foo":1, "bar":2}`).String(),
		`
SuperJSONOf({
              "bar": 2,
              "foo": 1
            })`[1:])

	test.EqualStr(t,
		testdeep.SuperJSONOf(`{"label": $value, "zip": $^NotZero}`,
			testdeep.Tag("value", testdeep.Bag(
				testdeep.SuperJSONOf(`{"name": $1,"age":$2}`,
					testdeep.HasPrefix("Bob"),
					testdeep.Between(12, 24),
				),
				testdeep.SuperJSONOf(`{"name": $1}`, testdeep.HasPrefix("Alice")),
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
	equalTypes(t, testdeep.SuperJSONOf(`{"a":12}`), (map[string]interface{})(nil))
}
