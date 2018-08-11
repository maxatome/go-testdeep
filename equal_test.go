// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"testing"

	. "github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

type ComplexStruct struct { // nolint: megacheck
	ItemsByName map[string]*ComplexStructItem
	ItemsById   map[uint32]*ComplexStructItem
	Items       []*ComplexStructItem
	Label       string
	Weight      float64
}

type ComplexStructItem struct { // nolint: megacheck
	Name       string
	Id         uint32
	properties []ItemProperty
	propByName map[string]ItemProperty
	Enabled    bool
}

type ItemPropertyKind uint8

type ItemProperty struct {
	name  string
	kind  ItemPropertyKind
	value interface{}
}

//
// Array
func TestEqualArray(t *testing.T) {
	checkOK(t, [8]int{1, 2}, [8]int{1, 2})

	checkError(t, [8]int{1, 2}, [8]int{1, 3},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})

	oldMaxErrors := DefaultContextConfig.MaxErrors
	defer func() { DefaultContextConfig.MaxErrors = oldMaxErrors }()

	t.Run("DefaultContextConfig.MaxErrors = 2",
		func(t *testing.T) {
			DefaultContextConfig.MaxErrors = 2
			err := EqDeeplyError([8]int{1, 2, 3, 4}, [8]int{1, 42, 43, 44})

			// First error
			ok := t.Run("First error",
				func(t *testing.T) {
					if err == nil {
						t.Errorf("An Error should have occurred")
						return
					}
					if !matchError(t, err.(*ctxerr.Error), expectedError{
						Message:  mustBe("values differ"),
						Path:     mustBe("DATA[1]"),
						Got:      mustBe("(int) 2"),
						Expected: mustBe("(int) 42"),
					}, false) {
						return
					}
				})
			if !ok {
				return
			}

			// Second error
			eErr := err.(*ctxerr.Error).Next
			t.Run("Second error",
				func(t *testing.T) {
					if eErr == nil {
						t.Errorf("A second Error should have occurred")
						return
					}
					if !matchError(t, eErr, expectedError{
						Message:  mustBe("values differ"),
						Path:     mustBe("DATA[2]"),
						Got:      mustBe("(int) 3"),
						Expected: mustBe("(int) 43"),
					}, false) {
						return
					}
					if eErr.Next != ctxerr.ErrTooManyErrors {
						if eErr.Next == nil {
							t.Error("ErrTooManyErrors should follow the 2 errors")
						} else {
							t.Errorf("Only 2 Errors should have occurred. Found 3rd: %s",
								eErr.Next)
						}
						return
					}
				})
		})

	t.Run("DefaultContextConfig.MaxErrors = -1 (aka. all errors)",
		func(t *testing.T) {
			DefaultContextConfig.MaxErrors = -1
			err := EqDeeplyError([8]int{1, 2, 3, 4}, [8]int{1, 42, 43, 44})

			// First error
			ok := t.Run("First error",
				func(t *testing.T) {
					if err == nil {
						t.Errorf("An Error should have occurred")
						return
					}
					if !matchError(t, err.(*ctxerr.Error), expectedError{
						Message:  mustBe("values differ"),
						Path:     mustBe("DATA[1]"),
						Got:      mustBe("(int) 2"),
						Expected: mustBe("(int) 42"),
					}, false) {
						return
					}
				})
			if !ok {
				return
			}

			// Second error
			eErr := err.(*ctxerr.Error).Next
			ok = t.Run("Second error",
				func(t *testing.T) {
					if eErr == nil {
						t.Errorf("A second Error should have occurred")
						return
					}
					if !matchError(t, eErr, expectedError{
						Message:  mustBe("values differ"),
						Path:     mustBe("DATA[2]"),
						Got:      mustBe("(int) 3"),
						Expected: mustBe("(int) 43"),
					}, false) {
						return
					}
				})
			if !ok {
				return
			}

			// Third error
			eErr = eErr.Next
			t.Run("Third error",
				func(t *testing.T) {
					if eErr == nil {
						t.Errorf("A third Error should have occurred")
						return
					}
					if !matchError(t, eErr, expectedError{
						Message:  mustBe("values differ"),
						Path:     mustBe("DATA[3]"),
						Got:      mustBe("(int) 4"),
						Expected: mustBe("(int) 44"),
					}, false) {
						return
					}
					if eErr.Next != nil {
						t.Errorf("Only 3 Errors should have occurred")
						return
					}
				})
		})
}

