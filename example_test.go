// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

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
	"testing"
	"time"
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
	ok := Cmp(t, masters, All(
		Len(Gt(0)), // len(masters) should be > 0
		ArrayEach(
			// For each Master
			Struct(Master{}, StructFields{
				// Master Name should be composed of 2 words, with 1st letter uppercased
				"Name": Re(`^[A-Z][a-z]+ [A-Z][a-z]+\z`),
				// Annual income should be greater than $10000
				"AnnualIncome": Gt(10000),
				"Pets": ArrayEach(
					// For each Pet
					Struct(&Pet{}, StructFields{
						// Pet Name should be composed of 1 word, with 1st letter uppercased
						"Name": Re(`^[A-Z][a-z]+\z`),
						"Birthday": All(
							// Pet should be born after 2010, January 1st, but before now!
							Between(dateToTime("2010-01-01T00:00:00Z"), time.Now()),
							// AND minutes, seconds and nanoseconds should be 0
							Code(func(t time.Time) bool {
								return t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
							}),
						),
						// Only dogs and cats allowed
						"Family": Any(Canidae, Felidae),
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

	ok := Cmp(t, []int{1, 2, 3},
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
	ok := Cmp(t,
		got,
		All(Re("o/b"), HasSuffix("bar"), "foo/bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = Cmp(t,
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
	ok := Cmp(t, got, Any(Re("zip"), HasSuffix("bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = Cmp(t, got, Any(Re("zip"), HasSuffix("foo")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleArray_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := Cmp(t, got, Array([3]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArray_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := Cmp(t, got, Array(MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, Array(&MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Array(&MyArray{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Array((*MyArray)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleArrayEach_array() {
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, ArrayEach(Between(25, 60)),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, ArrayEach(Between(25, 60)),
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
	ok := Cmp(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = Cmp(t, got, Bag(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = Cmp(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = Cmp(t, got, Bag(1, 2, 3, 5, Gt(7)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
}

func ExampleBetween_int() {
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = Cmp(t, got, Between(154, 156, BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsOutOut),
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

	ok := Cmp(t, got, Between("aaa", "abc"),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = Cmp(t, got, Between("aaa", "abc", BoundsInIn),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsInOut),
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsOutIn),
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsOutOut),
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleCap() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := Cmp(t, got, Cap(12), "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = Cmp(t, got, Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCap_operator() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := Cmp(t, got, Cap(Between(10, 12)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Cap(Gt(10)),
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
	ok := Cmp(t, got,
		JSON(`{"age":$1,"fullname":"Bob"}`,
			Catch(&age, Between(40, 45))))
	fmt.Println("check got age+fullname:", ok)
	fmt.Println("caught age:", age)

	// Output:
	// check got age+fullname: true
	// caught age: 42
}

func ExampleCode() {
	t := &testing.T{}

	got := "12"

	ok := Cmp(t, got,
		Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = Cmp(t, got,
		Code(func(num string) (bool, string) {
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
	ok = Cmp(t, got,
		Code(func(num string) error {
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

func ExampleContains_arraySlice() {
	t := &testing.T{}

	ok := Cmp(t, [...]int{11, 22, 33, 44}, Contains(22))
	fmt.Println("array contains 22:", ok)

	ok = Cmp(t, [...]int{11, 22, 33, 44}, Contains(Between(20, 25)))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = Cmp(t, []int{11, 22, 33, 44}, Contains(22))
	fmt.Println("slice contains 22:", ok)

	ok = Cmp(t, []int{11, 22, 33, 44}, Contains(Between(20, 25)))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
}

func ExampleContains_nil() {
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := Cmp(t, got, Contains(nil))
	fmt.Println("array contains untyped nil:", ok)

	ok = Cmp(t, got, Contains((*int)(nil)))
	fmt.Println("array contains *int nil:", ok)

	ok = Cmp(t, got, Contains(Nil()))
	fmt.Println("array contains Nil():", ok)

	ok = Cmp(t, got, Contains((*byte)(nil)))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false
}

func ExampleContains_map() {
	t := &testing.T{}

	ok := Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, Contains(22))
	fmt.Println("map contains value 22:", ok)

	ok = Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33},
		Contains(Between(20, 25)))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true
}

func ExampleContains_string() {
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleContains_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false
}

func ExampleContains_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false
}

func ExampleContainsKey() {
	t := &testing.T{}

	ok := Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, ContainsKey("foo"))
	fmt.Println(`map contains key "foo":`, ok)

	ok = Cmp(t,
		map[int]bool{12: true, 24: false, 42: true, 51: false},
		ContainsKey(Between(40, 50)))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
}

func ExampleContainsKey_nil() {
	t := &testing.T{}

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := Cmp(t, got, ContainsKey(nil))
	fmt.Println("map contains untyped nil key:", ok)

	ok = Cmp(t, got, ContainsKey((*int)(nil)))
	fmt.Println("map contains *int nil key:", ok)

	ok = Cmp(t, got, ContainsKey(Nil()))
	fmt.Println("map contains Nil() key:", ok)

	ok = Cmp(t, got, ContainsKey((*byte)(nil)))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false
}

func ExampleEmpty() {
	t := &testing.T{}

	ok := Cmp(t, nil, Empty()) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = Cmp(t, (*int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, "", Empty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = Cmp(t, 0, Empty())
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, Empty())
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, []int{}, Empty())
	fmt.Println(ok)

	ok = Cmp(t, []int{3}, Empty()) // fails, as not empty
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, Empty()) // fails, Empty() is not Zero()!
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

	ok := Cmp(t, MySlice{}, Empty()) // Ptr() not needed
	fmt.Println(ok)

	ok = Cmp(t, &MySlice{}, Empty())
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = Cmp(t, &l3, Empty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = Cmp(t, &MyStruct{}, Empty()) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleGt_int() {
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Gt(155), "checks %v is > 155", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gt(156), "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleGt_string() {
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Gt("abb"), `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gt("abc"), `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleGte_int() {
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Gte(156), "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte(155), "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte(157), "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleGte_string() {
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Gte("abc"), `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte("abb"), `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte("abd"), `checks "%v" is ≥ "abd"`, got)
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

	ok := Cmp(t, got, Isa(TstStruct{}), "checks got is a TstStruct")
	fmt.Println(ok)

	ok = Cmp(t, got, Isa(&TstStruct{}),
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = Cmp(t, &got, Isa(&TstStruct{}),
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

	ok := Cmp(t, got, Isa((*fmt.Stringer)(nil)),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = Cmp(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
	// does not match
	errGot = nil

	ok = Cmp(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = Cmp(t, &errGot, Isa((*error)(nil)),
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

	ok := Cmp(t, got, JSON(`{"age":42,"fullname":"Bob"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob","age":42}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = Cmp(t, got, JSON(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with gender field:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob"}`))
	fmt.Println("check got with fullname only:", ok)

	ok = Cmp(t, true, JSON(`true`))
	fmt.Println("check boolean got is true:", ok)

	ok = Cmp(t, 42, JSON(`42`))
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = Cmp(t, got, JSON(`null`))
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

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := Cmp(t, got, JSON(`{"age": $1, "fullname": $2}`, 42, "Bob Foobar"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": $1, "fullname": $2}`,
			Between(40, 45),
			HasSuffix("Foobar")))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": "$1", "fullname": "$2"}`,
			Between(40, 45),
			HasSuffix("Foobar")))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": $age, "fullname": $name}`,
			Tag("age", Between(40, 45)),
			Tag("name", HasSuffix("Foobar"))))
	fmt.Println("check got with named placeholders:", ok)

	ok = Cmp(t, got, JSON(`{"age": $^NotZero, "fullname": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
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
	ok := Cmp(t, got,
		JSON(filename,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = Cmp(t, got,
		JSON(file,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleKeys() {
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := Cmp(t, got, Keys([]string{"bar", "foo", "zip"}))
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = Cmp(t, got, Keys([]string{"zip", "bar", "foo"}))
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = Cmp(t, got, Keys(Bag("zip", "bar", "foo")))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = Cmp(t, got, Keys(ArrayEach(Len(3))))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true
}

func ExampleLax() {
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := Between(1230, 1240) // int type here

	ok := Cmp(t, gotInt64, Lax(expected))
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = Cmp(t, gotInt32, Lax(expected))
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = Cmp(t, gotMyInt, Lax(expected))
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true
}

func ExampleLen_slice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := Cmp(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_map() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := Cmp(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_operatorSlice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := Cmp(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(Lt(5)), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLen_operatorMap() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := Cmp(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(Gte(3)), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLt_int() {
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Lt(157), "checks %v is < 157", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lt(156), "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleLt_string() {
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Lt("abd"), `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lt("abc"), `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleLte_int() {
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Lte(156), "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte(157), "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte(155), "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleLte_string() {
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Lte("abc"), `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte("abd"), `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte("abb"), `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleMap_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got,
		Map(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Map(map[string]int{},
			MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Map((map[string]int)(nil),
			MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
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

	ok := Cmp(t, got,
		Map(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map(&MyMap{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map((*MyMap)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
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

	ok := Cmp(t, got, MapEach(Between(10, 90)),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleMapEach_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got, MapEach(Between(10, 90)),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, MapEach(Between(10, 90)),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleN() {
	t := &testing.T{}

	got := 1.12345

	ok := Cmp(t, got, N(1.1234, 0.00006),
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleNaN_float32() {
	t := &testing.T{}

	got := float32(math.NaN())

	ok := Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false
}

func ExampleNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false
}

func ExampleNil() {
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = Cmp(t, got, Nil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = Cmp(t, got, Nil())
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

	ok := Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
}

func ExampleNotAny() {
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := Cmp(t, got, NotAny(3, 6, 8, 41, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = Cmp(t, got, NotAny(3, 6, 8, 42, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleNot() {
	t := &testing.T{}

	got := 42

	ok := Cmp(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Not(Between(10, 30)),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = Cmp(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleNotEmpty() {
	t := &testing.T{}

	ok := Cmp(t, nil, NotEmpty()) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = Cmp(t, "foobar", NotEmpty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = Cmp(t, 0, NotEmpty())
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{"foobar": 42}, NotEmpty())
	fmt.Println(ok)

	ok = Cmp(t, []int{1}, NotEmpty())
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, NotEmpty()) // succeeds, NotEmpty() is not NotZero()!
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

	ok := Cmp(t, MySlice{12}, NotEmpty())
	fmt.Println(ok)

	ok = Cmp(t, &MySlice{12}, NotEmpty()) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = Cmp(t, &l3, NotEmpty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = Cmp(t, &MyStruct{}, NotEmpty()) // fails, use NotZero() instead
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

	ok := Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
}

func ExampleNotNaN_float64() {
	t := &testing.T{}

	got := math.NaN()

	ok := Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
}

func ExampleNotNil() {
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = Cmp(t, got, NotNil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = Cmp(t, got, NotNil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNotZero() {
	t := &testing.T{}

	ok := Cmp(t, 0, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, float64(0), NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, 12, NotZero())
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = Cmp(t, []int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, [3]int{0, 1}, NotZero()) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = Cmp(t, bytes.Buffer{}, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, NotZero()) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(NotZero())) // fails as deref by Ptr()
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

	ok := Cmp(t, &got, PPtr(12))
	fmt.Println(ok)

	ok = Cmp(t, &got, PPtr(Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExamplePtr() {
	t := &testing.T{}

	got := 12

	ok := Cmp(t, &got, Ptr(12))
	fmt.Println(ok)

	ok = Cmp(t, &got, Ptr(Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleRe() {
	t := &testing.T{}

	got := "foo bar"
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleRe_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_error() {
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := Cmp(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = Cmp(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := Cmp(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = Cmp(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_captureComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := Cmp(t, got,
		ReAll(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = Cmp(t, got,
		ReAll(`(\d+)`, ArrayEach(Code(func(num string) bool {
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
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = Cmp(t, got, Re(expected), "checks value %s", got)
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
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledError() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := Cmp(t, got, Re(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = Cmp(t, got, Re(expected, Set("biz", "foo", "bar")),
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
	ok := Cmp(t, got, ReAll(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = Cmp(t, got, ReAll(expected, Set("biz", "foo", "bar")),
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
	ok := Cmp(t, got,
		ReAll(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = Cmp(t, got,
		ReAll(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleSet() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := Cmp(t, got, Set(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = Cmp(t, got, Set(1, 2, 2, 2, 2, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = Cmp(t, got, Set(Between(1, 4), 3, Between(2, 10)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
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

	ok := Cmp(t, got, Shallow(&data),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = Cmp(t, got, Shallow(&MyStruct{Value: 12}),
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

	ok := Cmp(t, a, Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = Cmp(t, b, Shallow(back))
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

	ok := Cmp(t, a, Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = Cmp(t, b, Shallow(a))
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := Cmp(t, got, Slice([]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Slice([]int{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Slice(([]int)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
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

	ok := Cmp(t, got, Slice(MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, Slice(&MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Slice(&MySlice{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Slice((*MySlice)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleSmuggle_convert() {
	t := &testing.T{}

	got := int64(123)

	ok := Cmp(t, got,
		Smuggle(func(n int64) int { return int(n) }, 123),
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = Cmp(t, "123",
		Smuggle(
			func(numStr string) (int, bool) {
				n, err := strconv.Atoi(numStr)
				return n, err == nil
			},
			Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = Cmp(t, "123",
		Smuggle(
			func(numStr string) (int, bool, string) {
				n, err := strconv.Atoi(numStr)
				if err != nil {
					return 0, false, "string must contain a number"
				}
				return n, true, ""
			},
			Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = Cmp(t, "123",
		Smuggle(
			func(numStr string) (int, error) {
				return strconv.Atoi(numStr)
			},
			Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = Cmp(t, "123",
		Smuggle(strconv.Atoi, Between(120, 130)),
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

	ok := Cmp(t, got,
		Smuggle(func(n int64) uint32 { return uint32(n) }, uint32(123)))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true
}

func ExampleSmuggle_auto_unmarshal() {
	t := &testing.T{}

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := Cmp(t, got,
		Smuggle(
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

func ExampleSmuggle_complex() {
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
		ok := Cmp(t, got,
			Smuggle(
				func(sd StartDuration) time.Time {
					return sd.StartDate.Add(sd.Duration)
				},
				Between(
					time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC))))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = Cmp(t, got,
			Smuggle(
				func(sd StartDuration) SmuggledGot {
					return SmuggledGot{
						Name: "ComputedEndDate",
						Got:  sd.StartDate.Add(sd.Duration),
					}
				},
				Between(
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
	ok := Cmp(t, gotTime,
		Smuggle(func(s fmt.Stringer) string {
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = Cmp(t, MyTime(gotTime),
		Smuggle(func(s fmt.Stringer) string {
			fmt.Println("Smuggle func called!")
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified MyTime OK:", ok)

	// Output
	// stringified time.Time OK: true
	// stringified MyTime OK: false
}

func ExampleSmuggle_field_path() {
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
	ok := Cmp(t, got,
		Smuggle(
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
			Between(100, 200)))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = Cmp(t, got, Smuggle("Request.Body.Value.Num", Between(100, 200)))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = Cmp(t, got, Smuggle("Body.Value.Num", Between(100, 200)))
	fmt.Println("check Num using an other fields-path:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
}

func ExampleString() {
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleString_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleString_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix() {
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix() {
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, HasSuffix("bar"), "checks %s", got)
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
	ok := Cmp(t, got,
		Struct(Person{Name: "Foobar"}, StructFields{
			"Age": Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = Cmp(t, got,
		Struct(Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = Cmp(t, &got,
		Struct(&Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = Cmp(t, &got,
		Struct((*Person)(nil), StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
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
	ok := Cmp(t, got,
		SStruct(Person{Name: "Foobar"}, StructFields{
			"Age": Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = Cmp(t, got,
		SStruct(Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = Cmp(t, &got,
		SStruct(&Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = Cmp(t, &got,
		SStruct((*Person)(nil), StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleSubBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := Cmp(t, got, SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = Cmp(t, got, SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = Cmp(t, got, SubBagOf(
		Between(0, 3),
		Between(0, 3),
		Between(0, 3),
		Between(0, 3),
		Gt(4),
		Gt(4)),
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// Output:
	// true
	// false
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

	ok := Cmp(t, got, SubJSONOf(`{"age":42,"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = Cmp(t, got, SubJSONOf(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = Cmp(t, got, SubJSONOf(`
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

	ok = Cmp(t, got, SubJSONOf(`{"fullname":"Bob","gender":"male"}`))
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

	ok := Cmp(t, got,
		SubJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			42, "Bob Foobar", "male"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = Cmp(t, got,
		SubJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			Between(40, 45),
			HasSuffix("Foobar"),
			NotEmpty()))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = Cmp(t, got,
		SubJSONOf(`{"age": "$1", "fullname": "$2", "gender": "$3"}`,
			Between(40, 45),
			HasSuffix("Foobar"),
			NotEmpty()))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = Cmp(t, got,
		SubJSONOf(`{"age": $age, "fullname": $name, "gender": $gender}`,
			Tag("age", Between(40, 45)),
			Tag("name", HasSuffix("Foobar")),
			Tag("gender", NotEmpty())))
	fmt.Println("check got with named placeholders:", ok)

	ok = Cmp(t, got,
		SubJSONOf(`{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`))
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
	ok := Cmp(t, got,
		SubJSONOf(filename,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = Cmp(t, got,
		SubJSONOf(file,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleSubMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := Cmp(t, got,
		SubMapOf(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleSubMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := Cmp(t, got,
		SubMapOf(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		SubMapOf(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
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
	ok := Cmp(t, got, SubSetOf(1, 2, 3, 4, 5, 6, 7, 8),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = Cmp(t, got, SubSetOf(Between(1, 4), 3, Between(2, 10), Gt(100)),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSuperBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := Cmp(t, got, SuperBagOf(8, 5, 8),
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = Cmp(t, got, SuperBagOf(Gt(5), Lte(2)),
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
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

	ok := Cmp(t, got, SuperJSONOf(`{"age":42,"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = Cmp(t, got, SuperJSONOf(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = Cmp(t, got, SuperJSONOf(`
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

	ok = Cmp(t, got,
		SuperJSONOf(`{"fullname":"Bob","gender":"male","details":{}}`))
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

	ok := Cmp(t, got,
		SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			42, "Bob Foobar", "male"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = Cmp(t, got,
		SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			Between(40, 45),
			HasSuffix("Foobar"),
			NotEmpty()))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = Cmp(t, got,
		SuperJSONOf(`{"age": "$1", "fullname": "$2", "gender": "$3"}`,
			Between(40, 45),
			HasSuffix("Foobar"),
			NotEmpty()))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = Cmp(t, got,
		SuperJSONOf(`{"age": $age, "fullname": $name, "gender": $gender}`,
			Tag("age", Between(40, 45)),
			Tag("name", HasSuffix("Foobar")),
			Tag("gender", NotEmpty())))
	fmt.Println("check got with named placeholders:", ok)

	ok = Cmp(t, got,
		SuperJSONOf(`{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`))
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
	ok := Cmp(t, got,
		SuperJSONOf(filename,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = Cmp(t, got,
		SuperJSONOf(file,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleSuperMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got,
		SuperMapOf(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15)}),
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleSuperMapOf_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got,
		SuperMapOf(MyMap{"bar": 42}, MapEntries{"foo": Lt(15)}),
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		SuperMapOf(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15)}),
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSuperSetOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := Cmp(t, got, SuperSetOf(1, 2, 3),
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = Cmp(t, got, SuperSetOf(Gt(5), Lte(2)),
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// Output:
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
	ok := Cmp(t, got, TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = Cmp(t, got, TruncTime(expected, 24*time.Hour),
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = Cmp(t, got, TruncTime(expected),
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
	ok := Cmp(t, got, Values([]int{1, 2, 3}))
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = Cmp(t, got, Values([]int{3, 1, 2}))
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = Cmp(t, got, Values(Bag(3, 1, 2)))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = Cmp(t, got, Values(ArrayEach(Between(1, 3))))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true
}

func ExampleZero() {
	t := &testing.T{}

	ok := Cmp(t, 0, Zero())
	fmt.Println(ok)

	ok = Cmp(t, float64(0), Zero())
	fmt.Println(ok)

	ok = Cmp(t, 12, Zero()) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), Zero())
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), Zero())
	fmt.Println(ok)

	ok = Cmp(t, []int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, Zero())
	fmt.Println(ok)

	ok = Cmp(t, [3]int{0, 1}, Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = Cmp(t, bytes.Buffer{}, Zero())
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Zero()) // fails, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(Zero())) // OK with the help of Ptr()
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
