// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/maxatome/go-testdeep/internal/location"
	"github.com/maxatome/go-testdeep/internal/str"
)

type Visit struct {
	A1  unsafe.Pointer
	A2  unsafe.Pointer
	Typ reflect.Type
}

// Context is used internally to keep track of the CmpDeeply in-Depth
// traversal.
type Context struct {
	Path        string
	Depth       int
	Visited     map[Visit]bool
	CurOperator location.GetLocationer
	// If true, the contents of the returned *Error will not be
	// checked. Can be used to avoid filling Error{} with expensive
	// computations.
	BooleanError bool
	// 0 ≤ MaxErrors ≤ 1 stops when first error encoutered (without the
	// "Too many errors" error);
	// MaxErrors > 1 stops when MaxErrors'th error encoutered (with a
	// last "Too many errors" error);
	// < 0 do not stop until comparison ends.
	MaxErrors int
	Errors    *[]*Error
	// See ContexConfig.FailureIsFatal for details
	FailureIsFatal bool
}

func (c *Context) InitErrors() {
	if c.MaxErrors != 0 && c.MaxErrors != 1 {
		var errors []*Error
		c.Errors = &errors
	}
}

func (c Context) ResetErrors() (new Context) {
	new = c
	new.InitErrors()
	return
}

// CollectError collects an error in the context. It returns an error
// if the collector is full, nil otherwise.
func (c Context) CollectError(err *Error) *Error {
	if err == nil {
		return nil
	}

	// Error context not initialized yet
	if err.Context.Depth == 0 {
		err.Context = c
	}

	if !err.Location.IsInitialized() && c.CurOperator != nil {
		err.Location = c.CurOperator.GetLocation()
	}

	// Stop when first error encoutered
	if c.Errors == nil {
		return err
	}

	// Else, accumulate...
	*c.Errors = append(*c.Errors, err)
	if c.MaxErrors >= 0 && len(*c.Errors) >= c.MaxErrors {
		*c.Errors = append(*c.Errors, ErrTooManyErrors)
		return c.MergeErrors()
	}
	return nil
}

func (c Context) MergeErrors() *Error {
	if c.Errors == nil || len(*c.Errors) == 0 {
		return nil
	}

	if len(*c.Errors) > 1 {
		for idx, last := 0, len(*c.Errors)-2; idx <= last; idx++ {
			(*c.Errors)[idx].Next = (*c.Errors)[idx+1]
		}
	}
	return (*c.Errors)[0]
}

// AddDepth creates a new Context from current one plus pathAdd.
func (c Context) AddDepth(pathAdd string) (new Context) {
	new = c
	if strings.HasPrefix(new.Path, "*") {
		new.Path = "(" + new.Path + ")" + pathAdd
	} else {
		new.Path += pathAdd
	}
	new.Depth++
	return
}

// AddArrayIndex creates a new Context from current one plus an array
// dereference for index-th item.
func (c Context) AddArrayIndex(index int) Context {
	return c.AddDepth(fmt.Sprintf("[%d]", index))
}

// AddMapKey creates a new Context from current one plus a map
// dereference for key key.
func (c Context) AddMapKey(key interface{}) Context {
	return c.AddDepth("[" + str.ToString(key) + "]")
}

// AddPtr creates a new Context from current one plus a pointer dereference.
func (c Context) AddPtr(num int) (new Context) {
	new = c
	new.Path = strings.Repeat("*", num) + new.Path
	new.Depth++
	return
}

// AddFunctionCall creates a new Context from current one inside a
// function call.
func (c Context) AddFunctionCall(fn string) (new Context) {
	new = c
	new.Path = fn + "(" + new.Path + ")"
	new.Depth++
	return
}

// ResetPath creates a new Context from current one but reinitializing Path.
func (c Context) ResetPath(newPath string) (new Context) {
	new = c
	new.Path = newPath
	new.Depth++
	return
}
