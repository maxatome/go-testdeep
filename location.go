// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"runtime"
	"strings"
)

// Location record a place in a source file.
type Location struct {
	File string // File name
	Func string // Function name
	Line int    // Line number inside file
}

// NewLocation returns a new Location. "callDepth" is the number of
// stack frames to ascend to get the calling function (Func field),
// added to 1 to get the File & Line fields.
//
// If the location can not be determined, false is returned and
// the Location is not valid.
func NewLocation(callDepth int) (loc Location, ok bool) {
	_, loc.File, loc.Line, ok = runtime.Caller(callDepth + 1)
	if !ok {
		return
	}

	if index := strings.LastIndexAny(loc.File, `/\`); index >= 0 {
		loc.File = loc.File[index+1:]
	}

	pc, _, _, ok := runtime.Caller(callDepth)
	if !ok {
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn != nil {
		loc.Func = fn.Name()
	} else {
		ok = false
	}
	return
}

// IsInitialized returns true if the Location is initialized
// (eg. NewLocation() called without an error), false otherwise.
func (l Location) IsInitialized() bool {
	return l.File != ""
}

// Implements fmt.Stringer.
func (l Location) String() string {
	return fmt.Sprintf("%s at %s:%d", l.Func, l.File, l.Line)
}
