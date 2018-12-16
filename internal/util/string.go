// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/types"

	"github.com/davecgh/go-spew/spew"
)

// ToString does it best to stringify val.
func ToString(val interface{}) string {
	if val == nil {
		return "nil"
	}

typeSwitch:
	switch tval := val.(type) {
	case reflect.Value:
		newVal, ok := dark.GetInterface(tval, true)
		if ok {
			return ToString(newVal)
		}

		// no "(string) " prefix for printable strings
	case string:
		for _, chr := range tval {
			if !unicode.IsPrint(chr) {
				break typeSwitch
			}
		}
		return `"` + tval + `"`

		// no "(int) " prefix for ints
	case int:
		return strconv.Itoa(tval)

	case types.TestDeepStringer:
		return tval.String()
	}

	return strings.TrimRight(spew.Sdump(val), "\n")
}

// IndentString indents str lines (from 2nd one = 1st line is not
// indented) by indent.
func IndentString(str string, indent string) string {
	return strings.Replace(str, "\n", "\n"+indent, -1)
}

// SliceToBuffer stringifies items slice into buf then returns buf.
func SliceToBuffer(buf *bytes.Buffer, items []reflect.Value) *bytes.Buffer {
	buf.WriteByte('(')
	if len(items) < 2 {
		if len(items) > 0 {
			buf.WriteString(ToString(items[0]))
		}
	} else {
		begLine := bytes.LastIndexByte(buf.Bytes(), '\n') + 1

		prefix := strings.Repeat(" ", buf.Len()-begLine)

		for idx, item := range items {
			if idx != 0 {
				buf.WriteString(prefix)
			}
			buf.WriteString(ToString(item))
			buf.WriteString(",\n")
		}
		buf.Truncate(buf.Len() - 2)
	}
	buf.WriteByte(')')

	return buf
}
