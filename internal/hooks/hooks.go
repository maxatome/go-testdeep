// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package hooks

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/maxatome/go-testdeep/internal/types"
)

type hooks map[reflect.Type]reflect.Value

// Info gathers all hooks information.
type Info struct {
	sync.Mutex
	cmp     hooks
	smuggle hooks
}

// NewInfo returns a new instance of *Info.
func NewInfo() *Info {
	return &Info{}
}

var ErrBoolean = errors.New("CmpHook(got, expected) failed")

func copyHooks(from hooks) hooks {
	if len(from) == 0 {
		return nil
	}

	to := make(hooks, len(from))
	for t, v := range from {
		to[t] = v
	}
	return to
}

// Copy returns a new instance of *Info with the same hooks as i. As a
// special case, if i is nil, returned instance is non-nil.
func (i *Info) Copy() *Info {
	ni := NewInfo()

	if i == nil {
		return ni
	}

	i.Lock()
	defer i.Unlock()

	ni.cmp = copyHooks(i.cmp)
	ni.smuggle = copyHooks(i.smuggle)

	return ni
}

// AddCmpHooks records new Cmp hooks using functions contained in "fns".
//
// Each function in "fns" has to be a function with the following
// possible signatures:
//   func (A, A) bool
//   func (A, A) error
// First arg is always "got", and second is always "expected".
//
// A cannot be an interface. This retriction can be removed in the
// future, if really needed.
//
// It returns an error if an item of "fns" is not a function or if its
// signature does not match the expected ones.
func (i *Info) AddCmpHooks(fns []interface{}) error {
	for n, fn := range fns {
		vfn := reflect.ValueOf(fn)

		if vfn.Kind() != reflect.Func {
			return fmt.Errorf("expects a function, not a %s (@%d)", vfn.Kind(), n)
		}

		ft := vfn.Type()
		if !ft.IsVariadic() &&
			ft.NumIn() == 2 &&
			ft.NumOut() == 1 &&
			ft.In(0) == ft.In(1) &&
			ft.In(0).Kind() != reflect.Interface &&
			(ft.Out(0) == types.Bool || ft.Out(0) == types.Error) {
			i.Lock()
			if i.cmp == nil {
				i.cmp = hooks{}
			}
			i.cmp[ft.In(0)] = vfn
			i.Unlock()
			continue
		}

		return fmt.Errorf("expects: func (T, T) bool|error not %s (@%d)", ft, n)
	}
	return nil
}

// Cmp checks if a Cmp hook exists matching "got" and "expected" types.
//
// If no, it returns (false, nil)
//
// If yes, it calls it and returns (true, nil) if it succeeds,
// (true, <an error>) if it fails. If the hook returns a false bool, the
// error returned is ErrBoolean.
func (i *Info) Cmp(got, expected reflect.Value) (bool, error) {
	if i == nil {
		return false, nil
	}

	tg := got.Type()

	i.Lock()
	vfn, ok := i.cmp[tg]
	i.Unlock()
	if !ok {
		return false, nil
	}

	if !expected.Type().AssignableTo(vfn.Type().In(1)) {
		return false, nil
	}

	res := vfn.Call([]reflect.Value{got, expected})[0]
	if res.Kind() == reflect.Bool {
		if res.Bool() {
			return true, nil
		}
		return true, ErrBoolean
	}
	err, _ := res.Interface().(error)
	return true, err
}

// AddSmuggleHooks records new Smuggle hooks using functions contained
// in "fns".
//
// Each function in "fns" has to be a function with the following
// possible signatures:
//   func (A) B
//   func (A) (B, error)
//
// A cannot be an interface. This retriction can be removed in the
// future, if really needed.
//
// B can be an interface.
//
// It returns an error if an item of "fns" is not a function or if its
// signature does not match the expected ones.
func (i *Info) AddSmuggleHooks(fns []interface{}) error {
	for n, fn := range fns {
		vfn := reflect.ValueOf(fn)

		if vfn.Kind() != reflect.Func {
			return fmt.Errorf("expects a function, not a %s (@%d)", vfn.Kind(), n)
		}

		ft := vfn.Type()
		if !ft.IsVariadic() &&
			ft.NumIn() == 1 &&
			ft.In(0).Kind() != reflect.Interface &&
			(ft.NumOut() == 1 || (ft.NumOut() == 2 && ft.Out(1) == types.Error)) &&
			ft.Out(0).Kind() != reflect.Interface {
			i.Lock()
			if i.smuggle == nil {
				i.smuggle = hooks{}
			}
			i.smuggle[ft.In(0)] = vfn
			i.Unlock()
			continue
		}

		return fmt.Errorf("expects: func (A) (B[, error]) not %s (@%d)", ft, n)
	}
	return nil
}

// Smuggle checks if a Smuggle hook exists matching "*got" type.
//
// If no, it returns (false, nil)
//
// If yes, it calls it and returns (true, nil) if it succeeds,
// (true, <an error>) if it fails.
func (i *Info) Smuggle(got *reflect.Value) (bool, error) {
	if i == nil {
		return false, nil
	}

	tg := got.Type()

	i.Lock()
	vfn, ok := i.smuggle[tg]
	i.Unlock()
	if !ok {
		return false, nil
	}

	res := vfn.Call([]reflect.Value{*got})
	if len(res) == 2 {
		if err, _ := res[1].Interface().(error); err != nil {
			return true, err
		}
	}

	*got = res[0]
	return true, nil
}
