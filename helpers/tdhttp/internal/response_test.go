// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp/internal"
	"github.com/maxatome/go-testdeep/td"
)

func newResponseRecorder(body string) *httptest.ResponseRecorder {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, body) //nolint: errcheck
	}

	rec := httptest.NewRecorder()
	handler(rec, &http.Request{})
	return rec
}

func TestResponse(t *testing.T) {
	assert, require := td.AssertRequire(t)

	r := internal.NewResponse(newResponseRecorder("12345"))
	require.NotNil(r)

	resp := r.Response()
	require.NotNil(resp)
	assert.Cmp(resp.StatusCode, http.StatusOK)
	assert.Cmp(resp.Body, td.Smuggle(io.ReadAll, td.String("12345")))

	v, err := r.UnmarshalJSON()
	assert.CmpNoError(err)
	assert.CmpLax(v, 12345)

	// Second call is cached
	v, err = r.UnmarshalJSON()
	assert.CmpNoError(err)
	assert.CmpLax(v, 12345)

	r = internal.NewResponse(newResponseRecorder("bad json"))
	require.NotNil(r)
	_, err = r.UnmarshalJSON()
	assert.CmpError(err)
}
