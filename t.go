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

// All is a shortcut for:
//
//   t.CmpDeeply(got, All(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) All(got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, All(expectedValues...), args...)
}

// Any is a shortcut for:
//
//   t.CmpDeeply(got, Any(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Any(got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Any(expectedValues...), args...)
}

// Array is a shortcut for:
//
//   t.CmpDeeply(got, Array(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Array(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Array(model, expectedEntries), args...)
}

// ArrayEach is a shortcut for:
//
//   t.CmpDeeply(got, ArrayEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) ArrayEach(got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, ArrayEach(expectedValue), args...)
}

// Bag is a shortcut for:
//
//   t.CmpDeeply(got, Bag(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Bag(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Bag(expectedItems...), args...)
}

// Between is a shortcut for:
//
//   t.CmpDeeply(got, Between(from, to, bounds), args...)
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
func (t *T) Between(got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Between(from, to, bounds), args...)
}

// Cap is a shortcut for:
//
//   t.CmpDeeply(got, Cap(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Cap(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Cap(val), args...)
}

// Code is a shortcut for:
//
//   t.CmpDeeply(got, Code(fn), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Code(got interface{}, fn interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Code(fn), args...)
}

// Contains is a shortcut for:
//
//   t.CmpDeeply(got, Contains(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Contains(got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Contains(expectedValue), args...)
}

// Empty is a shortcut for:
//
//   t.CmpDeeply(got, Empty(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Empty(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Empty(), args...)
}

// Gt is a shortcut for:
//
//   t.CmpDeeply(got, Gt(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Gt(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Gt(val), args...)
}

// Gte is a shortcut for:
//
//   t.CmpDeeply(got, Gte(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Gte(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Gte(val), args...)
}

// HasPrefix is a shortcut for:
//
//   t.CmpDeeply(got, HasPrefix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) HasPrefix(got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, HasPrefix(expected), args...)
}

// HasSuffix is a shortcut for:
//
//   t.CmpDeeply(got, HasSuffix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) HasSuffix(got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, HasSuffix(expected), args...)
}

// Isa is a shortcut for:
//
//   t.CmpDeeply(got, Isa(model), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Isa(got interface{}, model interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Isa(model), args...)
}

// Len is a shortcut for:
//
//   t.CmpDeeply(got, Len(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Len(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Len(val), args...)
}

// Lt is a shortcut for:
//
//   t.CmpDeeply(got, Lt(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Lt(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Lt(val), args...)
}

// Lte is a shortcut for:
//
//   t.CmpDeeply(got, Lte(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Lte(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Lte(val), args...)
}

// Map is a shortcut for:
//
//   t.CmpDeeply(got, Map(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Map(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Map(model, expectedEntries), args...)
}

// MapEach is a shortcut for:
//
//   t.CmpDeeply(got, MapEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) MapEach(got interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, MapEach(expectedValue), args...)
}

// N is a shortcut for:
//
//   t.CmpDeeply(got, N(num, tolerance), args...)
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
func (t *T) N(got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, N(num, tolerance), args...)
}

// NaN is a shortcut for:
//
//   t.CmpDeeply(got, NaN(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NaN(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NaN(), args...)
}

// Nil is a shortcut for:
//
//   t.CmpDeeply(got, Nil(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Nil(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Nil(), args...)
}

// None is a shortcut for:
//
//   t.CmpDeeply(got, None(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) None(got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, None(expectedValues...), args...)
}

// Not is a shortcut for:
//
//   t.CmpDeeply(got, Not(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Not(got interface{}, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Not(expected), args...)
}

// NotAny is a shortcut for:
//
//   t.CmpDeeply(got, NotAny(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NotAny(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NotAny(expectedItems...), args...)
}

// NotEmpty is a shortcut for:
//
//   t.CmpDeeply(got, NotEmpty(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NotEmpty(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NotEmpty(), args...)
}

// NotNaN is a shortcut for:
//
//   t.CmpDeeply(got, NotNaN(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NotNaN(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NotNaN(), args...)
}

// NotNil is a shortcut for:
//
//   t.CmpDeeply(got, NotNil(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NotNil(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NotNil(), args...)
}

// NotZero is a shortcut for:
//
//   t.CmpDeeply(got, NotZero(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) NotZero(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, NotZero(), args...)
}

// PPtr is a shortcut for:
//
//   t.CmpDeeply(got, PPtr(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) PPtr(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, PPtr(val), args...)
}

// Ptr is a shortcut for:
//
//   t.CmpDeeply(got, Ptr(val), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Ptr(got interface{}, val interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Ptr(val), args...)
}

// Re is a shortcut for:
//
//   t.CmpDeeply(got, Re(reg, capture), args...)
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
func (t *T) Re(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Re(reg, capture), args...)
}

// ReAll is a shortcut for:
//
//   t.CmpDeeply(got, ReAll(reg, capture), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) ReAll(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, ReAll(reg, capture), args...)
}

// Set is a shortcut for:
//
//   t.CmpDeeply(got, Set(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Set(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Set(expectedItems...), args...)
}

// Shallow is a shortcut for:
//
//   t.CmpDeeply(got, Shallow(expectedPtr), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Shallow(got interface{}, expectedPtr interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Shallow(expectedPtr), args...)
}

// Slice is a shortcut for:
//
//   t.CmpDeeply(got, Slice(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Slice(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Slice(model, expectedEntries), args...)
}

// Smuggle is a shortcut for:
//
//   t.CmpDeeply(got, Smuggle(fn, expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Smuggle(got interface{}, fn interface{}, expectedValue interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Smuggle(fn, expectedValue), args...)
}

// String is a shortcut for:
//
//   t.CmpDeeply(got, String(expected), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) String(got interface{}, expected string, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, String(expected), args...)
}

// Struct is a shortcut for:
//
//   t.CmpDeeply(got, Struct(model, expectedFields), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Struct(got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Struct(model, expectedFields), args...)
}

// SubBagOf is a shortcut for:
//
//   t.CmpDeeply(got, SubBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SubBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SubBagOf(expectedItems...), args...)
}

// SubMapOf is a shortcut for:
//
//   t.CmpDeeply(got, SubMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SubMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SubMapOf(model, expectedEntries), args...)
}

// SubSetOf is a shortcut for:
//
//   t.CmpDeeply(got, SubSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SubSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SubSetOf(expectedItems...), args...)
}

// SuperBagOf is a shortcut for:
//
//   t.CmpDeeply(got, SuperBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SuperBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SuperBagOf(expectedItems...), args...)
}

// SuperMapOf is a shortcut for:
//
//   t.CmpDeeply(got, SuperMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SuperMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SuperMapOf(model, expectedEntries), args...)
}

// SuperSetOf is a shortcut for:
//
//   t.CmpDeeply(got, SuperSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) SuperSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, SuperSetOf(expectedItems...), args...)
}

// TruncTime is a shortcut for:
//
//   t.CmpDeeply(got, TruncTime(expectedTime, trunc), args...)
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
func (t *T) TruncTime(got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, TruncTime(expectedTime, trunc), args...)
}

// Zero is a shortcut for:
//
//   t.CmpDeeply(got, Zero(), args...)
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
func (t *T) Zero(got interface{}, args ...interface{}) bool {
	t.Helper()
	return t.CmpDeeply(got, Zero(), args...)
}
