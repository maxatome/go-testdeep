// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/td"
)

func server() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/any", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-TestDeep-Method", req.Method)
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "%s!", req.Method)
		if req.ContentLength != 0 {
			w.Write([]byte("\n---\n")) //nolint: errcheck
			io.Copy(w, req.Body)       //nolint: errcheck
		}
	})

	mux.HandleFunc("/any/json", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-TestDeep-Method", req.Method)
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		m := map[string]any{
			"method": req.Method,
		}
		if req.ContentLength != 0 {
			var body any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m["body"] = body
		}
		json.NewEncoder(w).Encode(m) //nolint: errcheck
	})

	mux.HandleFunc("/mirror/json", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-TestDeep-Method", req.Method)
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		io.Copy(w, req.Body) //nolint: errcheck
	})

	mux.HandleFunc("/hq/json", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m := map[string]any{
			"header":       req.Header,
			"query_params": req.URL.Query(),
		}
		json.NewEncoder(w).Encode(m) //nolint: errcheck
	})

	mux.HandleFunc("/any/xml", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-TestDeep-Method", req.Method)
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `<XResp><method>%s</method>`, req.Method)
		if req.ContentLength != 0 {
			io.Copy(w, req.Body) //nolint: errcheck
		}
		w.Write([]byte(`</XResp>`)) //nolint: errcheck
	})

	mux.HandleFunc("/any/cookies", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-TestDeep-Method", req.Method)
		if req.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "text/plain")

		http.SetCookie(w, &http.Cookie{
			Name:    "first",
			Value:   "cookie1",
			MaxAge:  123456,
			Expires: time.Date(2021, time.August, 12, 11, 22, 33, 0, time.UTC),
		})
		http.SetCookie(w, &http.Cookie{
			Name:   "second",
			Value:  "cookie2",
			MaxAge: 654321,
		})

		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "%s!", req.Method)
		if req.ContentLength != 0 {
			w.Write([]byte("\n---\n")) //nolint: errcheck
			io.Copy(w, req.Body)       //nolint: errcheck
		}
	})

	mux.HandleFunc("/any/trailer", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Trailer", "X-TestDeep-Method")
		w.Header().Add("Trailer", "X-TestDeep-Foo")

		io.WriteString(w, "Hey!") //nolint: errcheck

		w.Header().Set("X-TestDeep-Method", req.Method)
		w.Header().Set("X-TestDeep-Foo", "bar")
	})

	return mux
}

