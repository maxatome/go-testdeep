// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestTuple(t *testing.T) {
	multi := func() (a int, b string, err error) {
		return 12, "test", errors.New("err")
	}

	tuple := td.TupleFrom(multi())
	test.EqualInt(t, tuple.Len(), 3)
	test.EqualInt(t, tuple.Index(0).(int), 12)
	test.EqualStr(t, tuple.Index(1).(string), "test")
	test.EqualStr(t, tuple.Index(2).(error).Error(), "err")

	td.Cmp(t,
		td.TupleFrom(multi()),
		td.TupleFrom(12, "test", td.Not(nil)),
	)

	price := func(p float64) (float64, string, error) {
		if p < 0 {
			return 0, "", errors.New("negative price not supported")
		}
		return p * 1.2, "€", nil
	}

	td.Cmp(t,
		td.TupleFrom(price(10)),
		td.TupleFrom(float64(12), "€", nil),
	)
	td.Cmp(t,
		td.TupleFrom(price(-10)),
		td.TupleFrom(float64(0), "", td.Not(nil)),
	)

	// With Flatten
	td.Cmp(t,
		td.TupleFrom(td.Flatten([]int64{1, 2, 3}), "OK", nil),
		td.TupleFrom(int64(1), int64(2), int64(3), "OK", nil),
	)
}
