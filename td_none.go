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

// Not is the same as None() with only one argument. It is provided as
// a more readable function when only one argument is needed.
func Not(item interface{}) TestDeep {
	return &tdNone{
		tdList: newList(item),
	}
}

func (n *tdNone) Match(ctx Context, got reflect.Value) *Error {
	for idx, item := range n.items {
		if deepValueEqualOK(got, item) {
			if ctx.booleanError {
				return booleanError
			}

			var mesg string
			if n.GetLocation().Func == "Not" {
				mesg = "comparing with Not"
			} else {
				mesg = fmt.Sprintf("comparing with None (part %d of %d is OK)",
					idx+1, len(n.items))
			}
			return &Error{
				Context:  ctx,
				Message:  mesg,
				Got:      got,
				Expected: n,
				Location: n.GetLocation(),
			}
		}
	}
	return nil
}