func TestNewTestAPI(t *testing.T) {
	mux := server()

	containsKey := td.ContainsKey("X-Testdeep-Method")

	t.Run("No error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpHeader(td.SuperMapOf(http.Header{}, td.MapEntries{
					"X-Testdeep-Method": td.Bag(td.Re(`(?i)^head\z`)),
				})).
				NoBody().
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.Empty())
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody(td.Empty()).
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.Empty())
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("GET!").
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.String("GET!"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody(td.Contains("GET")).
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.Contains("GET"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Options("/any", strings.NewReader("OPTIONS body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("OPTIONS!\n---\nOPTIONS body").
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.String("OPTIONS!\n---\nOPTIONS body"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Post("/any", strings.NewReader("POST body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("POST!\n---\nPOST body").
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.String("POST!\n---\nPOST body"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostForm("/any", url.Values{"p1": {"v1"}, "p2": {"v2"}}).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("POST!\n---\np1=v1&p2=v2").
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.String("POST!\n---\np1=v1&p2=v2"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostForm("/any", tdhttp.Q{"p1": "v1", "p2": "v2"}).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("POST!\n---\np1=v1&p2=v2").
				CmpResponse(td.Code(func(assert *td.T, resp *http.Response) {
					assert.Cmp(resp.StatusCode, 200)
					assert.Cmp(resp.Header, containsKey)
					assert.Smuggle(resp.Body, io.ReadAll, td.String("POST!\n---\np1=v1&p2=v2"))
				})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostMultipartFormData("/any", &tdhttp.MultipartBody{
					Boundary: "BoUnDaRy",
					Parts: []*tdhttp.MultipartPart{
						tdhttp.NewMultipartPartString("pipo", "bingo"),
					},
				}).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody(strings.ReplaceAll(
					`POST!
---
--BoUnDaRy%CR
Content-Disposition: form-data; name="pipo"%CR
Content-Type: text/plain; charset=utf-8%CR
%CR
bingo%CR
--BoUnDaRy--%CR
`,
					"%CR", "\r")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Put("/any", strings.NewReader("PUT body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("PUT!\n---\nPUT body").
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Patch("/any", strings.NewReader("PATCH body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("PATCH!\n---\nPATCH body").
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Delete("/any", strings.NewReader("DELETE body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("DELETE!\n---\nDELETE body").
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())
	})

	t.Run("No JSON error", func(t *testing.T) {
		requestBody := map[string]any{"hey": 123}
		expectedBody := func(m string) td.TestDeep {
			return td.JSON(`{"method": $1, "body": {"hey": 123}}`, m)
		}

		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				NewJSONRequest("GET", "/mirror/json", json.RawMessage(`null`)).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(nil).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				NewJSONRequest("ZIP", "/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(expectedBody("ZIP")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				NewJSONRequest("ZIP", "/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(td.JSONPointer("/body/hey", 123)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostJSON("/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(expectedBody("POST")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PutJSON("/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(expectedBody("PUT")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PatchJSON("/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(expectedBody("PATCH")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				DeleteJSON("/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(expectedBody("DELETE")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// With anchors
		type ReqBody struct {
			Hey int `json:"hey"`
		}
		type Resp struct {
			Method  string  `json:"method"`
			ReqBody ReqBody `json:"body"`
		}
		mockT = tdutil.NewT("test")
		tt := td.NewT(mockT)
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				DeleteJSON("/any/json", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(Resp{
					Method: tt.A(td.Re(`^(?i)delete\z`), "").(string),
					ReqBody: ReqBody{
						Hey: tt.A(td.Between(120, 130)).(int),
					},
				}).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// JSON and root operator (here SuperMapOf)
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostJSON("/any/json", true).
				CmpStatus(200).
				CmpJSONBody(td.JSON(`SuperMapOf({"body":Ignore()})`)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// td.Bag+td.JSON
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostJSON("/mirror/json",
					json.RawMessage(`[{"name":"Bob"},{"name":"Alice"}]`)).
				CmpStatus(200).
				CmpJSONBody(td.Bag(
					td.JSON(`{"name":"Alice"}`),
					td.JSON(`{"name":"Bob"}`),
				)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// td.Bag+literal
		type People struct {
			Name string `json:"name"`
		}
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostJSON("/mirror/json",
					json.RawMessage(`[{"name":"Bob"},{"name":"Alice"}]`)).
				CmpStatus(200).
				CmpJSONBody(td.Bag(People{"Alice"}, People{"Bob"})).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())
	})

	t.Run("No XML error", func(t *testing.T) {
		type XBody struct {
			Hey int `xml:"hey"`
		}
		type XResp struct {
			Method  string `xml:"method"`
			ReqBody *XBody `xml:"XBody"`
		}

		requestBody := XBody{Hey: 123}
		expectedBody := func(m string) XResp {
			return XResp{
				Method:  m,
				ReqBody: &requestBody,
			}
		}

		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				NewXMLRequest("ZIP", "/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(expectedBody("ZIP")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostXML("/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(expectedBody("POST")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PutXML("/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(expectedBody("PUT")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PatchXML("/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(expectedBody("PATCH")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				DeleteXML("/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(expectedBody("DELETE")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// With anchors
		mockT = tdutil.NewT("test")
		tt := td.NewT(mockT)
		td.CmpFalse(tt,
			tdhttp.NewTestAPI(mockT, mux).
				DeleteXML("/any/xml", requestBody).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(XResp{
					Method: tt.A(td.Re(`^(?i)delete\z`), "").(string),
					ReqBody: &XBody{
						Hey: tt.A(td.Between(120, 130)).(int),
					},
				}).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())
	})

	t.Run("Cookies", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies([]*http.Cookie{
					{
						Name:    "first",
						Value:   "cookie1",
						MaxAge:  123456,
						Expires: time.Date(2021, time.August, 12, 11, 22, 33, 0, time.UTC),
					},
					{
						Name:   "second",
						Value:  "cookie2",
						MaxAge: 654321,
					},
				}).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies([]*http.Cookie{
					{
						Name:    "first",
						Value:   "cookie1",
						MaxAge:  123456,
						Expires: time.Date(2021, time.August, 12, 11, 22, 33, 0, time.UTC),
					},
				}).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'cookies should match'")
		td.CmpContains(t, mockT.LogBuf(),
			"Response.Cookie: comparing slices, from index #1")

		// 2 cookies are here whatever their order is using Bag
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies(td.Bag(
					td.Smuggle("Name", "second"),
					td.Smuggle("Name", "first"),
				)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// Testing only Name & Value whatever their order is using Bag
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies(td.Bag(
					td.Struct(&http.Cookie{Name: "first", Value: "cookie1"}, nil),
					td.Struct(&http.Cookie{Name: "second", Value: "cookie2"}, nil),
				)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// Testing the presence of only one using SuperBagOf
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies(td.SuperBagOf(
					td.Struct(&http.Cookie{Name: "first", Value: "cookie1"}, nil),
				)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// Testing only the number of cookies
		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/cookies").
				CmpCookies(td.Len(2)).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpCookies(td.Len(100)). // fails
				CmpCookies(td.Len(2)).   // succeeds
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'cookies should match'")

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any/cookies").
				Name("my test").
				CmpCookies(td.Len(100)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: cookies should match'")
		td.CmpContains(t, mockT.LogBuf(), "Response.Cookie: bad length")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// Request not sent
		mockT = tdutil.NewT("test")
		ta := tdhttp.NewTestAPI(mockT, mux).
			Name("my test").
			CmpCookies(td.Len(2))
		td.CmpTrue(t, ta.Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("No response received yet\n"))
	})

	t.Run("Trailer", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(200).
				CmpTrailer(nil). // No trailer at all
				Failed())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/trailer").
				CmpStatus(200).
				CmpTrailer(containsKey).
				Failed())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/trailer").
				CmpStatus(200).
				CmpTrailer(http.Header{
					"X-Testdeep-Method": {"GET"},
					"X-Testdeep-Foo":    {"bar"},
				}).
				Failed())

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any/trailer").
				Name("my test").
				CmpTrailer(http.Header{}).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: trailer should match'")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// OrDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/trailer").
				Name("my test").
				CmpTrailer(http.Header{}).
				OrDumpResponse().
				OrDumpResponse(). // only one log
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: trailer should match'")
		logPos := strings.Index(mockT.LogBuf(), "Received response:\n")
		if td.Cmp(t, logPos, td.Gte(0)) {
			// Only one occurrence
			td.Cmp(t,
				strings.Index(mockT.LogBuf()[logPos+1:], "Received response:\n"),
				-1)
		}

		mockT = tdutil.NewT("test")
		ta := tdhttp.NewTestAPI(mockT, mux).
			Name("my test").
			CmpTrailer(http.Header{})
		td.CmpTrue(t, ta.Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("No response received yet\n"))

		end := len(mockT.LogBuf())
		ta.OrDumpResponse()
		td.CmpContains(t, mockT.LogBuf()[end:], "No response received yet\n")
	})

	t.Run("Status error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(400).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'status code should match'")

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(400). // fails
				CmpStatus(200). // succeeds
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'status code should match'")

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpStatus(400).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: status code should match'")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any").
				Name("my test").
				CmpStatus(400).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: status code should match'")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// OrDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpStatus(400).
				OrDumpResponse().
				OrDumpResponse(). // only one log
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: status code should match'")
		logPos := strings.Index(mockT.LogBuf(), "Received response:\n")
		if td.Cmp(t, logPos, td.Gte(0)) {
			// Only one occurrence
			td.Cmp(t,
				strings.Index(mockT.LogBuf()[logPos+1:], "Received response:\n"),
				-1)
		}

		mockT = tdutil.NewT("test")
		ta := tdhttp.NewTestAPI(mockT, mux).
			Name("my test").
			CmpStatus(400)
		td.CmpTrue(t, ta.Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("No response received yet\n"))

		end := len(mockT.LogBuf())
		ta.OrDumpResponse()
		td.CmpContains(t, mockT.LogBuf()[end:], "No response received yet\n")
	})

	t.Run("Header error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'header should match'")

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpHeader(td.Not(containsKey)). // fails
				CmpHeader(td.Ignore()).         // succeeds
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'header should match'")

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: header should match'")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any").
				Name("my test").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: header should match'")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
	})

	t.Run("Body error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'body contents is OK'")
		td.CmpContains(t, mockT.LogBuf(), "Response.Body: values differ\n")
		td.CmpContains(t, mockT.LogBuf(), `expected: "xxx"`)
		td.CmpContains(t, mockT.LogBuf(), `got: "GET!"`)

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpBody("xxx").       // fails
				CmpBody(td.Ignore()). // succeeds
				Failed())

		// Without AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: body contents is OK'")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any").
				Name("my test").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: body contents is OK'")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")

		// NoBody
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				NoBody().
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				Head("/any").
				CmpBody("fail"). // fails
				NoBody().        // succeeds
				Failed())

		// No JSON body
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(json.RawMessage(`{}`)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'body should not be empty'")
		td.CmpContains(t, mockT.LogBuf(), "Response body is empty!")
		td.CmpContains(t, mockT.LogBuf(), "Body cannot be empty when using CmpJSONBody")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any/json").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(json.RawMessage(`{}`)). // fails
				CmpJSONBody(td.Ignore()).           // succeeds
				Failed())

		// No JSON body + AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpJSONBody(json.RawMessage(`{}`)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'body should not be empty'")
		td.CmpContains(t, mockT.LogBuf(), "Response body is empty!")
		td.CmpContains(t, mockT.LogBuf(), "Body cannot be empty when using CmpJSONBody")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// No XML body
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpXMLBody(struct{ Test string }{}).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'body should not be empty'")
		td.CmpContains(t, mockT.LogBuf(), "Response body is empty!")
		td.CmpContains(t, mockT.LogBuf(), "Body cannot be empty when using CmpXMLBody")
	})

	t.Run("Response error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpResponse(nil).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'full response should match'")
		td.CmpContains(t, mockT.LogBuf(), "Response: values differ")
		td.CmpContains(t, mockT.LogBuf(), "got: (*http.Response)(")
		td.CmpContains(t, mockT.LogBuf(), "expected: nil")

		// Error followed by a success: Failed() should return true anyway
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpResponse(nil).         // fails
				CmpResponse(td.Ignore()). // succeeds
				Failed())

		// Without AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpResponse(nil).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: full response should match'")
		td.CmpNot(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		// AutoDumpResponse
		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				AutoDumpResponse().
				Get("/any").
				Name("my test").
				CmpResponse(nil).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: full response should match'")
		td.Cmp(t, mockT.LogBuf(), td.Contains("Received response:\n"))

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpResponse(nil).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header, body or full response\n")
	})

	t.Run("Request error", func(t *testing.T) {
		var ta *tdhttp.TestAPI
		checkFatal := func(fn func()) {
			mockT := tdutil.NewT("test")
			td.CmpTrue(t, mockT.CatchFailNow(func() {
				ta = tdhttp.NewTestAPI(mockT, mux)
				fn()
			}))
			td.Cmp(t,
				mockT.LogBuf(),
				td.Contains("newRequestParams... can only contains string, http.Header, ([]*|*|)http.Cookie, url.Values, tdhttp.Q and func(*http.Request) error, not bool"),
			)
		}

		empty := strings.NewReader("")

		checkFatal(func() { ta.Get("/path", true) })
		checkFatal(func() { ta.Head("/path", true) })
		checkFatal(func() { ta.Options("/path", empty, true) })
		checkFatal(func() { ta.Post("/path", empty, true) })
		checkFatal(func() { ta.PostForm("/path", nil, true) })
		checkFatal(func() { ta.PostMultipartFormData("/path", &tdhttp.MultipartBody{}, true) })
		checkFatal(func() { ta.Put("/path", empty, true) })
		checkFatal(func() { ta.Patch("/path", empty, true) })
		checkFatal(func() { ta.Delete("/path", empty, true) })

		checkFatal(func() { ta.NewJSONRequest("ZIP", "/path", nil, true) })
		checkFatal(func() { ta.PostJSON("/path", nil, true) })
		checkFatal(func() { ta.PutJSON("/path", nil, true) })
		checkFatal(func() { ta.PatchJSON("/path", nil, true) })
		checkFatal(func() { ta.DeleteJSON("/path", nil, true) })

		checkFatal(func() { ta.NewXMLRequest("ZIP", "/path", nil, true) })
		checkFatal(func() { ta.PostXML("/path", nil, true) })
		checkFatal(func() { ta.PutXML("/path", nil, true) })
		checkFatal(func() { ta.PatchXML("/path", nil, true) })
		checkFatal(func() { ta.DeleteXML("/path", nil, true) })
	})

	t.Run("Request params", func(t *testing.T) {
		ta := tdhttp.NewTestAPI(t, mux)

		ta.Get("/hq/json",
			"X-Test", "pipo",
			tdhttp.Q{"a": "b"},
			http.Cookie{Name: "cook1", Value: "val1"}).
			CmpJSONBody(td.JSON(`{
				"header": SuperMapOf({
					"X-Test": ["pipo"]
				}),
				"query_params": {
					"a": ["b"]
				}
			}`))

		ta.DefaultRequestParams(
			"X-Zip", "test",
			tdhttp.Q{"x": "y"},
			http.Cookie{Name: "cook9", Value: "val9"},
			func(r *http.Request) error {
				r.Header.Add("X-Hook", "1")
				return nil
			},
			func(r *http.Request) error {
				r.Header.Add("X-Hook", "2")
				return nil
			})

		ta.Get("/hq/json").
			CmpJSONBody(td.JSON(`{
				"header": SuperMapOf({
					"X-Hook":  ["1", "2"],
					"X-Zip":  ["test"],
					"Cookie": ["cook9=val9"]
				}),
				"query_params": {
					"x": ["y"]
				}
			}`))

		ta.Get("/hq/json",
			"X-Test", "pipo",
			tdhttp.Q{"a": "b"},
			http.Cookie{Name: "cook1", Value: "val1"}).
			CmpJSONBody(td.JSON(`{
				"header": SuperMapOf({
					"X-Hook":  ["1", "2"],
					"X-Test": ["pipo"],
					"X-Zip":  ["test"],
					"Cookie": ["cook1=val1; cook9=val9"]
				}),
				"query_params": {
					"x": ["y"],
					"a": ["b"]
				}
			}`))

		ta.Get("/hq/json",
			"X-Zip", "override",
			tdhttp.Q{"x": "override"},
			http.Cookie{Name: "cook9", Value: "override"},
			func(r *http.Request) error {
				r.Header.Add("X-Hook", "3")
				return nil
			}).
			CmpJSONBody(td.JSON(`{
				"header": SuperMapOf({
					"X-Hook":  ["3", "1", "2"],
					"X-Zip":  ["override"],
					"Cookie": ["cook9=override"]
				}),
				"query_params": {
					"x": ["override"]
				}
			}`))

		t.Run("hook failures", func(t *testing.T) {
			tt := tdutil.NewT("test")
			ta := tdhttp.NewTestAPI(tt, mux)

			hook := func(r *http.Request) error {
				return fmt.Errorf("hook failure")
			}

			if td.CmpTrue(t, tt.CatchFailNow(func() { ta.Get("/hq/json", hook) })) {
				td.CmpContains(t, tt.LogBuf(), "hook failed: hook failure")
			}

			ta = ta.DefaultRequestParams(hook)
			if td.CmpTrue(t, tt.CatchFailNow(func() { ta.Get("/hq/json") })) {
				td.CmpContains(t, tt.LogBuf(), "hook failed: hook failure")
			}
		})
	})
}

func TestWith(t *testing.T) {
	mux := server()

	ta := tdhttp.NewTestAPI(tdutil.NewT("test1"), mux)

	td.CmpFalse(t, ta.Head("/any").CmpStatus(200).Failed())

	nt := tdutil.NewT("test2")

	nta := ta.With(nt)

	td.Cmp(t, nta.T(), td.Not(td.Shallow(ta.T())))

	td.CmpTrue(t, nta.CmpStatus(200).Failed()) // as no request sent yet
	td.CmpContains(t, nt.LogBuf(),
		"A request must be sent before testing status, header, body or full response")

	td.CmpFalse(t, ta.CmpStatus(200).Failed()) // request already sent, so OK

	nt = tdutil.NewT("test3")
	nta = ta.With(nt)

	td.CmpTrue(t, nta.Head("/any").
		CmpStatus(400).
		OrDumpResponse().
		Failed())
	td.CmpContains(t, nt.LogBuf(), "Response.Status: values differ")
	td.CmpContains(t, nt.LogBuf(), "X-Testdeep-Method: HEAD") // Header dumped
}

func TestClone(t *testing.T) {
	mux := server()

	ta := tdhttp.NewTestAPI(t, mux)
	nta := ta.Clone().DefaultRequestParams(
		"X-Test", "bingo", "X-Zip", "OK",
		tdhttp.Q{"a": "pizza", "b": "kiss"},
		http.Cookie{Name: "cook1", Value: "VAL1"},
	)

	td.Cmp(t, nta.T(), td.Shallow(ta.T()), "ta and nta share the same testing.TB")
	td.Cmp(t, nta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"bingo"}, "X-Zip": {"OK"}},
		"defaultQParams": url.Values{"a": {"pizza"}, "b": {"kiss"}},
		"defaultCookies": []*http.Cookie{
			{Name: "cook1", Value: "VAL1"},
		},
	}), "nta modified successfully")
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  nil,
		"defaultQParams": nil,
		"defaultCookies": nil,
	}), "ta was not modified")
}

func TestDefaultRequestParams(t *testing.T) {
	mux := server()

	ta := tdhttp.NewTestAPI(tdutil.NewT("test1"), mux)

	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  nil,
		"defaultQParams": nil,
		"defaultCookies": nil,
	}))

	ta.DefaultRequestParams("X-Test", "pipo")
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"pipo"}},
		"defaultQParams": nil,
		"defaultCookies": nil,
	}))

	ta.DefaultRequestParams(tdhttp.Q{"a": "zip"})
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  nil,
		"defaultQParams": url.Values{"a": {"zip"}},
		"defaultCookies": nil,
	}))

	ta.DefaultRequestParams(http.Cookie{Name: "cook1", Value: "val1"})
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  nil,
		"defaultQParams": nil,
		"defaultCookies": []*http.Cookie{{Name: "cook1", Value: "val1"}},
	}))

	ta.DefaultRequestParams()
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  nil,
		"defaultQParams": nil,
		"defaultCookies": nil,
	}))

	ta.DefaultRequestParams(
		"X-Test", "pipo",
		tdhttp.Q{"a": "zip"},
		http.Cookie{Name: "cook1", Value: "val1"})
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"pipo"}},
		"defaultQParams": url.Values{"a": {"zip"}},
		"defaultCookies": []*http.Cookie{{Name: "cook1", Value: "val1"}},
	}))

	ta.AddDefaultRequestParams()
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"pipo"}},
		"defaultQParams": url.Values{"a": {"zip"}},
		"defaultCookies": []*http.Cookie{{Name: "cook1", Value: "val1"}},
	}))

	ta.AddDefaultRequestParams(
		"X-Zip", "OK",
		tdhttp.Q{"b": "kiss"},
		http.Cookie{Name: "cook2", Value: "val2"})
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"pipo"}, "X-Zip": {"OK"}},
		"defaultQParams": url.Values{"a": {"zip"}, "b": {"kiss"}},
		"defaultCookies": []*http.Cookie{
			{Name: "cook2", Value: "val2"},
			{Name: "cook1", Value: "val1"},
		},
	}))

	ta.AddDefaultRequestParams(
		"X-Test", "bingo", "X-Zip", "OK",
		tdhttp.Q{"a": "pizza", "b": "kiss"},
		http.Cookie{Name: "cook1", Value: "VAL1"},
		http.Cookie{Name: "cook2", Value: "val2"})
	td.Cmp(t, ta, td.Struct(nil, td.StructFields{
		"defaultHeader":  http.Header{"X-Test": {"bingo"}, "X-Zip": {"OK"}},
		"defaultQParams": url.Values{"a": {"pizza"}, "b": {"kiss"}},
		"defaultCookies": []*http.Cookie{
			{Name: "cook1", Value: "VAL1"},
			{Name: "cook2", Value: "val2"},
		},
	}))

	t.Run("DefaultRequestParams fatal", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t, mockT.CatchFailNow(func() {
			ta = tdhttp.NewTestAPI(mockT, mux)
			ta.DefaultRequestParams(123)
		}))
		td.Cmp(t,
			mockT.LogBuf(),
			td.Contains("newRequestParams... can only contains"))
	})

	t.Run("AddDefaultRequestParams fatal", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t, mockT.CatchFailNow(func() {
			ta = tdhttp.NewTestAPI(mockT, mux)
			ta.AddDefaultRequestParams(123)
		}))
		td.Cmp(t,
			mockT.LogBuf(),
			td.Contains("newRequestParams... can only contains"))
	})
}

