// Copyright (c) 2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.18
// +build go1.18

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestAnchor(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())
	t := td.NewT(ttt)
	type MyStruct struct {
		PNum  *int
		Num   int64
		Str   string
		Slice []int
		Map   map[string]bool
		Time  time.Time
	}
	n := 42
	got := MyStruct{
		PNum: &n,
		Num:  136,
		Str:  "Pipo bingo",
		Time: timeParse(tt, "2019-01-02T11:22:33.123456Z"),
	}

	td.CmpTrue(tt,
		t.Cmp(got, MyStruct{
			PNum: td.Anchor[*int](t, td.Ptr(td.Between(40, 45))),
			Num:  td.Anchor[int64](t, td.Between(int64(135), int64(137))),
			Str:  td.Anchor[string](t, td.HasPrefix("Pipo")),
			Time: td.Anchor[time.Time](t, td.TruncTime(timeParse(tt, "2019-01-02T11:22:00Z"), time.Minute)),
		}))

	td.CmpTrue(tt,
		t.Cmp(got, MyStruct{
			PNum: td.A[*int](t, td.Ptr(td.Between(40, 45))),
			Num:  td.A[int64](t, td.Between(int64(135), int64(137))),
			Str:  td.A[string](t, td.HasPrefix("Pipo")),
			Time: td.A[time.Time](t, td.TruncTime(timeParse(tt, "2019-01-02T11:22:00Z"), time.Minute)),
		}))
}
