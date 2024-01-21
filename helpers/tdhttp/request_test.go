// Copyright (c) 2019-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"errors"
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
		http.Header{"Authorization": {"Basic bWF4OjVlY3IzVA=="}})
}

func TestNewRequest(tt *testing.T) {
	t := td.NewT(tt)

	t.Run("NewRequest", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip", "Test")

		t.Cmp(req.Header, http.Header{
			"Foo": {"Bar"},
			"Zip": {"Test"},
		})
	})

	t.Run("NewRequest last header value less", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"Foo", "Bar",
			"Zip")

		t.Cmp(req.Header, http.Header{
			"Foo": {"Bar"},
			"Zip": {""},
		})
	})

	t.Run("NewRequest header http.Header", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			http.Header{
				"Foo": {"Bar"},
				"Zip": {"Test"},
			})

		t.Cmp(req.Header, http.Header{
			"Foo": {"Bar"},
			"Zip": {"Test"},
		})
	})

	t.Run("NewRequest header http.Cookie", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			&http.Cookie{Name: "cook1", Value: "val1"},
			http.Cookie{Name: "cook2", Value: "val2"},
			[]*http.Cookie{
				{Name: "cook3", Value: "val3"},
				{Name: "cook4", Value: "val4"},
			},
		)

		t.Cmp(req.Header, http.Header{
			"Cookie": {"cook1=val1; cook2=val2; cook3=val3; cook4=val4"},
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
			"Foo":  {"Bar"},
			"Zip":  {"Test"},
			"Pipo": {"Bingo"},
			"Hey":  {"Yo"},
		})
	})

	t.Run("NewRequest header combined", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			"H1", "V1",
			http.Header{
				"H1": {"V2"},
				"H2": {"V1", "V2"},
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
			"H1": {"V1", "V2"},
			"H2": {"V1", "V2", "V3", "V4", "V5"},
			"H3": {"V1", "V2", "V3"},
		})
	})

	t.Run("NewRequest and query params", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			url.Values{"p1": {"a", "b"}},
			url.Values{"p2": {"a", "b"}},
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

	t.Run("NewRequest host", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "https://pipo.com:123/path", nil)
		td.Cmp(t, req.Host, "pipo.com:123")
		td.Cmp(t, req.URL, td.String("/path"))

		req = tdhttp.NewRequest("GET", "/path", nil, "Host", "pipo.com:456")
		td.Cmp(t, req.Host, "pipo.com:456")
		td.Cmp(t, req.URL, td.String("/path"))

		req = tdhttp.NewRequest(
			"GET", "https://pipo.com:123/path", nil,
			"Host", "bingo.com:456")
		td.Cmp(t, req.Host, "pipo.com:123", "URL wins")
		td.Cmp(t, req.URL, td.String("/path"))
	})

	t.Run("NewRequest and hooks", func(t *td.T) {
		req := tdhttp.NewRequest("GET", "/path", nil,
			func(req *http.Request) error {
				req.Header.Set("H1", "V1")
				return nil
			},
			func(req *http.Request) error {
				req.Header.Set("H2", "V2")
				return nil
			})

		t.Cmp(req.Header, http.Header{
			"H1": {"V1"},
			"H2": {"V2"},
		})
	})

	t.Run("NewRequest panics", func(t *td.T) {
		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", "/path", nil, "H", "V", true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[2])"))

		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", "/path", nil, "H1", true) },
			td.HasPrefix(`header "H1" should have a string value, not a bool (@ newRequestParams[1])`))

		// Hook error
		t.CmpPanic(
			func() {
				tdhttp.NewRequest("GET", "/path", nil, func(*http.Request) error {
					return errors.New("hook error")
				})
			},
			td.String("hook failed: hook error"))

		t.CmpPanic(
			func() { tdhttp.Get("/path", true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Head("/path", true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Options("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Post("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.PostForm("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.PostMultipartFormData("/path", &tdhttp.MultipartBody{}, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Patch("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Put("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		t.CmpPanic(
			func() { tdhttp.Delete("/path", nil, true) },
			td.HasPrefix("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool (@ newRequestParams[0])"))

		// Bad target
		t.CmpPanic(
			func() { tdhttp.NewRequest("GET", ":/badpath", nil) },
			td.HasPrefix(`target is not a valid path: `))

		// Q error
		t.CmpPanic(
			func() { tdhttp.Get("/", tdhttp.Q{"bad": map[string]bool{}}) },
			td.HasPrefix(`newRequestParams... tdhttp.Q bad parameter: don't know how to add type map[string]bool (map) to param "bad" (@ newRequestParams[0])`))
	})

	// Get
	t.Cmp(tdhttp.Get("/path", "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "GET",
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Head
	t.Cmp(tdhttp.Head("/path", "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "HEAD",
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Options
	t.Cmp(tdhttp.Options("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "OPTIONS",
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Post
	t.Cmp(tdhttp.Post("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// PostForm - url.Values
	t.Cmp(
		tdhttp.PostForm("/path",
			url.Values{
				"param1": {"val1", "val2"},
				"param2": {"zip"},
			},
			"Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "POST",
				Header: http.Header{
					"Content-Type": {"application/x-www-form-urlencoded"},
					"Foo":          {"Bar"},
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
					"Content-Type": {"application/x-www-form-urlencoded"},
					"Foo":          {"Bar"},
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
					"Content-Type": {"application/x-www-form-urlencoded"},
					"Foo":          {"Bar"},
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
					"Content-Type": {`multipart/form-data; boundary="BoUnDaRy"`},
					"Foo":          {"Bar"},
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
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Patch
	t.Cmp(tdhttp.Patch("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "PATCH",
				Header: http.Header{"Foo": {"Bar"}},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))

	// Delete
	t.Cmp(tdhttp.Delete("/path", nil, "Foo", "Bar"),
		td.Struct(
			&http.Request{
				Method: "DELETE",
				Header: http.Header{"Foo": {"Bar"}},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/json"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/json"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/json"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/json"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/xml"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/xml"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/xml"},
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
					"Foo":          {"Bar"},
					"Content-Type": {"application/xml"},
				},
			},
			td.StructFields{
				"URL": td.String("/path"),
			}))
}
