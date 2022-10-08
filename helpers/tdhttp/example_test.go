// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func Example() {
	t := &testing.T{}

	// Our API handle Persons with 3 routes:
	// - POST /person
	// - GET /person/{personID}
	// - DELETE /person/{personID}

	// Person describes a person.
	type Person struct {
		ID        int64      `json:"id,omitempty" xml:"ID,omitempty"`
		Name      string     `json:"name" xml:"Name"`
		Age       int        `json:"age" xml:"Age"`
		CreatedAt *time.Time `json:"created_at,omitempty" xml:"CreatedAt,omitempty"`
	}

	// Error is returned to the client in case of error.
	type Error struct {
		Mesg string `json:"message" xml:"Message"`
		Code int    `json:"code" xml:"Code"`
	}

	// Our µDB :)
	var mu sync.Mutex
	personByID := map[int64]*Person{}
	personByName := map[string]*Person{}
	var lastID int64

	// reply is a helper to send responses.
	reply := func(w http.ResponseWriter, status int, contentType string, body any) {
		if body == nil {
			w.WriteHeader(status)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		switch contentType {
		case "application/json":
			json.NewEncoder(w).Encode(body) //nolint: errcheck
		case "application/xml":
			xml.NewEncoder(w).Encode(body) //nolint: errcheck
		default: // text/plain
			fmt.Fprintf(w, "%+v", body)
		}
	}

	// Our API
	mux := http.NewServeMux()

	// POST /person
	mux.HandleFunc("/person", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if req.Body == nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		var in Person
		var contentType string

		switch req.Header.Get("Content-Type") {
		case "application/json":
			err := json.NewDecoder(req.Body).Decode(&in)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(req.Body).Decode(&in)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		case "application/x-www-form-urlencoded":
			b, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			v, err := url.ParseQuery(string(b))
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			in.Name = v.Get("name")
			in.Age, err = strconv.Atoi(v.Get("age"))
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		contentType = req.Header.Get("Accept")

		if in.Name == "" || in.Age <= 0 {
			reply(w, http.StatusBadRequest, contentType, Error{
				Mesg: "Empty name or bad age",
				Code: http.StatusBadRequest,
			})
			return
		}

		mu.Lock()
		defer mu.Unlock()
		if personByName[in.Name] != nil {
			reply(w, http.StatusConflict, contentType, Error{
				Mesg: "Person already exists",
				Code: http.StatusConflict,
			})
			return
		}
		lastID++
		in.ID = lastID
		now := time.Now()
		in.CreatedAt = &now
		personByID[in.ID] = &in
		personByName[in.Name] = &in
		reply(w, http.StatusCreated, contentType, in)
	})

	// GET /person/{id}
	// DELETE /person/{id}
	mux.HandleFunc("/person/", func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(strings.TrimPrefix(req.URL.Path, "/person/"), 10, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		accept := req.Header.Get("Accept")

		mu.Lock()
		defer mu.Unlock()
		if personByID[id] == nil {
			reply(w, http.StatusNotFound, accept, Error{
				Mesg: "Person does not exist",
				Code: http.StatusNotFound,
			})
			return
		}

		switch req.Method {
		case http.MethodGet:
			reply(w, http.StatusOK, accept, personByID[id])
		case http.MethodDelete:
			delete(personByID, id)
			reply(w, http.StatusNoContent, "", nil)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	//
	// Let's test our API
	//
	ta := tdhttp.NewTestAPI(t, mux)

	// Re-usable custom operator to check Content-Type header
	contentTypeIs := func(ct string) td.TestDeep {
		return td.SuperMapOf(http.Header{"Content-Type": []string{ct}}, nil)
	}

	//
	// Person not found
	//
	ta.Get("/person/42", "Accept", "application/json").
		Name("GET /person/42 - JSON").
		CmpStatus(404).
		CmpHeader(contentTypeIs("application/json")).
		CmpJSONBody(Error{
			Mesg: "Person does not exist",
			Code: 404,
		})
	fmt.Println("GET /person/42 - JSON:", !ta.Failed())

	ta.Get("/person/42", "Accept", "application/xml").
		Name("GET /person/42 - XML").
		CmpStatus(404).
		CmpHeader(contentTypeIs("application/xml")).
		CmpXMLBody(Error{
			Mesg: "Person does not exist",
			Code: 404,
		})
	fmt.Println("GET /person/42 - XML:", !ta.Failed())

	ta.Get("/person/42", "Accept", "text/plain").
		Name("GET /person/42 - raw").
		CmpStatus(404).
		CmpHeader(contentTypeIs("text/plain")).
		CmpBody("{Mesg:Person does not exist Code:404}")
	fmt.Println("GET /person/42 - raw:", !ta.Failed())

	//
	// Create a Person
	//
	var bobID int64
	ta.PostXML("/person", Person{Name: "Bob", Age: 32},
		"Accept", "application/xml").
		Name("POST /person - XML").
		CmpStatus(201).
		CmpHeader(contentTypeIs("application/xml")).
		CmpXMLBody(Person{ // using operator anchoring directly in literal
			ID:        ta.A(td.Catch(&bobID, td.NotZero()), int64(0)).(int64),
			Name:      "Bob",
			Age:       32,
			CreatedAt: ta.A(td.Ptr(td.Between(ta.SentAt(), time.Now()))).(*time.Time),
		})
	fmt.Printf("POST /person - XML: %t → Bob ID=%d\n", !ta.Failed(), bobID)

	var aliceID int64
	ta.PostJSON("/person", Person{Name: "Alice", Age: 35},
		"Accept", "application/json").
		Name("POST /person - JSON").
		CmpStatus(201).
		CmpHeader(contentTypeIs("application/json")).
		CmpJSONBody(td.JSON(` // using JSON operator (yes comment allowed in JSON!)
{
  "id":         $1,
  "name":       "Alice",
  "age":        35,
  "created_at": $2
}`,
			td.Catch(&aliceID, td.NotZero()),
			td.Smuggle(func(date string) (time.Time, error) {
				return time.Parse(time.RFC3339Nano, date)
			}, td.Between(ta.SentAt(), time.Now()))))
	fmt.Printf("POST /person - JSON: %t → Alice ID=%d\n", !ta.Failed(), aliceID)

	var brittID int64
	ta.PostForm("/person",
		url.Values{
			"name": []string{"Britt"},
			"age":  []string{"29"},
		},
		"Accept", "text/plain").
		Name("POST /person - raw").
		CmpStatus(201).
		CmpHeader(contentTypeIs("text/plain")).
		// using Re (= Regexp) operator
		CmpBody(td.Re(`\{ID:(\d+) Name:Britt Age:29 CreatedAt:.*\}\z`,
			td.Smuggle(func(groups []string) (int64, error) {
				return strconv.ParseInt(groups[0], 10, 64)
			}, td.Catch(&brittID, td.NotZero()))))
	fmt.Printf("POST /person - raw: %t → Britt ID=%d\n", !ta.Failed(), brittID)

	//
	// Get a Person
	//
	ta.Get(fmt.Sprintf("/person/%d", aliceID), "Accept", "application/xml").
		Name("GET Alice - XML (ID #%d)", aliceID).
		CmpStatus(200).
		CmpHeader(contentTypeIs("application/xml")).
		CmpXMLBody(td.SStruct( // using SStruct operator
			Person{
				ID:   aliceID,
				Name: "Alice",
				Age:  35,
			},
			td.StructFields{
				"CreatedAt": td.Ptr(td.NotZero()),
			},
		))
	fmt.Println("GET XML Alice:", !ta.Failed())

	ta.Get(fmt.Sprintf("/person/%d", aliceID), "Accept", "application/json").
		Name("GET Alice - JSON (ID #%d)", aliceID).
		CmpStatus(200).
		CmpHeader(contentTypeIs("application/json")).
		CmpJSONBody(td.JSON(` // using JSON operator (yes comment allowed in JSON!)
{
  "id":         $1,
  "name":       "Alice",
  "age":        35,
  "created_at": $2
}`,
			aliceID,
			td.Not(td.Re(`^0001-01-01`)), // time is not 0001-01-01… aka zero time.Time
		))
	fmt.Println("GET JSON Alice:", !ta.Failed())

	//
	// Delete a Person
	//
	ta.Delete(fmt.Sprintf("/person/%d", aliceID), nil).
		Name("DELETE Alice (ID #%d)", aliceID).
		CmpStatus(204).
		CmpHeader(td.Not(td.ContainsKey("Content-Type"))).
		NoBody()
	fmt.Println("DELETE Alice:", !ta.Failed())

	// Check Alice is deleted
	ta.Get(fmt.Sprintf("/person/%d", aliceID), "Accept", "application/json").
		Name("GET (deleted) Alice - JSON (ID #%d)", aliceID).
		CmpStatus(404).
		CmpHeader(contentTypeIs("application/json")).
		CmpJSONBody(td.JSON(`
{
  "message": "Person does not exist",
  "code":    404
}`))
	fmt.Println("Alice is not found anymore:", !ta.Failed())

	// Output:
	// GET /person/42 - JSON: true
	// GET /person/42 - XML: true
	// GET /person/42 - raw: true
	// POST /person - XML: true → Bob ID=1
	// POST /person - JSON: true → Alice ID=2
	// POST /person - raw: true → Britt ID=3
	// GET XML Alice: true
	// GET JSON Alice: true
	// DELETE Alice: true
	// Alice is not found anymore: true
}
