// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build !js,!appengine,!safe,!disableunsafe,!race

package dark

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

var stack = runtime.Stack

// GetFatalizer returns a Fatalizer based on stack trace. It first
// searches for a FatalizerBarrier then for a *testing.T instance. It
// defaults to a FatalPanic instance.
//
// Based on Dave Cheney idea at
// https://github.com/pkg/expect/blob/1fe4c9394a8a0ef4038dc997ea676f44784b7bae/t.go
func GetFatalizer() Fatalizer {
	var buf [8192]byte
	n := stack(buf[:], false)
	sc := bufio.NewScanner(bytes.NewReader(buf[:n]))
	for sc.Scan() {
		var bi int64
		n, _ := fmt.Sscanf(sc.Text(), "github.com/maxatome/go-testdeep/internal/dark.setFatalizerBarrier(%v", &bi)
		if n == 1 {
			return newFatalPanicBarrier(bi)
		}
		var p uintptr
		n, _ = fmt.Sscanf(sc.Text(), "testing.tRunner(%v", &p)
		if n != 1 {
			continue
		}
		return (*testing.T)(unsafe.Pointer(p))
	}
	return FatalPanic("")
}

// fatalPanicBarrier implements a catchable Fatalizer. See
// FatalizerBarrier for details.
type fatalPanicBarrier struct {
	idx int64
	err error
}

func newFatalPanicBarrier(bi int64) Fatalizer {
	barriersMu.Lock()
	defer barriersMu.Unlock()
	barriers[bi] = struct{}{}

	return fatalPanicBarrier{idx: bi}
}

func (p fatalPanicBarrier) Helper() {}
func (p fatalPanicBarrier) Fatal(args ...interface{}) {
	p.err = errors.New(fmt.Sprint(args...))
	panic(p)
}

var barrierIdx int64
var barriersMu sync.Mutex
var barriers = map[int64]struct{}{}

// FatalizerBarrier catch a Fatalizer.Fatal() call occurred in fn, and
// returns its content as an error. Other panics are left as is.
func FatalizerBarrier(fn func()) error {
	return setFatalizerBarrier(atomic.AddInt64(&barrierIdx, 1), fn)
}

func setFatalizerBarrier(key int64, fn func()) (err error) {
	defer func() {
		barriersMu.Lock()
		_, exists := barriers[key]
		if !exists {
			// No panic occurred or not a fatalPanicBarrier panic, so pass
			// without recovering anything
			barriersMu.Unlock()
			return
		}
		delete(barriers, key)
		barriersMu.Unlock()

		if x := recover(); x != nil {
			// Not expected, probably a double-panic
			p, ok := x.(fatalPanicBarrier)
			if !ok {
				panic(x)
			}
			err = p.err
		}
	}()

	fn()
	return
}
