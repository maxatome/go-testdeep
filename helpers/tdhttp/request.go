// Copyright (c) 2019-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/types"
)

func newRequest(method string, target string, body io.Reader, headersQueryParams []any) (*http.Request, error) {
	header := http.Header{}
	qp := url.Values{}
	var cookies []*http.Cookie

	headersQueryParams = flat.Interfaces(headersQueryParams...)
	for i := 0; i < len(headersQueryParams); i++ {
		switch cur := headersQueryParams[i].(type) {
		case string:
			i++
			var val string
			if i < len(headersQueryParams) {
				var ok bool
				if val, ok = headersQueryParams[i].(string); !ok {
					return nil, errors.New(color.Bad(
						`header "%s" should have a string value, not a %T (@ headersQueryParams[%d])`,
						cur, headersQueryParams[i], i))
				}
			}
			header.Add(cur, val)

		case http.Header:
			for k, v := range cur {
				k = http.CanonicalHeaderKey(k)
				header[k] = append(header[k], v...)
			}

		case *http.Cookie:
			cookies = append(cookies, cur)

		case http.Cookie:
			cookies = append(cookies, &cur)

		case url.Values:
			for k, v := range cur {
				qp[k] = append(qp[k], v...)
			}

		case Q:
			err := cur.AddTo(qp)
			if err != nil {
				return nil, errors.New(color.Bad(
					"headersQueryParams... tdhttp.Q bad parameter: %s (@ headersQueryParams[%d])",
					err, i))
			}

		default:
			return nil, errors.New(color.Bad(
				"headersQueryParams... can only contains string, http.Header, http.Cookie, url.Values and tdhttp.Q, not %T (@ headersQueryParams[%d])",
				cur, i))
		}
	}

	// Parse path even when no query params to have consistent error
	// messages when using query params or not
	u, err := url.Parse(target)
	if err != nil {
		return nil, errors.New(color.Bad("target is not a valid path: %s", err))
	}
	if len(qp) > 0 {
		if u.RawQuery != "" {
			u.RawQuery += "&"
		}
		u.RawQuery += qp.Encode()
		target = u.String()
	}

	req := httptest.NewRequest(method, target, body)

	for k, v := range header {
		req.Header[k] = append(req.Header[k], v...)
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	return req, nil
}

// BasicAuthHeader returns a new [http.Header] with only Authorization
// key set, compliant with HTTP Basic Authentication using user and
// password. It is provided as a facility to build request in one
// line:
//
//	ta.Get("/path", tdhttp.BasicAuthHeader("max", "5ecr3T"))
//
// instead of:
//
//	req := tdhttp.Get("/path")
//	req.SetBasicAuth("max", "5ecr3T")
//	ta.Request(req)
//
// See [http.Request.SetBasicAuth] for details.
func BasicAuthHeader(user, password string) http.Header {
	return http.Header{
		"Authorization": []string{
			"Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password)),
		},
	}
}

func get(target string, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodGet, target, nil, headersQueryParams)
}

func head(target string, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodHead, target, nil, headersQueryParams)
}

func options(target string, body io.Reader, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodOptions, target, body, headersQueryParams)
}

func post(target string, body io.Reader, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodPost, target, body, headersQueryParams)
}

func postForm(target string, data URLValuesEncoder, headersQueryParams ...any) (*http.Request, error) {
	var body string
	if data != nil {
		body = data.Encode()
	}

	return newRequest(
		http.MethodPost, target, strings.NewReader(body),
		append(headersQueryParams, "Content-Type", "application/x-www-form-urlencoded"),
	)
}

func postMultipartFormData(target string, data *MultipartBody, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(
		http.MethodPost, target, data,
		append(headersQueryParams, "Content-Type", data.ContentType()),
	)
}

func put(target string, body io.Reader, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodPut, target, body, headersQueryParams)
}

func patch(target string, body io.Reader, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodPatch, target, body, headersQueryParams)
}

func del(target string, body io.Reader, headersQueryParams ...any) (*http.Request, error) {
	return newRequest(http.MethodDelete, target, body, headersQueryParams)
}

