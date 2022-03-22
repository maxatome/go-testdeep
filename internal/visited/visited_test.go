// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package visited_test

import (
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/visited"
)

func TestVisited(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		v := visited.NewVisited()
		a, b := 1, 2
		test.IsFalse(t, v.Record(reflect.ValueOf(a), reflect.ValueOf(b)))
		test.IsFalse(t, v.Record(reflect.ValueOf(a), reflect.ValueOf(b)))
	})

	t.Run("map", func(t *testing.T) {
		v := visited.NewVisited()

		a, b := map[string]bool{}, map[string]bool{}

		f := func(m map[string]bool) reflect.Value {
			return reflect.ValueOf(m)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))

		// nil maps are not recorded
		b = nil
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(b), f(a)))
		test.IsFalse(t, v.Record(f(b), f(a)))
	})

	t.Run("pointer", func(t *testing.T) {
		v := visited.NewVisited()

		type S struct {
			p  *S
			ok bool
		}
		a, b := &S{}, &S{}
		a.p = &S{ok: true}
		b.p = &S{ok: false}

		f := func(m *S) reflect.Value {
			return reflect.ValueOf(m)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))

		test.IsFalse(t, v.Record(f(a.p), f(b.p)))
		test.IsTrue(t, v.Record(f(a.p), f(b.p)))
		test.IsTrue(t, v.Record(f(b.p), f(a.p)))

		// nil pointers are not recorded
		b = nil
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(b), f(a)))
		test.IsFalse(t, v.Record(f(b), f(a)))
	})

	// Visited.Record() needs its slice or interface param be
	// addressable, that's why we use a struct pointer below

	t.Run("slice", func(t *testing.T) {
		v := visited.NewVisited()

		type vSlice struct{ s []string }
		a, b := &vSlice{s: []string{}}, &vSlice{[]string{}}

		f := func(vm *vSlice) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))

		// nil slices are not recorded
		b = &vSlice{}
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(b), f(a)))
		test.IsFalse(t, v.Record(f(b), f(a)))
	})

	t.Run("interface", func(t *testing.T) {
		v := visited.NewVisited()

		type vIf struct{ i any }
		a, b := &vIf{i: 42}, &vIf{i: 24}

		f := func(vm *vIf) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))

		// nil interfaces are not recorded
		b = &vIf{}
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsFalse(t, v.Record(f(b), f(a)))
		test.IsFalse(t, v.Record(f(b), f(a)))
	})
}
