// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestMapEach(t *testing.T) {
	type MyMap map[string]int

	checkOKForEach(t,
		[]interface{}{
			map[string]int{"foo": 1, "bar": 1},
			&map[string]int{"foo": 1, "bar": 1},
			MyMap{"foo": 1, "bar": 1},
			&MyMap{"foo": 1, "bar": 1},
		},
		MapEach(1))

	checkOKForEach(t,
		[]interface{}{
			map[string]int{},
			&map[string]int{},
			MyMap{},
			&MyMap{},
		},
		MapEach(1))

	checkOK(t, (map[string]int)(nil), MapEach(1))
	checkOK(t, (MyMap)(nil), MapEach(1))
	checkError(t, (*MyMap)(nil), MapEach(4),
		expectedError{
			Message:  mustBe("nil pointer"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil *testdeep_test.MyMap"),
			Expected: mustBe("Map OR *Map"),
		})

	checkOKForEach(t,
		[]interface{}{
			map[string]int{"foo": 20, "bar": 22, "test": 29},
			&map[string]int{"foo": 20, "bar": 22, "test": 29},
			MyMap{"foo": 20, "bar": 22, "test": 29},
			&MyMap{"foo": 20, "bar": 22, "test": 29},
		},
		MapEach(Between(20, 30)))

	checkError(t, nil, MapEach(4), expectedError{
		Message:  mustBe("nil value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("Map OR *Map"),
	})

	checkErrorForEach(t,
		[]interface{}{
			map[string]int{"foo": 4, "bar": 5, "test": 4},
			&map[string]int{"foo": 4, "bar": 5, "test": 4},
			MyMap{"foo": 4, "bar": 5, "test": 4},
			&MyMap{"foo": 4, "bar": 5, "test": 4},
		},
		MapEach(4),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA[(string) (len=3) "bar"]`),
			Got:      mustBe("(int) 5"),
			Expected: mustBe("(int) 4"),
		})

	checkError(t, 666, MapEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("Map OR *Map"),
		})
	num := 666
	checkError(t, &num, MapEach(4),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*int"),
			Expected: mustBe("Map OR *Map"),
		})

	checkOK(t, map[string]interface{}{"a": nil, "b": nil, "c": nil}, MapEach(nil))
	checkError(t, map[string]interface{}{"a": nil, "b": nil, "c": nil, "d": 66},
		MapEach(nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA[(string) (len=1) "d"]`),
			Got:      mustBe("(int) 66"),
			Expected: mustBe("nil"),
		})

	//
	// String
	equalStr(t, MapEach(4).String(), "MapEach((int) 4)")
	equalStr(t, MapEach(All(1, 2)).String(),
		`MapEach(All((int) 1,
            (int) 2))`)
}

func TestMapEachTypeBehind(t *testing.T) {
	equalTypes(t, MapEach(4), nil)
}
