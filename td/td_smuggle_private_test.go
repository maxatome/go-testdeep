// Copyright (c) 2021-2024, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestFieldsPath(t *testing.T) {
	check := func(in string, expected ...string) []smuggleField {
		t.Helper()

		got, err := splitFieldsPath(in)
		test.NoError(t, err)

		var gotStr []string
		for _, s := range got {
			if s.Method {
				gotStr = append(gotStr, s.Name+"()")
			} else {
				gotStr = append(gotStr, s.Name)
			}
		}

		if !reflect.DeepEqual(gotStr, expected) {
			t.Errorf("Failed:\n       got: %v\n  expected: %v", got, expected)
		}

		test.EqualStr(t, in, joinFieldsPath(got))

		return got
	}

	check("test", "test")
	check("Test.Foo().bar", "Test", "Foo()", "bar")
	check("test.foo.bar", "test", "foo", "bar")
	check("test[foo.bar]", "test", "foo.bar")
	check("test[foo][bar]", "test", "foo", "bar")
	fp := check("test[foo][bar].zip", "test", "foo", "bar", "zip")

	// "." can be omitted just after "]"
	got, err := splitFieldsPath("test[foo][bar]zip")
	test.NoError(t, err)
	if !reflect.DeepEqual(got, fp) {
		t.Errorf("Failed:\n       got: %v\n  expected: %v", got, fp)
	}

	check("[foo][bar]", "foo", "bar")
	check("[0][foo][bar]", "0", "foo", "bar")

	//
	// Errors
	checkErr := func(in, expectedErr string) {
		t.Helper()

		_, err := splitFieldsPath(in)

		if test.Error(t, err) {
			test.EqualStr(t, err.Error(), expectedErr)
		}
	}

	checkErr("", "FIELDS_PATH cannot be empty")
	checkErr(".test", `'.' cannot be the first rune in FIELDS_PATH ".test"`)
	checkErr("foo.bar.", `final '.' in FIELDS_PATH "foo.bar." is not allowed`)
	checkErr("foo..bar", `unexpected '.' after '.' in FIELDS_PATH "foo..bar"`)
	checkErr("foo.[bar]", `unexpected '[' after '.' in FIELDS_PATH "foo.[bar]"`)
	checkErr("foo[bar", `cannot find final ']' in FIELDS_PATH "foo[bar"`)
	checkErr("test.%foo", `unexpected '%' in field name "%foo" in FIELDS_PATH "test.%foo"`)
	checkErr("test.f%oo", `unexpected '%' in field name "f%oo" in FIELDS_PATH "test.f%oo"`)
	checkErr("Foo().()", `missing method name before () in FIELDS_PATH "Foo().()"`)
	checkErr("abc.foo()", `method name "foo()" is not public in FIELDS_PATH "abc.foo()"`)
	checkErr("Fo%o().abc", `unexpected '%' in method name "Fo%o()" in FIELDS_PATH "Fo%o().abc"`)
	checkErr("Pipo.bingo.zzz.Foo.Zip().abc", `cannot call method Zip() as it is based on an unexported field "bingo" in FIELDS_PATH "Pipo.bingo.zzz.Foo.Zip().abc"`)
	checkErr("foo[bar", `cannot find final ']' in FIELDS_PATH "foo[bar"`)
}

type SmuggleBuild struct {
	Field struct {
		Path string
	}
	Iface any
	Next  *SmuggleBuild
}

func (s SmuggleBuild) FollowIface() any {
	return s.Iface
}

func (s *SmuggleBuild) PtrFollowIface() any {
	return s.Iface
}

func (s SmuggleBuild) MayFollowIface() (any, error) {
	if s.Iface == nil {
		return nil, errors.New("Iface is nil")
	}
	return s.Iface, nil
}

func (s SmuggleBuild) FollowNext() *SmuggleBuild {
	return s.Next
}

func (s *SmuggleBuild) PtrFollowNext() *SmuggleBuild {
	return s.Next
}

func (s SmuggleBuild) MayFollowNext() (*SmuggleBuild, error) {
	if s.Next == nil {
		return nil, errors.New("Next is nil")
	}
	return s.Next, nil
}

func (s SmuggleBuild) Error() (bool, error) {
	return false, errors.New("an error occurred")
}

