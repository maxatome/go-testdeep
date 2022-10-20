// Copyright (c) 2020-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/helpers/tdhttp/internal"
	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/td"
)

type failed uint8

const (
	responseFailed failed = 1 << iota
	statusFailed
	headerFailed
	trailerFailed
	cookiesFailed
	bodyFailed
)

// TestAPI allows to test one HTTP API. See [NewTestAPI] function to
// create a new instance and get some examples of use.
type TestAPI struct {
	t       *td.T
	handler http.Handler
	name    string

	sentAt   time.Time
	response *httptest.ResponseRecorder
	failed   failed

	// autoDumpResponse dumps the received response when a test fails.
	autoDumpResponse bool
	responseDumped   bool
}

// NewTestAPI creates a [TestAPI] that can be used to test routes of the
// API behind handler.
//
//	tdhttp.NewTestAPI(t, mux).
//	  Get("/test").
//	  CmpStatus(200).
//	  CmpBody("OK!")
//
// Several routes can be tested with the same instance as in:
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/test").
//	  CmpStatus(200).
//	  CmpBody("OK!")
//
//	ta.Get("/ping").
//	  CmpStatus(200).
//	  CmpBody("pong")
//
// Note that tb can be a [*testing.T] as well as a [*td.T].
func NewTestAPI(tb testing.TB, handler http.Handler) *TestAPI {
	return &TestAPI{
		t:       td.NewT(tb),
		handler: handler,
	}
}

// With creates a new [*TestAPI] instance copied from t, but resetting
// the [testing.TB] instance the tests are based on to tb. The
// returned instance is independent from t, sharing only the same
// handler.
//
// It is typically used when the [TestAPI] instance is “reused” in
// sub-tests, as in:
//
//	func TestMyAPI(t *testing.T) {
//	  ta := tdhttp.NewTestAPI(t, MyAPIHandler())
//
//	  ta.Get("/test").CmpStatus(200)
//
//	  t.Run("errors", func (t *testing.T) {
//	    ta := ta.With(t)
//
//	    ta.Get("/test?bad=1").CmpStatus(400)
//	    ta.Get("/test?bad=buzz").CmpStatus(400)
//	  }
//
//	  ta.Get("/next").CmpStatus(200)
//	}
//
// Note that tb can be a [*testing.T] as well as a [*td.T].
//
// See [TestAPI.Run] for another way to handle subtests.
func (ta *TestAPI) With(tb testing.TB) *TestAPI {
	return &TestAPI{
		t:                td.NewT(tb),
		handler:          ta.handler,
		autoDumpResponse: ta.autoDumpResponse,
	}
}

// T returns the internal instance of [*td.T].
func (ta *TestAPI) T() *td.T {
	return ta.t
}

// Run runs f as a subtest of t called name.
func (ta *TestAPI) Run(name string, f func(ta *TestAPI)) bool {
	return ta.t.Run(name, func(tdt *td.T) {
		f(NewTestAPI(tdt, ta.handler))
	})
}

// AutoDumpResponse allows to dump the HTTP response when the first
// error is encountered after a request.
//
//	ta.AutoDumpResponse()
//	ta.AutoDumpResponse(true)
//
// both enable the dump.
func (ta *TestAPI) AutoDumpResponse(enable ...bool) *TestAPI {
	ta.autoDumpResponse = len(enable) == 0 || enable[0]
	return ta
}

// Name allows to name the series of tests that follow. This name is
// used as a prefix for all following tests, in case of failure to
// qualify each test. If len(args) > 1 and the first item of args is
// a string and contains a '%' rune then [fmt.Fprintf] is used to
// compose the name, else args are passed to [fmt.Fprint].
func (ta *TestAPI) Name(args ...any) *TestAPI {
	ta.name = tdutil.BuildTestName(args...)
	if ta.name != "" {
		ta.name += ": "
	}
	return ta
}

// Request sends a new HTTP request to the tested API. Any Cmp* or
// [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
func (ta *TestAPI) Request(req *http.Request) *TestAPI {
	ta.response = httptest.NewRecorder()

	ta.failed = 0
	ta.sentAt = time.Now().Truncate(0)
	ta.responseDumped = false

	ta.handler.ServeHTTP(ta.response, req)

	return ta
}

