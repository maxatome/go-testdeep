package testdeep

import (
	"reflect"
)

type tdIgnore struct {
	TestDeepBaseOKNil
}

var ignoreSingleton TestDeep = &tdIgnore{}

func Ignore() TestDeep {
	// NewTestDeepBase() useless
	return ignoreSingleton
}

func (i *tdIgnore) Match(ctx Context, got reflect.Value) *Error {
	return nil
}

func (i *tdIgnore) String() string {
	return "Ignore()"
}
