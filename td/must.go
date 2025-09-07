// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.18
// +build go1.18

package td

import "testing"

// Must is a helper that wraps a call to a function returning (T, error)
// and aborts the current test if the error is non-nil. It is intended
// for use in variable initializations such as:
//
//	tm := td.Must[time.Time](t)(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z"))
//
// As Must uses [testing.TB.FailNow] behind the scenes in case of
// error, it must be called from the goroutine running the test or
// benchmark function, not from other goroutines created during the
// test.
func Must1[X any](t testing.TB) func(v X, err error) X {
	return func(v X, err error) X {
		t.Helper()
		Require(t).CmpNoError(err, "Must")
		return v
	}
}

// tm := td.Must(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z"))(t)
func Must2[X any](v X, err error) func(testing.TB) X {
	if err == nil {
		return func(testing.TB) X { return v }
	}
	return func(t testing.TB) X {
		t.Helper()
		Require(t).CmpNoError(err)
		return v
	}
}

type Called[X any] interface {
	Result() X
	Err() error
	Must(testing.TB) X
}

type called[X any] struct {
	v   X
	err error
}

func (c *called[X]) Err() error {
	return c.err
}

func (c *called[X]) Result() X {
	return c.v
}

func (c *called[X]) Must(t testing.TB) X {
	t.Helper()
	Require(t).CmpNoError(c.err)
	return c.v
}

func Call[X any](v X, err error) Called[X] {
	return &called[X]{v: v, err: err}
}

// tm := Must3(t, td.Call(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")))
func Must3[X any](t testing.TB, res Called[X]) X {
	t.Helper()
	Require(t).CmpNoError(res.Err())
	return res.Result()
}

// tm := assert.Must(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")).(time.Time)
func (t *T) Must(v any, err error) any {
	t.Helper()
	t.Require().CmpNoError(err)
	return v
}
