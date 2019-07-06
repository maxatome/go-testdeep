// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
)

func addHeaders(req *http.Request, headers []string) *http.Request {
	i := 0
	for ; i < len(headers)-1; i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}
	if i < len(headers) {
		req.Header.Add(headers[len(headers)-1], "")
	}
	return req
}

// NewRequest creates a new HTTP request as
// net/http/httptest.NewRequest does, with the ability to immediately
// add some headers.
//
//   req := NewRequest("POST", "/pdf", body,
//     "Content-type", "application/pdf",
//   )
func NewRequest(method, target string, body io.Reader, headers ...string) *http.Request {
	return addHeaders(httptest.NewRequest(method, target, body), headers)
}

// NewJSONRequest creates a new HTTP request with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Other headers can be added via headers, as in:
//
//   req := NewJSONRequest("POST", "/data", body,
//     "X-Foo", "Foo-value",
//     "X-Zip", "Zip-value",
//   )
func NewJSONRequest(method, target string, body interface{}, headers ...string) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		panic("JSON encoding failed: " + err.Error())
	}

	return addHeaders(NewRequest(method, target, bytes.NewBuffer(b)),
		append(headers[:len(headers):len(headers)],
			"Content-Type", "application/json"))
}

// NewXMLRequest creates a new HTTP request with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Other headers can be added via headers, as in:
//
//   req := NewXMLRequest("POST", "/data", body,
//     "X-Foo", "Foo-value",
//     "X-Zip", "Zip-value",
//   )
func NewXMLRequest(method, target string, body interface{}, headers ...string) *http.Request {
	b, err := xml.Marshal(body)
	if err != nil {
		panic("XML encoding failed: " + err.Error())
	}

	return addHeaders(NewRequest(method, target, bytes.NewBuffer(b)),
		append(headers[:len(headers):len(headers)],
			"Content-Type", "application/xml"))
}
