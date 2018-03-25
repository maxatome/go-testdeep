package testdeep

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func toString(val interface{}) string {
	if val == nil {
		return "nil"
	}

	switch tval := val.(type) {
	case testDeepStringer:
		return tval.String()

	case reflect.Value:
		var ok bool
		val, ok = getInterface(tval)
		if ok {
			if val == nil {
				return "nil"
			}
			if tdVal, ok := val.(TestDeep); ok {
				return tdVal.String()
			}
		}
	}

	return strings.TrimRight(spew.Sdump(val), "\n")
}

func indentString(str string, indent string) string {
	return strings.Replace(str, "\n", "\n"+indent, -1)
}

func sliceToBuffer(buf *bytes.Buffer, items []reflect.Value) *bytes.Buffer {
	buf.WriteByte('(')
	if len(items) < 2 {
		if len(items) > 0 {
			buf.WriteString(toString(items[0]))
		}
	} else {
		begLine := bytes.LastIndexByte(buf.Bytes(), '\n') + 1

		prefix := strings.Repeat(" ", buf.Len()-begLine)

		for idx, item := range items {
			if idx != 0 {
				buf.WriteString(prefix)
			}
			buf.WriteString(toString(item))
			buf.WriteString(",\n")
		}
		buf.Truncate(buf.Len() - 2)
	}
	buf.WriteByte(')')

	return buf
}
