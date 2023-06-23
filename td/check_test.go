// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestMain(m *testing.M) {
	color.SaveState()
	os.Exit(m.Run())
}

type MyStructBase struct {
	ValBool bool
}

type MyStructMid struct {
	MyStructBase
	ValStr string
}

type MyStruct struct {
	MyStructMid
	ValInt int
	Ptr    *int
}

func (s *MyStruct) MyString() string {
	return "!"
}

type MyInterface interface {
	MyString() string
}

type MyStringer struct{}

func (s MyStringer) String() string { return "pipo bingo" }

type expectedError struct {
	Path     expectedErrorMatch
	Message  expectedErrorMatch
	Got      expectedErrorMatch
	Expected expectedErrorMatch
	Summary  expectedErrorMatch
	Located  bool
	Under    expectedErrorMatch
	Origin   *expectedError
}

type expectedErrorMatch struct {
	Exact   string
	Match   *regexp.Regexp
	Contain string
}

func mustBe(str string) expectedErrorMatch {
	return expectedErrorMatch{Exact: str}
}

func mustMatch(str string) expectedErrorMatch {
	return expectedErrorMatch{Match: regexp.MustCompile(str)}
}

func mustContain(str string) expectedErrorMatch {
	return expectedErrorMatch{Contain: str}
}

func indent(str string, numSpc int) string {
	return strings.ReplaceAll(str, "\n", "\n\t"+strings.Repeat(" ", numSpc))
}

func fullError(err *ctxerr.Error) string {
	return strings.ReplaceAll(err.Error(), "\n", "\n\t> ")
}

func cmpErrorStr(t *testing.T, err *ctxerr.Error,
	got string, expected expectedErrorMatch, fieldName string,
	args ...any,
) bool {
	t.Helper()

	if expected.Exact != "" && got != expected.Exact {
		t.Errorf(`%sError.%s mismatch
	     got: %s
	expected: %s
	Full error:
	> %s`,
			tdutil.BuildTestName(args...),
			fieldName, indent(got, 10), indent(expected.Exact, 10),
			fullError(err))
		return false
	}

	if expected.Contain != "" && !strings.Contains(got, expected.Contain) {
		t.Errorf(`%sError.%s mismatch
	           got: %s
	should contain: %s
	Full error:
	> %s`,
			tdutil.BuildTestName(args...),
			fieldName,
			indent(got, 16), indent(expected.Contain, 16),
			fullError(err))
		return false
	}

	if expected.Match != nil && !expected.Match.MatchString(got) {
		t.Errorf(`%sError.%s mismatch
	         got: %s
	should match: %s
	Full error:
	> %s`,
			tdutil.BuildTestName(args...),
			fieldName,
			indent(got, 14), indent(expected.Match.String(), 14),
			fullError(err))
		return false
	}

	return true
}

func matchError(t *testing.T, err *ctxerr.Error, expectedError expectedError,
	expectedIsTestDeep bool, args ...any,
) bool {
	t.Helper()

	if !cmpErrorStr(t, err, err.Message, expectedError.Message,
		"Message", args...) {
		return false
	}

	if !cmpErrorStr(t, err, err.Context.Path.String(), expectedError.Path,
		"Context.Path", args...) {
		return false
	}

	if !cmpErrorStr(t, err, err.GotString(), expectedError.Got, "Got", args...) {
		return false
	}

	if !cmpErrorStr(t, err,
		err.ExpectedString(), expectedError.Expected, "Expected", args...) {
		return false
	}

	if !cmpErrorStr(t, err,
		err.SummaryString(), expectedError.Summary, "Summary", args...) {
		return false
	}

	// under
	serr, under := err.Error(), ""
	if pos := strings.Index(serr, "\n[under operator "); pos > 0 {
		under = serr[pos+2:]
		under = under[:strings.IndexByte(under, ']')]
	}
	if !cmpErrorStr(t, err, under, expectedError.Under, "[under operator …]", args...) {
		return false
	}

	// If expected is a TestDeep, the Location should be set
	if expectedIsTestDeep {
		expectedError.Located = true
	}
	if expectedError.Located != err.Location.IsInitialized() {
		t.Errorf(`%sLocation of the origin of the error
	     got: %v
	expected: %v`,
			tdutil.BuildTestName(args...), err.Location.IsInitialized(), expectedError.Located)
		return false
	}

	if expectedError.Located &&
		!strings.HasSuffix(err.Location.File, "_test.go") {
		t.Errorf(`%sFile of the origin of the error
	     got: line %d of %s
	expected: *_test.go`,
			tdutil.BuildTestName(args...), err.Location.Line, err.Location.File)
		return false
	}

	if expectedError.Origin != nil {
		if err.Origin == nil {
			t.Errorf(`%sError should originate from another Error`,
				tdutil.BuildTestName(args...))
			return false
		}

		return matchError(t, err.Origin, *expectedError.Origin,
			expectedIsTestDeep, args...)
	}
	if err.Origin != nil {
		t.Errorf(`%sError should NOT originate from another Error`,
			tdutil.BuildTestName(args...))
		return false
	}

	return true
}

