// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestMultipartPart(t *testing.T) {
	assert, require := td.AssertRequire(t)

	check := func(part *tdhttp.MultipartPart, expected string) {
		t.Helper()
		var final bytes.Buffer
		// Read in 2 times to be sure Read() can be called several times
		_, err := io.CopyN(&final, part, 5)
		if assert.CmpNoError(err) {
			_, err := io.Copy(&final, part)
			if assert.CmpNoError(err) {
				assert.Cmp(final.String(), strings.ReplaceAll(expected, "%CR", "\r"))
			}
		}
	}

	// Full empty
	b, err := io.ReadAll(&tdhttp.MultipartPart{})
	assert.CmpNoError(err)
	assert.Len(b, 0)

	// Without name
	part := tdhttp.MultipartPart{
		Content: strings.NewReader("hey!\nyo!"),
	}
	check(&part, `Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// Without body
	part = tdhttp.MultipartPart{
		Name: "nobody",
	}
	check(&part, `Content-Disposition: form-data; name="nobody"%CR
`)

	// Without header
	part = tdhttp.MultipartPart{
		Name:    "pipo",
		Content: strings.NewReader("hey!\nyo!"),
	}
	check(&part, `Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// With header
	part = tdhttp.MultipartPart{
		Name:    "pipo",
		Content: strings.NewReader("hey!\nyo!"),
		Header: http.Header{
			"Pipo":         []string{"bingo"},
			"Content-Type": []string{"text/rococo; charset=utf-8"},
		},
	}
	check(&part, `Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/rococo; charset=utf-8%CR
Pipo: bingo%CR
%CR
hey!
yo!`)

	// Without name & body, but with header
	part = tdhttp.MultipartPart{
		Header: http.Header{
			"Pipo": []string{"bingo"},
		},
	}
	check(&part, `Pipo: bingo%CR
`)

	// io.Reader
	check(tdhttp.NewMultipartPart("io", strings.NewReader("hey!\nyo!")),
		`Content-Disposition: form-data; name="io"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// io.Reader + Content-Type
	check(tdhttp.NewMultipartPart("io", strings.NewReader("hey!\nyo!"), "text/rococo; charset=utf-8"),
		`Content-Disposition: form-data; name="io"%CR
Content-Type: text/rococo; charset=utf-8%CR
%CR
hey!
yo!`)

	// String
	check(tdhttp.NewMultipartPartString("pipo", "hey!\nyo!"),
		`Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// String + Content-Type
	check(tdhttp.NewMultipartPartString("pipo", "hey!\nyo!", "text/rococo; charset=utf-8"),
		`Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/rococo; charset=utf-8%CR
%CR
hey!
yo!`)

	// Bytes
	check(tdhttp.NewMultipartPartBytes("pipo", []byte("hey!\nyo!")),
		`Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// Bytes + Content-Type
	check(tdhttp.NewMultipartPartBytes("pipo", []byte("hey!\nyo!"), "text/rococo; charset=utf-8"),
		`Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/rococo; charset=utf-8%CR
%CR
hey!
yo!`)

	// With file name
	dir, err := os.MkdirTemp("", "multipart")
	require.CmpNoError(err)
	defer os.RemoveAll(dir)
	filePath := filepath.Join(dir, "body.txt")
	require.CmpNoError(os.WriteFile(filePath, []byte("hey!\nyo!"), 0666))

	check(tdhttp.NewMultipartPartFile("pipo", filePath),
		`Content-Disposition: form-data; name="pipo"; filename="body.txt"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!`)

	// With file name + Content-Type
	check(tdhttp.NewMultipartPartFile("pipo", filePath, "text/rococo; charset=utf-8"),
		`Content-Disposition: form-data; name="pipo"; filename="body.txt"%CR
Content-Type: text/rococo; charset=utf-8%CR
%CR
hey!
yo!`)

	// Error during os.Open
	_, err = io.ReadAll(
		tdhttp.NewMultipartPartFile("pipo", filepath.Join(dir, "unknown.xxx")),
	)
	assert.CmpError(err)
}

func TestMultipartBody(t *testing.T) {
	assert, require := td.AssertRequire(t)

	dir, err := os.MkdirTemp("", "multipart")
	require.CmpNoError(err)
	defer os.RemoveAll(dir)
	filePath := filepath.Join(dir, "body.txt")
	require.CmpNoError(os.WriteFile(filePath, []byte("hey!\nyo!"), 0666))

	for _, boundary := range []struct{ in, out string }{
		{in: "", out: "go-testdeep-42"},
		{in: "BoUnDaRy", out: "BoUnDaRy"},
	} {
		multi := tdhttp.MultipartBody{
			Boundary: boundary.in,
			Parts: []*tdhttp.MultipartPart{
				{
					Name:    "pipo",
					Content: strings.NewReader("pipo!\nbingo!"),
				},
				tdhttp.NewMultipartPartFile("file", filePath),
				tdhttp.NewMultipartPartString("string", "zip!\nzap!"),
				tdhttp.NewMultipartPartBytes("bytes", []byte(`{"ola":"hello"}`), "application/json"),
				tdhttp.NewMultipartPart("io", nil),
				tdhttp.NewMultipartPart("", nil),
			},
		}

		expected := `--` + boundary.out + `%CR
Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
pipo!
bingo!%CR
--` + boundary.out + `%CR
Content-Disposition: form-data; name="file"; filename="body.txt"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
hey!
yo!%CR
--` + boundary.out + `%CR
Content-Disposition: form-data; name="string"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
zip!
zap!%CR
--` + boundary.out + `%CR
Content-Disposition: form-data; name="bytes"%CR
Content-Type: application/json%CR
%CR
{"ola":"hello"}%CR
--` + boundary.out + `%CR
Content-Disposition: form-data; name="io"%CR
%CR
--` + boundary.out + `%CR
%CR
--` + boundary.out + `--%CR
`

		var final bytes.Buffer
		// Read in 2 times to be sure Read() can be called several times
		_, err = io.CopyN(&final, &multi, 10)
		if !assert.CmpNoError(err) {
			continue
		}
		_, err := io.Copy(&final, &multi)
		if !assert.CmpNoError(err) {
			continue
		}
		if !assert.Cmp(final.String(), strings.ReplaceAll(expected, "%CR", "\r")) {
			continue
		}

		rd := multipart.NewReader(&final, boundary.out)

		// 0
		part, err := rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "pipo")
			assert.Cmp(part.FileName(), "")
			assert.Smuggle(part, io.ReadAll, td.String("pipo!\nbingo!"))
		}

		// 1
		part, err = rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "file")
			assert.Cmp(part.FileName(), "body.txt")
			assert.Smuggle(part, io.ReadAll, td.String("hey!\nyo!"))
		}

		// 2
		part, err = rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "string")
			assert.Cmp(part.FileName(), "")
			assert.Smuggle(part, io.ReadAll, td.String("zip!\nzap!"))
		}

		// 3
		part, err = rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "bytes")
			assert.Cmp(part.FileName(), "")
			assert.Smuggle(part, io.ReadAll, td.String(`{"ola":"hello"}`))
		}

		// 4
		part, err = rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "io")
			assert.Cmp(part.FileName(), "")
			assert.Smuggle(part, io.ReadAll, td.String(""))
		}

		// 5
		part, err = rd.NextPart()
		if assert.CmpNoError(err) {
			assert.Cmp(part.FormName(), "")
			assert.Cmp(part.FileName(), "")
			assert.Smuggle(part, io.ReadAll, td.String(""))
		}

		// EOF
		_, err = rd.NextPart()
		assert.Cmp(err, io.EOF)
	}

	multi := tdhttp.MultipartBody{}
	td.Cmp(t, multi.ContentType(), `multipart/form-data; boundary="go-testdeep-42"`)
	td.Cmp(t, multi.Boundary, "go-testdeep-42",
		"Boundary field set with default value")
	td.CmpEmpty(t, multi.MediaType, "MediaType field NOT set")

	multi.Boundary = "BoUnDaRy"
	td.Cmp(t, multi.ContentType(), `multipart/form-data; boundary="BoUnDaRy"`)

	multi.MediaType = "multipart/mixed"
	td.Cmp(t, multi.ContentType(), `multipart/mixed; boundary="BoUnDaRy"`)
}