func (ta *TestAPI) checkRequestSent() bool {
	ta.t.Helper()

	// If no request has been sent, display a nice error message
	return ta.t.RootName("Request").
		Code(ta.response != nil,
			func(sent bool) error {
				if sent {
					return nil
				}
				return &ctxerr.Error{
					Message: "%% not sent!",
					Summary: ctxerr.NewSummary("A request must be sent before testing status, header, body or full response"),
				}
			},
			ta.name+"request is sent")
}

// Failed returns true if any Cmp* or [TestAPI.NoBody] method failed since last
// request sending.
func (ta *TestAPI) Failed() bool {
	return ta.failed != 0
}

// Get sends a HTTP GET to the tested API. Any Cmp* or [TestAPI.NoBody] methods
// can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Get(target string, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := get(target, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Head sends a HTTP HEAD to the tested API. Any Cmp* or [TestAPI.NoBody] methods
// can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Head(target string, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := head(target, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Options sends a HTTP OPTIONS to the tested API. Any Cmp* or
// [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Options(target string, body io.Reader, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := options(target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Post sends a HTTP POST to the tested API. Any Cmp* or
// [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Post(target string, body io.Reader, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := post(target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PostForm sends a HTTP POST with data's keys and values URL-encoded
// as the request body to the tested API. "Content-Type" header is
// automatically set to "application/x-www-form-urlencoded". Any Cmp*
// or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PostForm(target string, data URLValuesEncoder, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := postForm(target, data, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PostMultipartFormData sends a HTTP POST multipart request, like
// multipart/form-data one for example. See [MultipartBody] type for
// details. "Content-Type" header is automatically set depending on
// data.MediaType (defaults to "multipart/form-data") and
// data.Boundary (defaults to "go-testdeep-42"). Any Cmp* or
// [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
//	ta.PostMultipartFormData("/data",
//	  &tdhttp.MultipartBody{
//	    // "multipart/form-data" by default
//	    Parts: []*tdhttp.MultipartPart{
//	      tdhttp.NewMultipartPartString("type", "Sales"),
//	      tdhttp.NewMultipartPartFile("report", "report.json", "application/json"),
//	    },
//	  },
//	  "X-Foo", "Foo-value",
//	  "X-Zip", "Zip-value",
//	)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PostMultipartFormData(target string, data *MultipartBody, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := postMultipartFormData(target, data, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Put sends a HTTP PUT to the tested API. Any Cmp* or [TestAPI.NoBody] methods
// can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Put(target string, body io.Reader, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := put(target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Patch sends a HTTP PATCH to the tested API. Any Cmp* or [TestAPI.NoBody] methods
// can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Patch(target string, body io.Reader, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := patch(target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// Delete sends a HTTP DELETE to the tested API. Any Cmp* or [TestAPI.NoBody] methods
// can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) Delete(target string, body io.Reader, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := del(target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// NewJSONRequest sends a HTTP request with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) NewJSONRequest(method, target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newJSONRequest(method, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PostJSON sends a HTTP POST with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PostJSON(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newJSONRequest(http.MethodPost, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PutJSON sends a HTTP PUT with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PutJSON(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newJSONRequest(http.MethodPut, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PatchJSON sends a HTTP PATCH with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PatchJSON(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newJSONRequest(http.MethodPatch, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// DeleteJSON sends a HTTP DELETE with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) DeleteJSON(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newJSONRequest(http.MethodDelete, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// NewXMLRequest sends a HTTP request with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) NewXMLRequest(method, target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newXMLRequest(method, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PostXML sends a HTTP POST with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PostXML(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newXMLRequest(http.MethodPost, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PutXML sends a HTTP PUT with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PutXML(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newXMLRequest(http.MethodPut, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// PatchXML sends a HTTP PATCH with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) PatchXML(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newXMLRequest(http.MethodPatch, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// DeleteXML sends a HTTP DELETE with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or [TestAPI.NoBody] methods can now be called.
//
// Note that [TestAPI.Failed] status is reset just after this call.
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func (ta *TestAPI) DeleteXML(target string, body any, headersQueryParams ...any) *TestAPI {
	ta.t.Helper()
	req, err := newXMLRequest(http.MethodDelete, target, body, headersQueryParams...)
	if err != nil {
		ta.t.Fatal(err)
	}
	return ta.Request(req)
}

