// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"github.com/maxatome/go-testdeep/td"
)

// TestingT is a deprecated alias of github.com/maxatome/go-testdeep/td.TestingT
type TestingT = td.TestingT

// TestingFT is a deprecated alias of github.com/maxatome/go-testdeep/td.TestingFT.
type TestingFT = td.TestingFT

// TestDeep is a deprecated alias of github.com/maxatome/go-testdeep/td.TestDeep.
type TestDeep = td.TestDeep

// ContextConfig is a deprecated alias of github.com/maxatome/go-testdeep/td.ContextConfig.
type ContextConfig = td.ContextConfig

// T is a deprecated alias of github.com/maxatome/go-testdeep/td.T.
type T = td.T

// ArrayEntries is a deprecated alias of github.com/maxatome/go-testdeep/td.ArrayEntries.
type ArrayEntries = td.ArrayEntries

// BoundsKind is a deprecated alias of github.com/maxatome/go-testdeep/td.BoundsKind.
type BoundsKind = td.BoundsKind

// MapEntries is a deprecated alias of github.com/maxatome/go-testdeep/td.MapEntries.
type MapEntries = td.MapEntries

// SmuggledGot is a deprecated alias of github.com/maxatome/go-testdeep/td.SmuggledGot.
type SmuggledGot = td.SmuggledGot

// StructFields is a deprecated alias of github.com/maxatome/go-testdeep/td.StructFields.
type StructFields = td.StructFields

// Incompatible change: testdeep.DefaultContextConfig must be replaced
// by td.DefaultContextConfig. Commented here to raise an error if used.
// var DefaultContextConfig = td.DefaultContextConfig

