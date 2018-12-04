// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
)

// SmuggledGot can be returned by a Smuggle function to name the
// transformed / returned value.
type SmuggledGot struct {
	Name string
	Got  interface{}
}

const smuggled = "<smuggled>"

func (s SmuggledGot) contextAndGot(ctx ctxerr.Context) (ctxerr.Context, reflect.Value) {
	// If the Name starts with a Letter, prefix it by a "."
	var name string
	if s.Name != "" {
		first, _ := utf8.DecodeRuneInString(s.Name)
		if unicode.IsLetter(first) {
			name = "."
		}
		name += s.Name
	} else {
		name = smuggled
	}
	return ctx.AddDepth(name), reflect.ValueOf(s.Got)
}

type tdSmuggle struct {
	tdSmugglerBase
	function reflect.Value
	argType  reflect.Type
}

var _ TestDeep = &tdSmuggle{}

// Smuggle operator allows to change data contents or mutate it into
// another type before stepping down in favor of generic comparison
// process. So "fn" is a function that must take one parameter whose
// type must be convertible to the type of the compared value.
//
// "fn" must return at least one value, these value will be compared as is
// to "expectedValue", here integer 28:
//
//   Smuggle(func (value string) int {
//       num, _ := strconv.Atoi(value)
//       return num
//     },
//     28)
//
// or using an other TestDeep operator, here Between(28, 30):
//
//   Smuggle(func (value string) int {
//       num, _ := strconv.Atoi(value)
//       return num
//     },
//     Between(28, 30))
//
// "fn" can return a second boolean value, used to tell that a problem
// occurred and so stop the comparison:
//
//   Smuggle(func (value string) (int, bool) {
//       num, err := strconv.Atoi(value)
//       return num, err == nil
//     },
//     Between(28, 30))
//
// "fn" can return a third string value which is used to describe the
// test when a problem occurred (false second boolean value):
//
//   Smuggle(func (value string) (int, bool, string) {
//       num, err := strconv.Atoi(value)
//       if err != nil {
//         return 0, false, "string must contain a number"
//       }
//       return num, true, ""
//     },
//     Between(28, 30))
//
// Instead of returning (X, bool) or (X, bool, string), "fn" can
// return (X, error). When a problem occurs, the returned error is
// non-nil, as in:
//
//   Smuggle(func (value string) (int, error) {
//       num, err := strconv.Atoi(value)
//       return num, err
//     },
//     Between(28, 30))
//
// Which can be simplified to:
//
//   Smuggle(strconv.Atoi, Between(28, 30))
//
// Imagine you want to compare that the Year of a date is between 2010
// and 2020:
//
//   Smuggle(func (date time.Time) int {
//       return date.Year()
//     },
//     Between(2010, 2020))
//
// In this case the data location forwarded to next test will be
// something like DATA.MyTimeField<smuggled>, but you can act on it
// too by returning a SmuggledGot struct (by value or by address):
//
//   Smuggle(func (date time.Time) SmuggledGot {
//       return SmuggledGot{
//         Name: "Year",
//         Got:  date.Year(),
//       }
//     },
//     Between(2010, 2020))
//
// then the data location forwarded to next test will be something like
// DATA.MyTimeField.Year. The "."  between the current path (here
// "DATA.MyTimeField") and the returned Name "Year" is automatically
// added when Name starts with a Letter.
//
// Note that SmuggledGot and *SmuggledGot returns are treated equally,
// and they are only used when "fn" has only one returned value or
// when the second boolean returned value is true.
//
// Of course, all cases can go together:
//
//   // Accepts a "YYYY/mm/DD HH:MM:SS" string to produce a time.Time and tests
//   // whether this date is contained between 2 hours before now and now.
//   Smuggle(func (date string) (*SmuggledGot, bool, string) {
//       date, err := time.Parse("2006/01/02 15:04:05", date)
//       if err != nil {
//         return nil, false, `date must conform to "YYYY/mm/DD HH:MM:SS" format`
//       }
//       return &SmuggledGot{
//         Name: "Date",
//         Got:  date,
//       }, true, ""
//     },
//     Between(time.Now().Add(-2*time.Hour), time.Now()))
//
// or:
//
//   // Accepts a "YYYY/mm/DD HH:MM:SS" string to produce a time.Time and tests
//   // whether this date is contained between 2 hours before now and now.
//   Smuggle(func (date string) (*SmuggledGot, error) {
//       date, err := time.Parse("2006/01/02 15:04:05", date)
//       if err != nil {
//         return nil, err
//       }
//       return &SmuggledGot{
//         Name: "Date",
//         Got:  date,
//       }, nil
//     },
//     Between(time.Now().Add(-2*time.Hour), time.Now()))
//
// The difference between Smuggle and Code operators is that Code is
// used to do a final comparison while Smuggle transforms the data and
// then steps down in favor of generic comparison process. Moreover,
// the type accepted as input for the function is lax to facilitate
// the tests writing (eg. the function can accept an float64 and the
// got value be an int). See examples. On the other hand, the output
// type is strict and must match exactly the expected value type.
//
// TypeBehind method returns the reflect.Type of only parameter of "fn".
func Smuggle(fn interface{}, expectedValue interface{}) TestDeep {
	vfn := reflect.ValueOf(fn)

	const usage = "Smuggle(FUNC, TESTDEEP_OPERATOR|EXPECTED_VALUE)"

	if vfn.Kind() != reflect.Func {
		panic("usage: " + usage)
	}

	fnType := vfn.Type()
	if fnType.NumIn() != 1 {
		panic(usage + ": FUNC must take only one argument")
	}

	switch fnType.NumOut() {
	case 3: // (value, bool, string)
		if fnType.Out(2).Kind() != reflect.String {
			break
		}
		fallthrough

	case 2:
		// (value, *bool*) or (value, *bool*, string)
		if fnType.Out(1).Kind() != reflect.Bool &&
			// (value, *error*)
			(fnType.NumOut() > 2 ||
				fnType.Out(1) != errorInterface) {
			break
		}
		fallthrough

	case 1: // (value)
		s := tdSmuggle{
			tdSmugglerBase: newSmugglerBase(expectedValue),
			function:       vfn,
			argType:        fnType.In(0),
		}
		if !s.isTestDeeper {
			s.expectedValue = reflect.ValueOf(expectedValue)
		}
		return &s
	}

	panic(usage +
		": FUNC must return value or (value, bool) or (value, bool, string) or (value, error)")
}

