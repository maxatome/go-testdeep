// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestMap(t *testing.T) {
	type MyMap map[string]int

	//
	// Map
	checkOK(t, (map[string]int)(nil), td.Map(map[string]int{}, nil))

	checkError(t, nil, td.Map(map[string]int{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA`),
			Got:      mustBe("nil"),
			Expected: mustBe("map[string]int{}"),
		})

	gotMap := map[string]int{"foo": 1, "bar": 2}

	checkOK(t, gotMap, td.Map(map[string]int{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotMap,
		td.Map(map[string]int{"foo": 1}, td.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		td.Map(map[string]int{}, td.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotMap,
		td.Map((map[string]int)(nil), td.MapEntries{"foo": 1, "bar": 2}))

	one := 1
	checkOK(t, map[string]*int{"foo": nil, "bar": &one},
		td.Map(map[string]*int{}, td.MapEntries{"foo": nil, "bar": &one}))

	checkError(t, gotMap, td.Map(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, td.Map(map[string]int{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotMap, td.Map(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotMap,
		td.Map(map[string]int{}, td.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotMap,
		td.Map(map[string]int{}, td.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")`),
		})
	checkError(t, gotMap,
		td.Map(MyMap{}, td.MapEntries{"foo": 1, "bar": 2}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("td_test.MyMap"),
		})

	//
	// Map type
	gotTypedMap := MyMap{"foo": 1, "bar": 2}

	checkOK(t, gotTypedMap, td.Map(MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotTypedMap,
		td.Map(MyMap{"foo": 1}, td.MapEntries{"bar": 2}))
	checkOK(t, gotTypedMap,
		td.Map(MyMap{}, td.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotTypedMap,
		td.Map((MyMap)(nil), td.MapEntries{"foo": 1, "bar": 2}))

	checkOK(t, &gotTypedMap, td.Map(&MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, &gotTypedMap,
		td.Map(&MyMap{"foo": 1}, td.MapEntries{"bar": 2}))
	checkOK(t, &gotTypedMap,
		td.Map(&MyMap{}, td.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, &gotTypedMap,
		td.Map((*MyMap)(nil), td.MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotTypedMap, td.Map(MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotTypedMap, td.Map(MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotTypedMap, td.Map(MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotTypedMap, td.Map(MyMap{}, td.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, gotTypedMap,
		td.Map(MyMap{}, td.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")`),
		})

	checkError(t, &gotTypedMap, td.Map(&MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, &gotTypedMap, td.Map(&MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, &gotTypedMap, td.Map(&MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, &gotTypedMap, td.Map(&MyMap{}, td.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: \("test"\)\nExtra 2 keys: \("bar",\s+"foo"\)\z`),
		})
	checkError(t, &gotTypedMap,
		td.Map(&MyMap{}, td.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")`),
		})

	checkError(t, &gotMap, td.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*map[string]int"),
			Expected: mustBe("*td_test.MyMap"),
		})
	checkError(t, gotMap, td.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("*td_test.MyMap"),
		})
	checkError(t, nil, td.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("*td_test.MyMap{}"),
		})
	checkError(t, nil, td.Map(MyMap{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("td_test.MyMap{}"),
		})

	//
	// nil cases
	var (
		gotNilMap      map[string]int
		gotNilTypedMap MyMap
	)

	checkOK(t, gotNilMap, td.Map(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, td.Map(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, td.Map(&MyMap{}, nil))

	// Be lax...
	// Without Lax → error
	checkError(t, MyMap{}, td.Map(map[string]int{}, nil),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkError(t, map[string]int{}, td.Map(MyMap{}, nil),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	// With Lax → OK
	checkOK(t, MyMap{}, td.Lax(td.Map(map[string]int{}, nil)))
	checkOK(t, map[string]int{}, td.Lax(td.Map(MyMap{}, nil)))

	//
	// SuperMapOf
	checkOK(t, gotMap, td.SuperMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotMap,
		td.SuperMapOf(map[string]int{"foo": 1}, td.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		td.SuperMapOf(map[string]int{}, td.MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotMap,
		td.SuperMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, td.SuperMapOf(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")`),
		})
	checkError(t, gotMap,
		td.SuperMapOf(map[string]int{}, td.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")`),
		})

	checkOK(t, gotNilMap, td.SuperMapOf(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, td.SuperMapOf(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, td.SuperMapOf(&MyMap{}, nil))

	//
	// SubMapOf
	checkOK(t, gotMap,
		td.SubMapOf(map[string]int{"foo": 1, "bar": 2, "tst": 3}, nil))
	checkOK(t, gotMap,
		td.SubMapOf(map[string]int{"foo": 1, "tst": 3}, td.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		td.SubMapOf(map[string]int{}, td.MapEntries{"foo": 1, "bar": 2, "tst": 3}))

	checkError(t, gotMap,
		td.SubMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, td.SubMapOf(map[string]int{"foo": 1}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Extra key: ("bar")`),
		})
	checkError(t, gotMap,
		td.SubMapOf(map[string]int{}, td.MapEntries{"foo": 1, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing key: ("test")
  Extra key: ("bar")`),
		})

	checkOK(t, gotNilMap, td.SubMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotNilTypedMap, td.SubMapOf(MyMap{"foo": 1}, nil))
	checkOK(t, &gotNilTypedMap, td.SubMapOf(&MyMap{"foo": 1}, nil))

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Map("test", nil) }, "usage: Map(")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SuperMapOf("test", nil) },
		"usage: SuperMapOf(")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SubMapOf("test", nil) },
		"usage: SubMapOf(")

	num := 12
	dark.CheckFatalizerBarrierErr(t, func() { td.Map(&num, nil) }, "usage: Map(")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SuperMapOf(&num, nil) },
		"usage: SuperMapOf(")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SubMapOf(&num, nil) },
		"usage: SubMapOf(")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Map(&MyMap{}, td.MapEntries{1: 2}) },
		"Map(): expected key 1 type mismatch: int != model key type (string)")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SuperMapOf(&MyMap{}, td.MapEntries{1: 2}) },
		"SuperMapOf(): expected key 1 type mismatch: int != model key type (string)")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.SubMapOf(&MyMap{}, td.MapEntries{1: 2}) },
		"SubMapOf(): expected key 1 type mismatch: int != model key type (string)")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Map(&MyMap{}, td.MapEntries{"foo": nil}) },
		`Map(): expected key "foo" value cannot be nil as entries value type is int`)

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Map(&MyMap{}, td.MapEntries{"foo": uint16(2)}) },
		`Map(): expected key "foo" value type mismatch: uint16 != model key type (int)`)

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Map(&MyMap{"foo": 1}, td.MapEntries{"foo": 1}) },
		`Map(): "foo" entry exists in both model & expectedEntries`)

	//
	// String
	test.EqualStr(t, td.Map(MyMap{}, nil).String(),
		"td_test.MyMap{}")
	test.EqualStr(t, td.Map(&MyMap{}, nil).String(),
		"*td_test.MyMap{}")
	test.EqualStr(t, td.Map(&MyMap{"foo": 2}, nil).String(),
		`*td_test.MyMap{
  "foo": 2,
}`)

	test.EqualStr(t, td.SubMapOf(MyMap{}, nil).String(),
		"SubMapOf(td_test.MyMap{})")
	test.EqualStr(t, td.SubMapOf(&MyMap{}, nil).String(),
		"SubMapOf(*td_test.MyMap{})")
	test.EqualStr(t, td.SubMapOf(&MyMap{"foo": 2}, nil).String(),
		`SubMapOf(*td_test.MyMap{
  "foo": 2,
})`)

	test.EqualStr(t, td.SuperMapOf(MyMap{}, nil).String(),
		"SuperMapOf(td_test.MyMap{})")
	test.EqualStr(t, td.SuperMapOf(&MyMap{}, nil).String(),
		"SuperMapOf(*td_test.MyMap{})")
	test.EqualStr(t, td.SuperMapOf(&MyMap{"foo": 2}, nil).String(),
		`SuperMapOf(*td_test.MyMap{
  "foo": 2,
})`)
}

func TestMapTypeBehind(t *testing.T) {
	type MyMap map[string]int

	// Map
	equalTypes(t, td.Map(map[string]int{}, nil), map[string]int{})
	equalTypes(t, td.Map(MyMap{}, nil), MyMap{})
	equalTypes(t, td.Map(&MyMap{}, nil), &MyMap{})

	// SubMap
	equalTypes(t, td.SubMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, td.SubMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, td.SubMapOf(&MyMap{}, nil), &MyMap{})

	// SuperMap
	equalTypes(t, td.SuperMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, td.SuperMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, td.SuperMapOf(&MyMap{}, nil), &MyMap{})
}
