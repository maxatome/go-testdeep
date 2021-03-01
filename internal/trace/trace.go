// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package trace

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ignorePkg = map[string]struct{}{}
	goPaths   []string
	goModDir  string
)

type Level struct {
	Func     string
	FileLine string
}

func getPackage(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			pkg, _ := SplitPackageFunc(fn.Name())
			return pkg
		}
	}
	return ""
}

// IgnorePackage record the calling package as ignored one in trace.
func IgnorePackage(skip int) bool {
	if pkg := getPackage(skip); pkg != "" {
		ignorePkg[pkg] = struct{}{}
		return true
	}
	return false
}

// UnignorePackage cancels a previous use of IgnorePackage, so the
// calling package is no longer ignored. Only intended to be used in
// go-testdeep internal tests.
func UnignorePackage(skip int) bool {
	if pkg := getPackage(skip); pkg != "" {
		delete(ignorePkg, pkg)
		return true
	}
	return false
}

// ResetIgnoredPackages reset the ignored packages map. Only intended
// to be used in go-testdeep internal tests.
func ResetIgnoredPackages() {
	ignorePkg = map[string]struct{}{}
}

// IsIgnoredPackage returns true if pkg is ignored, false
// otherwise. Only intended to be used in go-testdeep internal tests.
func IsIgnoredPackage(pkg string) (ok bool) {
	_, ok = ignorePkg[pkg]
	return
}

// FindGoModDir finds the closest directory containing go.mod file
// starting from directory in.
func FindGoModDir(in string) string {
	for {
		_, err := os.Stat(filepath.Join(in, "go.mod"))
		if err == nil {
			// Do not accept /tmp/go.mod
			if in != os.TempDir() {
				return in + string(filepath.Separator)
			}
			return ""
		}

		nd := filepath.Dir(in)
		if nd == in {
			return ""
		}
		in = nd
	}
}

// Init initializes trace global variables.
func Init() {
	// GOPATH directories
	goPaths = nil
	for _, dir := range filepath.SplitList(build.Default.GOPATH) {
		dir = filepath.Clean(dir)
		goPaths = append(goPaths,
			filepath.Join(dir, "pkg", "mod")+string(filepath.Separator),
			filepath.Join(dir, "src")+string(filepath.Separator),
		)
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// go.mod directory
	goModDir = FindGoModDir(filepath.Clean(wd))
	if goModDir == "" {
		lwd, err := filepath.EvalSymlinks(wd)
		if err == nil && lwd != wd {
			goModDir = FindGoModDir(filepath.Clean(lwd))
		}
	}
}

// Retrieve retrieves a trace and returns it.
func Retrieve(skip int, endFunction string) []Level {
	var trace []Level
	var pc [40]uintptr
	if num := runtime.Callers(skip, pc[:]); num > 0 {
		checkIgnore := true
		frames := runtime.CallersFrames(pc[:num])
		for {
			frame, more := frames.Next()

			fn := frame.Function
			if fn == endFunction {
				break
			}

			var pkg string
			if fn == "" {
				fn = "<unknown function>"
			} else {
				pkg, fn = SplitPackageFunc(fn)
				if checkIgnore && IsIgnoredPackage(pkg) {
					if more {
						continue
					}
					break
				}
				checkIgnore = false
			}

			file := strings.TrimPrefix(frame.File, goModDir)
			if file == frame.File {
				for _, dir := range goPaths {
					file = strings.TrimPrefix(frame.File, dir)
					if file != frame.File {
						break
					}
				}

				if file == frame.File {
					file = strings.TrimPrefix(frame.File, build.Default.GOROOT)
					if file != frame.File {
						file = filepath.Join("$GOROOT", file)
					}
				}
			}

			trace = append(trace, Level{
				Func:     fn,
				FileLine: fmt.Sprintf("%s:%d", file, frame.Line),
			})
			if !more {
				break
			}
		}
	}
	return trace
}

// SplitPackageFunc split a fully qualified function name into its
// package and function parts:
//   "foo/bar/test.fn" → "foo/bar/test", "fn"
//   "foo/bar/test.X.fn" → "foo/bar/test", "X.fn"
//   "foo/bar/test.(*X).fn" → "foo/bar/test", "(*X).fn"
//   "foo/bar/test.(*X).fn.func1" → "foo/bar/test", "(*X).fn.func1"
//   "weird" → "", "weird"
func SplitPackageFunc(fn string) (string, string) {
	sp := strings.LastIndexByte(fn, '/')
	if sp < 0 {
		sp = 0 // std package
	}

	dp := strings.IndexByte(fn[sp:], '.')
	if dp < 0 {
		return "", fn
	}

	return fn[:sp+dp], fn[sp+dp+1:]
}
