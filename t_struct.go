// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import "testing"

// T is a type that encapsulates *testing.T (in fact TestingFT
// interface which is implemented by *testing.T) allowing to easily
// use *testing.T methods as well as T ones.
type T struct {
	TestingFT
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
//     if t.CmpNoError(err) {
//       t.Log("No error, can now check struct contents")
//
//       ok := t.Struct(record,
//         &Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         td.StructFields{
//           Id:        td.NotZero(),
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
//   		DATA.Name: values differ
//   			     got: (string) (len=3) "Bob"
//   			expected: (string) (len=4) "John"
//   		DATA.Age: values differ
//   			     got: (int) 12
//   			expected: (int) 28
//   FAIL
//
// Now with a special configuration:
//
//   t := NewT(tt, ContextConfig{
//       RootName:  "RECORD", // got data named "RECORD" instead of "DATA"
//       MaxErrors: 2,        // stops after 2 errors instead of default 10
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
//   		Too many errors (use TESTDEEP_MAX_ERRORS=-1 to see all)
//   FAIL
//
// See RootName method to configure RootName in a more specific fashion.
//
// Note that setting MaxErrors to a negative value produces a dump
// with all errors.
//
// If MaxErrors is not set (or set to 0), it is set to
// DefaultContextConfig.MaxErrors which is potentially dependent from
// the TESTDEEP_MAX_ERRORS environment variable (else defaults to 10.)
// See ContextConfig documentation for details.
func NewT(t TestingFT, config ...ContextConfig) *T {
	switch len(config) {
	case 0:
		return &T{
			TestingFT: t,
			Config:    DefaultContextConfig,
		}

	case 1:
		config[0].sanitize()
		return &T{
			TestingFT: t,
			Config:    config[0],
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
//         Id:        td.NotZero(),
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
//   CmpDeeply(t.TestingFT, got, expected, args...)
//
// with the exception that t.Config is used to configure the test
// Context.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) CmpDeeply(got, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return cmpDeeply(NewContextWithConfig(t.Config),
		t.TestingFT, got, expected, args...)
}

// True is shortcut for:
//
//   t.CmpDeeply(got, true, args...)
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) True(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, true, args...)
}

// False is shortcut for:
//
//   t.CmpDeeply(got, false, args...)
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) False(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, false, args...)
}

// CmpError checks that "got" is non-nil error.
//
//   _, err := MyFunction(1, 2, 3)
//   t.CmpError(err, "MyFunction(1, 2, 3) should return an error")
//
// CmpError and not Error to avoid collision with t.TestingFT.Error method.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) CmpError(got error, args ...interface{}) bool {
	t.Helper()
	return cmpError(NewContextWithConfig(t.Config), t.TestingFT, got, args...)
}

// CmpNoError checks that "got" is nil error.
//
//   value, err := MyFunction(1, 2, 3)
//   if t.CmpNoError(err) {
//     // one can now check value...
//   }
//
// CmpNoError and not NoError to be consistent with CmpError method.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) CmpNoError(got error, args ...interface{}) bool {
	t.Helper()
	return cmpNoError(NewContextWithConfig(t.Config), t.TestingFT, got, args...)
}

// Run runs "f" as a subtest of t called "name". It runs "f" in a separate
// goroutine and blocks until "f" returns or calls t.Parallel to become
// a parallel test. Run reports whether "f" succeeded (or at least did
// not fail before calling t.Parallel).
//
// Run may be called simultaneously from multiple goroutines, but all
// such calls must return before the outer test function for t
// returns.
//
// Under the hood, Run delegates all this stuff to testing.Run. That
// is why this documentation is a copy/paste of testing.Run one.
func (t *T) Run(name string, f func(t *T)) bool {
	t.Helper()
	return t.TestingFT.Run(name, func(tt *testing.T) { f(NewT(tt)) })
}