// CmpResponse tests the last request response status against
// expectedResponse. expectedResponse can be a *http.Response or more
// probably a [td.TestDeep] operator.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/test").
//	  CmpResponse(td.Struct(
//	    &http.Response{Status: http.StatusOK}, td.StructFields{
//	      "Header":        td.SuperMapOf(http.Header{"X-Test": {"pipo"}}),
//	      "ContentLength": td.Gt(10),
//	    }))
//
// Some tests can be hard to achieve using operators chaining. In this
// case, the [td.Code] operator can be used to take the full control
// over the extractions and comparisons to do:
//
//	ta.Get("/test").
//	  CmpResponse(td.Code(func (assert, require *td.T, r *http.Response) {
//	    token, err := ParseToken(r.Header.Get("X-Token"))
//	    require.CmpNoError(err)
//
//	    baseURL,err := url.Parse(r.Header.Get("X-Base-URL"))
//	    require.CmpNoError(err)
//
//	    assert.Cmp(baseURL.Query().Get("id"), token.ID)
//	  }))
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpResponse(expectedResponse any) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= responseFailed
		return ta
	}

	if !ta.t.RootName("Response").
		Cmp(ta.response.Result(), expectedResponse, ta.name+"full response should match") {
		ta.failed |= responseFailed

		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
	}

	return ta
}

// CmpStatus tests the last request response status against
// expectedStatus. expectedStatus can be an int to match a fixed HTTP
// status code, or a [td.TestDeep] operator.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/test").
//	  CmpStatus(http.StatusOK)
//
//	ta.PostJSON("/new", map[string]string{"name": "Bob"}).
//	  CmpStatus(td.Between(200, 202))
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpStatus(expectedStatus any) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= statusFailed
		return ta
	}

	if !ta.t.RootName("Response.Status").
		CmpLax(ta.response.Code, expectedStatus, ta.name+"status code should match") {
		ta.failed |= statusFailed

		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
	}

	return ta
}

// CmpHeader tests the last request response header against
// expectedHeader. expectedHeader can be a [http.Header] or a
// [td.TestDeep] operator. Keep in mind that if it is a [http.Header],
// it has to match exactly the response header. Often only the
// presence of a header key is needed:
//
//	ta := tdhttp.NewTestAPI(t, mux).
//	  PostJSON("/new", map[string]string{"name": "Bob"}).
//	  CmdStatus(201).
//	  CmpHeader(td.ContainsKey("X-Custom"))
//
// or some specific key, value pairs:
//
//	ta.CmpHeader(td.SuperMapOf(
//	  http.Header{
//	    "X-Account": []string{"Bob"},
//	  },
//	  td.MapEntries{
//	    "X-Token": td.Bag(td.Re(`^[a-z0-9-]{32}\z`)),
//	  }),
//	)
//
// Note that CmpHeader calls can be chained:
//
//	ta.CmpHeader(td.ContainsKey("X-Account")).
//	  CmpHeader(td.ContainsKey("X-Token"))
//
// instead of doing all tests in one call as [td.All] operator allows it:
//
//	ta.CmpHeader(td.All(
//	  td.ContainsKey("X-Account"),
//	  td.ContainsKey("X-Token"),
//	))
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpHeader(expectedHeader any) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= headerFailed
		return ta
	}

	if !ta.t.RootName("Response.Header").
		CmpLax(ta.response.Result().Header, expectedHeader, ta.name+"header should match") {
		ta.failed |= headerFailed

		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
	}

	return ta
}

