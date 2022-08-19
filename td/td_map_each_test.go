// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestMapEach(t *testing.T) {
	type MyMap map[string]int

	checkOKForEach(t,
		[]any{
			map[string]int{"foo": 1, "bar": 1},
			&map[string]int{"foo": 1, "bar": 1},
			MyMap{"foo": 1, "bar": 1},
			&MyMap{"foo": 1, "bar": 1},
		},
		td.MapEach(1))

	checkOKForEach(t,
		[]any{
			map[int]string{1: "foo", 2: "bar"},
			&map[int]string{1: "foo", 2: "bar"},
		},
		td.MapEach(td.Len(3)))

	checkOKForEach(t,
		[]any{
			map[string]int{},
			&map[string]int{},
			MyMap{},
			&MyMap{},
		},
		td.MapEach(1))

	checkOK(t, (map[string]int)(nil), td.MapEach(1))
	checkOK(t, (MyMap)(nil), td.MapEach(1))
	checkError(t, (*MyMap)(nil), td.MapEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *map (*td_test.MyMap type)"),
			Expected: mustBe("non-nil *map"),
		})

	checkOKForEach(t,
		[]any{
			map[string]int{"foo": 20, "bar": 22, "test": 29},
			&map[string]int{"foo": 20, "bar": 22, "test": 29},
			MyMap{"foo": 20, "bar": 22, "test": 29},
			&MyMap{"foo": 20, "bar": 22, "test": 29},
		},
		td.MapEach(td.Between(20, 30)))

	checkError(t, nil, td.MapEach(4),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("map OR *map"),
		})

	checkErrorForEach(t,
		[]any{
			map[string]int{"foo": 4, "bar": 5, "test": 4},
			&map[string]int{"foo": 4, "bar": 5, "test": 4},
			MyMap{"foo": 4, "bar": 5, "test": 4},
			&MyMap{"foo": 4, "bar": 5, "test": 4},
		},
		td.MapEach(4),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("5"),
			Expected: mustBe("4"),
		})

	checkError(t, 666, td.MapEach(4),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("map OR *map"),
		})
	num := 666
	checkError(t, &num, td.MapEach(4),
		expectedError{
			Message:  mustBe("bad kind"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("map OR *map"),
		})

	checkOK(t, map[string]any{"a": nil, "b": nil, "c": nil},
		td.MapEach(nil))
	checkError(t,
		map[string]any{"a": nil, "b": nil, "c": nil, "d": 66},
		td.MapEach(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["d"]`),
			Got:      mustBe("66"),
			Expected: mustBe("nil"),
		})

	//
	// String
	test.EqualStr(t, td.MapEach(4).String(), "MapEach(4)")
	test.EqualStr(t, td.MapEach(td.All(1, 2)).String(),
		`MapEach(All(1,
            2))`)
}

func TestMapEachTypeBehind(t *testing.T) {
	equalTypes(t, td.MapEach(4), nil)
}
