// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package trace

import (
	"fmt"
	"io"
	"strings"
)

// Level is a level when retrieving a stack trace.
type Level struct {
	Package  string
	Func     string
	FileLine string
}

// Stack is a simple stack trace.
type Stack []Level

// Match returns true if the ith level of s matches pkg (if not empty)
// and any function in anyFunc.
//
// If anyFunc is empty, only the package is tested.
//
// If a function in anyFunc ends with "*", only the prefix is checked.
func (s Stack) Match(i int, pkg string, anyFunc ...string) bool {
	if i < 0 {
		i = len(s) + i
	}
	if i < 0 || i >= len(s) {
		return false
	}

	level := s[i]

	if pkg != "" && level.Package != pkg {
		return false
	}

	if len(anyFunc) == 0 {
		return true
	}

	for _, fn := range anyFunc {
		if strings.HasSuffix(fn, "*") {
			if strings.HasPrefix(level.Func, fn[:len(fn)-1]) {
				return true
			}
		} else if level.Func == fn {
			return true
		}
	}
	return false
}

// IsRelevant returns true if the stack contains more than one level,
// or if the single level has a path with at least one directory.
func (s Stack) IsRelevant() bool {
	return len(s) > 1 || (len(s) > 0 && strings.ContainsAny(s[0].FileLine, `/\`))
}

// Dump writes the stack to "w".
func (s Stack) Dump(w io.Writer) {
	fnMaxLen := 0
	for _, level := range s {
		if len(level.Func) > fnMaxLen {
			fnMaxLen = len(level.Func)
		}
	}
	fnMaxLen += 2

	nl := ""
	for _, level := range s {
		fmt.Fprintf(w, "%s\t%-*s %s", nl, fnMaxLen, level.Func+"()", level.FileLine)
		nl = "\n"
	}
}
