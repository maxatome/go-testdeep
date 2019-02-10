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

	// Visited.Record() needs its param be addressable, that's why we
	// use a struct pointer below

	t.Run("map", func(t *testing.T) {
		v := visited.NewVisited()

		type vMap struct{ m map[string]bool }
		a, b := &vMap{m: map[string]bool{}}, &vMap{m: map[string]bool{}}

		f := func(vm *vMap) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))
	})

	t.Run("slice", func(t *testing.T) {
		v := visited.NewVisited()

		type vSlice struct{ s []string }
		a, b := &vSlice{s: []string{}}, &vSlice{}

		f := func(vm *vSlice) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))
	})

	t.Run("ptr", func(t *testing.T) {
		v := visited.NewVisited()

		type vPtr struct{ p *int }
		n := 42
		a, b := &vPtr{p: &n}, &vPtr{}

		f := func(vm *vPtr) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))
	})

	t.Run("interface", func(t *testing.T) {
		v := visited.NewVisited()

		type vIf struct{ i interface{} }
		a, b := &vIf{i: 42}, &vIf{}

		f := func(vm *vIf) reflect.Value {
			return reflect.ValueOf(vm).Elem().Field(0)
		}

		test.IsFalse(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(a), f(b)))
		test.IsTrue(t, v.Record(f(b), f(a)))
	})
}
