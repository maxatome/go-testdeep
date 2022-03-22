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
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/td"
)

func TestMain(m *testing.M) {
	color.SaveState()
	os.Exit(m.Run())
}

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

	cookie := http.Cookie{Name: "Cookies-Testdeep", Value: "foobar"}
	handlerWithCokies := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestDeep", "cookies")
		http.SetCookie(w, &cookie)
		w.WriteHeader(242)
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
			Name:    "ignore body explicitly",
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
				`~ Failed test 'status code should match'
\s+Response.Status: values differ
\s+got: 242
\s+expected: 243`,
				`~ Received response:
\s+\x60(?s:.+?)
\s+
\s+text response
\s+\x60
`, // check the complete body is shown
			},
		},
		{
			Name:         "bad body",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Body: "BAD!"},
			ExpectedLogs: []string{
				`~ Failed test 'body contents is OK'
\s+Response.Body: values differ
\s+got: \x60text response
\s+\x60
\s+expected: "BAD!"`,
			},
		},
		{
			Name:         "bad body, as non-empty",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{},
			ExpectedLogs: []string{
				`~ Failed test 'body should be empty'
\s+Response.Body is not empty
\s+got: not empty
\s+expected: empty`,
				`~ Received response:
\s+\x60(?s:.+?)
\s+
\s+text response
\s+\x60
`, // check the complete body is shown
			},
		},
		{
			Name:    "bad header",
			Handler: handlerNonEmpty,
			Success: false,
			ExpectedResp: tdhttp.Response{
				Header: http.Header{"X-Testdeep": []string{"zzz"}},
				Body:   td.NotEmpty(),
			},
			ExpectedLogs: []string{
				`~ Failed test 'header should match'
\s+Response.Header\["X-Testdeep"\]\[0\]: values differ
\s+got: "foobar"
\s+expected: "zzz"`,
			},
		},
		{
			Name:    "bad cookies",
			Handler: handlerWithCokies,
			Success: false,
			ExpectedResp: tdhttp.Response{
				Cookies: []*http.Cookie{{Name: "Cookies-Testdeep", Value: "squalala"}},
				Body:    td.Empty(),
			},
			ExpectedLogs: []string{
				`~ Failed test 'cookies should match'
\s+Response.Cookie\[0\]\.Value: values differ
\s+got: "foobar"
\s+expected: "squalala"`,
			},
		},
		{
			Name:    "good cookies",
			Handler: handlerWithCokies,
			Success: true,
			ExpectedResp: tdhttp.Response{
				Header: http.Header{
					"X-Testdeep": []string{"cookies"},
					"Set-Cookie": []string{cookie.String()},
				},
				Cookies: []*http.Cookie{&cookie},
				Body:    td.Empty(),
			},
		},
		{
			Name:         "cannot unmarshal",
			Handler:      handlerNonEmpty,
			Success:      false,
			ExpectedResp: tdhttp.Response{Body: 123},
			ExpectedLogs: []string{
				`~ Failed test 'body unmarshaling'
\s+unmarshal\(Response\.Body\): should NOT be an error
\s+got: .*Cmp(Response|Body) only accepts expected(Resp\.)?Body be a \[\]byte, a string or a TestDeep operator allowing to match these types, but not type int.*
\s+expected: nil`,
				`~ Received response:
\s+\x60(?s:.+?)
\s+
\s+text response
\s+\x60
`, // check the complete body is shown
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
				`~ Failed test 'body contents is OK'
\s+Response.Body: values differ
\s+got: ""
\s+expected: "NOT EMPTY!"`,
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpResponse, "CmpResponse", curTest)
			})

		t.Run(curTest.Name+" TestAPI",
			func(t *td.T) {
				testTestAPI(t, (*tdhttp.TestAPI).CmpBody, "CmpBody", curTest)
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
				`~ Failed test 'body unmarshaling'
\s+unmarshal\(Response\.Body\): should NOT be an error
\s+got: .*cannot unmarshal object into Go value of type int.*
\s+expected: nil`,
				`~ Received response:
\s+\x60(?s:.+?)
\s+
\s+\{"name":"Bob"\}
\s+\x60
`, // check the complete body is shown
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpJSONResponse, "CmpJSONResponse", curTest)
			})

		t.Run(curTest.Name+" TestAPI",
			func(t *td.T) {
				testTestAPI(t, (*tdhttp.TestAPI).CmpJSONBody, "CmpJSONBody", curTest)
			})
	}
}

