// Copyright (c) 2018-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

package td

import (
	"time"
)

// All is a shortcut for:
//
//	t.Cmp(got, td.All(expectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#All for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) All(got any, expectedValues []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, All(expectedValues...), args...)
}

// Any is a shortcut for:
//
//	t.Cmp(got, td.Any(expectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Any for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Any(got any, expectedValues []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Any(expectedValues...), args...)
}

// Array is a shortcut for:
//
//	t.Cmp(got, td.Array(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Array for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Array(got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Array(model, expectedEntries), args...)
}

// ArrayEach is a shortcut for:
//
//	t.Cmp(got, td.ArrayEach(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ArrayEach for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) ArrayEach(got, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, ArrayEach(expectedValue), args...)
}

// Bag is a shortcut for:
//
//	t.Cmp(got, td.Bag(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Bag for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Bag(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Bag(expectedItems...), args...)
}

// Between is a shortcut for:
//
//	t.Cmp(got, td.Between(from, to, bounds), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Between for details.
//
// Between() optional parameter "bounds" is here mandatory.
// [BoundsInIn] value should be passed to mimic its absence in
// original Between() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Between(got, from, to any, bounds BoundsKind, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Between(from, to, bounds), args...)
}

// Cap is a shortcut for:
//
//	t.Cmp(got, td.Cap(expectedCap), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cap for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Cap(got, expectedCap any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Cap(expectedCap), args...)
}

// Code is a shortcut for:
//
//	t.Cmp(got, td.Code(fn), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Code for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Code(got, fn any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Code(fn), args...)
}

// Contains is a shortcut for:
//
//	t.Cmp(got, td.Contains(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Contains for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Contains(got, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Contains(expectedValue), args...)
}

// ContainsKey is a shortcut for:
//
//	t.Cmp(got, td.ContainsKey(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContainsKey for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) ContainsKey(got, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, ContainsKey(expectedValue), args...)
}

// Empty is a shortcut for:
//
//	t.Cmp(got, td.Empty(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Empty for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Empty(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Empty(), args...)
}

// Gt is a shortcut for:
//
//	t.Cmp(got, td.Gt(minExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Gt for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Gt(got, minExpectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Gt(minExpectedValue), args...)
}

// Gte is a shortcut for:
//
//	t.Cmp(got, td.Gte(minExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Gte for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Gte(got, minExpectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Gte(minExpectedValue), args...)
}

// HasPrefix is a shortcut for:
//
//	t.Cmp(got, td.HasPrefix(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasPrefix for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) HasPrefix(got any, expected string, args ...any) bool {
	t.Helper()
	return t.Cmp(got, HasPrefix(expected), args...)
}

// HasSuffix is a shortcut for:
//
//	t.Cmp(got, td.HasSuffix(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasSuffix for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) HasSuffix(got any, expected string, args ...any) bool {
	t.Helper()
	return t.Cmp(got, HasSuffix(expected), args...)
}

// Isa is a shortcut for:
//
//	t.Cmp(got, td.Isa(model), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Isa for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Isa(got, model any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Isa(model), args...)
}

// JSON is a shortcut for:
//
//	t.Cmp(got, td.JSON(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSON for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) JSON(got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, JSON(expectedJSON, params...), args...)
}

// JSONPointer is a shortcut for:
//
//	t.Cmp(got, td.JSONPointer(pointer, expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSONPointer for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) JSONPointer(got any, pointer string, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, JSONPointer(pointer, expectedValue), args...)
}

// Keys is a shortcut for:
//
//	t.Cmp(got, td.Keys(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Keys for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Keys(got, val any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Keys(val), args...)
}

// CmpLax is a shortcut for:
//
//	t.Cmp(got, td.Lax(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lax for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) CmpLax(got, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Lax(expectedValue), args...)
}

// Len is a shortcut for:
//
//	t.Cmp(got, td.Len(expectedLen), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Len for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Len(got, expectedLen any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Len(expectedLen), args...)
}