//
// Slice
func TestEqualSlice(t *testing.T) {
	checkOK(t, []int{1, 2}, []int{1, 2})

	// Same pointer
	array := [2]int{1, 2}
	checkOK(t, array[:], array[:])
	checkOK(t, ([]int)(nil), ([]int)(nil))

	checkError(t, []int{1, 2}, []int{1, 2, 3},
		expectedError{
			Message: mustBe("comparing slices, from index #2"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing items: ((int) 3)`),
		})

	checkError(t, []int{1, 2, 3}, []int{1, 2},
		expectedError{
			Message: mustBe("comparing slices, from index #2"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Extra items: ((int) 3)`),
		})

	checkError(t, []int{1, 2}, ([]int)(nil),
		expectedError{
			Message:  mustBe("nil slice"),
			Path:     mustBe("DATA"),
			Got:      mustBe("not nil"),
			Expected: mustBe("nil"),
		})

	checkError(t, ([]int)(nil), []int{1, 2},
		expectedError{
			Message:  mustBe("nil slice"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	checkError(t, []int{1, 2}, []int{1, 3},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("(int) 2"),
			Expected: mustBe("(int) 3"),
		})
}

//
// Interface
func TestEqualInterface(t *testing.T) {
	checkOK(t, []interface{}{1, "foo"}, []interface{}{1, "foo"})
	checkOK(t, []interface{}{1, nil}, []interface{}{1, nil})

	checkError(t, []interface{}{1, nil}, []interface{}{1, "foo"},
		expectedError{
			Message:  mustBe("nil interface"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	checkError(t, []interface{}{1, "foo"}, []interface{}{1, nil},
		expectedError{
			Message:  mustBe("nil interface"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("not nil"),
			Expected: mustBe("nil"),
		})

	checkError(t, []interface{}{1, "foo"}, []interface{}{1, 12},
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("string"),
			Expected: mustBe("int"),
		})
}

//
// Ptr
func TestEqualPtr(t *testing.T) {
	expected := 12
	gotOK := expected
	gotBad := 13

	checkOK(t, &gotOK, &expected)
	checkOK(t, &expected, &expected) // Same pointer

	checkError(t, &gotBad, &expected,
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("*DATA"),
			Got:      mustBe("(int) 13"),
			Expected: mustBe("(int) 12"),
		})
}

//
// Struct
func TestEqualStruct(t *testing.T) {
	checkOK(t,
		ItemProperty{ // got
			name:  "foo",
			kind:  12,
			value: "bar",
		},
		ItemProperty{ // expected
			name:  "foo",
			kind:  12,
			value: "bar",
		})

	checkError(t,
		ItemProperty{ // got
			name:  "foo",
			kind:  12,
			value: 12,
		},
		ItemProperty{ // expected
			name:  "foo",
			kind:  12,
			value: "bar",
		},
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA.value"),
			Got:      mustBe("int"),
			Expected: mustBe("string"),
		})
}

//
// Map
func TestEqualMap(t *testing.T) {
	checkOK(t, map[string]int{}, map[string]int{})
	checkOK(t, (map[string]int)(nil), (map[string]int)(nil))

	expected := map[string]int{"foo": 1, "bar": 4}
	checkOK(t, map[string]int{"foo": 1, "bar": 4}, expected)
	checkOK(t, expected, expected) // Same pointer

	checkError(t, map[string]int{"foo": 1, "bar": 4}, (map[string]int)(nil),
		expectedError{
			Message:  mustBe("nil map"),
			Path:     mustBe("DATA"),
			Got:      mustBe("not nil"),
			Expected: mustBe("nil"),
		})

	checkError(t, (map[string]int)(nil), map[string]int{"foo": 1, "bar": 4},
		expectedError{
			Message:  mustBe("nil map"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4},
		map[string]int{"foo": 1, "bar": 5},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe(`DATA["bar"]`),
			Got:      mustBe("(int) 4"),
			Expected: mustBe("(int) 5"),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4, "test": 12},
		map[string]int{"foo": 1, "bar": 4},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`Extra keys:[^"]+"test"`),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4},
		map[string]int{"foo": 1, "bar": 4, "test": 12},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`Missing keys:[^"]+"test"`),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4, "test+": 12},
		map[string]int{"foo": 1, "bar": 4, "test-": 12},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`Missing keys:[^"]+"test-".*
  Extra keys:[^"]+"test\+"`),
		})
}

