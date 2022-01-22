// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

type (
	MyTime    time.Time
	MyTimeStr time.Time
)

func (t MyTimeStr) String() string {
	return "<<" + time.Time(t).Format(time.RFC3339Nano) + ">>"
}

func TestTruncTime(t *testing.T) {
	//
	// Monotonic
	now := time.Now()
	nowWithoutMono := now.Truncate(0)

	// If monotonic clock available, check without TruncTime()
	if now != nowWithoutMono {
		// OK now contains a monotonic part != 0, so fail coz "==" used inside
		checkError(t, now, nowWithoutMono,
			expectedError{
				Message: mustBe("values differ"),
				Path:    mustContain("DATA"),
			})
	}
	checkOK(t, now, td.TruncTime(nowWithoutMono))

	//
	// time.Time
	gotDate := time.Date(2018, time.March, 9, 1, 2, 3, 4, time.UTC)

	// Time zone / location does not matter
	UTCp2 := time.FixedZone("UTC+2", 2)
	UTCm2 := time.FixedZone("UTC-2", 2)
	checkOK(t, gotDate, td.TruncTime(gotDate.In(UTCp2)))
	checkOK(t, gotDate, td.TruncTime(gotDate.In(UTCm2)))

	checkOK(t, gotDate.In(UTCm2), td.TruncTime(gotDate.In(UTCp2)))
	checkOK(t, gotDate.In(UTCp2), td.TruncTime(gotDate.In(UTCm2)))

	expDate := gotDate

	checkOK(t, gotDate, td.TruncTime(expDate))
	checkOK(t, gotDate, td.TruncTime(expDate, time.Second))
	checkOK(t, gotDate, td.TruncTime(expDate, time.Minute))

	expDate = expDate.Add(time.Second)
	checkError(t, gotDate, td.TruncTime(expDate, time.Second),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got: mustBe("2018-03-09 01:02:03.000000004 +0000 UTC\n" +
				"truncated to:\n" +
				"2018-03-09 01:02:03 +0000 UTC"),
			Expected: mustBe("2018-03-09 01:02:04 +0000 UTC"),
		})
	checkOK(t, gotDate, td.TruncTime(expDate, time.Minute))

	checkError(t, gotDate, td.TruncTime(MyTime(gotDate)),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("time.Time"),
			Expected: mustBe("td_test.MyTime"),
		})

	//
	// Type convertible to time.Time NOT implementing fmt.Stringer
	gotMyDate := MyTime(gotDate)
	expMyDate := MyTime(gotDate)

	checkOK(t, gotMyDate, td.TruncTime(expMyDate))
	checkOK(t, gotMyDate, td.TruncTime(expMyDate, time.Second))
	checkOK(t, gotMyDate, td.TruncTime(expMyDate, time.Minute))

	expMyDate = MyTime(gotDate.Add(time.Second))
	checkError(t, gotMyDate, td.TruncTime(expMyDate, time.Second),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got: mustBe("2018-03-09 01:02:03.000000004 +0000 UTC\n" +
				"truncated to:\n" +
				"2018-03-09 01:02:03 +0000 UTC"),
			Expected: mustBe("2018-03-09 01:02:04 +0000 UTC"),
		})
	checkOK(t, gotMyDate, td.TruncTime(expMyDate, time.Minute))

	checkError(t, MyTime(gotDate), td.TruncTime(gotDate),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MyTime"),
			Expected: mustBe("time.Time"),
		})

	//
	// Type convertible to time.Time implementing fmt.Stringer
	gotMyStrDate := MyTimeStr(gotDate)
	expMyStrDate := MyTimeStr(gotDate)

	checkOK(t, gotMyStrDate, td.TruncTime(expMyStrDate))
	checkOK(t, gotMyStrDate, td.TruncTime(expMyStrDate, time.Second))
	checkOK(t, gotMyStrDate, td.TruncTime(expMyStrDate, time.Minute))

	expMyStrDate = MyTimeStr(gotDate.Add(time.Second))
	checkError(t, gotMyStrDate, td.TruncTime(expMyStrDate, time.Second),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got: mustBe("<<2018-03-09T01:02:03.000000004Z>>\n" +
				"truncated to:\n" +
				"<<2018-03-09T01:02:03Z>>"),
			Expected: mustBe("<<2018-03-09T01:02:04Z>>"),
		})
	checkOK(t, gotMyStrDate, td.TruncTime(expMyStrDate, time.Minute))

	checkError(t, MyTimeStr(gotDate), td.TruncTime(gotDate),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("td_test.MyTimeStr"),
			Expected: mustBe("time.Time"),
		})

	//
	// Bad usage
	checkError(t, "never tested",
		td.TruncTime("test"),
		expectedError{
			Message: mustBe("bad usage of TruncTime operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: TruncTime(time.Time[, time.Duration]), 1st parameter must be time.Time or convertible to time.Time, but not string"),
		})

	checkError(t, "never tested",
		td.TruncTime(1, 2, 3),
		expectedError{
			Message: mustBe("bad usage of TruncTime operator"),
			Path:    mustBe("DATA"),
			Summary: mustBe("usage: TruncTime(time.Time[, time.Duration]), too many parameters"),
		})

	// Erroneous op
	test.EqualStr(t, td.TruncTime("test").String(), "TruncTime(<ERROR>)")
}

func TestTruncTimeTypeBehind(t *testing.T) {
	type MyTime time.Time

	equalTypes(t, td.TruncTime(time.Time{}), time.Time{})
	equalTypes(t, td.TruncTime(MyTime{}), MyTime{})

	// Erroneous op
	equalTypes(t, td.TruncTime("test"), nil)
}
