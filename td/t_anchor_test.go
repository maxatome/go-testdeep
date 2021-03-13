// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestAnchor(tt *testing.T) {
	timeParse := func(s string) time.Time {
		dt, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			tt.Helper()
			tt.Fatalf("Cannot parse `%s`: %s", s, err)
		}
		return dt
	}

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
		Time: timeParse("2019-01-02T11:22:33.123456Z"),
	}

	// Using T.Anchor()
	td.CmpTrue(tt,
		t.Cmp(got, MyStruct{
			PNum: t.Anchor(td.Ptr(td.Between(40, 45))).(*int),
			Num:  t.Anchor(td.Between(int64(135), int64(137))).(int64),
			Str:  t.Anchor(td.HasPrefix("Pipo"), "").(string),
			Time: t.Anchor(td.TruncTime(timeParse("2019-01-02T11:22:00Z"), time.Minute)).(time.Time),
		}))

	// Using T.A()
	td.CmpTrue(tt,
		t.Cmp(got, MyStruct{
			PNum: t.A(td.Ptr(td.Between(40, 45))).(*int),
			Num:  t.A(td.Between(int64(135), int64(137))).(int64),
			Str:  t.A(td.HasPrefix("Pipo"), "").(string),
			Time: t.A(td.TruncTime(timeParse("2019-01-02T11:22:00Z"), time.Minute)).(time.Time),
		}))

	// Testing persistence
	got = MyStruct{Num: 136}

	tt.Run("without persistence", func(tt *testing.T) {
		numOp := t.Anchor(td.Between(int64(135), int64(137))).(int64)

		td.CmpTrue(tt, t.Cmp(got, MyStruct{Num: numOp}))
		td.CmpFalse(tt, t.Cmp(got, MyStruct{Num: numOp}))
	})

	tt.Run("with persistence", func(tt *testing.T) {
		numOp := t.Anchor(td.Between(int64(135), int64(137))).(int64)
		defer t.AnchorsPersistTemporarily()()

		td.CmpTrue(tt, t.Cmp(got, MyStruct{Num: numOp}))
		td.CmpTrue(tt, t.Cmp(got, MyStruct{Num: numOp}))

		t.ResetAnchors() // force reset anchored operators
		td.CmpFalse(tt, t.Cmp(got, MyStruct{Num: numOp}))
	})

	// Errors
	tt.Run("errors", func(tt *testing.T) {
		td.Cmp(tt, ttt.CatchFatal(func() { t.Anchor(nil) }),
			"Cannot anchor a nil TestDeep operator")

		td.Cmp(tt, ttt.CatchFatal(func() { t.Anchor(td.Ignore(), 1, 2) }),
			"usage: Anchor(OPERATOR[, MODEL]), too many parameters")

		td.Cmp(tt, ttt.CatchFatal(func() { t.Anchor(td.Ignore(), nil) }),
			"Untyped nil value is not valid as model for an anchor")

		td.Cmp(tt, ttt.CatchFatal(func() { t.Anchor(td.Between(1, 2), 12.3) }),
			"Operator Between TypeBehind() returned int which differs from model type float64. Omit model or ensure its type is int")

		td.Cmp(tt, ttt.CatchFatal(func() { t.Anchor(td.Ignore()) }),
			"Cannot anchor operator Ignore as TypeBehind() returned nil. Use model parameter to specify the type to return")
	})
}

type privStruct struct {
	num int64
}

func (p privStruct) Num() int64 {
	return p.num
}

func TestAddAnchorableStructType(tt *testing.T) {
	type MyStruct struct {
		Priv privStruct
	}

	ttt := test.NewTestingTB(tt.Name())
	t := td.NewT(ttt)

	// We want to anchor this operator
	op := td.Smuggle((privStruct).Num, int64(42))

	// Without making privStruct anchorable, it does not work
	td.Cmp(tt, ttt.CatchFatal(func() { t.A(op, privStruct{}) }),
		"td_test.privStruct struct type is not supported as an anchor. Try AddAnchorableStructType")

	// Make privStruct anchorable
	td.AddAnchorableStructType(func(nextAnchor int) privStruct {
		return privStruct{num: int64(2e9 - nextAnchor)}
	})

	td.CmpTrue(tt,
		t.Cmp(MyStruct{Priv: privStruct{num: 42}},
			MyStruct{
				Priv: t.A(op, privStruct{}).(privStruct), // ← now it works
			}))

	// Error
	dark.CheckFatalizerBarrierErr(tt, func() { td.AddAnchorableStructType(123) },
		"usage: AddAnchorableStructType(func (nextAnchor int) STRUCT_TYPE)")
}

func TestAnchorsPersist(tt *testing.T) {
	ttt := test.NewTestingTB(tt.Name())

	t1 := td.NewT(ttt)
	t2 := td.NewT(ttt)
	t3 := td.NewT(t1)

	tt.Run("without anchors persistence", func(tt *testing.T) {
		// Anchors persistence is shared for a same testing.TB
		td.CmpFalse(tt, t1.DoAnchorsPersist())
		td.CmpFalse(tt, t2.DoAnchorsPersist())
		td.CmpFalse(tt, t3.DoAnchorsPersist())

		func() {
			defer t1.AnchorsPersistTemporarily()()
			td.CmpTrue(tt, t1.DoAnchorsPersist())
			td.CmpTrue(tt, t2.DoAnchorsPersist())
			td.CmpTrue(tt, t3.DoAnchorsPersist())
		}()
		td.CmpFalse(tt, t1.DoAnchorsPersist())
		td.CmpFalse(tt, t2.DoAnchorsPersist())
		td.CmpFalse(tt, t3.DoAnchorsPersist())
	})

	tt.Run("with anchors persistence", func(tt *testing.T) {
		t3.SetAnchorsPersist(true)

		td.CmpTrue(tt, t1.DoAnchorsPersist())
		td.CmpTrue(tt, t2.DoAnchorsPersist())
		td.CmpTrue(tt, t3.DoAnchorsPersist())

		func() {
			defer t1.AnchorsPersistTemporarily()()
			td.CmpTrue(tt, t1.DoAnchorsPersist())
			td.CmpTrue(tt, t2.DoAnchorsPersist())
			td.CmpTrue(tt, t3.DoAnchorsPersist())
		}()
		td.CmpTrue(tt, t1.DoAnchorsPersist())
		td.CmpTrue(tt, t2.DoAnchorsPersist())
		td.CmpTrue(tt, t3.DoAnchorsPersist())
	})
}
