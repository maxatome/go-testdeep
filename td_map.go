package testdeep

import (
	"bytes"
	"fmt"
	"reflect"
)

type mapKind uint8

const (
	allMap mapKind = iota
	subMap
	superMap
)

type tdMap struct {
	Base
	expectedModel   reflect.Value
	expectedEntries []mapEntryInfo
	kind            mapKind
	isPtr           bool
}

var _ TestDeep = &tdMap{}

type mapEntryInfo struct {
	key      reflect.Value
	expected reflect.Value
}

type MapEntries map[interface{}]interface{}

func newMap(model interface{}, entries MapEntries, kind mapKind) *tdMap {
	vmodel := reflect.ValueOf(model)

	m := tdMap{
		Base: NewBase(4),
		kind: kind,
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Map {
			break
		}
		m.isPtr = true
		fallthrough

	case reflect.Map:
		m.expectedModel = vmodel
		m.populateExpectedEntries(entries)
		return &m
	}

	panic(fmt.Sprintf("usage: %s(MAP|&MAP, EXPECTED_ENTRIES)",
		m.GetLocation().Func))
}

func (m *tdMap) populateExpectedEntries(entries MapEntries) {
	expectedKeys := m.expectedModel.MapKeys()

	m.expectedEntries = make([]mapEntryInfo, 0, len(expectedKeys)+len(entries))
	checkedEntries := make(map[interface{}]bool, len(entries))

	keyType := m.expectedModel.Type().Key()
	valueType := m.expectedModel.Type().Elem()

	var entryInfo mapEntryInfo

	for key, expectedValue := range entries {
		vkey := reflect.ValueOf(key)
		if !vkey.Type().AssignableTo(keyType) {
			panic(fmt.Sprintf(
				"Expected key %s type mismatch: %s != model key type (%s)",
				toString(key),
				vkey.Type(),
				keyType))
		}

		entryInfo.expected = reflect.ValueOf(expectedValue)

		if _, ok := expectedValue.(TestDeep); !ok {
			if !entryInfo.expected.Type().AssignableTo(valueType) {
				panic(fmt.Sprintf(
					"Expected key %s value type mismatch: %s != model key type (%s)",
					toString(key),
					entryInfo.expected.Type(),
					valueType))
			}
		}

		entryInfo.key = vkey
		m.expectedEntries = append(m.expectedEntries, entryInfo)
		checkedEntries[vkey.Interface()] = true
	}

	// Check entries in model
	for _, vkey := range expectedKeys {
		entryInfo.expected = m.expectedModel.MapIndex(vkey)

		if checkedEntries[vkey.Interface()] {
			panic(fmt.Sprintf(
				"%s entry exists in both model & expectedEntries", toString(vkey)))
		}

		entryInfo.key = vkey
		m.expectedEntries = append(m.expectedEntries, entryInfo)
	}
}

func Map(model interface{}, entries MapEntries) TestDeep {
	return newMap(model, entries, allMap)
}

func SubMapOf(model interface{}, entries MapEntries) TestDeep {
	return newMap(model, entries, subMap)
}

func SuperMapOf(model interface{}, entries MapEntries) TestDeep {
	return newMap(model, entries, superMap)
}

func (m *tdMap) Match(ctx Context, got reflect.Value) (err *Error) {
	if m.isPtr {
		if got.Kind() != reflect.Ptr {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "type mismatch",
				Got:      rawString(got.Type().String()),
				Expected: rawString(m.expectedTypeStr()),
				Location: m.GetLocation(),
			}
		}
		got = got.Elem()
	}

	if got.Type() != m.expectedModel.Type() {
		if ctx.booleanError {
			return booleanError
		}
		var gotType rawString
		if m.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      gotType,
			Expected: rawString(m.expectedTypeStr()),
			Location: m.GetLocation(),
		}
	}

	var notFoundKeys []reflect.Value
	foundKeys := map[interface{}]bool{}

	for _, entryInfo := range m.expectedEntries {
		gotValue := got.MapIndex(entryInfo.key)
		if !gotValue.IsValid() {
			notFoundKeys = append(notFoundKeys, entryInfo.key)
			continue
		}

		err = deepValueEqual(ctx.AddDepth("["+toString(entryInfo.key)+"]"),
			got.MapIndex(entryInfo.key), entryInfo.expected)
		if err != nil {
			return err.SetLocationIfMissing(m)
		}
		foundKeys[entryInfo.key.Interface()] = true
	}

	const errorMessage = "comparing hash keys of %%"

	// For SuperMapOf we don't care about extra keys
	if m.kind == superMap {
		if len(notFoundKeys) == 0 {
			return nil
		}

		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context: ctx,
			Message: errorMessage,
			Summary: tdSetResult{
				Kind:    keysSetResult,
				Missing: notFoundKeys,
			},
			Location: m.GetLocation(),
		}
	}

	// No extra key to search, all got keys have been found
	if got.Len() == len(foundKeys) {
		if m.kind == subMap {
			return nil
		}
		// allMap

		if len(notFoundKeys) == 0 {
			return nil
		}

		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context: ctx,
			Message: errorMessage,
			Summary: tdSetResult{
				Kind:    keysSetResult,
				Missing: notFoundKeys,
			},
			Location: m.GetLocation(),
		}
	}

	if ctx.booleanError {
		return booleanError
	}

	// Retrieve extra keys
	res := tdSetResult{
		Kind:    keysSetResult,
		Missing: notFoundKeys,
		Extra:   make([]reflect.Value, 0, got.Len()-len(foundKeys)),
	}

	for _, vkey := range got.MapKeys() {
		if !foundKeys[vkey.Interface()] {
			res.Extra = append(res.Extra, vkey)
		}
	}

	return &Error{
		Context:  ctx,
		Message:  errorMessage,
		Summary:  res,
		Location: m.GetLocation(),
	}
}

func (m *tdMap) String() string {
	buf := &bytes.Buffer{}

	if m.kind != allMap {
		buf.WriteString(m.GetLocation().Func)
		buf.WriteByte('(')
	}

	buf.WriteString(m.expectedTypeStr())

	if len(m.expectedEntries) == 0 {
		buf.WriteString("{}")
	} else {
		buf.WriteString("{\n")

		for _, entryInfo := range m.expectedEntries {
			fmt.Fprintf(buf, "  %s: %s,\n",
				toString(entryInfo.key),
				toString(entryInfo.expected))
		}

		buf.WriteByte('}')
	}

	if m.kind != allMap {
		buf.WriteByte(')')
	}

	return buf.String()
}

func (m *tdMap) expectedTypeStr() string {
	if m.isPtr {
		return "*" + m.expectedModel.Type().String()
	}
	return m.expectedModel.Type().String()
}
