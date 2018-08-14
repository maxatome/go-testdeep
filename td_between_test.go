// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
)

func TestBetween(t *testing.T) {
	checkOK(t, 12, testdeep.Between(9, 13))
	checkOK(t, 12, testdeep.Between(13, 9))
	checkOK(t, 12, testdeep.Between(9, 12, testdeep.BoundsOutIn))
	checkOK(t, 12, testdeep.Between(12, 13, testdeep.BoundsInOut))

	checkError(t, 10, testdeep.Between(10, 15, testdeep.BoundsOutIn),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("10 < got ≤ 15"),
		})
	checkError(t, 10, testdeep.Between(10, 15, testdeep.BoundsOutOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("10 < got < 15"),
		})
	checkError(t, 15, testdeep.Between(10, 15, testdeep.BoundsInOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("15"),
			Expected: mustBe("10 ≤ got < 15"),
		})
	checkError(t, 15, testdeep.Between(10, 15, testdeep.BoundsOutOut),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("15"),
			Expected: mustBe("10 < got < 15"),
		})

	checkError(t, 15, testdeep.Between(uint(10), uint(15), testdeep.BoundsOutOut),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint"),
		})

	checkOK(t, uint16(12), testdeep.Between(uint16(9), uint16(13)))
	checkOK(t, uint16(12), testdeep.Between(uint16(13), uint16(9)))
	checkOK(t, uint16(12),
		testdeep.Between(uint16(9), uint16(12), testdeep.BoundsOutIn))
	checkOK(t, uint16(12),
		testdeep.Between(uint16(12), uint16(13), testdeep.BoundsInOut))

	checkOK(t, 12.1, testdeep.Between(9.5, 13.1))
	checkOK(t, 12.1, testdeep.Between(13.1, 9.5))
	checkOK(t, 12.1, testdeep.Between(9.5, 12.1, testdeep.BoundsOutIn))
	checkOK(t, 12.1, testdeep.Between(12.1, 13.1, testdeep.BoundsInOut))

	checkOK(t, 12*time.Hour, testdeep.Between(60*time.Second, 24*time.Hour))

	//
	// Bad usage
	checkPanic(t, func() { testdeep.Between("test", "test") },
		"usage: Between(")
	checkPanic(t, func() { testdeep.Between(12, "test") },
		"from and to params must have the same type")
	checkPanic(t, func() { testdeep.Between("test", 12) },
		"from and to params must have the same type")
	checkPanic(t, func() { testdeep.Between(1, 2, testdeep.BoundsInIn, testdeep.BoundsInOut) },
		"usage: Between(")
}