// NewRequest creates a new HTTP request as [httptest.NewRequest]
// does, with the ability to immediately add some headers and/or some
// query parameters.
//
// Headers can be added using string pairs as in:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  "Content-type", "application/pdf",
//	  "X-Test", "value",
//	)
//
// or using [http.Header] as in:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  http.Header{"Content-type": []string{"application/pdf"}},
//	)
//
// or using [BasicAuthHeader] as in:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  tdhttp.BasicAuthHeader("max", "5ecr3T"),
//	)
//
// or using [http.Cookie] (pointer or not, behind the scene,
// [http.Request.AddCookie] is used) as in:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  http.Cookie{Name: "cook1", Value: "val1"},
//	  &http.Cookie{Name: "cook2", Value: "val2"},
//	)
//
// Several header sources are combined:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  "Content-type", "application/pdf",
//	  http.Header{"X-Test": []string{"value1"}},
//	  "X-Test", "value2",
//	  http.Cookie{Name: "cook1", Value: "val1"},
//	  tdhttp.BasicAuthHeader("max", "5ecr3T"),
//	  &http.Cookie{Name: "cook2", Value: "val2"},
//	)
//
// Produces the following [http.Header]:
//
//	http.Header{
//	  "Authorization": []string{"Basic bWF4OjVlY3IzVA=="},
//	  "Content-type":  []string{"application/pdf"},
//	  "Cookie":        []string{"cook1=val1; cook2=val2"},
//	  "X-Test":        []string{"value1", "value2"},
//	}
//
// A string slice or a map can be flatened as well. As [NewRequest] expects
// ...any, [td.Flatten] can help here too:
//
//	strHeaders := map[string]string{
//	  "X-Length": "666",
//	  "X-Foo":    "bar",
//	}
//	req := tdhttp.NewRequest("POST", "/pdf", body, td.Flatten(strHeaders))
//
// Or combined with forms seen above:
//
//	req := tdhttp.NewRequest("POST", "/pdf", body,
//	  "Content-type", "application/pdf",
//	  http.Header{"X-Test": []string{"value1"}},
//	  td.Flatten(strHeaders),
//	  "X-Test", "value2",
//	  http.Cookie{Name: "cook1", Value: "val1"},
//	  tdhttp.BasicAuthHeader("max", "5ecr3T"),
//	  &http.Cookie{Name: "cook2", Value: "val2"},
//	)
//
// Header keys are always canonicalized using [http.CanonicalHeaderKey].
//
// Query parameters can be added using [url.Values] or more flexible
// [Q], as in:
//
//	req := tdhttp.NewRequest("GET", "/pdf",
//	  url.Values{
//	    "param": {"val"},
//	    "names": {"bob", "alice"},
//	  },
//	  "X-Test": "a header in the middle",
//	  tdhttp.Q{
//	    "limit":   20,
//	    "ids":     []int64{456, 789},
//	    "details": true,
//	  },
//	)
//
// All [url.Values] and [Q] instances are combined to produce the
// final query string to use. The previous example produces the
// following target:
//
//	/pdf?details=true&ids=456&ids=789&limit=20&names=bob&names=alice&param=val
//
// If target already contains a query string, it is reused:
//
//	req := tdhttp.NewRequest("GET", "/pdf?limit=10", tdhttp.Q{"details": true})
//
// produces the following target:
//
//	/path?details=true&limit=10
//
// Behind the scene, [url.Values.Encode] is used, so the parameters
// are always sorted by key. If you want a specific order, then do not
// use [url.Values] nor [Q] instances, but compose target by yourself.
//
// See [Q] documentation to learn how values are stringified.
func NewRequest(method, target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := newRequest(method, target, body, headersQueryParams)
	if err != nil {
		panic(err)
	}
	return req
}

