// Copyright (c) 2024, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/maxatome/go-testdeep/internal/compare"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/visited"
)

const sortUsage = "(SORT_FUNC|int|string|[]string, TESTDEEP_OPERATOR|EXPECTED_VALUE)"

type tdSort struct {
	tdSmugglerBase
	how      any
	mkSortFn func(reflect.Type) (reflect.Value, error)
}

var _ TestDeep = &tdSort{}

func mkSortAsc(typ reflect.Type) (reflect.Value, error) {
	v := visited.NewVisited()
	return reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{typ, typ}, []reflect.Type{types.Bool}, false),
		func(args []reflect.Value) []reflect.Value {
			less := compare.Compare(v, args[0], args[1]) < 0
			return []reflect.Value{reflect.ValueOf(less)}
		}), nil
}

func mkSortDesc(typ reflect.Type) (reflect.Value, error) {
	v := visited.NewVisited()
	return reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{typ, typ}, []reflect.Type{types.Bool}, false),
		func(args []reflect.Value) []reflect.Value {
			less := compare.Compare(v, args[1], args[0]) < 0
			return []reflect.Value{reflect.ValueOf(less)}
		}), nil
}

func derefStruct(v reflect.Value) (dv reflect.Value, isStruct, isNil bool) {
	for {
		switch v.Kind() {
		case reflect.Interface:
			if v.IsNil() {
				return v, false, true
			}
			v = v.Elem()
		case reflect.Ptr:
			if v.IsNil() {
				for t := v.Type(); ; {
					switch t.Kind() {
					case reflect.Struct:
						return v, true, true
					case reflect.Ptr:
						t = t.Elem()
					default:
						return v, false, true
					}
				}
			}
			v = v.Elem()
		case reflect.Struct:
			return v, true, false
		default:
			return v, false, false
		}
	}
}

func mkSortStructFields(typ reflect.Type, fields []string) (reflect.Value, error) {
	v := visited.NewVisited()
	// Could be optimized if typ is a (Ptr)*Struct
	return reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{typ, typ}, []reflect.Type{types.Bool}, false),
		func(args []reflect.Value) []reflect.Value {
			a, aOK, aIsNil := derefStruct(args[0])
			b, bOK, bIsNil := derefStruct(args[1])
			if !aOK || !bOK {
				// non-struct is greater than any struct value. If both are
				// not structs, then use Compare
				if aOK || bOK {
					return []reflect.Value{reflect.ValueOf(aOK)}
				}
				less := compare.Compare(v, args[0], args[1]) < 0
				return []reflect.Value{reflect.ValueOf(less)}
			}
			if aIsNil || bIsNil {
				less := compare.Compare(v, args[0], args[1]) < 0
				return []reflect.Value{reflect.ValueOf(less)}
			}
			for _, name := range fields {
				asc := true
				if strings.HasPrefix(name, "-") {
					name = name[1:]
					asc = false
				}
				fa, fb := a.FieldByName(name), b.FieldByName(name)
				if fa.IsValid() && fb.IsValid() {
					cmp := compare.Compare(v, fa, fb)
					if cmp == 0 {
						continue
					}
					return []reflect.Value{reflect.ValueOf(cmp < 0 == asc)}
				}
				// at least a nonexistent field, use Compare
				break
			}
			less := compare.Compare(v, a, b) < 0
			return []reflect.Value{reflect.ValueOf(less)}
		}), nil
}

func sortFunc(how any) (func(reflect.Type) (reflect.Value, error), error) {
	switch v := how.(type) {
	case int:
		if v < 0 {
			return mkSortDesc, nil
		}
		return mkSortAsc, nil
	case float64: // to be used in JSON, SubJSONOf & SuperJSONOf
		if v < 0 {
			return mkSortDesc, nil
		}
		return mkSortAsc, nil
	case string: // one struct field
		return func(typ reflect.Type) (reflect.Value, error) {
			return mkSortStructFields(typ, []string{v})
		}, nil
	case []string: // struct fields list
		return func(typ reflect.Type) (reflect.Value, error) {
			return mkSortStructFields(typ, v)
		}, nil
	default:
		vv := reflect.ValueOf(how)
		if vv.Kind() != reflect.Func {
			return nil, fmt.Errorf("SORT_FUNC must be an int, string, []string or func(T, T) bool")
		}
		ft := vv.Type()
		if ft.IsVariadic() || ft.NumIn() != 2 || ft.In(0) != ft.In(1) ||
			ft.NumOut() != 1 || ft.Out(0) != types.Bool {
			return nil, fmt.Errorf("SORT_FUNC must match func(T, T) bool signature, not %T", v)
		}
		return func(typ reflect.Type) (reflect.Value, error) {
			if !typ.AssignableTo(ft.In(0)) {
				return reflect.Value{}, fmt.Errorf("%s is not assignable to %s", typ, ft.In(0))
			}
			return vv, nil
		}, nil
	}
}

// summary(Sort): sorts a slice or an array before comparing its content
// input(Sort): array,slice,ptr(ptr on array/slice)

// Sort is a smuggler operator. XXX.
func Sort(how any, expectedValue any) TestDeep {
	s := tdSort{}
	s.tdSmugglerBase = newSmugglerBase(expectedValue, 0)
	if !s.isTestDeeper {
		s.expectedValue = reflect.ValueOf(expectedValue)
	}

	var err error
	s.mkSortFn, err = sortFunc(how)
	if err != nil {
		s.err = ctxerr.OpBad("Sort", "usage: Sort%s, %s", sortUsage, err)
	} else if !s.isTestDeeper {
		switch s.expectedValue.Kind() {
		case reflect.Slice, reflect.Array:
		default:
			s.err = ctxerr.OpBad("Sort",
				"usage: Sort%s, EXPECTED_VALUE must be a slice not a %s",
				sortUsage, types.KindType(s.expectedValue))
		}
	}
	s.how = how
	return &s
}

func (s *tdSort) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if s.err != nil {
		return ctx.CollectError(s.err)
	}

	if rErr := grepResolvePtr(ctx, &got); rErr != nil {
		return rErr
	}

	switch got.Kind() {
	case reflect.Slice, reflect.Array:
		const sorted = "<sorted>"

		l := got.Len()
		if l <= 1 {
			return deepValueEqual(ctx.AddCustomLevel(sorted), got, s.expectedValue)
		}
		itemType := got.Type().Elem()
		fn, err := s.mkSortFn(itemType)
		if err != nil {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(&ctxerr.Error{
				Message:  "incompatible parameter type",
				Got:      types.RawString(itemType.String()),
				Expected: types.RawString(fn.Type().In(0).String()),
			})
		}

		out := reflect.MakeSlice(reflect.SliceOf(itemType), 0, l)
		for idx := 0; idx < l; idx++ {
			out = reflect.Append(out, got.Index(idx))
		}

		sort.SliceStable(out.Interface(), func(i, j int) bool {
			return fn.Call([]reflect.Value{out.Index(i), out.Index(j)})[0].Bool()
		})
		return deepValueEqual(ctx.AddCustomLevel(sorted), out, s.expectedValue)
	}

	return grepBadKind(ctx, got)
}

func (s *tdSort) String() string {
	if s.err != nil {
		return s.stringError()
	}
	how, typ := s.how, reflect.TypeOf(s.how)
	if typ.Kind() == reflect.Func {
		how = typ.String()
	}
	return fmt.Sprintf("Sort(%v)", how)
}
