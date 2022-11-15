// Copyright (c) 2018-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type ItemPropertyKind uint8

type ItemProperty struct {
	name  string
	kind  ItemPropertyKind
	value any
}

// Array.
func TestEqualArray(t *testing.T) {
	checkOK(t, [8]int{1, 2}, [8]int{1, 2})

	checkError(t, [8]int{1, 2}, [8]int{1, 3},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})

	oldMaxErrors := td.DefaultContextConfig.MaxErrors
	defer func() { td.DefaultContextConfig.MaxErrors = oldMaxErrors }()

	t.Run("DefaultContextConfig.MaxErrors = 2",
		func(t *testing.T) {
			td.DefaultContextConfig.MaxErrors = 2
			err := td.EqDeeplyError([8]int{1, 2, 3, 4}, [8]int{1, 42, 43, 44})

			// First error
			ok := t.Run("First error",
				func(t *testing.T) {
					if err == nil {
						t.Errorf("An Error should have occurred")
						return
					}
					if !matchError(t, err.(*ctxerr.Error),
						expectedError{
							Message:  mustBe("values differ"),
							Path:     mustBe("DATA[1]"),
							Got:      mustBe("2"),
							Expected: mustBe("42"),
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
					if !matchError(t, eErr,
						expectedError{
							Message:  mustBe("values differ"),
							Path:     mustBe("DATA[2]"),
							Got:      mustBe("3"),
							Expected: mustBe("43"),
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

	t.Run("DefaultContextConfig.MaxErrors = -1 (aka all errors)",
		func(t *testing.T) {
			td.DefaultContextConfig.MaxErrors = -1
			err := td.EqDeeplyError([8]int{1, 2, 3, 4}, [8]int{1, 42, 43, 44})

			// First error
			ok := t.Run("First error",
				func(t *testing.T) {
					if err == nil {
						t.Errorf("An Error should have occurred")
						return
					}
					if !matchError(t, err.(*ctxerr.Error),
						expectedError{
							Message:  mustBe("values differ"),
							Path:     mustBe("DATA[1]"),
							Got:      mustBe("2"),
							Expected: mustBe("42"),
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
					if !matchError(t, eErr,
						expectedError{
							Message:  mustBe("values differ"),
							Path:     mustBe("DATA[2]"),
							Got:      mustBe("3"),
							Expected: mustBe("43"),
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
					if !matchError(t, eErr,
						expectedError{
							Message:  mustBe("values differ"),
							Path:     mustBe("DATA[3]"),
							Got:      mustBe("4"),
							Expected: mustBe("44"),
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

// Slice.
func TestEqualSlice(t *testing.T) {
	checkOK(t, []int{1, 2}, []int{1, 2})

	// Same pointer
	array := [...]int{2, 1, 4, 3}
	checkOK(t, array[:], array[:])
	checkOK(t, ([]int)(nil), ([]int)(nil))

	// Same pointer, but not same len
	checkError(t, array[:2], array[:],
		expectedError{
			Message: mustBe("comparing slices, from index #2"),
			Path:    mustBe("DATA"),
			// Missing items are not sorted
			Summary: mustBe(`Missing 2 items: (4,
                  3)`),
		})

	checkError(t, []int{1, 2}, []int{1, 2, 3},
		expectedError{
			Message: mustBe("comparing slices, from index #2"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing item: (3)`),
		})

	checkError(t, []int{1, 2, 3}, []int{1, 2},
		expectedError{
			Message: mustBe("comparing slices, from index #2"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Extra item: (3)`),
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
			Got:      mustBe("2"),
			Expected: mustBe("3"),
		})
}

// Interface.
func TestEqualInterface(t *testing.T) {
	checkOK(t, []any{1, "foo"}, []any{1, "foo"})
	checkOK(t, []any{1, nil}, []any{1, nil})

	checkError(t, []any{1, nil}, []any{1, "foo"},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("nil"),
			Expected: mustBe(`"foo"`),
		})

	checkError(t, []any{1, "foo"}, []any{1, nil},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe(`"foo"`),
			Expected: mustBe("nil"),
		})

	checkError(t, []any{1, "foo"}, []any{1, 12},
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA[1]"),
			Got:      mustBe("string"),
			Expected: mustBe("int"),
		})
}

// Ptr.
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
			Got:      mustBe("13"),
			Expected: mustBe("12"),
		})
}

// Struct.
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

	type SType struct {
		Public  int
		private string
	}

	checkOK(t,
		SType{Public: 42, private: "test"},
		SType{Public: 42, private: "test"})

	checkError(t,
		SType{Public: 42, private: "test"},
		SType{Public: 42},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.private"),
			Got:      mustBe(`"test"`),
			Expected: mustBe(`""`),
		})

	defer func() { td.DefaultContextConfig.IgnoreUnexported = false }()
	td.DefaultContextConfig.IgnoreUnexported = true

	checkOK(t,
		SType{Public: 42, private: "test"},
		SType{Public: 42})

	// Be careful with structs containing only private fields
	checkOK(t,
		ItemProperty{
			name:  "foo",
			kind:  12,
			value: "bar",
		},
		ItemProperty{})
}

// Map.
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
			Got:      mustBe("4"),
			Expected: mustBe("5"),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4, "test": 12},
		map[string]int{"foo": 1, "bar": 4},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`Extra key:[^"]+"test"`),
		})

	checkError(t, map[string]int{"foo": 1, "bar": 4},
		map[string]int{"foo": 1, "bar": 4, "test": 12},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustMatch(`Missing key:[^"]+"test"`),
		})

	// Extra and missing keys are sorted
	checkError(t, map[string]int{"foo": 1, "bar": 4, "test1+": 12, "test2+": 13},
		map[string]int{"foo": 1, "bar": 4, "test1-": 12, "test2-": 13},
		expectedError{
			Message: mustBe("comparing map"),
			Path:    mustBe("DATA"),
			Summary: mustBe(`Missing 2 keys: ("test1-",
                 "test2-")
  Extra 2 keys: ("test1+",
                 "test2+")`),
		})
}

// Func.
func TestEqualFunc(t *testing.T) {
	checkOK(t, (func())(nil), (func())(nil))

	checkError(t, func() {}, func() {},
		expectedError{
			Message: mustBe("functions mismatch"),
			Path:    mustBe("DATA"),
			Summary: mustBe("<can not be compared>"),
		})
}

// Channel.
func TestEqualChannel(t *testing.T) {
	var gotCh, expectedCh chan int

	checkOK(t, gotCh, expectedCh) // nil channels

	gotCh = make(chan int, 1)
	checkOK(t, gotCh, gotCh) // exactly the same

	checkError(t, gotCh, make(chan int, 1),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("0x"), // hexadecimal pointer
			Expected: mustContain("0x"), // hexadecimal pointer
		})
}

// Others.
func TestEqualOthers(t *testing.T) {
	type Private struct { //nolint: maligned
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
	}

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
		})

	checkError(t, Private{num: 1}, Private{num: 2},
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.num"),
			Got:      mustBe("1"),
			Expected: mustBe("2"),
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
			Got:      mustBe("1.0"),
			Expected: mustBe("2.0"),
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
			Got:      mustBe("true"),
			Expected: mustBe("false"),
		})
}

