package testdeep

import (
	"fmt"
	"reflect"
)

type tdNone struct {
	tdList
}

var _ TestDeep = &tdNone{}

// None operator compares data against several expected values. During
// a match, none of them have to match to succeed.
func None(expectedValues ...interface{}) TestDeep {
	return &tdNone{
		tdList: newList(expectedValues...),
	}
}

// Not operator compares data against the expected value. During a
// match, it must not match to succeed.
//
// Not is the same operator as None() with only one argument. It is
// provided as a more readable function when only one argument is
// needed.
func Not(expected interface{}) TestDeep {
	return &tdNone{
		tdList: newList(expected),
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