var (
	// Cmp is a deprecated alias of github.com/maxatome/go-testdeep/td.Cmp.
	Cmp = td.Cmp
	// CmpDeeply is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpDeeply.
	CmpDeeply = td.CmpDeeply

	// CmpTrue is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpTrue.
	CmpTrue = td.CmpTrue
	// CmpFalse is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpFalse.
	CmpFalse = td.CmpFalse
	// CmpError is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpError.
	CmpError = td.CmpError
	// CmpNoError is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNoError.
	CmpNoError = td.CmpNoError
	// CmpPanic is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpPanic.
	CmpPanic = td.CmpPanic
	// CmpNotPanic is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotPanic.
	CmpNotPanic = td.CmpNotPanic

	// EqDeeply is a deprecated alias of github.com/maxatome/go-testdeep/td.EqDeeply.
	EqDeeply = td.EqDeeply
	// EqDeeplyError is a deprecated alias of github.com/maxatome/go-testdeep/td.EqDeeplyError.
	EqDeeplyError = td.EqDeeplyError

	// AddAnchorableStructType is a deprecated alias of github.com/maxatome/go-testdeep/td.AddAnchorableStructType.
	AddAnchorableStructType = td.AddAnchorableStructType

	// NewT is a deprecated alias of github.com/maxatome/go-testdeep/td.NewT.
	NewT = td.NewT
	// Assert is a deprecated alias of github.com/maxatome/go-testdeep/td.Assert.
	Assert = td.Assert
	// Require is a deprecated alias of github.com/maxatome/go-testdeep/td.Require.
	Require = td.Require
	// AssertRequire is a deprecated alias of github.com/maxatome/go-testdeep/td.AssertRequire.
	AssertRequire = td.AssertRequire

	// CmpAll is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpAll.
	CmpAll = td.CmpAll
	// CmpAny is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpAny.
	CmpAny = td.CmpAny
	// CmpArray is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpArray.
	CmpArray = td.CmpArray
	// CmpArrayEach is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpArrayEach.
	CmpArrayEach = td.CmpArrayEach
	// CmpBag is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpBag.
	CmpBag = td.CmpBag
	// CmpBetween is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpBetween.
	CmpBetween = td.CmpBetween
	// CmpCap is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpCap.
	CmpCap = td.CmpCap
	// CmpCode is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpCode.
	CmpCode = td.CmpCode
	// CmpContains is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpContains.
	CmpContains = td.CmpContains
	// CmpContainsKey is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpContainsKey.
	CmpContainsKey = td.CmpContainsKey
	// CmpEmpty is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpEmpty.
	CmpEmpty = td.CmpEmpty
	// CmpGt is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpGt.
	CmpGt = td.CmpGt
	// CmpGte is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpGte.
	CmpGte = td.CmpGte
	// CmpHasPrefix is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpHasPrefix.
	CmpHasPrefix = td.CmpHasPrefix
	// CmpHasSuffix is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpHasSuffix.
	CmpHasSuffix = td.CmpHasSuffix
	// CmpIsa is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpIsa.
	CmpIsa = td.CmpIsa
	// CmpJSON is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpJSON.
	CmpJSON = td.CmpJSON
	// CmpKeys is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpKeys.
	CmpKeys = td.CmpKeys
	// CmpLax is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpLax.
	CmpLax = td.CmpLax
	// CmpLen is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpLen.
	CmpLen = td.CmpLen
	// CmpLt is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpLt.
	CmpLt = td.CmpLt
	// CmpLte is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpLte.
	CmpLte = td.CmpLte
	// CmpMap is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpMap.
	CmpMap = td.CmpMap
	// CmpMapEach is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpMapEach.
	CmpMapEach = td.CmpMapEach
	// CmpN is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpN.
	CmpN = td.CmpN
	// CmpNaN is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNaN.
	CmpNaN = td.CmpNaN
	// CmpNil is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNil.
	CmpNil = td.CmpNil
	// CmpNone is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNone.
	CmpNone = td.CmpNone
	// CmpNot is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNot.
	CmpNot = td.CmpNot
	// CmpNotAny is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotAny.
	CmpNotAny = td.CmpNotAny
	// CmpNotEmpty is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotEmpty.
	CmpNotEmpty = td.CmpNotEmpty
	// CmpNotNaN is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotNaN.
	CmpNotNaN = td.CmpNotNaN
	// CmpNotNil is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotNil.
	CmpNotNil = td.CmpNotNil
	// CmpNotZero is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpNotZero.
	CmpNotZero = td.CmpNotZero
	// CmpPPtr is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpPPtr.
	CmpPPtr = td.CmpPPtr
	// CmpPtr is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpPtr.
	CmpPtr = td.CmpPtr
	// CmpRe is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpRe.
	CmpRe = td.CmpRe
	// CmpReAll is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpReAll.
	CmpReAll = td.CmpReAll
	// CmpSet is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSet.
	CmpSet = td.CmpSet
	// CmpShallow is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpShallow.
	CmpShallow = td.CmpShallow
	// CmpSlice is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSlice.
	CmpSlice = td.CmpSlice
	// CmpSmuggle is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSmuggle.
	CmpSmuggle = td.CmpSmuggle
	// CmpSStruct is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSStruct.
	CmpSStruct = td.CmpSStruct
	// CmpString is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpString.
	CmpString = td.CmpString
	// CmpStruct is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpStruct.
	CmpStruct = td.CmpStruct
	// CmpSubBagOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSubBagOf.
	CmpSubBagOf = td.CmpSubBagOf
	// CmpSubJSONOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSubJSONOf.
	CmpSubJSONOf = td.CmpSubJSONOf
	// CmpSubMapOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSubMapOf.
	CmpSubMapOf = td.CmpSubMapOf
	// CmpSubSetOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSubSetOf.
	CmpSubSetOf = td.CmpSubSetOf
	// CmpSuperBagOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSuperBagOf.
	CmpSuperBagOf = td.CmpSuperBagOf
	// CmpSuperJSONOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSuperJSONOf.
	CmpSuperJSONOf = td.CmpSuperJSONOf
	// CmpSuperMapOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSuperMapOf.
	CmpSuperMapOf = td.CmpSuperMapOf
	// CmpSuperSetOf is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpSuperSetOf.
	CmpSuperSetOf = td.CmpSuperSetOf
	// CmpTruncTime is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpTruncTime.
	CmpTruncTime = td.CmpTruncTime
	// CmpValues is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpValues.
	CmpValues = td.CmpValues
	// CmpZero is a deprecated alias of github.com/maxatome/go-testdeep/td.CmpZero.
	CmpZero = td.CmpZero

	// All is a deprecated alias of github.com/maxatome/go-testdeep/td.All.
	All = td.All
	// Any is a deprecated alias of github.com/maxatome/go-testdeep/td.Any.
	Any = td.Any
	// Array is a deprecated alias of github.com/maxatome/go-testdeep/td.Array.
	Array = td.Array
	// ArrayEach is a deprecated alias of github.com/maxatome/go-testdeep/td.ArrayEach.
	ArrayEach = td.ArrayEach
	// Bag is a deprecated alias of github.com/maxatome/go-testdeep/td.Bag.
	Bag = td.Bag
	// Between is a deprecated alias of github.com/maxatome/go-testdeep/td.Between.
	Between = td.Between
	// Cap is a deprecated alias of github.com/maxatome/go-testdeep/td.Cap.
	Cap = td.Cap
	// Catch is a deprecated alias of github.com/maxatome/go-testdeep/td.Catch.
	Catch = td.Catch
	// Code is a deprecated alias of github.com/maxatome/go-testdeep/td.Code.
	Code = td.Code
	// Contains is a deprecated alias of github.com/maxatome/go-testdeep/td.Contains.
	Contains = td.Contains
	// ContainsKey is a deprecated alias of github.com/maxatome/go-testdeep/td.ContainsKey.
	ContainsKey = td.ContainsKey
	// Delay is a deprecated alias of github.com/maxatome/go-testdeep/td.ContainsKey.
	Delay = td.Delay
	// Empty is a deprecated alias of github.com/maxatome/go-testdeep/td.Empty.
	Empty = td.Empty
	// Gt is a deprecated alias of github.com/maxatome/go-testdeep/td.Gt.
	Gt = td.Gt
	// Gte is a deprecated alias of github.com/maxatome/go-testdeep/td.Gte.
	Gte = td.Gte
	// HasPrefix is a deprecated alias of github.com/maxatome/go-testdeep/td.HasPrefix.
	HasPrefix = td.HasPrefix
	// HasSuffix is a deprecated alias of github.com/maxatome/go-testdeep/td.HasSuffix.
	HasSuffix = td.HasSuffix
	// Ignore is a deprecated alias of github.com/maxatome/go-testdeep/td.Ignore.
	Ignore = td.Ignore
	// Isa is a deprecated alias of github.com/maxatome/go-testdeep/td.Isa.
	Isa = td.Isa
	// JSON is a deprecated alias of github.com/maxatome/go-testdeep/td.JSON.
	JSON = td.JSON
	// Keys is a deprecated alias of github.com/maxatome/go-testdeep/td.Keys.
	Keys = td.Keys
	// Lax is a deprecated alias of github.com/maxatome/go-testdeep/td.Lax.
	Lax = td.Lax
	// Len is a deprecated alias of github.com/maxatome/go-testdeep/td.Len.
	Len = td.Len
	// Lt is a deprecated alias of github.com/maxatome/go-testdeep/td.Lt.
	Lt = td.Lt
	// Lte is a deprecated alias of github.com/maxatome/go-testdeep/td.Lte.
	Lte = td.Lte
	// Map is a deprecated alias of github.com/maxatome/go-testdeep/td.Map.
	Map = td.Map
	// MapEach is a deprecated alias of github.com/maxatome/go-testdeep/td.MapEach.
	MapEach = td.MapEach
	// N is a deprecated alias of github.com/maxatome/go-testdeep/td.N.
	N = td.N
	// NaN is a deprecated alias of github.com/maxatome/go-testdeep/td.NaN.
	NaN = td.NaN
	// Nil is a deprecated alias of github.com/maxatome/go-testdeep/td.Nil.
	Nil = td.Nil
	// None is a deprecated alias of github.com/maxatome/go-testdeep/td.None.
	None = td.None
	// Not is a deprecated alias of github.com/maxatome/go-testdeep/td.Not.
	Not = td.Not
	// NotAny is a deprecated alias of github.com/maxatome/go-testdeep/td.NotAny.
	NotAny = td.NotAny
	// NotEmpty is a deprecated alias of github.com/maxatome/go-testdeep/td.NotEmpty.
	NotEmpty = td.NotEmpty
	// NotNaN is a deprecated alias of github.com/maxatome/go-testdeep/td.NotNaN.
	NotNaN = td.NotNaN
	// NotNil is a deprecated alias of github.com/maxatome/go-testdeep/td.NotNil.
	NotNil = td.NotNil
	// NotZero is a deprecated alias of github.com/maxatome/go-testdeep/td.NotZero.
	NotZero = td.NotZero
	// Ptr is a deprecated alias of github.com/maxatome/go-testdeep/td.Ptr.
	Ptr = td.Ptr
	// PPtr is a deprecated alias of github.com/maxatome/go-testdeep/td.PPtr.
	PPtr = td.PPtr
	// Re is a deprecated alias of github.com/maxatome/go-testdeep/td.Re.
	Re = td.Re
	// ReAll is a deprecated alias of github.com/maxatome/go-testdeep/td.ReAll.
	ReAll = td.ReAll
	// Set is a deprecated alias of github.com/maxatome/go-testdeep/td.Set.
	Set = td.Set
	// Shallow is a deprecated alias of github.com/maxatome/go-testdeep/td.Shallow.
	Shallow = td.Shallow
	// Slice is a deprecated alias of github.com/maxatome/go-testdeep/td.Slice.
	Slice = td.Slice
	// Smuggle is a deprecated alias of github.com/maxatome/go-testdeep/td.Smuggle.
	Smuggle = td.Smuggle
	// String is a deprecated alias of github.com/maxatome/go-testdeep/td.String.
	String = td.String
	// SStruct is a deprecated alias of github.com/maxatome/go-testdeep/td.SStruct.
	SStruct = td.SStruct
	// Struct is a deprecated alias of github.com/maxatome/go-testdeep/td.Struct.
	Struct = td.Struct
	// SubBagOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SubBagOf.
	SubBagOf = td.SubBagOf
	// SubJSONOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SubJSONOf.
	SubJSONOf = td.SubJSONOf
	// SubMapOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SubMapOf.
	SubMapOf = td.SubMapOf
	// SubSetOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SubSetOf.
	SubSetOf = td.SubSetOf
	// SuperBagOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SuperBagOf.
	SuperBagOf = td.SuperBagOf
	// SuperJSONOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SuperJSONOf.
	SuperJSONOf = td.SuperJSONOf
	// SuperMapOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SuperMapOf.
	SuperMapOf = td.SuperMapOf
	// SuperSetOf is a deprecated alias of github.com/maxatome/go-testdeep/td.SuperSetOf.
	SuperSetOf = td.SuperSetOf
	// Tag is a deprecated alias of github.com/maxatome/go-testdeep/td.Tag.
	Tag = td.Tag
	// TruncTime is a deprecated alias of github.com/maxatome/go-testdeep/td.TruncTime.
	TruncTime = td.TruncTime
	// Values is a deprecated alias of github.com/maxatome/go-testdeep/td.Values.
	Values = td.Values
	// Zero is a deprecated alias of github.com/maxatome/go-testdeep/td.Zero.
	Zero = td.Zero
)

const (
	// BoundsInIn is a deprecated alias of github.com/maxatome/go-testdeep/td.BoundsInIn.
	BoundsInIn = td.BoundsInIn
	// BoundsInOut is a deprecated alias of github.com/maxatome/go-testdeep/td.BoundsInOut.
	BoundsInOut = td.BoundsInOut
	// BoundsOutIn is a deprecated alias of github.com/maxatome/go-testdeep/td.BoundsOutIn.
	BoundsOutIn = td.BoundsOutIn
	// BoundsOutOut is a deprecated alias of github.com/maxatome/go-testdeep/td.BoundsOutOut.
	BoundsOutOut = td.BoundsOutOut
)
