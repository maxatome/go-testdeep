package testdeep

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Context is used internally to keep track of the CmpDeeply in-depth
// traversal.
type Context struct {
	path    string
	depth   int
	visited map[visit]bool
	// If true, the content of the returned *Error will not be
	// checked. Can be used to avoid filling Error{} with expensive
	// computations.
	booleanError bool
}

// NewContext creates a new Context using path.
func NewContext(path string) Context {
	return Context{
		path:    path,
		visited: map[visit]bool{},
	}
}

// NewBooleanContext creates a new boolean Context.
func NewBooleanContext() Context {
	return Context{
		visited:      map[visit]bool{},
		booleanError: true,
	}
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

// Path returns the Context path.
func (c Context) Path() string {
	return c.path
}
