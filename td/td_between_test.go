// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/td"
)

func TestBetween(t *testing.T) {
	checkOK(t, 12, td.Between(9, 13))
	checkOK(t, 12, td.Between(13, 9))
	checkOK(t, 12, td.Between(9, 12, td.BoundsOutIn))
	checkOK(t, 12, td.Between(12, 13, td.BoundsInOut))

	checkError(t, 10, td.Between(10, 15, td.BoundsOutIn),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("10 < got ≤ 15"),
		})
	checkError(t, 10, td.Between(10, 15, td.BoundsOutOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("10 < got < 15"),
		})
	checkError(t, 15, td.Between(10, 15, td.BoundsInOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("15"),
			Expected: mustBe("10 ≤ got < 15"),
		})
	checkError(t, 15, td.Between(10, 15, td.BoundsOutOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("15"),
			Expected: mustBe("10 < got < 15"),
		})

	checkError(t, 15, td.Between(uint(10), uint(15), td.BoundsOutOut),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint"),
		})

	checkOK(t, uint16(12), td.Between(uint16(9), uint16(13)))
	checkOK(t, uint16(12), td.Between(uint16(13), uint16(9)))
	checkOK(t, uint16(12),
		td.Between(uint16(9), uint16(12), td.BoundsOutIn))
	checkOK(t, uint16(12),
		td.Between(uint16(12), uint16(13), td.BoundsInOut))

	checkOK(t, 12.1, td.Between(9.5, 13.1))
	checkOK(t, 12.1, td.Between(13.1, 9.5))
	checkOK(t, 12.1, td.Between(9.5, 12.1, td.BoundsOutIn))
	checkOK(t, 12.1, td.Between(12.1, 13.1, td.BoundsInOut))

	checkOK(t, "abc", td.Between("aaa", "bbb"))
	checkOK(t, "abc", td.Between("bbb", "aaa"))
	checkOK(t, "abc", td.Between("aaa", "abc", td.BoundsOutIn))
	checkOK(t, "abc", td.Between("abc", "bbb", td.BoundsInOut))

	checkOK(t, 12*time.Hour, td.Between(60*time.Second, 24*time.Hour))

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t,
		func() { td.Between([]byte("test"), []byte("test")) },
		"usage: Between(")
	dark.CheckFatalizerBarrierErr(t, func() { td.Between(12, "test") },
		"Between(FROM, TO): FROM and TO must have the same type: int ≠ string")
	dark.CheckFatalizerBarrierErr(t, func() { td.Between("test", 12) },
		"Between(FROM, TO): FROM and TO must have the same type: string ≠ int")
	dark.CheckFatalizerBarrierErr(t,
		func() { td.Between(1, 2, td.BoundsInIn, td.BoundsInOut) },
		"usage: Between(")
	dark.CheckFatalizerBarrierErr(t,
		func() {
			type notTime struct{}
			td.Between(notTime{}, notTime{})
		},
		"usage: Between(")
}

