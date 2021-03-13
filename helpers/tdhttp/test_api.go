// Copyright (c) 2020, 2021, Maxime Soulé
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
	"net/url"
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

// TestAPI allows to test one HTTP API. See NewTestAPI function to
// create a new instance and get some examples of use.
type TestAPI struct {
	t       *td.T
	handler http.Handler
	name    string

	sentAt       time.Time
	response     *httptest.ResponseRecorder
	statusFailed bool
	headerFailed bool
	bodyFailed   bool

	// autoDumpResponse dumps the received response when a test fails.
	autoDumpResponse bool
	responseDumped   bool
}

// NewTestAPI creates a TestAPI that can be used to test routes of the
// API behind "handler".
//
//   tdhttp.NewTestAPI(t, mux).
//     Get("/test").
//     CmpStatus(200).
//     CmpBody("OK!")
//
// Several routes can be tested with the same instance as in:
//
//   ta := tdhttp.NewTestAPI(t, mux)
//
//   ta.Get("/test").
//     CmpStatus(200).
//     CmpBody("OK!")
//
//   ta.Get("/ping").
//     CmpStatus(200).
//     CmpBody("pong")
//
// Note that "tb" can be a *testing.T as well as a *td.T.
func NewTestAPI(tb testing.TB, handler http.Handler) *TestAPI {
	return &TestAPI{
		t:       td.NewT(tb),
		handler: handler,
	}
}

// With creates a new *TestAPI instance copied from "t", but resetting
// the testing.TB instance the tests are based on to "tb". The
// returned instance is independent from "t", sharing only the same
// handler.
//
// It is typically used when the *TestAPI instance is "reused" in
// sub-tests, as in:
//
//   func TestMyAPI(t *testing.T) {
//     ta := tdhttp.NewTestAPI(t, MyAPIHandler())
//
//     ta.Get("/test").CmpStatus(200)
//
//     t.Run("errors", func (t *testing.T) {
//       ta := ta.With(t)
//
//       ta.Get("/test?bad=1").CmpStatus(400)
//       ta.Get("/test?bad=buzz").CmpStatus(400)
//     }
//
//     ta.Get("/next").CmpStatus(200)
//   }
//
// Note that "tb" can be a *testing.T as well as a *td.T.
//
// See Run method for another way to handle subtests.
func (t *TestAPI) With(tb testing.TB) *TestAPI {
	return &TestAPI{
		t:                td.NewT(tb),
		handler:          t.handler,
		autoDumpResponse: t.autoDumpResponse,
	}
}

// T returns the internal instance of *td.T.
func (t *TestAPI) T() *td.T {
	return t.t
}

// Run runs "f" as a subtest of t called "name".
func (t *TestAPI) Run(name string, f func(t *TestAPI)) bool {
	return t.t.Run(name, func(tdt *td.T) {
		f(NewTestAPI(tdt, t.handler))
	})
}

// AutoDumpResponse allows to dump the HTTP response when the first
// error is encountered after a request.
func (t *TestAPI) AutoDumpResponse() *TestAPI {
	t.autoDumpResponse = true
	return t
}

// Name allows to name the series of tests that follow. This name is
// used as a prefix for all following tests, in case of failure to
// qualify each test. If len(args) > 1 and the first item of "args" is
// a string and contains a '%' rune then fmt.Fprintf is used to
// compose the name, else "args" are passed to fmt.Fprint.
func (t *TestAPI) Name(args ...interface{}) *TestAPI {
	t.name = tdutil.BuildTestName(args...)
	if t.name != "" {
		t.name += ": "
	}
	return t
}

// Request sends a new HTTP request to the tested API. Any Cmp* or
// NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
func (t *TestAPI) Request(req *http.Request) *TestAPI {
	t.response = httptest.NewRecorder()

	t.statusFailed = false
	t.headerFailed = false
	t.bodyFailed = false
	t.sentAt = time.Now().Truncate(0)
	t.responseDumped = false

	t.handler.ServeHTTP(t.response, req)

	return t
}

func (t *TestAPI) checkRequestSent() bool {
	t.t.Helper()

	// If no request has been sent, display a nice error message
	return t.t.RootName("Request").
		Code(t.response != nil,
			func(sent bool) error {
				if sent {
					return nil
				}
				return &ctxerr.Error{
					Message: "%% not sent!",
					Summary: ctxerr.NewSummary("A request must be sent before testing status, header or body"),
				}
			},
			t.name+"request is sent")
}

