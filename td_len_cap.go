package testdeep

import (
	"fmt"
	"reflect"
)

type tdCapLen struct {
	TestDeepBase
	expectedMin int
	expectedMax int
	cap         bool
}

func (c *tdCapLen) initCapLen(min int, max ...int) bool {
	c.TestDeepBase = NewTestDeepBase(4)

	c.expectedMin = min

	if len(max) > 0 {
		if len(max) > 1 {
			return false
		}
		c.expectedMax = max[0]
	} else {
		c.expectedMax = c.expectedMin
	}

	if c.expectedMax < c.expectedMin {
		c.expectedMin, c.expectedMax = c.expectedMax, c.expectedMin
	}

	return true
}

func (c *tdCapLen) toString(name string) string {
	if c.expectedMin == c.expectedMax {
		return fmt.Sprintf("%s=%d", name, c.expectedMin)
	}
	return fmt.Sprintf("%d ≤ %s ≤ %d", c.expectedMin, name, c.expectedMax)
}

func (c *tdCapLen) checkVal(val int) bool {
	return val >= c.expectedMin && val <= c.expectedMax
}

func (c *tdCapLen) expectedRaw(name string) testDeepStringer {
	if c.expectedMin == c.expectedMax {
		return rawInt(c.expectedMin)
	}
	return rawString(c.toString(name))
}

//
//
type tdLen struct {
	tdCapLen
}

var _ TestDeep = &tdLen{}

func Len(min int, max ...int) TestDeep {
	l := tdLen{}

	if !l.tdCapLen.initCapLen(min, max...) {
		panic("usage: Len(MIN_LEN[, MAX_LEN])")
	}

	return &l
}

func (l *tdLen) String() string {
	return l.toString("len")
}

func (l *tdLen) Match(ctx Context, got reflect.Value) (err *Error) {
	switch got.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map,
		reflect.Slice, reflect.String:
		if l.checkVal(got.Len()) {
			return nil
		}
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad length",
			Got:      rawInt(got.Len()),
			Expected: l.expectedRaw("len"),
			Location: l.GetLocation(),
		}

	default:
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: rawString("Array, Chan, Map, Slice or string"),
			Location: l.GetLocation(),
		}
	}
}

//
//
type tdCap struct {
	tdCapLen
}

var _ TestDeep = &tdCap{}

func Cap(min int, max ...int) TestDeep {
	c := tdCap{}

	if !c.tdCapLen.initCapLen(min, max...) {
		panic("usage: Cap(MIN_LEN[, MAX_LEN])")
	}

	return &c
}
func (c *tdCap) String() string {
	return c.toString("cap")
}

func (c *tdCap) Match(ctx Context, got reflect.Value) (err *Error) {
	switch got.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if c.checkVal(got.Cap()) {
			return nil
		}
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad capacity",
			Got:      rawInt(got.Cap()),
			Expected: c.expectedRaw("cap"),
			Location: c.GetLocation(),
		}

	default:
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad type",
			Got:      rawString(got.Type().String()),
			Expected: rawString("Array, Chan or Slice"),
			Location: c.GetLocation(),
		}
	}
}