func TestN(t *testing.T) {
	//
	// Unsigned
	checkOK(t, uint(12), testdeep.N(uint(12)))
	checkOK(t, uint(11), testdeep.N(uint(12), uint(1)))
	checkOK(t, uint(13), testdeep.N(uint(12), uint(1)))
	checkError(t, 10, testdeep.N(uint(12), uint(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint"),
		})

	checkOK(t, uint8(12), testdeep.N(uint8(12)))
	checkOK(t, uint8(11), testdeep.N(uint8(12), uint8(1)))
	checkOK(t, uint8(13), testdeep.N(uint8(12), uint8(1)))
	checkError(t, 10, testdeep.N(uint8(12), uint8(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint8"),
		})

	checkOK(t, uint16(12), testdeep.N(uint16(12)))
	checkOK(t, uint16(11), testdeep.N(uint16(12), uint16(1)))
	checkOK(t, uint16(13), testdeep.N(uint16(12), uint16(1)))
	checkError(t, 10, testdeep.N(uint16(12), uint16(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint16"),
		})

	checkOK(t, uint32(12), testdeep.N(uint32(12)))
	checkOK(t, uint32(11), testdeep.N(uint32(12), uint32(1)))
	checkOK(t, uint32(13), testdeep.N(uint32(12), uint32(1)))
	checkError(t, 10, testdeep.N(uint32(12), uint32(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint32"),
		})

	checkOK(t, uint64(12), testdeep.N(uint64(12)))
	checkOK(t, uint64(11), testdeep.N(uint64(12), uint64(1)))
	checkOK(t, uint64(13), testdeep.N(uint64(12), uint64(1)))
	checkError(t, 10, testdeep.N(uint64(12), uint64(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("uint64"),
		})

	checkOK(t, uint64(math.MaxUint64),
		testdeep.N(uint64(math.MaxUint64), uint64(2)))
	checkError(t, uint64(0), testdeep.N(uint64(math.MaxUint64), uint64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				uint64(math.MaxUint64)-2, uint64(math.MaxUint64))),
		})

	checkOK(t, uint64(0), testdeep.N(uint64(0), uint64(2)))
	checkError(t, uint64(math.MaxUint64), testdeep.N(uint64(0), uint64(2)),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe(fmt.Sprintf("%v", uint64(math.MaxUint64))),
			Expected: mustBe("0 ≤ got ≤ 2"),
		})

	//
	// Signed
	checkOK(t, 12, testdeep.N(12))
	checkOK(t, 11, testdeep.N(12, 1))
	checkOK(t, 13, testdeep.N(12, 1))
	checkError(t, 10, testdeep.N(12, 1),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("11 ≤ got ≤ 13"),
		})

	checkError(t, 10, testdeep.N(12, 0),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("10"),
			Expected: mustBe("12"),
		})

	checkOK(t, int8(12), testdeep.N(int8(12)))
	checkOK(t, int8(11), testdeep.N(int8(12), int8(1)))
	checkOK(t, int8(13), testdeep.N(int8(12), int8(1)))
	checkError(t, 10, testdeep.N(int8(12), int8(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int8"),
		})

	checkOK(t, int16(12), testdeep.N(int16(12)))
	checkOK(t, int16(11), testdeep.N(int16(12), int16(1)))
	checkOK(t, int16(13), testdeep.N(int16(12), int16(1)))
	checkError(t, 10, testdeep.N(int16(12), int16(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int16"),
		})

	checkOK(t, int32(12), testdeep.N(int32(12)))
	checkOK(t, int32(11), testdeep.N(int32(12), int32(1)))
	checkOK(t, int32(13), testdeep.N(int32(12), int32(1)))
	checkError(t, 10, testdeep.N(int32(12), int32(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int32"),
		})

	checkOK(t, int64(12), testdeep.N(int64(12)))
	checkOK(t, int64(11), testdeep.N(int64(12), int64(1)))
	checkOK(t, int64(13), testdeep.N(int64(12), int64(1)))
	checkError(t, 10, testdeep.N(int64(12), int64(1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("int"),
			Expected: mustBe("int64"),
		})

	checkOK(t, int64(math.MaxInt64), testdeep.N(int64(math.MaxInt64), int64(2)))
	checkError(t, int64(0), testdeep.N(int64(math.MaxInt64), int64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				int64(math.MaxInt64)-2, int64(math.MaxInt64))),
		})

	checkOK(t, int64(math.MinInt64), testdeep.N(int64(math.MinInt64), int64(2)))
	checkError(t, int64(0), testdeep.N(int64(math.MinInt64), int64(2)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ %v",
				int64(math.MinInt64), int64(math.MinInt64)+2)),
		})

	//
	// Float
	checkOK(t, 12.1, testdeep.N(12.1))
	checkOK(t, 11.9, testdeep.N(12.0, 0.1))
	checkOK(t, 12.1, testdeep.N(12.0, 0.1))
	checkError(t, 11.8, testdeep.N(12.0, 0.1),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("11.8"),
			Expected: mustBe("11.9 ≤ got ≤ 12.1"),
		})

	checkOK(t, float32(12.1), testdeep.N(float32(12.1)))
	checkOK(t, float32(11.9), testdeep.N(float32(12), float32(0.1)))
	checkOK(t, float32(12.1), testdeep.N(float32(12), float32(0.1)))
	checkError(t, 11.8, testdeep.N(float32(12), float32(0.1)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("float64"),
			Expected: mustBe("float32"),
		})

	floatTol := 10e304
	checkOK(t, float64(math.MaxFloat64),
		testdeep.N(float64(math.MaxFloat64), floatTol))
	checkError(t, float64(0), testdeep.N(float64(math.MaxFloat64), floatTol),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("%v ≤ got ≤ +Inf",
				float64(math.MaxFloat64)-floatTol)),
		})

	checkOK(t, -float64(math.MaxFloat64),
		testdeep.N(-float64(math.MaxFloat64), float64(2)))
	checkError(t, float64(0), testdeep.N(-float64(math.MaxFloat64), floatTol),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("0"),
			Expected: mustBe(fmt.Sprintf("-Inf ≤ got ≤ %v",
				-float64(math.MaxFloat64)+floatTol)),
		})

	//
	// Bad usage
	checkPanic(t, func() { testdeep.N("test") }, "usage: N(")
	checkPanic(t, func() { testdeep.N(10, 1, 2) }, "usage: N(")
	checkPanic(t, func() { testdeep.N(10, "test") },
		"tolerance param must have the same type as num one")
}

