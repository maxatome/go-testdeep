// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdsuite_test

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type base struct {
	calls []string
}

func (b *base) rec(plus ...string) {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	pos := strings.LastIndexByte(name, '.')
	if name[pos+1:] == "func1" { // Cleanup()
		npos := strings.LastIndexByte(name[:pos], '.')
		name = name[npos+1:pos] + ".Cleanup"
	} else {
		name = name[pos+1:]
	}
	if len(plus) > 0 {
		name += "+" + strings.Join(plus, "+")
	}
	b.calls = append(b.calls, name)
}

func (b *base) clean() {
	b.calls = b.calls[:0]
}

// Mini has only tests, no hooks.
type Mini struct{ base }

func (m *Mini) Test1(t *td.T)                     { m.rec() }
func (m *Mini) Test2(assert *td.T, require *td.T) { m.rec() }

// Full has tests and all possible hooks.
type Full struct{ base }

func (f *Full) Setup(t *td.T) error               { f.rec(); return nil }
func (f *Full) PreTest(t *td.T, tn string) error  { f.rec(tn); return nil }
func (f *Full) PostTest(t *td.T, tn string) error { f.rec(tn); return nil }
func (f *Full) BetweenTests(t *td.T, prev, next string) error {
	f.rec(prev, next)
	return nil
}
func (f *Full) Destroy(t *td.T) error { f.rec(); return nil }

func (f *Full) Test1(t *td.T)                     { f.rec() }
func (f *Full) Test2(assert *td.T, require *td.T) { f.rec() }
func (f *Full) Test3(t *td.T)                     { f.rec() }
func (f *Full) Testimony(t *td.T)                 { f.rec() } // not a test method

var (
	_ tdsuite.Setup        = (*Full)(nil)
	_ tdsuite.PreTest      = (*Full)(nil)
	_ tdsuite.PostTest     = (*Full)(nil)
	_ tdsuite.BetweenTests = (*Full)(nil)
	_ tdsuite.Destroy      = (*Full)(nil)
)

// FullBrokenHooks has all possible hooks, but wrongly defined, as they don't
// match the hook interfaces.
type FullBrokenHooks struct{}

func (*FullBrokenHooks) Setup() error                                   { return nil }
func (*FullBrokenHooks) PreTest(t *td.T, testName *string) error        { return nil }
func (*FullBrokenHooks) PostTest(t *td.T, testName string)              {}
func (*FullBrokenHooks) BetweenTests(t *td.T, prev, next *string) error { return nil }
func (*FullBrokenHooks) Destroy(t *td.T)                                {}

func (*FullBrokenHooks) Test1(_ *td.T) {}

// FullNoPtr has hooks & tests as non-pointer & pointer methods.
type FullNoPtr struct{}

var traceFullNoPtr base

func (f FullNoPtr) Setup(t *td.T) error                { traceFullNoPtr.rec(); return nil }
func (f FullNoPtr) PreTest(t *td.T, tn string) error   { traceFullNoPtr.rec(tn); return nil }
func (f *FullNoPtr) PostTest(t *td.T, tn string) error { traceFullNoPtr.rec(tn); return nil }
func (f FullNoPtr) BetweenTests(t *td.T, prev, next string) error {
	traceFullNoPtr.rec(prev, next)
	return nil
}
func (f FullNoPtr) Destroy(t *td.T) error { traceFullNoPtr.rec(); return nil }

func (f FullNoPtr) Test1(t *td.T)                      { traceFullNoPtr.rec() }
func (f *FullNoPtr) Test2(assert *td.T, require *td.T) { traceFullNoPtr.rec() }
func (f FullNoPtr) Test3(t *td.T)                      { traceFullNoPtr.rec() }
func (f *FullNoPtr) Testimony(t *td.T)                 { traceFullNoPtr.rec() } // not a test method

// ErrNone has no tests.
type ErrNone struct{}

// ErrOut1 has a Test method with bad return type.
type ErrOut1 struct{}

func (f ErrOut1) Test(t *td.T) int { return 0 }

// ErrOut2a has a Test method with bad return types.
type ErrOut2a struct{}

func (f ErrOut2a) Test(t *td.T) (bool, int) { return false, 0 }

// ErrOut2b has a Test method with bad return types.
type ErrOut2b struct{}

func (f ErrOut2b) Test(t *td.T) (int, error) { return 0, nil }

