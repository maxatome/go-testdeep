// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/trace"
)

// TestingT is a type implementing td.TestingT intended to be used in
// tests.
type TestingT struct {
	Messages  []string
	IsFatal   bool
	HasFailed bool
}

type testingFatal string

// NewTestingT returns a new instance of [*TestingT].
func NewTestingT() *TestingT {
	return &TestingT{}
}

// Error mocks [testing.T] Error method.
func (t *TestingT) Error(args ...any) {
	t.Messages = append(t.Messages, fmt.Sprint(args...))
	t.IsFatal = false
	t.HasFailed = true
}

// Fatal mocks [testing.T.Fatal] method.
func (t *TestingT) Fatal(args ...any) {
	t.Messages = append(t.Messages, fmt.Sprint(args...))
	t.IsFatal = true
	t.HasFailed = true

	panic(testingFatal(t.Messages[len(t.Messages)-1]))
}

func (t *TestingT) CatchFatal(fn func()) (fatalStr string) {
	panicked := true
	trace.IgnorePackage()
	defer func() {
		trace.UnignorePackage()
		if panicked {
			if x := recover(); x != nil {
				if str, ok := x.(testingFatal); ok {
					fatalStr = string(str)
				} else {
					panic(x) // rethrow
				}
			}
		}
	}()

	fn()
	panicked = false
	return
}

// ContainsMessages checks expectedMsgs are all present in Messages, in
// this order. It stops when a message is not found and returns the
// remaining messages.
func (t *TestingT) ContainsMessages(expectedMsgs ...string) []string {
	curExp := 0
	for _, msg := range t.Messages {
		for {
			if curExp == len(expectedMsgs) {
				return nil
			}
			pos := strings.Index(msg, expectedMsgs[curExp])
			if pos < 0 {
				break
			}
			msg = msg[pos+len(expectedMsgs[curExp]):]
			curExp++
		}
	}
	return expectedMsgs[curExp:]
}

// Helper mocks [testing.T.Helper] method.
func (t *TestingT) Helper() {
	// Do nothing
}

// LastMessage returns the last message.
func (t *TestingT) LastMessage() string {
	if len(t.Messages) == 0 {
		return ""
	}
	return t.Messages[len(t.Messages)-1]
}

// ResetMessages resets the messages.
func (t *TestingT) ResetMessages() {
	t.Messages = t.Messages[:0]
}

// TestingTB is a type implementing [testing.TB] intended to be used in
// tests.
type TestingTB struct {
	TestingT
	name string
	testing.TB
	cleanup func()
}

// NewTestingTB returns a new instance of [*TestingTB].
func NewTestingTB(name string) *TestingTB {
	return &TestingTB{name: name}
}

// Cleanup mocks [testing.T.Cleanup] method. Not thread-safe but we
// don't care in tests.
func (t *TestingTB) Cleanup(fn func()) {
	old := t.cleanup
	t.cleanup = func() {
		if old != nil {
			defer old()
		}
		fn()
	}
	runtime.SetFinalizer(t, func(t *TestingTB) { t.cleanup() })
}

// Fatal mocks [testing.T.Error] method.
func (t *TestingTB) Error(args ...any) {
	t.TestingT.Error(args...)
}

// Errorf mocks [testing.T.Errorf] method.
func (t *TestingTB) Errorf(format string, args ...any) {
	t.TestingT.Error(fmt.Sprintf(format, args...))
}

// Fail mocks [testing.T.Fail] method.
func (t *TestingTB) Fail() {
	t.HasFailed = true
}

// FailNow mocks [testing.T.FailNow] method.
func (t *TestingTB) FailNow() {
	t.HasFailed = true
	t.IsFatal = true
}

// Failed mocks [testing.T.Failed] method.
func (t *TestingTB) Failed() bool {
	return t.HasFailed
}

// Fatal mocks [testing.T.Fatal] method.
func (t *TestingTB) Fatal(args ...any) {
	t.TestingT.Fatal(args...)
}

// Fatalf mocks [testing.T.Fatalf] method.
func (t *TestingTB) Fatalf(format string, args ...any) {
	t.TestingT.Fatal(fmt.Sprintf(format, args...))
}

// Helper mocks [testing.T.Helper] method.
func (t *TestingTB) Helper() {
	// Do nothing
}

// Log mocks [testing.T.Log] method.
func (t *TestingTB) Log(args ...any) {
	t.Messages = append(t.Messages, fmt.Sprint(args...))
}

// Logf mocks [testing.T.Logf] method.
func (t *TestingTB) Logf(format string, args ...any) {
	t.Log(fmt.Sprintf(format, args...))
}

// Name mocks [testing.T.Name] method.
func (t *TestingTB) Name() string {
	return t.name
}

// Skip mocks [testing.T.Skip] method.
func (t *TestingTB) Skip(args ...any) {}

// SkipNow mocks [testing.T.SkipNow] method.
func (t *TestingTB) SkipNow() {}

// Skipf mocks [testing.T.Skipf] method.
func (t *TestingTB) Skipf(format string, args ...any) {}

// Skipped mocks [testing.T.Skipped] method.
func (t *TestingTB) Skipped() bool {
	return false
}

// ParallelTestingTB is a type implementing [testing.TB] and a
// Parallel() method intended to be used in tests.
type ParallelTestingTB struct {
	IsParallel bool
	*TestingTB
}

// NewParallelTestingTB returns a new instance of [*ParallelTestingTB].
func NewParallelTestingTB(name string) *ParallelTestingTB {
	return &ParallelTestingTB{TestingTB: NewTestingTB(name)}
}

// Parallel mocks the [testing.T.Parallel] method. Not thread-safe.
func (t *ParallelTestingTB) Parallel() {
	if t.IsParallel {
		// testing.T.Parallel() panics if called multiple times for the
		// same test.
		panic("testing: t.Parallel called multiple times")
	}
	t.IsParallel = true
}
