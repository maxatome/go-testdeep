// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestStringifyPlaceholder(t *testing.T) {
	//
	// Success
	okStringifyPlaceholder := func(j string) string {
		t.Helper()

		pos := strings.IndexByte(j, '$')
		if pos < 0 {
			t.Fatalf("No $ found in `%s`", j)
		}

		// Force the capacity to len(j) to be sure we allocate extra room
		nb, err := stringifyPlaceholder([]byte(j)[0:len(j):len(j)], int64(pos))
		if err != nil {
			t.Fatalf("stringifyPlaceholder(`%s`, %d) failed: %s", j, pos, err)
		}
		return string(nb)
	}

	for _, curTest := range []struct {
		jsonOrig     string
		jsonExpected string
	}{
		// Numeric placeholder
		{jsonOrig: `[$456]`, jsonExpected: `["$456"]`},
		{jsonOrig: `[$456,12]`, jsonExpected: `["$456",12]`},
		{jsonOrig: `{"x":$456}`, jsonExpected: `{"x":"$456"}`},
		{jsonOrig: `[ $456 ]`, jsonExpected: `[ "$456" ]`},
		{jsonOrig: `$456`, jsonExpected: `"$456"`},
		{jsonOrig: "$456\r\n", jsonExpected: `"$456"` + "\r\n"},
		{jsonOrig: "$456\n", jsonExpected: `"$456"` + "\n"},
		{jsonOrig: "$456\t", jsonExpected: `"$456"` + "\t"},
		// Named placeholder
		{jsonOrig: `[$name]`, jsonExpected: `["$name"]`},
		{jsonOrig: `[$name,12]`, jsonExpected: `["$name",12]`},
		{jsonOrig: `{"x":$name}`, jsonExpected: `{"x":"$name"}`},
		{jsonOrig: `[ $name ]`, jsonExpected: `[ "$name" ]`},
		{jsonOrig: `$_name`, jsonExpected: `"$_name"`},
		{jsonOrig: `$foo_bar`, jsonExpected: `"$foo_bar"`},
		{jsonOrig: "$name\r\n", jsonExpected: `"$name"` + "\r\n"},
		{jsonOrig: "$name\n", jsonExpected: `"$name"` + "\n"},
		{jsonOrig: "$name\t", jsonExpected: `"$name"` + "\t"},
	} {
		test.EqualStr(t,
			okStringifyPlaceholder(curTest.jsonOrig),
			curTest.jsonExpected,
		)
	}

	//
	// Errors
	errStringifyPlaceholder := func(j string) string {
		t.Helper()

		pos := strings.IndexByte(j, '$')
		if pos < 0 {
			t.Fatalf("No $ found in `%s`", j)
		}

		_, err := stringifyPlaceholder([]byte(j), int64(pos))
		if err == nil {
			t.Fatalf("stringifyPlaceholder(`%s`, %d) succeeded", j, pos)
		}
		return err.Error()
	}

	for _, curTest := range []struct {
		jsonOrig    string
		errExpected string
	}{
		// Numeric placeholder
		{jsonOrig: `[$456a]`, errExpected: `invalid numeric placeholder at offset 2`},
		// Named placeholder
		{jsonOrig: `[$name%]`, errExpected: `invalid named placeholder at offset 2`},
		// Shortcut
		{jsonOrig: `[$^Op%]`, errExpected: `invalid operator shortcut at offset 2`},
		// Not a placeholder
		{jsonOrig: `[$%]`, errExpected: `invalid placeholder at offset 2`},
		{jsonOrig: `$`, errExpected: `invalid placeholder at offset 1`},
		{jsonOrig: `[12,$`, errExpected: `invalid placeholder at offset 5`},
	} {
		test.EqualStr(t,
			errStringifyPlaceholder(curTest.jsonOrig),
			curTest.errExpected,
		)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var target interface{}

	t.Run("clearComment", func(t *testing.T) {
		origErr := errors.New("orig error")

		from := []byte(`/* comment */`)
		err := clearComment(from, 0, origErr)
		if err != nil {
			t.Errorf("clearComment failed: %s", err)
		}
		if !reflect.DeepEqual(from, []byte(`             `)) {
			t.Errorf(`clearComment failed, unexpected buffer: "%s"`, string(from))
		}

		from = []byte(`foo // bar`)
		err = clearComment(from, 4, origErr)
		if err != nil {
			t.Errorf("clearComment failed: %s", err)
		}
		if !reflect.DeepEqual(from, []byte(`foo       `)) {
			t.Errorf(`clearComment failed, unexpected buffer: "%s"`, string(from))
		}

		from = []byte("foo // bar\nzip")
		err = clearComment(from, 4, origErr)
		if err != nil {
			t.Errorf("clearComment failed: %s", err)
		}
		if !reflect.DeepEqual(from, []byte("foo       \nzip")) {
			t.Errorf(`clearComment failed, unexpected buffer: "%s"`, string(from))
		}

		from = []byte("foo //\nzip")
		err = clearComment(from, 4, origErr)
		if err != nil {
			t.Errorf("clearComment failed: %s", err)
		}
		if !reflect.DeepEqual(from, []byte("foo   \nzip")) {
			t.Errorf(`clearComment failed, unexpected buffer: "%s"`, string(from))
		}

		from = []byte("foo /")
		err = clearComment(from, 4, origErr)
		if err != origErr {
			t.Errorf("got: %s, expected: %s", err, origErr)
		}

		from = []byte("foo /*")
		err = clearComment(from, 4, origErr)
		if err == nil || err.Error() != `unterminated comment at offset 5` {
			t.Errorf("got: %s, expected: %s", err, origErr)
		}
	})

	// First call to initialize jsonErrorMesg variable
	err := UnmarshalJSON([]byte(`{}`), &target)
	if err != nil {
		t.Fatalf("First UnmarshalJSON failed: %s", err)
	}
	if jsonErrPlaceholder == "" || jsonErrPlaceholder == "<NO_JSON_ERROR!>" {
		t.Fatal("json.SyntaxError error not found!")
	}
	t.Logf("OK json.SyntaxError error found: %s", jsonErrPlaceholder)

	// Normal case with several placeholders and operator shorcuts
	err = UnmarshalJSON([]byte(`
/* comment */ { /* comment
   */ "numeric_placeholders" /* comment */: [ $1, $2, $3 ], // comment
  "named_placeholders":   [ $foo, $^bar, /* ← op shortcut */ $zip /* comment */ ]
} // comment`), &target)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %s", err)
	}
	if !reflect.DeepEqual(target, map[string]interface{}{
		"numeric_placeholders": []interface{}{"$1", "$2", "$3"},
		"named_placeholders":   []interface{}{"$foo", "$^bar", "$zip"},
	}) {
		t.Errorf("UnmarshalJSON mismatch: %#+v", target)
	}

	//
	// Error cases
	//
	// Placeholders not allowed in map keys
	err = UnmarshalJSON([]byte(`{$foo: 12}`), &target)
	if err == nil {
		t.Errorf("Placeholders allowed in map keys!")
	}

	// Bad placeholder
	err = UnmarshalJSON([]byte(`{"foo": $8ar}`), &target)
	if err == nil {
		t.Errorf("Bad placeholders not detected!")
	}
}
