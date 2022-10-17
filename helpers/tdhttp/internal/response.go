// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"testing"
	"unicode/utf8"
)

// canBackquote is the same as strconv.CanBackquote but works on
// []byte and accepts '\n' and '\r'.
func canBackquote(b []byte) bool {
	for len(b) > 0 {
		r, wid := utf8.DecodeRune(b)
		b = b[wid:]
		if wid > 1 {
			if r == '\ufeff' {
				return false // BOMs are invisible and should not be quoted.
			}
			continue // All other multibyte runes are correctly encoded and assumed printable.
		}
		if r == utf8.RuneError {
			return false
		}
		if (r < ' ' && r != '\t' && r != '\n' && r != '\r') || r == '`' || r == '\u007F' {
			return false
		}
	}
	return true
}

func replaceCrLf(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

func backquote(b []byte) ([]byte, bool) {
	// if there is as many \r\n as \n, replace all occurrences by \n
	// so we can conveniently print the buffer inside `…`.
	crnl := bytes.Count(b, []byte("\r\n"))
	cr := bytes.Count(b, []byte("\r"))
	if crnl != 0 {
		nl := bytes.Count(b, []byte("\n"))
		if crnl != nl || crnl != cr {
			return nil, false
		}
		return replaceCrLf(b), true
	}

	return b, cr == 0
}

// DumpResponse logs "resp" using Logf method of "t".
//
// It tries to produce a result as readable as possible first using
// backquotes then falling back to double-quotes.
func DumpResponse(t testing.TB, resp *http.Response) {
	t.Helper()

	const label = "Received response:\n"
	b, _ := httputil.DumpResponse(resp, true)
	if canBackquote(b) {
		bodyPos := bytes.Index(b, []byte("\r\n\r\n"))

		if body, ok := backquote(b[bodyPos+4:]); ok {
			headers := replaceCrLf(b[:bodyPos])
			t.Logf(label+"`%s\n\n%s`", headers, body)
			return
		}
	}

	t.Logf(label+"%q", b)
}
