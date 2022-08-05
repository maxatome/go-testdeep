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

// allOperators lists the 62 operators.
// nil means not usable in JSON().
var allOperators = map[string]any{
	"All":          All,
	"Any":          Any,
	"Array":        nil,
	"ArrayEach":    ArrayEach,
	"Bag":          Bag,
	"Between":      Between,
	"Cap":          nil,
	"Catch":        nil,
	"Code":         nil,
	"Contains":     Contains,
	"ContainsKey":  ContainsKey,
	"Delay":        nil,
	"Empty":        Empty,
	"Gt":           Gt,
	"Gte":          Gte,
	"HasPrefix":    HasPrefix,
	"HasSuffix":    HasSuffix,
	"Ignore":       Ignore,
	"Isa":          nil,
	"JSON":         nil,
	"JSONPointer":  JSONPointer,
	"Keys":         Keys,
	"Lax":          nil,
	"Len":          Len,
	"Lt":           Lt,
	"Lte":          Lte,
	"Map":          nil,
	"MapEach":      MapEach,
	"N":            N,
	"NaN":          NaN,
	"Nil":          Nil,
	"None":         None,
	"Not":          Not,
	"NotAny":       NotAny,
	"NotEmpty":     NotEmpty,
	"NotNaN":       NotNaN,
	"NotNil":       NotNil,
	"NotZero":      NotZero,
	"PPtr":         nil,
	"Ptr":          nil,
	"Re":           Re,
	"ReAll":        ReAll,
	"SStruct":      nil,
	"Set":          Set,
	"Shallow":      nil,
	"Slice":        nil,
	"Smuggle":      nil,
	"String":       nil,
	"Struct":       nil,
	"SubBagOf":     SubBagOf,
	"SubJSONOf":    nil,
	"SubMapOf":     SubMapOf,
	"SubSetOf":     SubSetOf,
	"SuperBagOf":   SuperBagOf,
	"SuperJSONOf":  nil,
	"SuperMapOf":   SuperMapOf,
	"SuperSetOf":   SuperSetOf,
	"SuperSliceOf": nil,
	"Tag":          nil,
	"TruncTime":    nil,
	"Values":       Values,
	"Zero":         Zero,
}

// CmpAll is a shortcut for:
//
//	td.Cmp(t, got, td.All(expectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#All for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpAll(t TestingT, got any, expectedValues []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, All(expectedValues...), args...)
}

// CmpAny is a shortcut for:
//
//	td.Cmp(t, got, td.Any(expectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Any for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpAny(t TestingT, got any, expectedValues []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Any(expectedValues...), args...)
}

// CmpArray is a shortcut for:
//
//	td.Cmp(t, got, td.Array(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Array for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpArray(t TestingT, got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Array(model, expectedEntries), args...)
}

// CmpArrayEach is a shortcut for:
//
//	td.Cmp(t, got, td.ArrayEach(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ArrayEach for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpArrayEach(t TestingT, got, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, ArrayEach(expectedValue), args...)
}

// CmpBag is a shortcut for:
//
//	td.Cmp(t, got, td.Bag(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Bag for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpBag(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Bag(expectedItems...), args...)
}

// CmpBetween is a shortcut for:
//
//	td.Cmp(t, got, td.Between(from, to, bounds), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Between for details.
//
// Between() optional parameter "bounds" is here mandatory.
// [BoundsInIn] value should be passed to mimic its absence in
// original Between() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpBetween(t TestingT, got, from, to any, bounds BoundsKind, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Between(from, to, bounds), args...)
}

// CmpCap is a shortcut for:
//
//	td.Cmp(t, got, td.Cap(expectedCap), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cap for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpCap(t TestingT, got, expectedCap any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Cap(expectedCap), args...)
}

// CmpCode is a shortcut for:
//
//	td.Cmp(t, got, td.Code(fn), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Code for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpCode(t TestingT, got, fn any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Code(fn), args...)
}

