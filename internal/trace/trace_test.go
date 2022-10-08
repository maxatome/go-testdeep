// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package trace_test

import (
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/trace"
)

func TestIgnorePackage(t *testing.T) {
	const ourPkg = "github.com/maxatome/go-testdeep/internal/trace_test"

	trace.Reset()

	test.IsFalse(t, trace.IsIgnoredPackage(ourPkg))
	test.IsTrue(t, trace.IgnorePackage())
	test.IsTrue(t, trace.IsIgnoredPackage(ourPkg))

	test.IsTrue(t, trace.UnignorePackage())
	test.IsFalse(t, trace.IsIgnoredPackage(ourPkg))

	test.IsTrue(t, trace.IgnorePackage())
	test.IsTrue(t, trace.IsIgnoredPackage(ourPkg))

	test.IsFalse(t, trace.IgnorePackage(300))
	test.IsFalse(t, trace.UnignorePackage(300))
}

func TestFindGoModDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "go-testdeep")
	if err != nil {
		t.Fatalf("TempDir() failed: %s", err)
	}
	final := filepath.Join(tmp, "a", "b", "c", "d", "e")

	err = os.MkdirAll(final, 0755)
	if err != nil {
		t.Fatalf("MkdirAll(%s) failed: %s", final, err)
	}
	defer os.RemoveAll(tmp)

	test.EqualStr(t, trace.FindGoModDir(final), "")

	t.Run("/tmp/.../a/b/c/go.mod", func(t *testing.T) {
		goMod := filepath.Join(tmp, "a", "b", "c", "go.mod")

		err := os.WriteFile(goMod, nil, 0644)
		if err != nil {
			t.Fatalf("WriteFile(%s) failed: %s", goMod, err)
		}
		defer os.Remove(goMod)

		test.EqualStr(t,
			trace.FindGoModDir(final),
			filepath.Join(tmp, "a", "b", "c")+string(filepath.Separator),
		)
	})

	t.Run("/tmp/go.mod", func(t *testing.T) {
		goMod := filepath.Join(os.TempDir(), "go.mod")

		if _, err := os.Stat(goMod); err != nil {
			if !os.IsNotExist(err) {
				t.Fatalf("Stat(%s) failed: %s", goMod, err)
			}
			err := os.WriteFile(goMod, nil, 0644)
			if err != nil {
				t.Fatalf("WriteFile(%s) failed: %s", goMod, err)
			}
			defer os.Remove(goMod)
		}

		test.EqualStr(t, trace.FindGoModDir(final), "")
	})
}

func TestFindGoModDirLinks(t *testing.T) {
	tmp, err := os.MkdirTemp("", "go-testdeep")
	if err != nil {
		t.Fatalf("TempDir() failed: %s", err)
	}

	goModDir := filepath.Join(tmp, "a", "b", "c")
	truePath := filepath.Join(goModDir, "d", "e")
	linkPath := filepath.Join(tmp, "a", "b", "e")

	err = os.MkdirAll(truePath, 0755)
	if err != nil {
		t.Fatalf("MkdirAll(%s) failed: %s", truePath, err)
	}
	defer os.RemoveAll(tmp)

	err = os.Symlink(truePath, linkPath)
	if err != nil {
		t.Fatalf("Symlink(%s, %s) failed: %s", truePath, linkPath, err)
	}

	goMod := filepath.Join(goModDir, "go.mod")

	err = os.WriteFile(goMod, nil, 0644)
	if err != nil {
		t.Fatalf("WriteFile(%s) failed: %s", goMod, err)
	}
	defer os.Remove(goMod)

	goModDir += string(filepath.Separator)

	// Simple FindGoModDir
	test.EqualStr(t, trace.FindGoModDir(truePath), goModDir)
	test.EqualStr(t, trace.FindGoModDir(linkPath), "") // not found

	// FindGoModDirLinks
	test.EqualStr(t, trace.FindGoModDirLinks(truePath), goModDir)
	test.EqualStr(t, trace.FindGoModDirLinks(linkPath), goModDir)

	test.EqualStr(t, trace.FindGoModDirLinks(tmp), "")
}

func TestSplitPackageFunc(t *testing.T) {
	pkg, fn := trace.SplitPackageFunc("testing.Fatal")
	test.EqualStr(t, pkg, "testing")
	test.EqualStr(t, fn, "Fatal")

	pkg, fn = trace.SplitPackageFunc("github.com/maxatome/go-testdeep/td.Cmp")
	test.EqualStr(t, pkg, "github.com/maxatome/go-testdeep/td")
	test.EqualStr(t, fn, "Cmp")

	pkg, fn = trace.SplitPackageFunc("foo/bar/test.(*T).Cmp")
	test.EqualStr(t, pkg, "foo/bar/test")
	test.EqualStr(t, fn, "(*T).Cmp")

	pkg, fn = trace.SplitPackageFunc("foo/bar/test.(*X).c.func1")
	test.EqualStr(t, pkg, "foo/bar/test")
	test.EqualStr(t, fn, "(*X).c.func1")

	pkg, fn = trace.SplitPackageFunc("foo/bar/test.(*X).c.func1")
	test.EqualStr(t, pkg, "foo/bar/test")
	test.EqualStr(t, fn, "(*X).c.func1")

	pkg, fn = trace.SplitPackageFunc("foobar")
	test.EqualStr(t, pkg, "")
	test.EqualStr(t, fn, "foobar")

	pkg, fn = trace.SplitPackageFunc("")
	test.EqualStr(t, pkg, "")
	test.EqualStr(t, fn, "")
}

