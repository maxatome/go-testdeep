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
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "nil pointer",
				Got:      rawString("nil " + got.Type().String()),
				Expected: rawString("Slice OR Array OR *Slice OR *Array"),
				Location: s.GetLocation(),
			}
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
					if ctx.booleanError {
						return booleanError
					}
					res.Missing = missingItems
				}
			}

			if len(foundGotIdxes) < gotLen && s.kind != superSet {
				if ctx.booleanError {
					return booleanError
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
			if ctx.booleanError {
				return booleanError
			}
			res.Extra = foundItems
		}

		if res.IsEmpty() {
			return nil
		}
		return &Error{
			Context:  ctx,
			Message:  "comparing %% as a " + s.GetLocation().Func,
			Summary:  res,
			Location: s.GetLocation(),
		}
	}

	if ctx.booleanError {
		return booleanError
	}

	var gotStr rawString
	if got.IsValid() {
		gotStr = rawString(got.Type().String())
	} else {
		gotStr = "nil"
	}

	return &Error{
		Context:  ctx,
		Message:  "bad type",
		Got:      gotStr,
		Expected: rawString("Slice OR Array OR *Slice OR *Array"),
		Location: s.GetLocation(),
	}
}

func (s *tdSetBase) String() string {
	return sliceToBuffer(
		bytes.NewBufferString(s.GetLocation().Func), s.expectedItems).String()
}