// ErrOut has a Test method with bad return types.
type ErrOut struct{}

func (f ErrOut) Test(t *td.T) (int, int, int) { return 1, 2, 3 }

// Skip has several skipped Test methods.
type Skip struct{ base }

func (s *Skip) Test1Param(i int)                  {}
func (s *Skip) Test2ParamsA(i, j int)             {}
func (s *Skip) Test2ParamsB(i int, require *td.T) {}
func (s *Skip) Test2ParamsC(assert *td.T, i int)  {}
func (s *Skip) Test3Params(t *td.T, i, j int)     {}
func (s *Skip) TestNoParams()                     {}
func (s *Skip) TestOK(t *td.T)                    { s.rec() }
func (s *Skip) TestVariadic(t ...*td.T)           {}

func TestRun(t *testing.T) {
	t.Run("Mini", func(t *testing.T) {
		suite := Mini{}
		td.CmpTrue(t, tdsuite.Run(t, &suite))
		td.Cmp(t, suite.calls, []string{"Test1", "Test2"})
	})

	t.Run("Full ptr", func(t *testing.T) {
		suite := Full{}
		td.CmpTrue(t, tdsuite.Run(t, &suite))
		ok := td.Cmp(t, suite.calls, []string{
			"Setup",
			/**/ "PreTest+Test1",
			/**/ "Test1",
			/**/ "PostTest+Test1",
			"BetweenTests+Test1+Test2",
			/**/ "PreTest+Test2",
			/**/ "Test2",
			/**/ "PostTest+Test2",
			"BetweenTests+Test2+Test3",
			/**/ "PreTest+Test3",
			/**/ "Test3",
			/**/ "PostTest+Test3",
			"Destroy",
		})
		if !ok {
			for _, c := range suite.calls {
				switch c[0] {
				case 'S', 'B', 'D':
					t.Log(c)
				default:
					t.Log("  ", c)
				}
			}
		}
	})

	t.Run("Without ptr: only non-ptr methods", func(t *testing.T) {
		defer traceFullNoPtr.clean()
		suite := FullNoPtr{}
		tb := test.NewTestingTB("TestWithoutPtr")
		td.CmpTrue(t, tdsuite.Run(tb, suite)) // non-ptr
		ok := td.Cmp(t, traceFullNoPtr.calls, []string{
			"Setup",
			/**/ "PreTest+Test1",
			/**/ "Test1",
			// /**/ "PostTest+Test1", // PostTest is a ptr method
			// Test2 is a ptr method
			// "BetweenTests+Test1+Test2",
			// /**/ "PreTest+Test2",
			// /**/ "Test2",
			// /**/ "PostTest+Test2",
			// "BetweenTests+Test2+Test3",
			"BetweenTests+Test1+Test3",
			/**/ "PreTest+Test3",
			/**/ "Test3",
			// /**/ "PostTest+Test3", // PostTest is a ptr method
			"Destroy",
		})
		if !ok {
			for _, c := range traceFullNoPtr.calls {
				switch c[0] {
				case 'S', 'B', 'D':
					t.Log(c)
				default:
					t.Log("  ", c)
				}
			}
		}
		// Yes it is a bit ugly
		td.CmpEmpty(t, tb.ContainsMessages("Run(): several methods are not accessible as suite is not a pointer but tdsuite_test.FullNoPtr: PostTest, Test2"))
	})

	t.Run("With ptr: all ptr & non-ptr methods", func(t *testing.T) {
		defer traceFullNoPtr.clean()
		suite := FullNoPtr{}
		td.CmpTrue(t, tdsuite.Run(t, &suite)) // ptr
		ok := td.Cmp(t, traceFullNoPtr.calls, []string{
			"Setup",
			/**/ "PreTest+Test1",
			/**/ "Test1",
			/**/ "PostTest+Test1",
			"BetweenTests+Test1+Test2",
			/**/ "PreTest+Test2",
			/**/ "Test2",
			/**/ "PostTest+Test2",
			"BetweenTests+Test2+Test3",
			/**/ "PreTest+Test3",
			/**/ "Test3",
			/**/ "PostTest+Test3",
			"Destroy",
		})
		if !ok {
			for _, c := range traceFullNoPtr.calls {
				switch c[0] {
				case 'S', 'B', 'D':
					t.Log(c)
				default:
					t.Log("  ", c)
				}
			}
		}
	})

	t.Run("ErrNil", func(t *testing.T) {
		tb := test.NewTestingTB("TestNil")
		tb.CatchFatal(func() { tdsuite.Run(tb, nil) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): suite parameter cannot be nil")
	})

	t.Run("ErrNone", func(t *testing.T) {
		suite := ErrNone{}
		tb := test.NewTestingTB("TestErrNone")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): no test methods found for type tdsuite_test.ErrNone")
	})

	t.Run("Full-no-ptr", func(t *testing.T) {
		suite := Full{}
		tb := test.NewTestingTB("Full-no-ptr")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.Messages, []string{
			"Run(): several methods are not accessible as suite is not a pointer but tdsuite_test.Full: BetweenTests, Destroy, PostTest, PreTest, Setup, Test1, Test2, Test3",
			"Run(): no test methods found for type tdsuite_test.Full",
		})
	})

	t.Run("ErrOut1", func(t *testing.T) {
		suite := ErrOut1{}
		tb := test.NewTestingTB("TestErrOut1")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): method tdsuite_test.ErrOut1.Test returns int value. Only bool or error are allowed")
	})

	t.Run("ErrOut2a", func(t *testing.T) {
		suite := ErrOut2a{}
		tb := test.NewTestingTB("TestErrOut2a")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): method tdsuite_test.ErrOut2a.Test returns (bool, int) values. Only (bool, error) is allowed")
	})

	t.Run("ErrOut2b", func(t *testing.T) {
		suite := ErrOut2b{}
		tb := test.NewTestingTB("TestErrOut2b")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): method tdsuite_test.ErrOut2b.Test returns (int, error) values. Only (bool, error) is allowed")
	})

	t.Run("ErrOut", func(t *testing.T) {
		suite := ErrOut{}
		tb := test.NewTestingTB("TestErrOut")
		tb.CatchFatal(func() { tdsuite.Run(tb, suite) })
		td.CmpTrue(t, tb.IsFatal)
		td.Cmp(t, tb.LastMessage(), "Run(): method tdsuite_test.ErrOut.Test returns 3 values. Only 0, 1 (bool or error) or 2 (bool, error) values are allowed")
	})

	t.Run("Skip", func(t *testing.T) {
		suite := Skip{}
		tb := test.NewTestingTB("TestSkip")
		tdsuite.Run(tb, &suite)
		test.IsFalse(t, tb.IsFatal)
		td.Cmp(t, suite.calls, []string{"TestOK"})

		const p = "Run(): method *tdsuite_test.Skip."
		td.Cmp(t, tb.Messages, []string{
			p + "Test1Param skipped, unrecognized parameter type int. Only *td.T allowed",
			p + "Test2ParamsA skipped, unrecognized parameters types (int, int). Only (*td.T, *td.T) allowed",
			p + "Test2ParamsB skipped, unrecognized first parameter type int. Only (*td.T, *td.T) allowed",
			p + "Test2ParamsC skipped, unrecognized second parameter type int. Only (*td.T, *td.T) allowed",
			p + "Test3Params skipped, too many parameters",
			p + "TestNoParams skipped, no input parameters",
			p + "TestVariadic skipped, variadic parameters not supported",
			"++++ TestOK", // (*T).Run() log as test.TestingTB has no Run() method
		})
	})
}

