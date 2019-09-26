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
		Name string `json:"name"`
		Age  uint   `json:"age"`
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
	// No placeholder
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(`{"name":"Bob","age":42}`))

	// Numeric placeholders
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(`{"name":"$1","age":$2}`, "Bob", 42)) // raw values

	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(`{"name":"$1","age":$2}`,
			testdeep.Re(`^Bob`),
			testdeep.Between(40, 45)))

	// Tag placeholders
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(`{"name":"$name","age":$age}`,
			testdeep.Tag("name", testdeep.Re(`^Bob`)),
			testdeep.Tag("age", testdeep.Between(40, 45))))

	// Mixed placeholders
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(`{"name":"$name","age":$1}`,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// []byte
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON([]byte(`{"name":"$name","age":$1}`),
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
	if err = ioutil.WriteFile(filename, []byte(`{"name":$name,"age":$1}`), 0644); err != nil {
		t.Fatal(err)
	}
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
		testdeep.JSON(filename,
			testdeep.Tag("age", testdeep.Between(40, 45)),
			testdeep.Tag("name", testdeep.Re(`^Bob`))))

	//
	// Reading (a file)
	tmpfile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	checkOK(t, MyStruct{Name: "Bob", Age: 42},
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

	// named placeholders
	test.CheckPanic(t, func() { testdeep.JSON(`[1, "$bad%"]`) },
		`JSON obj[1] contains a bad placeholder "$bad%"`)
	test.CheckPanic(t, func() { testdeep.JSON(`[1, $unknown]`) },
		`JSON obj[1] contains a unknown placeholder "$unknown"`)

	//
	// Stringification
	test.EqualStr(t, testdeep.JSON(`[ 1, 2, 3 ]`).String(), `JSON("[1,2,3]")`)
	test.EqualStr(t, testdeep.JSON(` null `).String(), `JSON("null")`)
	test.EqualStr(t,
		testdeep.JSON(`[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17]`).String(),
		`JSON("[1,2,3,4,5,6,7,8,9,10,11,12…")`)
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
