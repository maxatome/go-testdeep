// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil

import (
	"reflect"
	"sort"
	"testing"
)

func TestKvSlice(t *testing.T) {
	t.Run("len=0", func(t *testing.T) {
		kvs := newKvSlice(0)
		if kvs.s != nil || kvs.v != nil {
			t.Errorf("newKvSlice failed: %v", *kvs)
		}
		sort.Sort(kvs)
	})

	t.Run("len=1", func(t *testing.T) {
		kvs := newKvSlice(1)
		if kvs.s == nil || kvs.v != nil {
			t.Errorf("newKvSlice failed: %v", *kvs)
		}
		kvs.s = append(kvs.s, kv{
			key:   reflect.ValueOf("a"),
			value: reflect.ValueOf(1),
		})
		sort.Sort(kvs)
	})

	t.Run("len>1", func(t *testing.T) {
		kvs := newKvSlice(3)
		if kvs.s == nil || kvs.v == nil {
			t.Errorf("newKvSlice failed: %v", *kvs)
		}
		kvs.s = append(kvs.s,
			kv{
				key:   reflect.ValueOf("b"),
				value: reflect.ValueOf(2),
			},
			kv{
				key:   reflect.ValueOf("c"),
				value: reflect.ValueOf(3),
			},
			kv{
				key:   reflect.ValueOf("a"),
				value: reflect.ValueOf(1),
			},
		)
		sort.Sort(kvs)

		if kvs.s[0].key.String() != "a" ||
			kvs.s[1].key.String() != "b" ||
			kvs.s[2].key.String() != "c" {
			t.Errorf("Sort failed: [%v, %v, %v]",
				kvs.s[0].key, kvs.s[1].key, kvs.s[2].key)
		}
	})
}
