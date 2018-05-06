package testdeep

// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!

import (
	"regexp"
	"testing"
	"time"
)

// CmpAll is a shortcut for:
//   CmpDeeply(t, got, All(items...), args...)
func CmpAll(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, All(items...), args...)
}

// CmpAny is a shortcut for:
//   CmpDeeply(t, got, Any(items...), args...)
func CmpAny(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Any(items...), args...)
}

// CmpArray is a shortcut for:
//   CmpDeeply(t, got, Array(model, entries), args...)
func CmpArray(t *testing.T, got interface{}, model interface{}, entries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Array(model, entries), args...)
}

// CmpArrayEach is a shortcut for:
//   CmpDeeply(t, got, ArrayEach(item), args...)
func CmpArrayEach(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, ArrayEach(item), args...)
}

// CmpBag is a shortcut for:
//   CmpDeeply(t, got, Bag(items...), args...)
func CmpBag(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Bag(items...), args...)
}

// CmpBetween is a shortcut for:
//   CmpDeeply(t, got, Between(from, to, bounds), args...)
//
// Between() optional parameter "bounds" is here mandatory.
// BoundsInIn value should be passed to mimic its absence in
// original Between() call.
func CmpBetween(t *testing.T, got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool {
	return CmpDeeply(t, got, Between(from, to, bounds), args...)
}

// CmpCap is a shortcut for:
//   CmpDeeply(t, got, Cap(min), args...)
func CmpCap(t *testing.T, got interface{}, min int, args ...interface{}) bool {
	return CmpDeeply(t, got, Cap(min), args...)
}

// CmpCapBetween is a shortcut for:
//   CmpDeeply(t, got, CapBetween(min, max), args...)
func CmpCapBetween(t *testing.T, got interface{}, min int, max int, args ...interface{}) bool {
	return CmpDeeply(t, got, CapBetween(min, max), args...)
}

// CmpCode is a shortcut for:
//   CmpDeeply(t, got, Code(fn), args...)
func CmpCode(t *testing.T, got interface{}, fn interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Code(fn), args...)
}

// CmpContains is a shortcut for:
//   CmpDeeply(t, got, Contains(expected), args...)
func CmpContains(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, Contains(expected), args...)
}

// CmpGt is a shortcut for:
//   CmpDeeply(t, got, Gt(val), args...)
func CmpGt(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Gt(val), args...)
}

// CmpGte is a shortcut for:
//   CmpDeeply(t, got, Gte(val), args...)
func CmpGte(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Gte(val), args...)
}

// CmpHasPrefix is a shortcut for:
//   CmpDeeply(t, got, HasPrefix(expected), args...)
func CmpHasPrefix(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, HasPrefix(expected), args...)
}

// CmpHasSuffix is a shortcut for:
//   CmpDeeply(t, got, HasSuffix(expected), args...)
func CmpHasSuffix(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, HasSuffix(expected), args...)
}

// CmpIsa is a shortcut for:
//   CmpDeeply(t, got, Isa(model), args...)
func CmpIsa(t *testing.T, got interface{}, model interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Isa(model), args...)
}

// CmpLen is a shortcut for:
//   CmpDeeply(t, got, Len(min), args...)
func CmpLen(t *testing.T, got interface{}, min int, args ...interface{}) bool {
	return CmpDeeply(t, got, Len(min), args...)
}

// CmpLenBetween is a shortcut for:
//   CmpDeeply(t, got, LenBetween(min, max), args...)
func CmpLenBetween(t *testing.T, got interface{}, min int, max int, args ...interface{}) bool {
	return CmpDeeply(t, got, LenBetween(min, max), args...)
}

// CmpLt is a shortcut for:
//   CmpDeeply(t, got, Lt(val), args...)
func CmpLt(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Lt(val), args...)
}

// CmpLte is a shortcut for:
//   CmpDeeply(t, got, Lte(val), args...)
func CmpLte(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Lte(val), args...)
}

// CmpMap is a shortcut for:
//   CmpDeeply(t, got, Map(model, entries), args...)
func CmpMap(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Map(model, entries), args...)
}

// CmpMapEach is a shortcut for:
//   CmpDeeply(t, got, MapEach(item), args...)
func CmpMapEach(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, MapEach(item), args...)
}

// CmpN is a shortcut for:
//   CmpDeeply(t, got, N(num, tolerance), args...)
//
// N() optional parameter "tolerance" is here mandatory.
// 0 value should be passed to mimic its absence in
// original N() call.
func CmpN(t *testing.T, got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, N(num, tolerance), args...)
}