// Failed returns true if any Cmp* or NoBody method failed since last
// request sending.
func (t *TestAPI) Failed() bool {
	return t.statusFailed || t.headerFailed || t.bodyFailed
}

// Get sends a HTTP GET to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Get(target string, headers ...interface{}) *TestAPI {
	return t.Request(Get(target, headers...))
}

// Head sends a HTTP HEAD to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Head(target string, headers ...interface{}) *TestAPI {
	return t.Request(Head(target, headers...))
}

// Post sends a HTTP POST to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Post(target string, body io.Reader, headers ...interface{}) *TestAPI {
	return t.Request(Post(target, body, headers...))
}

// PostForm sends a HTTP POST with data's keys and values URL-encoded
// as the request body to the tested API.. "Content-Type" header is
// automatically set to "application/x-www-form-urlencoded". Any Cmp*
// or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PostForm(target string, data url.Values, headers ...interface{}) *TestAPI {
	return t.Request(PostForm(target, data, headers...))
}

// Put sends a HTTP PUT to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Put(target string, body io.Reader, headers ...interface{}) *TestAPI {
	return t.Request(Put(target, body, headers...))
}

// Patch sends a HTTP PATCH to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Patch(target string, body io.Reader, headers ...interface{}) *TestAPI {
	return t.Request(Patch(target, body, headers...))
}

// Delete sends a HTTP DELETE to the tested API. Any Cmp* or NoBody methods
// can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) Delete(target string, body io.Reader, headers ...interface{}) *TestAPI {
	return t.Request(Delete(target, body, headers...))
}

// NewJSONRequest sends a HTTP request with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) NewJSONRequest(method, target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(NewJSONRequest(method, target, body, headers...))
}

// PostJSON sends a HTTP POST with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PostJSON(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PostJSON(target, body, headers...))
}

// PutJSON sends a HTTP PUT with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PutJSON(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PutJSON(target, body, headers...))
}

// PatchJSON sends a HTTP PATCH with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PatchJSON(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PatchJSON(target, body, headers...))
}

// DeleteJSON sends a HTTP DELETE with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) DeleteJSON(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(DeleteJSON(target, body, headers...))
}

// NewXMLRequest sends a HTTP request with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) NewXMLRequest(method, target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(NewXMLRequest(method, target, body, headers...))
}

// PostXML sends a HTTP POST with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PostXML(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PostXML(target, body, headers...))
}

// PutXML sends a HTTP PUT with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PutXML(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PutXML(target, body, headers...))
}

// PatchXML sends a HTTP PATCH with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) PatchXML(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(PatchXML(target, body, headers...))
}

// DeleteXML sends a HTTP DELETE with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Any Cmp* or NoBody methods can now be called.
//
// Note that Failed() status is reset just after this call.
//
// See NewRequest for all possible formats accepted in headers.
func (t *TestAPI) DeleteXML(target string, body interface{}, headers ...interface{}) *TestAPI {
	return t.Request(DeleteXML(target, body, headers...))
}

// CmpStatus tests the last request response status against
// expectedStatus. expectedStatus can be an int to match a fixed HTTP
// status code, or a TestDeep operator.
//
//   ta := tdhttp.NewTestAPI(t, mux)
//
//   ta.Get("/test").
//     CmpStatus(http.StatusOK)
//
//   ta.PostJSON("/new", map[string]string{"name": "Bob"}).
//     CmpStatus(td.Between(200, 202))
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpStatus(expectedStatus interface{}) *TestAPI {
	defer t.t.AnchorsPersistTemporarily()()

	t.t.Helper()

	if !t.checkRequestSent() {
		t.statusFailed = true
		return t
	}

	t.statusFailed = !t.t.RootName("Response.Status").
		CmpLax(t.response.Code, expectedStatus, t.name+"status code should match")

	if t.statusFailed && t.autoDumpResponse {
		t.dumpResponse()
	}

	return t
}

