// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
)

type Response struct {
	sync.Mutex

	name     string
	response *httptest.ResponseRecorder

	asJSON      any
	jsonDecoded bool
}

func NewResponse(resp *httptest.ResponseRecorder) *Response {
	return &Response{
		response: resp,
	}
}

func (r *Response) Response() *http.Response {
	// No lock needed here
	return r.response.Result()
}

func (r *Response) UnmarshalJSON() (any, error) {
	r.Lock()
	defer r.Unlock()

	if !r.jsonDecoded {
		err := json.Unmarshal(r.response.Body.Bytes(), &r.asJSON)
		if err != nil {
			return nil, err
		}
		r.jsonDecoded = true
	}

	return r.asJSON, nil
}
