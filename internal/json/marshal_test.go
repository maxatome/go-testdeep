// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package json_test

import (
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/maxatome/go-testdeep/internal/json"
	"github.com/maxatome/go-testdeep/internal/test"
)

type marshalTest bool

func (m marshalTest) MarshalJSON() ([]byte, error) {
	if m {
		return []byte("marshal\ntest"), nil
	}
	return nil, errors.New("marshalling error")
}

func TestMarshal(t *testing.T) {
	for i, tst := range []struct {
		in       any
		expected string
	}{
		{
			in:       float64(123),
			expected: "123",
		},
		{
			in:       math.NaN(),
			expected: "NaN",
		},
		{
			in:       math.Inf(1),
			expected: "+Inf",
		},
		{
			in:       1e-7,
			expected: "1e-7",
		},
		{
			in:       1e22,
			expected: "1e+22",
		},
		{
			in:       "foobar",
			expected: `"foobar"`,
		},
		{
			in:       true,
			expected: `true`,
		},
		{
			in:       false,
			expected: `false`,
		},
		{
			in:       nil,
			expected: `null`,
		},
		{
			in:       (map[string]any)(nil),
			expected: `null`,
		},
		{
			in:       map[string]any{},
			expected: `{}`,
		},
		{
			in: map[string]any{"z": float64(123), "a": float64(890)},
			expected: `{
  "a": 890,
  "z": 123
}`,
		},
		{
			in: map[string]any{
				"label": map[string]any{"age": float64(12), "name": "Bob"},
				"zip":   float64(456),
			},
			expected: `{
  "label": {
             "age": 12,
             "name": "Bob"
           },
  "zip": 456
}`,
		},
		{
			in:       ([]any)(nil),
			expected: `null`,
		},
		{
			in:       []any{},
			expected: `[]`,
		},
		{
			in: []any{"a", float64(123)},
			expected: `[
  "a",
  123
]`,
		},
		{
			in:       marshalTest(true),
			expected: "marshal\ntest",
		},
		{
			in: []any{float64(1), marshalTest(true), float64(3)},
			expected: `[
  1,
  marshal
  test,
  3
]`,
		},
	} {
		b, err := json.Marshal(tst.in, 0)
		test.NoError(t, err, "#%d", i)
		test.EqualStr(t, string(b), tst.expected, "#%d", i)
	}

	for i, in := range []any{
		marshalTest(false),
		map[string]any{"z": float64(123), "a": marshalTest(false)},
		[]any{"a", marshalTest(false)},
	} {
		_, err := json.Marshal(in, 0)
		if test.Error(t, err, "#%d", i) {
			test.EqualStr(t, err.Error(), "marshalling error", "#%d", i)
		}
	}

	_, err := json.Marshal(123, 0)
	if test.Error(t, err) {
		test.EqualStr(t, err.Error(), "Cannot marshal int")
	}
}

func TestAppendMarshal(t *testing.T) {
	var buf bytes.Buffer

	buf.WriteString("<<")
	err := json.AppendMarshal(&buf, "foo", 0)
	test.NoError(t, err)
	buf.WriteString(">>")

	test.EqualStr(t, buf.String(), `<<"foo">>`)
}
