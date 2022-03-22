// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package anchors_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/anchors"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestInfo(t *testing.T) {
	i := anchors.NewInfo()

	test.IsFalse(t, i.DoAnchorsPersist())

	i.SetAnchorsPersist(true)
	test.IsTrue(t, i.DoAnchorsPersist())

	i.SetAnchorsPersist(false)
	test.IsFalse(t, i.DoAnchorsPersist())
}

func TestBuildResolveAnchor(t *testing.T) {
	var i anchors.Info

	checkResolveAnchor := func(t *testing.T, val any, opName string) {
		t.Helper()
		v1, err := i.AddAnchor(reflect.TypeOf(val), reflect.ValueOf(opName+" (1)"))
		if !test.NoError(t, err, "first anchor") {
			return
		}
		v2, err := i.AddAnchor(reflect.TypeOf(val), reflect.ValueOf(opName+" (2)"))
		if !test.NoError(t, err, "second anchor") {
			return
		}

		op, found := i.ResolveAnchor(v1)
		test.IsTrue(t, found, "first anchor found")
		test.EqualStr(t, op.String(), opName+" (1)", "first anchor operator OK")

		op, found = i.ResolveAnchor(v2)
		test.IsTrue(t, found, "second anchor found")
		test.EqualStr(t, op.String(), opName+" (2)", "second anchor operator OK")
	}

	t.Run("AddAnchor basic types", func(t *testing.T) {
		checkResolveAnchor(t, 0, "int")
		checkResolveAnchor(t, int8(0), "int8")
		checkResolveAnchor(t, int16(0), "int16")
		checkResolveAnchor(t, int32(0), "int32")
		checkResolveAnchor(t, int64(0), "int64")

		checkResolveAnchor(t, uint(0), "uint")
		checkResolveAnchor(t, uint8(0), "uint8")
		checkResolveAnchor(t, uint16(0), "uint16")
		checkResolveAnchor(t, uint32(0), "uint32")
		checkResolveAnchor(t, uint64(0), "uint64")

		checkResolveAnchor(t, uintptr(0), "uintptr")

		checkResolveAnchor(t, float32(0), "float32")
		checkResolveAnchor(t, float64(0), "float64")

		checkResolveAnchor(t, complex(float32(0), 0), "complex64")
		checkResolveAnchor(t, complex(float64(0), 0), "complex128")

		checkResolveAnchor(t, "", "string")
		checkResolveAnchor(t, (chan int)(nil), "chan")
		checkResolveAnchor(t, (map[string]bool)(nil), "map")
		checkResolveAnchor(t, ([]int)(nil), "slice")
		checkResolveAnchor(t, (*time.Time)(nil), "pointer")
	})

	t.Run("AddAnchor", func(t *testing.T) {
		oldAnchorableTypes := anchors.AnchorableTypes
		defer func() { anchors.AnchorableTypes = oldAnchorableTypes }()

		type ok struct{ index int }

		// AddAnchor for ok type
		err := anchors.AddAnchorableStructType(func(nextAnchor int) ok {
			return ok{index: 1000 + nextAnchor}
		})
		if err != nil {
			t.Fatalf("AddAnchorableStructType failed: %s", err)
		}
		checkResolveAnchor(t, ok{}, "ok{}")

		// AddAnchor for ok convertible type
		type okConvert ok
		checkResolveAnchor(t, okConvert{}, "okConvert{}")

		// Replace ok type
		err = anchors.AddAnchorableStructType(func(nextAnchor int) ok {
			return ok{index: 2000 + nextAnchor}
		})
		if err != nil {
			t.Fatalf("AddAnchorableStructType failed: %s", err)
		}
		if len(anchors.AnchorableTypes) != 2 {
			t.Fatalf("Bad number of anchored type: got=%d expected=2",
				len(anchors.AnchorableTypes))
		}
		checkResolveAnchor(t, ok{}, "ok{}")

		// AddAnchor for builtin time.Time type
		checkResolveAnchor(t, time.Time{}, "time.Time{}")

		// AddAnchor for unknown type
		_, err = i.AddAnchor(reflect.TypeOf(func() {}), reflect.ValueOf(123))
		if test.Error(t, err) {
			test.EqualStr(t, err.Error(), "func kind is not supported as an anchor")
		}

		// AddAnchor for unknown struct type
		_, err = i.AddAnchor(reflect.TypeOf(struct{}{}), reflect.ValueOf(123))
		if test.Error(t, err) {
			test.EqualStr(t,
				err.Error(),
				"struct {} struct type is not supported as an anchor. Try AddAnchorableStructType")
		}

		// Struct not comparable
		type notComparable struct{ s []int }
		v := reflect.ValueOf(notComparable{s: []int{42}})
		op, found := i.ResolveAnchor(v)
		test.IsFalse(t, found)
		if !reflect.DeepEqual(v.Interface(), op.Interface()) {
			test.EqualErrorMessage(t, op.Interface(), v.Interface())
		}

		// Struct comparable but not anchored
		v = reflect.ValueOf(struct{}{})
		op, found = i.ResolveAnchor(v)
		test.IsFalse(t, found)
		if !reflect.DeepEqual(v.Interface(), op.Interface()) {
			test.EqualErrorMessage(t, op.Interface(), v.Interface())
		}

		// Struct anchored once, but not for this value
		v = reflect.ValueOf(ok{index: 42424242})
		op, found = i.ResolveAnchor(v)
		test.IsFalse(t, found)
		if !reflect.DeepEqual(v.Interface(), op.Interface()) {
			test.EqualErrorMessage(t, op.Interface(), v.Interface())
		}

		// Kind not supported
		v = reflect.ValueOf(true)
		op, found = i.ResolveAnchor(v)
		test.IsFalse(t, found)
		if !reflect.DeepEqual(v.Interface(), op.Interface()) {
			test.EqualErrorMessage(t, op.Interface(), v.Interface())
		}
	})

	t.Run("ResetAnchors", func(t *testing.T) {
		v, err := i.AddAnchor(reflect.TypeOf(12), reflect.ValueOf("zip"))
		if !test.NoError(t, err) {
			return
		}

		op, found := i.ResolveAnchor(v)
		test.IsTrue(t, found)
		test.EqualStr(t, op.String(), "zip")

		i.SetAnchorsPersist(true)
		i.ResetAnchors(false)

		op, found = i.ResolveAnchor(v)
		test.IsTrue(t, found)
		test.EqualStr(t, op.String(), "zip")

		i.ResetAnchors(true)
		_, found = i.ResolveAnchor(reflect.ValueOf(42))
		test.IsFalse(t, found)

		i.SetAnchorsPersist(false)
		v, err = i.AddAnchor(reflect.TypeOf(12), reflect.ValueOf("xxx"))
		if !test.NoError(t, err) {
			return
		}

		op, found = i.ResolveAnchor(v)
		test.IsTrue(t, found)
		test.EqualStr(t, op.String(), "xxx")

		i.ResetAnchors(false)
		_, found = i.ResolveAnchor(reflect.ValueOf(42))
		test.IsFalse(t, found)
	})

	t.Run("skip", func(t *testing.T) {
		var i *anchors.Info

		_, found := i.ResolveAnchor(reflect.ValueOf(42))
		test.IsFalse(t, found)

		i = &anchors.Info{}
		_, found = i.ResolveAnchor(reflect.ValueOf(42))
		test.IsFalse(t, found)
	})
}
