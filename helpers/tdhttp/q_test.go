// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp_test

import (
	"net/url"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

type qTest1 struct{}

func (qTest1) String() string { return "qTest1!" }

type qTest2 struct{}

func (*qTest2) String() string { return "qTest2!" }

func TestQ(t *testing.T) {
	q := tdhttp.Q{
		"str1":   "v1",
		"str2":   []string{"v20", "v21"},
		"int1":   1234,
		"int2":   []int{1, 2, 3},
		"uint1":  uint(1234),
		"uint2":  [3]uint{1, 2, 3},
		"float1": 1.2,
		"float2": []float64{1.2, 3.4},
		"bool1":  true,
		"bool2":  [2]bool{true, false},
	}
	td.Cmp(t, q.Values(), url.Values{
		"str1":   []string{"v1"},
		"str2":   []string{"v20", "v21"},
		"int1":   []string{"1234"},
		"int2":   []string{"1", "2", "3"},
		"uint1":  []string{"1234"},
		"uint2":  []string{"1", "2", "3"},
		"float1": []string{"1.2"},
		"float2": []string{"1.2", "3.4"},
		"bool1":  []string{"true"},
		"bool2":  []string{"true", "false"},
	})

	// Auto deref pointers
	num := 123
	pnum := &num
	ppnum := &pnum
	q = tdhttp.Q{
		"pnum":   pnum,
		"ppnum":  ppnum,
		"pppnum": &ppnum,
		"slice":  []***int{&ppnum, &ppnum},
		"pslice": &[]***int{&ppnum, &ppnum},
		"array":  [2]***int{&ppnum, &ppnum},
		"parray": &[2]***int{&ppnum, &ppnum},
	}
	td.Cmp(t, q.Values(), url.Values{
		"pnum":   []string{"123"},
		"ppnum":  []string{"123"},
		"pppnum": []string{"123"},
		"slice":  []string{"123", "123"},
		"pslice": []string{"123", "123"},
		"array":  []string{"123", "123"},
		"parray": []string{"123", "123"},
	})

	// Auto deref interfaces
	q = tdhttp.Q{
		"all": []any{
			"string",
			-1,
			int8(-2),
			int16(-3),
			int32(-4),
			int64(-5),
			uint(1),
			uint8(2),
			uint16(3),
			uint32(4),
			uint64(5),
			float32(6),
			float64(7),
			true,
			ppnum,
			(*int)(nil), // ignored
			nil,         // ignored
			qTest1{},
			&qTest1{}, // does not implement fmt.Stringer, but qTest does
			// qTest2{} panics as it does not implement fmt.Stringer, see Errors below
			&qTest2{},
		},
	}
	td.Cmp(t, q.Values(), url.Values{
		"all": []string{
			"string",
			"-1",
			"-2",
			"-3",
			"-4",
			"-5",
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"true",
			"123",
			"qTest1!",
			"qTest1!",
			"qTest2!",
		},
	})

	// nil case
	pnum = nil
	q = tdhttp.Q{
		"nil1": &ppnum,
		"nil2": (*int)(nil),
		"nil3": nil,
		"nil4": []*int{nil, nil},
		"nil5": ([]int)(nil),
		"nil6": []any{nil, nil},
	}
	td.Cmp(t, q.Values(), url.Values{})

	q = tdhttp.Q{
		"id":    []int{12, 34},
		"draft": true,
	}
	td.Cmp(t, q.Encode(), "draft=true&id=12&id=34")

	// Errors
	td.CmpPanic(t, func() { (tdhttp.Q{"panic": map[string]bool{}}).Values() },
		td.Contains(`don't know how to add type map[string]bool (map) to param "panic"`))
	td.CmpPanic(t, func() { (tdhttp.Q{"panic": qTest2{}}).Values() },
		td.Contains(`don't know how to add type tdhttp_test.qTest2 (struct) to param "panic"`))

	td.CmpPanic(t,
		func() { (tdhttp.Q{"panic": []any{[]int{}}}).Values() },
		td.Contains(`slice is only allowed at the root level for param "panic"`))
}