// CmpHeader tests the last request response header against
// expectedHeader. expectedHeader can be a http.Header or a TestDeep
// operator. Keep in mind that if it is a http.Header, it has to match
// exactly the response header. Often only the presence of a
// header key is needed:
//
//   ta := tdhttp.NewTestAPI(t, mux).
//     PostJSON("/new", map[string]string{"name": "Bob"}).
//     CmdStatus(201).
//     CmpHeader(td.ContainsKey("X-Custom"))
//
// or some specific key, value pairs:
//
//   ta.CmpHeader(td.SuperMapOf(
//     http.Header{
//       "X-Account": []string{"Bob"},
//     },
//     td.MapEntries{
//       "X-Token": td.Re(`^[a-z0-9-]{32}\z`),
//     }),
//   )
//
// Note that CmpHeader calls can be chained:
//
//   ta.CmpHeader(td.ContainsKey("X-Account")).
//     CmpHeader(td.ContainsKey("X-Token"))
//
// instead of doing all tests in one call as All operator allows it:
//
//   ta.CmpHeader(td.All(
//     td.ContainsKey("X-Account"),
//     td.ContainsKey("X-Token")),
//   )
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpHeader(expectedHeader interface{}) *TestAPI {
	defer t.t.AnchorsPersistTemporarily()()

	t.t.Helper()

	if !t.checkRequestSent() {
		t.headerFailed = true
		return t
	}

	t.headerFailed = !t.t.RootName("Response.Header").
		CmpLax(t.response.Header(), expectedHeader, t.name+"header should match")

	if t.headerFailed && t.autoDumpResponse {
		t.dumpResponse()
	}

	return t
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

func (t *TestAPI) cmpMarshaledBody(
	acceptEmptyBody bool,
	unmarshal func([]byte, interface{}) error,
	expectedBody interface{},
) *TestAPI {
	defer t.t.AnchorsPersistTemporarily()()

	t.t.Helper()

	if t.bodyFailed = !t.checkRequestSent(); t.bodyFailed {
		return t
	}

	if !acceptEmptyBody &&
		!t.t.RootName("Response body").Code(t.response.Body.Bytes(),
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
			t.name+"body should not be empty") {
		t.bodyFailed = true
		if t.autoDumpResponse {
			t.dumpResponse()
		}
		return t
	}

	tt := t.t.RootName("Response.Body")

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
			// unmarshal in an interface{}
			bodyType = types.Interface
			unknownExpectedType = true

			// Special case for Ignore & NotEmpty operators
			switch op.GetLocation().Func {
			case "Ignore", "NotEmpty":
				showRawBody = t.statusFailed // Show real body if status failed
			}
		}
	} else {
		bodyType = reflect.TypeOf(expectedBody)
	}

	// For unmarshaling below, body must be a pointer
	bodyPtr := reflect.New(bodyType)

	// Try to unmarshal body
	if !tt.RootName("unmarshal(Response.Body)").
		CmpNoError(unmarshal(t.response.Body.Bytes(), bodyPtr.Interface()), t.name+"body unmarshaling") {
		// If unmarshal failed, perhaps it's coz the expected body type
		// is unknown?
		if unknownExpectedType {
			tt.Logf("Cannot guess the body expected type as %[1]s TestDeep\n"+
				"operator does not know the type behind it.\n"+
				"You can try All(Isa(EXPECTED_TYPE), %[1]s(…)) to disambiguate…",
				op.GetLocation().Func)
		}
		showRawBody = true // let's show its real body contents
		t.bodyFailed = true
	} else if !tt.Cmp(bodyPtr.Elem().Interface(), expectedBody, t.name+"body contents is OK") {
		// Try to catch bad body expected type when nothing has been set
		// to non-zero during unmarshaling body. In this case, require
		// to show raw body contents.
		if len(t.response.Body.Bytes()) > 0 &&
			td.EqDeeply(bodyPtr.Interface(), reflect.New(bodyType).Interface()) {
			showRawBody = true
			tt.Log("Hmm… It seems nothing has been set during unmarshaling…")
		}
		t.bodyFailed = true
	}

	if showRawBody || (t.bodyFailed && t.autoDumpResponse) {
		t.dumpResponse()
	}

	return t
}

// CmpMarshaledBody tests that the last request response body can be
// unmarshalled using unmarhsal function and then, that it matches
// expectedBody. expectedBody can be any type unmarshal function can
// handle, or a TestDeep operator.
//
// See CmpJSONBody and CmpXMLBody sources for examples of use.
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpMarshaledBody(unmarshal func([]byte, interface{}) error, expectedBody interface{}) *TestAPI {
	t.t.Helper()
	return t.cmpMarshaledBody(false, unmarshal, expectedBody)
}

