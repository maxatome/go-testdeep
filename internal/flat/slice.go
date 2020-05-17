// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package flat

import (
	"reflect"
)

var SliceType = reflect.TypeOf(Slice{})

type Slice struct {
	Slice interface{}
}

// isFlat returns true if no flat.Slice items can be contained in
// f.Slice, so this Slice is already flattened.
func (f Slice) isFlat() bool {
	t := reflect.TypeOf(f.Slice).Elem()
	return t != SliceType && t.Kind() != reflect.Interface
}

func (f Slice) len() int {
	if f.isFlat() {
		return reflect.ValueOf(f.Slice).Len()
	}

	fv := reflect.ValueOf(f.Slice)
	fvLen := fv.Len()
	l := fvLen

	for i := 0; i < fvLen; i++ {
		if subf, ok := fv.Index(i).Interface().(Slice); ok {
			l += subf.len() - 1
		}
	}
	return l
}

func (f Slice) appendValuesTo(sv []reflect.Value) []reflect.Value {
	fv := reflect.ValueOf(f.Slice)
	fvLen := fv.Len()

	if f.isFlat() {
		for i := 0; i < fvLen; i++ {
			sv = append(sv, fv.Index(i))
		}
		return sv
	}

	for i := 0; i < fvLen; i++ {
		cv := fv.Index(i)
		if cv.Kind() == reflect.Interface {
			cv = cv.Elem()
		}
		if subf, ok := cv.Interface().(Slice); ok {
			sv = subf.appendValuesTo(sv)
		} else {
			sv = append(sv, cv)
		}
	}
	return sv
}

func (f Slice) appendTo(si []interface{}) []interface{} {
	fv := reflect.ValueOf(f.Slice)
	fvLen := fv.Len()

	if f.isFlat() {
		for i := 0; i < fvLen; i++ {
			si = append(si, fv.Index(i).Interface())
		}
		return si
	}

	for i := 0; i < fvLen; i++ {
		item := fv.Index(i).Interface()
		if subf, ok := item.(Slice); ok {
			si = subf.appendTo(si)
		} else {
			si = append(si, item)
		}
	}
	return si
}

func Len(items []interface{}) (int, bool) {
	l := len(items)
	flattened := true

	for _, item := range items {
		if subf, ok := item.(Slice); ok {
			l += subf.len() - 1
			flattened = false
		}
	}
	return l, flattened
}

func Values(items []interface{}) []reflect.Value {
	l, flattened := Len(items)
	if flattened {
		sv := make([]reflect.Value, l)
		for i, item := range items {
			sv[i] = reflect.ValueOf(item)
		}
		return sv
	}

	sv := make([]reflect.Value, 0, l)
	for _, item := range items {
		if f, ok := item.(Slice); ok {
			sv = f.appendValuesTo(sv)
		} else {
			sv = append(sv, reflect.ValueOf(item))
		}
	}
	return sv
}

func Interfaces(items ...interface{}) []interface{} {
	l, flattened := Len(items)
	if flattened {
		return items
	}

	si := make([]interface{}, 0, l)
	for _, item := range items {
		if f, ok := item.(Slice); ok {
			si = f.appendTo(si)
		} else {
			si = append(si, item)
		}
	}
	return si
}
