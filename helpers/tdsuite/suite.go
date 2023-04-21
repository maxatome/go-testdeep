// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdsuite

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/td"
)

var tType = reflect.TypeOf((*td.T)(nil))

// Setup is an interface a tests suite can implement. When running the
// tests suite, Setup method is called once before any test runs. If
// Setup returns an error, the tests suite aborts: no tests are run.
//
// t.Cleanup() can be called in Setup method. It can replace the
// definition of a [Destroy] method. It can also be used together, in
// this case cleanup registered functions are called after [Destroy].
type Setup interface {
	Setup(t *td.T) error
}

// PreTest is an interface a tests suite can implement. PreTest method
// is called before each test is run, in the same subtest as the test
// itself. If PreTest returns an error, the subtest aborts: the test
// is not run.
//
// t.Cleanup() can be called in PreTest method. It can replace the
// definition of a [PostTest] method. It can also be used together, in
// this case cleanup registered functions are called after [PostTest].
type PreTest interface {
	PreTest(t *td.T, testName string) error
}

// PostTest is an interface a tests suite can implement. PostTest
// method is called after each test is run, in the same subtest as
// the test itself. If [PreTest] interface is implemented and [PreTest]
// method returned an error, PostTest is never called.
type PostTest interface {
	PostTest(t *td.T, testName string) error
}

// BetweenTests is an interface a tests suite can
// implement. BetweenTests method is called between 2 tests. If it
// returns an error, the tests suite aborts: no more tests are run.
type BetweenTests interface {
	BetweenTests(t *td.T, previousTestName, nextTestName string) error
}

// Destroy is an interface a tests suite can implement. When running
// the tests suite, Destroy method is called once after all tests
// ran. If [Setup] interface is implemented and [Setup] method returned an
// error, Destroy is never called.
type Destroy interface {
	Destroy(t *td.T) error
}

func emptyPrePostTest(t *td.T, testName string) error    { return nil }
func emptyBetweenTests(t *td.T, prev, next string) error { return nil }

// isTest returns true if "name" is a valid test name.
// Derived from go sources in cmd/go/internal/load/test.go.
func isTest(name string) bool {
	if !strings.HasPrefix(name, "Test") {
		return false
	}
	if len(name) == 4 { // "Test" is ok
		return true
	}
	r, _ := utf8.DecodeRuneInString(name[4:])
	return !unicode.IsLower(r)
}

// shouldContinue returns true if the tests suite should continue
// based on ret, the value(s) returned by a test call.
func shouldContinue(t *td.T, testName string, ret []reflect.Value) bool {
	var (
		err  error
		cont bool
	)

	switch len(ret) {
	case 0:
		return true

	case 1:
		switch v := ret[0].Interface().(type) {
		case bool:
			return v
		case error: // non-nil error
			cont, err = false, v
		default: // nil error
			return true
		}

	default:
		cont = ret[0].Interface().(bool)
		err, _ = ret[1].Interface().(error) // nil error fails conversion
	}

	if err != nil {
		t.Helper()
		t.Errorf("%s error: %s", testName, err)
	}
	return cont
}

