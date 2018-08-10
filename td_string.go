// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"
	"strings"
)

type tdStringBase struct {
	Base
	expected string
}

func newStringBase(expected string) tdStringBase {
	return tdStringBase{
		Base:     NewBase(4),
		expected: expected,
	}
}

func getString(ctx Context, got reflect.Value) (string, *Error) {
	switch got.Kind() {
	case reflect.String:
		return got.String(), nil

	default:
		if got.CanInterface() {
			switch iface := got.Interface().(type) {
			case error:
				return iface.Error(), nil
			case fmt.Stringer:
				return iface.String(), nil
			}
		}
	}

	if ctx.BooleanError {
		return "", BooleanError
	}
	return "", ctx.CollectError(&Error{
		Message: "bad type",
		Got:     rawString(got.Type().String()),
		Expected: rawString(
			"string (convertible) OR fmt.Stringer OR error"),
	})
}

type tdString struct {
	tdStringBase
}

var _ TestDeep = &tdString{}

// String operator allows to compare a string (or convertible), error
// or fmt.Stringer interface (error interface is tested before
// fmt.Stringer.)
//
//   err := errors.New("error!")
//   CmpDeeply(t, err, String("error!")) // succeeds
//
//   bstr := bytes.NewBufferString("fmt.Stringer!")
//   CmpDeeply(t, bstr, String("fmt.Stringer!")) // succeeds
func String(expected string) TestDeep {
	return &tdString{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdString) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		return err
	}

	if str == s.expected {
		return nil
	}
	if ctx.BooleanError {
		return BooleanError
	}
	return ctx.CollectError(&Error{
		Message:  "does not match",
		Got:      str,
		Expected: s,
	})
}

func (s *tdString) String() string {
	return toString(s.expected)
}

type tdHasPrefix struct {
	tdStringBase
}

var _ TestDeep = &tdHasPrefix{}

// HasPrefix operator allows to compare the prefix of a string (or
// convertible), error or fmt.Stringer interface (error interface is
// tested before fmt.Stringer.)
//
//   type Foobar string
//   CmpDeeply(t, Foobar("foobar"), HasPrefix("foo")) // succeeds
//
//   err := errors.New("error!")
//   CmpDeeply(t, err, HasPrefix("err")) // succeeds
//
//   bstr := bytes.NewBufferString("fmt.Stringer!")
//   CmpDeeply(t, bstr, HasPrefix("fmt")) // succeeds
func HasPrefix(expected string) TestDeep {
	return &tdHasPrefix{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdHasPrefix) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		return err
	}

	if strings.HasPrefix(str, s.expected) {
		return nil
	}
	if ctx.BooleanError {
		return BooleanError
	}
	return ctx.CollectError(&Error{
		Message:  "has not prefix",
		Got:      str,
		Expected: s,
	})
}

func (s *tdHasPrefix) String() string {
	return "HasPrefix(" + toString(s.expected) + ")"
}

type tdHasSuffix struct {
	tdStringBase
}

var _ TestDeep = &tdHasSuffix{}

// HasSuffix operator allows to compare the suffix of a string (or
// convertible), error or fmt.Stringer interface (error interface is
// tested before fmt.Stringer.)
//
//   type Foobar string
//   CmpDeeply(t, Foobar("foobar"), HasSuffix("bar")) // succeeds
//
//   err := errors.New("error!")
//   CmpDeeply(t, err, HasSuffix("!")) // succeeds
//
//   bstr := bytes.NewBufferString("fmt.Stringer!")
//   CmpDeeply(t, bstr, HasSuffix("!")) // succeeds
func HasSuffix(expected string) TestDeep {
	return &tdHasSuffix{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdHasSuffix) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		return err
	}

	if strings.HasSuffix(str, s.expected) {
		return nil
	}
	if ctx.BooleanError {
		return BooleanError
	}
	return ctx.CollectError(&Error{
		Message:  "has not suffix",
		Got:      str,
		Expected: s,
	})
}

func (s *tdHasSuffix) String() string {
	return "HasSuffix(" + toString(s.expected) + ")"
}
