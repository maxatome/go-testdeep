// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

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

func ExampleT_All() {
	t := NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := t.All(got, []interface{}{Re("o/b"), HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = t.All(got, []interface{}{Re("o/b"), HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Any() {
	t := NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := t.Any(got, []interface{}{Re("zip"), HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = t.Any(got, []interface{}{Re("zip"), HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Array_array() {
	t := NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.Array(got, [3]int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Array_typedArray() {
	t := NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.Array(got, MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, &MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_ArrayEach_array() {
	t := NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_ArrayEach_typedArray() {
	t := NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, Between(25, 60),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_ArrayEach_slice() {
	t := NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_ArrayEach_typedSlice() {
	t := NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, Between(25, 60),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Bag() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := t.Bag(got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = t.Bag(got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = t.Bag(got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = t.Bag(got, []interface{}{1, 2, 3, 5, Gt(7)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true
}

func ExampleT_Between() {
	t := NewT(&testing.T{})

	got := 156

	ok := t.Between(got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleT_Cap() {
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	got := make([]int, 0, 12)

	ok := t.Cap(got, Between(10, 12),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = t.Cap(got, Gt(10),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Code() {
	t := NewT(&testing.T{})

	got := "12"

	ok := t.Code(got, func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Contains() {
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Contains_stringer() {
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Contains_error() {
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Gt() {
	t := NewT(&testing.T{})

	got := 156

	ok := t.Gt(got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = t.Gt(got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Gte() {
	t := NewT(&testing.T{})

	got := 156

	ok := t.Gte(got, 156, "checks %v is ≥ 156", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasPrefix() {
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasPrefix_stringer() {
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasPrefix_error() {
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasSuffix() {
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasSuffix_stringer() {
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_HasSuffix_error() {
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Isa() {
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	got := bytes.NewBufferString("foobar")

	ok := t.Isa(got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
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

func ExampleT_Len_slice() {
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	got := []int{11, 22, 33}

	ok := t.Len(got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, Lt(5), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Len_operatorMap() {
	t := NewT(&testing.T{})

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := t.Len(got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, Gte(3), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Lt() {
	t := NewT(&testing.T{})

	got := 156

	ok := t.Lt(got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = t.Lt(got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Lte() {
	t := NewT(&testing.T{})

	got := 156

	ok := t.Lte(got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Map_map() {
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Map_typedMap() {
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_MapEach_map() {
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_MapEach_typedMap() {
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = t.MapEach(&got, Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_N() {
	t := NewT(&testing.T{})

	got := 1.12345

	ok := t.N(got, 1.1234, 0.00006,
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Nil() {
	t := NewT(&testing.T{})

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := t.CmpDeeply(got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = t.Nil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = t.CmpDeeply(got, nil)
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
	t := NewT(&testing.T{})

	got := 18

	ok := t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
}

func ExampleT_NoneOf() {
	t := NewT(&testing.T{})

	got := []int{4, 5, 9, 42}

	ok := t.NoneOf(got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NoneOf()", got)
	fmt.Println(ok)

	ok = t.NoneOf(got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NoneOf()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Not() {
	t := NewT(&testing.T{})

	got := 42

	ok := t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = t.Not(got, Between(10, 30),
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

func ExampleT_NotNil() {
	t := NewT(&testing.T{})

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := t.CmpDeeply(got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = t.NotNil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = t.CmpDeeply(got, Not(nil))
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

func ExampleT_PPtr() {
	t := NewT(&testing.T{})

	num := 12
	got := &num

	ok := t.PPtr(&got, 12)
	fmt.Println(ok)

	ok = t.PPtr(&got, Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Ptr() {
	t := NewT(&testing.T{})

	got := 12

	ok := t.Ptr(&got, 12)
	fmt.Println(ok)

	ok = t.Ptr(&got, Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_Re() {
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_error() {
	t := NewT(&testing.T{})

	got := errors.New("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_capture() {
	t := NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.Re(got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Re_compiled() {
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_compiledError() {
	t := NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Re_compiledCapture() {
	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := t.Re(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_capture() {
	t := NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.ReAll(got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_captureComplex() {
	t := NewT(&testing.T{})

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, `(\d+)`, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, `(\d+)`, ArrayEach(Code(func(num string) bool {
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
	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := t.ReAll(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_ReAll_compiledCaptureComplex() {
	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleT_Set() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := t.Set(got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = t.Set(got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = t.Set(got, []interface{}{Between(1, 4), 3, Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_Shallow() {
	t := NewT(&testing.T{})

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

func ExampleT_Slice_slice() {
	t := NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.Slice(got, []int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Slice_typedSlice() {
	t := NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.Slice(got, MySlice{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, &MySlice{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_String() {
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_String_stringer() {
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_String_error() {
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_Struct() {
	t := NewT(&testing.T{})

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
	ok := t.Struct(got, Person{Name: "Foobar"}, StructFields{
		"Age": Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = t.Struct(got, Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = t.Struct(&got, &Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
}

func ExampleT_SubBagOf() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SubBagOf(got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = t.SubBagOf(got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = t.SubBagOf(got, []interface{}{Between(0, 3), Between(0, 3), Between(0, 3), Between(0, 3), Gt(4), Gt(4)},
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
}

func ExampleT_SubMapOf_map() {
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_SubMapOf_typedMap() {
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = t.SubMapOf(&got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SubSetOf() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := t.SubSetOf(got, []interface{}{1, 2, 3, 4, 5, 6, 7, 8},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = t.SubSetOf(got, []interface{}{Between(1, 4), 3, Between(2, 10), Gt(100)},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SuperBagOf() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperBagOf(got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = t.SuperBagOf(got, []interface{}{Gt(5), Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SuperMapOf_map() {
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15)},
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
}

func ExampleT_SuperMapOf_typedMap() {
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15)},
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = t.SuperMapOf(&got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15)},
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_SuperSetOf() {
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperSetOf(got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = t.SuperSetOf(got, []interface{}{Gt(5), Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
}

func ExampleT_TruncTime() {
	t := NewT(&testing.T{})

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

func ExampleT_Zero() {
	t := NewT(&testing.T{})

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

	ok = t.CmpDeeply([3]int{0, 1}, Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = t.Zero(bytes.Buffer{})
	fmt.Println(ok)

	ok = t.Zero(&bytes.Buffer{}) // fails, as pointer not nil
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
}
