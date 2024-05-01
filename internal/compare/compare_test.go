// Copyright (c) 2019-2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package compare_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/compare"
	"github.com/maxatome/go-testdeep/internal/visited"
)

type cmp1 int

func (cmp1) Compare(cmp1) bool { return false }

type cmp2 int

func (cmp2) Compare(int) int { return 0 }

type cmp3 int

func (cmp3) Compare(...int) int { return 0 }

type cmpPanic int

func (cmpPanic) Compare(cmpPanic) { panic("boom!") }

type cmpOK int

func (a cmpOK) Compare(b cmpOK) int {
	if a == b {
		return 0
	}
	if a > b {
		return -1
	}
	return 1
}

type cmpOKp int

func newCmpOKp(x cmpOKp) *cmpOKp {
	return &x
}

func (a *cmpOKp) Compare(b *cmpOKp) int {
	if *a == *b {
		return 0
	}
	if *a > *b {
		return -1
	}
	return 1
}

func TestCompare(t *testing.T) {
	type myStruct struct {
		n int
		s string
		p *myStruct
	}
	a, b := 12, 42
	ma, mb := map[int]bool{12: true}, map[int]bool{12: true, 13: false}

	checkCompare := func(t *testing.T, a, b any, expected int) {
		t.Helper()
		got := compare.Compare(visited.NewVisited(), reflect.ValueOf(a), reflect.ValueOf(b))
		if got != expected {
			t.Errorf("Compare() failed: got=%d expected=%d\n", got, expected)
		}
	}

	testCases := []struct {
		a, b     any
		expected int
	}{
		// nil
		{nil, 12, -1},
		{nil, nil, 0},
		{12, nil, 1},
		// type mismatch: int is before string
		{42, "str", -1},
		{"str", 42, 1},
		// Compare
		{cmpOK(99), cmpOK(11), -1},
		{cmpOK(11), cmpOK(99), 1},
		{cmpOK(99), cmpOK(99), 0},
		{newCmpOKp(99), newCmpOKp(11), -1},
		{newCmpOKp(11), newCmpOKp(99), 1},
		{newCmpOKp(99), newCmpOKp(99), 0},
		// Compare does not match expected signature
		{cmp1(99), cmp1(11), 1},
		{cmp1(11), cmp1(99), -1},
		{cmp1(99), cmp1(99), 0},
		{cmp2(99), cmp2(11), 1},
		{cmp2(11), cmp2(99), -1},
		{cmp2(99), cmp2(99), 0},
		{cmp3(99), cmp3(11), 1},
		{cmp3(11), cmp3(99), -1},
		{cmp3(99), cmp3(99), 0},
		// Compare panics fallback
		{cmpPanic(99), cmpPanic(11), 1},
		{cmpPanic(11), cmpPanic(99), -1},
		{cmpPanic(99), cmpPanic(99), 0},
		// bool
		{true, true, 0},
		{true, false, 1},
		{false, true, -1},
		{false, false, 0},
		// int
		{12, 42, -1},
		{42, 12, 1},
		{12, 12, 0},
		{int8(12), int8(42), -1},
		{int8(42), int8(12), 1},
		{int8(12), int8(12), 0},
		{int16(12), int16(42), -1},
		{int16(42), int16(12), 1},
		{int16(12), int16(12), 0},
		{int32(12), int32(42), -1},
		{int32(42), int32(12), 1},
		{int32(12), int32(12), 0},
		{int64(12), int64(42), -1},
		{int64(42), int64(12), 1},
		{int64(12), int64(12), 0},
		// uint
		{uint(12), uint(42), -1},
		{uint(42), uint(12), 1},
		{uint(12), uint(12), 0},
		{uint8(12), uint8(42), -1},
		{uint8(42), uint8(12), 1},
		{uint8(12), uint8(12), 0},
		{uint16(12), uint16(42), -1},
		{uint16(42), uint16(12), 1},
		{uint16(12), uint16(12), 0},
		{uint32(12), uint32(42), -1},
		{uint32(42), uint32(12), 1},
		{uint32(12), uint32(12), 0},
		{uint64(12), uint64(42), -1},
		{uint64(42), uint64(12), 1},
		{uint64(12), uint64(12), 0},
		{uintptr(12), uintptr(42), -1},
		{uintptr(42), uintptr(12), 1},
		{uintptr(12), uintptr(12), 0},
		// float
		{float32(12), float32(42), -1},
		{float32(42), float32(12), 1},
		{float32(12), float32(12), 0},
		{float64(12), float64(42), -1},
		{float64(42), float64(12), 1},
		{float64(12), float64(12), 0},
		{float64(12), float64(12), 0},
		{math.NaN(), float64(12), -1},
		{math.NaN(), math.NaN(), -1},
		{float64(12), math.NaN(), 1},
		// complex
		{complex(12, 0), complex(42, 0), -1},
		{complex(42, 0), complex(12, 0), 1},
		{complex(0, 12), complex(0, 42), -1},
		{complex(0, 42), complex(0, 12), 1},
		{complex(12, 0), complex(12, 0), 0},
		{complex(float32(12), 0), complex(float32(42), 0), -1},
		{complex(float32(42), 0), complex(float32(12), 0), 1},
		{complex(float32(0), 12), complex(float32(0), 42), -1},
		{complex(float32(0), 42), complex(float32(0), 12), 1},
		{complex(float32(12), 0), complex(float32(12), 0), 0},
		// string
		{"aaa", "bbb", -1},
		{"bbb", "aaa", 1},
		{"aaa", "aaa", 0},
		// array
		{[3]byte{1, 2, 3}, [3]byte{3, 1, 2}, -1},
		{[3]byte{3, 1, 2}, [3]byte{1, 2, 3}, 1},
		{[3]byte{1, 2, 3}, [3]byte{1, 2, 3}, 0},
		// slice
		{[]byte{1, 2, 3}, []byte{3, 1, 2}, -1},
		{[]byte{3, 1, 2}, []byte{1, 2, 3}, 1},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, 0},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 4}, -1},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, 1},
		// interface
		{[]any{1}, []any{3}, -1},
		{[]any{3}, []any{1}, 1},
		{[]any{1}, []any{1}, 0},
		{[]any{nil}, []any{nil}, 0},
		{[]any{nil}, []any{1}, -1},
		{[]any{1}, []any{nil}, 1},
		// struct
		{myStruct{n: 12, s: "a"}, myStruct{n: 12, s: "b"}, -1},
		{myStruct{n: 12, s: "b"}, myStruct{n: 12, s: "a"}, 1},
		{myStruct{n: 12, s: "a"}, myStruct{n: 12, s: "a"}, 0},
		// ptr
		{&a, &a, 0},
		{(*int)(nil), (*int)(nil), 0},
		{&a, &b, -1},
		{(*int)(nil), &b, -1},
		{&b, &a, 1},
		{&b, (*int)(nil), 1},
		// map
		{ma, mb, -1},
		{mb, ma, 1},
		{ma, ma, 0},
		{(map[int]bool)(nil), (map[int]bool)(nil), 0},
		{(map[int]bool)(nil), ma, -1},
		{ma, (map[int]bool)(nil), 1},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d %[1]T(%[1]v) %[1]T(%[1]v)", i, tc.a, tc.b),
			func(t *testing.T) {
				checkCompare(t, tc.a, tc.b, tc.expected)
			})
	}

	// cyclic references protection
	pa := &myStruct{n: 42, p: &myStruct{n: 18}}
	pa.p.p = pa.p
	pb := &myStruct{n: 42, p: &myStruct{n: 18}}
	pb.p.p = pb.p
	checkCompare(t, pa, pb, 0)
}
