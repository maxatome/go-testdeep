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
)

// SmuggledGot can be returned by a Smuggle function to name the
// transformed / returned value.
type SmuggledGot struct {
	Name string
	Got  interface{}
}

const smuggled = "<smuggled>"

func (s SmuggledGot) contextAndGot(ctx Context) (Context, reflect.Value) {
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
// type must be the same as the type of the compared value.
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
// Imagine you want to compare that the Year of a date is between 2010
// and 2020:
//
//   Smuggle(func (date time.Time) int {
//       return date.Year()
//     },
//     Between(2010, 2020))
//
// In this case the data location forwarded to next test will be
// somthing like DATA.MyTimeField<smuggled>, but you can act on it too
// by returning a SmuggledGot struct (by value or by address):
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
// The difference between Smuggle and Code operators is that Code is
// used to do a final comparison while Smuggle transforms the data and
// then steps down in favor of generic comparison process.
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

	case 2: // (value, bool)
		if fnType.Out(1).Kind() != reflect.Bool {
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
		": FUNC must return value or (value, bool) or (value, bool, string)")
}

func (s *tdSmuggle) Match(ctx Context, got reflect.Value) *Error {
	if !got.Type().AssignableTo(s.argType) {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message:  "incompatible parameter type",
			Got:      rawString(got.Type().String()),
			Expected: rawString(s.argType.String()),
		})
	}

	// Refuse to override unexported fields access in this case. It is a
	// choice, as we think it is better to work on surrounding struct
	// instead.
	if !got.CanInterface() {
		if ctx.booleanError {
			return booleanError
		}
		return ctx.CollectError(&Error{
			Message: "cannot smuggle unexported field",
			Summary: rawString("work on surrounding struct instead"),
		})
	}

	ret := s.function.Call([]reflect.Value{got})
	if len(ret) == 1 || ret[1].Bool() {
		newGot := ret[0]

		var newCtx Context
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

	if ctx.booleanError {
		return booleanError
	}

	err := Error{
		Message: "ran smuggle code with %% as argument",
	}

	if len(ret) > 2 {
		err.Summary = tdCodeResult{
			Value:  got,
			Reason: ret[2].String(),
		}
	} else {
		err.Summary = tdCodeResult{
			Value: got,
		}
	}

	return ctx.CollectError(&err)
}

func (s *tdSmuggle) String() string {
	return "Smuggle(" + s.function.Type().String() + ")"
}
