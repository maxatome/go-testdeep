package testdeep_test

import (
	"fmt"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func ExampleCmpTrue() {
	t := &testing.T{}

	got := true
	ok := CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	got = false
	ok = CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpFalse() {
	t := &testing.T{}

	got := false
	ok := CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	got = true
	ok = CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}