func TestCmpJSONResponseAnchor(tt *testing.T) {
	t := td.NewT(tt)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(242)
		fmt.Fprintln(w, `{"name":"Bob"}`)
	})
	req := httptest.NewRequest("GET", "/path", nil)

	type JResp struct {
		Name string `json:"name"`
	}

	// With *td.T
	tdhttp.CmpJSONResponse(t, req, handler,
		tdhttp.Response{
			Status: 242,
			Body: JResp{
				Name: t.A(td.Re("(?i)bob"), "").(string),
			},
		})

	// With *testing.T
	tdhttp.CmpJSONResponse(tt, req, handler,
		tdhttp.Response{
			Status: 242,
			Body: JResp{
				Name: t.A(td.Re("(?i)bob"), "").(string),
			},
		})

	func() {
		defer t.AnchorsPersistTemporarily()()

		op := t.A(td.Re("(?i)bob"), "").(string)

		// All calls should succeed, as op persists
		tdhttp.CmpJSONResponse(t, req, handler,
			tdhttp.Response{
				Status: 242,
				Body:   JResp{Name: op},
			})

		tdhttp.CmpJSONResponse(t, req, handler,
			tdhttp.Response{
				Status: 242,
				Body:   JResp{Name: op},
			})

		// Even with the original *testing.T instance (here tt)
		tdhttp.CmpJSONResponse(tt, req, handler,
			tdhttp.Response{
				Status: 242,
				Body:   JResp{Name: op},
			})
	}()

	// Failures
	t.FailureIsFatal().False(t.DoAnchorsPersist()) // just to be sure

	mt := td.NewT(tdutil.NewT("tdhttp_persistence_test"))
	op := mt.A(td.Re("(?i)bob"), "").(string)

	// First call should succeed
	t.True(tdhttp.CmpJSONResponse(mt, req, handler,
		tdhttp.Response{
			Status: 242,
			Body:   JResp{Name: op},
		}))

	// Second one should fail, as previously anchored operator has been reset
	t.False(tdhttp.CmpJSONResponse(mt, req, handler,
		tdhttp.Response{
			Status: 242,
			Body:   JResp{Name: op},
		}))
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
				`~ Failed test 'body unmarshaling'
\s+unmarshal\(Response\.Body\): should NOT be an error
\s+got: .*unknown type func\(\).*
\s+expected: nil`,
				`~ Received response:
\s+\x60(?s:.+?)
\s+
\s+<XResp><name>Bob</name></XResp>
\s+\x60
`, // check the complete body is shown
			},
		},
	} {
		t.Run(curTest.Name,
			func(t *td.T) {
				testCmpResponse(t, tdhttp.CmpXMLResponse, "CmpXMLResponse", curTest)
			})

		t.Run(curTest.Name+" TestAPI",
			func(t *td.T) {
				testTestAPI(t, (*tdhttp.TestAPI).CmpXMLBody, "CmpXMLBody", curTest)
			})
	}
}

var logsViz = strings.NewReplacer(
	" ", "·",
	"\t", "→",
	"\r", "<cr>",
)

func testLogs(t *td.T, mockT *tdutil.T, curTest CmpResponseTest) {
	t.Helper()

	dumpLogs := !t.Cmp(mockT.Failed(), !curTest.Success, "test failure")

	for _, expectedLog := range curTest.ExpectedLogs {
		if strings.HasPrefix(expectedLog, "~") {
			re := regexp.MustCompile(expectedLog[1:])
			if !re.MatchString(mockT.LogBuf()) {
				t.Errorf(`logs do not match "%s" regexp`, re)
				dumpLogs = true
			}
		} else if !strings.Contains(mockT.LogBuf(), expectedLog) {
			t.Errorf(`"%s" not found in test logs`, expectedLog)
			dumpLogs = true
		}
	}

	if dumpLogs {
		t.Errorf(`Test logs: "%s"`, logsViz.Replace(mockT.LogBuf()))
	}
}

func testCmpResponse(t *td.T,
	cmp func(testing.TB, *http.Request, func(http.ResponseWriter, *http.Request), tdhttp.Response, ...any) bool,
	cmpName string,
	curTest CmpResponseTest,
) {
	t.Helper()

	mockT := tdutil.NewT(cmpName)

	t.Cmp(cmp(mockT,
		httptest.NewRequest("GET", "/path", nil),
		curTest.Handler,
		curTest.ExpectedResp),
		curTest.Success)

	testLogs(t, mockT, curTest)
}

func testTestAPI(t *td.T,
	cmpBody func(*tdhttp.TestAPI, any) *tdhttp.TestAPI,
	cmpName string,
	curTest CmpResponseTest,
) {
	t.Helper()

	mockT := tdutil.NewT(cmpName)

	ta := tdhttp.NewTestAPI(mockT, http.HandlerFunc(curTest.Handler)).
		Get("/path")

	if curTest.ExpectedResp.Status != nil {
		ta.CmpStatus(curTest.ExpectedResp.Status)
	}

	if curTest.ExpectedResp.Header != nil {
		ta.CmpHeader(curTest.ExpectedResp.Header)
	}

	if curTest.ExpectedResp.Cookies != nil {
		ta.CmpCookies(curTest.ExpectedResp.Cookies)
	}

	cmpBody(ta, curTest.ExpectedResp.Body)

	t.Cmp(ta.Failed(), !curTest.Success)

	testLogs(t, mockT, curTest)
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
			func(body []byte, target any) error {
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

			mockT := tdutil.NewT("zeroed_body")
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
			t.Contains(mockT.LogBuf(), "Received response:")
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

			mockT := tdutil.NewT("Unmarshal_is_failing")
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
				"You can try All(Isa(EXPECTED_TYPE), Any(…)) to disambiguate…")
			t.Contains(mockT.LogBuf(), "Received response:")
		})
	})
}
