// Copyright (c) 2018, 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

package td_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
)

func ExampleCmpAll() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := td.CmpAll(t, got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = td.CmpAll(t, got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []interface{} copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = td.CmpAll(t, got, []interface{}{td.HasPrefix("foo"), regOps, td.HasSuffix("bar")},
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpAny() {
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := td.CmpAny(t, got, []interface{}{td.Re("zip"), td.HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = td.CmpAny(t, got, []interface{}{td.Re("zip"), td.HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []interface{} copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("a/c"), td.Re(`^xx`), td.Re(`ar$`)})
	ok = td.CmpAny(t, got, []interface{}{td.HasPrefix("xxx"), regOps, td.HasSuffix("zip")},
		"check at least one operator matches value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpArray_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.CmpArray(t, got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArray_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.CmpArray(t, got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpArrayEach_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.CmpArrayEach(t, got, td.Between(25, 60),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.CmpArrayEach(t, got, td.Between(25, 60),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = td.CmpArrayEach(t, &got, td.Between(25, 60),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.CmpArrayEach(t, got, td.Between(25, 60),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.CmpArrayEach(t, got, td.Between(25, 60),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = td.CmpArrayEach(t, &got, td.Between(25, 60),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpBag() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := td.CmpBag(t, got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = td.CmpBag(t, got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = td.CmpBag(t, got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = td.CmpBag(t, got, []interface{}{1, 2, 3, 5, td.Gt(7)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5}
	ok = td.CmpBag(t, got, []interface{}{td.Flatten(expected), td.Gt(7)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
	// true
}

func ExampleCmpBetween_int() {
	t := &testing.T{}

	got := 156

	ok := td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleCmpBetween_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleCmpCap() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.CmpCap(t, got, 12, "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = td.CmpCap(t, got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.CmpCap(t, got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpCap_operator() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.CmpCap(t, got, td.Between(10, 12),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = td.CmpCap(t, got, td.Gt(10),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpCode() {
	t := &testing.T{}

	got := "12"

	ok := td.CmpCode(t, got, func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = td.CmpCode(t, got, func(num string) (bool, string) {
		n, err := strconv.Atoi(num)
		if err != nil {
			return false, "not a number"
		}
		if n > 10 && n < 100 {
			return true, ""
		}
		return false, "not in ]10 .. 100["
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason thanks to error
	ok = td.CmpCode(t, got, func(num string) error {
		n, err := strconv.Atoi(num)
		if err != nil {
			return err
		}
		if n > 10 && n < 100 {
			return nil
		}
		return fmt.Errorf("%d not in ]10 .. 100[", n)
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpContains_arraySlice() {
	t := &testing.T{}

	ok := td.CmpContains(t, [...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = td.CmpContains(t, [...]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, []int{22, 33})
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true
}

func ExampleCmpContains_nil() {
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := td.CmpContains(t, got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = td.CmpContains(t, got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = td.CmpContains(t, got, td.Nil())
	fmt.Println("array contains Nil():", ok)

	ok = td.CmpContains(t, got, (*byte)(nil))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false
}

func ExampleCmpContains_map() {
	t := &testing.T{}

	ok := td.CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = td.CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true
}

func ExampleCmpContains_string() {
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, []byte("oob"), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleCmpContains_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleCmpContains_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleCmpContainsKey() {
	t := &testing.T{}

	ok := td.CmpContainsKey(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, "foo")
	fmt.Println(`map contains key "foo":`, ok)

	ok = td.CmpContainsKey(t, map[int]bool{12: true, 24: false, 42: true, 51: false}, td.Between(40, 50))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = td.CmpContainsKey(t, map[string]int{"FOO": 11, "bar": 22, "zip": 33}, td.Smuggle(strings.ToLower, "foo"))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true
}

func ExampleCmpContainsKey_nil() {
	t := &testing.T{}

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := td.CmpContainsKey(t, got, nil)
	fmt.Println("map contains untyped nil key:", ok)

	ok = td.CmpContainsKey(t, got, (*int)(nil))
	fmt.Println("map contains *int nil key:", ok)

	ok = td.CmpContainsKey(t, got, td.Nil())
	fmt.Println("map contains Nil() key:", ok)

	ok = td.CmpContainsKey(t, got, (*byte)(nil))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false
}

func ExampleCmpEmpty() {
	t := &testing.T{}

	ok := td.CmpEmpty(t, nil) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = td.CmpEmpty(t, (*int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, "")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = td.CmpEmpty(t, 0)
	fmt.Println(ok)

	ok = td.CmpEmpty(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, map[string]int{})
	fmt.Println(ok)

	ok = td.CmpEmpty(t, ([]int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, []int{})
	fmt.Println(ok)

	ok = td.CmpEmpty(t, []int{3}) // fails, as not empty
	fmt.Println(ok)

	ok = td.CmpEmpty(t, [3]int{}) // fails, Empty() is not Zero()!
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

func ExampleCmpEmpty_pointers() {
	t := &testing.T{}

	type MySlice []int

	ok := td.CmpEmpty(t, MySlice{}) // Ptr() not needed
	fmt.Println(ok)

	ok = td.CmpEmpty(t, &MySlice{})
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = td.CmpEmpty(t, &l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.CmpEmpty(t, &MyStruct{}) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleCmpGt_int() {
	t := &testing.T{}

	got := 156

	ok := td.CmpGt(t, got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = td.CmpGt(t, got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpGt_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.CmpGt(t, got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = td.CmpGt(t, got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpGte_int() {
	t := &testing.T{}

	got := 156

	ok := td.CmpGte(t, got, 156, "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = td.CmpGte(t, got, 155, "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = td.CmpGte(t, got, 157, "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleCmpGte_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.CmpGte(t, got, "abc", `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = td.CmpGte(t, got, "abb", `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = td.CmpGte(t, got, "abd", `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleCmpHasPrefix() {
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleCmpHasPrefix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpHasPrefix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpHasSuffix() {
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleCmpHasSuffix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpHasSuffix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpIsa() {
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := td.CmpIsa(t, got, TstStruct{}, "checks got is a TstStruct")
	fmt.Println(ok)

	ok = td.CmpIsa(t, got, &TstStruct{},
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = td.CmpIsa(t, &got, &TstStruct{},
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpIsa_interface() {
	t := &testing.T{}

	got := bytes.NewBufferString("foobar")

	ok := td.CmpIsa(t, got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = td.CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = td.CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = td.CmpIsa(t, &errGot, (*error)(nil),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleCmpJSON_basic() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.CmpJSON(t, got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = td.CmpJSON(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = td.CmpJSON(t, true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = td.CmpJSON(t, 42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = td.CmpJSON(t, got, `null`, nil)
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

func ExampleCmpJSON_placeholders() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.CmpJSON(t, got, `{"age": "$1", "fullname": "$2"}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.CmpJSON(t, got, `{"age": $age, "fullname": $name}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	ok = td.CmpJSON(t, got, `{"age": $^NotZero, "fullname": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleCmpJSON_file() {
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

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.CmpJSON(t, got, filename, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.CmpJSON(t, got, file, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleCmpKeys() {
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := td.CmpKeys(t, got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = td.CmpKeys(t, got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.CmpKeys(t, got, td.Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = td.CmpKeys(t, got, td.ArrayEach(td.Len(3)))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true
}

func ExampleCmpLax() {
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

	ok := td.CmpLax(t, gotInt64, expected)
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = td.CmpLax(t, gotInt32, expected)
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = td.CmpLax(t, gotMyInt, expected)
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true
}

func ExampleCmpLen_slice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := td.CmpLen(t, got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = td.CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpLen_map() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := td.CmpLen(t, got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = td.CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpLen_operatorSlice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := td.CmpLen(t, got, td.Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = td.CmpLen(t, got, td.Lt(5), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpLen_operatorMap() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := td.CmpLen(t, got, td.Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = td.CmpLen(t, got, td.Gte(3), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpLt_int() {
	t := &testing.T{}

	got := 156

	ok := td.CmpLt(t, got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = td.CmpLt(t, got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpLt_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.CmpLt(t, got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = td.CmpLt(t, got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpLte_int() {
	t := &testing.T{}

	got := 156

	ok := td.CmpLte(t, got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleCmpLte_string() {
	t := &testing.T{}

	got := "abc"

	ok := td.CmpLte(t, got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleCmpMap_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMap(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, got, map[string]int{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, got, (map[string]int)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpMap_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMap(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, &MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, (*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpMapEach_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMapEach(t, got, td.Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpMapEach_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMapEach(t, got, td.Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = td.CmpMapEach(t, &got, td.Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpN() {
	t := &testing.T{}

	got := 1.12345

	ok := td.CmpN(t, got, 1.1234, 0.00006,
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpNaN_float32() {
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false
}

func ExampleCmpNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false
}

func ExampleCmpNil() {
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := td.Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = td.CmpNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = td.Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = td.CmpNil(t, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleCmpNone() {
	t := &testing.T{}

	got := 18

	ok := td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = td.CmpNone(t, got, []interface{}{prime, even, td.Gt(14)},
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

func ExampleCmpNot() {
	t := &testing.T{}

	got := 42

	ok := td.CmpNot(t, got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = td.CmpNot(t, got, td.Between(10, 30),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = td.CmpNot(t, got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleCmpNotAny() {
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.CmpNotAny(t, got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.CmpNotAny(t, got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = td.CmpNotAny(t, got, []interface{}{td.Flatten(notExpected)},
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCmpNotEmpty() {
	t := &testing.T{}

	ok := td.CmpNotEmpty(t, nil) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, "foobar")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = td.CmpNotEmpty(t, 0)
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, map[string]int{"foobar": 42})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, []int{1})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, [3]int{}) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true
}

func ExampleCmpNotEmpty_pointers() {
	t := &testing.T{}

	type MySlice []int

	ok := td.CmpNotEmpty(t, MySlice{12})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, &MySlice{12}) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = td.CmpNotEmpty(t, &l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.CmpNotEmpty(t, &MyStruct{}) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleCmpNotNaN_float32() {
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
}

func ExampleCmpNotNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
}

func ExampleCmpNotNil() {
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = td.CmpNotNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = td.CmpNotNil(t, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleCmpNotZero() {
	t := &testing.T{}

	ok := td.CmpNotZero(t, 0) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, float64(0)) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, 12)
	fmt.Println(ok)

	ok = td.CmpNotZero(t, (map[string]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, map[string]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, ([]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, []int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, [3]int{}) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, [3]int{0, 1}) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.CmpNotZero(t, bytes.Buffer{}) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, &bytes.Buffer{}) // succeeds, as pointer not nil
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

func ExampleCmpPPtr() {
	t := &testing.T{}

	num := 12
	got := &num

	ok := td.CmpPPtr(t, &got, 12)
	fmt.Println(ok)

	ok = td.CmpPPtr(t, &got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpPtr() {
	t := &testing.T{}

	got := 12

	ok := td.CmpPtr(t, &got, 12)
	fmt.Println(ok)

	ok = td.CmpPtr(t, &got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpRe() {
	t := &testing.T{}

	got := "foo bar"
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpRe_error() {
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpRe_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_compiled() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpRe_compiledStringer() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpRe_compiledError() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpRe_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := td.CmpRe(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.CmpRe(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpReAll_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.CmpReAll(t, got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.CmpReAll(t, got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpReAll_captureComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := td.CmpReAll(t, got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.CmpReAll(t, got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpReAll_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := td.CmpReAll(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.CmpReAll(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpReAll_compiledCaptureComplex() {
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := td.CmpReAll(t, got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.CmpReAll(t, got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpSet() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := td.CmpSet(t, got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = td.CmpSet(t, got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = td.CmpSet(t, got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = td.CmpSet(t, got, []interface{}{td.Flatten(expected)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpShallow() {
	t := &testing.T{}

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := td.CmpShallow(t, got, &data,
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = td.CmpShallow(t, got, &MyStruct{Value: 12},
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpShallow_slice() {
	t := &testing.T{}

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := td.CmpShallow(t, a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.CmpShallow(t, b, back)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleCmpShallow_string() {
	t := &testing.T{}

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := td.CmpShallow(t, a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.CmpShallow(t, b, a)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleCmpSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.CmpSlice(t, got, []int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, got, []int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, got, ([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpSlice_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.CmpSlice(t, got, MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, &MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, &MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, (*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpSmuggle_convert() {
	t := &testing.T{}

	got := int64(123)

	ok := td.CmpSmuggle(t, got, func(n int64) int { return int(n) }, 123,
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, bool) {
		n, err := strconv.Atoi(numStr)
		return n, err == nil
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, bool, string) {
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, false, "string must contain a number"
		}
		return n, true, ""
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, error) {
		return strconv.Atoi(numStr)
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = td.CmpSmuggle(t, "123", strconv.Atoi, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true
}

func ExampleCmpSmuggle_lax() {
	t := &testing.T{}

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := td.CmpSmuggle(t, got, func(n int64) uint32 { return uint32(n) }, uint32(123))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true
}

func ExampleCmpSmuggle_auto_unmarshal() {
	t := &testing.T{}

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := td.CmpSmuggle(t, got, func(b json.RawMessage) (r map[string]int, err error) {
		err = json.Unmarshal(b, &r)
		return
	}, map[string]int{
		"a": 1,
		"b": 2,
	})
	fmt.Println("JSON contents is OK:", ok)

	// Output:
	// JSON contents is OK: true
}

func ExampleCmpSmuggle_complex() {
	t := &testing.T{}

	// No end date but a start date and a duration
	type StartDuration struct {
		StartDate time.Time
		Duration  time.Duration
	}

	// Checks that end date is between 17th and 19th February both at 0h
	// for each of these durations in hours

	for _, duration := range []time.Duration{48, 72, 96} {
		got := StartDuration{
			StartDate: time.Date(2018, time.February, 14, 12, 13, 14, 0, time.UTC),
			Duration:  duration * time.Hour,
		}

		// Simplest way, but in case of Between() failure, error will be bound
		// to DATA<smuggled>, not very clear...
		ok := td.CmpSmuggle(t, got, func(sd StartDuration) time.Time {
			return sd.StartDate.Add(sd.Duration)
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = td.CmpSmuggle(t, got, func(sd StartDuration) td.SmuggledGot {
			return td.SmuggledGot{
				Name: "ComputedEndDate",
				Got:  sd.StartDate.Add(sd.Duration),
			}
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
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

func ExampleCmpSmuggle_interface() {
	t := &testing.T{}

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := td.CmpSmuggle(t, gotTime, func(s fmt.Stringer) string {
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = td.CmpSmuggle(t, MyTime(gotTime), func(s fmt.Stringer) string {
		fmt.Println("Smuggle func called!")
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified MyTime OK:", ok)

	// Output
	// stringified time.Time OK: true
	// stringified MyTime OK: false
}

func ExampleCmpSmuggle_field_path() {
	t := &testing.T{}

	type Body struct {
		Name  string
		Value interface{}
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
	ok := td.CmpSmuggle(t, got, func(t *Transaction) (int, error) {
		if t.Request.Body == nil ||
			t.Request.Body.Value == nil {
			return 0, errors.New("Request.Body or Request.Body.Value is nil")
		}
		if v, ok := t.Request.Body.Value.(*ValueNum); ok && v != nil {
			return v.Num, nil
		}
		return 0, errors.New("Request.Body.Value isn't *ValueNum or nil")
	}, td.Between(100, 200))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = td.CmpSmuggle(t, got, "Request.Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = td.CmpSmuggle(t, got, "Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using an other fields-path:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
}

func ExampleCmpSStruct() {
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
	ok := td.CmpSStruct(t, got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = td.CmpSStruct(t, got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.CmpSStruct(t, &got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.CmpSStruct(t, &got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpString() {
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.String("foobar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleCmpString_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpString_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpStruct() {
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
	ok := td.CmpStruct(t, got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = td.CmpStruct(t, got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.CmpStruct(t, &got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.CmpStruct(t, &got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleCmpSubBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSubBagOf(t, got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = td.CmpSubBagOf(t, got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = td.CmpSubBagOf(t, got, []interface{}{td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Gt(4), td.Gt(4)},
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 9, 8}
	ok = td.CmpSubBagOf(t, got, []interface{}{td.Flatten(expected)},
		"checks at least all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// true
}

func ExampleCmpSubJSONOf_basic() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.CmpSubJSONOf(t, got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = td.CmpSubJSONOf(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // This field is ignored as SubJSONOf
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got without age field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got without age field: false
}

func ExampleCmpSubJSONOf_placeholders() {
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.CmpSubJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []interface{}{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"age": $age, "fullname": $name, "gender": $gender}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = td.CmpSubJSONOf(t, got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleCmpSubJSONOf_file() {
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

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
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
	ok := td.CmpSubJSONOf(t, got, filename, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.CmpSubJSONOf(t, got, file, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleCmpSubMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := td.CmpSubMapOf(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpSubMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := td.CmpSubMapOf(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = td.CmpSubMapOf(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpSubSetOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := td.CmpSubSetOf(t, got, []interface{}{1, 2, 3, 4, 5, 6, 7, 8},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = td.CmpSubSetOf(t, got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ok = td.CmpSubSetOf(t, got, []interface{}{td.Flatten(expected)},
		"checks at least all expected items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpSuperBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSuperBagOf(t, got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.CmpSuperBagOf(t, got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = td.CmpSuperBagOf(t, got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpSuperJSONOf_basic() {
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

	ok := td.CmpSuperJSONOf(t, got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = td.CmpSuperJSONOf(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // The gender!
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"fullname":"Bob","gender":"male","details":{}}`, nil)
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false
}

func ExampleCmpSuperJSONOf_placeholders() {
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

	ok := td.CmpSuperJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []interface{}{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $age, "fullname": $name, "gender": $gender}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleCmpSuperJSONOf_file() {
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

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.CmpSuperJSONOf(t, got, filename, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.CmpSuperJSONOf(t, got, file, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleCmpSuperMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpSuperMapOf(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleCmpSuperMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpSuperMapOf(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = td.CmpSuperMapOf(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCmpSuperSetOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSuperSetOf(t, got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = td.CmpSuperSetOf(t, got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = td.CmpSuperSetOf(t, got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
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

	// Compare dates ignoring nanoseconds and monotonic parts
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := td.CmpTruncTime(t, got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = td.CmpTruncTime(t, got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = td.CmpTruncTime(t, got, expected, 0,
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleCmpValues() {
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := td.CmpValues(t, got, []int{1, 2, 3})
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = td.CmpValues(t, got, []int{3, 1, 2})
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.CmpValues(t, got, td.Bag(3, 1, 2))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = td.CmpValues(t, got, td.ArrayEach(td.Between(1, 3)))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true
}

func ExampleCmpZero() {
	t := &testing.T{}

	ok := td.CmpZero(t, 0)
	fmt.Println(ok)

	ok = td.CmpZero(t, float64(0))
	fmt.Println(ok)

	ok = td.CmpZero(t, 12) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = td.CmpZero(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = td.CmpZero(t, map[string]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = td.CmpZero(t, ([]int)(nil))
	fmt.Println(ok)

	ok = td.CmpZero(t, []int{}) // fails, as not nil
	fmt.Println(ok)

	ok = td.CmpZero(t, [3]int{})
	fmt.Println(ok)

	ok = td.CmpZero(t, [3]int{0, 1}) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.CmpZero(t, bytes.Buffer{})
	fmt.Println(ok)

	ok = td.CmpZero(t, &bytes.Buffer{}) // fails, as pointer not nil
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
