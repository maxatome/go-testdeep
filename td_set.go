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

// summary(Set): compares the contents of an array or a slice ignoring
// duplicates and without taking care of the order of items
// input(Set): array,slice,ptr(ptr on array/slice)

// Set operator compares the contents of an array or a slice (or a
// pointer on array/slice) ignoring duplicates and without taking care
// of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice, and each array/slice item should be matched by an
// expected item to succeed.
//
//   Cmp(t, []int{1, 1, 2}, Set(1, 2))    // succeeds
//   Cmp(t, []int{1, 1, 2}, Set(2, 1))    // succeeds
//   Cmp(t, []int{1, 1, 2}, Set(1, 2, 3)) // fails, 3 is missing
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

// summary(SubSetOf): compares the contents of an array or a slice
// ignoring duplicates and without taking care of the order of items
// but with potentially some exclusions
// input(SubSetOf): array,slice,ptr(ptr on array/slice)

// SubSetOf operator compares the contents of an array or a slice (or a
// pointer on array/slice) ignoring duplicates and without taking care
// of the order of items.
//
// During a match, each array/slice item should be matched by an
// expected item to succeed. But some expected items can be missing
// from the compared array/slice.
//
//   Cmp(t, []int{1, 1}, SubSetOf(1, 2))    // succeeds
//   Cmp(t, []int{1, 1, 2}, SubSetOf(1, 3)) // fails, 2 is an extra item
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

// summary(SuperSetOf): compares the contents of an array or a slice
// ignoring duplicates and without taking care of the order of items
// but with potentially some extra items
// input(SuperSetOf): array,slice,ptr(ptr on array/slice)

// SuperSetOf operator compares the contents of an array or a slice (or
// a pointer on array/slice) ignoring duplicates and without taking
// care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice. But some items in the compared array/slice may not be
// expected.
//
//   Cmp(t, []int{1, 1, 2}, SuperSetOf(1))    // succeeds
//   Cmp(t, []int{1, 1, 2}, SuperSetOf(1, 3)) // fails, 3 is missing
func SuperSetOf(expectedItems ...interface{}) TestDeep {
	set := &tdSuperSetOf{
		tdSetBase: newSetBase(superSet, true),
	}
	set.Add(expectedItems...)
	return set
}

type tdNotAny struct {
	tdSetBase
}

var _ TestDeep = &tdNotAny{}

// summary(NotAny): compares the contents of an array or a slice, no
// values have to match
// input(NotAny): array,slice,ptr(ptr on array/slice)

// NotAny operator checks that the contents of an array or a slice (or
// a pointer on array/slice) does not contain any of "expectedItems".
//
//   Cmp(t, []int{1}, NotAny(1, 2, 3)) // fails
//   Cmp(t, []int{5}, NotAny(1, 2, 3)) // succeeds
func NotAny(expectedItems ...interface{}) TestDeep {
	set := &tdNotAny{
		tdSetBase: newSetBase(noneSet, true),
	}
	set.Add(expectedItems...)
	return set
}
