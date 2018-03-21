package testdeep

import (
	"fmt"
	"reflect"
	"regexp"
)

type tdRe struct {
	TestDeepBase
	re         *regexp.Regexp
	captures   reflect.Value
	numMatches int
}

var _ TestDeep = &tdRe{}

func newRe(usage string, opts ...interface{}) (r *tdRe) {
	r = &tdRe{
		TestDeepBase: NewTestDeepBase(4),
		numMatches:   1, // only one match by default, see Global() for all
	}

	switch len(opts) {
	case 0:
		return

	case 2:
		isGlobal, ok := opts[1].(bool)
		if !ok {
			break
		}
		if isGlobal {
			r.numMatches = -1
		}
		fallthrough

	case 1:
		if opts[0] != nil {
			r.captures = reflect.ValueOf(opts[0])
			return
		}
	}

	panic(usage)
}

func Re(reg string, opts ...interface{}) TestDeep {
	r := newRe("usage: Re(REGEXP_STR[, NON_NIL_CAPTURE[, IS_GLOBAL]])", opts...)
	r.re = regexp.MustCompile(reg)
	return r
}

func Rex(re *regexp.Regexp, opts ...interface{}) TestDeep {
	r := newRe("usage: Rex(*Regexp, NON_NIL_CAPTURE[, IS_GLOBAL]])", opts...)
	r.re = re
	return r
}

func (r *tdRe) needCaptures() bool {
	return r.captures.IsValid()
}

func (r *tdRe) matchByteCaptures(ctx Context, got []byte, result [][][]byte) *Error {
	if len(result) == 0 {
		return r.doesNotMatch(ctx, got)
	}

	num := 0
	for _, set := range result {
		num += len(set) - 1
	}

	// Not perfect but cast captured groups to string
	captures := make([]string, 0, num)
	for _, set := range result {
		for _, match := range set[1:] {
			captures = append(captures, string(match))
		}
	}

	return r.matchCaptures(ctx, captures)
}

func (r *tdRe) matchStringCaptures(ctx Context, got string, result [][]string) *Error {
	if len(result) == 0 {
		return r.doesNotMatch(ctx, got)
	}

	num := 0
	for _, set := range result {
		num += len(set) - 1
	}

	captures := make([]string, 0, num)
	for _, set := range result {
		captures = append(captures, set[1:]...)
	}

	return r.matchCaptures(ctx, captures)
}

func (r *tdRe) matchCaptures(ctx Context, captures []string) *Error {
	newCtx := NewContext("(" + ctx.Path + " =~ " + r.String() + ")")
	if ctx.booleanError {
		newCtx.booleanError = true
	}

	return deepValueEqual(reflect.ValueOf(captures), r.captures, newCtx)
}

func (r *tdRe) matchBool(ctx Context, got interface{}, result bool) *Error {
	if result {
		return nil
	}
	return r.doesNotMatch(ctx, got)
}

func (r *tdRe) doesNotMatch(ctx Context, got interface{}) *Error {
	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "does not match Regexp",
		Got:      got,
		Expected: rawString(r.re.String()),
		Location: r.GetLocation(),
	}
}

func (r *tdRe) Match(ctx Context, got reflect.Value) *Error {
	var str string
	switch got.Kind() {
	case reflect.String:
		str = got.String()

	case reflect.Slice:
		if got.Type().Elem().Kind() == reflect.Uint8 {
			gotBytes := got.Bytes()
			if r.needCaptures() {
				return r.matchByteCaptures(ctx,
					gotBytes, r.re.FindAllSubmatch(gotBytes, r.numMatches))
			}
			return r.matchBool(ctx, gotBytes, r.re.Match(gotBytes))
		}
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "bad slice type",
			Got:      rawString("[]" + got.Type().Elem().Kind().String()),
			Expected: rawString("[]uint8"),
			Location: r.GetLocation(),
		}

	default:
		var strOK bool
		if got.CanInterface() {
			switch iface := got.Interface().(type) {
			case fmt.Stringer:
				str = iface.String()
				strOK = true
			case error:
				str = iface.Error()
				strOK = true
			default:
			}
		}

		if !strOK {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context: ctx,
				Message: "bad type",
				Got:     rawString(got.Type().String()),
				Expected: rawString(
					"string (convertible) OR fmt.Stringer OR error OR []uint8"),
				Location: r.GetLocation(),
			}
		}
	}

	if r.needCaptures() {
		return r.matchStringCaptures(ctx,
			str, r.re.FindAllStringSubmatch(str, r.numMatches))
	}
	return r.matchBool(ctx, str, r.re.MatchString(str))
}

func (r *tdRe) String() string {
	return r.re.String()
}