func d(end string) []trace.Level { return trace.Retrieve(0, end) }
func c(end string) []trace.Level { return d(end) }
func b(end string) []trace.Level { return c(end) }
func a(end string) []trace.Level { return b(end) }

func TestZRetrieve(t *testing.T) {
	trace.Reset()

	levels := a("testing.tRunner")
	if !test.EqualInt(t, len(levels), 5) ||
		!test.EqualStr(t, levels[0].Func, "d") ||
		!test.EqualStr(t, levels[0].Package, "github.com/maxatome/go-testdeep/internal/trace_test") ||
		!test.EqualStr(t, levels[1].Func, "c") ||
		!test.EqualStr(t, levels[1].Package, "github.com/maxatome/go-testdeep/internal/trace_test") ||
		!test.EqualStr(t, levels[2].Func, "b") ||
		!test.EqualStr(t, levels[2].Package, "github.com/maxatome/go-testdeep/internal/trace_test") ||
		!test.EqualStr(t, levels[3].Func, "a") ||
		!test.EqualStr(t, levels[3].Package, "github.com/maxatome/go-testdeep/internal/trace_test") ||
		!test.EqualStr(t, levels[4].Func, "TestZRetrieve") ||
		!test.EqualStr(t, levels[4].Package, "github.com/maxatome/go-testdeep/internal/trace_test") {
		t.Errorf("%#v", levels)
	}

	levels = trace.Retrieve(0, "unknown.unknown")
	maxLevels := len(levels)
	test.IsTrue(t, maxLevels > 2)
	test.EqualStr(t, levels[len(levels)-1].Func, "goexit") // runtime.goexit

	for i := range levels {
		test.IsTrue(t, trace.IgnorePackage(i))
	}
	levels = trace.Retrieve(0, "unknown.unknown")
	test.EqualInt(t, len(levels), 0)

	// Init GOPATH filter
	trace.Reset()
	trace.Init()

	test.IsTrue(t, trace.IgnorePackage())
	levels = trace.Retrieve(0, "unknown.unknown")
	test.EqualInt(t, len(levels), maxLevels-1)
}

type FakeFrames struct {
	frames []runtime.Frame
	cur    int
}

func (f *FakeFrames) Next() (runtime.Frame, bool) {
	if f.cur >= len(f.frames) {
		return runtime.Frame{}, false
	}
	f.cur++
	return f.frames[f.cur-1], f.cur < len(f.frames)
}

func TestZRetrieveFake(t *testing.T) {
	saveCallersFrames, saveGOPATH := trace.CallersFrames, build.Default.GOPATH
	defer func() {
		trace.CallersFrames, build.Default.GOPATH = saveCallersFrames, saveGOPATH
	}()

	var fakeFrames FakeFrames
	trace.CallersFrames = func(_ []uintptr) trace.Frames { return &fakeFrames }
	build.Default.GOPATH = "/foo/bar"

	trace.Reset()
	trace.Init()

	fakeFrames = FakeFrames{
		frames: []runtime.Frame{
			{},
			{Function: "", File: "/foo/bar/src/zip/zip.go", Line: 23},
			{Function: "", File: "/foo/bar/pkg/mod/zzz/zzz.go", Line: 42},
			{Function: "", File: "/bar/foo.go", Line: 34},
			{Function: "pkg.MyFunc"},
			{},
		},
	}
	levels := trace.Retrieve(0, "pipo")
	if test.EqualInt(t, len(levels), 4) {
		test.EqualStr(t, levels[0].Func, "<unknown function>")
		test.EqualStr(t, levels[0].Package, "")
		test.EqualStr(t, levels[0].FileLine, "zip/zip.go:23")

		test.EqualStr(t, levels[1].Func, "<unknown function>")
		test.EqualStr(t, levels[1].Package, "")
		test.EqualStr(t, levels[1].FileLine, "zzz/zzz.go:42")

		test.EqualStr(t, levels[2].Func, "<unknown function>")
		test.EqualStr(t, levels[2].Package, "")
		test.EqualStr(t, levels[2].FileLine, "/bar/foo.go:34")

		test.EqualStr(t, levels[3].Func, "MyFunc")
		test.EqualStr(t, levels[3].Package, "pkg")
		test.EqualStr(t, levels[3].FileLine, "")
	} else {
		t.Errorf("%#v", levels)
	}
}
