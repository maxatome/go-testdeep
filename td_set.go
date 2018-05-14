package testdeep

type tdSet struct {
	tdSetBase
}

var _ TestDeep = &tdSet{}

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

func NoneOf(expectedItems ...interface{}) TestDeep {
	set := &tdNoneOf{
		tdSetBase: newSetBase(noneSet, true),
	}
	set.Add(expectedItems...)
	return set
}
