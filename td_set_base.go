// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"reflect"
)

type setKind uint8

const (
	allSet setKind = iota
	subSet
	superSet
	noneSet
)

type tdSetBase struct {
	BaseOKNil
	kind       setKind
	ignoreDups bool

	expectedItems []reflect.Value
}

func newSetBase(kind setKind, ignoreDups bool) tdSetBase {
	return tdSetBase{
		BaseOKNil:  NewBaseOKNil(4),
		kind:       kind,
		ignoreDups: ignoreDups,
	}
}

func (s *tdSetBase) Add(items ...interface{}) {
	for _, item := range items {
		s.expectedItems = append(s.expectedItems, reflect.ValueOf(item))
	}
}

func (s *tdSetBase) Match(ctx Context, got reflect.Value) *Error {
	switch got.Kind() {
	case reflect.Ptr:
		gotElem := got.Elem()
		if !gotElem.IsValid() {
			if ctx.BooleanError {
				return BooleanError
			}
			return ctx.CollectError(&Error{
				Message:  "nil pointer",
				Got:      rawString("nil " + got.Type().String()),
				Expected: rawString("Slice OR Array OR *Slice OR *Array"),
			})
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

				if deepValueEqualOK(got.Index(idx), expected) {
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
							if deepValueEqualOK(got.Index(idxGot), expected) {
								continue nextExpected
							}
						}

						newMissingItems = append(newMissingItems, expected)
					}

					missingItems = newMissingItems
				}

				if len(missingItems) > 0 {
					if ctx.BooleanError {
						return BooleanError
					}
					res.Missing = missingItems
				}
			}

			if len(foundGotIdxes) < gotLen && s.kind != superSet {
				if ctx.BooleanError {
					return BooleanError
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
				return BooleanError
			}
			res.Extra = foundItems
		}

		if res.IsEmpty() {
			return nil
		}
		return ctx.CollectError(&Error{
			Message: "comparing %% as a " + s.GetLocation().Func,
			Summary: res,
		})
	}

	if ctx.BooleanError {
		return BooleanError
	}

	var gotStr rawString
	if got.IsValid() {
		gotStr = rawString(got.Type().String())
	} else {
		gotStr = "nil"
	}

	return ctx.CollectError(&Error{
		Message:  "bad type",
		Got:      gotStr,
		Expected: rawString("Slice OR Array OR *Slice OR *Array"),
	})
}

func (s *tdSetBase) String() string {
	return sliceToBuffer(
		bytes.NewBufferString(s.GetLocation().Func), s.expectedItems).String()
}
