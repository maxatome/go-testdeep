// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMediaType = "multipart/form-data"
	defaultBoundary  = "go-testdeep-42"
)

// MultipartBody is a body of a multipart/form-data HTTP request (by
// default, or any other multipart/… body, see MediaType field) as
// defined in [RFC 2046] to be used as a [io.Reader] body of
// [http.Request] and so compliant with [RFC 2388]. It implements
// [io.Reader] and can only be read once. See [PostMultipartFormData]
// and [TestAPI.PostMultipartFormData] for examples of use.
//
// [RFC 2046]: https://tools.ietf.org/html/rfc2046
// [RFC 2388]: https://tools.ietf.org/html/rfc2388
type MultipartBody struct {
	MediaType string           // type to use instead of default "multipart/form-data"
	Boundary  string           // boundary to use between parts. Automatically initialized when calling ContentType().
	Parts     []*MultipartPart // parts composing this multipart/… body.
	content   io.Reader
}

// Read implements [io.Reader] interface.
func (b *MultipartBody) Read(p []byte) (n int, err error) {
	if b.content == nil {
		if b.Boundary == "" {
			b.Boundary = defaultBoundary
		}

		between := []byte("\r\n--" + b.Boundary + "\r\n")
		first := between[2:]
		end := []byte("\r\n--" + b.Boundary + "--\r\n")

		readers := make([]io.Reader, 0, len(b.Parts)*2+3)

		readers = append(readers, bytes.NewReader(first))
		for i, part := range b.Parts {
			if i > 0 {
				readers = append(readers, bytes.NewReader(between))
			}
			readers = append(readers, part)
		}
		readers = append(readers, bytes.NewReader(end))

		b.content = io.MultiReader(readers...)
	}

	return b.content.Read(p)
}

// ContentType returns the Content-Type header to use. As it contains
// the boundary, it is initialized first if it is still empty. By
// default the media type is multipart/form-data but it can be
// overridden using the MediaType field.
//
//	m.MediaType = "multipart/mixed"
//	ct := m.ContentType()
func (b *MultipartBody) ContentType() string {
	mt := b.MediaType
	if mt == "" {
		mt = defaultMediaType
	}

	if b.Boundary == "" {
		b.Boundary = defaultBoundary
	}

	return mt + `; boundary="` + b.Boundary + `"`
}

// MultipartPart is a part in a [MultipartBody] body. It implements io.Reader
// and can only be read once.
type MultipartPart struct {
	Name       string      // is "name" in Content-Disposition. If empty, Content-Disposition header is omitted.
	Filename   string      // is optional. If set it is "filename" in Content-Disposition.
	Content    io.Reader   // is the body section of the part.
	Header     http.Header // is the header of the part and is optional. It is automatically initialized when needed.
	headerDone bool
}

// NewMultipartPart returns a new [MultipartPart] based on body
// content. If body is nil, it means there is no body at all.
func NewMultipartPart(name string, body io.Reader, contentType ...string) *MultipartPart {
	p := MultipartPart{
		Name:    name,
		Content: body,
	}
	if len(contentType) > 0 {
		p.Header = http.Header{"Content-Type": contentType[:1]}
	}
	return &p
}

// NewMultipartPartFile returns a new [MultipartPart] based on
// filePath content. If filePath cannot be opened, an error is
// returned on first Read() call.
func NewMultipartPartFile(name string, filePath string, contentType ...string) *MultipartPart {
	p := NewMultipartPart(name, &fileReader{filePath: filePath}, contentType...)
	p.Filename = filepath.Base(filePath)
	return p
}

// NewMultipartPartString returns a new [MultipartPart] based on body content.
func NewMultipartPartString(name string, body string, contentType ...string) *MultipartPart {
	return NewMultipartPart(name, strings.NewReader(body), contentType...)
}

// NewMultipartPartBytes returns a new [MultipartPart] based on body content.
func NewMultipartPartBytes(name string, body []byte, contentType ...string) *MultipartPart {
	return NewMultipartPart(name, bytes.NewReader(body), contentType...)
}

// Read implements [io.Reader] interface.
func (p *MultipartPart) Read(b []byte) (n int, err error) {
	if !p.headerDone {
		// Header not yet computed
		if p.Header == nil {
			p.Header = http.Header{}
		}
		if p.Name != "" && p.Header.Get("Content-Disposition") == "" {
			val := `form-data; name="` + p.Name + `"`
			if p.Filename != "" {
				val += `; filename="` + p.Filename + `"`
			}
			p.Header.Set("Content-Disposition", val)
		}

		readers := make([]io.Reader, 1, 3)
		if p.Content != nil {
			if p.Header.Get("Content-Type") == "" {
				var head bytes.Buffer
				copied, err := io.CopyN(&head, p.Content, 512)
				if err != nil && err != io.EOF {
					return 0, err
				}
				if copied > 0 {
					p.Header.Set("Content-Type", http.DetectContentType(head.Bytes()))
					readers = append(readers, &head, p.Content)
				}
			} else {
				readers = append(readers, p.Content)
			}
		}

		var header bytes.Buffer
		p.Header.Write(&header) //nolint: errcheck
		if len(readers) > 1 {
			header.WriteString("\r\n")
		}
		readers[0] = &header

		p.Content = io.MultiReader(readers...)

		p.headerDone = true
	}

	return p.Content.Read(b)
}

type fileReader struct {
	filePath string
	file     *os.File
	err      error
}

func (f *fileReader) Read(b []byte) (n int, err error) {
	if f.err != nil {
		return 0, f.err
	}
	if f.file == nil {
		file, err := os.Open(f.filePath)
		if err != nil {
			f.err = err
			return 0, err
		}
		f.file = file
	}
	n, err = f.file.Read(b)
	if err != nil { // At EOF, (*os.File).Read() returns 0, io.EOF
		f.err = err
		f.file.Close()
		f.file = nil
	}
	return
}
