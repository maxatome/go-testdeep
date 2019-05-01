// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	t.Run("MapEach", func(t *testing.T) {
		type kv struct {
			key   string
			value int
		}
		var s []kv
		ok := tdutil.MapEach(reflect.ValueOf(m), func(k, v reflect.Value) bool {
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

	t.Run("MapEach short circuit", func(t *testing.T) {
		called := 0
		ok := tdutil.MapEach(reflect.ValueOf(m), func(k, v reflect.Value) bool {
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

	t.Run("MapEachValue", func(t *testing.T) {
		var s []int
		ok := tdutil.MapEachValue(reflect.ValueOf(m), func(v reflect.Value) bool {
			s = append(s, v.Interface().(int))
			return true
		})
		if !ok {
			t.Error("MapEachValue returned false")
		}

		sort.Ints(s)

		if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
			t.Errorf("MapEachValue failed: %v", s)
		}
	})

	t.Run("MapEachValue short circuit", func(t *testing.T) {
		called := 0
		ok := tdutil.MapEachValue(reflect.ValueOf(m), func(v reflect.Value) bool {
			called++
			return false
		})
		if ok {
			t.Error("MapEachValue returned true")
		}
		if called != 1 {
			t.Errorf("MapEachValue callback called %d times instead of 1", called)
		}
	})

	t.Run("MapSortedValues", func(t *testing.T) {
		vs := tdutil.MapSortedValues(reflect.ValueOf(m))

		if len(vs) != 3 ||
			vs[0].Int() != 1 || vs[1].Int() != 2 || vs[2].Int() != 3 {
			t.Errorf("MapSortedValues failed: %v", vs)
		}

		// nil map
		var mn map[string]int
		vs = tdutil.MapSortedKeys(reflect.ValueOf(mn))
		if len(vs) != 0 {
			t.Errorf("MapSortedValues failed: %v", vs)
		}
	})

	t.Run("MapSortedKeys", func(t *testing.T) {
		ks := tdutil.MapSortedKeys(reflect.ValueOf(m))

		if len(ks) != 3 ||
			ks[0].String() != "a" || ks[1].String() != "b" || ks[2].String() != "c" {
			t.Errorf("MapSortedKeys failed: %v", ks)
		}

		// nil map
		var mn map[string]int
		ks = tdutil.MapSortedKeys(reflect.ValueOf(mn))
		if len(ks) != 0 {
			t.Errorf("MapSortedKeys failed: %v", ks)
		}
	})
}
