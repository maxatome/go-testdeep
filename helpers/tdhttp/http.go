// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package tdhttp provides some functions to easily test HTTP handlers.
package tdhttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	td "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

// Response is used by Cmp*Response functions to make the HTTP
// response match easier. Each field, can be a TestDeep operator as
// well as the exact expected value.
type Response struct {
	Status interface{} // Status is the expected status (ignored if nil)
	Header interface{} // Header is the expected header (ignored if nil)
	Body   interface{} // Body is the expected body (expected to be empty if nil)
}

// CmpMarshaledResponse is the base function used by all others in
// tdhttp package. The req *http.Request is launched against
// handler. The response body is unmarshaled using unmarshal. The
// response is then tested against expectedResp.
//
// It returns true if the tests succeed, false otherwise.
func CmpMarshaledResponse(tt td.TestingFT,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	unmarshal func([]byte, interface{}) error,
	expectedResp Response,
	args ...interface{},
) bool {
	tt.Helper()

	if testName := tdutil.BuildTestName(args...); testName != "" {
		tt.Log(testName)
	}

	var statusMismatch, headerMismatch bool

	t := td.NewT(tt)
	defer t.AnchorsPersistTemporarily()()

	w := httptest.NewRecorder()

	handler(w, req)

	// Check status, nil = ignore
	if expectedResp.Status != nil {
		statusMismatch = !t.RootName("Response.Status").
			Cmp(w.Code, expectedResp.Status, "status code should match")
	}

	// Check header, nil = ignore
	if expectedResp.Header != nil {
		headerMismatch = !t.RootName("Response.Header").
			Cmp(w.Header(), expectedResp.Header, "header should match")
	}

	t = t.RootName("Response.Body")

	// Body, nil = no body expected
	if expectedResp.Body == nil {
		ok := t.Empty(w.Body.Bytes(), "body should be empty")
		return !statusMismatch && !headerMismatch && ok
	}

	if !t.NotEmpty(w.Body.Bytes(), "body should not be empty") {
		return false
	}

	var (
		bodyType reflect.Type
		body     interface{}
	)

	// If expectedBody is a TestDeep operator, try to ask it the type
	// behind it. It should work in most cases (typically Struct(),
	// Map() & Slice()).
	var unknownExpectedType, showRawBody bool
	op, ok := expectedResp.Body.(td.TestDeep)
	if ok {
		bodyType = op.TypeBehind()
		if bodyType == nil {
			// As the expected body type cannot be guessed, try to
			// unmarshal in an interface{}
			bodyType = reflect.TypeOf(&body).Elem()
			unknownExpectedType = true

			// Special case for Ignore & NotEmpty operators
			switch op.GetLocation().Func {
			case "Ignore", "NotEmpty":
				showRawBody = statusMismatch // Show real body if status failed
			}
		}
	} else {
		bodyType = reflect.TypeOf(expectedResp.Body)
	}

	// For unmarshaling below, body must be a pointer
	body = reflect.New(bodyType).Interface()

	success := !statusMismatch && !headerMismatch

	// Try to unmarshal body
	if !t.RootName("unmarshal(Response.Body)").
		CmpNoError(unmarshal(w.Body.Bytes(), body), "body unmarshaling") {
		// If unmarshal failed, perhaps it's coz the expected body type
		// is unknown?
		if unknownExpectedType {
			t.Logf("Cannot guess the body expected type as %[1]s TestDeep\n"+
				"operator does not know the type behind it.\n"+
				"You can try All(Isa(EXPECTED_TYPE), %[1]s(...)) to disambiguate...",
				op.GetLocation().Func)
		}
		success = false
		showRawBody = true // let's show its real body contents
	} else if !t.Cmp(body, td.Ptr(expectedResp.Body), "body contents is OK") {
		// If the body comparison fails
		success = false

		// Try to catch bad body expected type when nothing has been set
		// to non-zero during unmarshaling body. In this case, require
		// to show raw body contents.
		if td.EqDeeply(body, reflect.New(bodyType).Interface()) {
			showRawBody = true
			t.Log("Hmm… It seems nothing has been set during unmarshaling…")
		}
	}

	if showRawBody {
		t.Logf("Raw received body: %s", tdutil.FormatString(w.Body.String()))
	}

	return success
}