// CmpTrailer tests the last request response trailer against
// expectedTrailer. expectedTrailer can be a [http.Header] or a
// [td.TestDeep] operator. Keep in mind that if it is a [http.Header],
// it has to match exactly the response trailer. Often only the
// presence of a trailer key is needed:
//
//	ta := tdhttp.NewTestAPI(t, mux).
//	  PostJSON("/new", map[string]string{"name": "Bob"}).
//	  CmdStatus(201).
//	  CmpTrailer(td.ContainsKey("X-Custom"))
//
// or some specific key, value pairs:
//
//	ta.CmpTrailer(td.SuperMapOf(
//	  http.Header{
//	    "X-Account": []string{"Bob"},
//	  },
//	  td.MapEntries{
//	    "X-Token": td.Re(`^[a-z0-9-]{32}\z`),
//	  }),
//	)
//
// Note that CmpTrailer calls can be chained:
//
//	ta.CmpTrailer(td.ContainsKey("X-Account")).
//	  CmpTrailer(td.ContainsKey("X-Token"))
//
// instead of doing all tests in one call as [td.All] operator allows it:
//
//	ta.CmpTrailer(td.All(
//	  td.ContainsKey("X-Account"),
//	  td.ContainsKey("X-Token"),
//	))
//
// It fails if no request has been sent yet.
//
// Note that until go1.19, it does not handle multiple values in
// a single Trailer header field.
func (ta *TestAPI) CmpTrailer(expectedTrailer any) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= trailerFailed
		return ta
	}

	if !ta.t.RootName("Response.Trailer").
		CmpLax(ta.response.Result().Trailer, expectedTrailer, ta.name+"trailer should match") {
		ta.failed |= trailerFailed

		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
	}

	return ta
}

// CmpCookies tests the last request response cookies against
// expectedCookies. expectedCookies can be a [][*http.Cookie] or a
// [td.TestDeep] operator. Keep in mind that if it is a
// [][*http.Cookie], it has to match exactly the response
// cookies. Often only the presence of a cookie key is needed:
//
//	ta := tdhttp.NewTestAPI(t, mux).
//	  PostJSON("/login", map[string]string{"name": "Bob", "password": "Sponge"}).
//	  CmdStatus(200).
//	  CmpCookies(td.SuperBagOf(td.Struct(&http.Cookie{Name: "cookie_session"}, nil))).
//	  CmpCookies(td.SuperBagOf(td.Smuggle("Name", "cookie_session"))) // shorter
//
// To make tests easier, [http.Cookie.Raw] and [http.Cookie.RawExpires] fields
// of each [*http.Cookie] are zeroed before doing the comparison. So no need
// to fill them when comparing against a simple literal as in:
//
//	ta := tdhttp.NewTestAPI(t, mux).
//	  PostJSON("/login", map[string]string{"name": "Bob", "password": "Sponge"}).
//	  CmdStatus(200).
//	  CmpCookies([]*http.Cookies{
//	    {Name: "cookieName1", Value: "cookieValue1"},
//	    {Name: "cookieName2", Value: "cookieValue2"},
//	  })
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpCookies(expectedCookies any) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= cookiesFailed
		return ta
	}

	// Empty Raw* fields to make comparisons easier
	cookies := ta.response.Result().Cookies()
	for _, c := range cookies {
		c.RawExpires, c.Raw = "", ""
	}

	if !ta.t.RootName("Response.Cookie").
		CmpLax(cookies, expectedCookies, ta.name+"cookies should match") {
		ta.failed |= cookiesFailed

		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
	}

	return ta
}

// findCmpXBodyCaller finds the oldest Cmp* method called.
func findCmpXBodyCaller() string {
	var (
		fn    string
		pc    [20]uintptr
		found bool
	)
	if num := runtime.Callers(5, pc[:]); num > 0 {
		frames := runtime.CallersFrames(pc[:num])
		for {
			frame, more := frames.Next()
			if pos := strings.Index(frame.Function, "tdhttp.(*TestAPI).Cmp"); pos > 0 {
				fn = frame.Function[pos+18:]
				found = true
			} else if found {
				more = false
			}
			if !more {
				break
			}
		}
	}
	return fn
}

