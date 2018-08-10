// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestMap(t *testing.T) {
	type MyMap map[string]int

	//
	// Map
	checkOK(t, (map[string]int)(nil), Map(map[string]int{}, nil))

	checkError(t, nil, Map(map[string]int{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA`),
			Got:      mustBe("nil"),
			Expected: mustBe("map[string]int{}"),
		})

	gotMap := map[string]int{"foo": 1, "bar": 2}

	checkOK(t, gotMap, Map(map[string]int{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotMap, Map(map[string]int{"foo": 1}, MapEntries{"bar": 2}))
	checkOK(t, gotMap, Map(map[string]int{}, MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotMap, Map((map[string]int)(nil), MapEntries{"foo": 1, "bar": 2}))

	one := 1
	checkOK(t, map[string]*int{"foo": nil, "bar": &one},
		Map(map[string]*int{}, MapEntries{"foo": nil, "bar": &one}))

	checkError(t, gotMap, Map(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	checkError(t, gotMap, Map(map[string]int{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap, Map(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap, Map(map[string]int{}, MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotMap,
		Map(map[string]int{}, MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test[^\n]+\z`),
		})
	checkError(t, gotMap, Map(MyMap{}, MapEntries{"foo": 1, "bar": 2}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("testdeep_test.MyMap"),
		})

	//
	// Map type
	gotTypedMap := MyMap{"foo": 1, "bar": 2}

	checkOK(t, gotTypedMap, Map(MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, gotTypedMap, Map(MyMap{"foo": 1}, MapEntries{"bar": 2}))
	checkOK(t, gotTypedMap, Map(MyMap{}, MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, gotTypedMap, Map((MyMap)(nil), MapEntries{"foo": 1, "bar": 2}))

	checkOK(t, &gotTypedMap, Map(&MyMap{"foo": 1, "bar": 2}, nil))
	checkOK(t, &gotTypedMap, Map(&MyMap{"foo": 1}, MapEntries{"bar": 2}))
	checkOK(t, &gotTypedMap, Map(&MyMap{}, MapEntries{"foo": 1, "bar": 2}))
	checkOK(t, &gotTypedMap, Map((*MyMap)(nil), MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotTypedMap, Map(MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	checkError(t, gotTypedMap, Map(MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap, Map(MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap, Map(MyMap{}, MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, gotTypedMap,
		Map(MyMap{}, MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test[^\n]+\z`),
		})

	checkError(t, &gotTypedMap, Map(&MyMap{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	checkError(t, &gotTypedMap, Map(&MyMap{}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap, Map(&MyMap{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap, Map(&MyMap{}, MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(
				`^Missing keys: .*test(.|\n)*\n  Extra keys: .*(foo(.|\n)*bar|bar(.|\n)*foo)`),
		})
	checkError(t, &gotTypedMap,
		Map(&MyMap{}, MapEntries{"foo": 1, "bar": 2, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test[^\n]+\z`),
		})

	checkError(t, &gotMap, Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*map[string]int"),
			Expected: mustBe("*testdeep_test.MyMap"),
		})
	checkError(t, gotMap, Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("map[string]int"),
			Expected: mustBe("*testdeep_test.MyMap"),
		})
	checkError(t, nil, Map(&MyMap{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("*testdeep_test.MyMap{}"),
		})
	checkError(t, nil, Map(MyMap{}, nil),
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

	checkOK(t, gotNilMap, Map(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, Map(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, Map(&MyMap{}, nil))

	//
	// SuperMapOf
	checkOK(t, gotMap, SuperMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotMap,
		SuperMapOf(map[string]int{"foo": 1}, MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		SuperMapOf(map[string]int{}, MapEntries{"foo": 1, "bar": 2}))

	checkError(t, gotMap, SuperMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	checkError(t, gotMap, SuperMapOf(map[string]int{"test": 2}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test`),
		})
	checkError(t, gotMap, SuperMapOf(map[string]int{}, MapEntries{"test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test`),
		})

	checkOK(t, gotNilMap, SuperMapOf(map[string]int{}, nil))
	checkOK(t, gotNilTypedMap, SuperMapOf(MyMap{}, nil))
	checkOK(t, &gotNilTypedMap, SuperMapOf(&MyMap{}, nil))

	//
	// SubMapOf
	checkOK(t, gotMap,
		SubMapOf(map[string]int{"foo": 1, "bar": 2, "tst": 3}, nil))
	checkOK(t, gotMap,
		SubMapOf(map[string]int{"foo": 1, "tst": 3}, MapEntries{"bar": 2}))
	checkOK(t, gotMap,
		SubMapOf(map[string]int{}, MapEntries{"foo": 1, "bar": 2, "tst": 3}))

	checkError(t, gotMap, SubMapOf(map[string]int{"foo": 1, "bar": 3}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	checkError(t, gotMap, SubMapOf(map[string]int{"foo": 1}, nil),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Extra keys: .*bar`),
		})
	checkError(t, gotMap,
		SubMapOf(map[string]int{}, MapEntries{"foo": 1, "test": 2}),
		expectedError{
			Message: mustBe("comparing hash keys of %%"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`^Missing keys: .*test(.|\n)*\n  Extra keys: .*bar`),
		})

	checkOK(t, gotNilMap, SubMapOf(map[string]int{"foo": 1}, nil))
	checkOK(t, gotNilTypedMap, SubMapOf(MyMap{"foo": 1}, nil))
	checkOK(t, &gotNilTypedMap, SubMapOf(&MyMap{"foo": 1}, nil))

	//
	// Bad usage
	checkPanic(t, func() { Map("test", nil) }, "usage: Map(")
	checkPanic(t, func() { SuperMapOf("test", nil) }, "usage: SuperMapOf(")
	checkPanic(t, func() { SubMapOf("test", nil) }, "usage: SubMapOf(")

	num := 12
	checkPanic(t, func() { Map(&num, nil) }, "usage: Map(")
	checkPanic(t, func() { SuperMapOf(&num, nil) }, "usage: SuperMapOf(")
	checkPanic(t, func() { SubMapOf(&num, nil) }, "usage: SubMapOf(")

	checkPanic(t, func() { Map(&MyMap{}, MapEntries{1: 2}) },
		"expected key (int) 1 type mismatch: int != model key type (string)")

	checkPanic(t, func() { Map(&MyMap{}, MapEntries{"foo": nil}) },
		`expected key "foo" value cannot be nil as entries value type is int`)

	checkPanic(t, func() { Map(&MyMap{}, MapEntries{"foo": uint16(2)}) },
		`expected key "foo" value type mismatch: uint16 != model key type (int)`)

	checkPanic(t, func() { Map(&MyMap{"foo": 1}, MapEntries{"foo": 1}) },
		`"foo" entry exists in both model & expectedEntries`)

	//
	// String
	test.EqualStr(t, Map(MyMap{}, nil).String(), "testdeep_test.MyMap{}")
	test.EqualStr(t, Map(&MyMap{}, nil).String(), "*testdeep_test.MyMap{}")
	test.EqualStr(t, Map(&MyMap{"foo": 2}, nil).String(),
		`*testdeep_test.MyMap{
  "foo": (int) 2,
}`)

	test.EqualStr(t, SubMapOf(MyMap{}, nil).String(),
		"SubMapOf(testdeep_test.MyMap{})")
	test.EqualStr(t, SubMapOf(&MyMap{}, nil).String(),
		"SubMapOf(*testdeep_test.MyMap{})")
	test.EqualStr(t, SubMapOf(&MyMap{"foo": 2}, nil).String(),
		`SubMapOf(*testdeep_test.MyMap{
  "foo": (int) 2,
})`)

	test.EqualStr(t, SuperMapOf(MyMap{}, nil).String(),
		"SuperMapOf(testdeep_test.MyMap{})")
	test.EqualStr(t, SuperMapOf(&MyMap{}, nil).String(),
		"SuperMapOf(*testdeep_test.MyMap{})")
	test.EqualStr(t, SuperMapOf(&MyMap{"foo": 2}, nil).String(),
		`SuperMapOf(*testdeep_test.MyMap{
  "foo": (int) 2,
})`)
}

func TestMapTypeBehind(t *testing.T) {
	type MyMap map[string]int

	// Map
	equalTypes(t, Map(map[string]int{}, nil), map[string]int{})
	equalTypes(t, Map(MyMap{}, nil), MyMap{})
	equalTypes(t, Map(&MyMap{}, nil), &MyMap{})

	// SubMap
	equalTypes(t, SubMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, SubMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, SubMapOf(&MyMap{}, nil), &MyMap{})

	// SuperMap
	equalTypes(t, SuperMapOf(map[string]int{}, nil), map[string]int{})
	equalTypes(t, SuperMapOf(MyMap{}, nil), MyMap{})
	equalTypes(t, SuperMapOf(&MyMap{}, nil), &MyMap{})
}