func TestN(t *testing.T) {
	//
	// Unsigned
	checkOK(t, uint(12), td.N(uint(12)))
	checkOK(t, uint(11), td.N(uint(12), uint(1)))
	checkOK(t, uint(13), td.N(uint(12), uint(1)))
	checkError(t, 10, td.N(uint(12), uint(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint"),
		})

	checkOK(t, uint8(12), td.N(uint8(12)))
	checkOK(t, uint8(11), td.N(uint8(12), uint8(1)))
	checkOK(t, uint8(13), td.N(uint8(12), uint8(1)))
	checkError(t, 10, td.N(uint8(12), uint8(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint8"),
		})

	checkOK(t, uint16(12), td.N(uint16(12)))
	checkOK(t, uint16(11), td.N(uint16(12), uint16(1)))
	checkOK(t, uint16(13), td.N(uint16(12), uint16(1)))
	checkError(t, 10, td.N(uint16(12), uint16(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint16"),
		})

	checkOK(t, uint32(12), td.N(uint32(12)))
	checkOK(t, uint32(11), td.N(uint32(12), uint32(1)))
	checkOK(t, uint32(13), td.N(uint32(12), uint32(1)))
	checkError(t, 10, td.N(uint32(12), uint32(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint32"),
		})

	checkOK(t, uint64(12), td.N(uint64(12)))
	checkOK(t, uint64(11), td.N(uint64(12), uint64(1)))
	checkOK(t, uint64(13), td.N(uint64(12), uint64(1)))
	checkError(t, 10, td.N(uint64(12), uint64(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint64"),
		})

	checkOK(t, uint64(math.MaxUint64),
		td.N(uint64(math.MaxUint64), uint64(2)))
	checkError(t, uint64(0), td.N(uint64(math.MaxUint64), uint64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				uint64(math.MaxUint64)-2, uint64(math.MaxUint64))),
		})

	checkOK(t, uint64(0), td.N(uint64(0), uint64(2)))
	checkError(t, uint64(math.MaxUint64), td.N(uint64(0), uint64(2)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(fmt.Sprintf("%v", uint64(math.MaxUint64))),
			Expected: mustBe("0 ≤ got ≤ 2"),
		})

	//
	// Signed
	checkOK(t, 12, td.N(12))
	checkOK(t, 11, td.N(12, 1))
	checkOK(t, 13, td.N(12, 1))
	checkError(t, 10, td.N(12, 1),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("11 ≤ got ≤ 13"),
		})

	checkError(t, 10, td.N(12, 0),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("12"),
		})

	checkOK(t, int8(12), td.N(int8(12)))
	checkOK(t, int8(11), td.N(int8(12), int8(1)))
	checkOK(t, int8(13), td.N(int8(12), int8(1)))
	checkError(t, 10, td.N(int8(12), int8(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int8"),
		})

	checkOK(t, int16(12), td.N(int16(12)))
	checkOK(t, int16(11), td.N(int16(12), int16(1)))
	checkOK(t, int16(13), td.N(int16(12), int16(1)))
	checkError(t, 10, td.N(int16(12), int16(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int16"),
		})

	checkOK(t, int32(12), td.N(int32(12)))
	checkOK(t, int32(11), td.N(int32(12), int32(1)))
	checkOK(t, int32(13), td.N(int32(12), int32(1)))
	checkError(t, 10, td.N(int32(12), int32(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int32"),
		})

	checkOK(t, int64(12), td.N(int64(12)))
	checkOK(t, int64(11), td.N(int64(12), int64(1)))
	checkOK(t, int64(13), td.N(int64(12), int64(1)))
	checkError(t, 10, td.N(int64(12), int64(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int64"),
		})

	checkOK(t, int64(math.MaxInt64), td.N(int64(math.MaxInt64), int64(2)))
	checkError(t, int64(0), td.N(int64(math.MaxInt64), int64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				int64(math.MaxInt64)-2, int64(math.MaxInt64))),
		})

	checkOK(t, int64(math.MinInt64), td.N(int64(math.MinInt64), int64(2)))
	checkError(t, int64(0), td.N(int64(math.MinInt64), int64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				int64(math.MinInt64), int64(math.MinInt64)+2)),
		})

	//
	// Float
	checkOK(t, 12.1, td.N(12.1))
	checkOK(t, 11.9, td.N(12.0, 0.1))
	checkOK(t, 12.1, td.N(12.0, 0.1))
	checkError(t, 11.8, td.N(12.0, 0.1),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("11.8"),
			Expected: mustBe("11.9 ≤ got ≤ 12.1"),
		})

	checkOK(t, float32(12.1), td.N(float32(12.1)))
	checkOK(t, float32(11.9), td.N(float32(12), float32(0.1)))
	checkOK(t, float32(12.1), td.N(float32(12), float32(0.1)))
	checkError(t, 11.8, td.N(float32(12), float32(0.1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("float64"),
			Expected: mustBe("float32"),
		})

	floatTol := 10e304
	checkOK(t, float64(math.MaxFloat64),
		td.N(float64(math.MaxFloat64), floatTol))
	checkError(t, float64(0), td.N(float64(math.MaxFloat64), floatTol),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ +Inf",
				float64(math.MaxFloat64)-floatTol)),
		})

	checkOK(t, -float64(math.MaxFloat64),
		td.N(-float64(math.MaxFloat64), float64(2)))
	checkError(t, float64(0), td.N(-float64(math.MaxFloat64), floatTol),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("-Inf ≤ got ≤ %v",
				-float64(math.MaxFloat64)+floatTol)),
		})

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.N("test") }, "usage: N(")
	dark.CheckFatalizerBarrierErr(t, func() { td.N(10, 1, 2) }, "usage: N(")
	dark.CheckFatalizerBarrierErr(t, func() { td.N(10, "test") },
		"N(NUM, TOLERANCE): NUM and TOLERANCE must have the same type: int ≠ string")
}

