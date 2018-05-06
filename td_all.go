package testdeep

import (
	"fmt"
	"reflect"
)

type tdAll struct {
	tdList
}

var _ TestDeep = &tdAll{}

func All(items ...interface{}) TestDeep {
	return &tdAll{
		tdList: newList(items...),
	}
}

func (a *tdAll) Match(ctx Context, got reflect.Value) (err *Error) {
	for idx, item := range a.items {
		origErr := deepValueEqual(
			ctx.AddDepth(fmt.Sprintf("<All#%d/%d>", idx+1, len(a.items))),
			got, item)
		if origErr != nil {
			if ctx.booleanError {
				return booleanError
			}
			err = &Error{
				Context:  ctx,
				Message:  fmt.Sprintf("compared (part %d of %d)", idx+1, len(a.items)),
				Got:      got,
				Expected: item,
				Location: a.GetLocation(),
			}

			if item.IsValid() && item.Type().Implements(testDeeper) {
				err.Origin = origErr
			}
			return
		}
	}
	return
}
