// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build !js,!appengine,!safe,!disableunsafe,!race

package dark

import (
	"runtime"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestGetFatalizer(t *testing.T) {
	f := GetFatalizer()
	tt, ok := f.(*testing.T)
	if test.IsTrue(t, ok) {
		test.IsTrue(t, t == tt)
	}

	defer func() { stack = runtime.Stack }()
	stack = func(buf []byte, all bool) int {
		return copy(buf, `goroutine 21 [running]:
github.com/maxatome/go-testdeep/internal/dark.GetFatalizer(0x469c9b, 0x202e3258719)
	/home/max/Projet/go/src/github.com/maxatome/go-testdeep/internal/dark/fatalizer_unsafe.go:33 +0x76
github.com/maxatome/go-testdeep/internal/dark_test.TestFatalizer(0xc000082d80)
	/home/max/Projet/go/src/github.com/maxatome/go-testdeep/internal/dark/fatalizer_test.go:17 +0x26
testing.ZZZ(0xc000082d80, 0x589550)
	/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1238 +0x2b3
`)
	}
	f = GetFatalizer()
	_, ok = f.(FatalPanic)
	test.IsTrue(t, ok)
}
