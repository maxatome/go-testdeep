package location

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestNew(t *testing.T) {
	for _, curTest := range []struct {
		callDepth int
		expectOK  bool
	}{
		{callDepth: 0, expectOK: true},
		{callDepth: 1, expectOK: true},
		{callDepth: 100, expectOK: false},
	} {
		loc, ok := New(curTest.callDepth)
		test.EqualBool(t, ok, curTest.expectOK, "callDepth=%d", curTest.callDepth)
		if ok {
			test.IsTrue(t, loc.File != "", "File should not be empty for callDepth=%d", curTest.callDepth)
			test.IsTrue(t, loc.Func != "", "Func should not be empty for callDepth=%d", curTest.callDepth)
			test.IsTrue(t, loc.Line > 0, "Line should be >0 for callDepth=%d", curTest.callDepth)
		}
	}
}

func TestIsInitialized(t *testing.T) {
	for _, curTest := range []struct {
		loc      Location
		expected bool
	}{
		{loc: Location{}, expected: false},
		{loc: Location{File: ""}, expected: false},
		{loc: Location{File: "test.go"}, expected: true},
		{loc: Location{File: "test.go", Func: "Test", Line: 10}, expected: true},
	} {
		test.EqualBool(t, curTest.loc.IsInitialized(), curTest.expected, "loc=%+v", curTest.loc)
	}
}

func TestString(t *testing.T) {
	for _, curTest := range []struct {
		loc      Location
		expected string
	}{
		{
			loc:      Location{File: "test.go", Func: "TestFunc", Line: 10},
			expected: "TestFunc at test.go:10",
		},
		{
			loc:      Location{File: "main.go", Func: "main", Line: 5, Inside: "inside something "},
			expected: "main  inside somethingat main.go:5",
		},
		{
			loc:      Location{File: "pkg/utils.go", Func: "utils.Helper", Line: 42, Inside: "in map "},
			expected: "utils.Helper  in mapat pkg/utils.go:42",
		},
	} {
		test.EqualStr(t, curTest.loc.String(), curTest.expected, "loc=%+v", curTest.loc)
	}
}