// CmpBody tests the last request response body against
// expectedBody. expectedBody can be a []byte, a string or a TestDeep
// operator.
//
//   ta := tdhttp.NewTestAPI(t, mux)
//
//   ta.Get("/test").
//     CmpStatus(http.StatusOK).
//     CmpBody("OK!\n")
//
//   ta.Get("/test").
//     CmpStatus(http.StatusOK).
//     CmpBody(td.Contains("OK"))
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpBody(expectedBody interface{}) *TestAPI {
	t.t.Helper()

	if expectedBody == nil {
		return t.NoBody()
	}

	return t.cmpMarshaledBody(
		true, // accept empty body
		func(body []byte, target interface{}) error {
			switch target := target.(type) {
			case *string:
				*target = string(body)
			case *[]byte:
				*target = body
			case *interface{}:
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
// encoding/json.Unmarshall'ed and that it matches
// expectedBody. expectedBody can be any type encoding/json can
// Unmarshal into, or a TestDeep operator.
//
//   ta := tdhttp.NewTestAPI(t, mux)
//
//   ta.Get("/person/42").
//     CmpStatus(http.StatusOK).
//     CmpJSONBody(Person{
//       ID:   42,
//       Name: "Bob",
//       Age:  26,
//     })
//
//   ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//     CmpStatus(http.StatusCreated).
//     CmpJSONBody(td.SStruct(
//       Person{
//         Name: "Bob",
//         Age:  26,
//       },
//       td.StructFields{
//         "ID": td.NotZero(),
//       }))
//
// The same with anchoring, and so without td.SStruct():
//
//   ta := tdhttp.NewTestAPI(tt, mux)
//
//   ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//     CmpStatus(http.StatusCreated).
//     CmpJSONBody(Person{
//       ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//       Name: "Bob",
//       Age:  26,
//     })
//
// The same using td.JSON():
//
//   ta.PostJSON("/person", Person{Name: "Bob", Age: 23}).
//     CmpStatus(http.StatusCreated).
//     CmpJSONBody(td.JSON(`
//   {
//     "id":   NotZero(),
//     "name": "Bob",
//     "age":  26
//   }`))
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpJSONBody(expectedBody interface{}) *TestAPI {
	t.t.Helper()
	return t.CmpMarshaledBody(json.Unmarshal, expectedBody)
}

// CmpXMLBody tests that the last request response body can be
// encoding/xml.Unmarshall'ed and that it matches
// expectedBody. expectedBody can be any type encoding/xml can
// Unmarshal into, or a TestDeep operator.
//
//   ta := tdhttp.NewTestAPI(t, mux)
//
//   ta.Get("/person/42").
//     CmpStatus(http.StatusOK).
//     CmpXMLBody(Person{
//       ID:   42,
//       Name: "Bob",
//       Age:  26,
//     })
//
//   ta.Get("/person/43").
//     CmpStatus(http.StatusOK).
//     CmpXMLBody(td.SStruct(
//       Person{
//         Name: "Bob",
//         Age:  26,
//       },
//       td.StructFields{
//         "ID": td.NotZero(),
//       }))
//
// The same with anchoring:
//
//   ta := tdhttp.NewTestAPI(tt, mux)
//
//   ta.Get("/person/42").
//     CmpStatus(http.StatusOK).
//     CmpXMLBody(Person{
//       ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//       Name: "Bob",
//       Age:  26,
//     })
//
// It fails if no request has been sent yet.
func (t *TestAPI) CmpXMLBody(expectedBody interface{}) *TestAPI {
	t.t.Helper()
	return t.CmpMarshaledBody(xml.Unmarshal, expectedBody)
}

// NoBody tests that the last request response body is empty.
//
// It fails if no request has been sent yet.
func (t *TestAPI) NoBody() *TestAPI {
	defer t.t.AnchorsPersistTemporarily()()

	t.t.Helper()

	if !t.checkRequestSent() {
		t.bodyFailed = true
		return t
	}

	t.bodyFailed = !t.t.RootName("Response.Body").
		Code(len(t.response.Body.Bytes()) == 0,
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

	if t.bodyFailed { // Systematically dump response, no AutoDumpResponse needed
		t.dumpResponse()
	}

	return t
}

// Or executes function "fn" if t.Failed() is true at the moment it is called.
//
// "fn" can have several types:
//   - func(body string) or func(t *td.T, body string)
//     → "fn" is called with response body as a string.
//       If no response has been received yet, body is "";
//   - func(body []byte) or func(t *td.T, body []byte)
//     → "fn" is called with response body as a []byte.
//       If no response has been received yet, body is nil;
//   - func(t *td.T, resp *httptest.ResponseRecorder)
//     → "fn" is called with the internal object containing the response.
//       See net/http/httptest for details.
//       If no response has been received yet, resp is nil.
//
// If "fn" type is not one of these types, it calls t.T().Fatal().
func (t *TestAPI) Or(fn interface{}) *TestAPI {
	t.t.Helper()
	switch fn := fn.(type) {
	case func(string):
		if t.Failed() {
			var body string
			if t.response != nil && t.response.Body != nil {
				body = t.response.Body.String()
			}
			fn(body)
		}

	case func(*td.T, string):
		if t.Failed() {
			var body string
			if t.response != nil && t.response.Body != nil {
				body = t.response.Body.String()
			}
			fn(t.t, body)
		}

	case func([]byte):
		if t.Failed() {
			var body []byte
			if t.response != nil && t.response.Body != nil {
				body = t.response.Body.Bytes()
			}
			fn(body)
		}

	case func(*td.T, []byte):
		if t.Failed() {
			var body []byte
			if t.response != nil && t.response.Body != nil {
				body = t.response.Body.Bytes()
			}
			fn(t.t, body)
		}

	case func(*td.T, *httptest.ResponseRecorder):
		if t.Failed() {
			fn(t.t, t.response)
		}

	default:
		t.t.Fatal(color.BadUsage(
			"Or(func([*td.T,]string) | func([*td.T,][]byte) | func(*td.T,*httptest.ResponseRecorder))",
			fn, 1, true))
	}

	return t
}

// OrDumpResponse dumps the response if at least one previous test failed.
//
//   ta := tdhttp.NewTestAPI(t, handler)
//
//   ta.Get("/foo").
//     CmpStatus(200).
//     OrDumpResponse(). // if status check failed, dumps the response
//     CmpBody("bar")    // if it fails, the response is not dumped
//
//   ta.Get("/foo").
//     CmpStatus(200).
//     CmpBody("bar").
//     OrDumpResponse() // dumps the response if status and/or body checks fail
//
// See AutoDumpResponse method to automatize this dump.
func (t *TestAPI) OrDumpResponse() *TestAPI {
	if t.Failed() {
		t.dumpResponse()
	}
	return t
}

func (t *TestAPI) dumpResponse() {
	if t.responseDumped {
		return
	}

	t.t.Helper()
	if t.response != nil {
		t.responseDumped = true
		internal.DumpResponse(t.t, t.response.Result())
		return
	}

	t.t.Logf("No response received yet")
}

// Anchor returns a typed value allowing to anchor the TestDeep
// operator "operator" in a go classic litteral like a struct, slice,
// array or map value.
//
//   ta := tdhttp.NewTestAPI(tt, mux)
//
//   ta.Get("/person/42").
//     CmpStatus(http.StatusOK).
//     CmpJSONBody(Person{
//       ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//       Name: "Bob",
//       Age:  26,
//     })
//
// See (*td.T).Anchor documentation for details
// https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor
//
// See A method for a shorter synonym of Anchor.
func (t *TestAPI) Anchor(operator td.TestDeep, model ...interface{}) interface{} {
	return t.t.Anchor(operator, model...)
}

// A is a synonym for Anchor. It returns a typed value allowing to
// anchor the TestDeep operator "operator" in a go classic litteral
// like a struct, slice, array or map value.
//
//   ta := tdhttp.NewTestAPI(tt, mux)
//
//   ta.Get("/person/42").
//     CmpStatus(http.StatusOK).
//     CmpJSONBody(Person{
//       ID:   ta.A(td.NotZero(), uint64(0)).(uint64),
//       Name: "Bob",
//       Age:  26,
//     })
//
// See (*td.T).Anchor documentation
// https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor
// for details.
func (t *TestAPI) A(operator td.TestDeep, model ...interface{}) interface{} {
	return t.Anchor(operator, model...)
}

// SentAt returns the time just before the last request is handled. It
// can be used to check the time a route sets and returns, as in:
//
//   ta.PostJSON("/person/42", Person{Name: "Bob", Age: 23}).
//     CmpStatus(http.StatusCreated).
//     CmpJSONBody(Person{
//       ID:        ta.A(td.NotZero(), uint64(0)).(uint64),
//       Name:      "Bob",
//       Age:       23,
//       CreatedAt: ta.A(td.Between(ta.SentAt(), time.Now())).(time.Time),
//     })
//
// checks that CreatedAt is included between the time when the request
// has been sent, and the time when the comparison occurs.
func (t *TestAPI) SentAt() time.Time {
	return t.sentAt
}
