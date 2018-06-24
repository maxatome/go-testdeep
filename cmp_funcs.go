// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

package testdeep

import (
	"time"
)

// CmpAll is a shortcut for:
//
//   CmpDeeply(t, got, All(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpAll(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, All(expectedValues...), args...)
}

// CmpAny is a shortcut for:
//
//   CmpDeeply(t, got, Any(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpAny(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Any(expectedValues...), args...)
}

// CmpArray is a shortcut for:
//
//   CmpDeeply(t, got, Array(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpArray(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Array(model, expectedEntries), args...)
}

// CmpArrayEach is a shortcut for:
//
//   CmpDeeply(t, got, ArrayEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpArrayEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, ArrayEach(expectedValue), args...)
}

// CmpBag is a shortcut for:
//
//   CmpDeeply(t, got, Bag(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpBag(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Bag(expectedItems...), args...)
}

// CmpBetween is a shortcut for:
//
//   CmpDeeply(t, got, Between(from, to, bounds), args...)
//
// Between() optional parameter "bounds" is here mandatory.
// BoundsInIn value should be passed to mimic its absence in
// original Between() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpBetween(t TestingT, got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Between(from, to, bounds), args...)
}

// CmpCap is a shortcut for:
//
//   CmpDeeply(t, got, Cap(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpCap(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Cap(val), args...)
}

// CmpCode is a shortcut for:
//
//   CmpDeeply(t, got, Code(fn), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpCode(t TestingT, got interface{}, fn interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Code(fn), args...)
}

// CmpContains is a shortcut for:
//
//   CmpDeeply(t, got, Contains(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpContains(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Contains(expectedValue), args...)
}

// CmpEmpty is a shortcut for:
//
//   CmpDeeply(t, got, Empty(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpEmpty(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Empty(), args...)
}

// CmpGt is a shortcut for:
//
//   CmpDeeply(t, got, Gt(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpGt(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Gt(val), args...)
}

// CmpGte is a shortcut for:
//
//   CmpDeeply(t, got, Gte(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpGte(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Gte(val), args...)
}

// CmpHasPrefix is a shortcut for:
//
//   CmpDeeply(t, got, HasPrefix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpHasPrefix(t TestingT, got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, HasPrefix(expected), args...)
}

// CmpHasSuffix is a shortcut for:
//
//   CmpDeeply(t, got, HasSuffix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpHasSuffix(t TestingT, got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, HasSuffix(expected), args...)
}

// CmpIsa is a shortcut for:
//
//   CmpDeeply(t, got, Isa(model), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpIsa(t TestingT, got interface{}, model interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Isa(model), args...)
}

// CmpLen is a shortcut for:
//
//   CmpDeeply(t, got, Len(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpLen(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Len(val), args...)
}

// CmpLt is a shortcut for:
//
//   CmpDeeply(t, got, Lt(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpLt(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Lt(val), args...)
}

// CmpLte is a shortcut for:
//
//   CmpDeeply(t, got, Lte(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpLte(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Lte(val), args...)
}

// CmpMap is a shortcut for:
//
//   CmpDeeply(t, got, Map(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpMap(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Map(model, expectedEntries), args...)
}

// CmpMapEach is a shortcut for:
//
//   CmpDeeply(t, got, MapEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpMapEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, MapEach(expectedValue), args...)
}

// CmpN is a shortcut for:
//
//   CmpDeeply(t, got, N(num, tolerance), args...)
//
// N() optional parameter "tolerance" is here mandatory.
// 0 value should be passed to mimic its absence in
// original N() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpN(t TestingT, got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, N(num, tolerance), args...)
}

// CmpNil is a shortcut for:
//
//   CmpDeeply(t, got, Nil(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNil(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Nil(), args...)
}

// CmpNone is a shortcut for:
//
//   CmpDeeply(t, got, None(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNone(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, None(expectedValues...), args...)
}

// CmpNot is a shortcut for:
//
//   CmpDeeply(t, got, Not(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNot(t TestingT, got interface{}, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Not(expected), args...)
}

// CmpNotAny is a shortcut for:
//
//   CmpDeeply(t, got, NotAny(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNotAny(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, NotAny(expectedItems...), args...)
}

// CmpNotEmpty is a shortcut for:
//
//   CmpDeeply(t, got, NotEmpty(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNotEmpty(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, NotEmpty(), args...)
}

// CmpNotNil is a shortcut for:
//
//   CmpDeeply(t, got, NotNil(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNotNil(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, NotNil(), args...)
}

// CmpNotZero is a shortcut for:
//
//   CmpDeeply(t, got, NotZero(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpNotZero(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, NotZero(), args...)
}

// CmpPPtr is a shortcut for:
//
//   CmpDeeply(t, got, PPtr(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpPPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, PPtr(val), args...)
}

// CmpPtr is a shortcut for:
//
//   CmpDeeply(t, got, Ptr(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Ptr(val), args...)
}

// CmpRe is a shortcut for:
//
//   CmpDeeply(t, got, Re(reg, capture), args...)
//
// Re() optional parameter "capture" is here mandatory.
// nil value should be passed to mimic its absence in
// original Re() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpRe(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Re(reg, capture), args...)
}

// CmpReAll is a shortcut for:
//
//   CmpDeeply(t, got, ReAll(reg, capture), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpReAll(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, ReAll(reg, capture), args...)
}

// CmpSet is a shortcut for:
//
//   CmpDeeply(t, got, Set(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSet(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Set(expectedItems...), args...)
}

// CmpShallow is a shortcut for:
//
//   CmpDeeply(t, got, Shallow(expectedPtr), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpShallow(t TestingT, got interface{}, expectedPtr interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Shallow(expectedPtr), args...)
}

// CmpSlice is a shortcut for:
//
//   CmpDeeply(t, got, Slice(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSlice(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Slice(model, expectedEntries), args...)
}

// CmpSmuggle is a shortcut for:
//
//   CmpDeeply(t, got, Smuggle(fn, expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSmuggle(t TestingT, got interface{}, fn interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Smuggle(fn, expectedValue), args...)
}

// CmpString is a shortcut for:
//
//   CmpDeeply(t, got, String(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpString(t TestingT, got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, String(expected), args...)
}

// CmpStruct is a shortcut for:
//
//   CmpDeeply(t, got, Struct(model, expectedFields), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpStruct(t TestingT, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Struct(model, expectedFields), args...)
}

// CmpSubBagOf is a shortcut for:
//
//   CmpDeeply(t, got, SubBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSubBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SubBagOf(expectedItems...), args...)
}

// CmpSubMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SubMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSubMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SubMapOf(model, expectedEntries), args...)
}

// CmpSubSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SubSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSubSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SubSetOf(expectedItems...), args...)
}

// CmpSuperBagOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSuperBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SuperBagOf(expectedItems...), args...)
}

// CmpSuperMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSuperMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SuperMapOf(model, expectedEntries), args...)
}

// CmpSuperSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpSuperSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, SuperSetOf(expectedItems...), args...)
}

// CmpTruncTime is a shortcut for:
//
//   CmpDeeply(t, got, TruncTime(expectedTime, trunc), args...)
//
// TruncTime() optional parameter "trunc" is here mandatory.
// 0 value should be passed to mimic its absence in
// original TruncTime() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpTruncTime(t TestingT, got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, TruncTime(expectedTime, trunc), args...)
}

// CmpZero is a shortcut for:
//
//   CmpDeeply(t, got, Zero(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func CmpZero(t TestingT, got interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t, got, Zero(), args...)
}
