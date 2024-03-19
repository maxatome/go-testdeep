// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestWithCmpHooks(tt *testing.T) {
	na, nb := 1234, 1234
	date, _ := time.Parse(time.RFC3339, "2020-09-08T22:13:54+02:00")

	for _, tst := range []struct {
		name          string
		cmp           any
		got, expected any
	}{
		{
			name: "reflect.Value",
			cmp: func(got, expected reflect.Value) bool {
				return td.EqDeeply(got.Interface(), expected.Interface())
			},
			got:      reflect.ValueOf(&na),
			expected: reflect.ValueOf(&nb),
		},
		{
			name:     "time.Time",
			cmp:      (time.Time).Equal,
			got:      date,
			expected: date.UTC(),
		},
		{
			name: "numify",
			cmp: func(got, expected string) error {
				ngot, err := strconv.Atoi(got)
				if err != nil {
					return fmt.Errorf("strconv.Atoi(got) failed: %s", err)
				}
				nexpected, err := strconv.Atoi(expected)
				if err != nil {
					return fmt.Errorf("strconv.Atoi(expected) failed: %s", err)
				}
				if ngot != nexpected {
					return errors.New("values differ")
				}
				return nil
			},
			got:      "0000001234",
			expected: "1234",
		},
		{
			name: "false test :)",
			cmp: func(got, expected int) bool {
				return got == -expected
			},
			got:      1,
			expected: -1,
		},
	} {
		tt.Run(tst.name, func(tt *testing.T) {
			ttt := test.NewTestingTB(tt.Name())

			t := td.NewT(ttt)

			td.CmpFalse(tt, func() bool {
				// A panic can occur when -tags safe:
				//   dark.GetInterface() does not handle private unsafe.Pointer kind
				defer func() { recover() }() //nolint: errcheck
				return t.Cmp(tst.got, tst.expected)
			}())

			t = t.WithCmpHooks(tst.cmp)

			td.CmpTrue(tt, t.Cmp(tst.got, tst.expected))
		})
	}

	tt.Run("Error", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt).
			WithCmpHooks(func(got, expected int) error {
				return errors.New("never equal")
			})

		td.CmpFalse(tt, t.Cmp(1, 1))

		if !strings.Contains(ttt.LastMessage(), "DATA: never equal\n") {
			tt.Errorf(`<%s> does not contain "DATA: never equal\n"`, ttt.LastMessage())
		}
	})

	for _, tst := range []struct {
		name  string
		cmp   any
		fatal string
	}{
		{
			name:  "not a function",
			cmp:   "Booh",
			fatal: "WithCmpHooks expects a function, not a string",
		},
		{
			name:  "wrong signature",
			cmp:   func(a []int, b ...int) bool { return false },
			fatal: "WithCmpHooks expects: func (T, T) bool|error not ",
		},
	} {
		tt.Run("panic: "+tst.name, func(tt *testing.T) {
			ttt := test.NewTestingTB(tt.Name())

			t := td.NewT(ttt)

			fatalMesg := ttt.CatchFatal(func() { t.WithCmpHooks(tst.cmp) })
			test.IsTrue(tt, ttt.IsFatal)
			if !strings.Contains(fatalMesg, tst.fatal) {
				tt.Errorf(`<%s> does not contain %q`, fatalMesg, tst.fatal)
			}
		})
	}
}

