// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package tdhttp

import "github.com/maxatome/go-testdeep/td"

// AnchorT returns a typed value allowing to anchor the [td.TestDeep]
// operator operator in a go classic literal like a struct, slice,
// array or map value.
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpJSONBody(Person{
//	    ID:   ta.AnchorT[uint64](td.NotZero()),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// See [td.T.AnchorT] for details.
//
// See [TestAPI.AT] method for a shorter synonym of AnchorT.
func (ta *TestAPI) AnchorT[X any](operator td.TestDeep) X {
	return ta.t.AnchorT[X](operator)
}

// AT is a synonym for [TestAPI.AnchorT]. It returns a typed value allowing to
// anchor the [td.TestDeep] operator in a go classic literal
// like a struct, slice, array or map value.
//
//	ta := tdhttp.NewTestAPI(tt, mux)
//
//	ta.Get("/person/42").
//	  CmpStatus(http.StatusOK).
//	  CmpJSONBody(Person{
//	    ID:   ta.AT[uint64](td.NotZero()),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
// See [td.T.AnchorT] for details.
func (ta *TestAPI) AT[X any](operator td.TestDeep) X {
	return ta.AnchorT[X](operator)
}
