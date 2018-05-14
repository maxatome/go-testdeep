package testdeep

// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

import (
	"regexp"
	"testing"
	"time"
)

// CmpAll is a shortcut for:
//
//   CmpDeeply(t, got, All(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpAll(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, All(items...), args...)
}

// CmpAny is a shortcut for:
//
//   CmpDeeply(t, got, Any(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpAny(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Any(items...), args...)
}

// CmpArray is a shortcut for:
//
//   CmpDeeply(t, got, Array(model, entries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpArray(t *testing.T, got interface{}, model interface{}, entries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Array(model, entries), args...)
}

// CmpArrayEach is a shortcut for:
//
//   CmpDeeply(t, got, ArrayEach(item), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpArrayEach(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, ArrayEach(item), args...)
}

// CmpBag is a shortcut for:
//
//   CmpDeeply(t, got, Bag(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpBag(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Bag(items...), args...)
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
//   CmpDeeply(t, got, Map(model, entries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpMap(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Map(model, entries), args...)
}

// CmpMapEach is a shortcut for:
//
//   CmpDeeply(t, got, MapEach(item), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpMapEach(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, MapEach(item), args...)
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
//   CmpDeeply(t, got, None(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNone(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, None(items...), args...)
}

// CmpNoneOf is a shortcut for:
//
//   CmpDeeply(t, got, NoneOf(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNoneOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, NoneOf(items...), args...)
}

// CmpNot is a shortcut for:
//
//   CmpDeeply(t, got, Not(item), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpNot(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Not(item), args...)
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
//   CmpDeeply(t, got, Re(reg, opts...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpRe(t *testing.T, got interface{}, reg string, opts []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Re(reg, opts...), args...)
}

// CmpRex is a shortcut for:
//
//   CmpDeeply(t, got, Rex(re, opts...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpRex(t *testing.T, got interface{}, re *regexp.Regexp, opts []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Rex(re, opts...), args...)
}

// CmpSet is a shortcut for:
//
//   CmpDeeply(t, got, Set(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSet(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Set(items...), args...)
}

// CmpShallow is a shortcut for:
//
//   CmpDeeply(t, got, Shallow(ptr), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpShallow(t *testing.T, got interface{}, ptr interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Shallow(ptr), args...)
}

// CmpSlice is a shortcut for:
//
//   CmpDeeply(t, got, Slice(model, entries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSlice(t *testing.T, got interface{}, model interface{}, entries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Slice(model, entries), args...)
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
//   CmpDeeply(t, got, SubBagOf(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubBagOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubBagOf(items...), args...)
}

// CmpSubMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SubMapOf(model, entries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubMapOf(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SubMapOf(model, entries), args...)
}

// CmpSubSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SubSetOf(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSubSetOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubSetOf(items...), args...)
}

// CmpSuperBagOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperBagOf(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperBagOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperBagOf(items...), args...)
}

// CmpSuperMapOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperMapOf(model, entries), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperMapOf(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperMapOf(model, entries), args...)
}

// CmpSuperSetOf is a shortcut for:
//
//   CmpDeeply(t, got, SuperSetOf(items...), args...)
//
// Returns true if the test is OK, false if it fails.
func CmpSuperSetOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperSetOf(items...), args...)
}

// CmpTruncTime is a shortcut for:
//
//   CmpDeeply(t, got, TruncTime(val, trunc), args...)
//
// TruncTime() optional parameter "trunc" is here mandatory.
// 0 value should be passed to mimic its absence in
// original TruncTime() call.
//
// Returns true if the test is OK, false if it fails.
func CmpTruncTime(t *testing.T, got interface{}, val interface{}, trunc time.Duration, args ...interface{}) bool {
	return CmpDeeply(t, got, TruncTime(val, trunc), args...)
}
