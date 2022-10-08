// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp/internal"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func newResponse(body string) *http.Response {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"A": []string{"foo"},
			"B": []string{"bar"},
		},
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
}

func inBQ(s string) string {
	return "`" + s + "`"
}

func TestDumpResponse(t *testing.T) {
	tb := test.NewTestingTB("TestDumpResponse")
	internal.DumpResponse(tb, newResponse("one-line"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
`+inBQ(`HTTP/1.0 200 OK
A: foo
B: bar

one-line`))

	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("multi\r\nlines\r\nand\ttabs héhé"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
`+inBQ(`HTTP/1.0 200 OK
A: foo
B: bar

multi
lines
`+"and\ttabs héhé"))

	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("multi\nlines\nand\ttabs héhé"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
`+inBQ(`HTTP/1.0 200 OK
A: foo
B: bar

multi
lines
`+"and\ttabs héhé"))

	// one \r more in body
	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("multi\r\nline\r"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
"HTTP/1.0 200 OK\r\nA: foo\r\nB: bar\r\n\r\nmulti\r\nline\r"`)

	// BOM
	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("\ufeff"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
"HTTP/1.0 200 OK\r\nA: foo\r\nB: bar\r\n\r\n\ufeff"`)

	// Rune error
	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("\xf4\x9f\xbf\xbf"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
"HTTP/1.0 200 OK\r\nA: foo\r\nB: bar\r\n\r\n\xf4\x9f\xbf\xbf"`)

	// `
	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("he`o"))
	td.Cmp(t, tb.LastMessage(),
		`Received response:
"HTTP/1.0 200 OK\r\nA: foo\r\nB: bar\r\n\r\nhe`+"`"+`o"`)

	// 0x7f
	tb.ResetMessages()
	internal.DumpResponse(tb, newResponse("\x7f"))
	td.Cmp(t, tb.LastMessage(),
		td.Re(`Received response:
"HTTP/1.0 200 OK\\r\\nA: foo\\r\\nB: bar\\r\\n\\r\\n(\\u007f|\\x7f)"`))
}
