// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/trace"
	"github.com/maxatome/go-testdeep/td"
)

func init() {
	trace.IgnorePackage()
}

// Response is used by Cmp*Response functions to make the HTTP
// response match easier. Each field, can be a [td.TestDeep] operator
// as well as the exact expected value.
type Response struct {
	Status  any // is the expected status (ignored if nil)
	Header  any // is the expected header (ignored if nil)
	Cookies any // is the expected cookies (ignored if nil)
	Body    any // is the expected body (expected to be empty if nil)
}

func cmpMarshaledResponse(tb testing.TB,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	acceptEmptyBody bool,
	unmarshal func([]byte, any) error,
	expectedResp Response,
	args ...any,
) bool {
	tb.Helper()

	if testName := tdutil.BuildTestName(args...); testName != "" {
		tb.Log(testName)
	}

	t := td.NewT(tb)
	defer t.AnchorsPersistTemporarily()()

	ta := NewTestAPI(t, http.HandlerFunc(handler)).Request(req)

	// Check status, nil = ignore
	if expectedResp.Status != nil {
		ta.CmpStatus(expectedResp.Status)
	}

	// Check header, nil = ignore
	if expectedResp.Header != nil {
		ta.CmpHeader(expectedResp.Header)
	}

	// Check cookie, nil = ignore
	if expectedResp.Cookies != nil {
		ta.CmpCookies(expectedResp.Cookies)
	}

	if expectedResp.Body == nil {
		ta.NoBody()
	} else {
		ta.cmpMarshaledBody(acceptEmptyBody, unmarshal, expectedResp.Body)
	}

	return !ta.Failed()
}

// CmpMarshaledResponse is the base function used by some others in
// tdhttp package. req is launched against handler. The response body
// is unmarshaled using unmarshal. The response is then tested against
// expectedResp.
//
// args... are optional and allow to name the test, a t.Log() done
// before starting any test. If len(args) > 1 and the first item of
// args is a string and contains a '%' rune then [fmt.Fprintf] is used
// to compose the name, else args are passed to [fmt.Fprint].
//
// It returns true if the tests succeed, false otherwise.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpMarshaledResponse(t testing.TB,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	unmarshal func([]byte, any) error,
	expectedResp Response,
	args ...any,
) bool {
	t.Helper()
	return cmpMarshaledResponse(t, req, handler, false, unmarshal, expectedResp, args...)
}

// CmpResponse is used to match a []byte or string response body. req
// is launched against handler. If expectedResp.Body is non-nil, the
// response body is converted to []byte or string, depending on the
// expectedResp.Body type. The response is then tested against
// expectedResp.
//
// args... are optional and allow to name the test, a t.Log() done
// before starting any test. If len(args) > 1 and the first item of
// args is a string and contains a '%' rune then [fmt.Fprintf] is used
// to compose the name, else args are passed to [fmt.Fprint].
//
// It returns true if the tests succeed, false otherwise.
//
//	ok := tdhttp.CmpResponse(t,
//	  tdhttp.Get("/test"),
//	  myAPI.ServeHTTP,
//	  Response{
//	    Status: http.StatusOK,
//	    Header: td.ContainsKey("X-Custom-Header"),
//	    Body:   "OK!\n",
//	  },
//	  "/test route")
//
// Response.Status, Response.Header and Response.Body fields can all
// be [td.TestDeep] operators as it is for Response.Header field
// here. Otherwise, Response.Status should be an int, Response.Header
// a [http.Header] and Response.Body a []byte or a string.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpResponse(t testing.TB,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...any) bool {
	t.Helper()
	return cmpMarshaledResponse(t,
		req,
		handler,
		true,
		func(body []byte, target any) error {
			switch t := target.(type) {
			case *string:
				*t = string(body)
			case *[]byte:
				*t = body
			case *any:
				*t = body
			default:
				// cmpMarshaledBody (behind cmpMarshaledResponse) always calls
				// us with target as a pointer
				return fmt.Errorf(
					"CmpResponse only accepts expectedResp.Body be a []byte, a string or a TestDeep operator allowing to match these types, but not type %s",
					reflect.TypeOf(target).Elem())
			}
			return nil
		},
		expectedResp,
		args...)
}

// CmpJSONResponse is used to match a JSON response body. req is
// launched against handler. If expectedResp.Body is non-nil, the
// response body is [json.Unmarshal]'ed. The response is then tested
// against expectedResp.
//
// args... are optional and allow to name the test, a t.Log() done
// before starting any test. If len(args) > 1 and the first item of
// args is a string and contains a '%' rune then [fmt.Fprintf] is used
// to compose the name, else args are passed to [fmt.Fprint].
//
// It returns true if the tests succeed, false otherwise.
//
//	ok := tdhttp.CmpJSONResponse(t,
//	  tdhttp.Get("/person/42"),
//	  myAPI.ServeHTTP,
//	  Response{
//	    Status: http.StatusOK,
//	    Header: td.ContainsKey("X-Custom-Header"),
//	    Body:   Person{
//	      ID:   42,
//	      Name: "Bob",
//	      Age:  26,
//	    },
//	  },
//	  "/person/{id} route")
//
// Response.Status, Response.Header and Response.Body fields can all
// be [td.TestDeep] operators as it is for Response.Header field
// here. Otherwise, Response.Status should be an int, Response.Header
// a [http.Header] and Response.Body any type one can
// [json.Unmarshal] into.
//
// If Response.Status and Response.Header are omitted (or nil), they
// are not tested.
//
// If Response.Body is omitted (or nil), it means the body response has to be
// empty. If you want to ignore the body response, use [td.Ignore]
// explicitly.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpJSONResponse(t testing.TB,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...any,
) bool {
	t.Helper()
	return CmpMarshaledResponse(t,
		req,
		handler,
		json.Unmarshal,
		expectedResp,
		args...)
}

