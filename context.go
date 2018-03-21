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

type Context struct {
	Path    string
	Depth   int
	visited map[visit]bool
	// If true, got can contain TestDeep values. Used internally.
	expectVsExpect bool
	// If true, the content of the returned *Error will not be
	// checked. Can be used to avoid filling Error{} with expensive
	// computations.
	booleanError bool
}

func NewContext(path string) Context {
	return Context{
		Path:    path,
		visited: map[visit]bool{},
	}
}
func NewBooleanContext() Context {
	return Context{
		visited:      map[visit]bool{},
		booleanError: true,
	}
}

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
func (c Context) AddArrayIndex(index int) (new Context) {
	return c.AddDepth(fmt.Sprintf("[%d]", index))
}

func (c Context) AddPtr(num int) (new Context) {
	new = c
	new.Path = strings.Repeat("*", num) + new.Path
	new.Depth++
	return
}
