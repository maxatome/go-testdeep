// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"io/ioutil"
	"testing"

	td "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
)

func TestNewRequest(tt *testing.T) {
	t := td.NewT(tt)

	t.RunT("NewRequest", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip", "Test")

		t.String(req.Header.Get("Foo"), "Bar")
		t.String(req.Header.Get("Zip"), "Test")
	})

	t.RunT("NewRequest last header value less", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip")

		t.String(req.Header.Get("Foo"), "Bar")
		t.String(req.Header.Get("Zip"), "")
	})
}

type TestStruct struct {
	Name string `json:"name" xml:"name"`
}

func TestNewJSONRequest(tt *testing.T) {
	t := td.NewT(tt)

	t.RunT("NewJSONRequest", func(t *td.T) {
		req := tdhttp.NewJSONRequest("GET", "/path",
			TestStruct{
				Name: "Bob",
			},
			"Foo", "Bar",
			"Zip", "Test")

		t.String(req.Header.Get("Content-Type"), "application/json")
		t.String(req.Header.Get("Foo"), "Bar")
		t.String(req.Header.Get("Zip"), "Test")

		body, err := ioutil.ReadAll(req.Body)
		if t.CmpNoError(err, "read request body") {
			t.String(string(body), `{"name":"Bob"}`)
		}
	})

	t.RunT("NewJSONRequest panic", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewJSONRequest("GET", "/path", func() {}) },
			td.NotEmpty(),
			"JSON encoding failed")
	})
}

func TestNewXMLRequest(tt *testing.T) {
	t := td.NewT(tt)

	t.RunT("NewXMLRequest", func(t *td.T) {
		req := tdhttp.NewXMLRequest("GET", "/path",
			TestStruct{
				Name: "Bob",
			},
			"Foo", "Bar",
			"Zip", "Test")

		t.String(req.Header.Get("Content-Type"), "application/xml")
		t.String(req.Header.Get("Foo"), "Bar")
		t.String(req.Header.Get("Zip"), "Test")

		body, err := ioutil.ReadAll(req.Body)
		if t.CmpNoError(err, "read request body") {
			t.String(string(body), `<TestStruct><name>Bob</name></TestStruct>`)
		}
	})

	t.RunT("NewXMLRequest panic", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewXMLRequest("GET", "/path", func() {}) },
			td.NotEmpty(),
			"XML encoding failed")
	})
}
