// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"
	"time"

	. "github.com/maxatome/go-testdeep"
)

type MyTime time.Time
type MyTimeStr time.Time

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
		checkError(t, now, nowWithoutMono, expectedError{
			Message: mustBe("values differ"),
		})
	}
	checkOK(t, now, TruncTime(nowWithoutMono))

	//
	// time.Time
	gotDate := time.Date(2018, time.March, 9, 1, 2, 3, 4, time.UTC)
	expDate := gotDate

	checkOK(t, gotDate, TruncTime(expDate))
	checkOK(t, gotDate, TruncTime(expDate, time.Second))
	checkOK(t, gotDate, TruncTime(expDate, time.Minute))

	expDate = expDate.Add(time.Second)
	checkError(t, gotDate, TruncTime(expDate, time.Second), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("DATA"),
		Got: mustBe("2018-03-09 01:02:03.000000004 +0000 UTC\n" +
			"truncated to:\n" +
			"2018-03-09 01:02:03 +0000 UTC"),
		Expected: mustBe("2018-03-09 01:02:04 +0000 UTC"),
	})
	checkOK(t, gotDate, TruncTime(expDate, time.Minute))

	checkError(t, gotDate, TruncTime(MyTime(gotDate)), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("time.Time"),
		Expected: mustBe("testdeep_test.MyTime"),
	})

	//
	// Type convertible to time.Time NOT implementing fmt.Stringer
	gotMyDate := MyTime(gotDate)
	expMyDate := MyTime(gotDate)

	checkOK(t, gotMyDate, TruncTime(expMyDate))
	checkOK(t, gotMyDate, TruncTime(expMyDate, time.Second))
	checkOK(t, gotMyDate, TruncTime(expMyDate, time.Minute))

	expMyDate = MyTime(gotDate.Add(time.Second))
	checkError(t, gotMyDate, TruncTime(expMyDate, time.Second), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("DATA"),
		Got: mustBe("2018-03-09 01:02:03.000000004 +0000 UTC\n" +
			"truncated to:\n" +
			"2018-03-09 01:02:03 +0000 UTC"),
		Expected: mustBe("2018-03-09 01:02:04 +0000 UTC"),
	})
	checkOK(t, gotMyDate, TruncTime(expMyDate, time.Minute))

	checkError(t, MyTime(gotDate), TruncTime(gotDate), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("testdeep_test.MyTime"),
		Expected: mustBe("time.Time"),
	})

	//
	// Type convertible to time.Time implementing fmt.Stringer
	gotMyStrDate := MyTimeStr(gotDate)
	expMyStrDate := MyTimeStr(gotDate)

	checkOK(t, gotMyStrDate, TruncTime(expMyStrDate))
	checkOK(t, gotMyStrDate, TruncTime(expMyStrDate, time.Second))
	checkOK(t, gotMyStrDate, TruncTime(expMyStrDate, time.Minute))

	expMyStrDate = MyTimeStr(gotDate.Add(time.Second))
	checkError(t, gotMyStrDate, TruncTime(expMyStrDate, time.Second),
		expectedError{
			Message: mustBe("values differ"),
			Path:    mustBe("DATA"),
			Got: mustBe("<<2018-03-09T01:02:03.000000004Z>>\n" +
				"truncated to:\n" +
				"<<2018-03-09T01:02:03Z>>"),
			Expected: mustBe("<<2018-03-09T01:02:04Z>>"),
		})
	checkOK(t, gotMyStrDate, TruncTime(expMyStrDate, time.Minute))

	checkError(t, MyTimeStr(gotDate), TruncTime(gotDate), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustBe("testdeep_test.MyTimeStr"),
		Expected: mustBe("time.Time"),
	})

	//
	// Bad usage
	checkPanic(t, func() { TruncTime("test") }, "usage: TruncTime(")
}
