// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
)

func Example() {
	t := &testing.T{}

	dateToTime := func(str string) time.Time {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			panic(err)
		}
		return t
	}

	type PetFamily uint8

	const (
		Canidae PetFamily = 1
		Felidae PetFamily = 2
	)

	type Pet struct {
		Name     string
		Birthday time.Time
		Family   PetFamily
	}

	type Master struct {
		Name         string
		AnnualIncome int
		Pets         []*Pet
	}

	// Imagine a function returning a Master slice...
	masters := []Master{
		{
			Name:         "Bob Smith",
			AnnualIncome: 25000,
			Pets: []*Pet{
				{
					Name:     "Quizz",
					Birthday: dateToTime("2010-11-05T10:00:00Z"),
					Family:   Canidae,
				},
				{
					Name:     "Charlie",
					Birthday: dateToTime("2013-05-11T08:00:00Z"),
					Family:   Canidae,
				},
			},
		},
		{
			Name:         "John Doe",
			AnnualIncome: 38000,
			Pets: []*Pet{
				{
					Name:     "Coco",
					Birthday: dateToTime("2015-08-05T18:00:00Z"),
					Family:   Felidae,
				},
				{
					Name:     "Lucky",
					Birthday: dateToTime("2014-04-17T07:00:00Z"),
					Family:   Canidae,
				},
			},
		},
	}

	// Let's check masters slice contents
	ok := td.Cmp(t, masters, td.All(
		td.Len(td.Gt(0)), // len(masters) should be > 0
		td.ArrayEach(
			// For each Master
			td.Struct(Master{}, td.StructFields{
				// Master Name should be composed of 2 words, with 1st letter uppercased
				"Name": td.Re(`^[A-Z][a-z]+ [A-Z][a-z]+\z`),
				// Annual income should be greater than $10000
				"AnnualIncome": td.Gt(10000),
				"Pets": td.ArrayEach(
					// For each Pet
					td.Struct(&Pet{}, td.StructFields{
						// Pet Name should be composed of 1 word, with 1st letter uppercased
						"Name": td.Re(`^[A-Z][a-z]+\z`),
						"Birthday": td.All(
							// Pet should be born after 2010, January 1st, but before now!
							td.Between(dateToTime("2010-01-01T00:00:00Z"), time.Now()),
							// AND minutes, seconds and nanoseconds should be 0
							td.Code(func(t time.Time) bool {
								return t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
							}),
						),
						// Only dogs and cats allowed
						"Family": td.Any(Canidae, Felidae),
					}),
				),
			}),
		),
	))
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleIgnore() {
	t := &testing.T{}

	ok := td.Cmp(t, []int{1, 2, 3},
		td.Slice([]int{}, td.ArrayEntries{
			0: 1,
			1: td.Ignore(), // do not care about this entry
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
	ok := td.Cmp(t,
		got,
		td.All(td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = td.Cmp(t,
		got,
		td.All(td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"),
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []any copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = td.Cmp(t,
		got,
		td.All(td.HasPrefix("foo"), regOps, td.HasSuffix("bar")),
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleAny() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := td.Cmp(t, got, td.Any(td.Re("zip"), td.HasSuffix("bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = td.Cmp(t, got, td.Any(td.Re("zip"), td.HasSuffix("foo")),
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []any copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("a/c"), td.Re(`^xx`), td.Re(`ar$`)})
	ok = td.Cmp(t,
		got,
		td.Any(td.HasPrefix("xxx"), regOps, td.HasSuffix("zip")),
		"check at least one operator matches value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleArray_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array([3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Simple array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&[3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.Array((*[3]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Array pointer, nil model:", ok)

	// Output:
	// Simple array: true
	// Array pointer: true
	// Array pointer, nil model: true
}

func ExampleArray_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array(MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks typed array %v", got)
	fmt.Println("Typed array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, empty model:", ok)

	ok = td.Cmp(t, &got,
		td.Array((*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, nil model:", ok)

	// Output:
	// Typed array: true
	// Pointer on a typed array: true
	// Pointer on a typed array, empty model: true
	// Pointer on a typed array, nil model: true
}

func ExampleArrayEach_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.Cmp(t, got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.Cmp(t, got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.Cmp(t, got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.Cmp(t, got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.ArrayEach(td.Between(25, 60)),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleBag() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := td.Cmp(t, got, td.Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = td.Cmp(t, got, td.Bag(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = td.Cmp(t, got, td.Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = td.Cmp(t, got, td.Bag(1, 2, 3, 5, td.Gt(7)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5}
	ok = td.Cmp(t, got, td.Bag(td.Flatten(expected), td.Gt(7)),
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
	// true
}

func ExampleBetween_int() {
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutOut),
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleBetween_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Between("aaa", "abc"),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInIn),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInOut),
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutIn),
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutOut),
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleBetween_time() {
	t := &testing.T{}

	before := time.Now()
	occurredAt := time.Now()
	after := time.Now()

	ok := td.Cmp(t, occurredAt, td.Between(before, after))
	fmt.Println("It occurred between before and after:", ok)

	type MyTime time.Time
	ok = td.Cmp(t, MyTime(occurredAt), td.Between(MyTime(before), MyTime(after)))
	fmt.Println("Same for convertible MyTime type:", ok)

	ok = td.Cmp(t, MyTime(occurredAt), td.Between(before, after))
	fmt.Println("MyTime vs time.Time:", ok)

	ok = td.Cmp(t, occurredAt, td.Between(before, 10*time.Second))
	fmt.Println("Using a time.Duration as TO:", ok)

	ok = td.Cmp(t, MyTime(occurredAt), td.Between(MyTime(before), 10*time.Second))
	fmt.Println("Using MyTime as FROM and time.Duration as TO:", ok)

	// Output:
	// It occurred between before and after: true
	// Same for convertible MyTime type: true
	// MyTime vs time.Time: false
	// Using a time.Duration as TO: true
	// Using MyTime as FROM and time.Duration as TO: true
}

func ExampleCap() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.Cmp(t, got, td.Cap(12), "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.Cmp(t, got, td.Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCap_operator() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.Cmp(t, got, td.Cap(td.Between(10, 12)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Cap(td.Gt(10)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCatch() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	var age int
	ok := td.Cmp(t, got,
		td.JSON(`{"age":$1,"fullname":"Bob"}`,
			td.Catch(&age, td.Between(40, 45))))
	fmt.Println("check got age+fullname:", ok)
	fmt.Println("caught age:", age)

	// Output:
	// check got age+fullname: true
	// caught age: 42
}

func ExampleCode() {
	t := &testing.T{}

	got := "12"

	ok := td.Cmp(t, got,
		td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = td.Cmp(t, got,
		td.Code(func(num string) (bool, string) {
			n, err := strconv.Atoi(num)
			if err != nil {
				return false, "not a number"
			}
			if n > 10 && n < 100 {
				return true, ""
			}
			return false, "not in ]10 .. 100["
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason thanks to error
	ok = td.Cmp(t, got,
		td.Code(func(num string) error {
			n, err := strconv.Atoi(num)
			if err != nil {
				return err
			}
			if n > 10 && n < 100 {
				return nil
			}
			return fmt.Errorf("%d not in ]10 .. 100[", n)
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCode_custom() {
	t := &testing.T{}

	got := 123

	ok := td.Cmp(t, got, td.Code(func(t *td.T, num int) {
		t.Cmp(num, 123)
	}))
	fmt.Println("with one *td.T:", ok)

	ok = td.Cmp(t, got, td.Code(func(assert, require *td.T, num int) {
		assert.Cmp(num, 123)
		require.Cmp(num, 123)
	}))
	fmt.Println("with assert & require *td.T:", ok)

	// Output:
	// with one *td.T: true
	// with assert & require *td.T: true
}

func ExampleContains_arraySlice() {
	t := &testing.T{}

	ok := td.Cmp(t, [...]int{11, 22, 33, 44}, td.Contains(22))
	fmt.Println("array contains 22:", ok)

	ok = td.Cmp(t, [...]int{11, 22, 33, 44}, td.Contains(td.Between(20, 25)))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains(22))
	fmt.Println("slice contains 22:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains(td.Between(20, 25)))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains([]int{22, 33}))
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true
}

func ExampleContains_nil() {
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := td.Cmp(t, got, td.Contains(nil))
	fmt.Println("array contains untyped nil:", ok)

	ok = td.Cmp(t, got, td.Contains((*int)(nil)))
	fmt.Println("array contains *int nil:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Nil()))
	fmt.Println("array contains Nil():", ok)

	ok = td.Cmp(t, got, td.Contains((*byte)(nil)))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false
}

func ExampleContains_map() {
	t := &testing.T{}

	ok := td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Contains(22))
	fmt.Println("map contains value 22:", ok)

	ok = td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33},
		td.Contains(td.Between(20, 25)))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true
}

func ExampleContains_string() {
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains([]byte("oob")), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleContains_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleContains_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleContainsKey() {
	t := &testing.T{}

	ok := td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.ContainsKey("foo"))
	fmt.Println(`map contains key "foo":`, ok)

	ok = td.Cmp(t,
		map[int]bool{12: true, 24: false, 42: true, 51: false},
		td.ContainsKey(td.Between(40, 50)))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = td.Cmp(t,
		map[string]int{"FOO": 11, "bar": 22, "zip": 33},
		td.ContainsKey(td.Smuggle(strings.ToLower, "foo")))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true
}

func ExampleContainsKey_nil() {
	t := &testing.T{}

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := td.Cmp(t, got, td.ContainsKey(nil))
	fmt.Println("map contains untyped nil key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey((*int)(nil)))
	fmt.Println("map contains *int nil key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey(td.Nil()))
	fmt.Println("map contains Nil() key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey((*byte)(nil)))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false
}

func ExampleDelay() {
	t := &testing.T{}

	cmpNow := func(expected td.TestDeep) bool {
		time.Sleep(time.Microsecond) // imagine a DB insert returning a CreatedAt
		return td.Cmp(t, time.Now(), expected)
	}

	before := time.Now()

	ok := cmpNow(td.Between(before, time.Now()))
	fmt.Println("Between called before compare:", ok)

	ok = cmpNow(td.Delay(func() td.TestDeep {
		return td.Between(before, time.Now())
	}))
	fmt.Println("Between delayed until compare:", ok)

	// Output:
	// Between called before compare: false
	// Between delayed until compare: true
}

func ExampleEmpty() {
	t := &testing.T{}

	ok := td.Cmp(t, nil, td.Empty()) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = td.Cmp(t, (*int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, "", td.Empty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = td.Cmp(t, 0, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{3}, td.Empty()) // fails, as not empty
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.Empty()) // fails, Empty() is not Zero()!
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// false
	// true
	// true
	// true
	// true
	// false
	// false
}

func ExampleEmpty_pointers() {
	t := &testing.T{}

	type MySlice []int

	ok := td.Cmp(t, MySlice{}, td.Empty()) // Ptr() not needed
	fmt.Println(ok)

	ok = td.Cmp(t, &MySlice{}, td.Empty())
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = td.Cmp(t, &l3, td.Empty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.Cmp(t, &MyStruct{}, td.Empty()) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleErrorIs() {
	t := &testing.T{}

	err1 := fmt.Errorf("failure1")
	err2 := fmt.Errorf("failure2: %w", err1)
	err3 := fmt.Errorf("failure3: %w", err2)
	err := fmt.Errorf("failure4: %w", err3)

	ok := td.Cmp(t, err, td.ErrorIs(err))
	fmt.Println("error is itself:", ok)

	ok = td.Cmp(t, err, td.ErrorIs(err1))
	fmt.Println("error is also err1:", ok)

	ok = td.Cmp(t, err1, td.ErrorIs(err))
	fmt.Println("err1 is err:", ok)

	// Output:
	// error is itself: true
	// error is also err1: true
	// err1 is err: false
}

func ExampleFirst_classic() {
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.First(td.Gt(0), 1))
	fmt.Println("first positive number is 1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.First(isEven, -2))
	fmt.Println("first even number is -2:", ok)

	ok = td.Cmp(t, got, td.First(isEven, td.Lt(0)))
	fmt.Println("first even number is < 0:", ok)

	ok = td.Cmp(t, got, td.First(isEven, td.Code(isEven)))
	fmt.Println("first even number is well even:", ok)

	// Output:
	// first positive number is 1: true
	// first even number is -2: true
	// first even number is < 0: true
	// first even number is well even: true
}

func ExampleFirst_empty() {
	t := &testing.T{}

	ok := td.Cmp(t, ([]int)(nil), td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in nil slice:", ok)

	ok = td.Cmp(t, []int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty slice:", ok)

	ok = td.Cmp(t, &[]int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty pointed slice:", ok)

	ok = td.Cmp(t, [0]int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty array:", ok)

	// Output:
	// first in nil slice: false
	// first in empty slice: false
	// first in empty pointed slice: false
	// first in empty array: false
}

func ExampleFirst_struct() {
	t := &testing.T{}

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      37,
		},
	}

	ok := td.Cmp(t, got, td.First(
		td.Smuggle("Age", td.Gt(30)),
		td.Smuggle("Fullname", "Bob Foobar")))
	fmt.Println("first person.Age > 30 → Bob:", ok)

	ok = td.Cmp(t, got, td.First(
		td.JSONPointer("/age", td.Gt(30)),
		td.SuperJSONOf(`{"fullname":"Bob Foobar"}`)))
	fmt.Println("first person.Age > 30 → Bob, using JSON:", ok)

	ok = td.Cmp(t, got, td.First(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSONPointer("/fullname", td.HasPrefix("Bob"))))
	fmt.Println("first person.Age > 30 → Bob, using JSONPointer:", ok)

	// Output:
	// first person.Age > 30 → Bob: true
	// first person.Age > 30 → Bob, using JSON: true
	// first person.Age > 30 → Bob, using JSONPointer: true
}

func ExampleFirst_json() {
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": First(Gt(2), 3)}`))
	fmt.Println("first number > 2:", ok)

	got = map[string]any{
		"persons": []map[string]any{
			{"id": 1, "name": "Joe"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Alice"},
			{"id": 4, "name": "Brian"},
			{"id": 5, "name": "Britt"},
		},
	}
	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": First(JSONPointer("/name", "Brian"), {"id": 4, "name": "Brian"})
}`))
	fmt.Println(`is "Brian" content OK:`, ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": First(JSONPointer("/name", "Brian"), JSONPointer("/id", 4))
}`))
	fmt.Println(`ID of "Brian" is 4:`, ok)

	// Output:
	// first number > 2: true
	// is "Brian" content OK: true
	// ID of "Brian" is 4: true
}

func ExampleGrep_classic() {
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.Grep(td.Gt(0), []int{1, 2, 3}))
	fmt.Println("check positive numbers:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.Grep(isEven, []int{-2, 0, 2}))
	fmt.Println("even numbers are -2, 0 and 2:", ok)

	ok = td.Cmp(t, got, td.Grep(isEven, td.Set(0, 2, -2)))
	fmt.Println("even numbers are also 0, 2 and -2:", ok)

	ok = td.Cmp(t, got, td.Grep(isEven, td.ArrayEach(td.Code(isEven))))
	fmt.Println("even numbers are each even:", ok)

	// Output:
	// check positive numbers: true
	// even numbers are -2, 0 and 2: true
	// even numbers are also 0, 2 and -2: true
	// even numbers are each even: true
}

func ExampleGrep_nil() {
	t := &testing.T{}

	var got []int
	ok := td.Cmp(t, got, td.Grep(td.Gt(0), ([]int)(nil)))
	fmt.Println("typed []int nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), ([]string)(nil)))
	fmt.Println("typed []string nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), td.Nil()))
	fmt.Println("td.Nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), []int{}))
	fmt.Println("empty non-nil slice:", ok)

	// Output:
	// typed []int nil: true
	// typed []string nil: false
	// td.Nil: true
	// empty non-nil slice: false
}

func ExampleGrep_struct() {
	t := &testing.T{}

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      27,
		},
	}

	ok := td.Cmp(t, got, td.Grep(
		td.Smuggle("Age", td.Gt(30)),
		td.All(
			td.Len(1),
			td.ArrayEach(td.Smuggle("Fullname", "Bob Foobar")),
		)))
	fmt.Println("person.Age > 30 → only Bob:", ok)

	ok = td.Cmp(t, got, td.Grep(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSON(`[ SuperMapOf({"fullname":"Bob Foobar"}) ]`)))
	fmt.Println("person.Age > 30 → only Bob, using JSON:", ok)

	// Output:
	// person.Age > 30 → only Bob: true
	// person.Age > 30 → only Bob, using JSON: true
}

func ExampleGrep_json() {
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": Grep(Gt(2), [3, 4])}`))
	fmt.Println("grep a number > 2:", ok)

	got = map[string]any{
		"persons": []map[string]any{
			{"id": 1, "name": "Joe"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Alice"},
			{"id": 4, "name": "Brian"},
			{"id": 5, "name": "Britt"},
		},
	}
	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": Grep(JSONPointer("/name", HasPrefix("Br")), [
    {"id": 4, "name": "Brian"},
    {"id": 5, "name": "Britt"},
  ])
}`))
	fmt.Println(`grep "Br" prefix:`, ok)

	// Output:
	// grep a number > 2: true
	// grep "Br" prefix: true
}

func ExampleGt_int() {
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Gt(155), "checks %v is > 155", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gt(156), "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleGt_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Gt("abb"), `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gt("abc"), `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleGte_int() {
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Gte(156), "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gte(155), "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gte(157), "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleGte_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Gte("abc"), `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gte("abb"), `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gte("abd"), `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleIsa() {
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := td.Cmp(t, got, td.Isa(TstStruct{}), "checks got is a TstStruct")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Isa(&TstStruct{}),
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.Isa(&TstStruct{}),
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleIsa_interface() {
	t := &testing.T{}

	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.Isa((*fmt.Stringer)(nil)),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = td.Cmp(t, errGot, td.Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = td.Cmp(t, errGot, td.Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = td.Cmp(t, &errGot, td.Isa((*error)(nil)),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleJSON_basic() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age":42,"fullname":"Bob"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob","age":42}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = td.Cmp(t, got, td.JSON(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with gender field:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob"}`))
	fmt.Println("check got with fullname only:", ok)

	ok = td.Cmp(t, true, td.JSON(`true`))
	fmt.Println("check boolean got is true:", ok)

	ok = td.Cmp(t, 42, td.JSON(`42`))
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = td.Cmp(t, got, td.JSON(`null`))
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true
}

func ExampleJSON_placeholders() {
	t := &testing.T{}

	type Person struct {
		Fullname string    `json:"fullname"`
		Age      int       `json:"age"`
		Children []*Person `json:"children,omitempty"`
	}

	got := &Person{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age": $1, "fullname": $2}`, 42, "Bob Foobar"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": $1, "fullname": $2}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar")))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": "$1", "fullname": "$2"}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar")))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": $age, "fullname": $name}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar"))))
	fmt.Println("check got with named placeholders:", ok)

	got.Children = []*Person{
		{Fullname: "Alice", Age: 28},
		{Fullname: "Brian", Age: 22},
	}
	ok = td.Cmp(t, got,
		td.JSON(`{"age": $age, "fullname": $name, "children": $children}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar")),
			td.Tag("children", td.Bag(
				&Person{Fullname: "Brian", Age: 22},
				&Person{Fullname: "Alice", Age: 28},
			))))
	fmt.Println("check got w/named placeholders, and children w/go structs:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": Between($1, $2), "fullname": HasSuffix($suffix), "children": Len(2)}`,
			40, 45,
			td.Tag("suffix", "Foobar")))
	fmt.Println("check got w/num & named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got w/named placeholders, and children w/go structs: true
	// check got w/num & named placeholders: true
}

func ExampleJSON_embedding() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age": NotZero(), "fullname": NotEmpty()}`))
	fmt.Println("check got with simple operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"age": $^NotZero, "fullname": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between(40, 42, "]]"), // in ]40; 42]
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar")  // ← comma is optional here
  )
}`))
	fmt.Println("check got with complex operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between(40, 42, "]["), // in ]40; 42[ → 42 excluded
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar"),
  )
}`))
	fmt.Println("check got with complex operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between($1, $2, $3), // in ]40; 42]
  "fullname": All(
    HasPrefix($4),
    HasSuffix("bar")  // ← comma is optional here
  )
}`,
		40, 42, td.BoundsOutIn,
		"Bob"))
	fmt.Println("check got with complex operators, w/placeholder args:", ok)

	// Output:
	// check got with simple operators: true
	// check got with operator shortcuts: true
	// check got with complex operators: true
	// check got with complex operators: false
	// check got with complex operators, w/placeholder args: true
}

func ExampleJSON_rawStrings() {
	t := &testing.T{}

	type details struct {
		Address string `json:"address"`
		Car     string `json:"car"`
	}

	got := &struct {
		Fullname string  `json:"fullname"`
		Age      int     `json:"age"`
		Details  details `json:"details"`
	}{
		Fullname: "Foo Bar",
		Age:      42,
		Details: details{
			Address: "something",
			Car:     "Peugeot",
		},
	}

	ok := td.Cmp(t, got,
		td.JSON(`
{
  "fullname": HasPrefix("Foo"),
  "age":      Between(41, 43),
  "details":  SuperMapOf({
    "address": NotEmpty, // () are optional when no parameters
    "car":     Any("Peugeot", "Tesla", "Jeep") // any of these
  })
}`))
	fmt.Println("Original:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`
{
  "fullname": "$^HasPrefix(\"Foo\")",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({\n\"address\": NotEmpty,\n\"car\": Any(\"Peugeot\", \"Tesla\", \"Jeep\")\n})"
}`))
	fmt.Println("JSON compliant:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`
{
  "fullname": "$^HasPrefix(\"Foo\")",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    \"address\": NotEmpty, // () are optional when no parameters
    \"car\":     Any(\"Peugeot\", \"Tesla\", \"Jeep\") // any of these
  })"
}`))
	fmt.Println("JSON multilines strings:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`
{
  "fullname": "$^HasPrefix(r<Foo>)",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    r<address>: NotEmpty, // () are optional when no parameters
    r<car>:     Any(r<Peugeot>, r<Tesla>, r<Jeep>) // any of these
  })"
}`))
	fmt.Println("Raw strings:", ok)

	// Output:
	// Original: true
	// JSON compliant: true
	// JSON multilines strings: true
	// Raw strings: true
}

func ExampleJSON_file() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.Cmp(t, got,
		td.JSON(filename,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.Cmp(t, got,
		td.JSON(file,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleJSONPointer_rfc6901() {
	t := &testing.T{}

	got := json.RawMessage(`
{
   "foo":  ["bar", "baz"],
   "":     0,
   "a/b":  1,
   "c%d":  2,
   "e^f":  3,
   "g|h":  4,
   "i\\j": 5,
   "k\"l": 6,
   " ":    7,
   "m~n":  8
}`)

	expected := map[string]any{
		"foo": []any{"bar", "baz"},
		"":    0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		`i\j`: 5,
		`k"l`: 6,
		" ":   7,
		"m~n": 8,
	}
	ok := td.Cmp(t, got, td.JSONPointer("", expected))
	fmt.Println("Empty JSON pointer means all:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/foo`, []any{"bar", "baz"}))
	fmt.Println("Extract `foo` key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/foo/0`, "bar"))
	fmt.Println("First item of `foo` key slice:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/`, 0))
	fmt.Println("Empty key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/a~1b`, 1))
	fmt.Println("Slash has to be escaped using `~1`:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/c%d`, 2))
	fmt.Println("% in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/e^f`, 3))
	fmt.Println("^ in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/g|h`, 4))
	fmt.Println("| in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/i\j`, 5))
	fmt.Println("Backslash in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/k"l`, 6))
	fmt.Println("Double-quote in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/ `, 7))
	fmt.Println("Space key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/m~0n`, 8))
	fmt.Println("Tilde has to be escaped using `~0`:", ok)

	// Output:
	// Empty JSON pointer means all: true
	// Extract `foo` key: true
	// First item of `foo` key slice: true
	// Empty key: true
	// Slash has to be escaped using `~1`: true
	// % in key: true
	// ^ in key: true
	// | in key: true
	// Backslash in key: true
	// Double-quote in key: true
	// Space key: true
	// Tilde has to be escaped using `~0`: true
}

func ExampleJSONPointer_struct() {
	t := &testing.T{}

	// Without json tags, encoding/json uses public fields name
	type Item struct {
		Name  string
		Value int64
		Next  *Item
	}

	got := Item{
		Name:  "first",
		Value: 1,
		Next: &Item{
			Name:  "second",
			Value: 2,
			Next: &Item{
				Name:  "third",
				Value: 3,
			},
		},
	}

	ok := td.Cmp(t, got, td.JSONPointer("/Next/Next/Name", "third"))
	fmt.Println("3rd item name is `third`:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next/Value", td.Gte(int64(3))))
	fmt.Println("3rd item value is greater or equal than 3:", ok)

	ok = td.Cmp(t, got,
		td.JSONPointer("/Next",
			td.JSONPointer("/Next",
				td.JSONPointer("/Value", td.Gte(int64(3))))))
	fmt.Println("3rd item value is still greater or equal than 3:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next/Next/Name", td.Ignore()))
	fmt.Println("4th item exists and has a name:", ok)

	// Struct comparison work with or without pointer: &Item{…} works too
	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next", Item{
		Name:  "third",
		Value: 3,
	}))
	fmt.Println("3rd item full comparison:", ok)

	// Output:
	// 3rd item name is `third`: true
	// 3rd item value is greater or equal than 3: true
	// 3rd item value is still greater or equal than 3: true
	// 4th item exists and has a name: false
	// 3rd item full comparison: true
}

func ExampleJSONPointer_has_hasnt() {
	t := &testing.T{}

	got := json.RawMessage(`
{
  "name": "Bob",
  "age": 42,
  "children": [
    {
      "name": "Alice",
      "age": 16
    },
    {
      "name": "Britt",
      "age": 21,
      "children": [
        {
          "name": "John",
          "age": 1
        }
      ]
    }
  ]
}`)

	// Has Bob some children?
	ok := td.Cmp(t, got, td.JSONPointer("/children", td.Len(td.Gt(0))))
	fmt.Println("Bob has at least one child:", ok)

	// But checking "children" exists is enough here
	ok = td.Cmp(t, got, td.JSONPointer("/children/0/children", td.Ignore()))
	fmt.Println("Alice has children:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/children/1/children", td.Ignore()))
	fmt.Println("Britt has children:", ok)

	// The reverse can be checked too
	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/0/children", td.Ignore())))
	fmt.Println("Alice hasn't children:", ok)

	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/1/children", td.Ignore())))
	fmt.Println("Britt hasn't children:", ok)

	// Output:
	// Bob has at least one child: true
	// Alice has children: false
	// Britt has children: true
	// Alice hasn't children: true
	// Britt hasn't children: false
}

func ExampleKeys() {
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := td.Cmp(t, got, td.Keys([]string{"bar", "foo", "zip"}))
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = td.Cmp(t, got, td.Keys([]string{"zip", "bar", "foo"}))
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.Cmp(t, got, td.Keys(td.Bag("zip", "bar", "foo")))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = td.Cmp(t, got, td.Keys(td.ArrayEach(td.Len(3))))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true
}

func ExampleLast_classic() {
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.Last(td.Lt(0), -1))
	fmt.Println("last negative number is -1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.Last(isEven, 2))
	fmt.Println("last even number is 2:", ok)

	ok = td.Cmp(t, got, td.Last(isEven, td.Gt(0)))
	fmt.Println("last even number is > 0:", ok)

	ok = td.Cmp(t, got, td.Last(isEven, td.Code(isEven)))
	fmt.Println("last even number is well even:", ok)

	// Output:
	// last negative number is -1: true
	// last even number is 2: true
	// last even number is > 0: true
	// last even number is well even: true
}

func ExampleLast_empty() {
	t := &testing.T{}

	ok := td.Cmp(t, ([]int)(nil), td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in nil slice:", ok)

	ok = td.Cmp(t, []int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty slice:", ok)

	ok = td.Cmp(t, &[]int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty pointed slice:", ok)

	ok = td.Cmp(t, [0]int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty array:", ok)

	// Output:
	// last in nil slice: false
	// last in empty slice: false
	// last in empty pointed slice: false
	// last in empty array: false
}

func ExampleLast_struct() {
	t := &testing.T{}

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      37,
		},
	}

	ok := td.Cmp(t, got, td.Last(
		td.Smuggle("Age", td.Gt(30)),
		td.Smuggle("Fullname", "Alice Bingo")))
	fmt.Println("last person.Age > 30 → Alice:", ok)

	ok = td.Cmp(t, got, td.Last(
		td.JSONPointer("/age", td.Gt(30)),
		td.SuperJSONOf(`{"fullname":"Alice Bingo"}`)))
	fmt.Println("last person.Age > 30 → Alice, using JSON:", ok)

	ok = td.Cmp(t, got, td.Last(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSONPointer("/fullname", td.HasPrefix("Alice"))))
	fmt.Println("first person.Age > 30 → Alice, using JSONPointer:", ok)

	// Output:
	// last person.Age > 30 → Alice: true
	// last person.Age > 30 → Alice, using JSON: true
	// first person.Age > 30 → Alice, using JSONPointer: true
}

func ExampleLast_json() {
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": Last(Lt(3), 2)}`))
	fmt.Println("last number < 3:", ok)

	got = map[string]any{
		"persons": []map[string]any{
			{"id": 1, "name": "Joe"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Alice"},
			{"id": 4, "name": "Brian"},
			{"id": 5, "name": "Britt"},
		},
	}
	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": Last(JSONPointer("/name", "Brian"), {"id": 4, "name": "Brian"})
}`))
	fmt.Println(`is "Brian" content OK:`, ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": Last(JSONPointer("/name", "Brian"), JSONPointer("/id", 4))
}`))
	fmt.Println(`ID of "Brian" is 4:`, ok)

	// Output:
	// last number < 3: true
	// is "Brian" content OK: true
	// ID of "Brian" is 4: true
}

func ExampleLax() {
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

	ok := td.Cmp(t, gotInt64, td.Lax(expected))
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = td.Cmp(t, gotInt32, td.Lax(expected))
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = td.Cmp(t, gotMyInt, td.Lax(expected))
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true
}

func ExampleLen_slice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := td.Cmp(t, got, td.Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.Cmp(t, got, td.Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_map() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := td.Cmp(t, got, td.Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.Cmp(t, got, td.Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_operatorSlice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := td.Cmp(t, got, td.Len(td.Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Len(td.Lt(5)), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLen_operatorMap() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := td.Cmp(t, got, td.Len(td.Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Len(td.Gte(3)), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLt_int() {
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Lt(157), "checks %v is < 157", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lt(156), "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleLt_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Lt("abd"), `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lt("abc"), `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleLte_int() {
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Lte(156), "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte(157), "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte(155), "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleLte_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Lte("abc"), `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte("abd"), `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte("abb"), `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleMap_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.Map(map[string]int{"bar": 42},
			td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Map(map[string]int{},
			td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Map((map[string]int)(nil),
			td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleMap_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.Map(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map(&MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map((*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleMapEach_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got, td.MapEach(td.Between(10, 90)),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleMapEach_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got, td.MapEach(td.Between(10, 90)),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.MapEach(td.Between(10, 90)),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleN() {
	t := &testing.T{}

	got := 1.12345

	ok := td.Cmp(t, got, td.N(1.1234, 0.00006),
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleNaN_float32() {
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false
}

func ExampleNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false
}

func ExampleNil() {
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := td.Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = td.Cmp(t, got, td.Nil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = td.Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = td.Cmp(t, got, td.Nil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleNone() {
	t := &testing.T{}

	got := 18

	ok := td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = td.Cmp(t, got, td.None(prime, even, td.Gt(14)),
			"checks %v is not prime number, nor an even number and not > 14")
		fmt.Printf("%d → %t\n", got, ok)
	}

	// Output:
	// true
	// false
	// false
	// 9 → true
	// 3 → false
	// 8 → false
	// 15 → false
}

func ExampleNotAny() {
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.Cmp(t, got, td.NotAny(3, 6, 8, 41, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.NotAny(3, 6, 8, 42, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []any slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = td.Cmp(t, got, td.NotAny(td.Flatten(notExpected)),
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleNot() {
	t := &testing.T{}

	got := 42

	ok := td.Cmp(t, got, td.Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Not(td.Between(10, 30)),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = td.Cmp(t, got, td.Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleNotEmpty() {
	t := &testing.T{}

	ok := td.Cmp(t, nil, td.NotEmpty()) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = td.Cmp(t, "foobar", td.NotEmpty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = td.Cmp(t, 0, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{"foobar": 42}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{1}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.NotEmpty()) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true
}

func ExampleNotEmpty_pointers() {
	t := &testing.T{}

	type MySlice []int

	ok := td.Cmp(t, MySlice{12}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, &MySlice{12}, td.NotEmpty()) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = td.Cmp(t, &l3, td.NotEmpty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.Cmp(t, &MyStruct{}, td.NotEmpty()) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNotNaN_float32() {
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
}

func ExampleNotNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
}

func ExampleNotNil() {
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = td.Cmp(t, got, td.NotNil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = td.Cmp(t, got, td.NotNil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNotZero() {
	t := &testing.T{}

	ok := td.Cmp(t, 0, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, float64(0), td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, 12, td.NotZero())
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{0, 1}, td.NotZero()) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.Cmp(t, bytes.Buffer{}, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.NotZero()) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.NotZero())) // fails as deref by Ptr()
	fmt.Println(ok)

	// Output:
	// false
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
}

func ExamplePPtr() {
	t := &testing.T{}

	num := 12
	got := &num

	ok := td.Cmp(t, &got, td.PPtr(12))
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.PPtr(td.Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExamplePtr() {
	t := &testing.T{}

	got := 12

	ok := td.Cmp(t, &got, td.Ptr(12))
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.Ptr(td.Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleRe() {
	t := &testing.T{}

	got := "foo bar"
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_error() {
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.Re(`^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.Cmp(t, got, td.Re(`^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.ReAll(`(\w+)`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.Cmp(t, got, td.ReAll(`(\w+)`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_captureComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := td.Cmp(t, got,
		td.ReAll(`(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.Cmp(t, got,
		td.ReAll(`(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_compiled() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_compiledStringer() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledError() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.Re(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.Cmp(t, got, td.Re(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.ReAll(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.Cmp(t, got, td.ReAll(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_compiledCaptureComplex() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := td.Cmp(t, got,
		td.ReAll(expected, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.Cmp(t, got,
		td.ReAll(expected, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRecv_basic() {
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.Cmp(t, got, td.Recv(1))
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false
}

func ExampleRecv_channelPointer() {
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.Cmp(t, &got, td.Recv(1))
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, &got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false
}

func ExampleRecv_withTimeout() {
	t := &testing.T{}

	got := make(chan int, 1)
	tick := make(chan struct{})

	go func() {
		// ①
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 0

		// ②
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 1

		// ③
		<-tick
		time.Sleep(100 * time.Millisecond)
		close(got)
	}()

	td.Cmp(t, got, td.Recv(td.RecvNothing))

	// ①
	tick <- struct{}{}
	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("① RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(0, 150*time.Millisecond))
	fmt.Println("① receive 0 w/150ms timeout:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("① RecvNothing:", ok)

	// ②
	tick <- struct{}{}
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("② RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(1, 150*time.Millisecond))
	fmt.Println("② receive 1 w/150ms timeout:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("② RecvNothing:", ok)

	// ③
	tick <- struct{}{}
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("③ RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvClosed, 150*time.Millisecond))
	fmt.Println("③ check closed w/150ms timeout:", ok)

	// Output:
	// ① RecvNothing: true
	// ① receive 0 w/150ms timeout: true
	// ① RecvNothing: true
	// ② RecvNothing: true
	// ② receive 1 w/150ms timeout: true
	// ② RecvNothing: true
	// ③ RecvNothing: true
	// ③ check closed w/150ms timeout: true
}

func ExampleRecv_nilChannel() {
	t := &testing.T{}

	var ch chan int

	ok := td.Cmp(t, ch, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive from nil channel:", ok)

	ok = td.Cmp(t, ch, td.Recv(42))
	fmt.Println("something to receive from nil channel:", ok)

	ok = td.Cmp(t, ch, td.Recv(td.RecvClosed))
	fmt.Println("is a nil channel closed:", ok)

	// Output:
	// nothing to receive from nil channel: true
	// something to receive from nil channel: false
	// is a nil channel closed: false
}

func ExampleSet() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := td.Cmp(t, got, td.Set(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = td.Cmp(t, got, td.Set(1, 2, 2, 2, 2, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = td.Cmp(t, got, td.Set(td.Between(1, 4), 3, td.Between(2, 10)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = td.Cmp(t, got, td.Set(td.Flatten(expected)),
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleShallow() {
	t := &testing.T{}

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := td.Cmp(t, got, td.Shallow(&data),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = td.Cmp(t, got, td.Shallow(&MyStruct{Value: 12}),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleShallow_slice() {
	t := &testing.T{}

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := td.Cmp(t, a, td.Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.Cmp(t, b, td.Shallow(back))
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleShallow_string() {
	t := &testing.T{}

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := td.Cmp(t, a, td.Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.Cmp(t, b, td.Shallow(a))
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.Cmp(t, got, td.Slice([]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Slice([]int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Slice(([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleSlice_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.Cmp(t, got, td.Slice(MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.Slice(&MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Slice(&MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Slice((*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleSuperSliceOf_array() {
	t := &testing.T{}

	got := [4]int{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf([4]int{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf([4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&[4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*[4]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true
}

func ExampleSuperSliceOf_typedArray() {
	t := &testing.T{}

	type MyArray [4]int

	got := MyArray{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf(MyArray{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf(MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*MyArray)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true
}

func ExampleSuperSliceOf_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf([]int{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf([]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&[]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*[]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true
}

func ExampleSuperSliceOf_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf(MySlice{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf(MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*MySlice)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true
}

func ExampleSmuggle_convert() {
	t := &testing.T{}

	got := int64(123)

	ok := td.Cmp(t, got,
		td.Smuggle(func(n int64) int { return int(n) }, 123),
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, bool) {
				n, err := strconv.Atoi(numStr)
				return n, err == nil
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, bool, string) {
				n, err := strconv.Atoi(numStr)
				if err != nil {
					return 0, false, "string must contain a number"
				}
				return n, true, ""
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, error) { //nolint: gocritic
				return strconv.Atoi(numStr)
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = td.Cmp(t, "123",
		td.Smuggle(strconv.Atoi, td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true
}

func ExampleSmuggle_lax() {
	t := &testing.T{}

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := td.Cmp(t, got,
		td.Smuggle(func(n int64) uint32 { return uint32(n) }, uint32(123)))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true
}

func ExampleSmuggle_auto_unmarshal() {
	t := &testing.T{}

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := td.Cmp(t, got,
		td.Smuggle(
			func(b json.RawMessage) (r map[string]int, err error) {
				err = json.Unmarshal(b, &r)
				return
			},
			map[string]int{
				"a": 1,
				"b": 2,
			}))
	fmt.Println("JSON contents is OK:", ok)

	// Output:
	// JSON contents is OK: true
}

func ExampleSmuggle_cast() {
	t := &testing.T{}

	// A string containing JSON
	got := `{ "foo": 123 }`

	// Automatically cast a string to a json.RawMessage so td.JSON can operate
	ok := td.Cmp(t, got,
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":123}`)))
	fmt.Println("JSON contents in string is OK:", ok)

	// Automatically read from io.Reader to a json.RawMessage
	ok = td.Cmp(t, bytes.NewReader([]byte(got)),
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":123}`)))
	fmt.Println("JSON contents just read is OK:", ok)

	// Output:
	// JSON contents in string is OK: true
	// JSON contents just read is OK: true
}

func ExampleSmuggle_complex() {
	t := &testing.T{}

	// No end date but a start date and a duration
	type StartDuration struct {
		StartDate time.Time
		Duration  time.Duration
	}

	// Checks that end date is between 17th and 19th February both at 0h
	// for each of these durations in hours

	for _, duration := range []time.Duration{48 * time.Hour, 72 * time.Hour, 96 * time.Hour} {
		got := StartDuration{
			StartDate: time.Date(2018, time.February, 14, 12, 13, 14, 0, time.UTC),
			Duration:  duration,
		}

		// Simplest way, but in case of Between() failure, error will be bound
		// to DATA<smuggled>, not very clear...
		ok := td.Cmp(t, got,
			td.Smuggle(
				func(sd StartDuration) time.Time {
					return sd.StartDate.Add(sd.Duration)
				},
				td.Between(
					time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC))))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = td.Cmp(t, got,
			td.Smuggle(
				func(sd StartDuration) td.SmuggledGot {
					return td.SmuggledGot{
						Name: "ComputedEndDate",
						Got:  sd.StartDate.Add(sd.Duration),
					}
				},
				td.Between(
					time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC))))
		fmt.Println(ok)
	}

	// Output:
	// false
	// false
	// true
	// true
	// true
	// true
}

func ExampleSmuggle_interface() {
	t := &testing.T{}

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := td.Cmp(t, gotTime,
		td.Smuggle(func(s fmt.Stringer) string {
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = td.Cmp(t, MyTime(gotTime),
		td.Smuggle(func(s fmt.Stringer) string {
			fmt.Println("Smuggle func called!")
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified MyTime OK:", ok)

	// Output:
	// stringified time.Time OK: true
	// stringified MyTime OK: false
}

func ExampleSmuggle_field_path() {
	t := &testing.T{}

	type Body struct {
		Name  string
		Value any
	}
	type Request struct {
		Body *Body
	}
	type Transaction struct {
		Request
	}
	type ValueNum struct {
		Num int
	}

	got := &Transaction{
		Request: Request{
			Body: &Body{
				Name:  "test",
				Value: &ValueNum{Num: 123},
			},
		},
	}

	// Want to check whether Num is between 100 and 200?
	ok := td.Cmp(t, got,
		td.Smuggle(
			func(t *Transaction) (int, error) {
				if t.Request.Body == nil ||
					t.Request.Body.Value == nil {
					return 0, errors.New("Request.Body or Request.Body.Value is nil")
				}
				if v, ok := t.Request.Body.Value.(*ValueNum); ok && v != nil {
					return v.Num, nil
				}
				return 0, errors.New("Request.Body.Value isn't *ValueNum or nil")
			},
			td.Between(100, 200)))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = td.Cmp(t, got, td.Smuggle("Request.Body.Value.Num", td.Between(100, 200)))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = td.Cmp(t, got, td.Smuggle("Body.Value.Num", td.Between(100, 200)))
	fmt.Println("check Num using an other fields-path:", ok)

	// Note that maps and array/slices are supported
	got.Request.Body.Value = map[string]any{
		"foo": []any{
			3: map[int]string{666: "bar"},
		},
	}
	ok = td.Cmp(t, got, td.Smuggle("Body.Value[foo][3][666]", "bar"))
	fmt.Println("check fields-path including maps/slices:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
	// check fields-path including maps/slices: true
}

func ExampleString() {
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.String("foobar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleString_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleString_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix() {
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleHasPrefix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix() {
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleHasSuffix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleStruct() {
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := td.Cmp(t, got,
		td.Struct(Person{Name: "Foobar"}, td.StructFields{
			"Age": td.Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	ok = td.Cmp(t, got,
		td.Struct(Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = td.Cmp(t, &got,
		td.Struct(&Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = td.Cmp(t, &got,
		td.Struct((*Person)(nil), td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true
}

func ExampleStruct_overwrite_model() {
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := td.Cmp(t, got,
		td.Struct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = td.Cmp(t, got,
		td.Struct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true
}

func ExampleStruct_patterns() {
	t := &testing.T{}

	type Person struct {
		Firstname string
		Lastname  string
		Surname   string
		Nickname  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	now := time.Now()
	got := Person{
		Firstname: "Maxime",
		Lastname:  "Foo",
		Surname:   "Max",
		Nickname:  "max",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // not deleted yet
	}

	ok := td.Cmp(t, got,
		td.Struct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`: nil,
			`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		}),
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = td.Cmp(t, got,
		td.Struct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`:  nil,
			`1 =  *name`: td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`2 =~ At\z`:  td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		}),
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true
}

func ExampleStruct_struct_fields() { // only operator
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := td.Cmp(t, got, td.Struct(Person{Name: "Foobar"}), "no StructFields")
	fmt.Println("Without any StructFields:", ok)

	ok = td.Cmp(t, got,
		td.Struct(Person{Name: "Bingo"},
			td.StructFields{
				"> Name": "pipo",
				"Age":    42,
			},
			td.StructFields{
				"> Name":      "bingo",
				"NumChildren": 10,
			},
			td.StructFields{
				">Name":       "Foobar",
				"NumChildren": 3,
			}),
		"merge several StructFields")
	fmt.Println("Merge several StructFields:", ok)

	// Output:
	// Without any StructFields: true
	// Merge several StructFields: true
}

func ExampleStruct_lazy_model() {
	t := &testing.T{}

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := td.Cmp(t, got, td.Struct(nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	}))
	fmt.Println("Lazy model:", ok)

	ok = td.Cmp(t, got, td.Struct(nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	}))
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false
}

func ExampleSStruct() {
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := td.Cmp(t, got,
		td.SStruct(Person{Name: "Foobar"}, td.StructFields{
			"Age": td.Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	got.NumChildren = 3
	ok = td.Cmp(t, got,
		td.SStruct(Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = td.Cmp(t, &got,
		td.SStruct(&Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = td.Cmp(t, &got,
		td.SStruct((*Person)(nil), td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true
}

func ExampleSStruct_overwrite_model() {
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := td.Cmp(t, got,
		td.SStruct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = td.Cmp(t, got,
		td.SStruct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true
}

func ExampleSStruct_patterns() {
	t := &testing.T{}

	type Person struct {
		Firstname string
		Lastname  string
		Surname   string
		Nickname  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
		id        int64
		secret    string
	}

	now := time.Now()
	got := Person{
		Firstname: "Maxime",
		Lastname:  "Foo",
		Surname:   "Max",
		Nickname:  "max",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // not deleted yet
		id:        2345,
		secret:    "5ecr3T",
	}

	ok := td.Cmp(t, got,
		td.SStruct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`: nil,
			`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
			`!  [A-Z]*`: td.Ignore(),        // private fields
		}),
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = td.Cmp(t, got,
		td.SStruct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`:   nil,
			`1 =  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`2 =~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
			`3 !~ ^[A-Z]`: td.Ignore(),        // private fields
		}),
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true
}

func ExampleSStruct_struct_fields() { // only operator
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// No added value here, but it works
	ok := td.Cmp(t, got,
		td.SStruct(Person{
			Name:        "Foobar",
			Age:         42,
			NumChildren: 3,
		}),
		"no StructFields")
	fmt.Println("Without any StructFields:", ok)

	ok = td.Cmp(t, got,
		td.SStruct(Person{Name: "Bingo"},
			td.StructFields{
				"> Name": "pipo",
				"Age":    42,
			},
			td.StructFields{
				"> Name":      "bingo",
				"NumChildren": 10,
			},
			td.StructFields{
				">Name":       "Foobar",
				"NumChildren": 3,
			}),
		"merge several StructFields")
	fmt.Println("Merge several StructFields:", ok)

	// Output:
	// Without any StructFields: true
	// Merge several StructFields: true
}

func ExampleSStruct_lazy_model() {
	t := &testing.T{}

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := td.Cmp(t, got, td.SStruct(nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	}))
	fmt.Println("Lazy model:", ok)

	ok = td.Cmp(t, got, td.SStruct(nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	}))
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false
}

func ExampleSubBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = td.Cmp(t, got, td.SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = td.Cmp(t, got, td.SubBagOf(
		td.Between(0, 3),
		td.Between(0, 3),
		td.Between(0, 3),
		td.Between(0, 3),
		td.Gt(4),
		td.Gt(4)),
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 9, 8}
	ok = td.Cmp(t, got, td.SubBagOf(td.Flatten(expected)),
		"checks at least all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// true
}

func ExampleSubJSONOf_basic() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.SubJSONOf(`{"age":42,"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = td.Cmp(t, got, td.SubJSONOf(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = td.Cmp(t, got, td.SubJSONOf(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // This field is ignored as SubJSONOf
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.Cmp(t, got, td.SubJSONOf(`{"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got without age field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got without age field: false
}

func ExampleSubJSONOf_placeholders() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.Cmp(t, got,
		td.SubJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			42, "Bob Foobar", "male"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.Cmp(t, got,
		td.SubJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SubJSONOf(`{"age": "$1", "fullname": "$2", "gender": "$3"}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SubJSONOf(`{"age": $age, "fullname": $name, "gender": $gender}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar")),
			td.Tag("gender", td.NotEmpty())))
	fmt.Println("check got with named placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SubJSONOf(`{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleSubJSONOf_file() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender",
  "details":  {
    "city": "TestCity",
    "zip":  666
  }
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.Cmp(t, got,
		td.SubJSONOf(filename,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.Cmp(t, got,
		td.SubJSONOf(file,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleSubMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := td.Cmp(t, got,
		td.SubMapOf(map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleSubMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := td.Cmp(t, got,
		td.SubMapOf(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.SubMapOf(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSubSetOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := td.Cmp(t, got, td.SubSetOf(1, 2, 3, 4, 5, 6, 7, 8),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = td.Cmp(t, got, td.SubSetOf(td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ok = td.Cmp(t, got, td.SubSetOf(td.Flatten(expected)),
		"checks at least all expected items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleSuperBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SuperBagOf(8, 5, 8),
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.SuperBagOf(td.Gt(5), td.Lte(2)),
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = td.Cmp(t, got, td.SuperBagOf(td.Flatten(expected)),
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleSuperJSONOf_basic() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.Cmp(t, got, td.SuperJSONOf(`{"age":42,"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = td.Cmp(t, got, td.SuperJSONOf(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = td.Cmp(t, got, td.SuperJSONOf(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // The gender!
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"fullname":"Bob","gender":"male","details":{}}`))
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false
}

func ExampleSuperJSONOf_placeholders() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			42, "Bob Foobar", "male"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": "$1", "fullname": "$2", "gender": "$3"}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $age, "fullname": $name, "gender": $gender}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar")),
			td.Tag("gender", td.NotEmpty())))
	fmt.Println("check got with named placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleSuperJSONOf_file() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.Cmp(t, got,
		td.SuperJSONOf(filename,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.Cmp(t, got,
		td.SuperJSONOf(file,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleSuperMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.SuperMapOf(map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks map %v contains at least all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleSuperMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.SuperMapOf(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks typed map %v contains at least all expected keys/values", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.SuperMapOf(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks pointed typed map %v contains at least all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSuperSetOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SuperSetOf(1, 2, 3),
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.SuperSetOf(td.Gt(5), td.Lte(2)),
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = td.Cmp(t, got, td.SuperSetOf(td.Flatten(expected)),
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
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

	// Compare dates ignoring nanoseconds and monotonic parts
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := td.Cmp(t, got, td.TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = td.Cmp(t, got, td.TruncTime(expected, 24*time.Hour),
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = td.Cmp(t, got, td.TruncTime(expected),
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleValues() {
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := td.Cmp(t, got, td.Values([]int{1, 2, 3}))
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = td.Cmp(t, got, td.Values([]int{3, 1, 2}))
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.Cmp(t, got, td.Values(td.Bag(3, 1, 2)))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = td.Cmp(t, got, td.Values(td.ArrayEach(td.Between(1, 3))))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true
}

func ExampleZero() {
	t := &testing.T{}

	ok := td.Cmp(t, 0, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, float64(0), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, 12, td.Zero()) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{0, 1}, td.Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.Cmp(t, bytes.Buffer{}, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Zero()) // fails, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.Zero())) // OK with the help of Ptr()
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
}
