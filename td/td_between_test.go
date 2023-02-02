// Copyright (c) 2018-2022, Maxime Soulé
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

	"github.com/maxatome/go-testdeep/internal/test"
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
	checkError(t, "never tested",
		td.Between([]byte("test"), []byte("test")),
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Between(NUM|STRING|TIME, NUM|STRING|TIME/DURATION[, BOUNDS_KIND]), but received []uint8 (slice) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Between(12, "test"),
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Between(FROM, TO): FROM and TO must have the same type: int ≠ string"),
		})

	checkError(t, "never tested",
		td.Between("test", 12),
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Between(FROM, TO): FROM and TO must have the same type: string ≠ int"),
		})

	checkError(t, "never tested",
		td.Between(1, 2, td.BoundsInIn, td.BoundsInOut),
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Between(NUM|STRING|TIME, NUM|STRING|TIME/DURATION[, BOUNDS_KIND]), too many parameters"),
		})

	type notTime struct{}
	checkError(t, "never tested",
		td.Between(notTime{}, notTime{}),
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Between(NUM|STRING|TIME, NUM|STRING|TIME/DURATION[, BOUNDS_KIND]), but received td_test.notTime (struct) as 1st parameter"),
		})

	// Erroneous op
	test.EqualStr(t, td.Between("test", 12).String(), "Between(<ERROR>)")
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
			Got:     mustBe("(uint64) 0"),
			Expected: mustBe(fmt.Sprintf("(uint64) %v ≤ got ≤ (uint64) %v",
				uint64(math.MaxUint64)-2, uint64(math.MaxUint64))),
		})

	checkOK(t, uint64(0), td.N(uint64(0), uint64(2)))
	checkError(t, uint64(math.MaxUint64), td.N(uint64(0), uint64(2)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(fmt.Sprintf("(uint64) %v", uint64(math.MaxUint64))),
			Expected: mustBe("(uint64) 0 ≤ got ≤ (uint64) 2"),
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
			Expected: mustBe("12 ≤ got ≤ 12"),
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
			Got:     mustBe("(int64) 0"),
			Expected: mustBe(fmt.Sprintf("(int64) %v ≤ got ≤ (int64) %v",
				int64(math.MaxInt64)-2, int64(math.MaxInt64))),
		})

	checkOK(t, int64(math.MinInt64), td.N(int64(math.MinInt64), int64(2)))
	checkError(t, int64(0), td.N(int64(math.MinInt64), int64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("(int64) 0"),
			Expected: mustBe(fmt.Sprintf("(int64) %v ≤ got ≤ (int64) %v",
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
			Got:     mustBe("0.0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ +Inf",
				float64(math.MaxFloat64)-floatTol)),
		})

	checkOK(t, -float64(math.MaxFloat64),
		td.N(-float64(math.MaxFloat64), float64(2)))
	checkError(t, float64(0), td.N(-float64(math.MaxFloat64), floatTol),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0.0"),
			Expected: mustBe(fmt.Sprintf("-Inf ≤ got ≤ %v",
				-float64(math.MaxFloat64)+floatTol)),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.N("test"),
		expectedError{
			Message: mustBe("bad usage of N operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: N({,U}INT{,8,16,32,64}|FLOAT{32,64}[, TOLERANCE]), but received string as 1st parameter"),
		})

	checkError(t, "never tested",
		td.N(10, 1, 2),
		expectedError{
			Message: mustBe("bad usage of N operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: N({,U}INT{,8,16,32,64}|FLOAT{32,64}[, TOLERANCE]), too many parameters"),
		})

	checkError(t, "never tested",
		td.N(10, "test"),
		expectedError{
			Message: mustBe("bad usage of N operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("N(NUM, TOLERANCE): NUM and TOLERANCE must have the same type: int ≠ string"),
		})

	// Erroneous op
	test.EqualStr(t, td.N(10, 1, 2).String(), "N(<ERROR>)")
}

func TestLGt(t *testing.T) {
	type MyTime time.Time

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
	checkOK(t, gotDate, td.Lax(td.Gte(MyTime(expectedDate))))
	checkOK(t, gotDate, td.Lax(td.Lte(MyTime(expectedDate))))

	checkError(t, gotDate, td.Gt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(time.Time) 2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("> (time.Time) 2018-03-04 01:02:03 +0000 UTC"),
		})
	checkError(t, gotDate, td.Lt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(time.Time) 2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("< (time.Time) 2018-03-04 01:02:03 +0000 UTC"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.Gt([]byte("test")),
		expectedError{
			Message: mustBe("bad usage of Gt operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Gt(NUM|STRING|TIME), but received []uint8 (slice) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Gte([]byte("test")),
		expectedError{
			Message: mustBe("bad usage of Gte operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Gte(NUM|STRING|TIME), but received []uint8 (slice) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Lt([]byte("test")),
		expectedError{
			Message: mustBe("bad usage of Lt operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Lt(NUM|STRING|TIME), but received []uint8 (slice) as 1st parameter"),
		})

	checkError(t, "never tested",
		td.Lte([]byte("test")),
		expectedError{
			Message: mustBe("bad usage of Lte operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: Lte(NUM|STRING|TIME), but received []uint8 (slice) as 1st parameter"),
		})

	// Erroneous op
	test.EqualStr(t, td.Gt([]byte("test")).String(), "Gt(<ERROR>)")
	test.EqualStr(t, td.Gte([]byte("test")).String(), "Gte(<ERROR>)")
	test.EqualStr(t, td.Lt([]byte("test")).String(), "Lt(<ERROR>)")
	test.EqualStr(t, td.Lte([]byte("test")).String(), "Lte(<ERROR>)")
}

func TestBetweenTime(t *testing.T) {
	type MyTime time.Time

	now := time.Now()

	checkOK(t, now, td.Between(now, now))
	checkOK(t, now, td.Between(now.Add(-time.Second), now.Add(time.Second)))
	checkOK(t, now, td.Between(now.Add(time.Second), now.Add(-time.Second)))

	// (TIME, DURATION)
	checkOK(t, now, td.Between(now.Add(-time.Second), 2*time.Second))
	checkOK(t, now, td.Between(now.Add(time.Second), -2*time.Second))

	checkOK(t, MyTime(now),
		td.Between(
			MyTime(now.Add(-time.Second)),
			MyTime(now.Add(time.Second))))

	// (TIME, DURATION)
	checkOK(t, MyTime(now),
		td.Between(
			MyTime(now.Add(-time.Second)),
			2*time.Second))
	checkOK(t, MyTime(now),
		td.Between(
			MyTime(now.Add(time.Second)),
			-2*time.Second))

	// Lax mode
	checkOK(t, MyTime(now),
		td.Lax(td.Between(
			now.Add(time.Second),
			now.Add(-time.Second))))
	checkOK(t, now,
		td.Lax(td.Between(
			MyTime(now.Add(time.Second)),
			MyTime(now.Add(-time.Second)))))
	checkOK(t, MyTime(now),
		td.Lax(td.Between(
			now.Add(-time.Second),
			2*time.Second)))
	checkOK(t, now,
		td.Lax(td.Between(
			MyTime(now.Add(-time.Second)),
			2*time.Second)))
	checkOK(t, now,
		td.Lax(td.Between(
			MyTime(now.Add(-time.Second)),
			2*time.Second,
			td.BoundsOutOut)))

	date := time.Date(2018, time.March, 4, 0, 0, 0, 0, time.UTC)
	checkError(t, date,
		td.Between(date.Add(-2*time.Second), date.Add(-time.Second)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("(time.Time) 2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("(time.Time) 2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got ≤ " +
				"(time.Time) 2018-03-03 23:59:59 +0000 UTC"),
		})

	checkError(t, MyTime(date),
		td.Between(MyTime(date.Add(-2*time.Second)), MyTime(date.Add(-time.Second))),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustContain("(td_test.MyTime) "),
			Expected: mustBe(
				"(time.Time) 2018-03-03 23:59:58 +0000 UTC" +
					" ≤ got ≤ " +
					"(time.Time) 2018-03-03 23:59:59 +0000 UTC"),
		})

	checkError(t, date,
		td.Between(date.Add(-2*time.Second), date, td.BoundsInOut),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("(time.Time) 2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("(time.Time) 2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got < " +
				"(time.Time) 2018-03-04 00:00:00 +0000 UTC"),
		})

	checkError(t, date,
		td.Between(date, date.Add(2*time.Second), td.BoundsOutIn),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("(time.Time) 2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("(time.Time) 2018-03-04 00:00:00 +0000 UTC" +
				" < got ≤ " +
				"(time.Time) 2018-03-04 00:00:02 +0000 UTC"),
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

	checkError(t, "never tested",
		td.Between(date, 12), // (Time, Time) or (Time, Duration)
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Between(FROM, TO): when FROM type is time.Time, TO must have the same type or time.Duration: int ≠ time.Time|time.Duration"),
		})

	checkError(t, "never tested",
		td.Between(MyTime(date), 12), // (MyTime, MyTime) or (MyTime, Duration)
		expectedError{
			Message: mustBe("bad usage of Between operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("Between(FROM, TO): when FROM type is td_test.MyTime, TO must have the same type or time.Duration: int ≠ td_test.MyTime|time.Duration"),
		})

	checkOK(t, now, td.Gt(now.Add(-time.Second)))
	checkOK(t, now, td.Lt(now.Add(time.Second)))
}

type compareType int

func (i compareType) Compare(j compareType) int {
	if i < j {
		return -1
	}
	if i > j {
		return 1
	}
	return 0
}

type lessType int

func (i lessType) Less(j lessType) bool {
	return i < j
}

func TestBetweenCmp(t *testing.T) {
	t.Run("compareType", func(t *testing.T) {
		checkOK(t, compareType(5), td.Between(compareType(4), compareType(6)))
		checkOK(t, compareType(5), td.Between(compareType(6), compareType(4)))
		checkOK(t, compareType(5), td.Between(compareType(5), compareType(6)))
		checkOK(t, compareType(5), td.Between(compareType(4), compareType(5)))

		checkOK(t, compareType(5),
			td.Between(compareType(4), compareType(6), td.BoundsOutOut))
		checkError(t, compareType(5),
			td.Between(compareType(5), compareType(6), td.BoundsOutIn),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("(td_test.compareType) 5"),
				Expected: mustBe("(td_test.compareType) 5 < got ≤ (td_test.compareType) 6"),
			})
		checkError(t, compareType(5),
			td.Between(compareType(4), compareType(5), td.BoundsInOut),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("(td_test.compareType) 5"),
				Expected: mustBe("(td_test.compareType) 4 ≤ got < (td_test.compareType) 5"),
			})

		// Other between forms
		checkOK(t, compareType(5), td.Gt(compareType(4)))
		checkOK(t, compareType(5), td.Gte(compareType(5)))
		checkOK(t, compareType(5), td.Lt(compareType(6)))
		checkOK(t, compareType(5), td.Lte(compareType(5)))

		// BeLax or not BeLax
		for i, op := range []td.TestDeep{
			td.Between(compareType(4), compareType(6)),
			td.Gt(compareType(4)),
			td.Gte(compareType(5)),
			td.Lt(compareType(6)),
			td.Lte(compareType(5)),
		} {
			// Type mismatch if BeLax not enabled
			checkError(t, 5, op,
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe("DATA"),
					Got:      mustBe("int"),
					Expected: mustBe("td_test.compareType"),
				},
				"Op #%d", i)

			// BeLax enabled is OK
			checkOK(t, 5, td.Lax(op), "Op #%d", i)
		}

		// In a private field
		type private struct {
			num compareType
		}
		checkOK(t, private{num: 5},
			td.Struct(private{},
				td.StructFields{
					"num": td.Between(compareType(4), compareType(6)),
				}))
	})

	t.Run("lessType", func(t *testing.T) {
		checkOK(t, lessType(5), td.Between(lessType(4), lessType(6)))
		checkOK(t, lessType(5), td.Between(lessType(6), lessType(4)))
		checkOK(t, lessType(5), td.Between(lessType(5), lessType(6)))
		checkOK(t, lessType(5), td.Between(lessType(4), lessType(5)))

		checkOK(t, lessType(5),
			td.Between(lessType(4), lessType(6), td.BoundsOutOut))
		checkError(t, lessType(5),
			td.Between(lessType(5), lessType(6), td.BoundsOutIn),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("(td_test.lessType) 5"),
				Expected: mustBe("(td_test.lessType) 5 < got ≤ (td_test.lessType) 6"),
			})
		checkError(t, lessType(5),
			td.Between(lessType(4), lessType(5), td.BoundsInOut),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("DATA"),
				Got:      mustBe("(td_test.lessType) 5"),
				Expected: mustBe("(td_test.lessType) 4 ≤ got < (td_test.lessType) 5"),
			})

		// Other between forms
		checkOK(t, lessType(5), td.Gt(lessType(4)))
		checkOK(t, lessType(5), td.Gte(lessType(5)))
		checkOK(t, lessType(5), td.Lt(lessType(6)))
		checkOK(t, lessType(5), td.Lte(lessType(5)))

		// BeLax or not BeLax
		for i, op := range []td.TestDeep{
			td.Between(lessType(4), lessType(6)),
			td.Gt(lessType(4)),
			td.Gte(lessType(5)),
			td.Lt(lessType(6)),
			td.Lte(lessType(5)),
		} {
			// Type mismatch if BeLax not enabled
			checkError(t, 5, op,
				expectedError{
					Message:  mustBe("type mismatch"),
					Path:     mustBe("DATA"),
					Got:      mustBe("int"),
					Expected: mustBe("td_test.lessType"),
				},
				"Op #%d", i)

			// BeLax enabled is OK
			checkOK(t, 5, td.Lax(op), "Op #%d", i)
		}

		// In a private field
		type private struct {
			num lessType
		}
		checkOK(t, private{num: 5},
			td.Struct(private{},
				td.StructFields{
					"num": td.Between(lessType(4), lessType(6)),
				}))
	})
}

func TestBetweenTypeBehind(t *testing.T) {
	type MyTime time.Time

	for _, typ := range []any{
		10,
		int64(23),
		int32(23),
		time.Time{},
		MyTime{},
		compareType(0),
		lessType(0),
	} {
		equalTypes(t, td.Between(typ, typ), typ)
		equalTypes(t, td.Gt(typ), typ)
		equalTypes(t, td.Gte(typ), typ)
		equalTypes(t, td.Lt(typ), typ)
		equalTypes(t, td.Lte(typ), typ)
	}
	equalTypes(t, td.N(int64(23), int64(5)), int64(0))

	// Erroneous op
	equalTypes(t, td.Between("test", 12), nil)
	equalTypes(t, td.N(10, 1, 2), nil)
	equalTypes(t, td.Gt([]byte("test")), nil)
	equalTypes(t, td.Gte([]byte("test")), nil)
	equalTypes(t, td.Lt([]byte("test")), nil)
	equalTypes(t, td.Lte([]byte("test")), nil)
}
