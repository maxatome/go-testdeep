// Copyright (c) 2018-2022, Maxime Soulé
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
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
)

func ExampleT_All() {
	t := td.NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := t.All(got, []any{td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = t.All(got, []any{td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []any copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = t.All(got, []any{td.HasPrefix("foo"), regOps, td.HasSuffix("bar")},
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Any() {
	t := td.NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := t.Any(got, []any{td.Re("zip"), td.HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = t.Any(got, []any{td.Re("zip"), td.HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []any copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("a/c"), td.Re(`^xx`), td.Re(`ar$`)})
	ok = t.Any(got, []any{td.HasPrefix("xxx"), regOps, td.HasSuffix("zip")},
		"check at least one operator matches value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Array_array() {
	t := td.NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.Array(got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Simple array:", ok)

	ok = t.Array(&got, &[3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer:", ok)

	ok = t.Array(&got, (*[3]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer, nil model:", ok)

	// Output:
	// Simple array: true
	// Array pointer: true
	// Array pointer, nil model: true
}

func ExampleT_Array_typedArray() {
	t := td.NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.Array(got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println("Typed array:", ok)

	ok = t.Array(&got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array:", ok)

	ok = t.Array(&got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, empty model:", ok)

	ok = t.Array(&got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, nil model:", ok)

	// Output:
	// Typed array: true
	// Pointer on a typed array: true
	// Pointer on a typed array, empty model: true
	// Pointer on a typed array, nil model: true
}

func ExampleT_ArrayEach_array() {
	t := td.NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.ArrayEach(got, td.Between(25, 60),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_ArrayEach_typedArray() {
	t := td.NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.ArrayEach(got, td.Between(25, 60),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, td.Between(25, 60),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_ArrayEach_slice() {
	t := td.NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.ArrayEach(got, td.Between(25, 60),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_ArrayEach_typedSlice() {
	t := td.NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.ArrayEach(got, td.Between(25, 60),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, td.Between(25, 60),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Bag() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := t.Bag(got, []any{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = t.Bag(got, []any{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = t.Bag(got, []any{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = t.Bag(got, []any{1, 2, 3, 5, td.Gt(7)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5}
	ok = t.Bag(got, []any{td.Flatten(expected), td.Gt(7)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
	// true
}

func ExampleT_Between_int() {
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleT_Between_string() {
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleT_Between_time() {
	t := td.NewT(&testing.T{})

	before := time.Now()
	occurredAt := time.Now()
	after := time.Now()

	ok := t.Between(occurredAt, before, after, td.BoundsInIn)
	fmt.Println("It occurred between before and after:", ok)

	type MyTime time.Time
	ok = t.Between(MyTime(occurredAt), MyTime(before), MyTime(after), td.BoundsInIn)
	fmt.Println("Same for convertible MyTime type:", ok)

	ok = t.Between(MyTime(occurredAt), before, after, td.BoundsInIn)
	fmt.Println("MyTime vs time.Time:", ok)

	ok = t.Between(occurredAt, before, 10*time.Second, td.BoundsInIn)
	fmt.Println("Using a time.Duration as TO:", ok)

	ok = t.Between(MyTime(occurredAt), MyTime(before), 10*time.Second, td.BoundsInIn)
	fmt.Println("Using MyTime as FROM and time.Duration as TO:", ok)

	// Output:
	// It occurred between before and after: true
	// Same for convertible MyTime type: true
	// MyTime vs time.Time: false
	// Using a time.Duration as TO: true
	// Using MyTime as FROM and time.Duration as TO: true
}

func ExampleT_Cap() {
	t := td.NewT(&testing.T{})

	got := make([]int, 0, 12)

	ok := t.Cap(got, 12, "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = t.Cap(got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Cap(got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Cap_operator() {
	t := td.NewT(&testing.T{})

	got := make([]int, 0, 12)

	ok := t.Cap(got, td.Between(10, 12),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = t.Cap(got, td.Gt(10),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Code() {
	t := td.NewT(&testing.T{})

	got := "12"

	ok := t.Code(got, func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = t.Code(got, func(num string) (bool, string) {
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
	ok = t.Code(got, func(num string) error {
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

func ExampleT_Code_custom() {
	t := td.NewT(&testing.T{})

	got := 123

	ok := t.Code(got, func(t *td.T, num int) {
		t.Cmp(num, 123)
	})
	fmt.Println("with one *td.T:", ok)

	ok = t.Code(got, func(assert, require *td.T, num int) {
		assert.Cmp(num, 123)
		require.Cmp(num, 123)
	})
	fmt.Println("with assert & require *td.T:", ok)

	// Output:
	// with one *td.T: true
	// with assert & require *td.T: true
}

func ExampleT_Contains_arraySlice() {
	t := td.NewT(&testing.T{})

	ok := t.Contains([...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = t.Contains([...]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, []int{22, 33})
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true
}

func ExampleT_Contains_nil() {
	t := td.NewT(&testing.T{})

	num := 123
	got := [...]*int{&num, nil}

	ok := t.Contains(got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = t.Contains(got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = t.Contains(got, td.Nil())
	fmt.Println("array contains Nil():", ok)

	ok = t.Contains(got, (*byte)(nil))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false
}

func ExampleT_Contains_map() {
	t := td.NewT(&testing.T{})

	ok := t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true
}

func ExampleT_Contains_string() {
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, []byte("oob"), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleT_Contains_stringer() {
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleT_Contains_error() {
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true
}

func ExampleT_ContainsKey() {
	t := td.NewT(&testing.T{})

	ok := t.ContainsKey(map[string]int{"foo": 11, "bar": 22, "zip": 33}, "foo")
	fmt.Println(`map contains key "foo":`, ok)

	ok = t.ContainsKey(map[int]bool{12: true, 24: false, 42: true, 51: false}, td.Between(40, 50))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = t.ContainsKey(map[string]int{"FOO": 11, "bar": 22, "zip": 33}, td.Smuggle(strings.ToLower, "foo"))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true
}

func ExampleT_ContainsKey_nil() {
	t := td.NewT(&testing.T{})

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := t.ContainsKey(got, nil)
	fmt.Println("map contains untyped nil key:", ok)

	ok = t.ContainsKey(got, (*int)(nil))
	fmt.Println("map contains *int nil key:", ok)

	ok = t.ContainsKey(got, td.Nil())
	fmt.Println("map contains Nil() key:", ok)

	ok = t.ContainsKey(got, (*byte)(nil))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false
}

func ExampleT_Empty() {
	t := td.NewT(&testing.T{})

	ok := t.Empty(nil) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = t.Empty((*int)(nil))
	fmt.Println(ok)

	ok = t.Empty("")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = t.Empty(0)
	fmt.Println(ok)

	ok = t.Empty((map[string]int)(nil))
	fmt.Println(ok)

	ok = t.Empty(map[string]int{})
	fmt.Println(ok)

	ok = t.Empty(([]int)(nil))
	fmt.Println(ok)

	ok = t.Empty([]int{})
	fmt.Println(ok)

	ok = t.Empty([]int{3}) // fails, as not empty
	fmt.Println(ok)

	ok = t.Empty([3]int{}) // fails, Empty() is not Zero()!
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

func ExampleT_Empty_pointers() {
	t := td.NewT(&testing.T{})

	type MySlice []int

	ok := t.Empty(MySlice{}) // Ptr() not needed
	fmt.Println(ok)

	ok = t.Empty(&MySlice{})
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = t.Empty(&l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = t.Empty(&MyStruct{}) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleT_CmpErrorIs() {
	t := td.NewT(&testing.T{})

	err1 := fmt.Errorf("failure1")
	err2 := fmt.Errorf("failure2: %w", err1)
	err3 := fmt.Errorf("failure3: %w", err2)
	err := fmt.Errorf("failure4: %w", err3)

	ok := t.CmpErrorIs(err, err)
	fmt.Println("error is itself:", ok)

	ok = t.CmpErrorIs(err, err1)
	fmt.Println("error is also err1:", ok)

	ok = t.CmpErrorIs(err1, err)
	fmt.Println("err1 is err:", ok)

	// Output:
	// error is itself: true
	// error is also err1: true
	// err1 is err: false
}

func ExampleT_First_classic() {
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.First(got, td.Gt(0), 1)
	fmt.Println("first positive number is 1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.First(got, isEven, -2)
	fmt.Println("first even number is -2:", ok)

	ok = t.First(got, isEven, td.Lt(0))
	fmt.Println("first even number is < 0:", ok)

	ok = t.First(got, isEven, td.Code(isEven))
	fmt.Println("first even number is well even:", ok)

	// Output:
	// first positive number is 1: true
	// first even number is -2: true
	// first even number is < 0: true
	// first even number is well even: true
}

func ExampleT_First_empty() {
	t := td.NewT(&testing.T{})

	ok := t.First(([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("first in nil slice:", ok)

	ok = t.First([]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty slice:", ok)

	ok = t.First(&[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty pointed slice:", ok)

	ok = t.First([0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty array:", ok)

	// Output:
	// first in nil slice: false
	// first in empty slice: false
	// first in empty pointed slice: false
	// first in empty array: false
}

func ExampleT_First_struct() {
	t := td.NewT(&testing.T{})

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

	ok := t.First(got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Bob Foobar"))
	fmt.Println("first person.Age > 30 → Bob:", ok)

	ok = t.First(got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Bob Foobar"}`))
	fmt.Println("first person.Age > 30 → Bob, using JSON:", ok)

	ok = t.First(got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Bob")))
	fmt.Println("first person.Age > 30 → Bob, using JSONPointer:", ok)

	// Output:
	// first person.Age > 30 → Bob: true
	// first person.Age > 30 → Bob, using JSON: true
	// first person.Age > 30 → Bob, using JSONPointer: true
}

func ExampleT_Grep_classic() {
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.Grep(got, td.Gt(0), []int{1, 2, 3})
	fmt.Println("check positive numbers:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.Grep(got, isEven, []int{-2, 0, 2})
	fmt.Println("even numbers are -2, 0 and 2:", ok)

	ok = t.Grep(got, isEven, td.Set(0, 2, -2))
	fmt.Println("even numbers are also 0, 2 and -2:", ok)

	ok = t.Grep(got, isEven, td.ArrayEach(td.Code(isEven)))
	fmt.Println("even numbers are each even:", ok)

	// Output:
	// check positive numbers: true
	// even numbers are -2, 0 and 2: true
	// even numbers are also 0, 2 and -2: true
	// even numbers are each even: true
}

func ExampleT_Grep_nil() {
	t := td.NewT(&testing.T{})

	var got []int
	ok := t.Grep(got, td.Gt(0), ([]int)(nil))
	fmt.Println("typed []int nil:", ok)

	ok = t.Grep(got, td.Gt(0), ([]string)(nil))
	fmt.Println("typed []string nil:", ok)

	ok = t.Grep(got, td.Gt(0), td.Nil())
	fmt.Println("td.Nil:", ok)

	ok = t.Grep(got, td.Gt(0), []int{})
	fmt.Println("empty non-nil slice:", ok)

	// Output:
	// typed []int nil: true
	// typed []string nil: false
	// td.Nil: true
	// empty non-nil slice: false
}

func ExampleT_Grep_struct() {
	t := td.NewT(&testing.T{})

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

	ok := t.Grep(got, td.Smuggle("Age", td.Gt(30)), td.All(
		td.Len(1),
		td.ArrayEach(td.Smuggle("Fullname", "Bob Foobar")),
	))
	fmt.Println("person.Age > 30 → only Bob:", ok)

	ok = t.Grep(got, td.JSONPointer("/age", td.Gt(30)), td.JSON(`[ SuperMapOf({"fullname":"Bob Foobar"}) ]`))
	fmt.Println("person.Age > 30 → only Bob, using JSON:", ok)

	// Output:
	// person.Age > 30 → only Bob: true
	// person.Age > 30 → only Bob, using JSON: true
}

func ExampleT_Gt_int() {
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Gt(got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = t.Gt(got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Gt_string() {
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Gt(got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = t.Gt(got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Gte_int() {
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Gte(got, 156, "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = t.Gte(got, 155, "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = t.Gte(got, 157, "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleT_Gte_string() {
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Gte(got, "abc", `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = t.Gte(got, "abb", `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = t.Gte(got, "abd", `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleT_HasPrefix() {
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = t.Cmp([]byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleT_HasPrefix_stringer() {
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasPrefix_error() {
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasSuffix() {
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = t.Cmp([]byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleT_HasSuffix_stringer() {
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasSuffix_error() {
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Isa() {
	t := td.NewT(&testing.T{})

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := t.Isa(got, TstStruct{}, "checks got is a TstStruct")
	fmt.Println(ok)

	ok = t.Isa(got, &TstStruct{},
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = t.Isa(&got, &TstStruct{},
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Isa_interface() {
	t := td.NewT(&testing.T{})

	got := bytes.NewBufferString("foobar")

	ok := t.Isa(got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = t.Isa(&errGot, (*error)(nil),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleT_JSON_basic() {
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := t.JSON(got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.JSON(got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = t.JSON(got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = t.JSON(true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = t.JSON(42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = t.JSON(got, `null`, nil)
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

func ExampleT_JSON_placeholders() {
	t := td.NewT(&testing.T{})

	type Person struct {
		Fullname string    `json:"fullname"`
		Age      int       `json:"age"`
		Children []*Person `json:"children,omitempty"`
	}

	got := &Person{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.JSON(got, `{"age": $1, "fullname": $2}`, []any{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.JSON(got, `{"age": $1, "fullname": $2}`, []any{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": "$1", "fullname": "$2"}`, []any{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": $age, "fullname": $name}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	got.Children = []*Person{
		{Fullname: "Alice", Age: 28},
		{Fullname: "Brian", Age: 22},
	}
	ok = t.JSON(got, `{"age": $age, "fullname": $name, "children": $children}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("children", td.Bag(
		&Person{Fullname: "Brian", Age: 22},
		&Person{Fullname: "Alice", Age: 28},
	))})
	fmt.Println("check got w/named placeholders, and children w/go structs:", ok)

	ok = t.JSON(got, `{"age": Between($1, $2), "fullname": HasSuffix($suffix), "children": Len(2)}`, []any{40, 45, td.Tag("suffix", "Foobar")})
	fmt.Println("check got w/num & named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got w/named placeholders, and children w/go structs: true
	// check got w/num & named placeholders: true
}

func ExampleT_JSON_embedding() {
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.JSON(got, `{"age": NotZero(), "fullname": NotEmpty()}`, nil)
	fmt.Println("check got with simple operators:", ok)

	ok = t.JSON(got, `{"age": $^NotZero, "fullname": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	ok = t.JSON(got, `
{
  "age":      Between(40, 42, "]]"), // in ]40; 42]
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = t.JSON(got, `
{
  "age":      Between(40, 42, "]["), // in ]40; 42[ → 42 excluded
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar"),
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = t.JSON(got, `
{
  "age":      Between($1, $2, $3), // in ]40; 42]
  "fullname": All(
    HasPrefix($4),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, []any{40, 42, td.BoundsOutIn, "Bob"})
	fmt.Println("check got with complex operators, w/placeholder args:", ok)

	// Output:
	// check got with simple operators: true
	// check got with operator shortcuts: true
	// check got with complex operators: true
	// check got with complex operators: false
	// check got with complex operators, w/placeholder args: true
}

func ExampleT_JSON_rawStrings() {
	t := td.NewT(&testing.T{})

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

	ok := t.JSON(got, `
{
  "fullname": HasPrefix("Foo"),
  "age":      Between(41, 43),
  "details":  SuperMapOf({
    "address": NotEmpty, // () are optional when no parameters
    "car":     Any("Peugeot", "Tesla", "Jeep") // any of these
  })
}`, nil)
	fmt.Println("Original:", ok)

	ok = t.JSON(got, `
{
  "fullname": "$^HasPrefix(\"Foo\")",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({\n\"address\": NotEmpty,\n\"car\": Any(\"Peugeot\", \"Tesla\", \"Jeep\")\n})"
}`, nil)
	fmt.Println("JSON compliant:", ok)

	ok = t.JSON(got, `
{
  "fullname": "$^HasPrefix(\"Foo\")",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    \"address\": NotEmpty, // () are optional when no parameters
    \"car\":     Any(\"Peugeot\", \"Tesla\", \"Jeep\") // any of these
  })"
}`, nil)
	fmt.Println("JSON multilines strings:", ok)

	ok = t.JSON(got, `
{
  "fullname": "$^HasPrefix(r<Foo>)",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    r<address>: NotEmpty, // () are optional when no parameters
    r<car>:     Any(r<Peugeot>, r<Tesla>, r<Jeep>) // any of these
  })"
}`, nil)
	fmt.Println("Raw strings:", ok)

	// Output:
	// Original: true
	// JSON compliant: true
	// JSON multilines strings: true
	// Raw strings: true
}

func ExampleT_JSON_file() {
	t := td.NewT(&testing.T{})

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
	ok := t.JSON(got, filename, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.JSON(got, file, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleT_JSONPointer_rfc6901() {
	t := td.NewT(&testing.T{})

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
	ok := t.JSONPointer(got, "", expected)
	fmt.Println("Empty JSON pointer means all:", ok)

	ok = t.JSONPointer(got, `/foo`, []any{"bar", "baz"})
	fmt.Println("Extract `foo` key:", ok)

	ok = t.JSONPointer(got, `/foo/0`, "bar")
	fmt.Println("First item of `foo` key slice:", ok)

	ok = t.JSONPointer(got, `/`, 0)
	fmt.Println("Empty key:", ok)

	ok = t.JSONPointer(got, `/a~1b`, 1)
	fmt.Println("Slash has to be escaped using `~1`:", ok)

	ok = t.JSONPointer(got, `/c%d`, 2)
	fmt.Println("% in key:", ok)

	ok = t.JSONPointer(got, `/e^f`, 3)
	fmt.Println("^ in key:", ok)

	ok = t.JSONPointer(got, `/g|h`, 4)
	fmt.Println("| in key:", ok)

	ok = t.JSONPointer(got, `/i\j`, 5)
	fmt.Println("Backslash in key:", ok)

	ok = t.JSONPointer(got, `/k"l`, 6)
	fmt.Println("Double-quote in key:", ok)

	ok = t.JSONPointer(got, `/ `, 7)
	fmt.Println("Space key:", ok)

	ok = t.JSONPointer(got, `/m~0n`, 8)
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

func ExampleT_JSONPointer_struct() {
	t := td.NewT(&testing.T{})

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

	ok := t.JSONPointer(got, "/Next/Next/Name", "third")
	fmt.Println("3rd item name is `third`:", ok)

	ok = t.JSONPointer(got, "/Next/Next/Value", td.Gte(int64(3)))
	fmt.Println("3rd item value is greater or equal than 3:", ok)

	ok = t.JSONPointer(got, "/Next", td.JSONPointer("/Next",
		td.JSONPointer("/Value", td.Gte(int64(3)))))
	fmt.Println("3rd item value is still greater or equal than 3:", ok)

	ok = t.JSONPointer(got, "/Next/Next/Next/Name", td.Ignore())
	fmt.Println("4th item exists and has a name:", ok)

	// Struct comparison work with or without pointer: &Item{…} works too
	ok = t.JSONPointer(got, "/Next/Next", Item{
		Name:  "third",
		Value: 3,
	})
	fmt.Println("3rd item full comparison:", ok)

	// Output:
	// 3rd item name is `third`: true
	// 3rd item value is greater or equal than 3: true
	// 3rd item value is still greater or equal than 3: true
	// 4th item exists and has a name: false
	// 3rd item full comparison: true
}

func ExampleT_JSONPointer_has_hasnt() {
	t := td.NewT(&testing.T{})

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
	ok := t.JSONPointer(got, "/children", td.Len(td.Gt(0)))
	fmt.Println("Bob has at least one child:", ok)

	// But checking "children" exists is enough here
	ok = t.JSONPointer(got, "/children/0/children", td.Ignore())
	fmt.Println("Alice has children:", ok)

	ok = t.JSONPointer(got, "/children/1/children", td.Ignore())
	fmt.Println("Britt has children:", ok)

	// The reverse can be checked too
	ok = t.Cmp(got, td.Not(td.JSONPointer("/children/0/children", td.Ignore())))
	fmt.Println("Alice hasn't children:", ok)

	ok = t.Cmp(got, td.Not(td.JSONPointer("/children/1/children", td.Ignore())))
	fmt.Println("Britt hasn't children:", ok)

	// Output:
	// Bob has at least one child: true
	// Alice has children: false
	// Britt has children: true
	// Alice hasn't children: true
	// Britt hasn't children: false
}

func ExampleT_Keys() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := t.Keys(got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = t.Keys(got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = t.Keys(got, td.Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = t.Keys(got, td.ArrayEach(td.Len(3)))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true
}

func ExampleT_Last_classic() {
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.Last(got, td.Lt(0), -1)
	fmt.Println("last negative number is -1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.Last(got, isEven, 2)
	fmt.Println("last even number is 2:", ok)

	ok = t.Last(got, isEven, td.Gt(0))
	fmt.Println("last even number is > 0:", ok)

	ok = t.Last(got, isEven, td.Code(isEven))
	fmt.Println("last even number is well even:", ok)

	// Output:
	// last negative number is -1: true
	// last even number is 2: true
	// last even number is > 0: true
	// last even number is well even: true
}

func ExampleT_Last_empty() {
	t := td.NewT(&testing.T{})

	ok := t.Last(([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("last in nil slice:", ok)

	ok = t.Last([]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty slice:", ok)

	ok = t.Last(&[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty pointed slice:", ok)

	ok = t.Last([0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty array:", ok)

	// Output:
	// last in nil slice: false
	// last in empty slice: false
	// last in empty pointed slice: false
	// last in empty array: false
}

func ExampleT_Last_struct() {
	t := td.NewT(&testing.T{})

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

	ok := t.Last(got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Alice Bingo"))
	fmt.Println("last person.Age > 30 → Alice:", ok)

	ok = t.Last(got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Alice Bingo"}`))
	fmt.Println("last person.Age > 30 → Alice, using JSON:", ok)

	ok = t.Last(got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Alice")))
	fmt.Println("first person.Age > 30 → Alice, using JSONPointer:", ok)

	// Output:
	// last person.Age > 30 → Alice: true
	// last person.Age > 30 → Alice, using JSON: true
	// first person.Age > 30 → Alice, using JSONPointer: true
}

func ExampleT_CmpLax() {
	t := td.NewT(&testing.T{})

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

	ok := t.CmpLax(gotInt64, expected)
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = t.CmpLax(gotInt32, expected)
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = t.CmpLax(gotMyInt, expected)
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true
}

func ExampleT_Len_slice() {
	t := td.NewT(&testing.T{})

	got := []int{11, 22, 33}

	ok := t.Len(got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Len_map() {
	t := td.NewT(&testing.T{})

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := t.Len(got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_Len_operatorSlice() {
	t := td.NewT(&testing.T{})

	got := []int{11, 22, 33}

	ok := t.Len(got, td.Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, td.Lt(5), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Len_operatorMap() {
	t := td.NewT(&testing.T{})

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := t.Len(got, td.Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, td.Gte(3), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Lt_int() {
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Lt(got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = t.Lt(got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Lt_string() {
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Lt(got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = t.Lt(got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Lte_int() {
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Lte(got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = t.Lte(got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = t.Lte(got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleT_Lte_string() {
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Lte(got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleT_Map_map() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, map[string]int{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, (map[string]int)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_Map_typedMap() {
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, (*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleT_MapEach_map() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, td.Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_MapEach_typedMap() {
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, td.Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = t.MapEach(&got, td.Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_N() {
	t := td.NewT(&testing.T{})

	got := 1.12345

	ok := t.N(got, 1.1234, 0.00006,
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_NaN_float32() {
	t := td.NewT(&testing.T{})

	got := float32(math.NaN())

	ok := t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false
}

func ExampleT_NaN_float64() {
	t := td.NewT(&testing.T{})

	got := math.NaN()

	ok := t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false
}

func ExampleT_Nil() {
	t := td.NewT(&testing.T{})

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := t.Cmp(got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = t.Nil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = t.Cmp(got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = t.Nil(got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleT_None() {
	t := td.NewT(&testing.T{})

	got := 18

	ok := t.None(got, []any{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = t.None(got, []any{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = t.None(got, []any{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = t.None(got, []any{prime, even, td.Gt(14)},
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

func ExampleT_Not() {
	t := td.NewT(&testing.T{})

	got := 42

	ok := t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = t.Not(got, td.Between(10, 30),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
}

func ExampleT_NotAny() {
	t := td.NewT(&testing.T{})

	got := []int{4, 5, 9, 42}

	ok := t.NotAny(got, []any{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = t.NotAny(got, []any{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []any slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = t.NotAny(got, []any{td.Flatten(notExpected)},
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_NotEmpty() {
	t := td.NewT(&testing.T{})

	ok := t.NotEmpty(nil) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = t.NotEmpty("foobar")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = t.NotEmpty(0)
	fmt.Println(ok)

	ok = t.NotEmpty(map[string]int{"foobar": 42})
	fmt.Println(ok)

	ok = t.NotEmpty([]int{1})
	fmt.Println(ok)

	ok = t.NotEmpty([3]int{}) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true
}

func ExampleT_NotEmpty_pointers() {
	t := td.NewT(&testing.T{})

	type MySlice []int

	ok := t.NotEmpty(MySlice{12})
	fmt.Println(ok)

	ok = t.NotEmpty(&MySlice{12}) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = t.NotEmpty(&l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = t.NotEmpty(&MyStruct{}) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleT_NotNaN_float32() {
	t := td.NewT(&testing.T{})

	got := float32(math.NaN())

	ok := t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
}

func ExampleT_NotNaN_float64() {
	t := td.NewT(&testing.T{})

	got := math.NaN()

	ok := t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
}

func ExampleT_NotNil() {
	t := td.NewT(&testing.T{})

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := t.Cmp(got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = t.NotNil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = t.Cmp(got, td.Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = t.NotNil(got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleT_NotZero() {
	t := td.NewT(&testing.T{})

	ok := t.NotZero(0) // fails
	fmt.Println(ok)

	ok = t.NotZero(float64(0)) // fails
	fmt.Println(ok)

	ok = t.NotZero(12)
	fmt.Println(ok)

	ok = t.NotZero((map[string]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = t.NotZero(map[string]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = t.NotZero(([]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = t.NotZero([]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = t.NotZero([3]int{}) // fails
	fmt.Println(ok)

	ok = t.NotZero([3]int{0, 1}) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = t.NotZero(bytes.Buffer{}) // fails
	fmt.Println(ok)

	ok = t.NotZero(&bytes.Buffer{}) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = t.Cmp(&bytes.Buffer{}, td.Ptr(td.NotZero())) // fails as deref by Ptr()
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

func ExampleT_PPtr() {
	t := td.NewT(&testing.T{})

	num := 12
	got := &num

	ok := t.PPtr(&got, 12)
	fmt.Println(ok)

	ok = t.PPtr(&got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Ptr() {
	t := td.NewT(&testing.T{})

	got := 12

	ok := t.Ptr(&got, 12)
	fmt.Println(ok)

	ok = t.Ptr(&got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Re() {
	t := td.NewT(&testing.T{})

	got := "foo bar"
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Re_stringer() {
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_error() {
	t := td.NewT(&testing.T{})

	got := errors.New("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_capture() {
	t := td.NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.Re(got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Re_compiled() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Re_compiledStringer() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_compiledError() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_compiledCapture() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := t.Re(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_capture() {
	t := td.NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.ReAll(got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_captureComplex() {
	t := td.NewT(&testing.T{})

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_compiledCapture() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := t.ReAll(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_compiledCaptureComplex() {
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Recv_basic() {
	t := td.NewT(&testing.T{})

	got := make(chan int, 3)

	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = t.Recv(got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = t.Cmp(got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false
}

func ExampleT_Recv_channelPointer() {
	t := td.NewT(&testing.T{})

	got := make(chan int, 3)

	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = t.Recv(&got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = t.Cmp(&got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false
}

func ExampleT_Recv_withTimeout() {
	t := td.NewT(&testing.T{})

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

	t.Recv(got, td.RecvNothing, 0)

	// ①
	tick <- struct{}{}
	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)
	ok = t.Recv(got, 0, 150*time.Millisecond)
	fmt.Println("① receive 0 w/150ms timeout:", ok)
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)

	// ②
	tick <- struct{}{}
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)
	ok = t.Recv(got, 1, 150*time.Millisecond)
	fmt.Println("② receive 1 w/150ms timeout:", ok)
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)

	// ③
	tick <- struct{}{}
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("③ RecvNothing:", ok)
	ok = t.Recv(got, td.RecvClosed, 150*time.Millisecond)
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

func ExampleT_Recv_nilChannel() {
	t := td.NewT(&testing.T{})

	var ch chan int

	ok := t.Recv(ch, td.RecvNothing, 0)
	fmt.Println("nothing to receive from nil channel:", ok)

	ok = t.Recv(ch, 42, 0)
	fmt.Println("something to receive from nil channel:", ok)

	ok = t.Recv(ch, td.RecvClosed, 0)
	fmt.Println("is a nil channel closed:", ok)

	// Output:
	// nothing to receive from nil channel: true
	// something to receive from nil channel: false
	// is a nil channel closed: false
}

func ExampleT_Set() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := t.Set(got, []any{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = t.Set(got, []any{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = t.Set(got, []any{td.Between(1, 4), 3, td.Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = t.Set(got, []any{td.Flatten(expected)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleT_Shallow() {
	t := td.NewT(&testing.T{})

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := t.Shallow(got, &data,
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = t.Shallow(got, &MyStruct{Value: 12},
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Shallow_slice() {
	t := td.NewT(&testing.T{})

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := t.Shallow(a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = t.Shallow(b, back)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleT_Shallow_string() {
	t := td.NewT(&testing.T{})

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := t.Shallow(a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = t.Shallow(b, a)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false
}

func ExampleT_Slice_slice() {
	t := td.NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.Slice(got, []int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(got, []int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(got, ([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_Slice_typedSlice() {
	t := td.NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.Slice(got, MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, &MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, &MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, (*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleT_Smuggle_convert() {
	t := td.NewT(&testing.T{})

	got := int64(123)

	ok := t.Smuggle(got, func(n int64) int { return int(n) }, 123,
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, bool) {
		n, err := strconv.Atoi(numStr)
		return n, err == nil
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, bool, string) {
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, false, "string must contain a number"
		}
		return n, true, ""
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, error) { //nolint: gocritic
		return strconv.Atoi(numStr)
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = t.Smuggle("123", strconv.Atoi, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true
}

func ExampleT_Smuggle_lax() {
	t := td.NewT(&testing.T{})

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := t.Smuggle(got, func(n int64) uint32 { return uint32(n) }, uint32(123))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true
}

func ExampleT_Smuggle_auto_unmarshal() {
	t := td.NewT(&testing.T{})

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := t.Smuggle(got, func(b json.RawMessage) (r map[string]int, err error) {
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

func ExampleT_Smuggle_cast() {
	t := td.NewT(&testing.T{})

	// A string containing JSON
	got := `{ "foo": 123 }`

	// Automatically cast a string to a json.RawMessage so td.JSON can operate
	ok := t.Smuggle(got, json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents in string is OK:", ok)

	// Automatically read from io.Reader to a json.RawMessage
	ok = t.Smuggle(bytes.NewReader([]byte(got)), json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents just read is OK:", ok)

	// Output:
	// JSON contents in string is OK: true
	// JSON contents just read is OK: true
}

func ExampleT_Smuggle_complex() {
	t := td.NewT(&testing.T{})

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
		ok := t.Smuggle(got, func(sd StartDuration) time.Time {
			return sd.StartDate.Add(sd.Duration)
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = t.Smuggle(got, func(sd StartDuration) td.SmuggledGot {
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

func ExampleT_Smuggle_interface() {
	t := td.NewT(&testing.T{})

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := t.Smuggle(gotTime, func(s fmt.Stringer) string {
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = t.Smuggle(MyTime(gotTime), func(s fmt.Stringer) string {
		fmt.Println("Smuggle func called!")
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified MyTime OK:", ok)

	// Output:
	// stringified time.Time OK: true
	// stringified MyTime OK: false
}

func ExampleT_Smuggle_field_path() {
	t := td.NewT(&testing.T{})

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
	ok := t.Smuggle(got, func(t *Transaction) (int, error) {
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
	ok = t.Smuggle(got, "Request.Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = t.Smuggle(got, "Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using an other fields-path:", ok)

	// Note that maps and array/slices are supported
	got.Request.Body.Value = map[string]any{
		"foo": []any{
			3: map[int]string{666: "bar"},
		},
	}
	ok = t.Smuggle(got, "Body.Value[foo][3][666]", "bar")
	fmt.Println("check fields-path including maps/slices:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
	// check fields-path including maps/slices: true
}

func ExampleT_SStruct() {
	t := td.NewT(&testing.T{})

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
	ok := t.SStruct(got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	got.NumChildren = 3
	ok = t.SStruct(got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = t.SStruct(&got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = t.SStruct(&got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true
}

func ExampleT_SStruct_overwrite_model() {
	t := td.NewT(&testing.T{})

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

	ok := t.SStruct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = t.SStruct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true
}

func ExampleT_SStruct_patterns() {
	t := td.NewT(&testing.T{})

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

	ok := t.SStruct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`: nil,
		`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		`!  [A-Z]*`: td.Ignore(),        // private fields
	},
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = t.SStruct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`:   nil,
		`1 =  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`2 =~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		`3 !~ ^[A-Z]`: td.Ignore(),        // private fields
	},
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true
}

func ExampleT_SStruct_lazy_model() {
	t := td.NewT(&testing.T{})

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := t.SStruct(got, nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	})
	fmt.Println("Lazy model:", ok)

	ok = t.SStruct(got, nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	})
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false
}

func ExampleT_String() {
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = t.Cmp([]byte(got), td.String("foobar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true
}

func ExampleT_String_stringer() {
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_String_error() {
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Struct() {
	t := td.NewT(&testing.T{})

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
	ok := t.Struct(got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	ok = t.Struct(got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = t.Struct(&got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = t.Struct(&got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true
}

func ExampleT_Struct_overwrite_model() {
	t := td.NewT(&testing.T{})

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

	ok := t.Struct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = t.Struct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true
}

func ExampleT_Struct_patterns() {
	t := td.NewT(&testing.T{})

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

	ok := t.Struct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`: nil,
		`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = t.Struct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`:  nil,
		`1 =  *name`: td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`2 =~ At\z`:  td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true
}

func ExampleT_Struct_lazy_model() {
	t := td.NewT(&testing.T{})

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := t.Struct(got, nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	})
	fmt.Println("Lazy model:", ok)

	ok = t.Struct(got, nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	})
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false
}

func ExampleT_SubBagOf() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SubBagOf(got, []any{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = t.SubBagOf(got, []any{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = t.SubBagOf(got, []any{td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Gt(4), td.Gt(4)},
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 9, 8}
	ok = t.SubBagOf(got, []any{td.Flatten(expected)},
		"checks at least all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// true
}

func ExampleT_SubJSONOf_basic() {
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := t.SubJSONOf(got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.SubJSONOf(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.SubJSONOf(got, `
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

	ok = t.SubJSONOf(got, `{"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got without age field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got without age field: false
}

func ExampleT_SubJSONOf_placeholders() {
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.SubJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.SubJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.SubJSONOf(got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.SubJSONOf(got, `{"age": $age, "fullname": $name, "gender": $gender}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = t.SubJSONOf(got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleT_SubJSONOf_file() {
	t := td.NewT(&testing.T{})

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
	ok := t.SubJSONOf(got, filename, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.SubJSONOf(got, file, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleT_SubMapOf_map() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_SubMapOf_typedMap() {
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = t.SubMapOf(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SubSetOf() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := t.SubSetOf(got, []any{1, 2, 3, 4, 5, 6, 7, 8},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = t.SubSetOf(got, []any{td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ok = t.SubSetOf(got, []any{td.Flatten(expected)},
		"checks at least all expected items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_SuperBagOf() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperBagOf(got, []any{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = t.SuperBagOf(got, []any{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = t.SuperBagOf(got, []any{td.Flatten(expected)},
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_SuperJSONOf_basic() {
	t := td.NewT(&testing.T{})

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

	ok := t.SuperJSONOf(got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.SuperJSONOf(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.SuperJSONOf(got, `
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

	ok = t.SuperJSONOf(got, `{"fullname":"Bob","gender":"male","details":{}}`, nil)
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false
}

func ExampleT_SuperJSONOf_placeholders() {
	t := td.NewT(&testing.T{})

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

	ok := t.SuperJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.SuperJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": $age, "fullname": $name, "gender": $gender}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true
}

func ExampleT_SuperJSONOf_file() {
	t := td.NewT(&testing.T{})

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
	ok := t.SuperJSONOf(got, filename, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.SuperJSONOf(got, file, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true
}

func ExampleT_SuperMapOf_map() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks map %v contains at least all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_SuperMapOf_typedMap() {
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks typed map %v contains at least all expected keys/values", got)
	fmt.Println(ok)

	ok = t.SuperMapOf(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks pointed typed map %v contains at least all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SuperSetOf() {
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperSetOf(got, []any{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = t.SuperSetOf(got, []any{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = t.SuperSetOf(got, []any{td.Flatten(expected)},
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_SuperSliceOf_array() {
	t := td.NewT(&testing.T{})

	got := [4]int{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, [4]int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, [4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &[4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = t.SuperSliceOf(&got, (*[4]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true
}

func ExampleT_SuperSliceOf_typedArray() {
	t := td.NewT(&testing.T{})

	type MyArray [4]int

	got := MyArray{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, MyArray{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = t.SuperSliceOf(&got, (*MyArray)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true
}

func ExampleT_SuperSliceOf_slice() {
	t := td.NewT(&testing.T{})

	got := []int{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, []int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, []int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &[]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = t.SuperSliceOf(&got, (*[]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true
}

func ExampleT_SuperSliceOf_typedSlice() {
	t := td.NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, MySlice{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = t.SuperSliceOf(&got, (*MySlice)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true
}

func ExampleT_TruncTime() {
	t := td.NewT(&testing.T{})

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
	ok := t.TruncTime(got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = t.TruncTime(got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = t.TruncTime(got, expected, 0,
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_Values() {
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := t.Values(got, []int{1, 2, 3})
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = t.Values(got, []int{3, 1, 2})
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = t.Values(got, td.Bag(3, 1, 2))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = t.Values(got, td.ArrayEach(td.Between(1, 3)))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true
}

func ExampleT_Zero() {
	t := td.NewT(&testing.T{})

	ok := t.Zero(0)
	fmt.Println(ok)

	ok = t.Zero(float64(0))
	fmt.Println(ok)

	ok = t.Zero(12) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = t.Zero((map[string]int)(nil))
	fmt.Println(ok)

	ok = t.Zero(map[string]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = t.Zero(([]int)(nil))
	fmt.Println(ok)

	ok = t.Zero([]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = t.Zero([3]int{})
	fmt.Println(ok)

	ok = t.Zero([3]int{0, 1}) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = t.Zero(bytes.Buffer{})
	fmt.Println(ok)

	ok = t.Zero(&bytes.Buffer{}) // fails, as pointer not nil
	fmt.Println(ok)

	ok = t.Cmp(&bytes.Buffer{}, td.Ptr(td.Zero())) // OK with the help of Ptr()
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
