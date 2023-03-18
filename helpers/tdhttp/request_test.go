// Copyright (c) 2019-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestBasicAuthHeader(t *testing.T) {
	td.Cmp(t,
		tdhttp.BasicAuthHeader("max", "5ecr3T"),
		http.Header{"Authorization": []string{"Basic bWF4OjVlY3IzVA=="}})
}

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

	t.Run("NewRequest header http.Cookie", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			&http.Cookie{Name: "cook1", Value: "val1"},
			http.Cookie{Name: "cook2", Value: "val2"},
		)

		t.Cmp(req.Header, http.Header{"Cookie": []string{"cook1=val1; cook2=val2"}})
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

	t.Run("NewRequest and query params", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			url.Values{"p1": []string{"a", "b"}},
			url.Values{"p2": []string{"a", "b"}},
			tdhttp.Q{"p1": "c", "p2": []string{"c", "d"}},
			tdhttp.Q{"p1": 123, "p3": true},
		)

		t.Cmp(req.URL.String(), "/path?p1=a&p1=b&p1=c&p1=123&p2=a&p2=b&p2=c&p2=d&p3=true")

		// Query param already set in path
		req = tdhttp.NewRequest("GET", "/path?already=true", nil,
			tdhttp.Q{"p1": 123, "p3": true},
		)
		t.Cmp(req.URL.String(), "/path?already=true&p1=123&p3=true")
	})

	t.Run("NewRequest panics", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", "/path", nil, "H", "V", true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[2])"))

		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", "/path", nil, "H1", true) },
			td.HasPrefix(`header "H1" should have a string value, not a bool (@ headersQueryParams[1])`))

		t.CmpPanic(
			func() { tdhttp.Get("/path", true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Head("/path", true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Options("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Post("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.PostForm("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.PostMultipartFormData("/path", &tdhttp.MultipartBody{}, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Patch("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Put("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Delete("/path", nil, true) },
			td.HasPrefix("headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not bool (@ headersQueryParams[0])"))

		// Bad target
		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", ":/badpath", nil) },
			td.HasPrefix(`target is not a valid path: `))

		// Q error
		t.CmpPanic(
			func() { tdhttp.Get("/", tdhttp.Q{"bad": map[string]bool{}}) },
			td.HasPrefix(`headersQueryParams... tdhttp.Q bad parameter: don't know how to add type map[string]bool (map) to param "bad" (@ headersQueryParams[0])`))
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

	// Options
	t.Cmp(tdhttp.Options("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "OPTIONS",
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

	// PostForm - url.Values
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
					io.ReadAll,
					[]byte("param1=val1&param1=val2&param2=zip"),
				),
			}))

	// PostForm - tdhttp.Q
	t.Cmp(
		tdhttp.PostForm("/path",
			tdhttp.Q{
				"param1": "val1",
				"param2": "val2",
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
					io.ReadAll,
					[]byte("param1=val1&param2=val2"),
				),
			}))

	// PostForm - nil data
	t.Cmp(
		tdhttp.PostForm("/path", nil, "Foo", "Bar"),
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
					io.ReadAll,
					[]byte{},
				),
			}))

	// PostMultipartFormData
	req := tdhttp.PostMultipartFormData("/path",
		&tdhttp.MultipartBody{
			Boundary: "BoUnDaRy",
			Parts: []*tdhttp.MultipartPart{
				tdhttp.NewMultipartPartString("p1", "body1!"),
				tdhttp.NewMultipartPartString("p2", "body2!"),
			},
		},
		"Foo", "Bar")
	t.Cmp(req,
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{
					"Content-Type": []string{`multipart/form-data; boundary="BoUnDaRy"`},
					"Foo":          []string{"Bar"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))
	if t.CmpNoError(req.ParseMultipartForm(10000)) {
		t.Cmp(req.PostFormValue("p1"), "body1!")
		t.Cmp(req.PostFormValue("p2"), "body2!")
	}

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

		body, err := io.ReadAll(req.Body)
		if t.CmpNoError(err, "read request body") {
			t.String(string(body), `{"name":"Bob"}`)
		}
	})

	t.Run("NewJSONRequest panic", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewJSONRequest("GET", "/path", func() {}) },
			td.Contains("json: unsupported type: func()"))

		t.CmpPanic(
			func() { tdhttp.PostJSON("/path", func() {}) },
			td.Contains("json: unsupported type: func()"))

		t.CmpPanic(
			func() { tdhttp.PutJSON("/path", func() {}) },
			td.Contains("json: unsupported type: func()"))

		t.CmpPanic(
			func() { tdhttp.PatchJSON("/path", func() {}) },
			td.Contains("json: unsupported type: func()"))

		t.CmpPanic(
			func() { tdhttp.DeleteJSON("/path", func() {}) },
			td.Contains("json: unsupported type: func()"))

		t.CmpPanic(
			func() { tdhttp.NewJSONRequest("GET", "/path", td.JSONPointer("/a", 0)) },
			td.Contains("JSON encoding failed: json: error calling MarshalJSON for type *td.tdJSONPointer: JSONPointer TestDeep operator cannot be json.Marshal'led"))

		// Common user mistake
		t.CmpPanic(
			func() { tdhttp.NewJSONRequest("GET", "/path", td.JSON(`{}`)) },
			td.Contains(`JSON encoding failed: json: error calling MarshalJSON for type *td.tdJSON: JSON TestDeep operator cannot be json.Marshal'led, use json.RawMessage() instead`))
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

		body, err := io.ReadAll(req.Body)
		if t.CmpNoError(err, "read request body") {
			t.String(string(body), `<TestStruct><name>Bob</name></TestStruct>`)
		}
	})

	t.Run("NewXMLRequest panic", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewXMLRequest("GET", "/path", func() {}) },
			td.Contains("XML encoding failed"))

		t.CmpPanic(
			func() { tdhttp.PostXML("/path", func() {}) },
			td.Contains("XML encoding failed"))

		t.CmpPanic(
			func() { tdhttp.PutXML("/path", func() {}) },
			td.Contains("XML encoding failed"))

		t.CmpPanic(
			func() { tdhttp.PatchXML("/path", func() {}) },
			td.Contains("XML encoding failed"))

		t.CmpPanic(
			func() { tdhttp.DeleteXML("/path", func() {}) },
			td.Contains("XML encoding failed"))
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
