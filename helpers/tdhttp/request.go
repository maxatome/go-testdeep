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
	"net/http"
	"net/http/httptest"
)

// NewRequest is an alias on net/http/httptest.NewRequest for
// convenience purpose.
var NewRequest = httptest.NewRequest

// NewJSONRequest creates a new HTTP request with body marshaled to JSON.
func NewJSONRequest(method, target string, body interface{}) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		panic("JSON encoding failed: " + err.Error())
	}

	req := NewRequest(method, target, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")

	return req
}

// NewXMLRequest creates a new HTTP request with body marshaled to XML.
func NewXMLRequest(method, target string, body interface{}) *http.Request {
	b, err := xml.Marshal(body)
	if err != nil {
		panic("XML encoding failed: " + err.Error())
	}

	req := NewRequest(method, target, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/xml")

	return req
}
