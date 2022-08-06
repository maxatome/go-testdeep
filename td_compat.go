// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"github.com/maxatome/go-testdeep/td"
)

// TestingT is a deprecated alias of [td.TestingT].
type TestingT = td.TestingT

// TestingFT is a deprecated alias of [td.TestingFT], which is itself
// a deprecated alias of [testing.TB].
type TestingFT = td.TestingFT

// TestDeep is a deprecated alias of [td.TestDeep].
type TestDeep = td.TestDeep

// ContextConfig is a deprecated alias of [td.ContextConfig].
type ContextConfig = td.ContextConfig

// T is a deprecated alias of [td.T].
type T = td.T

// ArrayEntries is a deprecated alias of [td.ArrayEntries].
type ArrayEntries = td.ArrayEntries

// BoundsKind is a deprecated alias of [td.BoundsKind].
type BoundsKind = td.BoundsKind

// MapEntries is a deprecated alias of [td.MapEntries].
type MapEntries = td.MapEntries

// SmuggledGot is a deprecated alias of [td.SmuggledGot].
type SmuggledGot = td.SmuggledGot

// StructFields is a deprecated alias of [td.StructFields].
type StructFields = td.StructFields

// Incompatible change: testdeep.DefaultContextConfig must be replaced
// by td.DefaultContextConfig. Commented here to raise an error if used.
// var DefaultContextConfig = td.DefaultContextConfig

// Cmp is a deprecated alias of [td.Cmp].
var Cmp = td.Cmp

// CmpDeeply is a deprecated alias of [td.CmpDeeply].
var CmpDeeply = td.CmpDeeply

// CmpTrue is a deprecated alias of [td.CmpTrue].
var CmpTrue = td.CmpTrue

// CmpFalse is a deprecated alias of [td.CmpFalse].
var CmpFalse = td.CmpFalse

// CmpError is a deprecated alias of [td.CmpError].
var CmpError = td.CmpError

// CmpNoError is a deprecated alias of [td.CmpNoError].
var CmpNoError = td.CmpNoError

// CmpPanic is a deprecated alias of [td.CmpPanic].
var CmpPanic = td.CmpPanic

// CmpNotPanic is a deprecated alias of [td.CmpNotPanic].
var CmpNotPanic = td.CmpNotPanic

// EqDeeply is a deprecated alias of [td.EqDeeply].
var EqDeeply = td.EqDeeply

// EqDeeplyError is a deprecated alias of [td.EqDeeplyError].
var EqDeeplyError = td.EqDeeplyError

// AddAnchorableStructType is a deprecated alias of [td.AddAnchorableStructType].
var AddAnchorableStructType = td.AddAnchorableStructType

// NewT is a deprecated alias of [td.NewT].
var NewT = td.NewT

// Assert is a deprecated alias of [td.Assert].
var Assert = td.Assert

// Require is a deprecated alias of [td.Require].
var Require = td.Require

// AssertRequire is a deprecated alias of [td.AssertRequire].
var AssertRequire = td.AssertRequire

// CmpAll is a deprecated alias of [td.CmpAll].
var CmpAll = td.CmpAll

// CmpAny is a deprecated alias of [td.CmpAny].
var CmpAny = td.CmpAny

// CmpArray is a deprecated alias of [td.CmpArray].
var CmpArray = td.CmpArray

// CmpArrayEach is a deprecated alias of [td.CmpArrayEach].
var CmpArrayEach = td.CmpArrayEach

// CmpBag is a deprecated alias of [td.CmpBag].
var CmpBag = td.CmpBag

// CmpBetween is a deprecated alias of [td.CmpBetween].
var CmpBetween = td.CmpBetween

// CmpCap is a deprecated alias of [td.CmpCap].
var CmpCap = td.CmpCap

// CmpCode is a deprecated alias of [td.CmpCode].
var CmpCode = td.CmpCode

// CmpContains is a deprecated alias of [td.CmpContains].
var CmpContains = td.CmpContains

// CmpContainsKey is a deprecated alias of [td.CmpContainsKey].
var CmpContainsKey = td.CmpContainsKey

// CmpEmpty is a deprecated alias of [td.CmpEmpty].
var CmpEmpty = td.CmpEmpty

// CmpGt is a deprecated alias of [td.CmpGt].
var CmpGt = td.CmpGt

// CmpGte is a deprecated alias of [td.CmpGte].
var CmpGte = td.CmpGte

// CmpHasPrefix is a deprecated alias of [td.CmpHasPrefix].
var CmpHasPrefix = td.CmpHasPrefix

// CmpHasSuffix is a deprecated alias of [td.CmpHasSuffix].
var CmpHasSuffix = td.CmpHasSuffix

// CmpIsa is a deprecated alias of [td.CmpIsa].
var CmpIsa = td.CmpIsa

// CmpJSON is a deprecated alias of [td.CmpJSON].
var CmpJSON = td.CmpJSON

// CmpKeys is a deprecated alias of [td.CmpKeys].
var CmpKeys = td.CmpKeys

// CmpLax is a deprecated alias of [td.CmpLax].
var CmpLax = td.CmpLax

