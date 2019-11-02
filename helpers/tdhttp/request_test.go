// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"io/ioutil"
	"net/http"
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

	// Get
	t.Cmp(tdhttp.Get("/path", "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "GET",
				Header: http.Header{"Foo": []string{"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Post
	t.Cmp(tdhttp.Post("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{"Foo": []string{"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Put
	t.Cmp(tdhttp.Put("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PUT",
				Header: http.Header{"Foo": []string{"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Patch
	t.Cmp(tdhttp.Patch("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PATCH",
				Header: http.Header{"Foo": []string{"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Delete
	t.Cmp(tdhttp.Delete("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "DELETE",
				Header: http.Header{"Foo": []string{"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))
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

	// Post
	t.Cmp(tdhttp.PostJSON("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/json"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Put
	t.Cmp(tdhttp.PutJSON("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PUT",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/json"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Patch
	t.Cmp(tdhttp.PatchJSON("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PATCH",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/json"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Delete
	t.Cmp(tdhttp.DeleteJSON("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "DELETE",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/json"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))
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

	// Post
	t.Cmp(tdhttp.PostXML("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/xml"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Put
	t.Cmp(tdhttp.PutXML("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PUT",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/xml"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Patch
	t.Cmp(tdhttp.PatchXML("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PATCH",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/xml"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Delete
	t.Cmp(tdhttp.DeleteXML("/path", 42, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "DELETE",
				Header: http.Header{
					"Foo":          []string{"Bar"},
					"Content-Type": []string{"application/xml"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))
}
