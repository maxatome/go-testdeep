// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package tdsuite adds tests suite feature to [go-testdeep] in a
// non-intrusive way, but easily and powerfully.
//
// A tests suite is a set of tests run sequentially that share some data.
//
// Some hooks can be set to be automatically called before the suite
// is run, before, after and/or between each test, and at the end of
// the suite.
//
// In addition, a test can discontinue the suite.
//
// Giving a suite using a MySuite type, the test methods have the form:
//
//	func (s *MySuite) TestXxx(t *td.T)
//	func (s *MySuite) TestXxx(assert, require *td.T)
//
// where Xxx does not start with a lowercase letter. Each test method
// is run in a subtest, the method name serves to identify the
// subtest.
//
// A test method can return a bool, as in:
//
//	func (s *MySuite) TestXxx(t *td.T) bool
//	func (s *MySuite) TestXxx(assert, require *td.T) bool
//
// in this case, returning false means discontinuing the suite without
// any error. Consider it as a skip feature.
//
// A test method can instead return an error, as in:
//
//	func (s *MySuite) TestXxx(t *td.T) error
//	func (s *MySuite) TestXxx(assert, require *td.T) error
//
// in this case, returning a non-nil error marks the test as having
// failed, logs the error and discontinues the suite.
//
// A test method can also return a tuple (bool, error), as in:
//
//	func (s *MySuite) TestXxx(t *td.T) (bool, error)
//	func (s *MySuite) TestXxx(assert, require *td.T) (bool, error)
//
// in this case, both returned values are independent. Returning a
// false boolean means discontinuing the suite while returning a
// non-nil error marks the test as having failed and logs the
// error. So:
//
//	Returning       do...
//	(false, nil)    continue the suite, do not log anything
//	(false, ERROR)  continue the suite, marks the test as failed & log ERROR
//	(true, nil)     discontinue the suite & log the discontinuation
//	(true, ERROR)   discontinue the suite & log the discontinuation, marks the
//	                test as failed & log ERROR
//
// Test methods are run in lexicographic order.
//
// # Very simple tests suite
//
// Used typically to group tests and benefit from already instanciated
// [*td.T] instances.
//
//	import (
//	  "testing"
//
//	  "github.com/maxatome/go-testdeep/td"
//	  "github.com/maxatome/go-testdeep/helpers/tdsuite"
//	)
//
//	type MySuite struct{}
//
//	func (s MySuite) TestDB(assert, require *td.T) {
//	  db, err := initDB()
//	  require.CmpNoError(err)
//	  assert.CmpNoError(db.Ping())
//	}
//
//	func (s MySuite) TestPerson(assert *td.T) {
//	  person := Getperson("Bob")
//	  assert.Cmp(person, Person{Name: "Bob", Age: 44})
//	}
//
//	// TestMySuite is the go test entry point.
//	func TestMySuite(t *testing.T) {
//	  tdsuite.Run(t, MySuite{})
//	}
//
// # Suite setup and other hooks
//
// In most cases, a suite is used for sharing information between
// tests. The type of the suite can implement several methods that are
// called before, after and/or between tests.
//
//	type SuiteDB struct{
//	  DB *sql.DB
//	}
//
//	// Setup is called once before any test runs.
//	func (s *SuiteDB) Setup(t *td.T) error {
//	  db, err := sql.Open(driver, dataSourceName)
//	  s.DB = db
//	  return err // automatically logged + failure if non-nil
//	}
//
//	// Destroy is called after all tests are run.
//	// Destroy is not called if Setup returned an error.
//	func (s *SuiteDB) Destroy(t *td.T) error {
//	  return s.DB.Close() // automatically logged + failure if non-nil
//	}
//
//	func (s *SuiteDB) TestPerson(assert, require *td.T) {
//	  person, err := GetPerson(s.DB, "Bob")
//	  require.CmpNoError(err)
//	  assert.Cmp(person, Person{Name: "Bob", Age: 44})
//	}
//
//	// TestMySuite is the go test entry point.
//	func TestSuiteDB(t *testing.T) {
//	  tdsuite.Run(t, &SuiteDB{})
//	}
//
// See documentation below for other possible hooks: [PreTest], [PostTest]
// and [BetweenTests].
//
// [go-testdeep]: https://go-testdeep.zetta.rocks/
package tdsuite