func (ta *TestAPI) cmpMarshaledBody(
	acceptEmptyBody bool,
	unmarshal func([]byte, any) error,
	expectedBody any,
) *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= bodyFailed
		return ta
	}

	if !acceptEmptyBody &&
		!ta.t.RootName("Response body").Code(ta.response.Body.Bytes(),
			func(b []byte) error {
				if len(b) > 0 {
					return nil
				}
				return &ctxerr.Error{
					Message: "%% is empty!",
					Summary: ctxerr.NewSummary(
						"Body cannot be empty when using " + findCmpXBodyCaller()),
				}
			},
			ta.name+"body should not be empty") {
		ta.failed |= bodyFailed
		if ta.autoDumpResponse {
			ta.dumpResponse()
		}
		return ta
	}

	tt := ta.t.RootName("Response.Body")

	var bodyType reflect.Type

	// If expectedBody is a TestDeep operator, try to ask it the type
	// behind it. It should work in most cases (typically Struct(),
	// Map() & Slice()).
	var unknownExpectedType, showRawBody bool
	op, ok := expectedBody.(td.TestDeep)
	if ok {
		bodyType = op.TypeBehind()
		if bodyType == nil {
			// As the expected body type cannot be guessed, try to
			// unmarshal in an any
			bodyType = types.Interface
			unknownExpectedType = true

			// Special case for Ignore & NotEmpty operators
			switch op.GetLocation().Func {
			case "Ignore", "NotEmpty":
				showRawBody = (ta.failed & statusFailed) != 0 // Show real body if status failed
			}
		}
	} else {
		bodyType = reflect.TypeOf(expectedBody)
		if bodyType == nil {
			bodyType = types.Interface
		}
	}

	// For unmarshaling below, body must be a pointer
	bodyPtr := reflect.New(bodyType)

	// Try to unmarshal body
	if !tt.RootName("unmarshal(Response.Body)").
		CmpNoError(unmarshal(ta.response.Body.Bytes(), bodyPtr.Interface()), ta.name+"body unmarshaling") {
		// If unmarshal failed, perhaps it's coz the expected body type
		// is unknown?
		if unknownExpectedType {
			tt.Logf("Cannot guess the body expected type as %[1]s TestDeep\n"+
				"operator does not know the type behind it.\n"+
				"You can try All(Isa(EXPECTED_TYPE), %[1]s(…)) to disambiguate…",
				op.GetLocation().Func)
		}
		showRawBody = true // let's show its real body contents
		ta.failed |= bodyFailed
	} else if !tt.Cmp(bodyPtr.Elem().Interface(), expectedBody, ta.name+"body contents is OK") {
		// Try to catch bad body expected type when nothing has been set
		// to non-zero during unmarshaling body. In this case, require
		// to show raw body contents.
		if len(ta.response.Body.Bytes()) > 0 &&
			td.EqDeeply(bodyPtr.Interface(), reflect.New(bodyType).Interface()) {
			showRawBody = true
			tt.Log("Hmm… It seems nothing has been set during unmarshaling…")
		}
		ta.failed |= bodyFailed
	}

	if showRawBody || ((ta.failed&bodyFailed) != 0 && ta.autoDumpResponse) {
		ta.dumpResponse()
	}

	return ta
}

// CmpMarshaledBody tests that the last request response body can be
// unmarshaled using unmarshal function and then, that it matches
// expectedBody. expectedBody can be any type unmarshal function can
// handle, or a [td.TestDeep] operator.
//
// See [TestAPI.CmpJSONBody] and [TestAPI.CmpXMLBody] sources for
// examples of use.
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpMarshaledBody(unmarshal func([]byte, any) error, expectedBody any) *TestAPI {
	ta.t.Helper()
	return ta.cmpMarshaledBody(false, unmarshal, expectedBody)
}