// Private non-copyable fields.
func TestEqualReallyPrivate(t *testing.T) {
	type Private struct {
		channel chan int
	}

	ch := make(chan int, 3)

	checkOKOrPanicIfUnsafeDisabled(t, Private{channel: ch}, Private{channel: ch})
}

func TestEqualRecursPtr(t *testing.T) {
	type S struct {
		Next *S
		OK   bool
	}

	expected1 := &S{}
	expected1.Next = expected1

	got := &S{}
	got.Next = got

	expected2 := &S{}
	expected2.Next = expected2

	checkOK(t, got, expected1)
	checkOK(t, got, expected2)

	got.Next = &S{OK: true}
	expected1.Next = &S{OK: false}
	checkError(t, got, expected1,
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.Next.OK"),
			Got:      mustBe("true"),
			Expected: mustBe("false"),
		})
}

func TestEqualRecursMap(t *testing.T) { // issue #101
	gen := func() any {
		type S struct {
			Map map[int]S
		}

		m := make(map[int]S)
		m[1] = S{
			Map: m,
		}

		return m
	}

	checkOK(t, gen(), gen())
}

func TestEqualPanic(t *testing.T) {
	test.CheckPanic(t,
		func() {
			td.EqDeeply(td.Ignore(), td.Ignore())
		},
		"Found a TestDeep operator in got param, can only use it in expected one!")

	type tdInside struct {
		Operator td.TestDeep
	}
	test.CheckPanic(t,
		func() {
			td.EqDeeply(&tdInside{}, &tdInside{})
		},
		"Found a TestDeep operator in got param, can only use it in expected one!")

	t.Cleanup(func() { td.DefaultContextConfig.TestDeepInGotOK = false })
	td.DefaultContextConfig.TestDeepInGotOK = true

	test.IsTrue(t, td.EqDeeply(td.Ignore(), td.Ignore()))
	test.IsTrue(t, td.EqDeeply(&tdInside{}, &tdInside{}))
}

