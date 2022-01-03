// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package anchors

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
)

type anchor struct {
	Anchor   reflect.Value // Anchor is the generated value used as anchor
	Operator reflect.Value // Operator is a td.TestDeep behind
}

// Info gathers all anchors information.
type Info struct {
	sync.Mutex
	index   int
	persist bool
	anchors map[interface{}]anchor
}

// NewInfo returns a new instance of *Info.
func NewInfo() *Info {
	return &Info{
		anchors: map[interface{}]anchor{},
	}
}

// AddAnchor anchors a new operator op, with type typ then returns the
// anchor value.
func (i *Info) AddAnchor(typ reflect.Type, op reflect.Value) (reflect.Value, error) {
	i.Lock()
	defer i.Unlock()

	anc, key, err := i.build(typ)
	if err != nil {
		return reflect.Value{}, err
	}

	if i.anchors == nil {
		i.anchors = map[interface{}]anchor{}
	}

	i.anchors[key] = anchor{
		Anchor:   anc,
		Operator: op,
	}

	return anc, nil
}

// DoAnchorsPersist returns true if anchors are persistent across tests.
func (i *Info) DoAnchorsPersist() bool {
	i.Lock()
	defer i.Unlock()
	return i.persist
}

// SetAnchorsPersist enables or disables anchors persistence.
func (i *Info) SetAnchorsPersist(persist bool) {
	i.Lock()
	defer i.Unlock()
	i.persist = persist
}

// ResetAnchors removes all anchors if persistence is disabled or
// force is true.
func (i *Info) ResetAnchors(force bool) {
	i.Lock()
	defer i.Unlock()

	if !i.persist || force {
		for k := range i.anchors {
			delete(i.anchors, k)
		}
		i.index = 0
	}
}

func (i *Info) nextIndex() (n int) {
	n = i.index
	i.index++
	return
}

// ResolveAnchor checks whether the passed value matches an anchored
// operator or not. If yes, this operator is returned with true. If
// no, the value is returned as is with false.
func (i *Info) ResolveAnchor(v reflect.Value) (reflect.Value, bool) {
	if i == nil || !v.CanInterface() {
		return v, false
	}
	// Shortcut
	i.Lock()
	la := len(i.anchors)
	i.Unlock()
	if la == 0 {
		return v, false
	}

	var key interface{}
sw:
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		key = v.Interface()

	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Ptr:
		key = v.Pointer()

	case reflect.Struct:
		typ := v.Type()
		if typ.Comparable() {
			// Check for anchorable types. No need of 2 passes here.
			for _, at := range AnchorableTypes {
				if typ == at.typ || at.typ.ConvertibleTo(typ) { // 1.17 ok as struct here
					key = v.Interface()
					break sw
				}
			}
		}
		fallthrough

	default:
		return v, false
	}

	i.Lock()
	defer i.Unlock()
	if anchor, ok := i.anchors[key]; ok {
		return anchor.Operator, true
	}
	return v, false
}

func (i *Info) setInt(typ reflect.Type, min int64) (reflect.Value, interface{}) {
	nvm := reflect.New(typ).Elem()
	nvm.SetInt(min + int64(i.nextIndex()))
	return nvm, nvm.Interface()
}

func (i *Info) setUint(typ reflect.Type, max uint64) (reflect.Value, interface{}) {
	nvm := reflect.New(typ).Elem()
	nvm.SetUint(max - uint64(i.nextIndex()))
	return nvm, nvm.Interface()
}

func (i *Info) setFloat(typ reflect.Type, min float64) (reflect.Value, interface{}) {
	nvm := reflect.New(typ).Elem()
	nvm.SetFloat(min + float64(i.nextIndex()))
	return nvm, nvm.Interface()
}

func (i *Info) setComplex(typ reflect.Type, min float64) (reflect.Value, interface{}) {
	nvm := reflect.New(typ).Elem()
	min += float64(i.nextIndex())
	nvm.SetComplex(complex(min, min))
	return nvm, nvm.Interface()
}