// CmpContains is a shortcut for:
//
//	td.Cmp(t, got, td.Contains(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Contains for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpContains(t TestingT, got, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Contains(expectedValue), args...)
}

// CmpContainsKey is a shortcut for:
//
//	td.Cmp(t, got, td.ContainsKey(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContainsKey for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpContainsKey(t TestingT, got, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, ContainsKey(expectedValue), args...)
}

// CmpEmpty is a shortcut for:
//
//	td.Cmp(t, got, td.Empty(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Empty for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpEmpty(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Empty(), args...)
}

// CmpGt is a shortcut for:
//
//	td.Cmp(t, got, td.Gt(minExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Gt for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpGt(t TestingT, got, minExpectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Gt(minExpectedValue), args...)
}

// CmpGte is a shortcut for:
//
//	td.Cmp(t, got, td.Gte(minExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Gte for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpGte(t TestingT, got, minExpectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Gte(minExpectedValue), args...)
}

// CmpHasPrefix is a shortcut for:
//
//	td.Cmp(t, got, td.HasPrefix(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasPrefix for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpHasPrefix(t TestingT, got any, expected string, args ...any) bool {
	t.Helper()
	return Cmp(t, got, HasPrefix(expected), args...)
}

// CmpHasSuffix is a shortcut for:
//
//	td.Cmp(t, got, td.HasSuffix(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasSuffix for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpHasSuffix(t TestingT, got any, expected string, args ...any) bool {
	t.Helper()
	return Cmp(t, got, HasSuffix(expected), args...)
}

// CmpIsa is a shortcut for:
//
//	td.Cmp(t, got, td.Isa(model), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Isa for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpIsa(t TestingT, got, model any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Isa(model), args...)
}

// CmpJSON is a shortcut for:
//
//	td.Cmp(t, got, td.JSON(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSON for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpJSON(t TestingT, got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, JSON(expectedJSON, params...), args...)
}

// CmpJSONPointer is a shortcut for:
//
//	td.Cmp(t, got, td.JSONPointer(pointer, expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSONPointer for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpJSONPointer(t TestingT, got any, pointer string, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, JSONPointer(pointer, expectedValue), args...)
}

// CmpKeys is a shortcut for:
//
//	td.Cmp(t, got, td.Keys(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Keys for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpKeys(t TestingT, got, val any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Keys(val), args...)
}

// CmpLax is a shortcut for:
//
//	td.Cmp(t, got, td.Lax(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lax for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpLax(t TestingT, got, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Lax(expectedValue), args...)
}

// CmpLen is a shortcut for:
//
//	td.Cmp(t, got, td.Len(expectedLen), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Len for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpLen(t TestingT, got, expectedLen any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Len(expectedLen), args...)
}

// CmpLt is a shortcut for:
//
//	td.Cmp(t, got, td.Lt(maxExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lt for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpLt(t TestingT, got, maxExpectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Lt(maxExpectedValue), args...)
}

// CmpLte is a shortcut for:
//
//	td.Cmp(t, got, td.Lte(maxExpectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lte for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpLte(t TestingT, got, maxExpectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Lte(maxExpectedValue), args...)
}

// CmpMap is a shortcut for:
//
//	td.Cmp(t, got, td.Map(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Map for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpMap(t TestingT, got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Map(model, expectedEntries), args...)
}

// CmpMapEach is a shortcut for:
//
//	td.Cmp(t, got, td.MapEach(expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#MapEach for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpMapEach(t TestingT, got, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, MapEach(expectedValue), args...)
}

// CmpN is a shortcut for:
//
//	td.Cmp(t, got, td.N(num, tolerance), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#N for details.
//
// N() optional parameter "tolerance" is here mandatory.
// 0 value should be passed to mimic its absence in
// original N() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpN(t TestingT, got, num, tolerance any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, N(num, tolerance), args...)
}

// CmpNaN is a shortcut for:
//
//	td.Cmp(t, got, td.NaN(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NaN for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNaN(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NaN(), args...)
}

