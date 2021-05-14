// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/util"
)

type tdJSONPointer struct {
	tdSmugglerBase
	pointer string
}

var _ TestDeep = &tdJSONPointer{}

// summary(JSONPointer): compares against JSON representation using a
// JSON pointer
// input(JSONPointer): nil,bool,str,int,float,array,slice,map,struct,ptr

// JSONPointer is a smuggler operator. It takes the JSON
// representation of data, gets the value corresponding to the JSON
// pointer "pointer" (as RFC 6901 specifies it) and compares it to
// "expectedValue".
//
// JSONPointer does its best to convert back the JSON pointed data to
// the type of "expectedValue" or to the type behind the
// "expectedValue" operator, if it is an operator. Allowing to do
// things like:
//
//   type Item struct {
//     Val  int   `json:"val"`
//     Next *Item `json:"next"`
//   }
//   got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}
//
//   td.Cmp(t, got, td.JSONPointer("/next/next", Item{Val: 3}))
//   td.Cmp(t, got, td.JSONPointer("/next/next", &Item{Val: 3}))
//   td.Cmp(t,
//     got,
//     td.JSONPointer("/next/next",
//       td.Struct(Item{}, td.StructFields{"Val": td.Gte(3)})),
//   )
//
// It does this conversion only if the expected type is a struct, a
// struct pointer or implements the encoding/json.Unmarshaler
// interface. In the case the conversion does not occur, the Lax mode
// is automatically enabled to simplify numeric tests.
//
//   got := map[string]int64{"zzz": 42} // 42 is int64 here
//   td.Cmp(t, got, td.JSONPointer("/zzz", 42))
//   td.Cmp(t, got, td.JSONPointer("/zzz", td.Between(40, 45)))
//
// Note that as any TestDeep operator can be used as "expectedValue",
// JSON operator works out of the box:
//
//   got := json.RawMessage(`{"foo":{"bar": {"zip": true}}}`)
//   td.Cmp(t, got, td.JSONPointer("/foo/bar", td.JSON(`{"zip": true}`)))
//
// It can be used with structs lacking json tags. In this case, fields
// names have to be used in JSON pointer:
//
//   type Item struct {
//     Val  int
//     Next *Item
//   }
//   got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}
//
//   td.Cmp(t, got, td.JSONPointer("/Next/Next", Item{Val: 3}))
//
// Contrary to Smuggle operator and its fields-path feature, only
// public fields can be followed, as private ones are never (un)marshalled.
//
// There is no JSONHas nor JSONHasnt operators to only check a JSON
// pointer exists or not, but they can easily be emulated:
//
//   JSONHas := func(pointer string) td.TestDeep {
//     return td.JSONPointer(pointer, td.Ignore())
//   }
//
//   JSONHasnt := func(pointer string) td.TestDeep {
//     return td.Not(td.JSONPointer(pointer, td.Ignore()))
//   }
//
// TypeBehind method always returns nil as the expected type cannot be
// guessed from a JSON pointer.
func JSONPointer(pointer string, expectedValue interface{}) TestDeep {
	if !strings.HasPrefix(pointer, "/") && pointer != "" {
		f := dark.GetFatalizer()
		f.Helper()
		dark.Fatal(f, color.Bad("JSONPointer(): bad JSON pointer %s", pointer))
	}

	p := tdJSONPointer{
		tdSmugglerBase: newSmugglerBase(expectedValue),
		pointer:        pointer,
	}
	if !p.isTestDeeper {
		p.expectedValue = reflect.ValueOf(expectedValue)
	}
	return &p
}

func (p *tdJSONPointer) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	vgot, eErr := jsonify(ctx, got)
	if eErr != nil {
		return ctx.CollectError(eErr)
	}

	newGot, err := util.JSONPointer(vgot, p.pointer)
	if err != nil {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		pErr := err.(*util.JSONPointerError)
		ctx = jsonPointerContext(ctx, pErr.Pointer)
		return ctx.CollectError(&ctxerr.Error{
			Message: "cannot retrieve value via JSON pointer",
			Summary: ctxerr.NewSummary(pErr.Type),
		})
	}

	ctx = jsonPointerContext(ctx, p.pointer)

	// Here, newGot type is either a bool, float64, string,
	// []interface{} or a map[string]interface{}

	// Check if we have to transform the new got into something
	// compatible with the type of expected
	if expectedType := p.internalTypeBehind(); expectedType != nil &&
		(expectedType.Implements(types.JsonUnmarshaler) ||
			reflect.PtrTo(expectedType).Implements(types.JsonUnmarshaler) ||
			types.IsStruct(expectedType)) {
		b, _ := json.Marshal(newGot) // No error can occur here

		got = reflect.New(expectedType)
		err := json.Unmarshal(b, got.Interface())
		if err != nil {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(&ctxerr.Error{
				Message: fmt.Sprintf(
					"an error occurred while unmarshalling JSON into %s", expectedType),
				Summary: ctxerr.NewSummary(err.Error()),
			})
		}
		got = got.Elem()
	} else {
		ctx.BeLax = true
		got = reflect.ValueOf(newGot)
	}

	return deepValueEqual(ctx, got, p.expectedValue)
}

func (p *tdJSONPointer) String() string {
	var expected string
	switch {
	case p.isTestDeeper:
		expected = p.expectedValue.Interface().(TestDeep).String()
	case p.expectedValue.IsValid():
		expected = util.ToString(p.expectedValue.Interface())
	default:
		expected = "nil"
	}
	return fmt.Sprintf("JSONPointer(%s, %s)", p.pointer, expected)
}

func (p *tdJSONPointer) internalTypeBehind() reflect.Type {
	if p.isTestDeeper {
		return p.expectedValue.Interface().(TestDeep).TypeBehind()
	}
	if p.expectedValue.IsValid() {
		return p.expectedValue.Type()
	}
	return nil
}

func (p *tdJSONPointer) HandleInvalid() bool {
	return true
}

func jsonPointerContext(ctx ctxerr.Context, pointer string) ctxerr.Context {
	return ctx.AddCustomLevel(".JSONPointer<" + pointer + ">")
}