// Error allows to raise errors.
type Error struct {
	base
	setup        bool
	destroy      bool
	betweenTests bool
	preTest      int
	postTest     int

	testBool          [2]bool
	testError         [2]bool
	testBoolErrorBool [2]bool
	testBoolErrorErr  [2]bool
}

func (e *Error) Setup(t *td.T) error {
	if e.setup {
		return errors.New("Setup error")
	}
	return nil
}

func (e *Error) PreTest(t *td.T, tn string) error {
	if e.preTest > 0 {
		e.preTest--
		if e.preTest == 0 {
			return errors.New("PreTest error")
		}
	}
	return nil
}

func (e *Error) PostTest(t *td.T, tn string) error {
	if e.postTest > 0 {
		e.postTest--
		if e.postTest == 0 {
			return errors.New("PostTest error")
		}
	}
	return nil
}

func (e *Error) BetweenTests(t *td.T, prev, next string) error {
	if e.betweenTests {
		return errors.New("BetweenTests error")
	}
	return nil
}

func (e *Error) Destroy(t *td.T) error {
	if e.destroy {
		return errors.New("Destroy error")
	}
	return nil
}

// 1 param methods.
func (e *Error) Test1Bool(t *td.T) bool {
	e.rec()
	return !e.testBool[0]
}

