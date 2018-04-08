package testdeep

import (
	"reflect"
)

type tdIgnore struct {
	BaseOKNil
}

var ignoreSingleton TestDeep = &tdIgnore{}

func Ignore() TestDeep {
	// NewBase() useless
	return ignoreSingleton
}

func (i *tdIgnore) Match(ctx Context, got reflect.Value) *Error {
	return nil
}

func (i *tdIgnore) String() string {
	return "Ignore()"
}
