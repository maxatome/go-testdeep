package testdeep

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type tdStruct struct {
	TestDeepBase
	expectedModel  reflect.Value
	expectedFields fieldInfoSlice
	isPtr          bool
}

var _ TestDeep = &tdStruct{}

type fieldInfo struct {
	name     string
	expected reflect.Value
	index    []int
}

type fieldInfoSlice []fieldInfo

func (e fieldInfoSlice) Len() int           { return len(e) }
func (e fieldInfoSlice) Less(i, j int) bool { return e[i].name < e[j].name }
func (e fieldInfoSlice) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type StructFields map[string]interface{}

func newStruct(model interface{}) *tdStruct {
	vmodel := reflect.ValueOf(model)

	st := tdStruct{
		TestDeepBase: NewTestDeepBase(4),
	}

	switch vmodel.Kind() {
	case reflect.Ptr:
		vmodel = vmodel.Elem()
		if vmodel.Kind() != reflect.Struct {
			break
		}
		st.isPtr = true
		fallthrough

	case reflect.Struct:
		st.expectedModel = vmodel
		return &st
	}

	panic("usage: Struct(STRUCT|&STRUCT, EXPECTED_FIELDS)")
}

func Struct(model interface{}, expectedFields StructFields) TestDeep {
	st := newStruct(model)

	st.expectedFields = make([]fieldInfo, 0, len(expectedFields))
	checkedFields := make(map[string]bool, len(expectedFields))

	vmodel := st.expectedModel

	// Check that all given fields are available in model
	stType := vmodel.Type()
	for fieldName, expectedValue := range expectedFields {
		field, found := stType.FieldByName(fieldName)
		if !found {
			panic(fmt.Sprintf("struct %s has no field `%s'",
				vmodel.Type(), fieldName))
		}

		vexpectedValue := reflect.ValueOf(expectedValue)

		if _, ok := expectedValue.(TestDeep); !ok {
			if !vexpectedValue.Type().AssignableTo(field.Type) {
				panic(fmt.Sprintf(
					"type %s of field expected value %s differs from struct one (%s)",
					vexpectedValue.Type(),
					fieldName,
					field.Type))
			}
		}

		st.expectedFields = append(st.expectedFields, fieldInfo{
			name:     fieldName,
			expected: vexpectedValue,
			index:    field.Index,
		})
		checkedFields[fieldName] = true
	}

	// Get all field names
	var allFields []string
	stType.FieldByNameFunc(func(fieldName string) bool {
		allFields = append(allFields, fieldName)
		return false
	})

	// Check initialized fields in model
	for _, fieldName := range allFields {
		field, _ := stType.FieldByName(fieldName)
		if field.Anonymous {
			continue
		}

		vfield := vmodel.FieldByIndex(field.Index)

		// If non-zero field
		if !reflect.DeepEqual(reflect.Zero(field.Type).Interface(),
			vfield.Interface()) {
			if checkedFields[fieldName] {
				panic(fmt.Sprintf(
					"non zero field %s in model already exists in expectedFields",
					fieldName))
			}

			st.expectedFields = append(st.expectedFields, fieldInfo{
				name:     fieldName,
				expected: vfield,
				index:    field.Index,
			})
		}
	}

	sort.Sort(st.expectedFields)

	return st
}

func (s *tdStruct) Match(ctx Context, got reflect.Value) (err *Error) {
	if s.isPtr {
		if got.Kind() != reflect.Ptr {
			if ctx.booleanError {
				return booleanError
			}
			return &Error{
				Context:  ctx,
				Message:  "type mismatch",
				Got:      rawString(got.Type().String()),
				Expected: rawString(s.expectedTypeStr()),
				Location: s.GetLocation(),
			}
		}
		got = got.Elem()
	}

	if got.Type() != s.expectedModel.Type() {
		if ctx.booleanError {
			return booleanError
		}
		var gotType rawString
		if s.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      gotType,
			Expected: rawString(s.expectedTypeStr()),
			Location: s.GetLocation(),
		}
	}

	for _, fieldInfo := range s.expectedFields {
		err = deepValueEqual(got.FieldByIndex(fieldInfo.index),
			fieldInfo.expected,
			ctx.AddDepth("."+fieldInfo.name))
		if err != nil {
			return err.SetLocationIfMissing(s)
		}
	}
	return nil
}

func (s *tdStruct) String() string {
	buf := bytes.NewBufferString("Struct(")

	if s.isPtr {
		buf.WriteByte('*')
	}

	buf.WriteString(s.expectedModel.Type().String())

	if len(s.expectedFields) == 0 {
		buf.WriteString("{})")
	} else {
		buf.WriteString("{\n")

		for _, fieldInfo := range s.expectedFields {
			fmt.Fprintf(buf, "  %s: %s\n",
				fieldInfo.name, toString(fieldInfo.expected))
		}

		buf.WriteString("})")
	}

	return buf.String()
}

func (s *tdStruct) expectedTypeStr() string {
	if s.isPtr {
		return "*" + s.expectedModel.Type().String()
	}
	return s.expectedModel.Type().String()
}