// CmpLen is a deprecated alias of [td.CmpLen].
var CmpLen = td.CmpLen

// CmpLt is a deprecated alias of [td.CmpLt].
var CmpLt = td.CmpLt

// CmpLte is a deprecated alias of [td.CmpLte].
var CmpLte = td.CmpLte

// CmpMap is a deprecated alias of [td.CmpMap].
var CmpMap = td.CmpMap

// CmpMapEach is a deprecated alias of [td.CmpMapEach].
var CmpMapEach = td.CmpMapEach

// CmpN is a deprecated alias of [td.CmpN].
var CmpN = td.CmpN

// CmpNaN is a deprecated alias of [td.CmpNaN].
var CmpNaN = td.CmpNaN

// CmpNil is a deprecated alias of [td.CmpNil].
var CmpNil = td.CmpNil

// CmpNone is a deprecated alias of [td.CmpNone].
var CmpNone = td.CmpNone

// CmpNot is a deprecated alias of [td.CmpNot].
var CmpNot = td.CmpNot

// CmpNotAny is a deprecated alias of [td.CmpNotAny].
var CmpNotAny = td.CmpNotAny

// CmpNotEmpty is a deprecated alias of [td.CmpNotEmpty].
var CmpNotEmpty = td.CmpNotEmpty

// CmpNotNaN is a deprecated alias of [td.CmpNotNaN].
var CmpNotNaN = td.CmpNotNaN

// CmpNotNil is a deprecated alias of [td.CmpNotNil].
var CmpNotNil = td.CmpNotNil

// CmpNotZero is a deprecated alias of [td.CmpNotZero].
var CmpNotZero = td.CmpNotZero

// CmpPPtr is a deprecated alias of [td.CmpPPtr].
var CmpPPtr = td.CmpPPtr

// CmpPtr is a deprecated alias of [td.CmpPtr].
var CmpPtr = td.CmpPtr

// CmpRe is a deprecated alias of [td.CmpRe].
var CmpRe = td.CmpRe

// CmpReAll is a deprecated alias of [td.CmpReAll].
var CmpReAll = td.CmpReAll

// CmpSet is a deprecated alias of [td.CmpSet].
var CmpSet = td.CmpSet

// CmpShallow is a deprecated alias of [td.CmpShallow].
var CmpShallow = td.CmpShallow

// CmpSlice is a deprecated alias of [td.CmpSlice].
var CmpSlice = td.CmpSlice

// CmpSmuggle is a deprecated alias of [td.CmpSmuggle].
var CmpSmuggle = td.CmpSmuggle

// CmpSStruct is a deprecated alias of [td.CmpSStruct].
var CmpSStruct = td.CmpSStruct

// CmpString is a deprecated alias of [td.CmpString].
var CmpString = td.CmpString

// CmpStruct is a deprecated alias of [td.CmpStruct].
var CmpStruct = td.CmpStruct

// CmpSubBagOf is a deprecated alias of [td.CmpSubBagOf].
var CmpSubBagOf = td.CmpSubBagOf

// CmpSubJSONOf is a deprecated alias of [td.CmpSubJSONOf].
var CmpSubJSONOf = td.CmpSubJSONOf

// CmpSubMapOf is a deprecated alias of [td.CmpSubMapOf].
var CmpSubMapOf = td.CmpSubMapOf

// CmpSubSetOf is a deprecated alias of [td.CmpSubSetOf].
var CmpSubSetOf = td.CmpSubSetOf

// CmpSuperBagOf is a deprecated alias of [td.CmpSuperBagOf].
var CmpSuperBagOf = td.CmpSuperBagOf

// CmpSuperJSONOf is a deprecated alias of [td.CmpSuperJSONOf].
var CmpSuperJSONOf = td.CmpSuperJSONOf

// CmpSuperMapOf is a deprecated alias of [td.CmpSuperMapOf].
var CmpSuperMapOf = td.CmpSuperMapOf

// CmpSuperSetOf is a deprecated alias of [td.CmpSuperSetOf].
var CmpSuperSetOf = td.CmpSuperSetOf

// CmpTruncTime is a deprecated alias of [td.CmpTruncTime].
var CmpTruncTime = td.CmpTruncTime

// CmpValues is a deprecated alias of [td.CmpValues].
var CmpValues = td.CmpValues

// CmpZero is a deprecated alias of [td.CmpZero].
var CmpZero = td.CmpZero

// All is a deprecated alias of [td.All].
var All = td.All

// Any is a deprecated alias of [td.Any].
var Any = td.Any

// Array is a deprecated alias of [td.Array].
var Array = td.Array

// ArrayEach is a deprecated alias of [td.ArrayEach].
var ArrayEach = td.ArrayEach

// Bag is a deprecated alias of [td.Bag].
var Bag = td.Bag

// Between is a deprecated alias of [td.Between].
var Between = td.Between

// Cap is a deprecated alias of [td.Cap].
var Cap = td.Cap

// Catch is a deprecated alias of [td.Catch].
var Catch = td.Catch

// Code is a deprecated alias of [td.Code].
var Code = td.Code

