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
	checkOK(t, "foo bar test", Rex(regexp.MustCompile("test$")))

	checkOK(t, "foo bar test",
		Re(`(\w+)`, Bag("bar", "test", "foo"), true))

	type MyString string
	checkOK(t, MyString("Ho zz hoho"),
		Re("(?i)(ho)", []string{"Ho", "ho", "ho"}, true))

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
		Re(`(\w+)`, Bag("bar", "test", "foo"), true))

	type MySlice []byte
	checkOK(t, MySlice("Ho zz hoho"),
		Re("(?i)(ho)", []string{"Ho", "ho", "ho"}, true))

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
		Re("(pi)(po)", []string{"pi", "po"}, false),
		expectedError{
			Message:  mustBe("does not match Regexp"),
			Path:     mustBe("DATA"),
			Got:      mustContain(`foo bar test`),
			Expected: mustBe("(pi)(po)"),
		})

	//
	// Bad usage
	const reUsage = "usage: Re("
	checkPanic(t, func() { Re("bar", nil) }, reUsage)
	checkPanic(t, func() { Re("bar", []string{}, 1) }, reUsage)
	checkPanic(t, func() { Re("bar", []string{}, true, 123) }, reUsage)

	const rexUsage = "usage: Rex("
	re := regexp.MustCompile("bar")
	checkPanic(t, func() { Rex(re, nil) }, rexUsage)
	checkPanic(t, func() { Rex(re, []string{}, 1) }, rexUsage)
	checkPanic(t, func() { Rex(re, []string{}, true, 123) }, rexUsage)
}
