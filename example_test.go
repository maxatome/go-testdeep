package testdeep_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	. "github.com/maxatome/go-testdeep"
)

func ExampleIgnore() {
	t := &testing.T{}

	ok := CmpDeeply(t, []int{1, 2, 3},
		Slice([]int{}, ArrayEntries{
			0: 1,
			1: Ignore(), // do not care about this entry
			2: 3,
		}))
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleAll() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := CmpDeeply(t,
		got,
		All(Re("o/b"), HasSuffix("bar"), "foo/bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = CmpDeeply(t,
		got,
		All(Re("o/b"), HasSuffix("bar"), "fooX/Ybar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleAny() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := CmpDeeply(t, got, Any(Re("zip"), HasSuffix("bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = CmpDeeply(t, got, Any(Re("zip"), HasSuffix("foo")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleArray_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpDeeply(t, got, Array([3]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArray_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpDeeply(t, got, Array(MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Array(&MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleArrayEach_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"check each item of array %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"check each item of typed array %v is in [25; 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
		"check each item of typed array pointer %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"check each item of slice %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"check each item of typed slice %v is in [25; 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
		"check each item of typed slice pointer %v is in [25; 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleBag() {
}

func ExampleBetween() {
}

func ExampleCap() {
}

func ExampleCode() {
}

func ExampleContains() {
}

func ExampleGt() {
}

func ExampleGte() {
}

func ExampleHasPrefix() {
}

func ExampleHasSuffix() {
}

func ExampleIsa() {
}

func ExampleLen() {
}

func ExampleLt() {
}

func ExampleLte() {
}

func ExampleMap() {
}

func ExampleMapEach() {
}

func ExampleN() {
}

func ExampleNil() {
}

func ExampleNone() {
}

func ExampleNoneOf() {
}

func ExampleNotNil() {
}

func ExamplePPtr() {
}

func ExamplePtr() {
}

func ExampleRe() {
	t := &testing.T{}

	got := "foo bar"
	ok := CmpDeeply(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpDeeply(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpDeeply(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpDeeply(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_captureAll() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpDeeply(t, got, Re(`(\w+)`, Set("biz", "foo", "bar"), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpDeeply(t, got, Re(`(\w+)`, Set("biz", "foo", "bar"), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_captureAllComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := CmpDeeply(t, got,
		Re(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		})), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpDeeply(t, got,
		Re(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		})), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRex() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := CmpDeeply(t, got, Rex(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpDeeply(t, got, Rex(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRex_capture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := CmpDeeply(t, got, Rex(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpDeeply(t, got, Rex(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRex_captureAll() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := CmpDeeply(t, got, Rex(expected, Set("biz", "foo", "bar"), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpDeeply(t, got, Rex(expected, Set("biz", "foo", "bar"), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRex_captureAllComplex() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := CmpDeeply(t, got,
		Rex(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		})), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpDeeply(t, got,
		Rex(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		})), true),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleSet() {
}

func ExampleShallow() {
}

func ExampleSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpDeeply(t, got, Slice([]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleSlice_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpDeeply(t, got, Slice(MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Slice(&MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleString() {
}

func ExampleStruct() {
}

func ExampleSubBagOf() {
}

func ExampleSubMapOf() {
}

func ExampleSubSetOf() {
}

func ExampleSuperBagOf() {
}

func ExampleSuperMapOf() {
}

func ExampleSuperSetOf() {
}

func ExampleTruncTime() {
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
	ok := CmpDeeply(t, got, TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic part
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = CmpDeeply(t, got, TruncTime(expected, 24*time.Hour),
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = CmpDeeply(t, got, TruncTime(expected),
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}
