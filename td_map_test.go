// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestMap(t *testing.T) {
	type MyMap map[string]int

	//
	// Map
	checkOK(t, (map[string]int)(nil), testdeep.Map(map[string]int{}, nil))

	checkError(t, nil, testdeep.Map(map[string]int{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA`),
			Got:      mustBe("nil"),
			Expected: mustBe("map[string]int{}"),
		})

	gotMap := map[string]int{"foo": 1, "bar": 2}

	checkOK(t, gotMap, testdeep.Map(map[string]int{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotMap,
		testdeep.Map(map[string]int{"foo": 1}, testdeep.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		testdeep.Map(map[string]int{}, testdeep.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotMap,
		testdeep.Map((map[string]int)(nil), testdeep.MapEntries{"foo": 1, "bar": 2}))

	one := 1
	checkOK(t, map[string]*int{"foo": nil, "bar": &one},
		testdeep.Map(map[string]*int{}, testdeep.MapEntries{"foo": nil, "bar": &one}))

	checkError(t, gotMap, testdeep.Map(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, testdeep.Map(map[string]int{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap, testdeep.Map(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap,
		testdeep.Map(map[string]int{}, testdeep.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap,
		testdeep.Map(map[string]int{}, testdeep.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test[^\n]+\z`),
		})
	checkError(t, gotMap,
		testdeep.Map(MyMap{}, testdeep.MapEntries{"foo": 1, "bar": 2}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("testdeep_test.MyMap"),
		})

	//
	// Map type
	gotTypedMap := MyMap{"foo": 1, "bar": 2}

	checkOK(t, gotTypedMap, testdeep.Map(MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotTypedMap,
		testdeep.Map(MyMap{"foo": 1}, testdeep.MapEntries{"bar": 2}))
	checkOK(t, gotTypedMap,
		testdeep.Map(MyMap{}, testdeep.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotTypedMap,
		testdeep.Map((MyMap)(nil), testdeep.MapEntries{"foo": 1, "bar": 2}))

	checkOK(t, &gotTypedMap, testdeep.Map(&MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, &gotTypedMap,
		testdeep.Map(&MyMap{"foo": 1}, testdeep.MapEntries{"bar": 2}))
	checkOK(t, &gotTypedMap,
		testdeep.Map(&MyMap{}, testdeep.MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, &gotTypedMap,
		testdeep.Map((*MyMap)(nil), testdeep.MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotTypedMap, testdeep.Map(MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotTypedMap, testdeep.Map(MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap, testdeep.Map(MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap, testdeep.Map(MyMap{}, testdeep.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap,
		testdeep.Map(MyMap{}, testdeep.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test[^\n]+\z`),
		})

	checkError(t, &gotTypedMap, testdeep.Map(&MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, &gotTypedMap, testdeep.Map(&MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap, testdeep.Map(&MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap, testdeep.Map(&MyMap{}, testdeep.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^ Missing key: .*test(.|\n)*\nExtra 2 keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap,
		testdeep.Map(&MyMap{}, testdeep.MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test[^\n]+\z`),
		})

	checkError(t, &gotMap, testdeep.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*map[string]int"),
			Expected: mustBe("*testdeep_test.MyMap"),
		})
	checkError(t, gotMap, testdeep.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("*testdeep_test.MyMap"),
		})
	checkError(t, nil, testdeep.Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("*testdeep_test.MyMap{}"),
		})
	checkError(t, nil, testdeep.Map(MyMap{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("testdeep_test.MyMap{}"),
		})

	//
	// nil cases
	var (
		gotNilMap      map[string]int
		gotNilTypedMap MyMap
	)

	checkOK(t, gotNilMap, testdeep.Map(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, testdeep.Map(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, testdeep.Map(&MyMap{}, nil))

	//
	// SuperMapOf
	checkOK(t, gotMap, testdeep.SuperMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotMap,
		testdeep.SuperMapOf(map[string]int{"foo": 1}, testdeep.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		testdeep.SuperMapOf(map[string]int{}, testdeep.MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotMap,
		testdeep.SuperMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, testdeep.SuperMapOf(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test`),
		})
	checkError(t, gotMap,
		testdeep.SuperMapOf(map[string]int{}, testdeep.MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test`),
		})

	checkOK(t, gotNilMap, testdeep.SuperMapOf(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, testdeep.SuperMapOf(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, testdeep.SuperMapOf(&MyMap{}, nil))

	//
	// SubMapOf
	checkOK(t, gotMap,
		testdeep.SubMapOf(map[string]int{"foo": 1, "bar": 2, "tst": 3}, nil))
	checkOK(t, gotMap,
		testdeep.SubMapOf(map[string]int{"foo": 1, "tst": 3}, testdeep.MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		testdeep.SubMapOf(map[string]int{}, testdeep.MapEntries{"foo": 1, "bar": 2, "tst": 3}))

	checkError(t, gotMap,
		testdeep.SubMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	checkError(t, gotMap, testdeep.SubMapOf(map[string]int{"foo": 1}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra key: .*bar`),
		})
	checkError(t, gotMap,
		testdeep.SubMapOf(map[string]int{}, testdeep.MapEntries{"foo": 1, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing key: .*test(.|\n)*\n  Extra key: .*bar`),
		})

	checkOK(t, gotNilMap, testdeep.SubMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotNilTypedMap, testdeep.SubMapOf(MyMap{"foo": 1}, nil))
	checkOK(t, &gotNilTypedMap, testdeep.SubMapOf(&MyMap{"foo": 1}, nil))

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Map("test", nil) }, "usage: Map(")
	test.CheckPanic(t,
		func() { testdeep.SuperMapOf("test", nil) },
		"usage: SuperMapOf(")
	test.CheckPanic(t,
		func() { testdeep.SubMapOf("test", nil) },
		"usage: SubMapOf(")

	num := 12
	test.CheckPanic(t, func() { testdeep.Map(&num, nil) }, "usage: Map(")
	test.CheckPanic(t,
		func() { testdeep.SuperMapOf(&num, nil) },
		"usage: SuperMapOf(")
	test.CheckPanic(t,
		func() { testdeep.SubMapOf(&num, nil) },
		"usage: SubMapOf(")

	test.CheckPanic(t,
		func() { testdeep.Map(&MyMap{}, testdeep.MapEntries{1: 2}) },
		"expected key 1 type mismatch: int != model key type (string)")

	test.CheckPanic(t,
		func() { testdeep.Map(&MyMap{}, testdeep.MapEntries{"foo": nil}) },
		`expected key "foo" value cannot be nil as entries value type is int`)

	test.CheckPanic(t,
		func() { testdeep.Map(&MyMap{}, testdeep.MapEntries{"foo": uint16(2)}) },
		`expected key "foo" value type mismatch: uint16 != model key type (int)`)

	test.CheckPanic(t,
		func() { testdeep.Map(&MyMap{"foo": 1}, testdeep.MapEntries{"foo": 1}) },
		`"foo" entry exists in both model & expectedEntries`)

	//
	// String
	test.EqualStr(t, testdeep.Map(MyMap{}, nil).String(),
		"testdeep_test.MyMap{}")
	test.EqualStr(t, testdeep.Map(&MyMap{}, nil).String(),
		"*testdeep_test.MyMap{}")
	test.EqualStr(t, testdeep.Map(&MyMap{"foo": 2}, nil).String(),
		`*testdeep_test.MyMap{
  "foo": 2,
}`)

	test.EqualStr(t, testdeep.SubMapOf(MyMap{}, nil).String(),
		"SubMapOf(testdeep_test.MyMap{})")
	test.EqualStr(t, testdeep.SubMapOf(&MyMap{}, nil).String(),
		"SubMapOf(*testdeep_test.MyMap{})")
	test.EqualStr(t, testdeep.SubMapOf(&MyMap{"foo": 2}, nil).String(),
		`SubMapOf(*testdeep_test.MyMap{
  "foo": 2,
})`)

	test.EqualStr(t, testdeep.SuperMapOf(MyMap{}, nil).String(),
		"SuperMapOf(testdeep_test.MyMap{})")
	test.EqualStr(t, testdeep.SuperMapOf(&MyMap{}, nil).String(),
		"SuperMapOf(*testdeep_test.MyMap{})")
	test.EqualStr(t, testdeep.SuperMapOf(&MyMap{"foo": 2}, nil).String(),
		`SuperMapOf(*testdeep_test.MyMap{
  "foo": 2,
})`)
}

func TestMapTypeBehind(t *testing.T) {
	type MyMap map[string]int

	// Map
	equalTypes(t, testdeep.Map(map[string]int{}, nil), map[string]int{})
	equalTypes(t, testdeep.Map(MyMap{}, nil), MyMap{})
	equalTypes(t, testdeep.Map(&MyMap{}, nil), &MyMap{})

	// SubMap
	equalTypes(t, testdeep.SubMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, testdeep.SubMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, testdeep.SubMapOf(&MyMap{}, nil), &MyMap{})

	// SuperMap
	equalTypes(t, testdeep.SuperMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, testdeep.SuperMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, testdeep.SuperMapOf(&MyMap{}, nil), &MyMap{})
}