// CmpNil is a shortcut for:
//
//	td.Cmp(t, got, td.Nil(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Nil for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNil(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Nil(), args...)
}

// CmpNone is a shortcut for:
//
//	td.Cmp(t, got, td.None(notExpectedValues...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#None for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNone(t TestingT, got any, notExpectedValues []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, None(notExpectedValues...), args...)
}

// CmpNot is a shortcut for:
//
//	td.Cmp(t, got, td.Not(notExpected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Not for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNot(t TestingT, got, notExpected any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Not(notExpected), args...)
}

// CmpNotAny is a shortcut for:
//
//	td.Cmp(t, got, td.NotAny(notExpectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotAny for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNotAny(t TestingT, got any, notExpectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NotAny(notExpectedItems...), args...)
}

// CmpNotEmpty is a shortcut for:
//
//	td.Cmp(t, got, td.NotEmpty(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotEmpty for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNotEmpty(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NotEmpty(), args...)
}

// CmpNotNaN is a shortcut for:
//
//	td.Cmp(t, got, td.NotNaN(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNaN for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNotNaN(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NotNaN(), args...)
}

// CmpNotNil is a shortcut for:
//
//	td.Cmp(t, got, td.NotNil(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNil for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNotNil(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NotNil(), args...)
}

// CmpNotZero is a shortcut for:
//
//	td.Cmp(t, got, td.NotZero(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotZero for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpNotZero(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, NotZero(), args...)
}

// CmpPPtr is a shortcut for:
//
//	td.Cmp(t, got, td.PPtr(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#PPtr for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpPPtr(t TestingT, got, val any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, PPtr(val), args...)
}

// CmpPtr is a shortcut for:
//
//	td.Cmp(t, got, td.Ptr(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Ptr for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpPtr(t TestingT, got, val any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Ptr(val), args...)
}

// CmpRe is a shortcut for:
//
//	td.Cmp(t, got, td.Re(reg, capture), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Re for details.
//
// Re() optional parameter "capture" is here mandatory.
// nil value should be passed to mimic its absence in
// original Re() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpRe(t TestingT, got, reg, capture any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Re(reg, capture), args...)
}

// CmpReAll is a shortcut for:
//
//	td.Cmp(t, got, td.ReAll(reg, capture), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ReAll for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpReAll(t TestingT, got, reg, capture any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, ReAll(reg, capture), args...)
}

// CmpSet is a shortcut for:
//
//	td.Cmp(t, got, td.Set(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Set for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSet(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Set(expectedItems...), args...)
}

// CmpShallow is a shortcut for:
//
//	td.Cmp(t, got, td.Shallow(expectedPtr), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Shallow for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpShallow(t TestingT, got, expectedPtr any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Shallow(expectedPtr), args...)
}

// CmpSlice is a shortcut for:
//
//	td.Cmp(t, got, td.Slice(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Slice for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSlice(t TestingT, got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Slice(model, expectedEntries), args...)
}

// CmpSmuggle is a shortcut for:
//
//	td.Cmp(t, got, td.Smuggle(fn, expectedValue), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Smuggle for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSmuggle(t TestingT, got, fn, expectedValue any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Smuggle(fn, expectedValue), args...)
}

// CmpSStruct is a shortcut for:
//
//	td.Cmp(t, got, td.SStruct(model, expectedFields), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SStruct for details.
//
// SStruct() optional parameter "expectedFields" is here mandatory.
// nil value should be passed to mimic its absence in
// original SStruct() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSStruct(t TestingT, got, model any, expectedFields StructFields, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SStruct(model, expectedFields), args...)
}

// CmpString is a shortcut for:
//
//	td.Cmp(t, got, td.String(expected), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#String for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpString(t TestingT, got any, expected string, args ...any) bool {
	t.Helper()
	return Cmp(t, got, String(expected), args...)
}

