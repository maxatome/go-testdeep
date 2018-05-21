package testdeep

import (
	"fmt"
	"reflect"
	"regexp"
)

type tdRe struct {
	Base
	re         *regexp.Regexp
	captures   reflect.Value
	numMatches int
}

var _ TestDeep = &tdRe{}

func newRe(regIf interface{}, capture ...interface{}) (r *tdRe) {
	r = &tdRe{
		Base: NewBase(4),
	}

	switch len(capture) {
	case 0:
	case 1:
		if capture[0] != nil {
			r.captures = reflect.ValueOf(capture[0])
		}
		break
	default:
		r.usage()
	}

	switch reg := regIf.(type) {
	case *regexp.Regexp:
		r.re = reg
	case string:
		r.re = regexp.MustCompile(reg)
	default:
		r.usage()
	}
	return
}

// Re operator allows to apply a regexp on a string (or convertible),
// []byte, error or fmt.Stringer interface (error interface is tested
// before fmt.Stringer.)
//
// "reg" is the regexp. It can be a string that is automatically
// compiled using regexp.MustCompile, or a *regexp.Regexp.
//
// Optional "capture" parameter can be used to match the contents of
// regexp groups. Groups are presented as a []string or [][]byte
// depending the original matched data. Note that an other operator
// can be used here.
//
//   CmpDeeply(t, "foobar zip!", Re(`^foobar`))     // is true
//   CmpDeeply(t, "John Doe",
//     Re(`^(\w+) (\w+)`, []string{"John", "Doe"})) // is true
//   CmpDeeply(t, "John Doe",
//     Re(`^(\w+) (\w+)`, Bag("Doe", "John"))       // is true
func Re(reg interface{}, capture ...interface{}) TestDeep {
	r := newRe(reg, capture...)
	r.numMatches = 1
	return r
}

// ReAll operator allows to successively apply a regexp on a string
// (or convertible), []byte, error or fmt.Stringer interface (error
// interface is tested before fmt.Stringer) and to match its groups
// contents.
//
// "reg" is the regexp. It can be a string that is automatically
// compiled using regexp.MustCompile, or a *regexp.Regexp.
//
// "capture" is used to match the contents of regexp groups. Groups
// are presented as a []string or [][]byte depending the original
// matched data. Note that an other operator can be used here.
//
//   CmpDeeply(t, "John Doe",
//     ReAll(`(\w+)(?: |\z)`, []string{"John", "Doe"})) // is true
//   CmpDeeply(t, "John Doe",
//     ReAll(`(\w+)(?: |\z)`, Bag("Doe", "John"))       // is true
func ReAll(reg interface{}, capture interface{}) TestDeep {
	r := newRe(reg, capture)
	r.numMatches = -1
	return r
}

func (r *tdRe) usage() {
	panic(fmt.Sprintf("usage: %s(STRING|*regexp.Regexp[, NON_NIL_CAPTURE])",
		r.location.Func))
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
	newCtx := NewContext("(" + ctx.path + " =~ " + r.String() + ")")
	if ctx.booleanError {
		newCtx.booleanError = true
	}

	return deepValueEqual(newCtx, reflect.ValueOf(captures), r.captures)
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
		if iface, ok := getInterface(got, true); ok {
			switch gotVal := iface.(type) {
			case error:
				str = gotVal.Error()
				strOK = true
			case fmt.Stringer:
				str = gotVal.String()
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
