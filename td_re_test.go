package testdeep_test

import (
	"errors"
	"regexp"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestRe(t *testing.T) {
	//
	// string
	checkOK(t, "foo bar test", Re("bar"))
	checkOK(t, "foo bar test", Re(regexp.MustCompile("test$")))

	checkOK(t, "foo bar test",
		ReAll(`(\w+)`, Bag("bar", "test", "foo")))

	type MyString string
	checkOK(t, MyString("Ho zz hoho"),
		ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))

	// error interface
	checkOK(t, errors.New("pipo bingo"), Re("bin"))
	// fmt.Stringer interface
	checkOK(t, MyStringer{}, Re("bin"))

	checkError(t, 12, Re("bar"), expectedError{
		Message: mustBe("bad type"),
		Path:    mustBe("DATA"),
		Got:     mustBe("int"),
		Expected: mustBe(
			"string (convertible) OR fmt.Stringer OR error OR []uint8"),
	})

	checkError(t, "foo bar test", Re("pipo"), expectedError{
		Message:  mustBe("does not match Regexp"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`"foo bar test"`),
		Expected: mustBe("pipo"),
	})

	checkError(t, "foo bar test", Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`"foo bar test"`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// bytes
	checkOK(t, []byte("foo bar test"), Re("bar"))

	checkOK(t, []byte("foo bar test"),
		ReAll(`(\w+)`, Bag("bar", "test", "foo")))

	type MySlice []byte
	checkOK(t, MySlice("Ho zz hoho"),
		ReAll("(?i)(ho)", []string{"Ho", "ho", "ho"}))

	checkError(t, []int{12}, Re("bar"), expectedError{
		Message:  mustBe("bad slice type"),
		Path:     mustBe("DATA"),
		Got:      mustBe("[]int"),
		Expected: mustBe("[]uint8"),
	})

	checkError(t, []byte("foo bar test"), Re("pipo"), expectedError{
		Message:  mustBe("does not match Regexp"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`foo bar test`),
		Expected: mustBe("pipo"),
	})

	checkError(t, []byte("foo bar test"),
		Re("(pi)(po)", []string{"pi", "po"}),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// Bad usage
	const reUsage = "usage: Re("
	checkPanic(t, func() { Re(123) }, reUsage)
	checkPanic(t, func() { Re("bar", []string{}, 1) }, reUsage)

	const reAllUsage = "usage: ReAll("
	checkPanic(t, func() { ReAll(123, 456) }, reAllUsage)
}
