// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestNewRequest(tt *testing.T) {
	t := td.NewT(tt)

	t.Run("NewRequest", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip", "Test")

		t.Cmp(req.Header, http.Header{
			"Foo": []string{"Bar"},
			"Zip": []string{"Test"},
		})
	})

	t.Run("NewRequest last header value less", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip")

		t.Cmp(req.Header, http.Header{
			"Foo": []string{"Bar"},
			"Zip": []string{""},
		})
	})

	t.Run("NewRequest header http.Header", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			http.Header{
				"Foo": []string{"Bar"},
				"Zip": []string{"Test"},
			})

		t.Cmp(req.Header, http.Header{
			"Foo": []string{"Bar"},
			"Zip": []string{"Test"},
		})
	})

	t.Run("NewRequest header flattened", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			td.Flatten([]string{
				"Foo", "Bar",
				"Zip", "Test",
			}),
			td.Flatten(map[string]string{
				"Pipo": "Bingo",
				"Hey":  "Yo",
			}),
		)

		t.Cmp(req.Header, http.Header{
			"Foo":  []string{"Bar"},
			"Zip":  []string{"Test"},
			"Pipo": []string{"Bingo"},
			"Hey":  []string{"Yo"},
		})
	})

	t.Run("NewRequest header combined", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"H1", "V1",
			http.Header{
				"H1": []string{"V2"},
				"H2": []string{"V1", "V2"},
			},
			"H2", "V3",
			td.Flatten([]string{
				"H2", "V4",
				"H3", "V1",
				"H3", "V2",
			}),
			td.Flatten(map[string]string{
				"H2": "V5",
				"H3": "V3",
			}),
		)

		t.Cmp(req.Header, http.Header{
			"H1": []string{"V1", "V2"},
			"H2": []string{"V1", "V2", "V3", "V4", "V5"},
			"H3": []string{"V1", "V2", "V3"},
		})
	})

	t.Run("NewRequest header panic", func(t *td.T) {
		t.CmpPanic(func() { tdhttp.NewRequest("GET", "/path", nil, "H", "V", true) },
			"headers... can only contains string and http.Header, not bool (@ headers[2])")

		t.CmpPanic(func() { tdhttp.NewRequest("GET", "/path", nil, "H1", true) },
			`header "H1" should have a string value, not a bool (@ headers[1])`)
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

	// Head
	t.Cmp(tdhttp.Head("/path", "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "HEAD",
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

	// PostForm
	t.Cmp(
		tdhttp.PostForm("/path",
			url.Values{
				"param1": []string{"val1", "val2"},
				"param2": []string{"zip"},
			},
			"Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{
					"Content-Type": []string{"application/x-www-form-urlencoded"},
					"Foo":          []string{"Bar"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
				"Body": td.Smuggle(
					ioutil.ReadAll,
					[]byte("param1=val1&param1=val2&param2=zip"),
				),
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

	t.Run("NewJSONRequest", func(t *td.T) {
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

	t.Run("NewJSONRequest panic", func(t *td.T) {
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

	t.Run("NewXMLRequest", func(t *td.T) {
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

	t.Run("NewXMLRequest panic", func(t *td.T) {
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