func TestLGt(t *testing.T) {
	checkOK(t, 12, testdeep.Gt(11))
	checkOK(t, 12, testdeep.Gte(12))
	checkOK(t, 12, testdeep.Lt(13))
	checkOK(t, 12, testdeep.Lte(12))

	checkOK(t, uint16(12), testdeep.Gt(uint16(11)))
	checkOK(t, uint16(12), testdeep.Gte(uint16(12)))
	checkOK(t, uint16(12), testdeep.Lt(uint16(13)))
	checkOK(t, uint16(12), testdeep.Lte(uint16(12)))

	checkOK(t, 12.3, testdeep.Gt(12.2))
	checkOK(t, 12.3, testdeep.Gte(12.3))
	checkOK(t, 12.3, testdeep.Lt(12.4))
	checkOK(t, 12.3, testdeep.Lte(12.3))

	checkError(t, 12, testdeep.Gt(12),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("> 12"),
		})
	checkError(t, 12, testdeep.Lt(12),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("< 12"),
		})
	checkError(t, 12, testdeep.Gte(13),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("≥ 13"),
		})
	checkError(t, 12, testdeep.Lte(11),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("12"),
			Expected: mustBe("≤ 11"),
		})

	gotDate := time.Date(2018, time.March, 4, 1, 2, 3, 0, time.UTC)
	expectedDate := gotDate
	checkOK(t, gotDate, testdeep.Gte(expectedDate))
	checkOK(t, gotDate, testdeep.Lte(expectedDate))

	checkError(t, gotDate, testdeep.Gt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("> 2018-03-04 01:02:03 +0000 UTC"),
		})
	checkError(t, gotDate, testdeep.Lt(expectedDate),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustBe("2018-03-04 01:02:03 +0000 UTC"),
			Expected: mustBe("< 2018-03-04 01:02:03 +0000 UTC"),
		})

	//
	// Bad usage
	checkPanic(t, func() { testdeep.Gt("test") }, "usage: Gt(")
	checkPanic(t, func() { testdeep.Gte("test") }, "usage: Gte(")
	checkPanic(t, func() { testdeep.Lt("test") }, "usage: Lt(")
	checkPanic(t, func() { testdeep.Lte("test") }, "usage: Lte(")
}

func TestBetweenTime(t *testing.T) {
	type MyTime time.Time

	now := time.Now()

	checkOK(t, now, testdeep.Between(now, now))
	checkOK(t, now, testdeep.Between(now.Add(-time.Second), now.Add(time.Second)))
	checkOK(t, now, testdeep.Between(now.Add(time.Second), now.Add(-time.Second)))

	checkOK(t, MyTime(now),
		testdeep.Between(MyTime(now.Add(-time.Second)),
			MyTime(now.Add(time.Second))))

	checkOK(t, MyTime(now),
		testdeep.Between(MyTime(now.Add(time.Second)),
			MyTime(now.Add(-time.Second))))

	date := time.Date(2018, time.March, 4, 0, 0, 0, 0, time.UTC)
	checkError(t, date,
		testdeep.Between(date.Add(-2*time.Second), date.Add(-time.Second)),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got ≤ " +
				"2018-03-03 23:59:59 +0000 UTC"),
		})

	checkError(t, date,
		testdeep.Between(date.Add(-2*time.Second), date, testdeep.BoundsInOut),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-03 23:59:58 +0000 UTC" +
				" ≤ got < " +
				"2018-03-04 00:00:00 +0000 UTC"),
		})

	checkError(t, date,
		testdeep.Between(date, date.Add(2*time.Second), testdeep.BoundsOutIn),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got:     mustBe("2018-03-04 00:00:00 +0000 UTC"),
			Expected: mustBe("2018-03-04 00:00:00 +0000 UTC" +
				" < got ≤ " +
				"2018-03-04 00:00:02 +0000 UTC"),
		})

	checkError(t, "string",
		testdeep.Between(date, date.Add(2*time.Second), testdeep.BoundsOutIn),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("time.Time"),
		})

	checkError(t, "string",
		testdeep.Between(MyTime(date), MyTime(date.Add(2*time.Second)), testdeep.BoundsOutIn),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("string"),
			Expected: mustBe("testdeep_test.MyTime"),
		})

	checkOK(t, now, testdeep.Gt(now.Add(-time.Second)))
	checkOK(t, now, testdeep.Lt(now.Add(time.Second)))
}

func TestBetweenTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Between(0, 10), 23)
	equalTypes(t, testdeep.Between(int64(0), int64(10)), int64(23))

	type MyTime time.Time

	equalTypes(t, testdeep.Between(time.Time{}, time.Time{}), time.Time{})
	equalTypes(t, testdeep.Between(MyTime{}, MyTime{}), MyTime{})

	equalTypes(t, testdeep.N(int64(23), int64(5)), int64(0))
	equalTypes(t, testdeep.Gt(int32(23)), int32(0))
	equalTypes(t, testdeep.Gte(int32(23)), int32(0))
	equalTypes(t, testdeep.Lt(int32(23)), int32(0))
	equalTypes(t, testdeep.Lte(int32(23)), int32(0))
}