// CmpResponse is used to match a []byte or string response body. The
// req *http.Request is launched against handler. If expectedResp.Body
// is non-nil, the response body is converted to []byte or string,
// depending on the expectedResp.Body type. The response is then
// tested against expectedResp.
//
// It returns true if the tests succeed, false otherwise.
func CmpResponse(t td.TestingFT,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...interface{}) bool {
	t.Helper()
	return CmpMarshaledResponse(t,
		req,
		handler,
		func(body []byte, target interface{}) error {
			switch t := target.(type) {
			case *string:
				*t = string(body)
			case *[]byte:
				*t = body
			case *interface{}:
				*t = body
			default:
				return fmt.Errorf(
					"CmpResponse does not handle %T body, only string & []byte",
					target)
			}
			return nil
		},
		expectedResp,
		args...)
}

// CmpJSONResponse is used to match a JSON response body. The req
// *http.Request is launched against handler. If expectedResp.Body is
// non-nil, the response body is json.Unmarshal'ed. The response is
// then tested against expectedResp.
//
// It returns true if the tests succeed, false otherwise.
func CmpJSONResponse(t td.TestingFT,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...interface{},
) bool {
	t.Helper()
	return CmpMarshaledResponse(t,
		req,
		handler,
		json.Unmarshal,
		expectedResp,
		args...)
}

// CmpXMLResponse is used to match an XML response body. The req
// *http.Request is launched against handler. If expectedResp.Body is
// non-nil, the response body is xml.Unmarshal'ed. The response is
// then tested against expectedResp.
//
// It returns true if the tests succeed, false otherwise.
func CmpXMLResponse(t td.TestingFT,
	req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response,
	args ...interface{},
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
// testing.Run, calling CmpMarshaledResponse behind the scene. As it
// is intended to be used in conjunction with testing.Run() which
// names the sub-test, the test name part (args...) is voluntary
// omitted.
//
//   t.Run("Subtest name", tdhttp.CmpMarshaledResponseFunc(
//     tdhttp.Get("/text"),
//     mux.ServeHTTP,
//     tdhttp.Response{
//       Status: http.StatusOK,
//     }))
func CmpMarshaledResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	unmarshal func([]byte, interface{}) error,
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpMarshaledResponse(t, req, handler, unmarshal, expectedResp)
	}
}

// CmpResponseFunc returns a function ready to be used with
// testing.Run, calling CmpResponse behind the scene. As it is
// intended to be used in conjunction with testing.Run() which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//   t.Run("Subtest name", tdhttp.CmpResponseFunc(
//     tdhttp.Get("/text"),
//     mux.ServeHTTP,
//     tdhttp.Response{
//       Status: http.StatusOK,
//     }))
func CmpResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpResponse(t, req, handler, expectedResp)
	}
}

// CmpJSONResponseFunc returns a function ready to be used with
// testing.Run, calling CmpJSONResponse behind the scene. As it is
// intended to be used in conjunction with testing.Run() which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//   t.Run("Subtest name", tdhttp.CmpJSONResponseFunc(
//     tdhttp.Get("/json"),
//     mux.ServeHTTP,
//     tdhttp.Response{
//       Status: http.StatusOK,
//       Body:   JResp{Comment: "expected comment!"},
//     }))
func CmpJSONResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpJSONResponse(t, req, handler, expectedResp)
	}
}

// CmpXMLResponseFunc returns a function ready to be used with
// testing.Run, calling CmpXMLResponse behind the scene. As it is
// intended to be used in conjunction with testing.Run() which names
// the sub-test, the test name part (args...) is voluntary omitted.
//
//   t.Run("Subtest name", tdhttp.CmpXMLResponseFunc(
//     tdhttp.Get("/xml"),
//     mux.ServeHTTP,
//     tdhttp.Response{
//       Status: http.StatusOK,
//       Body:   JResp{Comment: "expected comment!"},
//     }))
func CmpXMLResponseFunc(req *http.Request,
	handler func(w http.ResponseWriter, r *http.Request),
	expectedResp Response) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		CmpXMLResponse(t, req, handler, expectedResp)
	}
}
