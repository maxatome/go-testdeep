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

// summary(Bag): compares the contents of an array or a slice without taking
// care of the order of items
// input(Bag): array,slice,ptr(ptr on array/slice)

// Bag operator compares the contents of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice, and each array/slice item should be matched by an
// expected item to succeed.
//
//   Cmp(t, []int{1, 1, 2}, Bag(1, 1, 2))    // succeeds
//   Cmp(t, []int{1, 1, 2}, Bag(1, 2, 1))    // succeeds
//   Cmp(t, []int{1, 1, 2}, Bag(2, 1, 1))    // succeeds
//   Cmp(t, []int{1, 1, 2}, Bag(1, 2))       // fails, one 1 is missing
//   Cmp(t, []int{1, 1, 2}, Bag(1, 2, 1, 3)) // fails, 3 is missing
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

// summary(SubBagOf): compares the contents of an array or a slice
// without taking care of the order of items but with potentially some
// exclusions
// input(SubBagOf): array,slice,ptr(ptr on array/slice)

// SubBagOf operator compares the contents of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each array/slice item should be matched by an
// expected item to succeed. But some expected items can be missing
// from the compared array/slice.
//
//   Cmp(t, []int{1}, SubBagOf(1, 1, 2))       // succeeds
//   Cmp(t, []int{1, 1, 1}, SubBagOf(1, 1, 2)) // fails, one 1 is an extra item
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

// summary(SuperBagOf): compares the contents of an array or a slice
// without taking care of the order of items but with potentially some
// extra items
// input(SuperBagOf): array,slice,ptr(ptr on array/slice)

// SuperBagOf operator compares the contents of an array or a slice (or a
// pointer on array/slice) without taking care of the order of items.
//
// During a match, each expected item should match in the compared
// array/slice. But some items in the compared array/slice may not be
// expected.
//
//   Cmp(t, []int{1, 1, 2}, SuperBagOf(1))       // succeeds
//   Cmp(t, []int{1, 1, 2}, SuperBagOf(1, 1, 1)) // fails, one 1 is missing
func SuperBagOf(expectedItems ...interface{}) TestDeep {
	bag := &tdSuperBagOf{
		tdSetBase: newSetBase(superSet, false),
	}
	bag.Add(expectedItems...)
	return bag
}
