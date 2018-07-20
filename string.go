// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"reflect"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

func toString(val interface{}) string {
	if val == nil {
		return "nil"
	}

typeSwitch:
	switch tval := val.(type) {
	case reflect.Value:
		newVal, ok := getInterface(tval, true)
		if ok {
			return toString(newVal)
		}

	case string:
		for _, chr := range tval {
			if !unicode.IsPrint(chr) {
				break typeSwitch
			}
		}
		return `"` + tval + `"`

	case testDeepStringer:
		return tval.String()
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