//
// Func
func TestEqualFunc(t *testing.T) {
	checkOK(t, (func())(nil), (func())(nil))

	checkError(t, func() {}, func() {},
		expectedError{
			Message: mustBe("functions mismatch"),
			Path:    mustBe("DATA"),
			Summary: mustBe("<can not be compared>"),
		})
}

//
// Others
func TestEqualOthers(t *testing.T) {
	type Private struct {
		num   int
		num8  int8
		num16 int16
		num32 int32
		num64 int64

		numu   uint
		numu8  uint8
		numu16 uint16
		numu32 uint32
		numu64 uint64

		numf32 float32
		numf64 float64

		numc64  complex64
		numc128 complex128

		boolean bool

		channel chan int
	}

	channel := make(chan int, 1)
	checkOK(t,
		Private{ // got
			num:     1,
			num8:    8,
			num16:   16,
			num32:   32,
			num64:   64,
			numu:    1,
			numu8:   8,
			numu16:  16,
			numu32:  32,
			numu64:  64,
			numf32:  32,
			numf64:  64,
			numc64:  complex(64, 1),
			numc128: complex(128, -1),
			boolean: true,
			channel: channel,
		},
		Private{
			num:     1,
			num8:    8,
			num16:   16,
			num32:   32,
			num64:   64,
			numu:    1,
			numu8:   8,
			numu16:  16,
			numu32:  32,
			numu64:  64,
			numf32:  32,
			numf64:  64,
			numc64:  complex(64, 1),
			numc128: complex(128, -1),
			boolean: true,
			channel: channel,
		})

	checkError(t, Private{num: 1}, Private{num: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num"),
			Got:      mustBe("(int) 1"),
			Expected: mustBe("(int) 2"),
		})

	checkError(t, Private{num8: 1}, Private{num8: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num8"),
			Got:      mustBe("(int8) 1"),
			Expected: mustBe("(int8) 2"),
		})

	checkError(t, Private{num16: 1}, Private{num16: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num16"),
			Got:      mustBe("(int16) 1"),
			Expected: mustBe("(int16) 2"),
		})

	checkError(t, Private{num32: 1}, Private{num32: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num32"),
			Got:      mustBe("(int32) 1"),
			Expected: mustBe("(int32) 2"),
		})

	checkError(t, Private{num64: 1}, Private{num64: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num64"),
			Got:      mustBe("(int64) 1"),
			Expected: mustBe("(int64) 2"),
		})

	checkError(t, Private{numu: 1}, Private{numu: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numu"),
			Got:      mustBe("(uint) 1"),
			Expected: mustBe("(uint) 2"),
		})

	checkError(t, Private{numu8: 1}, Private{numu8: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numu8"),
			Got:      mustBe("(uint8) 1"),
			Expected: mustBe("(uint8) 2"),
		})

	checkError(t, Private{numu16: 1}, Private{numu16: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numu16"),
			Got:      mustBe("(uint16) 1"),
			Expected: mustBe("(uint16) 2"),
		})

	checkError(t, Private{numu32: 1}, Private{numu32: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numu32"),
			Got:      mustBe("(uint32) 1"),
			Expected: mustBe("(uint32) 2"),
		})

	checkError(t, Private{numu64: 1}, Private{numu64: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numu64"),
			Got:      mustBe("(uint64) 1"),
			Expected: mustBe("(uint64) 2"),
		})

	checkError(t, Private{numf32: 1}, Private{numf32: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numf32"),
			Got:      mustBe("(float32) 1"),
			Expected: mustBe("(float32) 2"),
		})

	checkError(t, Private{numf64: 1}, Private{numf64: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numf64"),
			Got:      mustBe("(float64) 1"),
			Expected: mustBe("(float64) 2"),
		})

	checkError(t, Private{numc64: complex(1, 2)}, Private{numc64: complex(2, 1)},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numc64"),
			Got:      mustBe("(complex64) (1+2i)"),
			Expected: mustBe("(complex64) (2+1i)"),
		})

	checkError(t, Private{numc128: complex(1, 2)},
		Private{numc128: complex(2, 1)},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.numc128"),
			Got:      mustBe("(complex128) (1+2i)"),
			Expected: mustBe("(complex128) (2+1i)"),
		})

	checkError(t, Private{boolean: true}, Private{boolean: false},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.boolean"),
			Got:      mustBe("(bool) true"),
			Expected: mustBe("(bool) false"),
		})

	var expectedChannel = make(chan int, 2)
	checkError(t, Private{channel: channel},
		Private{channel: expectedChannel},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.channel"),
			Got:      mustContain("0x"), // hexadecimal pointer
			Expected: mustContain("0x"), // hexadecimal pointer
		})
}