func (s SmuggleBuild) SetPath(path string) {
	s.Field.Path = path
}

func (s SmuggleBuild) Panic() string {
	panic("oops!")
}

func (s SmuggleBuild) Num() int {
	return 42
}

func (s *SmuggleBuild) PNum() int {
	return 42
}

func TestBuildFieldsPathFn(t *testing.T) {
	_, err := buildFieldsPathFn("bad[path")
	test.Error(t, err)

	t.Run("Struct", func(t *testing.T) {
		fn, err := buildFieldsPathFn("Field.Path.Bad")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Field.Path" is a string and should be a struct or a map[string]…`)
			}

			_, err = fn(123)
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					"it is a int and should be a struct or a map[string]…")
			}
		}

		fn, err = buildFieldsPathFn("Iface.Bad")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: 42})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface" is a int (after dereferencing) and should be a struct or a map[string]…`)
			}

			num := 42
			_, err = fn(&num)
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`it is a int (after dereferencing) and should be a struct or a map[string]…`)
			}
		}

		fn, err = buildFieldsPathFn("Field.Unknown")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(), `field "Field.Unknown" not found`)
			}
		}
	})

	t.Run("Map", func(t *testing.T) {
		fn, err := buildFieldsPathFn("Iface[str].Field")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: map[int]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" is not an integer and so cannot match int map key type`)
			}

			_, err = fn(SmuggleBuild{Iface: map[uint]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" is not an unsigned integer and so cannot match uint map key type`)
			}

			_, err = fn(SmuggleBuild{Iface: map[float32]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" is not a float and so cannot match float32 map key type`)
			}

			_, err = fn(SmuggleBuild{Iface: map[complex128]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" is not a complex number and so cannot match complex128 map key type`)
			}

			_, err = fn(SmuggleBuild{Iface: map[struct{ A int }]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" cannot match unsupported struct { A int } map key type`)
			}

			_, err = fn(SmuggleBuild{Iface: map[string]SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(), `field "Iface[str]", "str" map key not found`)
			}
		}
	})

	t.Run("Array-Slice", func(t *testing.T) {
		fn, err := buildFieldsPathFn("Iface[str].Field")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: []int{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[str]", "str" is not a slice/array index`)
			}
		}

		fn, err = buildFieldsPathFn("Iface[18].Field")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: []int{1, 2, 3}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface[18]", 18 is out of slice/array range (len 3)`)
			}

			_, err = fn(SmuggleBuild{Iface: 42})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field "Iface" is a int, but a map, array or slice is expected`)
			}
		}

		fn, err = buildFieldsPathFn("[18].Field")
		if test.NoError(t, err) {
			_, err = fn(42)
			test.EqualStr(t, err.Error(),
				`it is a int, but a map, array or slice is expected`)
		}
	})

	t.Run("Function", func(t *testing.T) {
		fn, err := buildFieldsPathFn("Iface.Unknown()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`field Iface (type td.SmuggleBuild) does not implement Unknown() method`)
			}
		}

		fn, err = buildFieldsPathFn("Iface.NilUnknown()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(), `field "Iface" is nil`)
			}
		}

		fn, err = buildFieldsPathFn("Unknown()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`type td.SmuggleBuild has no method Unknown()`)
			}
		}

		fn, err = buildFieldsPathFn("Iface.SetPath()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: &SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`cannot call Iface.SetPath(), signature func(string) not handled, only func() A or func() (A, error) allowed`)
			}
		}

		fn, err = buildFieldsPathFn("Iface.Panic()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: &SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`method Iface.Panic() panicked: oops!`)
			}
		}

		fn, err = buildFieldsPathFn("Iface.Error()")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: &SmuggleBuild{}})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`method Iface.Error() returned an error: an error occurred`)
			}
		}

		fn, err = buildFieldsPathFn("Num().Bad")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`method Num() returned a int and should be a struct or a map[string]…`)
			}
		}

		fn, err = buildFieldsPathFn("FollowIface().Bad")
		if test.NoError(t, err) {
			_, err = fn(SmuggleBuild{Iface: 42})
			if test.Error(t, err) {
				test.EqualStr(t, err.Error(),
					`method FollowIface() returned a int (after dereferencing) and should be a struct or a map[string]…`)
			}
		}
	})
}
