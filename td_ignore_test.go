package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestIgnore(t *testing.T) {
	checkOK(t, "any value!", Ignore())

	checkOK(t, nil, Ignore())
	checkOK(t, (*int)(nil), Ignore())

	//
	// String
	equalStr(t, Ignore().String(), "Ignore()")
}