// CmpBody tests the last request response body against
// expectedBody. expectedBody can be a []byte, a string or a
// [td.TestDeep] operator.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/test").
//	  CmpStatus(http.StatusOK).
//	  CmpBody("OK!\n")
//
//	ta.Get("/test").
//	  CmpStatus(http.StatusOK).
//	  CmpBody(td.Contains("OK"))
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpBody(expectedBody any) *TestAPI {
	ta.t.Helper()

	if expectedBody == nil {
		return ta.NoBody()
	}

	return ta.cmpMarshaledBody(
		true, // accept empty body
		func(body []byte, target any) error {
			switch target := target.(type) {
			case *string:
				*target = string(body)
			case *[]byte:
				*target = body
			case *any:
				*target = body
			default:
				// cmpMarshaledBody always calls us with target as a pointer
				return fmt.Errorf(
					"CmpBody only accepts expectedBody be a []byte, a string or a TestDeep operator allowing to match these types, but not type %s",
					reflect.TypeOf(target).Elem())
			}
			return nil
		},
		expectedBody)
}

// CmpJSONBody tests that the last request response body can be
// [json.Unmarshal]'ed and that it matches expectedBody. expectedBody
// can be any type one can [json.Unmarshal] into, or a [td.TestDeep]
// operator.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpJSONBody(Person{
//	    ID:   42,
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
//	ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//	  CmpStatus(http.StatusCreated).
//	  CmpJSONBody(td.SStruct(
//	    Person{
//	      Name: "Bob",
//	      Age:  26,
//	    },
//	    td.StructFields{
//	      "ID": td.NotZero(),
//	    }))
//
// The same with anchoring, and so without [td.SStruct]:
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//	  CmpStatus(http.StatusCreated).
//	  CmpJSONBody(Person{
//	    ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// The same using [td.JSON]:
//
//	ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//	  CmpStatus(http.StatusCreated).
//	  CmpJSONBody(td.JSON(`
//	{
//	  "id":   NotZero(),
//	  "name": "Bob",
//	  "age":  26
//	}`))
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpJSONBody(expectedBody any) *TestAPI {
	ta.t.Helper()
	return ta.CmpMarshaledBody(json.Unmarshal, expectedBody)
}