func (e *Error) Test1Error(t *td.T) error {
	e.rec()
	if e.testError[0] {
		return errors.New("Test1Error error")
	}
	return nil
}

func (e *Error) Test1BoolError(t *td.T) (b bool, err error) {
	e.rec()
	b = !e.testBoolErrorBool[0]
	if e.testBoolErrorErr[0] {
		err = errors.New("Test1BoolError error")
	}
	return
}
func (e *Error) Test1Z(t *td.T) { e.rec() }

// 2 params methods.
func (e *Error) Test2Bool(assert, require *td.T) bool {
	e.rec()
	return !e.testBool[1]
}

func (e *Error) Test2Error(assert, require *td.T) error {
	e.rec()
	if e.testError[1] {
		return errors.New("Test2Error error")
	}
	return nil
}

func (e *Error) Test2BoolError(assert, require *td.T) (b bool, err error) {
	e.rec()
	b = !e.testBoolErrorBool[1]
	if e.testBoolErrorErr[1] {
		err = errors.New("Test2BoolError error")
	}
	return
}
func (e *Error) Test2Z(assert, require *td.T) { e.rec() }

func TestRunErrors(t *testing.T) {
	t.Run("Setup", func(t *testing.T) {
		suite := Error{setup: true}
		tb := test.NewTestingTB("TestError")
		td.CmpFalse(t, tdsuite.Run(tb, &suite))
		td.CmpFalse(t, tb.IsFatal)
		td.Cmp(t, tb.Messages, []string{
			"*tdsuite_test.Error suite setup error: Setup error",
		})
	})

	t.Run("Destroy", func(t *testing.T) {
		suite := Error{destroy: true}
		tb := test.NewTestingTB("TestError")
		td.CmpFalse(t, tdsuite.Run(tb, &suite))
		td.CmpFalse(t, tb.IsFatal)
		td.Cmp(t, tb.Messages, []string{
			"++++ Test1Bool",
			"++++ Test1BoolError",
			"++++ Test1Error",
			"++++ Test1Z",
			//
			"++++ Test2Bool",
			"++++ Test2BoolError",
			"++++ Test2Error",
			"++++ Test2Z",
			"*tdsuite_test.Error suite destroy error: Destroy error",
		})
	})

	t.Run("PreTest", func(t *testing.T) {
		t.Run("1 param", func(t *testing.T) {
			suite := Error{preTest: 2}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"Test1BoolError pre-test error: PreTest error",
				"++++ Test1Error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"++++ Test2BoolError",
				"++++ Test2Error",
				"++++ Test2Z",
			})
		})

		t.Run("2 params", func(t *testing.T) {
			suite := Error{preTest: 6}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"++++ Test2BoolError",
				"Test2BoolError pre-test error: PreTest error",
				"++++ Test2Error",
				"++++ Test2Z",
			})
		})
	})

	t.Run("PostTest", func(t *testing.T) {
		t.Run("1 param", func(t *testing.T) {
			suite := Error{postTest: 3}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"Test1Error post-test error: PostTest error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"++++ Test2BoolError",
				"++++ Test2Error",
				"++++ Test2Z",
			})
		})

		t.Run("2 params", func(t *testing.T) {
			suite := Error{postTest: 7}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"++++ Test2BoolError",
				"++++ Test2Error",
				"Test2Error post-test error: PostTest error",
				"++++ Test2Z",
			})
		})
	})

	t.Run("BetweenTests", func(t *testing.T) {
		suite := Error{betweenTests: true}
		tb := test.NewTestingTB("TestError")
		td.CmpFalse(t, tdsuite.Run(tb, &suite))
		td.CmpFalse(t, tb.IsFatal)
		td.Cmp(t, tb.Messages, []string{
			"++++ Test1Bool",
			"Test1Bool / Test1BoolError between-tests error: BetweenTests error",
		})
	})

	t.Run("InvalidHooks", func(t *testing.T) {
		tb := test.NewTestingTB("TestError")
		td.CmpFalse(t, tdsuite.Run(tb, &FullBrokenHooks{}))
		td.CmpFalse(t, tb.IsFatal)
		name := "*tdsuite_test.FullBrokenHooks"
		td.Cmp(t, tb.Messages, []string{
			name + " suite has a Setup method but it does not match Setup(t *td.T) error",
			name + " suite has a Destroy method but it does not match Destroy(t *td.T) error",
			name + " suite has a PreTest method but it does not match PreTest(t *td.T, testName string) error",
			name + " suite has a PostTest method but it does not match PostTest(t *td.T, testName string) error",
			name + " suite has a BetweenTests method but it does not match BetweenTests(t *td.T, previousTestName, nextTestName string) error",
			"++++ Test1",
		})
	})

	t.Run("Stop_after_TestBool", func(t *testing.T) {
		t.Run("1 param", func(t *testing.T) {
			suite := Error{testBool: [2]bool{true, false}}
			tb := test.NewTestingTB("TestError")
			td.CmpTrue(t, tdsuite.Run(tb, &suite)) // returning false is not an error
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"Test1Bool required discontinuing suite tests",
			})
		})

		t.Run("2 params", func(t *testing.T) {
			suite := Error{testBool: [2]bool{false, true}}
			tb := test.NewTestingTB("TestError")
			td.CmpTrue(t, tdsuite.Run(tb, &suite)) // returning false is not an error
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"Test2Bool required discontinuing suite tests",
			})
		})
	})

	t.Run("TestBoolError", func(t *testing.T) {
		t.Run("Stop after", func(t *testing.T) {
			t.Run("1 param", func(t *testing.T) {
				suite := Error{testBoolErrorBool: [2]bool{true, false}}
				tb := test.NewTestingTB("TestError")
				td.CmpTrue(t, tdsuite.Run(tb, &suite)) // returning false is not an error
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"Test1BoolError required discontinuing suite tests",
				})
			})

			t.Run("2 params", func(t *testing.T) {
				suite := Error{testBoolErrorBool: [2]bool{false, true}}
				tb := test.NewTestingTB("TestError")
				td.CmpTrue(t, tdsuite.Run(tb, &suite)) // returning false is not an error
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"++++ Test1Error",
					"++++ Test1Z",
					//
					"++++ Test2Bool",
					"++++ Test2BoolError",
					"Test2BoolError required discontinuing suite tests",
				})
			})
		})

		t.Run("Error but continue", func(t *testing.T) {
			t.Run("1 param", func(t *testing.T) {
				suite := Error{testBoolErrorErr: [2]bool{true, false}}
				tb := test.NewTestingTB("TestError")
				td.CmpFalse(t, tdsuite.Run(tb, &suite))
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"Test1BoolError error: Test1BoolError error",
					"++++ Test1Error",
					"++++ Test1Z",
					//
					"++++ Test2Bool",
					"++++ Test2BoolError",
					"++++ Test2Error",
					"++++ Test2Z",
				})
			})

			t.Run("2 params", func(t *testing.T) {
				suite := Error{testBoolErrorErr: [2]bool{false, true}}
				tb := test.NewTestingTB("TestError")
				td.CmpFalse(t, tdsuite.Run(tb, &suite))
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"++++ Test1Error",
					"++++ Test1Z",
					//
					"++++ Test2Bool",
					"++++ Test2BoolError",
					"Test2BoolError error: Test2BoolError error",
					"++++ Test2Error",
					"++++ Test2Z",
				})
			})
		})

		t.Run("Error and stop after", func(t *testing.T) {
			t.Run("1 param", func(t *testing.T) {
				suite := Error{
					testBoolErrorBool: [2]bool{true, false},
					testBoolErrorErr:  [2]bool{true, false},
				}
				tb := test.NewTestingTB("TestError")
				td.CmpFalse(t, tdsuite.Run(tb, &suite))
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"Test1BoolError error: Test1BoolError error",
					"Test1BoolError required discontinuing suite tests",
				})
			})

			t.Run("2 params", func(t *testing.T) {
				suite := Error{
					testBoolErrorBool: [2]bool{false, true},
					testBoolErrorErr:  [2]bool{false, true},
				}
				tb := test.NewTestingTB("TestError")
				td.CmpFalse(t, tdsuite.Run(tb, &suite))
				td.CmpFalse(t, tb.IsFatal)
				td.Cmp(t, tb.Messages, []string{
					"++++ Test1Bool",
					"++++ Test1BoolError",
					"++++ Test1Error",
					"++++ Test1Z",
					//
					"++++ Test2Bool",
					"++++ Test2BoolError",
					"Test2BoolError error: Test2BoolError error",
					"Test2BoolError required discontinuing suite tests",
				})
			})
		})
	})

	t.Run("Error_for_TestError", func(t *testing.T) {
		t.Run("1 param", func(t *testing.T) {
			suite := Error{testError: [2]bool{true, false}}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"Test1Error error: Test1Error error",
				"Test1Error required discontinuing suite tests",
			})
		})

		t.Run("2 params", func(t *testing.T) {
			suite := Error{testError: [2]bool{false, true}}
			tb := test.NewTestingTB("TestError")
			td.CmpFalse(t, tdsuite.Run(tb, &suite))
			td.CmpFalse(t, tb.IsFatal)
			td.Cmp(t, tb.Messages, []string{
				"++++ Test1Bool",
				"++++ Test1BoolError",
				"++++ Test1Error",
				"++++ Test1Z",
				//
				"++++ Test2Bool",
				"++++ Test2BoolError",
				"++++ Test2Error",
				"Test2Error error: Test2Error error",
				"Test2Error required discontinuing suite tests",
			})
		})
	})
}