func TestOr(t *testing.T) {
	mux := server()

	t.Run("Success", func(t *testing.T) {
		var orCalled bool
		for i, fn := range []any{
			func(body string) { orCalled = true },
			func(t *td.T, body string) { orCalled = true },
			func(body []byte) { orCalled = true },
			func(t *td.T, body []byte) { orCalled = true },
			func(t *td.T, r *httptest.ResponseRecorder) { orCalled = true },
		} {
			orCalled = false
			// As CmpStatus succeeds, Or function is not called
			td.CmpFalse(t,
				tdhttp.NewTestAPI(tdutil.NewT("test"), mux).
					Head("/any").
					CmpStatus(200).
					Or(fn).
					Failed(),
				"Not failed #%d", i)
			td.CmpFalse(t, orCalled, "called #%d", i)
		}
	})

	t.Run("No request sent", func(t *testing.T) {
		var ok, orCalled bool
		for i, fn := range []any{
			func(body string) { orCalled = true; ok = body == "" },
			func(t *td.T, body string) { orCalled = true; ok = t != nil && body == "" },
			func(body []byte) { orCalled = true; ok = body == nil },
			func(t *td.T, body []byte) { orCalled = true; ok = t != nil && body == nil },
			func(t *td.T, r *httptest.ResponseRecorder) { orCalled = true; ok = t != nil && r == nil },
		} {
			orCalled, ok = false, false
			// Check status without sending a request → fail
			td.CmpTrue(t,
				tdhttp.NewTestAPI(tdutil.NewT("test"), mux).
					CmpStatus(123).
					Or(fn).
					Failed(),
				"Failed #%d", i)
			td.CmpTrue(t, orCalled, "called #%d", i)
			td.CmpTrue(t, ok, "OK #%d", i)
		}
	})

	t.Run("Empty bodies", func(t *testing.T) {
		var ok, orCalled bool
		for i, fn := range []any{
			func(body string) { orCalled = true; ok = body == "" },
			func(t *td.T, body string) { orCalled = true; ok = t != nil && body == "" },
			func(body []byte) { orCalled = true; ok = body == nil },
			func(t *td.T, body []byte) { orCalled = true; ok = t != nil && body == nil },
			func(t *td.T, r *httptest.ResponseRecorder) {
				orCalled = true
				ok = t != nil && r != nil && r.Body.Len() == 0
			},
		} {
			orCalled, ok = false, false
			// HEAD /any = no body + CmpStatus fails
			td.CmpTrue(t,
				tdhttp.NewTestAPI(tdutil.NewT("test"), mux).
					Head("/any").
					CmpStatus(123).
					Or(fn).
					Failed(),
				"Failed #%d", i)
			td.CmpTrue(t, orCalled, "called #%d", i)
			td.CmpTrue(t, ok, "OK #%d", i)
		}
	})

	t.Run("Body", func(t *testing.T) {
		var ok, orCalled bool
		for i, fn := range []any{
			func(body string) { orCalled = true; ok = body == "GET!" },
			func(t *td.T, body string) { orCalled = true; ok = t != nil && body == "GET!" },
			func(body []byte) { orCalled = true; ok = string(body) == "GET!" },
			func(t *td.T, body []byte) { orCalled = true; ok = t != nil && string(body) == "GET!" },
			func(t *td.T, r *httptest.ResponseRecorder) {
				orCalled = true
				ok = t != nil && r != nil && r.Body.String() == "GET!"
			},
		} {
			orCalled, ok = false, false
			// GET /any = "GET!" body + CmpStatus fails
			td.CmpTrue(t,
				tdhttp.NewTestAPI(tdutil.NewT("test"), mux).
					Get("/any").
					CmpStatus(123).
					Or(fn).
					Failed(),
				"Failed #%d", i)
			td.CmpTrue(t, orCalled, "called #%d", i)
			td.CmpTrue(t, ok, "OK #%d", i)
		}
	})

	tt := tdutil.NewT("test")
	ta := tdhttp.NewTestAPI(tt, mux)
	if td.CmpTrue(t, tt.CatchFailNow(func() { ta.Or(123) })) {
		td.CmpContains(t, tt.LogBuf(),
			"usage: Or(func([*td.T,]string) | func([*td.T,][]byte) | func(*td.T,*httptest.ResponseRecorder)), but received int as 1st parameter")
	}
}

func TestRun(t *testing.T) {
	mux := server()

	ta := tdhttp.NewTestAPI(tdutil.NewT("test"), mux)

	ok := ta.Run("Test", func(ta *tdhttp.TestAPI) {
		td.CmpFalse(t, ta.Get("/any").CmpStatus(200).Failed())
	})
	td.CmpTrue(t, ok)

	ok = ta.Run("Test", func(ta *tdhttp.TestAPI) {
		td.CmpTrue(t, ta.Get("/any").CmpStatus(123).Failed())
	})
	td.CmpFalse(t, ok)
}
