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
	"net/url"
	"strings"
	"testing"

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

		m := map[string]interface{}{
			"method": req.Method,
		}
		if req.ContentLength != 0 {
			var body interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m["body"] = body
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

	return mux
}

func TestNewTestAPI(t *testing.T) {
	mux := server()

	containsKey := td.ContainsKey(td.Smuggle(http.CanonicalHeaderKey, "X-Testdeep-Method"))

	t.Run("No error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				NoBody().
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Head("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody(td.Empty()).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("GET!").
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody(td.Contains("GET")).
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				Post("/any", strings.NewReader("POST body")).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("POST!\n---\nPOST body").
				Failed())
		td.CmpEmpty(t, mockT.LogBuf())

		mockT = tdutil.NewT("test")
		td.CmpFalse(t,
			tdhttp.NewTestAPI(mockT, mux).
				PostForm("/any", url.Values{"p1": []string{"v1"}, "p2": []string{"v2"}}).
				CmpStatus(200).
				CmpHeader(containsKey).
				CmpBody("POST!\n---\np1=v1&p2=v2").
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
		requestBody := map[string]interface{}{"hey": 123}
		expectedBody := func(m string) td.TestDeep {
			return td.JSON(`{"method": $1, "body": {"hey": 123}}`, m)
		}

		mockT := tdutil.NewT("test")
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
		td.CmpFalse(tt,
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

	t.Run("Status error", func(t *testing.T) {
		mockT := tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				CmpStatus(400).
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

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpStatus(400).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header or body\n")
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

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(),
			"Failed test 'my test: header should match'")

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpHeader(td.Not(containsKey)).
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header or body\n")
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

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Get("/any").
				Name("my test").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: body contents is OK'")

		mockT = tdutil.NewT("test")
		td.CmpTrue(t,
			tdhttp.NewTestAPI(mockT, mux).
				Name("my test").
				CmpBody("xxx").
				Failed())
		td.CmpContains(t, mockT.LogBuf(), "Failed test 'my test: request is sent'\n")
		td.CmpContains(t, mockT.LogBuf(), "Request not sent!\n")
		td.CmpContains(t, mockT.LogBuf(), "A request must be sent before testing status, header or body\n")

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
}