// Contains is a deprecated alias of [td.Contains].
var Contains = td.Contains

// ContainsKey is a deprecated alias of [td.ContainsKey].
var ContainsKey = td.ContainsKey

// Delay is a deprecated alias of [td.ContainsKey].
var Delay = td.Delay

// Empty is a deprecated alias of [td.Empty].
var Empty = td.Empty

// Gt is a deprecated alias of [td.Gt].
var Gt = td.Gt

// Gte is a deprecated alias of [td.Gte].
var Gte = td.Gte

// HasPrefix is a deprecated alias of [td.HasPrefix].
var HasPrefix = td.HasPrefix

// HasSuffix is a deprecated alias of [td.HasSuffix].
var HasSuffix = td.HasSuffix

// Ignore is a deprecated alias of [td.Ignore].
var Ignore = td.Ignore

// Isa is a deprecated alias of [td.Isa].
var Isa = td.Isa

// JSON is a deprecated alias of [td.JSON].
var JSON = td.JSON

// Keys is a deprecated alias of [td.Keys].
var Keys = td.Keys

// Lax is a deprecated alias of [td.Lax].
var Lax = td.Lax

// Len is a deprecated alias of [td.Len].
var Len = td.Len

// Lt is a deprecated alias of [td.Lt].
var Lt = td.Lt

// Lte is a deprecated alias of [td.Lte].
var Lte = td.Lte

// Map is a deprecated alias of [td.Map].
var Map = td.Map

// MapEach is a deprecated alias of [td.MapEach].
var MapEach = td.MapEach

// N is a deprecated alias of [td.N].
var N = td.N

// NaN is a deprecated alias of [td.NaN].
var NaN = td.NaN

// Nil is a deprecated alias of [td.Nil].
var Nil = td.Nil

// None is a deprecated alias of [td.None].
var None = td.None

// Not is a deprecated alias of [td.Not].
var Not = td.Not

// NotAny is a deprecated alias of [td.NotAny].
var NotAny = td.NotAny

// NotEmpty is a deprecated alias of [td.NotEmpty].
var NotEmpty = td.NotEmpty

// NotNaN is a deprecated alias of [td.NotNaN].
var NotNaN = td.NotNaN

// NotNil is a deprecated alias of [td.NotNil].
var NotNil = td.NotNil

// NotZero is a deprecated alias of [td.NotZero].
var NotZero = td.NotZero

// Ptr is a deprecated alias of [td.Ptr].
var Ptr = td.Ptr

// PPtr is a deprecated alias of [td.PPtr].
var PPtr = td.PPtr

// Re is a deprecated alias of [td.Re].
var Re = td.Re

// ReAll is a deprecated alias of [td.ReAll].
var ReAll = td.ReAll

// Set is a deprecated alias of [td.Set].
var Set = td.Set

// Shallow is a deprecated alias of [td.Shallow].
var Shallow = td.Shallow

// Slice is a deprecated alias of [td.Slice].
var Slice = td.Slice

// Smuggle is a deprecated alias of [td.Smuggle].
var Smuggle = td.Smuggle

// String is a deprecated alias of [td.String].
var String = td.String

// SStruct is a deprecated alias of [td.SStruct].
var SStruct = td.SStruct

// Struct is a deprecated alias of [td.Struct].
var Struct = td.Struct

// SubBagOf is a deprecated alias of [td.SubBagOf].
var SubBagOf = td.SubBagOf

// SubJSONOf is a deprecated alias of [td.SubJSONOf].
var SubJSONOf = td.SubJSONOf

// SubMapOf is a deprecated alias of [td.SubMapOf].
var SubMapOf = td.SubMapOf

// SubSetOf is a deprecated alias of [td.SubSetOf].
var SubSetOf = td.SubSetOf

// SuperBagOf is a deprecated alias of [td.SuperBagOf].
var SuperBagOf = td.SuperBagOf

// SuperJSONOf is a deprecated alias of [td.SuperJSONOf].
var SuperJSONOf = td.SuperJSONOf

// SuperMapOf is a deprecated alias of [td.SuperMapOf].
var SuperMapOf = td.SuperMapOf

// SuperSetOf is a deprecated alias of [td.SuperSetOf].
var SuperSetOf = td.SuperSetOf

// Tag is a deprecated alias of [td.Tag].
var Tag = td.Tag

// TruncTime is a deprecated alias of [td.TruncTime].
var TruncTime = td.TruncTime

// Values is a deprecated alias of [td.Values].
var Values = td.Values

// Zero is a deprecated alias of [td.Zero].
var Zero = td.Zero

// BoundsInIn is a deprecated alias of [td.BoundsInIn].
const BoundsInIn = td.BoundsInIn

// BoundsInOut is a deprecated alias of [td.BoundsInOut].
const BoundsInOut = td.BoundsInOut

// BoundsOutIn is a deprecated alias of [td.BoundsOutIn].
const BoundsOutIn = td.BoundsOutIn

// BoundsOutOut is a deprecated alias of [td.BoundsOutOut].
const BoundsOutOut = td.BoundsOutOut
