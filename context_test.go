package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestContext(t *testing.T) {
	equalStr(t, NewContext("test").Path, "test")
	equalStr(t, NewBooleanContext().Path, "")

	equalStr(t, NewContext("test").AddDepth(".foo").Path, "test.foo")

	equalStr(t, NewContext("test").AddDepth(".foo").AddDepth(".bar").Path,
		"test.foo.bar")

	equalStr(t, NewContext("*test").AddDepth(".foo").Path, "(*test).foo")

	equalStr(t, NewContext("test").AddArrayIndex(12).Path, "test[12]")
	equalStr(t, NewContext("*test").AddArrayIndex(12).Path, "(*test)[12]")

	equalStr(t, NewContext("test").AddPtr(2).Path, "**test")
	equalStr(t, NewContext("test.foo").AddPtr(1).Path, "*test.foo")
	equalStr(t, NewContext("test[3]").AddPtr(1).Path, "*test[3]")
}
