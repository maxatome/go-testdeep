// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/location"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestContext(t *testing.T) {
	for _, maxErrors := range []int{0, 1} {
		ctx := ctxerr.Context{
			MaxErrors: maxErrors,
		}

		ctx.InitErrors()
		if ctx.Errors != nil {
			t.Errorf("Errors is non-nil for MaxErrors %d", maxErrors)
		}
	}

	for _, maxErrors := range []int{-1, 2} {
		ctx := ctxerr.Context{
			MaxErrors: maxErrors,
		}

		ctx.InitErrors()
		if ctx.Errors == nil {
			t.Errorf("Errors is nil for MaxErrors %d", maxErrors)
			continue
		}

		*ctx.Errors = append(*ctx.Errors, &ctxerr.Error{})

		newc := ctx.ResetErrors()
		if newc.Errors == nil {
			t.Errorf("after ResetErrors, new Errors is nil for MaxErrors %d",
				maxErrors)
			continue
		}
		if len(*newc.Errors) > 0 {
			t.Errorf("after ResetErrors, new Errors is not empty for MaxErrors %d",
				maxErrors)
		}
		if ctx.Errors == nil {
			t.Errorf("after ResetErrors, old Errors is nil for MaxErrors %d",
				maxErrors)
			continue
		}
	}
}

type MyGetLocationer struct{}

func (g MyGetLocationer) GetLocation() location.Location {
	return location.Location{
		File: "context_test.go",
		Func: "MyFunc",
		Line: 42,
	}
}

func TestContextMergeErrors(t *testing.T) {
	// No errors to merge
	ctx := ctxerr.Context{}
	if ctx.MergeErrors() != nil {
		t.Error("ctx.MergeErrors() returned a *Error")
	}

	errors := []*ctxerr.Error{}
	ctx = ctxerr.Context{
		Errors: &errors,
	}
	if ctx.MergeErrors() != nil {
		t.Error("ctx.MergeErrors() returned a *Error")
	}

	// Only 1 error to merge => itself
	firstErr := &ctxerr.Error{}
	errors = []*ctxerr.Error{firstErr}
	ctx = ctxerr.Context{
		Errors: &errors,
	}
	if ctx.MergeErrors() != firstErr {
		t.Error("ctx.MergeErrors() did not return the only one error")
	}

	// Several errors to merge
	secondErr, thirdErr := &ctxerr.Error{}, &ctxerr.Error{}
	errors = []*ctxerr.Error{firstErr, secondErr, thirdErr}
	ctx = ctxerr.Context{
		Errors: &errors,
	}
	if ctx.MergeErrors() != firstErr {
		t.Error("ctx.MergeErrors() did not return the first error")
		return
	}
	if firstErr.Next != secondErr {
		t.Error("ctx.MergeErrors() second error is not linked to first one")
		return
	}
	if secondErr.Next != thirdErr {
		t.Error("ctx.MergeErrors() third error is not linked to second one")
		return
	}
	if thirdErr.Next != nil {
		t.Error("ctx.MergeErrors() third error has a non-nil Next!")
	}
}