func _checkError(t *testing.T, got, expected any,
	expectedError expectedError, args ...any,
) bool {
	t.Helper()

	err := td.EqDeeplyError(got, expected)
	if err == nil {
		t.Errorf("%sAn Error should have occurred", tdutil.BuildTestName(args...))
		return false
	}

	_, expectedIsTestDeep := expected.(td.TestDeep)
	if !matchError(t, err.(*ctxerr.Error), expectedError, expectedIsTestDeep, args...) {
		return false
	}

	if td.EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: true
	expected: false`, tdutil.BuildTestName(args...))
		return false
	}

	return true
}

func ifaceExpectedError(t *testing.T, expectedError expectedError) expectedError {
	t.Helper()

	if !strings.Contains(expectedError.Path.Exact, "DATA") {
		return expectedError
	}

	newExpectedError := expectedError
	newExpectedError.Path.Exact = strings.Replace(expectedError.Path.Exact,
		"DATA", "DATA.Iface", 1)

	if newExpectedError.Origin != nil {
		newOrigin := ifaceExpectedError(t, *newExpectedError.Origin)
		newExpectedError.Origin = &newOrigin
	}

	return newExpectedError
}

// checkError calls _checkError twice. The first time with the same
// parameters, the second time in an any context.
func checkError(t *testing.T, got, expected any,
	expectedError expectedError, args ...any,
) bool {
	t.Helper()

	if ok := _checkError(t, got, expected, expectedError, args...); !ok {
		return false
	}

	type tmpStruct struct {
		Iface any
	}

	return _checkError(t, tmpStruct{Iface: got},
		td.Struct(
			tmpStruct{},
			td.StructFields{
				"Iface": expected,
			}),
		ifaceExpectedError(t, expectedError),
		args...)
}

func checkErrorForEach(t *testing.T,
	gotList []any, expected any,
	expectedError expectedError, args ...any,
) (ret bool) {
	t.Helper()

	globalTestName := tdutil.BuildTestName(args...)
	ret = true

	for idx, got := range gotList {
		testName := fmt.Sprintf("Got #%d", idx)
		if globalTestName != "" {
			testName += ", " + globalTestName
		}
		ret = checkError(t, got, expected, expectedError, testName) && ret
	}
	return
}

// customCheckOK calls chk twice. The first time with the same
// parameters, the second time in an any context.
func customCheckOK(t *testing.T,
	chk func(t *testing.T, got, expected any, args ...any) bool,
	got, expected any,
	args ...any,
) bool {
	t.Helper()

	if ok := chk(t, got, expected, args...); !ok {
		return false
	}

	type tmpStruct struct {
		Iface any
	}

	// Dirty hack to force got be passed as an interface kind
	return chk(t, tmpStruct{Iface: got},
		td.Struct(
			tmpStruct{},
			td.StructFields{
				"Iface": expected,
			}),
		args...)
}

func _checkOK(t *testing.T, got, expected any,
	args ...any,
) bool {
	t.Helper()

	if !td.Cmp(t, got, expected, args...) {
		return false
	}

	if !td.EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: false
	expected: true`, tdutil.BuildTestName(args...))
		return false
	}

	if err := td.EqDeeplyError(got, expected); err != nil {
		t.Errorf(`%sEqDeeplyError returned an error: %s`,
			tdutil.BuildTestName(args...), err)
		return false
	}

	return true
}

// checkOK calls _checkOK twice. The first time with the same
// parameters, the second time in an any context.
func checkOK(t *testing.T, got, expected any,
	args ...any,
) bool {
	t.Helper()
	return customCheckOK(t, _checkOK, got, expected, args...)
}

func checkOKOrPanicIfUnsafeDisabled(t *testing.T, got, expected any,
	args ...any,
) (ret bool) {
	t.Helper()

	cmp := func() {
		t.Helper()
		ret = _checkOK(t, got, expected, args...)
	}

	// Should panic if unsafe package is not available
	if dark.UnsafeDisabled {
		return test.CheckPanic(t, cmp,
			"dark.GetInterface() does not handle private ")
	}

	cmp()
	return
}

func checkOKForEach(t *testing.T, gotList []any, expected any,
	args ...any,
) (ret bool) {
	t.Helper()

	globalTestName := tdutil.BuildTestName(args...)
	ret = true

	for idx, got := range gotList {
		testName := fmt.Sprintf("Got #%d", idx)
		if globalTestName != "" {
			testName += ", " + globalTestName
		}
		ret = checkOK(t, got, expected, testName) && ret
	}
	return
}

func equalTypes(t *testing.T, got td.TestDeep, expected any, args ...any) bool {
	gotType := got.TypeBehind()

	expectedType, ok := expected.(reflect.Type)
	if !ok {
		expectedType = reflect.TypeOf(expected)
	}

	if gotType == expectedType {
		return true
	}

	var gotStr, expectedStr string

	if gotType == nil {
		gotStr = "nil"
	} else {
		gotStr = gotType.String()
	}

	if expected == nil {
		expectedStr = "nil"
	} else {
		expectedStr = expectedType.String()
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %s
	expected: %s`,
		tdutil.BuildTestName(args...), gotStr, expectedStr)
	return false
}
