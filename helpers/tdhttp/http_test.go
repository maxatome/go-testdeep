// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	td "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

type CmpResponseTest struct {
	Name         string
	Handler      func(w http.ResponseWriter, r *http.Request)
	Success      bool
	ExpectedResp tdhttp.Response
	ExpectedLogs []string
}

func TestCmpResponse(tt *testing.T) {
	handlerNonEmpty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestDeep", "foobar")
		w.WriteHeader(242)
		fmt.Fprintln(w, "text response")
	})

	handlerEmpty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestDeep", "zip")
		w.WriteHeader(424)
	})

	t := td.NewT(tt)

	for _, curTest := range []CmpResponseTest{
		// Non-empty success
		{
			Name:         "string body only",
			Handler:      handlerNonEmpty,
			Success:      true,
			ExpectedResp: tdhttp.Response{Body: "text response\n"},
		},
		{
			Name:    "[]byte status + body",
			Handler: handlerNonEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Body:   []byte("text response\n"),
			},
		},
		{
			Name:    "[]byte status + header + body",
			Handler: handlerNonEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   []byte("text response\n"),
			},
		},
		{
			Name:    "with TestDeep operators",
			Handler: handlerNonEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: td.Between(200, 300),
				Header: td.ContainsKey("X-Testdeep"),
				Body: td.All(
					td.Isa(""), // enforces TypeBehind → string
					td.Contains("response"),
				),
			},
		},
		{
			Name:    "ignore body explicitely",
			Handler: handlerNonEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   td.Ignore(),
			},
		},
		{
			Name:    "body just not empty",
			Handler: handlerNonEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   td.NotEmpty(),
			},
		},
		// Non-empty failures
		{
			Name:         "bad status",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Status: 243, Body: td.NotEmpty()},
			ExpectedLogs: []string{
				"status code should match",
				"Raw received body:",
				"text response", // check the complete body is shown
			},
		},
		{
			Name:         "bad body",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Body: "BAD!"},
			ExpectedLogs: []string{
				"Failed test 'body contents is OK'",
				`*Response.Body: values differ`,
				"text response", // check the response is shown
			},
		},
		{
			Name:         "bad body, as non-empty",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{},
			ExpectedLogs: []string{
				"Failed test 'body should be empty'",
				"text response", // check the response is shown
			},
		},
		{
			Name:         "bad header",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Header: http.Header{"X-Testdeep": []string{"zzz"}}},
			ExpectedLogs: []string{
				"Failed test 'header should match'",
				`Response.Header["X-Testdeep"][0]: values differ`,
			},
		},
		{
			Name:         "cannot unmarshal",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Body: 123},
			ExpectedLogs: []string{
				"Failed test 'body unmarshaling'",
				"unmarshal(Response.Body): should NOT be an error",
				"Raw received body:",
				"text response", // check the complete body is shown
			},
		},
		// Empty success
		{
			Name:    "empty body",
			Handler: handlerEmpty,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 424,
				Header: http.Header{"X-Testdeep": []string{"zip"}},
				Body:   nil,
			},
		},
		// Empty failures
		{
			Name:         "should not be empty",
			Handler:      handlerEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Body: "NOT EMPTY!"},
			ExpectedLogs: []string{
				"Failed test 'body should not be empty'",
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpResponse, "CmpResponse", curTest)
			})
	}
}

func TestCmpJSONResponse(tt *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestDeep", "foobar")
		w.WriteHeader(242)
		fmt.Fprintln(w, `{"name":"Bob"}`)
	})

	type JResp struct {
		Name string `json:"name"`
	}

	t := td.NewT(tt)

	for _, curTest := range []CmpResponseTest{
		// Success
		{
			Name:    "JSON OK",
			Handler: handler,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   JResp{Name: "Bob"},
			},
		},
		{
			Name:    "JSON ptr OK",
			Handler: handler,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   &JResp{Name: "Bob"},
			},
		},
		// Failure
		{
			Name:    "JSON failure",
			Handler: handler,
			Success: false,
			ExpectedResp: tdhttp.Response{
				Body: 123,
			},
			ExpectedLogs: []string{
				"Failed test 'body unmarshaling'",
				"unmarshal(Response.Body): should NOT be an error",
				"Raw received body:",
				`{"name":"Bob"}`, // check the complete body is shown
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpJSONResponse, "CmpJSONResponse", curTest)
			})
	}
}