func TestContextCollectError(t *testing.T) {
	//
	// Only one error kept
	ctx := ctxerr.Context{}

	if ctx.CollectError(nil) != nil {
		t.Error("ctx.CollectError(nil) returned non-nil *Error")
	}

	err := ctxerr.Context{BooleanError: true}.CollectError(&ctxerr.Error{})
	if err != ctxerr.BooleanError {
		t.Error("boolean-ctx.CollectError(X) did not return BooleanError")
	}

	// !err.Location.IsInitialized() + ctx.CurOperator == nil
	origErr := &ctxerr.Error{}
	err = ctx.CollectError(origErr)
	if err != origErr {
		t.Error("ctx.CollectError(err) != err")
	}

	// !err.Location.IsInitialized() + ctx.CurOperator != nil
	ctx.CurOperator = MyGetLocationer{}
	origErr = &ctxerr.Error{}
	err = ctx.CollectError(origErr)
	if err != origErr {
		t.Error("ctx.CollectError(err) != err")
	}
	test.EqualInt(t, err.Location.Line, 42, // see MyGetLocationer.GetLocation()
		"ctx.CollectError(err) initialized err.Location")

	// err.Location.IsInitialized()
	origErr = &ctxerr.Error{
		Location: location.Location{
			File: "zz.go",
			Func: "ErrFunc",
			Line: 24,
		},
	}
	err = ctx.CollectError(origErr)
	if err != origErr {
		t.Error("ctx.CollectError(err) != err")
	}
	test.EqualInt(t, err.Location.Line, 24,
		"ctx.CollectError(err) did not touch err.Location")

	//
	// 2 errors kept max
	errors := []*ctxerr.Error{}
	ctx = ctxerr.Context{
		Errors:    &errors,
		MaxErrors: 2,
	}
	origErr = &ctxerr.Error{}
	if ctx.CollectError(origErr) != nil { // 1st error is accumulated
		t.Error("ctx.CollectError(err) != nil")
		return
	}

	secondErr := &ctxerr.Error{}
	if ctx.CollectError(secondErr) != origErr {
		t.Error("ctx.CollectError(err) != origErr")
		return
	}

	if origErr.Next != secondErr {
		t.Error("origErr.Next != secondErr")
		return
	}

	if secondErr.Next != ctxerr.ErrTooManyErrors {
		t.Error("secondErr.Next != ErrTooManyErrors")
		return
	}

	//
	// All errors kept
	errors = nil
	ctx = ctxerr.Context{
		Errors:    &errors,
		MaxErrors: -1,
	}
	for i := 0; i < 100; i++ {
		if ctx.CollectError(&ctxerr.Error{}) != nil { // 1st error is accumulated
			t.Errorf("#%d: ctx.CollectError(err) != nil", i)
			return
		}
	}
	if len(errors) != 100 {
		t.Errorf("Only %d errors accumulated instead of 100", len(errors))
	}

	//
	// Do not collect 2 times the same error
	errors = nil
	ctx = ctxerr.Context{
		Errors:    &errors,
		MaxErrors: -1,
	}
	ctx.CollectError(&ctxerr.Error{}) //nolint: errcheck
	x := &ctxerr.Error{}
	ctx.CollectError(x)               //nolint: errcheck
	ctx.CollectError(&ctxerr.Error{}) //nolint: errcheck
	ctx.CollectError(x)               //nolint: errcheck
	ctx.CollectError(x)               //nolint: errcheck
	ctx.CollectError(&ctxerr.Error{}) //nolint: errcheck
	if len(errors) != 4 {
		t.Errorf("%d errors accumulated instead of 4", len(errors))
	}
}

func TestCannotCompareError(t *testing.T) {
	ctx := ctxerr.Context{BooleanError: true}

	err := ctx.CannotCompareError()
	if err != ctxerr.BooleanError {
		t.Error("CannotCompareError does not return ctxerr.BooleanError")
	}

	ctx = ctxerr.Context{}
	err = ctx.CannotCompareError()
	test.EqualStr(t, err.Message, "cannot compare")
}

func TestContextPath(t *testing.T) {
	ctx := ctxerr.Context{Path: ctxerr.NewPath("DATA")}
	ctx = ctx.AddField("field")
	test.EqualStr(t, ctx.Path.String(), "DATA.field")
	test.EqualInt(t, ctx.Depth, 1)

	ctx = ctx.AddPtr(2)
	test.EqualStr(t, ctx.Path.String(), "**DATA.field")
	test.EqualInt(t, ctx.Depth, 2)

	ctx = ctx.AddField("another")
	test.EqualStr(t, ctx.Path.String(), "(*DATA.field).another")
	test.EqualInt(t, ctx.Depth, 3)

	ctx = ctx.AddCustomLevel("→cust")
	test.EqualStr(t, ctx.Path.String(), "(*DATA.field).another→cust")
	test.EqualInt(t, ctx.Depth, 4)

	ctx = ctxerr.Context{Path: ctxerr.NewPath("DATA")}
	ctx = ctx.AddArrayIndex(18)
	test.EqualStr(t, ctx.Path.String(), "DATA[18]")
	test.EqualInt(t, ctx.Depth, 1)

	ctx = ctxerr.Context{Path: ctxerr.NewPath("DATA")}
	ctx = ctx.AddMapKey("foo")
	test.EqualStr(t, ctx.Path.String(), `DATA["foo"]`) // special case of util.ToString()
	test.EqualInt(t, ctx.Depth, 1)

	ctx = ctxerr.Context{Path: ctxerr.NewPath("DATA")}
	ctx = ctx.AddMapKey(12)
	test.EqualStr(t, ctx.Path.String(), `DATA[12]`)
	test.EqualInt(t, ctx.Depth, 1)

	ctx = ctxerr.Context{Path: ctxerr.NewPath("DATA")}
	ctx = ctx.AddFunctionCall("foobar")
	test.EqualStr(t, ctx.Path.String(), "foobar(DATA)")
	test.EqualInt(t, ctx.Depth, 1)

	ctx = ctx.ResetPath("NEW")
	test.EqualStr(t, ctx.Path.String(), "NEW")
	test.EqualInt(t, ctx.Depth, 2)
}
