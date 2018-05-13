package testdeep

import (
	"testing"
)

func equalStr(t *testing.T, got, expected string) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`Failed test
	     got: %s
	expected: %s`, got, expected)
	return false
}

func TestContext(t *testing.T) {
	equalStr(t, NewContext("test").path, "test")
	equalStr(t, NewBooleanContext().path, "")

	equalStr(t, NewContext("test").AddDepth(".foo").path, "test.foo")

	equalStr(t, NewContext("test").AddDepth(".foo").AddDepth(".bar").path,
		"test.foo.bar")

	equalStr(t, NewContext("*test").AddDepth(".foo").path, "(*test).foo")

	equalStr(t, NewContext("test").AddArrayIndex(12).path, "test[12]")
	equalStr(t, NewContext("*test").AddArrayIndex(12).path, "(*test)[12]")

	equalStr(t, NewContext("test").AddPtr(2).path, "**test")
	equalStr(t, NewContext("test.foo").AddPtr(1).path, "*test.foo")
	equalStr(t, NewContext("test[3]").AddPtr(1).path, "*test[3]")
}