// Get creates a new HTTP GET. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodGet, target, nil, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Get(target string, headersQueryParams ...any) *http.Request {
	req, err := get(target, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Head creates a new HTTP HEAD. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodHead, target, nil, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Head(target string, headersQueryParams ...any) *http.Request {
	req, err := head(target, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Options creates a HTTP OPTIONS. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodOptions, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Options(target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := options(target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Post creates a HTTP POST. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodPost, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Post(target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := post(target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// URLValuesEncoder is an interface [PostForm] and [TestAPI.PostForm] data
// must implement.
// Encode can be called to generate a "URL encoded" form such as
// ("bar=baz&foo=quux") sorted by key.
//
// [url.Values] and [Q] implement this interface.
type URLValuesEncoder interface {
	Encode() string
}

// PostForm creates a HTTP POST with data's keys and values
// URL-encoded as the request body. "Content-Type" header is
// automatically set to "application/x-www-form-urlencoded". Other
// headers can be added via headersQueryParams, as in:
//
//	req := tdhttp.PostForm("/data",
//	  url.Values{
//	    "param1": []string{"val1", "val2"},
//	    "param2": []string{"zip"},
//	  },
//	  "X-Foo", "Foo-value",
//	  "X-Zip", "Zip-value",
//	)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PostForm(target string, data URLValuesEncoder, headersQueryParams ...any) *http.Request {
	req, err := postForm(target, data, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PostMultipartFormData creates a HTTP POST multipart request, like
// multipart/form-data one for example. See [MultipartBody] type for
// details. "Content-Type" header is automatically set depending on
// data.MediaType (defaults to "multipart/form-data") and data.Boundary
// (defaults to "go-testdeep-42"). Other headers can be added via
// headersQueryParams, as in:
//
//	req := tdhttp.PostMultipartFormData("/data",
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
// and with a different media type:
//
//	req := tdhttp.PostMultipartFormData("/data",
//	  &tdhttp.MultipartBody{
//	    MediaType: "multipart/mixed",
//	    Parts:     []*tdhttp.MultipartPart{
//	      tdhttp.NewMultipartPartString("type", "Sales"),
//	      tdhttp.NewMultipartPartFile("report", "report.json", "application/json"),
//	    },
//	  },
//	  "X-Foo", "Foo-value",
//	  "X-Zip", "Zip-value",
//	)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PostMultipartFormData(target string, data *MultipartBody, headersQueryParams ...any) *http.Request {
	req, err := postMultipartFormData(target, data, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Put creates a HTTP PUT. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodPut, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Put(target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := put(target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Patch creates a HTTP PATCH. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodPatch, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Patch(target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := patch(target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// Delete creates a HTTP DELETE. It is a shortcut for:
//
//	tdhttp.NewRequest(http.MethodDelete, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func Delete(target string, body io.Reader, headersQueryParams ...any) *http.Request {
	req, err := del(target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

func newJSONRequest(method, target string, body any, headersQueryParams ...any) (*http.Request, error) {
	b, err := json.Marshal(body)
	if err != nil {
		if opErr, ok := types.AsOperatorNotJSONMarshallableError(err); ok {
			var plus string
			switch op := opErr.Operator(); op {
			case "JSON", "SubJSONOf", "SuperJSONOf":
				plus = ", use json.RawMessage() instead"
			}
			return nil, errors.New(color.Bad("JSON encoding failed: %s%s", err, plus))
		}
		return nil, errors.New(color.Bad("%s", err))
	}

	return newRequest(
		method, target, bytes.NewReader(b),
		append(headersQueryParams, "Content-Type", "application/json"),
	)
}

// NewJSONRequest creates a new HTTP request with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". Other headers can be added via headersQueryParams, as in:
//
//	req := tdhttp.NewJSONRequest("POST", "/data", body,
//	  "X-Foo", "Foo-value",
//	  "X-Zip", "Zip-value",
//	)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func NewJSONRequest(method, target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newJSONRequest(method, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PostJSON creates a HTTP POST with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". It is a shortcut for:
//
//	tdhttp.NewJSONRequest(http.MethodPost, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PostJSON(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newJSONRequest(http.MethodPost, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PutJSON creates a HTTP PUT with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". It is a shortcut for:
//
//	tdhttp.NewJSONRequest(http.MethodPut, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PutJSON(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newJSONRequest(http.MethodPut, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PatchJSON creates a HTTP PATCH with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". It is a shortcut for:
//
//	tdhttp.NewJSONRequest(http.MethodPatch, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PatchJSON(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newJSONRequest(http.MethodPatch, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// DeleteJSON creates a HTTP DELETE with body marshaled to
// JSON. "Content-Type" header is automatically set to
// "application/json". It is a shortcut for:
//
//	tdhttp.NewJSONRequest(http.MethodDelete, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func DeleteJSON(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newJSONRequest(http.MethodDelete, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

func newXMLRequest(method, target string, body any, headersQueryParams ...any) (*http.Request, error) {
	b, err := xml.Marshal(body)
	if err != nil {
		return nil, errors.New(color.Bad("XML encoding failed: %s", err))
	}

	return newRequest(
		method, target, bytes.NewReader(b),
		append(headersQueryParams, "Content-Type", "application/xml"),
	)
}

// NewXMLRequest creates a new HTTP request with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". Other headers can be added via headersQueryParams, as in:
//
//	req := tdhttp.NewXMLRequest("POST", "/data", body,
//	  "X-Foo", "Foo-value",
//	  "X-Zip", "Zip-value",
//	)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func NewXMLRequest(method, target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newXMLRequest(method, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PostXML creates a HTTP POST with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". It is a shortcut for:
//
//	tdhttp.NewXMLRequest(http.MethodPost, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PostXML(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newXMLRequest(http.MethodPost, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PutXML creates a HTTP PUT with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". It is a shortcut for:
//
//	tdhttp.NewXMLRequest(http.MethodPut, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PutXML(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newXMLRequest(http.MethodPut, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// PatchXML creates a HTTP PATCH with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". It is a shortcut for:
//
//	tdhttp.NewXMLRequest(http.MethodPatch, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func PatchXML(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newXMLRequest(http.MethodPatch, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}

// DeleteXML creates a HTTP DELETE with body marshaled to
// XML. "Content-Type" header is automatically set to
// "application/xml". It is a shortcut for:
//
//	tdhttp.NewXMLRequest(http.MethodDelete, target, body, headersQueryParams...)
//
// See [NewRequest] for all possible formats accepted in headersQueryParams.
func DeleteXML(target string, body any, headersQueryParams ...any) *http.Request {
	req, err := newXMLRequest(http.MethodDelete, target, body, headersQueryParams...)
	if err != nil {
		panic(err)
	}
	return req
}
