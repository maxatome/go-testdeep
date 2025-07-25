// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/flat"
	"github.com/maxatome/go-testdeep/internal/util"
)

type setKind uint8

const (
	allSet setKind = iota
	subSet
	superSet
	noneSet
)

type tdSetBase struct {
	baseOKNil
	kind       setKind
	ignoreDups bool

	expectedItems []reflect.Value
}

func newSetBase(kind setKind, ignoreDups bool, expectedItems []any) *tdSetBase {
	return &tdSetBase{
		baseOKNil:     newBaseOKNil(4),
		kind:          kind,
		ignoreDups:    ignoreDups,
		expectedItems: flat.Values(expectedItems),
	}
}

func (s *tdSetBase) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	switch got.Kind() {
	case reflect.Ptr:
		gotElem := got.Elem()
		if !gotElem.IsValid() {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			return ctx.CollectError(ctxerr.NilPointer(got, "non-nil *slice OR *array"))
		}

		if gotElem.Kind() != reflect.Array && gotElem.Kind() != reflect.Slice {
			break
		}
		got = gotElem
		fallthrough

	case reflect.Array, reflect.Slice:
		var (
			gotLen = got.Len()

			foundItems    []reflect.Value
			missingItems  []reflect.Value
			foundGotIdxes = map[int]bool{}
		)

		for _, expected := range s.expectedItems {
			found := false

			for idx := 0; len(foundGotIdxes) < gotLen && idx < gotLen; idx++ {
				if foundGotIdxes[idx] {
					continue
				}

				ok, err := deepValueEqualFinalOK(ctx, got.Index(idx), expected)
				if err != nil { // user error, stop asap
					return err
				}
				if ok {
					foundItems = append(foundItems, expected)

					foundGotIdxes[idx] = true
					found = true

					if !s.ignoreDups {
						break
					}
				}
			}

			if !found {
				missingItems = append(missingItems, expected)
			}
		}

		res := tdSetResult{
			Kind: itemsSetResult,
			Sort: true,
		}

		if s.kind != noneSet {
			if s.kind != subSet {
				// In Set* cases with missing items, try a second pass. Perhaps
				// an already matching got item, matches another expected item?
				if s.ignoreDups && len(missingItems) > 0 {
					var newMissingItems []reflect.Value

				nextExpected:
					for _, expected := range missingItems {
						for idxGot := range foundGotIdxes {
							ok, _ := deepValueEqualFinalOK(ctx, got.Index(idxGot), expected)
							if ok {
								continue nextExpected
							}
						}

						newMissingItems = append(newMissingItems, expected)
					}

					missingItems = newMissingItems
				}

				if len(missingItems) > 0 {
					if ctx.BooleanError {
						return ctxerr.BooleanError
					}
					res.Missing = missingItems
				}
			}

			if len(foundGotIdxes) < gotLen && s.kind != superSet {
				if ctx.BooleanError {
					return ctxerr.BooleanError
				}
				notFoundRemain := gotLen - len(foundGotIdxes)
				res.Extra = make([]reflect.Value, 0, notFoundRemain)
				for idx := 0; notFoundRemain > 0; idx++ {
					if !foundGotIdxes[idx] {
						res.Extra = append(res.Extra, got.Index(idx))
						notFoundRemain--
					}
				}
			}
		} else if len(foundItems) > 0 {
			if ctx.BooleanError {
				return ctxerr.BooleanError
			}
			res.Extra = foundItems
		}

		if res.IsEmpty() {
			return nil
		}
		return ctx.CollectError(&ctxerr.Error{
			Message: "comparing %% as a " + s.GetLocation().Func,
			Summary: res.Summary(),
		})
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	return ctx.CollectError(ctxerr.BadKind(got, "slice OR array OR *slice OR *array"))
}

func (s *tdSetBase) String() string {
	var b strings.Builder
	b.WriteString(s.GetLocation().Func)
	return util.SliceToString(&b, s.expectedItems).String()
}

func (s *tdSetBase) TypeBehind() reflect.Type {
	typ := uniqTypeBehindSlice(s.expectedItems)
	if typ == nil {
		return nil
	}
	return reflect.SliceOf(typ)
}
