package testdeep_test

// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	. "github.com/maxatome/go-testdeep"
)

func ExampleCmpAll() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := CmpAll(t, got, []interface{}{Re("o/b"), HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = CmpAll(t, got, []interface{}{Re("o/b"), HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpAny() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := CmpAny(t, got, []interface{}{Re("zip"), HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = CmpAny(t, got, []interface{}{Re("zip"), HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpArray_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpArray(t, got, [3]int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArray_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpArray(t, got, MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Array(&MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpArrayEach_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"check each item of array %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"check each item of typed array %v is in [25; 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
		"check each item of typed array pointer %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"check each item of slice %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"check each item of typed slice %v is in [25; 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
		"check each item of typed slice pointer %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpBag() {
}

func ExampleCmpBetween() {
}

func ExampleCmpCap() {
}

func ExampleCmpCode() {
}

func ExampleCmpContains() {
}

func ExampleCmpGt() {
}

func ExampleCmpGte() {
}

func ExampleCmpHasPrefix() {
}

func ExampleCmpHasSuffix() {
}

func ExampleCmpIsa() {
}

func ExampleCmpLen() {
}

func ExampleCmpLt() {
}

func ExampleCmpLte() {
}

func ExampleCmpMap() {
}

func ExampleCmpMapEach() {
}

func ExampleCmpN() {
}

func ExampleCmpNil() {
}

func ExampleCmpNone() {
}

func ExampleCmpNoneOf() {
}

func ExampleCmpNotNil() {
}

func ExampleCmpPPtr() {
}

func ExampleCmpPtr() {
}

func ExampleCmpRe() {
	t := &testing.T{}

	got := "foo bar"
	ok := CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, []interface{}{Set("biz", "foo", "bar")},
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, []interface{}{Set("biz", "foo", "bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_captureAll() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpRe(t, got, `(\w+)`, []interface{}{Set("biz", "foo", "bar"), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpRe(t, got, `(\w+)`, []interface{}{Set("biz", "foo", "bar"), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_captureAllComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := CmpRe(t, got, `(\d+)`, []interface{}{ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpRe(t, got, `(\d+)`, []interface{}{ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRex() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := CmpRex(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpRex(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRex_capture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := CmpRex(t, got, expected, []interface{}{Set("biz", "foo", "bar")},
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpRex(t, got, expected, []interface{}{Set("biz", "foo", "bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRex_captureAll() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := CmpRex(t, got, expected, []interface{}{Set("biz", "foo", "bar"), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpRex(t, got, expected, []interface{}{Set("biz", "foo", "bar"), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRex_captureAllComplex() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := CmpRex(t, got, expected, []interface{}{ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpRex(t, got, expected, []interface{}{ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})), true},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpSet() {
}

func ExampleCmpShallow() {
}

func ExampleCmpSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpSlice(t, got, []int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpSlice_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpSlice(t, got, MySlice{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Slice(&MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpString() {
}

func ExampleCmpStruct() {
}

func ExampleCmpSubBagOf() {
}

func ExampleCmpSubMapOf() {
}

func ExampleCmpSubSetOf() {
}

func ExampleCmpSuperBagOf() {
}

func ExampleCmpSuperMapOf() {
}

func ExampleCmpSuperSetOf() {
}

func ExampleCmpTruncTime() {
	t := &testing.T{}

	dateToTime := func(str string) time.Time {
		t, err := time.Parse(time.RFC3339Nano, str)
		if err != nil {
			panic(err)
		}
		return t
	}

	got := dateToTime("2018-05-01T12:45:53.123456789Z")

	// Compare dates ignoring nanoseconds and monotonic part
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := CmpTruncTime(t, got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic part
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = CmpTruncTime(t, got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = CmpTruncTime(t, got, expected, 0,
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}
