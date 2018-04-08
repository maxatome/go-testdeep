package testdeep

import (
	"bytes"
	"reflect"
)

type tdList struct {
	BaseOKNil
	items []reflect.Value
}

func newList(items ...interface{}) (ret tdList) {
	ret.BaseOKNil = NewBaseOKNil(4)
	ret.items = make([]reflect.Value, len(items))

	for idx, item := range items {
		ret.items[idx] = reflect.ValueOf(item)
	}
	return
}

func (l *tdList) String() string {
	return sliceToBuffer(bytes.NewBufferString(l.GetLocation().Func), l.items).
		String()
}
