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
//   CmpDeeply(t, []int{1}, SubBagOf(1, 1, 2))       // is true
//   CmpDeeply(t, []int{1, 1, 1}, SubBagOf(1, 1, 2)) // is false
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
//   CmpDeeply(t, []int{1, 1, 2}, SuperBagOf(1))       // is true
//   CmpDeeply(t, []int{1, 1, 2}, SuperBagOf(1, 1, 1)) // is false
func SuperBagOf(expectedItems ...interface{}) TestDeep {
	bag := &tdSuperBagOf{
		tdSetBase: newSetBase(superSet, false),
	}
	bag.Add(expectedItems...)
	return bag
}
