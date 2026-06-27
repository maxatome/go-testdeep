// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestT_AnchorT(tt *testing.T) {
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
			PNum: t.AnchorT[*int](td.Ptr(td.Between(40, 45))),
			Num:  t.AnchorT[int64](td.Between(int64(135), int64(137))),
			Str:  t.AnchorT[string](td.HasPrefix("Pipo")),
			Time: t.AnchorT[time.Time](td.TruncTime(timeParse(tt, "2019-01-02T11:22:00Z"), time.Minute)),
		}))

	td.CmpTrue(tt,
		t.Cmp(got, MyStruct{
			PNum: t.AT[*int](td.Ptr(td.Between(40, 45))),
			Num:  t.AT[int64](td.Between(int64(135), int64(137))),
			Str:  t.AT[string](td.HasPrefix("Pipo")),
			Time: t.AT[time.Time](td.TruncTime(timeParse(tt, "2019-01-02T11:22:00Z"), time.Minute)),
		}))
}