// Lt is a shortcut for:
//
//	t.Cmp(got, td.Lt(maxExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lt for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Lt(got, maxExpectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Lt(maxExpectedValue), args...)
}

// Lte is a shortcut for:
//
//	t.Cmp(got, td.Lte(maxExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lte for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Lte(got, maxExpectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Lte(maxExpectedValue), args...)
}

// Map is a shortcut for:
//
//	t.Cmp(got, td.Map(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Map for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Map(got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Map(model, expectedEntries), args...)
}

// MapEach is a shortcut for:
//
//	t.Cmp(got, td.MapEach(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#MapEach for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) MapEach(got, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, MapEach(expectedValue), args...)
}

// N is a shortcut for:
//
//	t.Cmp(got, td.N(num, tolerance), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#N for details.
//
// N() optional parameter "tolerance" is here mandatory.
// 0 value should be passed to mimic its absence in
// original N() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) N(got, num, tolerance any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, N(num, tolerance), args...)
}

// NaN is a shortcut for:
//
//	t.Cmp(got, td.NaN(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NaN for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NaN(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NaN(), args...)
}

// Nil is a shortcut for:
//
//	t.Cmp(got, td.Nil(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Nil for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Nil(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Nil(), args...)
}

// None is a shortcut for:
//
//	t.Cmp(got, td.None(notExpectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#None for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) None(got any, notExpectedValues []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, None(notExpectedValues...), args...)
}

// Not is a shortcut for:
//
//	t.Cmp(got, td.Not(notExpected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Not for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Not(got, notExpected any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Not(notExpected), args...)
}

// NotAny is a shortcut for:
//
//	t.Cmp(got, td.NotAny(notExpectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotAny for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NotAny(got any, notExpectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NotAny(notExpectedItems...), args...)
}

// NotEmpty is a shortcut for:
//
//	t.Cmp(got, td.NotEmpty(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotEmpty for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NotEmpty(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NotEmpty(), args...)
}

// NotNaN is a shortcut for:
//
//	t.Cmp(got, td.NotNaN(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNaN for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NotNaN(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NotNaN(), args...)
}

// NotNil is a shortcut for:
//
//	t.Cmp(got, td.NotNil(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNil for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NotNil(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NotNil(), args...)
}

// NotZero is a shortcut for:
//
//	t.Cmp(got, td.NotZero(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotZero for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) NotZero(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, NotZero(), args...)
}

// PPtr is a shortcut for:
//
//	t.Cmp(got, td.PPtr(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#PPtr for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) PPtr(got, val any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, PPtr(val), args...)
}

// Ptr is a shortcut for:
//
//	t.Cmp(got, td.Ptr(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Ptr for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Ptr(got, val any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Ptr(val), args...)
}

// Re is a shortcut for:
//
//	t.Cmp(got, td.Re(reg, capture), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Re for details.
//
// Re() optional parameter "capture" is here mandatory.
// nil value should be passed to mimic its absence in
// original Re() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Re(got, reg, capture any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Re(reg, capture), args...)
}

// ReAll is a shortcut for:
//
//	t.Cmp(got, td.ReAll(reg, capture), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ReAll for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) ReAll(got, reg, capture any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, ReAll(reg, capture), args...)
}

// Set is a shortcut for:
//
//	t.Cmp(got, td.Set(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Set for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Set(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Set(expectedItems...), args...)
}

// Shallow is a shortcut for:
//
//	t.Cmp(got, td.Shallow(expectedPtr), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Shallow for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Shallow(got, expectedPtr any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Shallow(expectedPtr), args...)
}

// Slice is a shortcut for:
//
//	t.Cmp(got, td.Slice(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Slice for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Slice(got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Slice(model, expectedEntries), args...)
}

// Smuggle is a shortcut for:
//
//	t.Cmp(got, td.Smuggle(fn, expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Smuggle for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Smuggle(got, fn, expectedValue any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Smuggle(fn, expectedValue), args...)
}

