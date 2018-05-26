// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

type tdSet struct {
	tdSetBase
}

var _ TestDeep = &tdSet{}

// Set operator compares the contents of an array or a slice (or a
// pointer on array/slice) ignoring duplicates and without taking care
// of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice, and each array/slice item should be matched by an
// expected item to succeed.
//
//   CmpDeeply(t, []int{1, 1, 2}, Set(1, 2))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, Set(2, 1))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, Set(1, 2, 3)) // fails, 3 is missing
func Set(expectedItems ...interface{}) TestDeep {
	set := &tdSet{
		tdSetBase: newSetBase(allSet, true),
	}
	set.Add(expectedItems...)
	return set
}

type tdSubSetOf struct {
	tdSetBase
}

var _ TestDeep = &tdSubSetOf{}

// SubSetOf operator compares the contents of an array or a slice (or a
// pointer on array/slice) ignoring duplicates and without taking care
// of the order of items.
//
// During a match, each array/slice item should be matched by an
// expected item to succeed. But some expected items can be missing
// from the compared array/slice.
//
//   CmpDeeply(t, []int{1, 1}, SubSetOf(1, 2))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, SubSetOf(1, 3)) // fails, 2 is an extra item
func SubSetOf(expectedItems ...interface{}) TestDeep {
	set := &tdSubSetOf{
		tdSetBase: newSetBase(subSet, true),
	}
	set.Add(expectedItems...)
	return set
}

type tdSuperSetOf struct {
	tdSetBase
}

var _ TestDeep = &tdSuperSetOf{}

// SuperSetOf operator compares the contents of an array or a slice (or
// a pointer on array/slice) ignoring duplicates and without taking
// care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice. But some items in the compared array/slice may not be
// expected.
//
//   CmpDeeply(t, []int{1, 1, 2}, SuperSetOf(1))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, SuperSetOf(1, 3)) // fails, 3 is missing
func SuperSetOf(expectedItems ...interface{}) TestDeep {
	set := &tdSuperSetOf{
		tdSetBase: newSetBase(superSet, true),
	}
	set.Add(expectedItems...)
	return set
}

type tdNoneOf struct {
	tdSetBase
}

var _ TestDeep = &tdNoneOf{}

// NoneOf operator checks that the contents of an array or a slice (or
// a pointer on array/slice) does not contain any of "expectedItems".
//
//   CmpDeeply(t, []int{1}, NoneOf(1, 2, 3)) // fails
//   CmpDeeply(t, []int{5}, NoneOf(1, 2, 3)) // succeeds
func NoneOf(expectedItems ...interface{}) TestDeep {
	set := &tdNoneOf{
		tdSetBase: newSetBase(noneSet, true),
	}
	set.Add(expectedItems...)
	return set
}
