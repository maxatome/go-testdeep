// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil

import (
	"math"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/internal/visited"
)

func TestSortCmp(t *testing.T) {
	checkCmp := func(a, b any, expected int) {
		t.Helper()
		got := cmp(visited.NewVisited(), reflect.ValueOf(a), reflect.ValueOf(b))
		if got != expected {
			t.Errorf("cmp() failed: got=%d expected=%d\n", got, expected)
		}
	}

	// IsValid
	checkCmp(nil, 12, -1)
	checkCmp(nil, nil, 0)
	checkCmp(12, nil, 1)

	// type mismatch: int is before string
	checkCmp(42, "str", -1)
	checkCmp("str", 42, 1)

	// bool
	checkCmp(true, true, 0)
	checkCmp(true, false, 1)
	checkCmp(false, true, -1)
	checkCmp(false, false, 0)

	// int
	checkCmp(12, 42, -1)
	checkCmp(42, 12, 1)
	checkCmp(12, 12, 0)

	checkCmp(int8(12), int8(42), -1)
	checkCmp(int8(42), int8(12), 1)
	checkCmp(int8(12), int8(12), 0)

	checkCmp(int16(12), int16(42), -1)
	checkCmp(int16(42), int16(12), 1)
	checkCmp(int16(12), int16(12), 0)

	checkCmp(int32(12), int32(42), -1)
	checkCmp(int32(42), int32(12), 1)
	checkCmp(int32(12), int32(12), 0)

	checkCmp(int64(12), int64(42), -1)
	checkCmp(int64(42), int64(12), 1)
	checkCmp(int64(12), int64(12), 0)

	// uint
	checkCmp(uint(12), uint(42), -1)
	checkCmp(uint(42), uint(12), 1)
	checkCmp(uint(12), uint(12), 0)

	checkCmp(uint8(12), uint8(42), -1)
	checkCmp(uint8(42), uint8(12), 1)
	checkCmp(uint8(12), uint8(12), 0)

	checkCmp(uint16(12), uint16(42), -1)
	checkCmp(uint16(42), uint16(12), 1)
	checkCmp(uint16(12), uint16(12), 0)

	checkCmp(uint32(12), uint32(42), -1)
	checkCmp(uint32(42), uint32(12), 1)
	checkCmp(uint32(12), uint32(12), 0)

	checkCmp(uint64(12), uint64(42), -1)
	checkCmp(uint64(42), uint64(12), 1)
	checkCmp(uint64(12), uint64(12), 0)

	checkCmp(uintptr(12), uintptr(42), -1)
	checkCmp(uintptr(42), uintptr(12), 1)
	checkCmp(uintptr(12), uintptr(12), 0)

	// float
	checkCmp(float32(12), float32(42), -1)
	checkCmp(float32(42), float32(12), 1)
	checkCmp(float32(12), float32(12), 0)

	checkCmp(float64(12), float64(42), -1)
	checkCmp(float64(42), float64(12), 1)
	checkCmp(float64(12), float64(12), 0)
	checkCmp(float64(12), float64(12), 0)

	checkCmp(math.NaN(), float64(12), -1)
	checkCmp(math.NaN(), math.NaN(), -1)
	checkCmp(float64(12), math.NaN(), 1)

	// complex
	checkCmp(complex(12, 0), complex(42, 0), -1)
	checkCmp(complex(42, 0), complex(12, 0), 1)
	checkCmp(complex(0, 12), complex(0, 42), -1)
	checkCmp(complex(0, 42), complex(0, 12), 1)
	checkCmp(complex(12, 0), complex(12, 0), 0)

	checkCmp(complex(float32(12), 0), complex(float32(42), 0), -1)
	checkCmp(complex(float32(42), 0), complex(float32(12), 0), 1)
	checkCmp(complex(float32(0), 12), complex(float32(0), 42), -1)
	checkCmp(complex(float32(0), 42), complex(float32(0), 12), 1)
	checkCmp(complex(float32(12), 0), complex(float32(12), 0), 0)

	// string
	checkCmp("aaa", "bbb", -1)
	checkCmp("bbb", "aaa", 1)
	checkCmp("aaa", "aaa", 0)

	// array
	checkCmp([3]byte{1, 2, 3}, [3]byte{3, 1, 2}, -1)
	checkCmp([3]byte{3, 1, 2}, [3]byte{1, 2, 3}, 1)
	checkCmp([3]byte{1, 2, 3}, [3]byte{1, 2, 3}, 0)

	// slice
	checkCmp([]byte{1, 2, 3}, []byte{3, 1, 2}, -1)
	checkCmp([]byte{3, 1, 2}, []byte{1, 2, 3}, 1)
	checkCmp([]byte{1, 2, 3}, []byte{1, 2, 3}, 0)
	checkCmp([]byte{1, 2, 3}, []byte{1, 2, 3, 4}, -1)
	checkCmp([]byte{1, 2, 3, 4}, []byte{1, 2, 3}, 1)

	// interface
	checkCmp([]any{1}, []any{3}, -1)
	checkCmp([]any{3}, []any{1}, 1)
	checkCmp([]any{1}, []any{1}, 0)
	checkCmp([]any{nil}, []any{nil}, 0)
	checkCmp([]any{nil}, []any{1}, -1)
	checkCmp([]any{1}, []any{nil}, 1)

	// struct
	type myStruct struct {
		n int
		s string
		p *myStruct
	}
	checkCmp(myStruct{n: 12, s: "a"}, myStruct{n: 12, s: "b"}, -1)
	checkCmp(myStruct{n: 12, s: "b"}, myStruct{n: 12, s: "a"}, 1)
	checkCmp(myStruct{n: 12, s: "a"}, myStruct{n: 12, s: "a"}, 0)

	// ptr
	a, b := 12, 42
	checkCmp(&a, &a, 0)
	checkCmp((*int)(nil), (*int)(nil), 0)
	checkCmp(&a, &b, -1)
	checkCmp((*int)(nil), &b, -1)
	checkCmp(&b, &a, 1)
	checkCmp(&b, (*int)(nil), 1)

	// map
	ma, mb := map[int]bool{12: true}, map[int]bool{12: true, 13: false}
	checkCmp(ma, mb, -1)
	checkCmp(mb, ma, 1)
	checkCmp(ma, ma, 0)
	checkCmp((map[int]bool)(nil), (map[int]bool)(nil), 0)
	checkCmp((map[int]bool)(nil), ma, -1)
	checkCmp(ma, (map[int]bool)(nil), 1)

	// cyclic references protection
	pa := &myStruct{n: 42, p: &myStruct{n: 18}}
	pa.p.p = pa.p
	pb := &myStruct{n: 42, p: &myStruct{n: 18}}
	pb.p.p = pb.p
	checkCmp(pa, pb, 0)
}