func (s *tdSmuggle) laxConvert(got reflect.Value) (reflect.Value, bool) {
	if !got.Type().ConvertibleTo(s.argType) {
		if got.Kind() != reflect.Interface || got.IsNil() {
			return got, false
		}

		got = got.Elem()
		if !got.Type().ConvertibleTo(s.argType) {
			return got, false
		}
	}

	return got.Convert(s.argType), true
}

func (s *tdSmuggle) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	got, ok := s.laxConvert(got)
	if !ok {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "incompatible parameter type",
			Got:      types.RawString(got.Type().String()),
			Expected: types.RawString(s.argType.String()),
		})
	}

	// Refuse to override unexported fields access in this case. It is a
	// choice, as we think it is better to work on surrounding struct
	// instead.
	if !got.CanInterface() {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message: "cannot smuggle unexported field",
			Summary: types.RawString("work on surrounding struct instead"),
		})
	}

	ret := s.function.Call([]reflect.Value{got})
	if len(ret) == 1 ||
		(ret[1].Kind() == reflect.Bool && ret[1].Bool()) ||
		(ret[1].Kind() == reflect.Interface && ret[1].IsNil()) {
		newGot := ret[0]

		var newCtx ctxerr.Context
		if newGot.IsValid() {
			switch newGot.Type() {
			case smuggledGotType:
				newCtx, newGot = newGot.Interface().(SmuggledGot).contextAndGot(ctx)

			case smuggledGotPtrType:
				if smGot := newGot.Interface().(*SmuggledGot); smGot == nil {
					newCtx, newGot = ctx, reflect.ValueOf(nil)
				} else {
					newCtx, newGot = smGot.contextAndGot(ctx)
				}

			default:
				newCtx = ctx.AddDepth(smuggled)
			}
		}
		return deepValueEqual(newCtx, newGot, s.expectedValue)
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}

	summary := tdCodeResult{
		Value: got,
	}

	switch len(ret) {
	case 3: // (value, false, string)
		summary.Reason = ret[2].String()
	case 2:
		// (value, error)
		if ret[1].Kind() == reflect.Interface {
			summary.Reason = ret[1].Interface().(error).Error()
		}
		// (value, false)
	}
	return ctx.CollectError(&ctxerr.Error{
		Message: "ran smuggle code with %% as argument",
		Summary: summary,
	})
}

func (s *tdSmuggle) String() string {
	return "Smuggle(" + s.function.Type().String() + ")"
}

func (s *tdSmuggle) TypeBehind() reflect.Type {
	return s.argType
}
