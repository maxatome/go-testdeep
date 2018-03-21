package testdeep

import (
	"fmt"
	"reflect"
)

type tdNone struct {
	tdList
}

var _ TestDeep = &tdNone{}

func None(items ...interface{}) TestDeep {
	return &tdNone{
		tdList: newList(items...),
	}
}

func (n *tdNone) Match(ctx Context, got reflect.Value) *Error {
	for idx, item := range n.items {
		if deepValueEqualOK(got, item) {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context: ctx,
				Message: fmt.Sprintf(
					"comparing with None (part %d of %d is OK)", idx+1, len(n.items)),
				Got:      got,
				Expected: n,
				Location: n.GetLocation(),
			}
		}
	}
	return nil
}
