// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"testing"
)

// T is a type that encapsulates testing.T allowing to easily use
// testing.T methods as well as T ones.
type T struct {
	*testing.T
	Config ContextConfig // defaults to DefaultContextConfig
}

// NewT returns a new T instance. Typically used as:
//
//   import (
//     "testing"
//     td "github.com/maxatome/go-testdeep"
//   )
//
//   type Record struct {
//     Id        uint64
//     Name      string
//     Age       int
//     CreatedAt time.Time
//   }
//
//   func TestCreateRecord(tt *testing.T) {
//     t := NewT(tt, ContextConfig{
//       MaxErrors: 3, // in case of failure, will dump up to 3 errors
//     })
//
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if t.Nil(err) {
//       t.Log("No error, can now check struct contents")
//
//       ok := t.Struct(record,
//         &Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         td.StructFields{
//           Id:        td.Not(uint64(0)),
//           CreatedAt: td.Between(before, time.Now()),
//         },
//         "Newly created record")
//       if ok {
//         t.Log(Record created successfully!")
//       }
//     }
//   }
//
// "config" is an optional argument and, if passed, must be unique. It
// allows to configure how failures will be rendered during the life
// time of the returned instance.
//
//   t := NewT(tt)
//   t.CmpDeeply(
//     Record{Age: 12, Name: "Bob", Id: 12},  // got
//     Record{Age: 21, Name: "John", Id: 28}) // expected
//
// will produce:
//
//   === RUN   TestFoobar
//   --- FAIL: TestFoobar (0.00s)
//   	foobar_test.go:88: Failed test
//   		DATA.Id: values differ
//   			     got: (uint64) 12
//   			expected: (uint64) 28
//   FAIL
//
// Now with a special configuration:
//
//   t := NewT(tt, ContextConfig{
//       RootName:  "RECORD", // got data named "RECORD" instead of "DATA"
//       MaxErrors: 2,        // stops after 2 errors instead of 1
//     })
//   t.CmpDeeply(
//     Record{Age: 12, Name: "Bob", Id: 12},  // got
//     Record{Age: 21, Name: "John", Id: 28}) // expected
//
// will produce:
//
//   === RUN   TestFoobar
//   --- FAIL: TestFoobar (0.00s)
//   	foobar_test.go:96: Failed test
//   		RECORD.Id: values differ
//   			     got: (uint64) 12
//   			expected: (uint64) 28
//   		RECORD.Name: values differ
//   			     got: (string) (len=3) "Bob"
//   			expected: (string) (len=4) "John"
//   FAIL
//
// See RootName method to configure RootName in a more specific fashion.
//
// Note that setting MaxErrors to a negative value produces a dump
// with all errors.
//
// If MaxErrors is not set (or set to 0), it is set to
// DefaultContextConfig.MaxErrors which is potentially dependent from
// the TESTDEEP_MAX_ERRORS environment variable. See ContextConfig
// documentation for details.
func NewT(t *testing.T, config ...ContextConfig) *T {
	switch len(config) {
	case 0:
		return &T{
			T:      t,
			Config: DefaultContextConfig,
		}

	case 1:
		config[0].sanitize()
		return &T{
			T:      t,
			Config: config[0],
		}

	default:
		panic("usage: NewT(*testing.T[, ContextConfig]")
	}
}

// RootName changes the name of the got data. By default it is
// "DATA". For an HTTP response body, it could be "BODY" for example.
//
// It returns a new instance of *T so does not alter the original t
// and used as follows:
//
//   t.RootName("RECORD").
//     Struct(record,
//       &Record{
//         Name: "Bob",
//         Age:  23,
//       },
//       td.StructFields{
//         Id:        td.Not(uint64(0)),
//         CreatedAt: td.Between(before, time.Now()),
//       },
//       "Newly created record")
//
// In case of error for the field Age, the failure message will contain:
//
//   RECORD.Age: values differ
//
// Which is more readable than the generic:
//
//   DATA.Age: values differ
func (t *T) RootName(rootName string) *T {
	new := *t
	new.Config.RootName = rootName
	return &new
}

// CmpDeeply is mostly a shortcut for:
//
//   CmpDeeply(t.T, got, expected, args...)
//
// with the exception that t.Config is used to configure the test
// Context.
func (t *T) CmpDeeply(got, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return cmpDeeply(NewContextWithConfig(t.Config),
		t.T, got, expected, args...)
}

// True is shortcut for:
//
//   t.CmpDeeply(got, true, args...)
func (t *T) True(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, true, args...)
}

// False is shortcut for:
//
//   t.CmpDeeply(got, false, args...)
func (t *T) False(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, false, args...)
}
