// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"errors"
	"math"
	"testing"
	"time"

	td "github.com/maxatome/go-testdeep"
)

// These tests are only here to ensure that all obsolete but aliased
// functions are still callable from outside.
//
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td for real
// tests and examples.
func TestCompat(tt *testing.T) {
	type MyStruct struct {
		Num int64  `json:"num"`
		Str string `json:"str"`
	}

	tt.Run("Cmp", func(t *testing.T) {
		td.Cmp(t, 1, 1)
		td.CmpDeeply(t, 1, 1)

		td.CmpTrue(t, true)
		td.CmpFalse(t, false)
		td.CmpError(t, errors.New("Error"))
		td.CmpNoError(t, nil)
		td.CmpPanic(t, func() { panic("boom!") }, "boom!")
		td.CmpNotPanic(t, func() {})

		td.CmpTrue(t, td.EqDeeply(1, 1))
		td.CmpNoError(t, td.EqDeeplyError(1, 1))
	})

	td.AddAnchorableStructType(func(nextAnchor int) MyStruct {
		return MyStruct{Num: 999999999 - int64(nextAnchor)}
	})

	tt.Run("td.T", func(tt *testing.T) {
		t := td.NewT(tt)
		t.Cmp(1, 1)

		assert := td.Assert(tt)
		assert.Cmp(1, 1)

		require := td.Require(tt)
		require.Cmp(1, 1)

		assert, require = td.AssertRequire(tt)
		assert.Cmp(1, 1)
		require.Cmp(1, 1)
	})

	tt.Run("All", func(t *testing.T) {
		td.Cmp(t, 1, td.All(1))
		td.CmpAll(t, 1, []any{1})
	})

	tt.Run("Any", func(t *testing.T) {
		td.Cmp(t, 1, td.Any(3, 2, 1))
		td.CmpAny(t, 1, []any{3, 2, 1})
	})

	tt.Run("Array", func(t *testing.T) {
		td.Cmp(t, [2]int{1, 2}, td.Array([2]int{}, td.ArrayEntries{0: 1, 1: 2}))
		td.CmpArray(t, [2]int{1, 2}, [2]int{}, td.ArrayEntries{0: 1, 1: 2})
	})

	tt.Run("ArrayEach", func(t *testing.T) {
		got := []int{1, 1}
		td.Cmp(t, got, td.ArrayEach(1))
		td.CmpArrayEach(t, got, 1)
	})

	tt.Run("Bag", func(t *testing.T) {
		got := []int{1, 2}
		td.Cmp(t, got, td.Bag(1, 2))
		td.CmpBag(t, got, []any{1, 2})
	})

	tt.Run("Between", func(t *testing.T) {
		for _, bounds := range []td.BoundsKind{
			td.BoundsInIn, td.BoundsInOut, td.BoundsOutIn, td.BoundsOutOut,
		} {
			td.Cmp(t, 5, td.Between(0, 10, bounds))
			td.CmpBetween(t, 5, 0, 10, bounds)
		}
	})

	tt.Run("Cap", func(t *testing.T) {
		got := make([]int, 2, 3)
		td.Cmp(t, got, td.Cap(3))
		td.CmpCap(t, got, 3)
	})

	tt.Run("Catch", func(t *testing.T) {
		var num int
		td.Cmp(t, 12, td.Catch(&num, 12))
		td.Cmp(t, num, 12)
	})

	tt.Run("Code", func(t *testing.T) {
		fn := func(n int) bool { return n == 5 }
		td.Cmp(t, 5, td.Code(fn))
		td.CmpCode(t, 5, fn)
	})

	tt.Run("Contains", func(t *testing.T) {
		td.Cmp(t, "foobar", td.Contains("ob"))
		td.CmpContains(t, "foobar", "ob")
	})

	tt.Run("ContainsKey", func(t *testing.T) {
		got := map[string]bool{"a": false}
		td.Cmp(t, got, td.ContainsKey("a"))
		td.CmpContainsKey(t, got, "a")
	})

	tt.Run("Empty", func(t *testing.T) {
		td.Cmp(t, "", td.Empty())
		td.CmpEmpty(t, "")
	})

	tt.Run("Gt", func(t *testing.T) {
		td.Cmp(t, 5, td.Gt(3))
		td.CmpGt(t, 5, 3)
	})

	tt.Run("Gte", func(t *testing.T) {
		td.Cmp(t, 5, td.Gte(3))
		td.CmpGte(t, 5, 3)
	})

	tt.Run("HasPrefix", func(t *testing.T) {
		td.Cmp(t, "foobar", td.HasPrefix("foo"))
		td.CmpHasPrefix(t, "foobar", "foo")
	})

	tt.Run("HasSuffix", func(t *testing.T) {
		td.Cmp(t, "foobar", td.HasSuffix("bar"))
		td.CmpHasSuffix(t, "foobar", "bar")
	})

	td.Cmp(tt, 42, td.Ignore())

	tt.Run("Isa", func(t *testing.T) {
		td.Cmp(t, 2, td.Isa(0))
		td.CmpIsa(t, 2, 0)
	})

	tt.Run("JSON", func(t *testing.T) {
		td.Cmp(t, []int{1, 2}, td.JSON(`[1,$val]`, td.Tag("val", 2)))
		td.CmpJSON(t, []int{1, 2}, `[1,$val]`, []any{td.Tag("val", 2)})
	})

	tt.Run("Keys", func(t *testing.T) {
		got := map[string]bool{"a": false}
		td.Cmp(t, got, td.Keys([]string{"a"}))
		td.CmpKeys(t, got, []string{"a"})
	})

	tt.Run("Lax", func(t *testing.T) {
		td.Cmp(t, int64(42), td.Lax(42))
		td.CmpLax(t, int64(42), 42)
	})

	tt.Run("Len", func(t *testing.T) {
		got := make([]int, 2, 3)
		td.Cmp(t, got, td.Len(2))
		td.CmpLen(t, got, 2)
	})

	tt.Run("Lt", func(t *testing.T) {
		td.Cmp(t, 5, td.Lt(10))
		td.CmpLt(t, 5, 10)
	})

	tt.Run("Lte", func(t *testing.T) {
		td.Cmp(t, 5, td.Lte(10))
		td.CmpLte(t, 5, 10)
	})

	tt.Run("Map", func(t *testing.T) {
		got := map[string]bool{"a": false, "b": true}
		td.Cmp(t, got, td.Map(map[string]bool{"a": false}, td.MapEntries{"b": true}))
		td.CmpMap(t, got, map[string]bool{"a": false}, td.MapEntries{"b": true})
	})

	tt.Run("MapEach", func(t *testing.T) {
		got := map[string]int{"a": 1}
		td.Cmp(t, got, td.MapEach(1))
		td.CmpMapEach(t, got, 1)
	})

	tt.Run("N", func(t *testing.T) {
		td.Cmp(t, 12, td.N(10, 2))
		td.CmpN(t, 12, 10, 2)
	})

	tt.Run("NaN", func(t *testing.T) {
		td.Cmp(t, math.NaN(), td.NaN())
		td.CmpNaN(t, math.NaN())
	})

	tt.Run("Nil", func(t *testing.T) {
		td.Cmp(t, nil, td.Nil())
		td.CmpNil(t, nil)
	})

	tt.Run("None", func(t *testing.T) {
		td.Cmp(t, 28, td.None(3, 4, 5))
		td.CmpNone(t, 28, []any{3, 4, 5})
	})

	tt.Run("Not", func(t *testing.T) {
		td.Cmp(t, 28, td.Not(3))
		td.CmpNot(t, 28, 3)
	})

	tt.Run("NotAny", func(t *testing.T) {
		got := []int{5}
		td.Cmp(t, got, td.NotAny(1, 2, 3))
		td.CmpNotAny(t, got, []any{1, 2, 3})
	})

	tt.Run("NotEmpty", func(t *testing.T) {
		td.Cmp(t, "OOO", td.NotEmpty())
		td.CmpNotEmpty(t, "OOO")
	})

	tt.Run("NotNaN", func(t *testing.T) {
		td.Cmp(t, 12., td.NotNaN())
		td.CmpNotNaN(t, 12.)
	})

	tt.Run("NotNil", func(t *testing.T) {
		td.Cmp(t, 4, td.NotNil())
		td.CmpNotNil(t, 4)
	})

	tt.Run("NotZero", func(t *testing.T) {
		td.Cmp(t, 3, td.NotZero())
		td.CmpNotZero(t, 3)
	})

	tt.Run("Ptr", func(t *testing.T) {
		num := 12
		td.Cmp(t, &num, td.Ptr(12))
		td.CmpPtr(t, &num, 12)
	})

	tt.Run("PPtr", func(t *testing.T) {
		num := 12
		pnum := &num
		td.Cmp(t, &pnum, td.PPtr(12))
		td.CmpPPtr(t, &pnum, 12)
	})

	tt.Run("Re", func(t *testing.T) {
		td.Cmp(t, "foobar", td.Re(`o+`))
		td.CmpRe(t, "foobar", `o+`, nil)
	})

	tt.Run("ReAll", func(t *testing.T) {
		td.Cmp(t, "foo bar", td.ReAll(`([a-z]+)(?: |\z)`, td.Bag("bar", "foo")))
		td.CmpReAll(t, "foo bar", `([a-z]+)(?: |\z)`, td.Bag("bar", "foo"))
	})

	tt.Run("Set", func(t *testing.T) {
		got := []int{1, 1, 2}
		td.Cmp(t, got, td.Set(2, 1))
		td.CmpSet(t, got, []any{2, 1})
	})

	tt.Run("Shallow", func(t *testing.T) {
		got := []int{1, 2, 3}
		expected := got[:1]
		td.Cmp(t, got, td.Shallow(expected))
		td.CmpShallow(t, got, expected)
	})

	tt.Run("Slice", func(t *testing.T) {
		got := []int{1, 2}
		td.Cmp(t, got, td.Slice([]int{}, td.ArrayEntries{0: 1, 1: 2}))
		td.CmpSlice(t, got, []int{}, td.ArrayEntries{0: 1, 1: 2})
	})

	tt.Run("Smuggle", func(t *testing.T) {
		fn := func(v int) int { return v * 2 }
		td.Cmp(t, 5, td.Smuggle(fn, 10))
		td.CmpSmuggle(t, 5, fn, 10)
	})

	tt.Run("String", func(t *testing.T) {
		td.Cmp(t, "foo", td.String("foo"))
		td.CmpString(t, "foo", "foo")
	})

	tt.Run("SStruct", func(t *testing.T) {
		got := MyStruct{
			Num: 42,
			Str: "foo",
		}
		td.Cmp(t, got, td.SStruct(MyStruct{Num: 42}, td.StructFields{"Str": "foo"}))
		td.CmpSStruct(t, got, MyStruct{Num: 42}, td.StructFields{"Str": "foo"})
	})

	tt.Run("Struct", func(t *testing.T) {
		got := MyStruct{
			Num: 42,
			Str: "foo",
		}
		td.Cmp(t, got, td.Struct(MyStruct{Num: 42}, td.StructFields{"Str": "foo"}))
		td.CmpStruct(t, got, MyStruct{Num: 42}, td.StructFields{"Str": "foo"})
	})

	tt.Run("SubBagOf", func(t *testing.T) {
		got := []int{1}
		td.Cmp(t, got, td.SubBagOf(1, 1, 2))
		td.CmpSubBagOf(t, got, []any{1, 1, 2})
	})

	tt.Run("SubJSONOf", func(t *testing.T) {
		got := MyStruct{
			Num: 42,
			Str: "foo",
		}
		td.Cmp(t, got,
			td.SubJSONOf(`{"num":42,"str":$str,"zip":45600}`, td.Tag("str", "foo")))
		td.CmpSubJSONOf(t, got,
			`{"num":42,"str":$str,"zip":45600}`, []any{td.Tag("str", "foo")})
	})

	tt.Run("SubMapOf", func(t *testing.T) {
		got := map[string]int{"a": 1, "b": 2}
		td.Cmp(t, got,
			td.SubMapOf(map[string]int{"a": 1, "c": 3}, td.MapEntries{"b": 2}))
		td.CmpSubMapOf(t, got, map[string]int{"a": 1, "c": 3}, td.MapEntries{"b": 2})
	})

	tt.Run("SubSetOf", func(t *testing.T) {
		got := []int{1, 1}
		td.Cmp(t, got, td.SubSetOf(1, 2))
		td.CmpSubSetOf(t, got, []any{1, 2})
	})

	tt.Run("SuperBagOf", func(t *testing.T) {
		got := []int{1, 1, 2}
		td.Cmp(t, got, td.SuperBagOf(1))
		td.CmpSuperBagOf(t, got, []any{1})
	})

	tt.Run("SuperJSONOf", func(t *testing.T) {
		got := MyStruct{
			Num: 42,
			Str: "foo",
		}
		td.Cmp(t, got, td.SuperJSONOf(`{"str":$str}`, td.Tag("str", "foo")))
		td.CmpSuperJSONOf(t, got, `{"str":$str}`, []any{td.Tag("str", "foo")})
	})

	tt.Run("SuperMapOf", func(t *testing.T) {
		got := map[string]int{"a": 1, "b": 2, "c": 3}
		td.Cmp(t, got, td.SuperMapOf(map[string]int{"a": 1}, td.MapEntries{"b": 2}))
		td.CmpSuperMapOf(t, got, map[string]int{"a": 1}, td.MapEntries{"b": 2})
	})

	tt.Run("SuperSetOf", func(t *testing.T) {
		got := []int{1, 1, 2}
		td.Cmp(t, got, td.SuperSetOf(1))
		td.CmpSuperSetOf(t, got, []any{1})
	})

	tt.Run("TruncTime", func(t *testing.T) {
		got, err := time.Parse(time.RFC3339Nano, "2020-02-22T12:34:56.789Z")
		if err != nil {
			t.Fatal(err)
		}
		expected, err := time.Parse(time.RFC3339, "2020-02-22T12:34:56Z")
		if err != nil {
			t.Fatal(err)
		}
		td.Cmp(t, got, td.TruncTime(expected, time.Second))
		td.CmpTruncTime(t, got, expected, time.Second)
	})

	tt.Run("Values", func(t *testing.T) {
		got := map[string]string{"a": "b"}
		td.Cmp(t, got, td.Values([]string{"b"}))
		td.CmpValues(t, got, []string{"b"})
	})

	tt.Run("Zero", func(t *testing.T) {
		td.Cmp(t, 0, td.Zero())
		td.CmpZero(t, 0)
	})
}
