// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td

// AnchorT returns a typed value allowing to anchor the TestDeep
// operator operator in a go classic literal like a struct, slice,
// array or map value.
//
// X type must be compatible with operator, so if the TypeBehind
// method of operator returns a non-nil type, it has to match X.
//
// AnchorT returns a typed value ready to be embed in a go data
// structure to be compared using [T.Cmp] or [T.CmpLax]:
//
//	import (
//	  "testing"
//
//	  "github.com/maxatome/go-testdeep/td"
//	)
//
//	func TestFunc(tt *testing.T) {
//	  got := Func()
//
//	  t := td.NewT(tt)
//	  t.Cmp(got, &MyStruct{
//	    Name:    "Bob",
//	    Details: &MyDetails{
//	      Nick: t.AnchorT[string](td.HasPrefix("Bobby")),
//	      Age:  t.AnchorT[int](td.Between(40, 50)),
//	    },
//	  })
//	}
//
// Without operator anchoring feature, the previous example would have
// been:
//
//	import (
//	  "testing"
//
//	  "github.com/maxatome/go-testdeep/td"
//	)
//
//	func TestFunc(tt *testing.T) {
//	  got := Func()
//
//	  t := td.NewT(tt)
//	  t.Cmp(got, td.Struct(&MyStruct{Name: "Bob"},
//	    td.StructFields{
//	    "Details": td.Struct(&MyDetails{},
//	      td.StructFields{
//	        "Nick": td.HasPrefix("Bobby"),
//	        "Age":  td.Between(40, 50),
//	      }),
//	  }))
//	}
//
// using two times the [Struct] operator to work around the strict type
// checking of golang.
//
// AnchorT is a synonym of [T.AT].
//
// By default, the value returned by AnchorT can only be used in the
// next [T.Cmp] or [T.CmpLax] call. To make it persistent across calls,
// see [T.SetAnchorsPersist] and [T.AnchorsPersistTemporarily] methods.
//
// See also [T.AnchorsPersistTemporarily], [T.DoAnchorsPersist],
// [T.ResetAnchors], [T.SetAnchorsPersist] and [AddAnchorableStructType].
func (t *T) AnchorT[X any](operator TestDeep) X {
	t.Helper()
	return Anchor[X](t, operator)
}

// AT returns a typed value allowing to anchor the TestDeep
// operator operator in a go classic literal like a struct, slice,
// array or map value.
//
// X type must be compatible with operator, so if the TypeBehind
// method of operator returns a non-nil type, it has to match X.
//
// AT returns a typed value ready to be embed in a go data
// structure to be compared using [T.Cmp] or [T.CmpLax]:
//
//	import (
//	  "testing"
//
//	  "github.com/maxatome/go-testdeep/td"
//	)
//
//	func TestFunc(tt *testing.T) {
//	  got := Func()
//
//	  t := td.NewT(tt)
//	  t.Cmp(got, &MyStruct{
//	    Name:    "Bob",
//	    Details: &MyDetails{
//	      Nick: t.AT[string](td.HasPrefix("Bobby")),
//	      Age:  t.AT[int](td.Between(40, 50)),
//	    },
//	  })
//	}
//
// Without operator anchoring feature, the previous example would have
// been:
//
//	import (
//	  "testing"
//
//	  "github.com/maxatome/go-testdeep/td"
//	)
//
//	func TestFunc(tt *testing.T) {
//	  got := Func()
//
//	  t := td.NewT(tt)
//	  t.Cmp(got, td.Struct(&MyStruct{Name: "Bob"},
//	    td.StructFields{
//	    "Details": td.Struct(&MyDetails{},
//	      td.StructFields{
//	        "Nick": td.HasPrefix("Bobby"),
//	        "Age":  td.Between(40, 50),
//	      }),
//	  }))
//	}
//
// using two times the [Struct] operator to work around the strict type
// checking of golang.
//
// AT is a shorter synonym of [T.AnchorT].
//
// By default, the value returned by AT can only be used in the
// next [T.Cmp] or [T.CmpLax] call. To make it persistent across calls,
// see [T.SetAnchorsPersist] and [T.AnchorsPersistTemporarily] methods.
//
// See also [T.AnchorsPersistTemporarily], [T.DoAnchorsPersist],
// [T.ResetAnchors], [T.SetAnchorsPersist] and [AddAnchorableStructType].
func (t *T) AT[X any](operator TestDeep) X {
	t.Helper()
	return A[X](t, operator)
}