// CmpStruct is a shortcut for:
//
//	td.Cmp(t, got, td.Struct(model, expectedFields), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Struct for details.
//
// Struct() optional parameter "expectedFields" is here mandatory.
// nil value should be passed to mimic its absence in
// original Struct() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpStruct(t TestingT, got, model any, expectedFields StructFields, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Struct(model, expectedFields), args...)
}

// CmpSubBagOf is a shortcut for:
//
//	td.Cmp(t, got, td.SubBagOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubBagOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSubBagOf(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SubBagOf(expectedItems...), args...)
}

// CmpSubJSONOf is a shortcut for:
//
//	td.Cmp(t, got, td.SubJSONOf(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubJSONOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSubJSONOf(t TestingT, got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SubJSONOf(expectedJSON, params...), args...)
}

// CmpSubMapOf is a shortcut for:
//
//	td.Cmp(t, got, td.SubMapOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubMapOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSubMapOf(t TestingT, got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SubMapOf(model, expectedEntries), args...)
}

// CmpSubSetOf is a shortcut for:
//
//	td.Cmp(t, got, td.SubSetOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubSetOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSubSetOf(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SubSetOf(expectedItems...), args...)
}

// CmpSuperBagOf is a shortcut for:
//
//	td.Cmp(t, got, td.SuperBagOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperBagOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSuperBagOf(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SuperBagOf(expectedItems...), args...)
}

// CmpSuperJSONOf is a shortcut for:
//
//	td.Cmp(t, got, td.SuperJSONOf(expectedJSON, params...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperJSONOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSuperJSONOf(t TestingT, got, expectedJSON any, params []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SuperJSONOf(expectedJSON, params...), args...)
}

// CmpSuperMapOf is a shortcut for:
//
//	td.Cmp(t, got, td.SuperMapOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperMapOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSuperMapOf(t TestingT, got, model any, expectedEntries MapEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SuperMapOf(model, expectedEntries), args...)
}

// CmpSuperSetOf is a shortcut for:
//
//	td.Cmp(t, got, td.SuperSetOf(expectedItems...), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSetOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSuperSetOf(t TestingT, got any, expectedItems []any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SuperSetOf(expectedItems...), args...)
}

// CmpSuperSliceOf is a shortcut for:
//
//	td.Cmp(t, got, td.SuperSliceOf(model, expectedEntries), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSliceOf for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpSuperSliceOf(t TestingT, got, model any, expectedEntries ArrayEntries, args ...any) bool {
	t.Helper()
	return Cmp(t, got, SuperSliceOf(model, expectedEntries), args...)
}

// CmpTruncTime is a shortcut for:
//
//	td.Cmp(t, got, td.TruncTime(expectedTime, trunc), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#TruncTime for details.
//
// TruncTime() optional parameter "trunc" is here mandatory.
// 0 value should be passed to mimic its absence in
// original TruncTime() call.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpTruncTime(t TestingT, got, expectedTime any, trunc time.Duration, args ...any) bool {
	t.Helper()
	return Cmp(t, got, TruncTime(expectedTime, trunc), args...)
}

// CmpValues is a shortcut for:
//
//	td.Cmp(t, got, td.Values(val), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Values for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpValues(t TestingT, got, val any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Values(val), args...)
}

// CmpZero is a shortcut for:
//
//	td.Cmp(t, got, td.Zero(), args...)
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Zero for details.
//
// Returns true if the test is OK, false if it fails.
//
// If "t" is a *T then its Config is inherited.
//
// args... are optional and allow to name the test. This name is
// used in case of failure to qualify the test. If len(args) > 1 and
// the first item of args is a string and contains a '%' rune then
// [fmt.Fprintf] is used to compose the name, else args are passed to
// [fmt.Fprint]. Do not forget it is the name of the test, not the
// reason of a potential failure.
func CmpZero(t TestingT, got any, args ...any) bool {
	t.Helper()
	return Cmp(t, got, Zero(), args...)
}