func TestCmpXMLResponse(tt *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestDeep", "foobar")
		w.WriteHeader(242)
		fmt.Fprintln(w, `<XResp><name>Bob</name></XResp>`)
	})

	type XResp struct {
		Name string `xml:"name"`
	}

	t := td.NewT(tt)

	for _, curTest := range []CmpResponseTest{
		// Success
		{
			Name:    "XML OK",
			Handler: handler,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   XResp{Name: "Bob"},
			},
		},
		{
			Name:    "XML ptr OK",
			Handler: handler,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Status: 242,
				Header: http.Header{"X-Testdeep": []string{"foobar"}},
				Body:   &XResp{Name: "Bob"},
			},
		},
		// Failure
		{
			Name:    "XML failure",
			Handler: handler,
			Success: false,
			ExpectedResp: tdhttp.Response{
				// xml.Unmarshal does not raise an error when trying to
				// unmarshal in an int, as json does...
				Body: func() {},
			},
			ExpectedLogs: []string{
				"Failed test 'body unmarshaling'",
				"unmarshal(Response.Body): should NOT be an error",
				"Raw received body:",
				`<XResp><name>Bob</name></XResp>`, // check the complete body is shown
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpXMLResponse, "CmpXMLResponse", curTest)
			})
	}
}

func testCmpResponse(t *td.T,
	cmp func(td.TestingFT, *http.Request, func(http.ResponseWriter, *http.Request), tdhttp.Response, ...interface{}) bool,
	cmpName string,
	curTest CmpResponseTest,
) {
	t.Helper()

	mockT := &tdutil.T{}

	t.CmpDeeply(cmp(mockT,
		httptest.NewRequest("GET", "/path", nil),
		curTest.Handler,
		curTest.ExpectedResp),
		curTest.Success)

	dumpLogs := !t.CmpDeeply(mockT.Failed(), !curTest.Success)

	for _, expectedLog := range curTest.ExpectedLogs {
		if !strings.Contains(mockT.LogBuf(), expectedLog) {
			t.Errorf(`"%s" not found in test logs`, expectedLog)
			dumpLogs = true
		}
	}

	if dumpLogs {
		t.Errorf(`Test logs: "%s"`, mockT.LogBuf())
	}
}