// build builds a new value of type "typ" and returns it under two
// forms:
//   - the new value itself as a reflect.Value;
//   - an interface{} usable as a key in an AnchorsSet map.
//
// It returns an error if "typ" kind is not recognized or if it is a
// non-anchorable struct.
func (i *Info) build(typ reflect.Type) (reflect.Value, interface{}, error) {
	// For each numeric type, anchor the operator on a number close to
	// the limit of this type, but not at the extreme limit to avoid
	// edge cases where these limits are used in real world and so avoid
	// collisions
	switch typ.Kind() {
	case reflect.Int:
		nvm, iface := i.setInt(typ, int64(^int(^uint(0)>>1))+1004293)
		return nvm, iface, nil
	case reflect.Int8:
		nvm, iface := i.setInt(typ, math.MinInt8+13)
		return nvm, iface, nil
	case reflect.Int16:
		nvm, iface := i.setInt(typ, math.MinInt16+1049)
		return nvm, iface, nil
	case reflect.Int32:
		nvm, iface := i.setInt(typ, math.MinInt32+1004293)
		return nvm, iface, nil
	case reflect.Int64:
		nvm, iface := i.setInt(typ, math.MinInt64+1000424443)
		return nvm, iface, nil

	case reflect.Uint:
		nvm, iface := i.setUint(typ, uint64(^uint(0))-1004293)
		return nvm, iface, nil
	case reflect.Uint8:
		nvm, iface := i.setUint(typ, math.MaxUint8-29)
		return nvm, iface, nil
	case reflect.Uint16:
		nvm, iface := i.setUint(typ, math.MaxUint16-2099)
		return nvm, iface, nil
	case reflect.Uint32:
		nvm, iface := i.setUint(typ, math.MaxUint32-2008571)
		return nvm, iface, nil
	case reflect.Uint64:
		nvm, iface := i.setUint(typ, math.MaxUint64-2000848901)
		return nvm, iface, nil
	case reflect.Uintptr:
		nvm, iface := i.setUint(typ, uint64(^uintptr(0))-2000848901)
		return nvm, iface, nil

	case reflect.Float32:
		nvm, iface := i.setFloat(typ, -(1<<24)+104243)
		return nvm, iface, nil
	case reflect.Float64:
		nvm, iface := i.setFloat(typ, -(1<<53)+100004243)
		return nvm, iface, nil

	case reflect.Complex64:
		nvm, iface := i.setComplex(typ, -(1<<24)+104243)
		return nvm, iface, nil
	case reflect.Complex128:
		nvm, iface := i.setComplex(typ, -(1<<53)+100004243)
		return nvm, iface, nil

	case reflect.String:
		nvm := reflect.New(typ).Elem()
		nvm.SetString(fmt.Sprintf("<testdeep@anchor#%d>", i.nextIndex()))
		return nvm, nvm.Interface(), nil

	case reflect.Chan:
		nvm := reflect.MakeChan(typ, 0)
		return nvm, nvm.Pointer(), nil

	case reflect.Map:
		nvm := reflect.MakeMap(typ)
		return nvm, nvm.Pointer(), nil

	case reflect.Slice:
		nvm := reflect.MakeSlice(typ, 0, 1) // cap=1 to avoid same ptr below
		return nvm, nvm.Pointer(), nil

	case reflect.Ptr:
		nvm := reflect.New(typ.Elem())
		return nvm, nvm.Pointer(), nil

	case reflect.Struct:
		// First pass for the exact type
		for _, at := range AnchorableTypes {
			if typ == at.typ {
				nvm := at.builder.Call([]reflect.Value{reflect.ValueOf(i.nextIndex())})[0]
				return nvm, nvm.Interface(), nil
			}
		}
		// Second pass for convertible type
		for _, at := range AnchorableTypes {
			if at.typ.ConvertibleTo(typ) {
				nvm := at.builder.Call([]reflect.Value{reflect.ValueOf(i.nextIndex())})[0].
					Convert(typ)
				return nvm, nvm.Interface(), nil
			}
		}
		return reflect.Value{}, nil,
			errors.New(typ.String() + " struct type is not supported as an anchor. Try AddAnchorableStructType")

	default:
		return reflect.Value{}, nil,
			errors.New(typ.Kind().String() + " kind is not supported as an anchor")
	}
}