func TestWithSmuggleHooks(tt *testing.T) {
	for _, tst := range []struct {
		name          string
		cmp           any
		got, expected any
	}{
		{
			name: "abs",
			cmp: func(got int) int {
				if got < 0 {
					return -got
				}
				return got
			},
			got:      -1234,
			expected: 1234,
		},
		{
			name:     "int2bool",
			cmp:      func(got int) bool { return got != 0 },
			got:      1,
			expected: true,
		},
		{
			name:     "Atoi",
			cmp:      strconv.Atoi,
			got:      "1234",
			expected: 1234,
		},
		{
			name: "reflect2any",
			cmp: func(v reflect.Value) any {
				return v.Interface()
			},
			got:      reflect.ValueOf(123),
			expected: 123,
		},
	} {
		tt.Run(tst.name, func(tt *testing.T) {
			ttt := test.NewTestingTB(tt.Name())

			t := td.NewT(ttt)

			td.CmpFalse(tt, t.Cmp(tst.got, tst.expected))

			t = t.WithSmuggleHooks(tst.cmp)

			td.CmpTrue(tt, t.Cmp(tst.got, tst.expected))
		})
	}

	tt.Run("Array", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt).WithSmuggleHooks(func(got reflect.Value) any {
			if got.IsValid() {
				return got.Interface()
			}
			return nil
		})

		got := [3]reflect.Value{1: reflect.ValueOf(42)}

		td.CmpTrue(tt, td.Cmp(t, got, td.Array([3]reflect.Value{}, td.ArrayEntries{
			0: nil, // if omitted, td.Array sets it to reflect.Value{}
			1: 42,
			2: nil, // if omitted, td.Array sets it to reflect.Value{}
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got[:], td.Slice([]reflect.Value{}, td.ArrayEntries{
			0: nil, // if omitted, td.Array sets it to reflect.Value{}
			1: 42,
			2: nil, // if omitted, td.Array sets it to reflect.Value{}
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got[:], td.SuperSliceOf([]reflect.Value{}, td.ArrayEntries{
			1: 42,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()
	})

	tt.Run("Map", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt).WithSmuggleHooks(func(got reflect.Value) any {
			if got.IsValid() {
				return got.Interface()
			}
			return nil
		})

		got := map[string]reflect.Value{
			"pipo":  reflect.ValueOf(42),
			"bingo": reflect.ValueOf(666),
			"zip":   {},
		}

		td.CmpTrue(tt, td.Cmp(t, got, td.Map(map[string]reflect.Value{}, td.MapEntries{
			"pipo":  42,
			"bingo": 666,
			"zip":   nil,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got, td.SubMapOf(map[string]reflect.Value{}, td.MapEntries{
			"pipo":  42,
			"bingo": 666,
			"zip":   nil,
			"test":  "more",
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got, td.SuperMapOf(map[string]reflect.Value{}, td.MapEntries{
			"pipo": 42,
			"zip":  nil,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()
	})

	tt.Run("Struct", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt).WithSmuggleHooks(func(got reflect.Value) any {
			if got.IsValid() {
				return got.Interface()
			}
			return nil
		})

		type Rstruct struct {
			V1 reflect.Value
			V2 reflect.Value
		}

		got := Rstruct{V1: reflect.ValueOf(42)}

		td.CmpTrue(tt, td.Cmp(t, got, td.Struct(Rstruct{}, td.StructFields{
			"V1": 42,
			"V2": nil,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got, td.SStruct(Rstruct{}, td.StructFields{
			"V1": 42,
			"V2": nil,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()

		td.CmpTrue(tt, td.Cmp(t, got, td.Struct(nil, td.StructFields{
			"V1": 42,
			"V2": nil,
		})))
		test.EqualStr(tt, ttt.LastMessage(), "")
		ttt.ResetMessages()
	})

	tt.Run("Error", func(tt *testing.T) {
		ttt := test.NewTestingTB(tt.Name())

		t := td.NewT(ttt).WithSmuggleHooks(func(got int) (int, error) {
			return 0, errors.New("never equal")
		})

		td.CmpFalse(tt, t.Cmp(1, 1))

		if !strings.Contains(ttt.LastMessage(), "DATA: never equal\n") {
			tt.Errorf(`<%s> does not contain "DATA: never equal\n"`, ttt.LastMessage())
		}
	})

	for _, tst := range []struct {
		name  string
		cmp   any
		fatal string
	}{
		{
			name:  "not a function",
			cmp:   "Booh",
			fatal: "WithSmuggleHooks expects a function, not a string",
		},
		{
			name:  "wrong signature",
			cmp:   func(a []int, b ...int) bool { return false },
			fatal: "WithSmuggleHooks expects: func (A) (B[, error]) not ",
		},
	} {
		tt.Run("panic: "+tst.name, func(tt *testing.T) {
			ttt := test.NewTestingTB(tt.Name())

			t := td.NewT(ttt)

			fatalMesg := ttt.CatchFatal(func() { t.WithSmuggleHooks(tst.cmp) })
			test.IsTrue(tt, ttt.IsFatal)
			if !strings.Contains(fatalMesg, tst.fatal) {
				tt.Errorf(`<%s> does not contain %q`, fatalMesg, tst.fatal)
			}
		})
	}
}