func TestMux(t *testing.T) {
	mux := http.NewServeMux()

	// GET /text
	mux.HandleFunc("/text", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.NotFound(w, req)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Text result!")
	})

	// GET /json
	mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if req.Method != "GET" {
			fmt.Fprintf(w, `{"code":404,"message":"Not found"}`)
			return
		}
		fmt.Fprintf(w, `{"comment":"JSON result!"}`)
	})

	// GET /xml
	mux.HandleFunc("/xml", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.NotFound(w, req)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<XResp><comment>XML result!</comment></XResp>`)
	})

	//
	// Check GET /text route
	t.Run("/text route", func(t *testing.T) {
		tdhttp.CmpResponse(t, tdhttp.NewRequest("GET", "/text", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"text/plain"},
				}, nil),
				Body: "Text result!",
			},
			"GET /text should return 200 + text/plain + Text result!")

		tdhttp.CmpResponse(t, tdhttp.NewRequest("PATCH", "/text", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusNotFound,
				Body:   td.Ignore(),
			},
			"PATCH /text should return Not Found")

		t.Run("/text route via CmpResponseFunc", tdhttp.CmpResponseFunc(
			tdhttp.NewRequest("GET", "/text", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"text/plain"},
				}, nil),
				Body: "Text result!",
			}))

		t.Run("/text route via CmpMarshaledResponseFunc", tdhttp.CmpMarshaledResponseFunc(
			tdhttp.NewRequest("GET", "/text", nil),
			mux.ServeHTTP,
			func(body []byte, target interface{}) error {
				*target.(*string) = string(body)
				return nil
			},
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"text/plain"},
				}, nil),
				Body: "Text result!",
			}))
	})

	//
	// Check GET /json
	t.Run("/json route", func(t *testing.T) {
		type JResp struct {
			Comment string `json:"comment"`
		}
		tdhttp.CmpJSONResponse(t, tdhttp.NewRequest("GET", "/json", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/json"},
				}, nil),
				Body: JResp{Comment: "JSON result!"},
			},
			"GET /json should return 200 + application/json + comment=JSON result!")

		tdhttp.CmpJSONResponse(t, tdhttp.NewRequest("GET", "/json", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/json"},
				}, nil),
				Body: td.Struct(JResp{}, td.StructFields{
					"Comment": td.Contains("result!"),
				}),
			},
			"GET /json should return 200 + application/json + comment=~result!")

		t.Run("/json route via CmpJSONResponseFunc", tdhttp.CmpJSONResponseFunc(
			tdhttp.NewRequest("GET", "/json", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/json"},
				}, nil),
				Body: JResp{Comment: "JSON result!"},
			}))

		// We expect to receive a specific message, but a complete
		// different one is received (here a Not Found, as PUT is used).
		// In this case, a log message tell us that nothing has been set
		// during unmarshaling AND the received body should be dumped
		t.Run("zeroed body", func(tt *testing.T) {
			t := td.NewT(tt)

			mockT := &tdutil.T{}
			ok := tdhttp.CmpJSONResponse(mockT,
				tdhttp.NewRequest("PUT", "/json", nil),
				mux.ServeHTTP,
				tdhttp.Response{
					Status: http.StatusOK,
					Header: td.SuperMapOf(http.Header{
						"Content-Type": []string{"application/json"},
					}, nil),
					Body: JResp{Comment: "JSON result!"},
				})

			t.False(ok)
			t.True(mockT.Failed())
			t.Contains(mockT.LogBuf(), "nothing has been set during unmarshaling")
			t.Contains(mockT.LogBuf(), "Raw received body:")
		})
	})

	//
	// Check GET /xml
	t.Run("/xml route", func(t *testing.T) {
		type XResp struct {
			Comment string `xml:"comment"`
		}
		tdhttp.CmpXMLResponse(t, tdhttp.NewRequest("GET", "/xml", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/xml"},
				}, nil),
				Body: XResp{Comment: "XML result!"},
			},
			"GET /xml should return 200 + application/xml + comment=XML result!")

		tdhttp.CmpXMLResponse(t, tdhttp.NewRequest("GET", "/xml", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/xml"},
				}, nil),
				Body: td.Struct(XResp{}, td.StructFields{
					"Comment": td.Contains("result!"),
				}),
			},
			"GET /xml should return 200 + application/xml + comment=~result!")

		t.Run("/xml route via CmpXMLResponseFunc", tdhttp.CmpXMLResponseFunc(
			tdhttp.NewRequest("GET", "/xml", nil),
			mux.ServeHTTP,
			tdhttp.Response{
				Status: http.StatusOK,
				Header: td.SuperMapOf(http.Header{
					"Content-Type": []string{"application/xml"},
				}, nil),
				Body: td.Struct(XResp{}, td.StructFields{
					"Comment": td.Contains("result!"),
				}),
			}))

		// We expect to receive a specific message into a not specific
		// type (behind the TestDeep operator). The XML unmarshaling
		// fails, so a log message should tell us to be more clear
		// concerning the expected body type (in case it is the origin of
		// the problem).
		t.Run("Unmarshal is failing", func(tt *testing.T) {
			t := td.NewT(tt)

			mockT := &tdutil.T{}
			ok := tdhttp.CmpXMLResponse(mockT,
				tdhttp.NewRequest("PUT", "/xml", nil),
				mux.ServeHTTP,
				tdhttp.Response{
					Status: http.StatusOK,
					Header: td.SuperMapOf(http.Header{
						"Content-Type": []string{"application/xml"},
					}, nil),
					// This TestDeep operators combination is absurd. It is only
					// intended to avoid CmpXMLResponse detects the expected body
					// type
					Body: td.Any(XResp{Comment: "XML result!"}, 12),
				})

			t.False(ok)
			t.True(mockT.Failed())
			t.Contains(mockT.LogBuf(),
				"Cannot guess the body expected type as Any TestDeep")
			t.Contains(mockT.LogBuf(),
				"You can try All(Isa(EXPECTED_TYPE), Any(...)) to disambiguate")
			t.Contains(mockT.LogBuf(), "Raw received body:")
		})
	})
}