// CmpNil is a shortcut for:
//   CmpDeeply(t, got, Nil(), args...)
func CmpNil(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Nil(), args...)
}

// CmpNone is a shortcut for:
//   CmpDeeply(t, got, None(items...), args...)
func CmpNone(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, None(items...), args...)
}

// CmpNoneOf is a shortcut for:
//   CmpDeeply(t, got, NoneOf(items...), args...)
func CmpNoneOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, NoneOf(items...), args...)
}

// CmpNot is a shortcut for:
//   CmpDeeply(t, got, Not(item), args...)
func CmpNot(t *testing.T, got interface{}, item interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Not(item), args...)
}

// CmpNotNil is a shortcut for:
//   CmpDeeply(t, got, NotNil(), args...)
func CmpNotNil(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, NotNil(), args...)
}

// CmpPPtr is a shortcut for:
//   CmpDeeply(t, got, PPtr(val), args...)
func CmpPPtr(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, PPtr(val), args...)
}

// CmpPtr is a shortcut for:
//   CmpDeeply(t, got, Ptr(val), args...)
func CmpPtr(t *testing.T, got interface{}, val interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Ptr(val), args...)
}

// CmpRe is a shortcut for:
//   CmpDeeply(t, got, Re(reg, opts...), args...)
func CmpRe(t *testing.T, got interface{}, reg string, opts []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Re(reg, opts...), args...)
}

// CmpRex is a shortcut for:
//   CmpDeeply(t, got, Rex(re, opts...), args...)
func CmpRex(t *testing.T, got interface{}, re *regexp.Regexp, opts []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Rex(re, opts...), args...)
}

// CmpSet is a shortcut for:
//   CmpDeeply(t, got, Set(items...), args...)
func CmpSet(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Set(items...), args...)
}

// CmpShallow is a shortcut for:
//   CmpDeeply(t, got, Shallow(ptr), args...)
func CmpShallow(t *testing.T, got interface{}, ptr interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, Shallow(ptr), args...)
}

// CmpSlice is a shortcut for:
//   CmpDeeply(t, got, Slice(model, entries), args...)
func CmpSlice(t *testing.T, got interface{}, model interface{}, entries ArrayEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, Slice(model, entries), args...)
}

// CmpString is a shortcut for:
//   CmpDeeply(t, got, String(expected), args...)
func CmpString(t *testing.T, got interface{}, expected string, args ...interface{}) bool {
	return CmpDeeply(t, got, String(expected), args...)
}

// CmpStruct is a shortcut for:
//   CmpDeeply(t, got, Struct(model, expectedFields), args...)
func CmpStruct(t *testing.T, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool {
	return CmpDeeply(t, got, Struct(model, expectedFields), args...)
}

// CmpSubBagOf is a shortcut for:
//   CmpDeeply(t, got, SubBagOf(items...), args...)
func CmpSubBagOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubBagOf(items...), args...)
}

// CmpSubMapOf is a shortcut for:
//   CmpDeeply(t, got, SubMapOf(model, entries), args...)
func CmpSubMapOf(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SubMapOf(model, entries), args...)
}

// CmpSubSetOf is a shortcut for:
//   CmpDeeply(t, got, SubSetOf(items...), args...)
func CmpSubSetOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SubSetOf(items...), args...)
}

// CmpSuperBagOf is a shortcut for:
//   CmpDeeply(t, got, SuperBagOf(items...), args...)
func CmpSuperBagOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperBagOf(items...), args...)
}

// CmpSuperMapOf is a shortcut for:
//   CmpDeeply(t, got, SuperMapOf(model, entries), args...)
func CmpSuperMapOf(t *testing.T, got interface{}, model interface{}, entries MapEntries, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperMapOf(model, entries), args...)
}

// CmpSuperSetOf is a shortcut for:
//   CmpDeeply(t, got, SuperSetOf(items...), args...)
func CmpSuperSetOf(t *testing.T, got interface{}, items []interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, SuperSetOf(items...), args...)
}

// CmpTruncTime is a shortcut for:
//   CmpDeeply(t, got, TruncTime(val, trunc), args...)
//
// TruncTime() optional parameter "trunc" is here mandatory.
// 0 value should be passed to mimic its absence in
// original TruncTime() call.
func CmpTruncTime(t *testing.T, got interface{}, val interface{}, trunc time.Duration, args ...interface{}) bool {
	return CmpDeeply(t, got, TruncTime(val, trunc), args...)
}
