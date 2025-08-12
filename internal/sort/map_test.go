// Copyright (c) 2018-2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package sort_test

import (
	"reflect"
	"sort"
	"testing"

	tdsort "github.com/maxatome/go-testdeep/internal/sort"
)

func TestMapEach(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	t.Run("full", func(t *testing.T) {
		type kv struct {
			key   string
			value int
		}
		var s []kv
		ok := tdsort.MapEach(reflect.ValueOf(m), func(k, v reflect.Value) bool {
			s = append(s, kv{
				key:   k.Interface().(string),
				value: v.Interface().(int),
			})
			return true
		})
		if !ok {
			t.Error("MapEach returned false")
		}

		sort.Slice(s, func(i, j int) bool { return s[i].key < s[j].key })

		if len(s) != 3 ||
			s[0] != (kv{key: "a", value: 1}) ||
			s[1] != (kv{key: "b", value: 2}) ||
			s[2] != (kv{key: "c", value: 3}) {
			t.Errorf("MapEach failed: %v", s)
		}
	})

	t.Run("short circuit", func(t *testing.T) {
		called := 0
		ok := tdsort.MapEach(reflect.ValueOf(m), func(k, v reflect.Value) bool {
			called++
			return false
		})
		if ok {
			t.Error("MapEach returned true")
		}
		if called != 1 {
			t.Errorf("MapEach callback called %d times instead of 1", called)
		}
	})
}