func TestEqualRecurs(t *testing.T) {
	type S struct {
		Next *S
	}

	expected1 := &S{}
	expected1.Next = expected1

	got := &S{}
	got.Next = got

	expected2 := &S{}
	expected2.Next = expected2

	checkOK(t, got, expected1)
	checkOK(t, got, expected2)
}

func TestEqualPanic(t *testing.T) {
	test.CheckPanic(t,
		func() {
			EqDeeply(Ignore(), Ignore())
		},
		"Found a TestDeep operator in got param, can only use it in expected one!")
}

func TestCmpDeeply(t *testing.T) {
	mockT := &testing.T{}
	test.IsTrue(t, CmpDeeply(mockT, 1, 1))
	test.IsFalse(t, mockT.Failed())

	mockT = &testing.T{}
	test.IsFalse(t, CmpDeeply(mockT, 1, 2))
	test.IsTrue(t, mockT.Failed())

	mockT = &testing.T{}
	test.IsFalse(t, CmpDeeply(mockT, 1, 2, "Basic test"))
	test.IsTrue(t, mockT.Failed())

	mockT = &testing.T{}
	test.IsFalse(t, CmpDeeply(mockT, 1, 2, "Basic test with %d and %d", 1, 2))
	test.IsTrue(t, mockT.Failed())
}

func ExampleEqDeeply() {
	type MyStruct struct {
		Name  string
		Num   int
		Items []int
	}

	got := &MyStruct{
		Name:  "Foobar",
		Num:   12,
		Items: []int{4, 5, 9, 3, 8},
	}

	if EqDeeply(got,
		Struct(&MyStruct{},
			StructFields{
				"Name":  Re("^Foo"),
				"Num":   Between(10, 20),
				"Items": ArrayEach(Between(3, 9)),
			})) {
		fmt.Println("Match!")
	} else {
		fmt.Println("NO!")
	}

	// Output:
	// Match!
}

func ExampleEqDeeplyError() {
	type MyStruct struct {
		Name  string
		Num   int
		Items []int
	}

	got := &MyStruct{
		Name:  "Foobar",
		Num:   12,
		Items: []int{4, 5, 9, 3, 8},
	}

	err := EqDeeplyError(got,
		Struct(&MyStruct{},
			StructFields{
				"Name":  Re("^Foo"),
				"Num":   Between(10, 20),
				"Items": ArrayEach(Between(3, 8)),
			}))
	if err != nil {
		fmt.Println(err)
	}

	// Output something like:
	// DATA.Items[2]: values differ
	// 	     got: 9
	// 	expected: 3 ≤ got ≤ 8
	// [under TestDeep operator Between at equal_test.go:17]
}
