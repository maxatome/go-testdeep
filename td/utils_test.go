// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestGetTime(t *testing.T) {
	type MyTime time.Time

	oneTime := time.Date(2018, 7, 14, 12, 11, 10, 0, time.UTC)

	// OK cases
	for idx, curTest := range []struct {
		ParamGot         any
		ParamMustConvert bool
		ExpectedTime     time.Time
	}{
		{
			ParamGot:     oneTime,
			ExpectedTime: oneTime,
		},
		{
			ParamGot:         MyTime(oneTime),
			ParamMustConvert: true,
			ExpectedTime:     oneTime,
		},
	} {
		testName := fmt.Sprintf("Test #%d: ", idx)

		tm, err := getTime(newContext(nil),
			reflect.ValueOf(curTest.ParamGot), curTest.ParamMustConvert)

		if !tm.Equal(curTest.ExpectedTime) {
			test.EqualErrorMessage(t, tm, curTest.ExpectedTime,
				testName+"time")
		}

		if err != nil {
			test.EqualErrorMessage(t, err, "no error",
				testName+"should NOT return an error")
		}
	}

	// Simulate error return from dark.GetInterface
	oldGetInterface := dark.GetInterface
	defer func() { dark.GetInterface = oldGetInterface }()
	dark.GetInterface = func(val reflect.Value, force bool) (any, bool) {
		return nil, false
	}

	// Error cases
	for idx, ctx := range []ctxerr.Context{newContext(nil), newBooleanContext()} {
		testName := fmt.Sprintf("Test #%d: ", idx)

		tm, err := getTime(ctx, reflect.ValueOf(oneTime), false)

		if !tm.Equal(time.Time{}) {
			test.EqualErrorMessage(t, tm, time.Time{}, testName+"time")
		}

		if err == nil {
			test.EqualErrorMessage(t, "no error", err,
				testName+"should return an error")
		}
	}
}