// Run runs the tests suite suite using tb as base testing
// framework. tb is typically a [*testing.T] as in:
//
//	func TestSuite(t *testing.T) {
//	  tdsuite.Run(t, &Suite{})
//	}
//
// but it can also be a [*td.T] of course.
//
// config can be used to alter the internal [*td.T] instance. See
// [td.ContextConfig] for detailed options, as:
//
//	func TestSuite(t *testing.T) {
//	  tdsuite.Run(t, &Suite{}, td.ContextConfig{
//	    UseEqual: true, // use the Equal method to compare if available
//	    BeLax:    true, // able to compare different but convertible types
//	  })
//	}
//
// Run returns true if all the tests succeeded, false otherwise.
//
// Note that if suite is not an empty struct, it should be a pointer
// if its contents has to be altered by hooks & tests methods.
//
// If suite is a pointer, it has access to non-pointer & pointer
// methods hooks & tests. If suite is not a pointer, it only has
// access to non-pointer methods hooks & tests.
func Run(tb testing.TB, suite any, config ...td.ContextConfig) bool {
	t := td.NewT(tb, config...)

	t.Helper()
	if suite == nil {
		t.Fatal("Run(): suite parameter cannot be nil")
		return false // only for tests
	}

	typ := reflect.TypeOf(suite)

	// The suite is not a pointer and in its pointer version it has
	// access to more method. Check the user isn't making a mistake by
	// not passing a pointer
	if possibleMistakes := diffWithPtrMethods(typ); len(possibleMistakes) > 0 {
		t.Logf("Run(): several methods are not accessible as suite is not a pointer but %T: %s",
			suite, strings.Join(possibleMistakes, ", "))
	}

	var methods []int
	for i, num := 0, typ.NumMethod(); i < num; i++ {
		m := typ.Method(i)

		if isTest(m.Name) {
			mt := m.Type

			if mt.IsVariadic() {
				t.Logf("Run(): method %T.%s skipped, variadic parameters not supported",
					suite, m.Name)
				continue
			}

			// Check input parameters
			switch mt.NumIn() {
			case 2:
				// TestXxx(*td.T)
				if mt.In(1) != tType {
					t.Logf("Run(): method %T.%s skipped, unrecognized parameter type %s. Only *td.T allowed",
						suite, m.Name, mt.In(1))
					continue
				}

			case 3:
				// TestXxx(*td.T, *td.T)
				if mt.In(1) != tType || mt.In(2) != tType {
					var log string
					if mt.In(1) != tType {
						if mt.In(2) != tType {
							log = fmt.Sprintf("parameters types (%s, %s)", mt.In(1), mt.In(2))
						} else {
							log = fmt.Sprintf("first parameter type %s", mt.In(1))
						}
					} else {
						log = fmt.Sprintf("second parameter type %s", mt.In(2))
					}
					t.Logf("Run(): method %T.%s skipped, unrecognized %s. Only (*td.T, *td.T) allowed",
						suite, m.Name, log)
					continue
				}

			case 1:
				t.Logf("Run(): method %T.%s skipped, no input parameters",
					suite, m.Name)
				continue

			default:
				t.Logf("Run(): method %T.%s skipped, too many parameters",
					suite, m.Name)
				continue
			}

			// Check output parameters
			switch mt.NumOut() {
			case 0:
			case 1:
				switch mt.Out(0) {
				case types.Bool, types.Error:
				default:
					t.Fatalf("Run(): method %T.%s returns %s value. Only bool or error are allowed",
						suite, m.Name, mt.Out(0))
					return false // only for tests
				}
			case 2:
				if mt.Out(0) != types.Bool || mt.Out(1) != types.Error {
					t.Fatalf("Run(): method %T.%s returns (%s, %s) values. Only (bool, error) is allowed",
						suite, m.Name, mt.Out(0), mt.Out(1))
					return false // only for tests
				}
			default:
				t.Fatalf("Run(): method %T.%s returns %d values. Only 0, 1 (bool or error) or 2 (bool, error) values are allowed",
					suite, m.Name, mt.NumOut())
				return false // only for tests
			}

			methods = append(methods, i)
		}
	}

	if len(methods) == 0 {
		t.Fatalf("Run(): no test methods found for type %T", suite)
		return false // only for tests
	}

	run(t, suite, methods)

	return !t.Failed()
}