// SStruct is a shortcut for:
//
//	t.Cmp(got, td.SStruct(model, expectedFields), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SStruct for details.
//
// SStruct() optional parameter "expectedFields" is here mandatory.
// nil value should be passed to mimic its absence in
// original SStruct() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SStruct(got, model any, expectedFields StructFields, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SStruct(model, expectedFields), args...)
}

// String is a shortcut for:
//
//	t.Cmp(got, td.String(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#String for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) String(got any, expected string, args ...any) bool {
	t.Helper()
	return t.Cmp(got, String(expected), args...)
}

// Struct is a shortcut for:
//
//	t.Cmp(got, td.Struct(model, expectedFields), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Struct for details.
//
// Struct() optional parameter "expectedFields" is here mandatory.
// nil value should be passed to mimic its absence in
// original Struct() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Struct(got, model any, expectedFields StructFields, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Struct(model, expectedFields), args...)
}

// SubBagOf is a shortcut for:
//
//	t.Cmp(got, td.SubBagOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubBagOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SubBagOf(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SubBagOf(expectedItems...), args...)
}

// SubJSONOf is a shortcut for:
//
//	t.Cmp(got, td.SubJSONOf(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubJSONOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SubJSONOf(got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SubJSONOf(expectedJSON, params...), args...)
}

// SubMapOf is a shortcut for:
//
//	t.Cmp(got, td.SubMapOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubMapOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SubMapOf(got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SubMapOf(model, expectedEntries), args...)
}

// SubSetOf is a shortcut for:
//
//	t.Cmp(got, td.SubSetOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubSetOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SubSetOf(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SubSetOf(expectedItems...), args...)
}

// SuperBagOf is a shortcut for:
//
//	t.Cmp(got, td.SuperBagOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperBagOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SuperBagOf(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SuperBagOf(expectedItems...), args...)
}

// SuperJSONOf is a shortcut for:
//
//	t.Cmp(got, td.SuperJSONOf(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperJSONOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SuperJSONOf(got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SuperJSONOf(expectedJSON, params...), args...)
}

// SuperMapOf is a shortcut for:
//
//	t.Cmp(got, td.SuperMapOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperMapOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SuperMapOf(got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SuperMapOf(model, expectedEntries), args...)
}

// SuperSetOf is a shortcut for:
//
//	t.Cmp(got, td.SuperSetOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSetOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SuperSetOf(got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SuperSetOf(expectedItems...), args...)
}

// SuperSliceOf is a shortcut for:
//
//	t.Cmp(got, td.SuperSliceOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSliceOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) SuperSliceOf(got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return t.Cmp(got, SuperSliceOf(model, expectedEntries), args...)
}

// TruncTime is a shortcut for:
//
//	t.Cmp(got, td.TruncTime(expectedTime, trunc), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#TruncTime for details.
//
// TruncTime() optional parameter "trunc" is here mandatory.
// 0 value should be passed to mimic its absence in
// original TruncTime() call.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) TruncTime(got, expectedTime any, trunc time.Duration, args ...any) bool {
	t.Helper()
	return t.Cmp(got, TruncTime(expectedTime, trunc), args...)
}

// Values is a shortcut for:
//
//	t.Cmp(got, td.Values(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Values for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Values(got, val any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Values(val), args...)
}

// Zero is a shortcut for:
//
//	t.Cmp(got, td.Zero(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Zero for details.
//
// Returns true if the test is OK, false if it fails.
//
// "args..." are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of "args" is a string and contains a '%' rune then
// fmt.Fprintf is used to compose the name, else "args" are passed to
// fmt.Fprint. Do not forget it is the name of the test, not the
// reason of a potential failure.
func (t *T) Zero(got any, args ...any) bool {
	t.Helper()
	return t.Cmp(got, Zero(), args...)
}
