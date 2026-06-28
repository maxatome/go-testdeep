// Copyright (c) 2020-2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/util"
)

type tdJSONPointer struct {
	tdSmugglerBase
	pointer string
	options jsonv2Options
}

var _ TestDeep = &tdJSONPointer{}

// summary(JSONPointer): compares against JSON representation using a
// JSON pointer
// input(JSONPointer): nil,bool,str,int,float,array,slice,map,struct,ptr

// JSONPointer is a smuggler operator. It takes the JSON
// representation of data, gets the value corresponding to the JSON
// pointer ptr (as [RFC 6901] specifies it) and compares it to
// expectedValue.
//
// By default, [encoding/json] is used to marshal and unmarshal
// data. But if opts are passed using [encoding/json/v2.Options]
// values, then [encoding/json/v2] is used (go1.27 required). Note
// that multiple opts are automatically merged using
// [encoding/json/v2.JoinOptions] before use. So the following calls
// are equivalent:
//
//	td.Cmp(t, got, td.JSONPointer("/a/b", 123))
//	td.Cmp(t, got, td.JSONPointer("/a/b", 123, json.DefaultOptionsV1()))
//
// [Lax] mode is automatically enabled to simplify numeric tests.
//
// JSONPointer does its best to convert back the JSON pointed data to
// the type of expectedValue or to the type behind the
// expectedValue operator, if it is an operator. Allowing to do
// things like:
//
//	type Item struct {
//	  Val  int   `json:"val"`
//	  Next *Item `json:"next"`
//	}
//	got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}
//
//	td.Cmp(t, got, td.JSONPointer("/next/next", Item{Val: 3}))
//	td.Cmp(t, got, td.JSONPointer("/next/next", &Item{Val: 3}))
//	td.Cmp(t,
//	  got,
//	  td.JSONPointer("/next/next",
//	    td.Struct(Item{}, td.StructFields{"Val": td.Gte(3)})),
//	)
//
//	got := map[string]int64{"zzz": 42} // 42 is int64 here
//	td.Cmp(t, got, td.JSONPointer("/zzz", 42))
//	td.Cmp(t, got, td.JSONPointer("/zzz", td.Between(40, 45)))
//
// Of course, it does this conversion only if the expected type can be
// guessed. In the case the conversion cannot occur, data is compared
// as is, in its freshly unmarshaled JSON form (so as bool, float64,
// string, []any, map[string]any or simply nil).
//
// Note that as any [TestDeep] operator can be used as expectedValue,
// [JSON] operator works out of the box:
//
//	got := json.RawMessage(`{"foo":{"bar": {"zip": true}}}`)
//	td.Cmp(t, got, td.JSONPointer("/foo/bar", td.JSON(`{"zip": true}`)))
//
// It can be used with structs lacking json tags. In this case, fields
// names have to be used in JSON pointer:
//
//	type Item struct {
//	  Val  int
//	  Next *Item
//	}
//	got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}
//
//	td.Cmp(t, got, td.JSONPointer("/Next/Next", Item{Val: 3}))
//
// Contrary to [Smuggle] operator and its fields-path feature, only
// public fields can be followed, as private ones are never (un)marshaled.
//
// There is no JSONHas nor JSONHasnt operators to only check a JSON
// pointer exists or not, but they can easily be emulated:
//
//	JSONHas := func(pointer string) td.TestDeep {
//	  return td.JSONPointer(pointer, td.Ignore())
//	}
//
//	JSONHasnt := func(pointer string) td.TestDeep {
//	  return td.Not(td.JSONPointer(pointer, td.Ignore()))
//	}
//
// TypeBehind method always returns nil as the expected type cannot be
// guessed from a JSON pointer.
//
// See also [JSON], [SubJSONOf], [SuperJSONOf], [Smuggle] and [Flatten].
//
// [RFC 6901]: https://tools.ietf.org/html/rfc6901
func JSONPointer(ptr string, expectedValue any, opts ...jsonv2Options) TestDeep {
	p := tdJSONPointer{
		tdSmugglerBase: newSmugglerBase(expectedValue),
		pointer:        ptr,
	}
	for _, o := range opts {
		p.options = joinOptions(p.options, o)
	}

	if !strings.HasPrefix(ptr, "/") && ptr != "" {
		p.err = ctxerr.OpBad("JSONPointer", "bad JSON pointer %q", ptr)
		return &p
	}

	if !p.isTestDeeper {
		p.expectedValue = reflect.ValueOf(expectedValue)
	}
	return &p
}

func (p *tdJSONPointer) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if p.err != nil {
		return ctx.CollectError(p.err)
	}

	vgot, eErr := jsonify(ctx, got, p.options)
	if eErr != nil {
		return ctx.CollectError(eErr)
	}

	vgot, err := util.JSONPointer(vgot, p.pointer)
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

	// Here, vgot type is either a bool, float64, string,
	// []any, a map[string]any or simply nil

	ctx = jsonPointerContext(ctx, p.pointer)
	ctx.BeLax = true

	return p.jsonValueEqual(ctx, vgot, p.options)
}

func (p *tdJSONPointer) String() string {
	if p.err != nil {
		return p.stringError()
	}

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

func (p *tdJSONPointer) HandleInvalid() bool {
	return true
}

func jsonPointerContext(ctx ctxerr.Context, pointer string) ctxerr.Context {
	return ctx.AddCustomLevel(".JSONPointer<" + pointer + ">")
}