// CmpXMLResponse is used to match an XML response body. req
// is launched against handler. If expectedResp.Body is
// non-nil, the response body is [xml.Unmarshal]'ed. The response is
// then tested against expectedResp.
//
// args... are optional and allow to name the test, a t.Log() done
// before starting any test. If len(args) > 1 and the first item of
// args is a string and contains a '%' rune then [fmt.Fprintf] is used
// to compose the name, else args are passed to [fmt.Fprint].
//
// It returns true if the tests succeed, false otherwise.
//
//	ok := tdhttp.CmpXMLResponse(t,
//	  tdhttp.Get("/person/42"),
//	  myAPI.ServeHTTP,
//	  Response{
//	    Status: http.StatusOK,
//	    Header: td.ContainsKey("X-Custom-Header"),
//	    Body:   Person{
//	      ID:   42,
//	      Name: "Bob",
//	      Age:  26,
//	    },
//	  },
//	  "/person/{id} route")
//
// Response.Status, Response.Header and Response.Body fields can all
// be [td.TestDeep] operators as it is for Response.Header field
// here. Otherwise, Response.Status should be an int, Response.Header
// a [http.Header] and Response.Body any type one can [xml.Unmarshal]
// into.
//
// If Response.Status and Response.Header are omitted (or nil), they
// are not tested.
//
// If Response.Body is omitted (or nil), it means the body response
// has to be empty. If you want to ignore the body response, use
// [td.Ignore] explicitly.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpXMLResponse(t testing.TB,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...any,
) bool {
	t.Helper()
	return CmpMarshaledResponse(t,
		req,
		handler,
		xml.Unmarshal,
		expectedResp,
		args...)
}

// CmpMarshaledResponseFunc returns a function ready to be used with
// [testing.T.Run], calling [CmpMarshaledResponse] behind the scene. As it
// is intended to be used in conjunction with [testing.T.Run] which
// names the sub-test, the test name part (args...) is voluntary
// omitted.
//
//	t.Run("Subtest name", tdhttp.CmpMarshaledResponseFunc(
//	  tdhttp.Get("/text"),
//	  mux.ServeHTTP,
//	  tdhttp.Response{
//	    Status: http.StatusOK,
//	  }))
//
// See [CmpMarshaledResponse] for details.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpMarshaledResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	unmarshal func([]byte, any) error,
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpMarshaledResponse(t, req, handler, unmarshal, expectedResp)
	}
}

// CmpResponseFunc returns a function ready to be used with
// [testing.T.Run], calling [CmpResponse] behind the scene. As it is
// intended to be used in conjunction with [testing.T.Run] which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//	t.Run("Subtest name", tdhttp.CmpResponseFunc(
//	  tdhttp.Get("/text"),
//	  mux.ServeHTTP,
//	  tdhttp.Response{
//	    Status: http.StatusOK,
//	  }))
//
// See [CmpResponse] documentation for details.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpResponse(t, req, handler, expectedResp)
	}
}

// CmpJSONResponseFunc returns a function ready to be used with
// [testing.T.Run], calling [CmpJSONResponse] behind the scene. As it is
// intended to be used in conjunction with [testing.T.Run] which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//	t.Run("Subtest name", tdhttp.CmpJSONResponseFunc(
//	  tdhttp.Get("/json"),
//	  mux.ServeHTTP,
//	  tdhttp.Response{
//	    Status: http.StatusOK,
//	    Body:   JResp{Comment: "expected comment!"},
//	  }))
//
// See [CmpJSONResponse] documentation for details.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpJSONResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpJSONResponse(t, req, handler, expectedResp)
	}
}

// CmpXMLResponseFunc returns a function ready to be used with
// [testing.T.Run], calling [CmpXMLResponse] behind the scene. As it is
// intended to be used in conjunction with [testing.T.Run] which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//	t.Run("Subtest name", tdhttp.CmpXMLResponseFunc(
//	  tdhttp.Get("/xml"),
//	  mux.ServeHTTP,
//	  tdhttp.Response{
//	    Status: http.StatusOK,
//	    Body:   JResp{Comment: "expected comment!"},
//	  }))
//
// See [CmpXMLResponse] documentation for details.
//
// See [TestAPI] type and its methods for more flexible tests.
func CmpXMLResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpXMLResponse(t, req, handler, expectedResp)
	}
}
