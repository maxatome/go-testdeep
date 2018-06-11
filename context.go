// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// ContextConfig allows to configure finely how tests failures are rendered.
//
// See NewT function to use it.
type ContextConfig struct {
	// RootName is the string used to represent the root of got data. It
	// defaults to "DATA". For an HTTP response body, it could be "BODY"
	// for example.
	RootName string
	// MaxErrors is the maximal number of errors to dump in case of Cmp*
	// failure.
	//
	// It defaults to 1 except if the environment variable
	// TESTDEEP_MAX_ERRORS is set. In this latter case, the
	// TESTDEEP_MAX_ERRORS value is converted to an int and used as is.
	//
	// Setting it to 0 has the same effect as 1.
	//
	// Setting it to a negative number means no limit: all errors
	// will be dumped.
	MaxErrors int
}

const contextDefaultRootName = "DATA"

// DefaultContextConfig is the default configuration used to render
// tests failures. If overridden, new settings will impact all Cmp*
// functions and *T methods (if not specifically configured.)
var DefaultContextConfig = ContextConfig{
	RootName: contextDefaultRootName,
	MaxErrors: func() (n int) {
		n, err := strconv.Atoi(os.Getenv("TESTDEEP_MAX_ERRORS"))
		if err != nil || n == 0 {
			n = 1
		}
		return
	}(),
}

func (c *ContextConfig) sanitize() {
	if c.RootName == "" {
		c.RootName = DefaultContextConfig.RootName
	}
	if c.MaxErrors == 0 {
		c.MaxErrors = DefaultContextConfig.MaxErrors
	}
}

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Context is used internally to keep track of the CmpDeeply in-depth
// traversal.
type Context struct {
	path        string
	depth       int
	visited     map[visit]bool
	curOperator TestDeep
	// If true, the contents of the returned *Error will not be
	// checked. Can be used to avoid filling Error{} with expensive
	// computations.
	booleanError bool
	// 0 ≤ maxErrors ≤ 1 stops when first error encoutered, else accumulate
	// maxErrors > 1 stops when maxErrors'th error encoutered
	// < 0 do not stop until comparison ends
	maxErrors int
	errors    *[]*Error
}

// NewContext creates a new Context using DefaultContextConfig configuration.
func NewContext() Context {
	return NewContextWithConfig(DefaultContextConfig)
}

// NewContextWithConfig creates a new Context using a specific configuration.
func NewContextWithConfig(config ContextConfig) (ctx Context) {
	config.sanitize()

	ctx = Context{
		path:      config.RootName,
		visited:   map[visit]bool{},
		maxErrors: config.MaxErrors,
	}

	ctx.initErrors()
	return
}

// NewBooleanContext creates a new boolean Context.
func NewBooleanContext() Context {
	return Context{
		visited:      map[visit]bool{},
		booleanError: true,
	}
}

func (c *Context) initErrors() {
	if c.maxErrors != 0 && c.maxErrors != 1 {
		errors := make([]*Error, 0)
		c.errors = &errors
	}
}

func (c Context) resetErrors() (new Context) {
	new = c
	new.initErrors()
	return
}

// CollectError collects an error in the context. It returns an error
// if the collector is full, nil otherwise.
func (c Context) CollectError(err *Error) *Error {
	if err == nil {
		return nil
	}

	// Error context not initialized yet
	if err.Context.depth == 0 {
		err.Context = c
	}

	if !err.Location.IsInitialized() && c.curOperator != nil {
		err.Location = c.curOperator.GetLocation()
	}

	// Stop when first error encoutered
	if c.errors == nil {
		return err
	}

	// Else, accumulate...
	*c.errors = append(*c.errors, err)
	if c.maxErrors >= 0 && len(*c.errors) >= c.maxErrors {
		return c.mergeErrors()
	}
	return nil
}

func (c Context) mergeErrors() *Error {
	if c.errors == nil || len(*c.errors) == 0 {
		return nil
	}

	if len(*c.errors) > 1 {
		for idx, last := 0, len(*c.errors)-2; idx <= last; idx++ {
			(*c.errors)[idx].Next = (*c.errors)[idx+1]
		}
	}
	return (*c.errors)[0]
}

// AddDepth creates a new Context from current one plus pathAdd.
func (c Context) AddDepth(pathAdd string) (new Context) {
	new = c
	if strings.HasPrefix(new.path, "*") {
		new.path = "(" + new.path + ")" + pathAdd
	} else {
		new.path += pathAdd
	}
	new.depth++
	return
}

// AddArrayIndex creates a new Context from current one plus an array
// dereference for index-th item.
func (c Context) AddArrayIndex(index int) (new Context) {
	return c.AddDepth(fmt.Sprintf("[%d]", index))
}

// AddPtr creates a new Context from current one plus a pointer dereference.
func (c Context) AddPtr(num int) (new Context) {
	new = c
	new.path = strings.Repeat("*", num) + new.path
	new.depth++
	return
}

// AddFunctionCall creates a new Context from current one inside a
// function call.
func (c Context) AddFunctionCall(fn string) (new Context) {
	new = c
	new.path = fn + "(" + new.path + ")"
	new.depth++
	return
}

// ResetPath creates a new Context from current one but reinitializing path.
func (c Context) ResetPath(newPath string) (new Context) {
	new = c
	new.path = newPath
	new.depth++
	return
}

// Path returns the Context path.
func (c Context) Path() string {
	return c.path
}
