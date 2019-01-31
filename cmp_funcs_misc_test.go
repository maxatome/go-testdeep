// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
)

func ExampleCmpTrue() {
	t := &testing.T{}

	got := true
	ok := CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	got = false
	ok = CmpTrue(t, got, "check that got is true!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpFalse() {
	t := &testing.T{}

	got := false
	ok := CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	got = true
	ok = CmpFalse(t, got, "check that got is false!")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := CmpError(t, got, "An error occurred")
	fmt.Println(ok)

	got = nil
	ok = CmpError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCmpNoError() {
	t := &testing.T{}

	got := fmt.Errorf("Error #%d", 42)
	ok := CmpNoError(t, got, "An error occurred") // fails
	fmt.Println(ok)

	got = nil
	ok = CmpNoError(t, got, "An error occurred")
	fmt.Println(ok)

	// Output:
	// false
	// true
}

func ExampleCmpPanic() {
	t := &testing.T{}

	ok := CmpPanic(t,
		func() { panic("I am panicking!") }, "I am panicking!",
		"Checks for panic")
	fmt.Println("checks exact panic() string:", ok)

	// Can use TestDeep operator too
	ok = CmpPanic(t,
		func() { panic("I am panicking!") }, Contains("panicking!"),
		"Checks for panic")
	fmt.Println("checks panic() sub-string:", ok)

	// Can detect panic(nil)
	ok = CmpPanic(t, func() { panic(nil) }, nil, "Checks for panic(nil)")
	fmt.Println("checks for panic(nil):", ok)

	// As well as structured data panic
	type PanicStruct struct {
		Error string
		Code  int
	}

	ok = CmpPanic(t,
		func() {
			panic(PanicStruct{Error: "Memory violation", Code: 11})
		},
		PanicStruct{
			Error: "Memory violation",
			Code:  11,
		})
	fmt.Println("checks exact panic() struct:", ok)

	// or combined with TestDeep operators too
	ok = CmpPanic(t,
		func() {
			panic(PanicStruct{Error: "Memory violation", Code: 11})
		},
		Struct(PanicStruct{}, StructFields{
			"Code": Between(10, 20),
		}))
	fmt.Println("checks panic() struct against TestDeep operators:", ok)

	// Of course, do not panic = test failure, even for expected nil
	// panic parameter
	ok = CmpPanic(t, func() {}, nil)
	fmt.Println("checks a panic occurred:", ok)

	// Output:
	// checks exact panic() string: true
	// checks panic() sub-string: true
	// checks for panic(nil): true
	// checks exact panic() struct: true
	// checks panic() struct against TestDeep operators: true
	// checks a panic occurred: false
}

func ExampleCmpNotPanic() {
	t := &testing.T{}

	ok := CmpNotPanic(t, func() {}, nil)
	fmt.Println("checks a panic DID NOT occur:", ok)

	// Classic panic
	ok = CmpNotPanic(t, func() { panic("I am panicking!") },
		"Hope it does not panic!")
	fmt.Println("still no panic?", ok)

	// Can detect panic(nil)
	ok = CmpNotPanic(t, func() { panic(nil) }, "Checks for panic(nil)")
	fmt.Println("last no panic?", ok)

	// Output:
	// checks a panic DID NOT occur: true
	// still no panic? false
	// last no panic? false
}

// Just to test the case where t is an interface and not a *testing.T
// See t.Helper() issue in all tested functions.
func TestCmpWithInterfaceT(tt *testing.T) {
	ttt := &test.TestingFT{}

	test.IsTrue(tt, CmpTrue(ttt, true, true))
	test.IsTrue(tt, CmpFalse(ttt, false, false))
	test.IsFalse(tt, CmpError(ttt, nil))
	test.IsFalse(tt, CmpNoError(ttt, errors.New("error")))
	test.IsFalse(tt, CmpPanic(ttt, func() {}, "panic")) // no panic occurred
	test.IsTrue(tt, CmpNotPanic(ttt, func() {}))

	test.IsFalse(tt, CmpAll(ttt, 0, []interface{}{12}))
	test.IsFalse(tt, CmpAny(ttt, 0, nil))
	test.IsFalse(tt, CmpArray(ttt, 0, [2]int{}, nil))
	test.IsFalse(tt, CmpArrayEach(ttt, 0, nil))
	test.IsFalse(tt, CmpBag(ttt, 0, nil))
	test.IsFalse(tt, CmpBetween(ttt, 0, 1, 2, BoundsInIn))
	test.IsFalse(tt, CmpCap(ttt, nil, 12))
	test.IsFalse(tt, CmpCode(ttt, 0, func(n int) bool { return false }))
	test.IsFalse(tt, CmpContains(ttt, 0, nil))
	test.IsFalse(tt, CmpContainsKey(ttt, map[bool]int{}, true))
	test.IsFalse(tt, CmpEmpty(ttt, 0))
	test.IsFalse(tt, CmpGt(ttt, 0, 12))
	test.IsFalse(tt, CmpGte(ttt, 0, 12))
	test.IsFalse(tt, CmpHasPrefix(ttt, 0, "pipo"))
	test.IsFalse(tt, CmpHasSuffix(ttt, 0, "pipo"))
	test.IsFalse(tt, CmpIsa(ttt, 0, "string"))
	test.IsFalse(tt, CmpLen(ttt, nil, 12))
	test.IsFalse(tt, CmpLt(ttt, 0, -12))
	test.IsFalse(tt, CmpLte(ttt, 0, -12))
	test.IsFalse(tt, CmpMap(ttt, 0, map[int]bool{}, nil))
	test.IsFalse(tt, CmpMapEach(ttt, 0, nil))
	test.IsFalse(tt, CmpN(ttt, 0, 12, 0))
	test.IsFalse(tt, CmpNaN(ttt, 0, nil))
	test.IsFalse(tt, CmpNil(ttt, 0))
	test.IsFalse(tt, CmpNone(ttt, 0, []interface{}{0}))
	test.IsFalse(tt, CmpNot(ttt, 0, 0))
	test.IsFalse(tt, CmpNotAny(ttt, 0, nil))
	test.IsFalse(tt, CmpNotEmpty(ttt, 0, nil))
	test.IsFalse(tt, CmpNotNaN(ttt, 0, nil))
	test.IsFalse(tt, CmpNotNil(ttt, nil))
	test.IsFalse(tt, CmpNotZero(ttt, 0, nil))
	test.IsFalse(tt, CmpPPtr(ttt, 0, 12))
	test.IsFalse(tt, CmpPtr(ttt, 0, 12))
	test.IsFalse(tt, CmpRe(ttt, 0, "pipo", nil))
	test.IsFalse(tt, CmpReAll(ttt, 0, "pipo", nil))
	test.IsFalse(tt, CmpSet(ttt, 0, nil))
	test.IsFalse(tt, CmpShallow(ttt, 0, []int{}))
	test.IsFalse(tt, CmpSlice(ttt, 0, []int{}, nil))
	test.IsFalse(tt, CmpSmuggle(ttt, 0, func(n int) int { return 0 }, 12))
	test.IsFalse(tt, CmpString(ttt, 0, "pipo"))
	test.IsFalse(tt, CmpStruct(ttt, 0, struct{}{}, nil))
	test.IsFalse(tt, CmpSubBagOf(ttt, 0, nil))
	test.IsFalse(tt, CmpSubMapOf(ttt, 0, map[int]bool{}, nil))
	test.IsFalse(tt, CmpSubSetOf(ttt, 0, nil))
	test.IsFalse(tt, CmpSuperBagOf(ttt, 0, nil))
	test.IsFalse(tt, CmpSuperMapOf(ttt, 0, map[int]bool{}, nil))
	test.IsFalse(tt, CmpSuperSetOf(ttt, 0, nil))
	test.IsFalse(tt, CmpTruncTime(ttt, 0, time.Now(), time.Second))
	test.IsFalse(tt, CmpZero(ttt, 12))
}
