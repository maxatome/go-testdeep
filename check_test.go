package testdeep_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

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
}

func (s *MyStruct) MyString() string {
	return "!"
}

var gotStruct = MyStruct{
	MyStructMid: MyStructMid{
		MyStructBase: MyStructBase{
			ValBool: true,
		},
		ValStr: "foobar",
	},
	ValInt: 123,
}

type MyInterface interface {
	MyString() string
}

type MyStringer struct{}

func (s MyStringer) String() string { return "pipo bingo" }

func checkPanic(t *testing.T, fn func(), contains string) bool {
	t.Helper()

	var (
		panicked   bool
		panicParam interface{}
	)

	func() {
		defer func() { panicParam = recover() }()
		panicked = true
		fn()
		panicked = false
	}()

	if !panicked {
		t.Error("panic() did not occur")
		return false
	}

	panicStr, ok := panicParam.(string)
	if !ok {
		t.Errorf("panic() occurred but recover()d %T type instead of string",
			panicParam)
		return false
	}

	if !strings.Contains(panicStr, contains) {
		t.Errorf("panic() string `%s'\ndoes not contain `%s'", panicStr, contains)
		return false
	}
	return true
}

type expectedError struct {
	Path     expectedErrorMatch
	Message  expectedErrorMatch
	Got      expectedErrorMatch
	Expected expectedErrorMatch
	Summary  expectedErrorMatch
	Located  bool
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
	return strings.Replace(str, "\n", "\n\t"+strings.Repeat(" ", numSpc), -1)
}

func cmpErrorStr(t *testing.T, err *Error,
	got string, expected expectedErrorMatch, fieldName string,
	args ...interface{}) bool {
	t.Helper()

	if expected.Exact != "" && got != expected.Exact {
		t.Errorf(`%sError.%s mismatch
	     got: %s
	expected: %s
	Full error:
	> %s`,
			buildTestName(args),
			fieldName, indent(got, 10), indent(expected.Exact, 10),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	if expected.Contain != "" && !strings.Contains(got, expected.Contain) {
		t.Errorf(`%sError.%s mismatch
	           got: %s
	should contain: %s
	Full error:
	> %s`,
			buildTestName(args),
			fieldName,
			indent(got, 16), indent(expected.Contain, 16),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	if expected.Match != nil && !expected.Match.MatchString(got) {
		t.Errorf(`%sError.%s mismatch
	         got: %s
	should match: %s
	Full error:
	> %s`,
			buildTestName(args),
			fieldName,
			indent(got, 14), indent(expected.Match.String(), 14),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	return true
}

func matchError(t *testing.T, err *Error, expectedError expectedError,
	expectedIsTestDeep bool, args ...interface{}) bool {
	t.Helper()

	if !cmpErrorStr(t, err, err.Message, expectedError.Message,
		"Message", args...) {
		return false
	}

	if !cmpErrorStr(t, err, err.Context.Path, expectedError.Path,
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

	// If expected is a TestDeep, the Location should be set
	if expectedIsTestDeep {
		expectedError.Located = true
	}
	if expectedError.Located != err.Location.IsInitialized() {
		t.Errorf(`%sLocation of the origin of the error
	     got: %v
	expected: %v`,
			buildTestName(args), err.Location.IsInitialized(), expectedError.Located)
		return false
	}

	if expectedError.Located &&
		!strings.HasSuffix(err.Location.File, "_test.go") {
		t.Errorf(`%sFile of the origin of the error
	     got: line %d of %s
	expected: *_test.go`,
			buildTestName(args), err.Location.Line, err.Location.File)
		return false
	}

	if expectedError.Origin != nil {
		if err.Origin == nil {
			t.Errorf(`%sError should originate from another Error`,
				buildTestName(args))
			return false
		}

		return matchError(t, err.Origin, *expectedError.Origin,
			expectedIsTestDeep, args...)
	}
	if err.Origin != nil {
		t.Errorf(`%sError should NOT originate from another Error`,
			buildTestName(args))
		return false
	}

	return true
}

func checkError(t *testing.T, got, expected interface{},
	expectedError expectedError, args ...interface{}) bool {
	t.Helper()

	err := EqDeeplyError(got, expected)
	if err == nil {
		t.Errorf("%sAn Error should have occurred", buildTestName(args))
		return false
	}

	_, expectedIsTestDeep := expected.(TestDeep)
	if !matchError(t, err, expectedError, expectedIsTestDeep, args...) {
		return false
	}

	if EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: true
	expected: false`, buildTestName(args))
		return false
	}

	return true
}

func checkErrorForEach(t *testing.T,
	gotList []interface{}, expected interface{},
	expectedError expectedError, args ...interface{}) (ret bool) {
	t.Helper()

	globalTestName := buildTestName(args)
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

func checkOK(t *testing.T, got, expected interface{},
	args ...interface{}) bool {
	t.Helper()

	if !CmpDeeply(t, got, expected, args...) {
		return false
	}

	if !EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: false
	expected: true`, buildTestName(args))
		return false
	}

	return true
}

func checkOKForEach(t *testing.T, gotList []interface{}, expected interface{},
	args ...interface{}) (ret bool) {
	t.Helper()

	globalTestName := buildTestName(args)
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

func equalStr(t *testing.T, got, expected string, args ...interface{}) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %s
	expected: %s`,
		buildTestName(args), got, expected)
	return false
}

func buildTestName(args []interface{}) string {
	switch len(args) {
	case 0:
		return ""

	case 1:
		return args[0].(string) + "\n"

	default:
		return fmt.Sprintf(args[0].(string)+"\n", args[1:]...)
	}
}

/*
func checkTrue(t *testing.T, got, expected interface{}) bool {
	t.Helper()
	return isTrue(t, EqDeeply(got, expected))
}

func checkFalse(t *testing.T, got, expected interface{}) bool {
	t.Helper()
	return isFalse(t, EqDeeply(got, expected))
}

func isTrue(t *testing.T, got bool) bool {
	if got {
		return true
	}

	t.Helper()
	t.Error(`Failed test
	     got: false
	expected: true`)
	return false
}

func isFalse(t *testing.T, got bool) bool {
	if !got {
		return true
	}

	t.Helper()
	t.Error(`Failed test
	     got: true
	expected: false`)
	return false
}
*/
