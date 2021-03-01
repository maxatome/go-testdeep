// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package trace_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxatome/go-testdeep/internal/trace"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestIgnorePackage(t *testing.T) {
	const ourPkg = "github.com/maxatome/go-testdeep/internal/trace_test"

	test.IsFalse(t, trace.IsIgnoredPackage(ourPkg))
	test.IsTrue(t, trace.IgnorePackage(1))
	test.IsTrue(t, trace.IsIgnoredPackage(ourPkg))

	trace.ResetIgnoredPackages()
	test.IsFalse(t, trace.IsIgnoredPackage(ourPkg))

	test.IsFalse(t, trace.IgnorePackage(300))
}

func TestFindGoModDir(t *testing.T) {
	tmp, err := ioutil.TempDir("", "go-testdeep")
	if err != nil {
		t.Fatalf("TempDir() failed: %s", err)
	}
	final := filepath.Join(tmp, "a", "b", "c", "d", "e")

	err = os.MkdirAll(final, 0755)
	if err != nil {
		t.Fatalf("MkdirAll(%s) failed: %s", final, err)
	}
	defer os.RemoveAll(final)

	test.EqualStr(t, trace.FindGoModDir(final), "")

	t.Run("/tmp/.../a/b/c/go.mod", func(t *testing.T) {
		goMod := filepath.Join(tmp, "a", "b", "c", "go.mod")

		err := ioutil.WriteFile(goMod, nil, 0644)
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
			err := ioutil.WriteFile(goMod, nil, 0644)
			if err != nil {
				t.Fatalf("WriteFile(%s) failed: %s", goMod, err)
			}
			defer os.Remove(goMod)
		}

		test.EqualStr(t, trace.FindGoModDir(final), "")
	})
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