// FullCleanup has tests and all possible hooks.
type FullCleanup struct{ base }

func (f *FullCleanup) Setup(t *td.T) error { f.rec(); return nil }
func (f *FullCleanup) PreTest(t *td.T, tn string) error {
	f.rec(tn)
	t.Cleanup(func() { f.rec(tn) })
	return nil
}

func (f *FullCleanup) PostTest(t *td.T, tn string) error {
	f.rec(tn)
	t.Cleanup(func() { f.rec(tn) })
	return nil
}

func (f *FullCleanup) BetweenTests(t *td.T, prev, next string) error {
	f.rec(prev, next)
	return nil
}
func (f *FullCleanup) Destroy(t *td.T) error { f.rec(); return nil }

func (f *FullCleanup) Test1(t *td.T) {
	f.rec()
	t.Cleanup(func() { f.rec() })
}

func (f *FullCleanup) Test2(assert *td.T, require *td.T) {
	f.rec()
	assert.Cleanup(func() { f.rec() })
}

func (f *FullCleanup) Test3(t *td.T) {
	f.rec()
	t.Cleanup(func() { f.rec() })
}
func (f *FullCleanup) Testimony(t *td.T) {} // not a test method

var (
	_ tdsuite.Setup        = (*FullCleanup)(nil)
	_ tdsuite.PreTest      = (*FullCleanup)(nil)
	_ tdsuite.PostTest     = (*FullCleanup)(nil)
	_ tdsuite.BetweenTests = (*FullCleanup)(nil)
	_ tdsuite.Destroy      = (*FullCleanup)(nil)
)

