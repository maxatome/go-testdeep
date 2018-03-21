package testdeep

type tdBag struct {
	tdSetBase
}

var _ TestDeep = &tdBag{}

func Bag(items ...interface{}) TestDeep {
	bag := &tdBag{
		tdSetBase: newSetBase(allSet, false),
	}
	bag.Add(items...)
	return bag
}

type tdSubBagOf struct {
	tdSetBase
}

var _ TestDeep = &tdSubBagOf{}

func SubBagOf(items ...interface{}) TestDeep {
	bag := &tdSubBagOf{
		tdSetBase: newSetBase(subSet, false),
	}
	bag.Add(items...)
	return bag
}

type tdSuperBagOf struct {
	tdSetBase
}

var _ TestDeep = &tdSuperBagOf{}

func SuperBagOf(items ...interface{}) TestDeep {
	bag := &tdSuperBagOf{
		tdSetBase: newSetBase(superSet, false),
	}
	bag.Add(items...)
	return bag
}
