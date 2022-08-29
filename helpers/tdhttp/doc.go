// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package tdhttp, from [go-testdeep], provides some functions to easily
// test HTTP handlers.
//
// Combined to [td] package it provides powerful testing features.
//
// # TestAPI
//
// The better way to test HTTP APIs using this package.
//
//	ta := tdhttp.NewTestAPI(t, mux)
//
//	ta.Get("/person/42", "Accept", "application/xml").
//	  CmpStatus(http.StatusOK).
//	  CmpHeader(td.ContainsKey("X-Custom-Header")).
//	  CmpCookie(td.SuperBagOf(td.Smuggle("Name", "cookie_session"))).
//	  CmpXMLBody(Person{
//	    ID:   ta.Anchor(td.NotZero(), uint64(0)).(uint64),
//	    Name: "Bob",
//	    Age:  26,
//	  })
//
//	ta.Get("/person/42", "Accept", "application/json").
//	  CmpStatus(http.StatusOK).
//	  CmpHeader(td.ContainsKey("X-Custom-Header")).
//	  CmpCookies(td.SuperBagOf(td.Struct(&http.Cookie{Name: "cookie_session"}, nil))).
//	  CmpJSONBody(td.JSON(`
//	{
//	  "id":   $1,
//	  "name": "Bob",
//	  "age":  26
//	}`,
//	    td.NotZero()))
//
// See the full example below.
//
// # Cmp…Response functions
//
// Historically, it was the only way to test HTTP APIs using
// this package.
//
//	ok := tdhttp.CmpJSONResponse(t,
//	  tdhttp.Get("/person/42"),
//	  myAPI.ServeHTTP,
//	  Response{
//	    Status:  http.StatusOK,
//	    Header:  td.ContainsKey("X-Custom-Header"),
//	    Cookies: td.SuperBagOf(td.Smuggle("Name", "cookie_session")),
//	    Body: Person{
//	      ID:   42,
//	      Name: "Bob",
//	      Age:  26,
//	    },
//	  },
//	  "/person/{id} route")
//
// It now uses [TestAPI] behind the scene. It is better to directly
// use [TestAPI] and its methods instead, as it is more flexible and
// readable.
//
// [go-testdeep]: https://go-testdeep.zetta.rocks/
package tdhttp
