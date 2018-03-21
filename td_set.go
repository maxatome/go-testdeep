package testdeep

type tdSet struct {
	tdSetBase
}

var _ TestDeep = &tdSet{}

func Set(items ...interface{}) TestDeep {
	set := &tdSet{
		tdSetBase: newSetBase(allSet, true),
	}
	set.Add(items...)
	return set
}

type tdSubSetOf struct {
	tdSetBase
}

var _ TestDeep = &tdSubSetOf{}

func SubSetOf(items ...interface{}) TestDeep {
	set := &tdSubSetOf{
		tdSetBase: newSetBase(subSet, true),
	}
	set.Add(items...)
	return set
}

type tdSuperSetOf struct {
	tdSetBase
}

var _ TestDeep = &tdSuperSetOf{}

func SuperSetOf(items ...interface{}) TestDeep {
	set := &tdSuperSetOf{
		tdSetBase: newSetBase(superSet, true),
	}
	set.Add(items...)
	return set
}

type tdNoneOf struct {
	tdSetBase
}

var _ TestDeep = &tdNoneOf{}

func NoneOf(items ...interface{}) TestDeep {
	set := &tdNoneOf{
		tdSetBase: newSetBase(noneSet, true),
	}
	set.Add(items...)
	return set
}