type AssignableType1 struct{ x, Ignore int }

func (a AssignableType1) Equal(b AssignableType1) bool {
	return a.x == b.x
}

type AssignableType2 struct{ x, Ignore int }

func (a AssignableType2) Equal(b struct{ x, Ignore int }) bool {
	return a.x == b.x
}

type AssignablePtrType3 struct{ x, Ignore int }

func (a *AssignablePtrType3) Equal(b *AssignablePtrType3) bool {
	if a == nil {
		return b == nil
	}
	return b != nil && a.x == b.x
}

type BadEqual1 int

func (b BadEqual1) Equal(o ...BadEqual1) bool { return true } // IsVariadic

type BadEqual2 int

func (b BadEqual2) Equal() bool { return true } // NumIn() ≠ 2

type BadEqual3 int

func (b BadEqual3) Equal(o BadEqual3) (int, int) { return 1, 2 } // NumOut() ≠ 1

type BadEqual4 int

func (b BadEqual4) Equal(o string) int { return 1 } // !AssignableTo

type BadEqual5 int

func (b BadEqual5) Equal(o BadEqual5) int { return 1 } // Out=bool

func TestUseEqualGlobal(t *testing.T) {
	defer func() { td.DefaultContextConfig.UseEqual = false }()
	td.DefaultContextConfig.UseEqual = true

	// Real case with time.Time
	time1 := time.Now()
	time2 := time1.Truncate(0)
	if !time1.Equal(time2) || !time2.Equal(time1) {
		t.Fatal("time.Equal() does not work as expected")
	}

	checkOK(t, time1, time2)
	checkOK(t, time2, time1)

	// AssignableType1
	a1 := AssignableType1{x: 13, Ignore: 666}
	b1 := AssignableType1{x: 13, Ignore: 789}
	checkOK(t, a1, b1)
	checkOK(t, b1, a1)
	checkError(t, a1, AssignableType1{x: 14, Ignore: 666},
		expectedError{
			Message:  mustBe("got.Equal(expected) failed"),
			Path:     mustBe("DATA"),
			Got:      mustContain("x: (int) 13,"),
			Expected: mustContain("x: (int) 14,"),
		})

	bs := struct{ x, Ignore int }{x: 13, Ignore: 789}
	checkOK(t, a1, bs) // bs type is assignable to AssignableType1
	checkError(t, bs, a1,
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("struct { x int; Ignore int }"),
			Expected: mustBe("td_test.AssignableType1"),
		})

	// AssignableType2
	a2 := AssignableType2{x: 13, Ignore: 666}
	b2 := AssignableType2{x: 13, Ignore: 789}
	checkOK(t, a2, b2)
	checkOK(t, b2, a2)
	checkOK(t, a2, bs) // bs type is assignable to AssignableType2
	checkError(t, bs, a2,
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("struct { x int; Ignore int }"),
			Expected: mustBe("td_test.AssignableType2"),
		})

	// AssignablePtrType3
	a3 := &AssignablePtrType3{x: 13, Ignore: 666}
	b3 := &AssignablePtrType3{x: 13, Ignore: 789}
	checkOK(t, a3, b3)
	checkOK(t, b3, a3)
	checkError(t, a3, &bs, // &bs type not assignable to AssignablePtrType3
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.AssignablePtrType3"),
			Expected: mustBe("*struct { x int; Ignore int }"),
		})
	checkOK(t, (*AssignablePtrType3)(nil), (*AssignablePtrType3)(nil))
	checkError(t, (*AssignablePtrType3)(nil), b3,
		expectedError{
			Message:  mustBe("got.Equal(expected) failed"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(*td_test.AssignablePtrType3)(<nil>)"),
			Expected: mustContain("x: (int) 13,"),
		})
	checkError(t, b3, (*AssignablePtrType3)(nil),
		expectedError{
			Message:  mustBe("got.Equal(expected) failed"),
			Path:     mustBe("DATA"),
			Got:      mustContain("x: (int) 13,"),
			Expected: mustBe("(*td_test.AssignablePtrType3)(<nil>)"),
		})

	// (A) Equal(A) method not found
	checkError(t, BadEqual1(1), BadEqual1(2),
		expectedError{
			Message: mustBe("values differ"),
		})
	checkError(t, BadEqual2(1), BadEqual2(2),
		expectedError{
			Message: mustBe("values differ"),
		})
	checkError(t, BadEqual3(1), BadEqual3(2),
		expectedError{
			Message: mustBe("values differ"),
		})
	checkError(t, BadEqual4(1), BadEqual4(2),
		expectedError{
			Message: mustBe("values differ"),
		})
	checkError(t, BadEqual5(1), BadEqual5(2),
		expectedError{
			Message: mustBe("values differ"),
		})
}