func run(t *td.T, suite any, methods []int) {
	t.Helper()

	suiteType := reflect.TypeOf(suite)

	// setup
	if s, ok := suite.(Setup); ok {
		if err := s.Setup(t); err != nil {
			t.Errorf("%T suite setup error: %s", suite, err)
			return
		}
	} else if _, exists := suiteType.MethodByName("Setup"); exists {
		t.Errorf("%T suite has a Setup method but it does not match Setup(t *td.T) error", suite)
	}

	// destroy
	if s, ok := suite.(Destroy); ok {
		defer func() {
			if err := s.Destroy(t); err != nil {
				t.Errorf("%T suite destroy error: %s", suite, err)
			}
		}()
	} else if _, exists := suiteType.MethodByName("Destroy"); exists {
		t.Errorf("%T suite has a Destroy method but it does not match Destroy(t *td.T) error", suite)
	}

	preTest := emptyPrePostTest
	if s, ok := suite.(PreTest); ok {
		preTest = s.PreTest
	} else if _, exists := suiteType.MethodByName("PreTest"); exists {
		t.Errorf("%T suite has a PreTest method but it does not match PreTest(t *td.T, testName string) error", suite)
	}

	postTest := emptyPrePostTest
	if s, ok := suite.(PostTest); ok {
		postTest = s.PostTest
	} else if _, exists := suiteType.MethodByName("PostTest"); exists {
		t.Errorf("%T suite has a PostTest method but it does not match PostTest(t *td.T, testName string) error", suite)
	}

	between := emptyBetweenTests
	if s, ok := suite.(BetweenTests); ok {
		between = s.BetweenTests
	} else if _, exists := suiteType.MethodByName("BetweenTests"); exists {
		t.Errorf("%T suite has a BetweenTests method but it does not match BetweenTests(t *td.T, previousTestName, nextTestName string) error", suite)
	}

	vs := reflect.ValueOf(suite)
	typ := reflect.TypeOf(suite)

	for i, method := range methods {
		m := typ.Method(method)
		mt := m.Type

		call := vs.Method(method).Call

		cont := true
		if mt.NumIn() == 2 {
			t.Run(m.Name, func(t *td.T) {
				if err := preTest(t, m.Name); err != nil {
					t.Errorf("%s pre-test error: %s", m.Name, err)
					return
				}
				defer func() {
					if err := postTest(t, m.Name); err != nil {
						t.Errorf("%s post-test error: %s", m.Name, err)
					}
				}()

				cont = shouldContinue(t, m.Name, call([]reflect.Value{reflect.ValueOf(t)}))
			})
		} else {
			t.RunAssertRequire(m.Name, func(assert, require *td.T) {
				if err := preTest(assert, m.Name); err != nil {
					assert.Errorf("%s pre-test error: %s", m.Name, err)
					return
				}
				defer func() {
					if err := postTest(assert, m.Name); err != nil {
						assert.Errorf("%s post-test error: %s", m.Name, err)
					}
				}()

				cont = shouldContinue(assert, m.Name, call([]reflect.Value{
					reflect.ValueOf(assert),
					reflect.ValueOf(require),
				}))
			})
		}

		if !cont {
			t.Logf("%s required discontinuing suite tests", m.Name)
			break
		}

		if i != len(methods)-1 {
			next := typ.Method(methods[i+1]).Name
			if err := between(t, m.Name, next); err != nil {
				t.Errorf("%s / %s between-tests error: %s", m.Name, next, err)
				break
			}
		}
	}
}

func diffWithPtrMethods(typ reflect.Type) []string {
	if typ.Kind() == reflect.Ptr {
		return nil
	}

	ptyp := reflect.PtrTo(typ)
	if typ.NumMethod() == ptyp.NumMethod() {
		return nil
	}

	keep := func(m reflect.Method) bool {
		switch m.Name {
		case "Setup", "PreTest", "PostTest", "BetweenTests", "Destroy":
			return true
		default:
			return isTest(m.Name)
		}
	}

	var nonPtrMethods []string
	for i, num := 0, typ.NumMethod(); i < num; i++ {
		if m := typ.Method(i); keep(m) {
			nonPtrMethods = append(nonPtrMethods, m.Name)
		}
	}

	var onlyPtrMethods []string
	for ni, pi, num := 0, 0, ptyp.NumMethod(); pi < num; pi++ {
		if m := ptyp.Method(pi); keep(m) {
			if ni >= len(nonPtrMethods) || nonPtrMethods[ni] != m.Name {
				onlyPtrMethods = append(onlyPtrMethods, m.Name)
			} else {
				ni++
			}
		}
	}

	return onlyPtrMethods
}
