// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	. "github.com/maxatome/go-testdeep"
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
	ok := CmpDeeply(t, masters, All(
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

	ok = CmpDeeply(t, &got,
		Array(&MyArray{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
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

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedArray() {
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleArrayEach_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleArrayEach_typedSlice() {
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpDeeply(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, ArrayEach(Between(25, 60)),
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
	ok := CmpDeeply(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = CmpDeeply(t, got, Bag(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = CmpDeeply(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = CmpDeeply(t, got, Bag(1, 2, 3, 5, Gt(7)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
}

func ExampleBetween() {
	t := &testing.T{}

	got := 156

	ok := CmpDeeply(t, got, Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = CmpDeeply(t, got, Between(154, 156, BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Between(154, 156, BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Between(154, 156, BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Between(154, 156, BoundsOutOut),
		"checks %v is in ]154 .. 156[", got)
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

	ok := CmpDeeply(t, got, Cap(12), "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = CmpDeeply(t, got, Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleCap_operator() {
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := CmpDeeply(t, got, Cap(Between(10, 12)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Cap(Gt(10)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleCode() {
	t := &testing.T{}

	got := "12"

	ok := CmpDeeply(t, got,
		Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleEmpty() {
	t := &testing.T{}

	ok := CmpDeeply(t, nil, Empty()) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = CmpDeeply(t, (*int)(nil), Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, "", Empty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = CmpDeeply(t, 0, Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, (map[string]int)(nil), Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, map[string]int{}, Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, ([]int)(nil), Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, []int{}, Empty())
	fmt.Println(ok)

	ok = CmpDeeply(t, []int{3}, Empty()) // fails, as not empty
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{}, Empty()) // fails, Empty() is not Zero()!
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

	ok := CmpDeeply(t, MySlice{}, Empty()) // Ptr() not needed
	fmt.Println(ok)

	ok = CmpDeeply(t, &MySlice{}, Empty())
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = CmpDeeply(t, &l3, Empty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = CmpDeeply(t, &MyStruct{}, Empty()) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleGt() {
	t := &testing.T{}

	got := 156

	ok := CmpDeeply(t, got, Gt(155), "checks %v is > 155", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Gt(156), "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleGte() {
	t := &testing.T{}

	got := 156

	ok := CmpDeeply(t, got, Gte(156), "checks %v is ≥ 156", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleIsa() {
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := CmpDeeply(t, got, Isa(TstStruct{}), "checks got is a TstStruct")
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Isa(&TstStruct{}),
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Isa(&TstStruct{}),
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

	ok := CmpDeeply(t, got, Isa((*fmt.Stringer)(nil)),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = CmpDeeply(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
	// does not match
	errGot = nil

	ok = CmpDeeply(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = CmpDeeply(t, &errGot, Isa((*error)(nil)),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleLen_slice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := CmpDeeply(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = CmpDeeply(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_map() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := CmpDeeply(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = CmpDeeply(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleLen_operatorSlice() {
	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := CmpDeeply(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Len(Lt(5)), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLen_operatorMap() {
	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := CmpDeeply(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Len(Gte(3)), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleLt() {
	t := &testing.T{}

	got := 156

	ok := CmpDeeply(t, got, Lt(157), "checks %v is < 157", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Lt(156), "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleLte() {
	t := &testing.T{}

	got := 156

	ok := CmpDeeply(t, got, Lte(156), "checks %v is ≤ 156", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleMap_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpDeeply(t, got,
		Map(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got,
		Map(map[string]int{},
			MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got,
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

	ok := CmpDeeply(t, got,
		Map(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
		Map(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
		Map(&MyMap{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
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

	ok := CmpDeeply(t, got, MapEach(Between(10, 90)),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleMapEach_typedMap() {
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpDeeply(t, got, MapEach(Between(10, 90)),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, MapEach(Between(10, 90)),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleN() {
	t := &testing.T{}

	got := 1.12345

	ok := CmpDeeply(t, got, N(1.1234, 0.00006),
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleNil() {
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := CmpDeeply(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = CmpDeeply(t, got, Nil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = CmpDeeply(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = CmpDeeply(t, got, Nil())
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

	ok := CmpDeeply(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = CmpDeeply(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = CmpDeeply(t, got, None(0, 10, 20, 30, Between(100, 199)),
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

	ok := CmpDeeply(t, got, NotAny(3, 6, 8, 41, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, NotAny(3, 6, 8, 42, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleNot() {
	t := &testing.T{}

	got := 42

	ok := CmpDeeply(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got, Not(Between(10, 30)),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = CmpDeeply(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleNotEmpty() {
	t := &testing.T{}

	ok := CmpDeeply(t, nil, NotEmpty()) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = CmpDeeply(t, "foobar", NotEmpty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = CmpDeeply(t, 0, NotEmpty())
	fmt.Println(ok)

	ok = CmpDeeply(t, map[string]int{"foobar": 42}, NotEmpty())
	fmt.Println(ok)

	ok = CmpDeeply(t, []int{1}, NotEmpty())
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{}, NotEmpty()) // succeeds, NotEmpty() is not NotZero()!
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

	ok := CmpDeeply(t, MySlice{12}, NotEmpty())
	fmt.Println(ok)

	ok = CmpDeeply(t, &MySlice{12}, NotEmpty()) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = CmpDeeply(t, &l3, NotEmpty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = CmpDeeply(t, &MyStruct{}, NotEmpty()) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNotNil() {
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := CmpDeeply(t, got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = CmpDeeply(t, got, NotNil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = CmpDeeply(t, got, Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = CmpDeeply(t, got, NotNil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNotZero() {
	t := &testing.T{}

	ok := CmpDeeply(t, 0, NotZero()) // fails
	fmt.Println(ok)

	ok = CmpDeeply(t, float64(0), NotZero()) // fails
	fmt.Println(ok)

	ok = CmpDeeply(t, 12, NotZero())
	fmt.Println(ok)

	ok = CmpDeeply(t, (map[string]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = CmpDeeply(t, map[string]int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, ([]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = CmpDeeply(t, []int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{}, NotZero()) // fails
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{0, 1}, NotZero()) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = CmpDeeply(t, bytes.Buffer{}, NotZero()) // fails
	fmt.Println(ok)

	ok = CmpDeeply(t, &bytes.Buffer{}, NotZero()) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, &bytes.Buffer{}, Ptr(NotZero())) // fails as deref by Ptr()
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

	ok := CmpDeeply(t, &got, PPtr(12))
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, PPtr(Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExamplePtr() {
	t := &testing.T{}

	got := 12

	ok := CmpDeeply(t, &got, Ptr(12))
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Ptr(Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true
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

func ExampleRe_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := CmpDeeply(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_error() {
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := CmpDeeply(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
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

func ExampleReAll_capture() {
	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpDeeply(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpDeeply(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleReAll_captureComplex() {
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := CmpDeeply(t, got,
		ReAll(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpDeeply(t, got,
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
	ok := CmpDeeply(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpDeeply(t, got, Re(expected), "checks value %s", got)
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
	ok := CmpDeeply(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledError() {
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := CmpDeeply(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleRe_compiledCapture() {
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := CmpDeeply(t, got, Re(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpDeeply(t, got, Re(expected, Set("biz", "foo", "bar")),
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
	ok := CmpDeeply(t, got, ReAll(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpDeeply(t, got, ReAll(expected, Set("biz", "foo", "bar")),
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
	ok := CmpDeeply(t, got,
		ReAll(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpDeeply(t, got,
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
	ok := CmpDeeply(t, got, Set(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = CmpDeeply(t, got, Set(1, 2, 2, 2, 2, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = CmpDeeply(t, got, Set(Between(1, 4), 3, Between(2, 10)),
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

	ok := CmpDeeply(t, got, Shallow(&data),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = CmpDeeply(t, got, Shallow(&MyStruct{Value: 12}),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleSlice_slice() {
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpDeeply(t, got, Slice([]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got,
		Slice([]int{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, got,
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

	ok := CmpDeeply(t, got, Slice(MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got, Slice(&MySlice{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
		Slice(&MySlice{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
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

	ok := CmpDeeply(t, got,
		Smuggle(func(n int64) int { return int(n) }, 123),
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = CmpDeeply(t, "123",
		Smuggle(
			func(numStr string) (int, bool) {
				n, err := strconv.Atoi(numStr)
				return n, err == nil
			},
			Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = CmpDeeply(t, "123",
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

	// Output:
	// true
	// true
	// true
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
		ok := CmpDeeply(t, got,
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
		ok = CmpDeeply(t, got,
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

func ExampleString() {
	t := &testing.T{}

	got := "foobar"

	ok := CmpDeeply(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleString_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpDeeply(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleString_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpDeeply(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix() {
	t := &testing.T{}

	got := "foobar"

	ok := CmpDeeply(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpDeeply(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasPrefix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpDeeply(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix() {
	t := &testing.T{}

	got := "foobar"

	ok := CmpDeeply(t, got, HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpDeeply(t, got, HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleHasSuffix_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpDeeply(t, got, HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleContains() {
	t := &testing.T{}

	got := "foobar"

	ok := CmpDeeply(t, got, Contains("oob"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleContains_stringer() {
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpDeeply(t, got, Contains("oob"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleContains_error() {
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpDeeply(t, got, Contains("oob"), "checks %s", got)
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
	ok := CmpDeeply(t, got,
		Struct(Person{Name: "Foobar"}, StructFields{
			"Age": Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = CmpDeeply(t, got,
		Struct(Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = CmpDeeply(t, &got,
		Struct(&Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = CmpDeeply(t, &got,
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

func ExampleSubBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := CmpDeeply(t, got, SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = CmpDeeply(t, got, SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = CmpDeeply(t, got, SubBagOf(
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

func ExampleSubMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := CmpDeeply(t, got,
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

	ok := CmpDeeply(t, got,
		SubMapOf(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
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
	ok := CmpDeeply(t, got, SubSetOf(1, 2, 3, 4, 5, 6, 7, 8),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = CmpDeeply(t, got, SubSetOf(Between(1, 4), 3, Between(2, 10), Gt(100)),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSuperBagOf() {
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := CmpDeeply(t, got, SuperBagOf(8, 5, 8),
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = CmpDeeply(t, got, SuperBagOf(Gt(5), Lte(2)),
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleSuperMapOf_map() {
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpDeeply(t, got,
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

	ok := CmpDeeply(t, got,
		SuperMapOf(MyMap{"bar": 42}, MapEntries{"foo": Lt(15)}),
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = CmpDeeply(t, &got,
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

	ok := CmpDeeply(t, got, SuperSetOf(1, 2, 3),
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = CmpDeeply(t, got, SuperSetOf(Gt(5), Lte(2)),
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
	ok := CmpDeeply(t, got, TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
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

func ExampleZero() {
	t := &testing.T{}

	ok := CmpDeeply(t, 0, Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, float64(0), Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, 12, Zero()) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = CmpDeeply(t, (map[string]int)(nil), Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, map[string]int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, ([]int)(nil), Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, []int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{}, Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, [3]int{0, 1}, Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = CmpDeeply(t, bytes.Buffer{}, Zero())
	fmt.Println(ok)

	ok = CmpDeeply(t, &bytes.Buffer{}, Zero()) // fails, as pointer not nil
	fmt.Println(ok)

	ok = CmpDeeply(t, &bytes.Buffer{}, Ptr(Zero())) // OK with the help of Ptr()
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