func TestLGt(t *testing.T) {
	checkOK(t, 12, td.Gt(11))
	checkOK(t, 12, td.Gte(12))
	checkOK(t, 12, td.Lt(13))
	checkOK(t, 12, td.Lte(12))

	checkOK(t, uint16(12), td.Gt(uint16(11)))
	checkOK(t, uint16(12), td.Gte(uint16(12)))
	checkOK(t, uint16(12), td.Lt(uint16(13)))
	checkOK(t, uint16(12), td.Lte(uint16(12)))

	checkOK(t, 12.3, td.Gt(12.2))
	checkOK(t, 12.3, td.Gte(12.3))
	checkOK(t, 12.3, td.Lt(12.4))
	checkOK(t, 12.3, td.Lte(12.3))

	checkOK(t, "abc", td.Gt("abb"))
	checkOK(t, "abc", td.Gte("abc"))
	checkOK(t, "abc", td.Lt("abd"))
	checkOK(t, "abc", td.Lte("abc"))

	checkError(t, 12, td.Gt(12),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("> 12"),
		})
	checkError(t, 12, td.Lt(12),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("< 12"),
		})
	checkError(t, 12, td.Gte(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("≥ 13"),
		})
	checkError(t, 12, td.Lte(11),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("≤ 11"),
		})

	checkError(t, "abc", td.Gt("abc"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`"abc"`),
			Expected: mustBe(`> "abc"`),
		})
	checkError(t, "abc", td.Lt("abc"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`"abc"`),
			Expected: mustBe(`< "abc"`),
		})
	checkError(t, "abc", td.Gte("abd"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`"abc"`),
			Expected: mustBe(`≥ "abd"`),
		})
	checkError(t, "abc", td.Lte("abb"),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(`"abc"`),
			Expected: mustBe(`≤ "abb"`),
		})

	gotDate := time.Date(2018, time.March, 4, 1, 2, 3, 0, time.UTC)
	expectedDate := gotDate
	checkOK(t, gotDate, td.Gte(expectedDate))
	checkOK(t, gotDate, td.Lte(expectedDate))

	checkError(t, gotDate, td.Gt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("> 2018-03-04 01:02:03 +0000 UTC"),
		})
	checkError(t, gotDate, td.Lt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("< 2018-03-04 01:02:03 +0000 UTC"),
		})

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Gt([]byte("test")) }, "usage: Gt(")
	dark.CheckFatalizerBarrierErr(t, func() { td.Gte([]byte("test")) }, "usage: Gte(")
	dark.CheckFatalizerBarrierErr(t, func() { td.Lt([]byte("test")) }, "usage: Lt(")
	dark.CheckFatalizerBarrierErr(t, func() { td.Lte([]byte("test")) }, "usage: Lte(")
}

func TestBetweenTime(t *testing.T) {
	type MyTime time.Time

	now := time.Now()

	checkOK(t, now, td.Between(now, now))
	checkOK(t, now, td.Between(now.Add(-time.Second), now.Add(time.Second)))
	checkOK(t, now, td.Between(now.Add(time.Second), now.Add(-time.Second)))

	checkOK(t, MyTime(now),
		td.Between(
			MyTime(now.Add(-time.Second)),
			MyTime(now.Add(time.Second))))

	// Lax mode
	checkOK(t, MyTime(now),
		td.Lax(td.Between(
			now.Add(time.Second),
			now.Add(-time.Second))))
	checkOK(t, now,
		td.Lax(td.Between(
			MyTime(now.Add(time.Second)),
			MyTime(now.Add(-time.Second)))))

	date := time.Date(2018, time.March, 4, 0, 0, 0, 0, time.UTC)
	checkError(t, date,
		td.Between(date.Add(-2*time.Second), date.Add(-time.Second)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got ≤ " +
				"2018-03-03 23:59:59 +0000 UTC"),
		})

	checkError(t, date,
		td.Between(date.Add(-2*time.Second), date, td.BoundsInOut),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got < " +
				"2018-03-04 00:00:00 +0000 UTC"),
		})

	checkError(t, date,
		td.Between(date, date.Add(2*time.Second), td.BoundsOutIn),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-04 00:00:00 +0000 UTC" +
				" < got ≤ " +
				"2018-03-04 00:00:02 +0000 UTC"),
		})

	checkError(t, "string",
		td.Between(date, date.Add(2*time.Second), td.BoundsOutIn),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("time.Time"),
		})

	checkError(t, "string",
		td.Between(MyTime(date), MyTime(date.Add(2*time.Second)), td.BoundsOutIn),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("td_test.MyTime"),
		})

	checkOK(t, now, td.Gt(now.Add(-time.Second)))
	checkOK(t, now, td.Lt(now.Add(time.Second)))
}

func TestBetweenTypeBehind(t *testing.T) {
	equalTypes(t, td.Between(0, 10), 23)
	equalTypes(t, td.Between(int64(0), int64(10)), int64(23))

	type MyTime time.Time

	equalTypes(t, td.Between(time.Time{}, time.Time{}), time.Time{})
	equalTypes(t, td.Between(MyTime{}, MyTime{}), MyTime{})

	equalTypes(t, td.N(int64(23), int64(5)), int64(0))
	equalTypes(t, td.Gt(int32(23)), int32(0))
	equalTypes(t, td.Gte(int32(23)), int32(0))
	equalTypes(t, td.Lt(int32(23)), int32(0))
	equalTypes(t, td.Lte(int32(23)), int32(0))
}