// CmpXMLBody tests that the last request response body can be
// [xml.Unmarshal]'ed and that it matches expectedBody. expectedBody
// can be any type one can [xml.Unmarshal] into, or a [td.TestDeep]
// operator.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpXMLBody(Person{
//	    ID:   42,
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
//	ta.Get("/person/43").
//	  CmpStatus(http.StatusOK).
//	  CmpXMLBody(td.SStruct(
//	    Person{
//	      Name: "Bob",
//	      Age:  26,
//	    },
//	    td.StructFields{
//	      "ID": td.NotZero(),
//	    }))
//
// The same with anchoring:
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpXMLBody(Person{
//	    ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// It fails if no request has been sent yet.
func (ta *TestAPI) CmpXMLBody(expectedBody any) *TestAPI {
	ta.t.Helper()
	return ta.CmpMarshaledBody(xml.Unmarshal, expectedBody)
}

// NoBody tests that the last request response body is empty.
//
// It fails if no request has been sent yet.
func (ta *TestAPI) NoBody() *TestAPI {
	defer ta.t.AnchorsPersistTemporarily()()

	ta.t.Helper()

	if !ta.checkRequestSent() {
		ta.failed |= bodyFailed
		return ta
	}

	ok := ta.t.RootName("Response.Body").
		Code(len(ta.response.Body.Bytes()) == 0,
			func(empty bool) error {
				if empty {
					return nil
				}
				return &ctxerr.Error{
					Message:  "%% is not empty",
					Got:      types.RawString("not empty"),
					Expected: types.RawString("empty"),
				}
			},
			"body should be empty")
	if !ok {
		ta.failed |= bodyFailed

		// Systematically dump response, no AutoDumpResponse needed
		ta.dumpResponse()
	}

	return ta
}

// Or executes function fn if ta.Failed() is true at the moment it is called.
//
// fn can have several types:
//   - func(body string) or func(t *td.T, body string)
//     → fn is called with response body as a string.
//     If no response has been received yet, body is "";
//   - func(body []byte) or func(t *td.T, body []byte)
//     → fn is called with response body as a []byte.
//     If no response has been received yet, body is nil;
//   - func(t *td.T, resp *httptest.ResponseRecorder)
//     → fn is called with the internal object containing the response.
//     See net/http/httptest for details.
//     If no response has been received yet, resp is nil.
//
// If fn type is not one of these types, it calls ta.T().Fatal().
func (ta *TestAPI) Or(fn any) *TestAPI {
	ta.t.Helper()
	switch fn := fn.(type) {
	case func(string):
		if ta.Failed() {
			var body string
			if ta.response != nil && ta.response.Body != nil {
				body = ta.response.Body.String()
			}
			fn(body)
		}

	case func(*td.T, string):
		if ta.Failed() {
			var body string
			if ta.response != nil && ta.response.Body != nil {
				body = ta.response.Body.String()
			}
			fn(ta.t, body)
		}

	case func([]byte):
		if ta.Failed() {
			var body []byte
			if ta.response != nil && ta.response.Body != nil {
				body = ta.response.Body.Bytes()
			}
			fn(body)
		}

	case func(*td.T, []byte):
		if ta.Failed() {
			var body []byte
			if ta.response != nil && ta.response.Body != nil {
				body = ta.response.Body.Bytes()
			}
			fn(ta.t, body)
		}

	case func(*td.T, *httptest.ResponseRecorder):
		if ta.Failed() {
			fn(ta.t, ta.response)
		}

	default:
		ta.t.Fatal(color.BadUsage(
			"Or(func([*td.T,]string) | func([*td.T,][]byte) | func(*td.T,*httptest.ResponseRecorder))",
			fn, 1, true))
	}

	return ta
}

// OrDumpResponse dumps the response if at least one previous test failed.
//
//	ta := tdhttp.NewTestAPI(t, handler)
//
//	ta.Get("/foo").
//	  CmpStatus(200).
//	  OrDumpResponse(). // if status check failed, dumps the response
//	  CmpBody("bar")    // if it fails, the response is not dumped
//
//	ta.Get("/foo").
//	  CmpStatus(200).
//	  CmpBody("bar").
//	  OrDumpResponse() // dumps the response if status and/or body checks fail
//
// See [TestAPI.AutoDumpResponse] method to automatize this dump.
func (ta *TestAPI) OrDumpResponse() *TestAPI {
	if ta.Failed() {
		ta.dumpResponse()
	}
	return ta
}

func (ta *TestAPI) dumpResponse() {
	if ta.responseDumped {
		return
	}

	ta.t.Helper()
	if ta.response != nil {
		ta.responseDumped = true
		internal.DumpResponse(ta.t, ta.response.Result())
		return
	}

	ta.t.Logf("No response received yet")
}

// Anchor returns a typed value allowing to anchor the [td.TestDeep]
// operator operator in a go classic literal like a struct, slice,
// array or map value.
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpJSONBody(Person{
//	    ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// See [td.T.Anchor] for details.
//
// See [TestAPI.A] method for a shorter synonym of Anchor.
func (ta *TestAPI) Anchor(operator td.TestDeep, model ...any) any {
	return ta.t.Anchor(operator, model...)
}

// A is a synonym for [TestAPI.Anchor]. It returns a typed value allowing to
// anchor the [td.TestDeep] operator in a go classic literal
// like a struct, slice, array or map value.
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpJSONBody(Person{
//	    ID:   ta.A(td.NotZero(), uint64(0)).(uint64),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// See [td.T.Anchor] for details.
func (ta *TestAPI) A(operator td.TestDeep, model ...any) any {
	return ta.Anchor(operator, model...)
}

// SentAt returns the time just before the last request is handled. It
// can be used to check the time a route sets and returns, as in:
//
//	ta.PostJSON("/person/42", Person{Name: "Bob", Age: 23}).
//	  CmpStatus(http.StatusCreated).
//	  CmpJSONBody(Person{
//	    ID:        ta.A(td.NotZero(), uint64(0)).(uint64),
//	    Name:      "Bob",
//	    Age:       23,
//	    CreatedAt: ta.A(td.Between(ta.SentAt(), time.Now())).(time.Time),
//	  })
//
// checks that CreatedAt field is included between the time when the
// request has been sent, and the time when the comparison occurs.
func (ta *TestAPI) SentAt() time.Time {
	return ta.sentAt
}
