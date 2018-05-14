package testdeep

import (
	"reflect"
)

// tdSmuggler is the base class of all smuggler TestDeep operators.
type tdSmuggler struct {
	Base
	expectedValue reflect.Value
	isTestDeeper  bool
}

func newSmuggler(val interface{}) (ret tdSmuggler) {
	ret.Base = NewBase(4)

	// Initializes only if TestDeep operator. Other cases are specific.
	if _, ok := val.(TestDeep); ok {
		ret.expectedValue = reflect.ValueOf(val)
		ret.isTestDeeper = true
	}
	return
}
