package testdeep

// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

import (
	"testing"
	"time"
)

// CmpAll is a shortcut for:
//
//   CmpDeeply(t, got, All(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpAll(t *testing.T, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, All(expectedValues...), args...)
}

// CmpAny is a shortcut for:
//
//   CmpDeeply(t, got, Any(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpAny(t *testing.T, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Any(expectedValues...), args...)
}

// CmpArray is a shortcut for:
//
//   CmpDeeply(t, got, Array(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpArray(t *testing.T, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Array(model, expectedEntries), args...)
}

// CmpArrayEach is a shortcut for:
//
//   CmpDeeply(t, got, ArrayEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpArrayEach(t *testing.T, got interface{}, expectedValue interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, ArrayEach(expectedValue), args...)
}

// CmpBag is a shortcut for:
//
//   CmpDeeply(t, got, Bag(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpBag(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
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
func CmpBetween(t *testing.T, got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool {
	return CmpDeeply(t, got, Between(from, to, bounds), args...)
}

// CmpCap is a shortcut for:
//
//   CmpDeeply(t, got, Cap(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpCap(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Cap(val), args...)
}

// CmpCode is a shortcut for:
//
//   CmpDeeply(t, got, Code(fn), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpCode(t *testing.T, got interface{}, fn interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Code(fn), args...)
}

// CmpContains is a shortcut for:
//
//   CmpDeeply(t, got, Contains(expected), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpContains(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, Contains(expected), args...)
}

// CmpGt is a shortcut for:
//
//   CmpDeeply(t, got, Gt(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpGt(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Gt(val), args...)
}

// CmpGte is a shortcut for:
//
//   CmpDeeply(t, got, Gte(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpGte(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Gte(val), args...)
}

// CmpHasPrefix is a shortcut for:
//
//   CmpDeeply(t, got, HasPrefix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpHasPrefix(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, HasPrefix(expected), args...)
}

// CmpHasSuffix is a shortcut for:
//
//   CmpDeeply(t, got, HasSuffix(expected), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpHasSuffix(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, HasSuffix(expected), args...)
}

// CmpIsa is a shortcut for:
//
//   CmpDeeply(t, got, Isa(model), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpIsa(t *testing.T, got interface{}, model interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Isa(model), args...)
}

// CmpLen is a shortcut for:
//
//   CmpDeeply(t, got, Len(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpLen(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Len(val), args...)
}

// CmpLt is a shortcut for:
//
//   CmpDeeply(t, got, Lt(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpLt(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Lt(val), args...)
}

// CmpLte is a shortcut for:
//
//   CmpDeeply(t, got, Lte(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpLte(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Lte(val), args...)
}

// CmpMap is a shortcut for:
//
//   CmpDeeply(t, got, Map(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpMap(t *testing.T, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Map(model, expectedEntries), args...)
}

// CmpMapEach is a shortcut for:
//
//   CmpDeeply(t, got, MapEach(expectedValue), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpMapEach(t *testing.T, got interface{}, expectedValue interface{}, args ...interface{}) bool {
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
func CmpN(t *testing.T, got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, N(num, tolerance), args...)
}

// CmpNil is a shortcut for:
//
//   CmpDeeply(t, got, Nil(), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNil(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Nil(), args...)
}

// CmpNone is a shortcut for:
//
//   CmpDeeply(t, got, None(expectedValues...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNone(t *testing.T, got interface{}, expectedValues []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, None(expectedValues...), args...)
}

// CmpNoneOf is a shortcut for:
//
//   CmpDeeply(t, got, NoneOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNoneOf(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, NoneOf(expectedItems...), args...)
}

// CmpNot is a shortcut for:
//
//   CmpDeeply(t, got, Not(expected), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNot(t *testing.T, got interface{}, expected interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Not(expected), args...)
}

// CmpNotNil is a shortcut for:
//
//   CmpDeeply(t, got, NotNil(), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNotNil(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, NotNil(), args...)
}

// CmpPPtr is a shortcut for:
//
//   CmpDeeply(t, got, PPtr(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpPPtr(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, PPtr(val), args...)
}

// CmpPtr is a shortcut for:
//
//   CmpDeeply(t, got, Ptr(val), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpPtr(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
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
func CmpRe(t *testing.T, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Re(reg, capture), args...)
}

// CmpReAll is a shortcut for:
//
//   CmpDeeply(t, got, ReAll(reg, capture), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpReAll(t *testing.T, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, ReAll(reg, capture), args...)
}

// CmpSet is a shortcut for:
//
//   CmpDeeply(t, got, Set(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSet(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Set(expectedItems...), args...)
}

// CmpShallow is a shortcut for:
//
//   CmpDeeply(t, got, Shallow(expectedPtr), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpShallow(t *testing.T, got interface{}, expectedPtr interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Shallow(expectedPtr), args...)
}

// CmpSlice is a shortcut for:
//
//   CmpDeeply(t, got, Slice(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSlice(t *testing.T, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Slice(model, expectedEntries), args...)
}

// CmpString is a shortcut for:
//
//   CmpDeeply(t, got, String(expected), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpString(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, String(expected), args...)
}

// CmpStruct is a shortcut for:
//
//   CmpDeeply(t, got, Struct(model, expectedFields), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpStruct(t *testing.T, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool {
	return CmpDeeply(t, got, Struct(model, expectedFields), args...)
}

// CmpSubBagOf is a shortcut for:
//
//   CmpDeeply(t, got, SubBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubBagOf(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubBagOf(expectedItems...), args...)
}

// CmpSubMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SubMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubMapOf(t *testing.T, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SubMapOf(model, expectedEntries), args...)
}

// CmpSubSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SubSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubSetOf(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubSetOf(expectedItems...), args...)
}

// CmpSuperBagOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperBagOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperBagOf(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperBagOf(expectedItems...), args...)
}

// CmpSuperMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperMapOf(model, expectedEntries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperMapOf(t *testing.T, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperMapOf(model, expectedEntries), args...)
}

// CmpSuperSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperSetOf(expectedItems...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperSetOf(t *testing.T, got interface{}, expectedItems []interface{}, args ...interface{}) bool {
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
func CmpTruncTime(t *testing.T, got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool {
	return CmpDeeply(t, got, TruncTime(expectedTime, trunc), args...)
}
