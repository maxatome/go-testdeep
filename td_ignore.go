package testdeep

import (
	"reflect"
)

type tdIgnore struct {
	BaseOKNil
}

var ignoreSingleton TestDeep = &tdIgnore{}

// Ignore operator is always true, whatever data is. It is useful when
// comparing a slice and wanting to ignore some indexes, for example.
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
