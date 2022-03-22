// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/template"

	"github.com/maxatome/go-testdeep/helpers/tdhttp/internal"
)

func getResp(ta *TestAPI, name ...string) (*internal.Response, error) {
	var resp *internal.Response
	if len(name) == 0 {
		resp = ta.responses.Last()
		if resp == nil {
			return nil, errors.New("no last response found")
		}
	} else {
		resp = ta.responses.Get(name[0])
		if resp == nil {
			return nil, fmt.Errorf("no response recorded as %q found", name[0])
		}
	}
	return resp, nil
}

func tmplFuncs(ta *TestAPI) template.FuncMap {
	return template.FuncMap{
		//
		// Functions working on previous responses
		//
		// jsonp "/json/pointer"        → works on last response
		// jsonp "/json/pointer" "name" → works on response recorded as "name"
		"jsonp": func(pointer string, name ...string) (any, error) {
			resp, err := getResp(ta, name...)
			if err != nil {
				return nil, err
			}
			val, _, err := ta.respJSONPointer(resp, pointer)
			return val, err
		},
		// json "name" → returns response recorded as "name" JSON decoded
		"json": func(name string) (any, error) {
			resp, err := getResp(ta, name)
			if err != nil {
				return nil, err
			}
			return resp.UnmarshalJSON()
		},
		// header "X-Header-Foo"        → works on last response
		// header "X-Header-Foo" "name" → works on response recorded as "name"
		"header": func(key string, name ...string) (string, error) {
			resp, err := getResp(ta, name...)
			if err != nil {
				return "", err
			}
			return resp.Response().Header.Get(key), nil
		},
		// trailer "X-Trailer-Foo"        → works on last response
		// trailer "X-Trailer-Foo" "name" → works on response recorded as "name"
		"trailer": func(key string, name ...string) (string, error) {
			resp, err := getResp(ta, name...)
			if err != nil {
				return "", err
			}
			return resp.Response().Trailer.Get(key), nil
		},
		//
		// Basic util functions
		//
		// toJson VAL → returns the JSON representation of VAL
		"toJson": func(val any) (string, error) {
			b, err := json.Marshal(val)
			return string(b), err
		},
		"quote": strconv.Quote,
		"sub":   func(a, b int) int { return a - b },
		"add":   func(a, b int) int { return a + b },
	}
}

type TemplateJSON interface {
	io.ReadCloser
	json.Marshaler
	fmt.Stringer
	Err() error
}

type tmplJSON struct {
	cache *bytes.Buffer
	err   error
}

func newTemplateJSON(ta *TestAPI, s string) TemplateJSON {
	last := ta.responses.Last()
	if last == nil {
		return &tmplJSON{err: errors.New("no last response found")}
	}

	jlast, err := last.UnmarshalJSON()
	if err != nil {
		return &tmplJSON{err: fmt.Errorf("last response is not JSON formatted: %s", err)}
	}

	tmpl, err := template.New("").Funcs(tmplFuncs(ta)).Parse(s)
	if err != nil {
		return &tmplJSON{err: fmt.Errorf("template parsing failed: %s", err)}
	}

	var cache bytes.Buffer
	err = tmpl.Execute(&cache, jlast)
	if err != nil {
		return &tmplJSON{err: fmt.Errorf("template execution failed: %s", err)}
	}
	return &tmplJSON{cache: &cache}
}

// Read implements [io.ReadCloser] interface.
func (t *tmplJSON) Read(p []byte) (n int, err error) {
	if t.err != nil {
		return 0, t.err
	}
	n, err = t.cache.Read(p)
	t.err = err
	return
}

// Close implements [io.ReadCloser] interface. It always returns nil here.
func (t *tmplJSON) Close() error {
	return nil
}

// MarshalJSON implements [json.Marshaler] interface.
func (t *tmplJSON) MarshalJSON() (b []byte, err error) {
	if t.err != nil {
		return nil, t.err
	}
	b, err = json.RawMessage(t.cache.Bytes()).MarshalJSON()
	t.err = err
	return
}

// String implements [fmt.Stringer] interface.
func (t *tmplJSON) String() string {
	if t.err != nil {
		return ""
	}
	return t.cache.String()
}

// Err returns an error if something got wrong during the construction
// of the instance, typically the last response unmarshalling, the
// template parsing or its execution, or during [MarshalJSON] or
// [Read] calls.
func (t *tmplJSON) Err() error {
	return t.err
}