func TestUseEqualGlobalVsAnchor(t *testing.T) {
	defer func() { td.DefaultContextConfig.UseEqual = false }()
	td.DefaultContextConfig.UseEqual = true

	tt := test.NewTestingTB(t.Name())

	assert := td.Assert(tt)

	type timeAnchored struct {
		Time time.Time
	}
	td.CmpTrue(t,
		assert.Cmp(
			timeAnchored{Time: timeParse(t, "2022-05-31T06:00:00Z")},
			timeAnchored{
				Time: assert.A(td.Between(
					timeParse(t, "2022-05-31T00:00:00Z"),
					timeParse(t, "2022-05-31T12:00:00Z"),
				)).(time.Time),
			}))
}

func TestBeLaxGlobalt(t *testing.T) {
	defer func() { td.DefaultContextConfig.BeLax = false }()
	td.DefaultContextConfig.BeLax = true

	// expected float64 value first converted to int64 before comparison
	checkOK(t, int64(123), float64(123.56))

	type MyInt int32
	checkOK(t, int64(123), MyInt(123))
	checkOK(t, MyInt(123), int64(123))

	type gotStruct struct {
		name string
		age  int
	}
	type expectedStruct struct {
		name string
		age  int
	}
	checkOK(t,
		gotStruct{
			name: "bob",
			age:  42,
		},
		expectedStruct{
			name: "bob",
			age:  42,
		})
	checkOK(t,
		&gotStruct{
			name: "bob",
			age:  42,
		},
		&expectedStruct{
			name: "bob",
			age:  42,
		})
}
