package testdeep

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type tdArray struct {
	TestDeepBase
	expectedModel   reflect.Value
	expectedEntries entryInfoSlice
	isPtr           bool
}

var _ TestDeep = &tdArray{}

type entryInfo struct {
	index    int
	expected reflect.Value
}

type entryInfoSlice []entryInfo

func (e entryInfoSlice) Len() int           { return len(e) }
func (e entryInfoSlice) Less(i, j int) bool { return e[i].index < e[j].index }
func (e entryInfoSlice) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type ArrayEntries map[int]interface{}

func Array(model interface{}, entries ArrayEntries) TestDeep {
	vmodel := reflect.ValueOf(model)

	a := tdArray{
		TestDeepBase: NewTestDeepBase(3),
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Array {
			break
		}
		a.isPtr = true
		fallthrough

	case reflect.Array:
		a.expectedModel = vmodel
		a.populateExpectedEntries(entries)
		return &a
	}

	panic("usage: Array(ARRAY|&ARRAY, EXPECTED_ENTRIES)")
}

func Slice(model interface{}, entries ArrayEntries) TestDeep {
	vmodel := reflect.ValueOf(model)

	a := tdArray{
		TestDeepBase: NewTestDeepBase(3),
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Slice {
			break
		}
		a.isPtr = true
		fallthrough

	case reflect.Slice:
		a.expectedModel = vmodel
		a.populateExpectedEntries(entries)
		return &a
	}

	panic("usage: Slice(SLICE|&SLICE, EXPECTED_ENTRIES)")
}

func (a *tdArray) populateExpectedEntries(entries ArrayEntries) {
	var length int
	if a.expectedModel.Kind() == reflect.Array {
		length = a.expectedModel.Len()
	} else {
		length = -1
	}

	a.expectedEntries = make([]entryInfo, 0, len(entries))
	checkedIndexes := make(map[int]bool, len(entries))

	elemType := a.expectedModel.Type().Elem()
	for index, expectedValue := range entries {
		if length >= 0 && index >= length {
			panic(fmt.Sprintf(
				"array length is %d, so cannot have #%d expected index",
				length,
				index))
		}

		vexpectedValue := reflect.ValueOf(expectedValue)

		if _, ok := expectedValue.(TestDeep); !ok {
			if !vexpectedValue.Type().AssignableTo(elemType) {
				arrayOrSlice := "array"
				if length < 0 {
					arrayOrSlice = "slice"
				}
				panic(fmt.Sprintf(
					"type %s of #%d expected value differs from %s content (%s)",
					vexpectedValue.Type(),
					index,
					arrayOrSlice,
					elemType))
			}
		}

		a.expectedEntries = append(a.expectedEntries, entryInfo{
			index:    index,
			expected: vexpectedValue,
		})
		checkedIndexes[index] = true
	}

	// Check initialized entries in model
	zero := reflect.Zero(elemType).Interface()
	for index := a.expectedModel.Len() - 1; index >= 0; index-- {
		ventry := a.expectedModel.Index(index)

		// If non-zero entry
		if !reflect.DeepEqual(zero, ventry.Interface()) {
			if checkedIndexes[index] {
				panic(fmt.Sprintf(
					"non zero #%d entry in model already exists in expectedEntries",
					index))
			}

			a.expectedEntries = append(a.expectedEntries, entryInfo{
				index:    index,
				expected: ventry,
			})
		}
	}

	sort.Sort(a.expectedEntries)
}

func (a *tdArray) Match(ctx Context, got reflect.Value) (err *Error) {
	if a.isPtr {
		if got.Kind() != reflect.Ptr {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "type mismatch",
				Got:      rawString(got.Type().String()),
				Expected: rawString(a.expectedTypeStr()),
				Location: a.GetLocation(),
			}
		}
		got = got.Elem()
	}

	if got.Type() != a.expectedModel.Type() {
		if ctx.booleanError {
			return booleanError
		}
		var gotType rawString
		if a.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      gotType,
			Expected: rawString(a.expectedTypeStr()),
			Location: a.GetLocation(),
		}
	}

	gotLen := got.Len()
	for _, entryInfo := range a.expectedEntries {
		curCtx := ctx.AddArrayIndex(entryInfo.index)

		if entryInfo.index >= gotLen {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  curCtx,
				Message:  "expected value out of range",
				Got:      rawString("<non-existent value>"),
				Expected: entryInfo.expected,
				Location: a.GetLocation(),
			}
		}

		err = deepValueEqual(got.Index(entryInfo.index), entryInfo.expected, curCtx)
		if err != nil {
			return err.SetLocationIfMissing(a)
		}
	}
	return nil
}

func (a *tdArray) String() string {
	buf := bytes.NewBufferString(ternStr(a.expectedModel.Kind() == reflect.Array,
		"Array(", "Slice("))

	buf.WriteString(a.expectedTypeStr())

	if len(a.expectedEntries) == 0 {
		buf.WriteString("{})")
	} else {
		buf.WriteString("{\n")

		for _, entryInfo := range a.expectedEntries {
			fmt.Fprintf(buf, "  %d: %s\n",
				entryInfo.index, toString(entryInfo.expected))
		}

		buf.WriteString("})")
	}
	return buf.String()
}

func (a *tdArray) expectedTypeStr() string {
	if a.isPtr {
		return "*" + a.expectedModel.Type().String()
	}
	return a.expectedModel.Type().String()
}
