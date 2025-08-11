package or_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/or"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestOrPanic(t *testing.T) {
	test.CheckPanic(t,
		func() { or.Panic(false, "big ", reflect.Uint16, " bang") },
		"big uint16 bang")

	or.Panic(true, "no panic at all")
}