func TestRunCleanup(t *testing.T) {
	t.Run("Full", func(t *testing.T) {
		suite := FullCleanup{}
		td.CmpTrue(t, tdsuite.Run(t, &suite))
		ok := td.Cmp(t, suite.calls, []string{
			"Setup",
			/**/ "PreTest+Test1",
			/**/ "Test1",
			/**/ "PostTest+Test1",
			/**/ "PostTest.Cleanup+Test1",
			/**/ "Test1.Cleanup",
			/**/ "PreTest.Cleanup+Test1",
			"BetweenTests+Test1+Test2",
			/**/ "PreTest+Test2",
			/**/ "Test2",
			/**/ "PostTest+Test2",
			/**/ "PostTest.Cleanup+Test2",
			/**/ "Test2.Cleanup",
			/**/ "PreTest.Cleanup+Test2",
			"BetweenTests+Test2+Test3",
			/**/ "PreTest+Test3",
			/**/ "Test3",
			/**/ "PostTest+Test3",
			/**/ "PostTest.Cleanup+Test3",
			/**/ "Test3.Cleanup",
			/**/ "PreTest.Cleanup+Test3",
			"Destroy",
		})
		if !ok {
			for _, c := range suite.calls {
				switch c[0] {
				case 'S', 'B', 'D':
					t.Log(c)
				default:
					t.Log("  ", c)
				}
			}
		}
	})
}
