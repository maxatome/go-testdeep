// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

type tdBag struct {
	tdSetBase
}

var _ TestDeep = &tdBag{}

// Bag operator compares the content of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice, and each array/slice item should be matched by an
// expected item to succeed.
//
//   CmpDeeply(t, []int{1, 1, 2}, Bag(1, 1, 2))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, Bag(1, 2, 1))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, Bag(2, 1, 1))    // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, Bag(1, 2))       // fails, one 1 is missing
//   CmpDeeply(t, []int{1, 1, 2}, Bag(1, 2, 1, 3)) // fails, 3 is missing
func Bag(expectedItems ...interface{}) TestDeep {
	bag := &tdBag{
		tdSetBase: newSetBase(allSet, false),
	}
	bag.Add(expectedItems...)
	return bag
}

type tdSubBagOf struct {
	tdSetBase
}

var _ TestDeep = &tdSubBagOf{}

// SubBagOf operator compares the content of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each array/slice item should be matched by an
// expected item to succeed. But some expected items can be missing
// from the compared array/slice.
//
//   CmpDeeply(t, []int{1}, SubBagOf(1, 1, 2))       // succeeds
//   CmpDeeply(t, []int{1, 1, 1}, SubBagOf(1, 1, 2)) // fails, one 1 is an extra item
func SubBagOf(expectedItems ...interface{}) TestDeep {
	bag := &tdSubBagOf{
		tdSetBase: newSetBase(subSet, false),
	}
	bag.Add(expectedItems...)
	return bag
}

type tdSuperBagOf struct {
	tdSetBase
}

var _ TestDeep = &tdSuperBagOf{}

// SuperBagOf operator compares the content of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice. But some items in the compared array/slice may not be
// expected.
//
//   CmpDeeply(t, []int{1, 1, 2}, SuperBagOf(1))       // succeeds
//   CmpDeeply(t, []int{1, 1, 2}, SuperBagOf(1, 1, 1)) // fails, one 1 is missing
func SuperBagOf(expectedItems ...interface{}) TestDeep {
	bag := &tdSuperBagOf{
		tdSetBase: newSetBase(superSet, false),
	}
	bag.Add(expectedItems...)
	return bag
}
