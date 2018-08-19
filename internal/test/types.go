// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package test

import (
	"fmt"
	"testing"
)

// TestingT is a type implementing testdeep.TestingT intended to be
// used in tests.
type TestingT struct {
	LastMessage string
	IsFatal     bool
	HasFailed   bool
}

// Fatal mocks testing.T Error method.
func (t *TestingT) Error(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = false
	t.HasFailed = true
}

// Fatal mocks testing.T Fatal method.
func (t *TestingT) Fatal(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = true
	t.HasFailed = true
}

// Helper mocks testing.T Helper method.
func (t *TestingT) Helper() {
	// Do nothing
}

// TestingFT is a type implementing testdeep.TestingFT intended to be
// used in tests.
type TestingFT struct {
	TestingT
}

// Errorf mocks testing.T Errorf method.
func (t *TestingFT) Errorf(format string, args ...interface{}) {
	t.Error(fmt.Sprintf(format, args...))
}

// Fail mocks testing.T Fail method.
func (t *TestingFT) Fail() {
	t.HasFailed = true
}

// FailNow mocks testing.T FailNow method.
func (t *TestingFT) FailNow() {
	t.HasFailed = true
	t.IsFatal = true
}

// Failed mocks testing.T Failed method.
func (t *TestingFT) Failed() bool {
	return t.HasFailed
}

// Fatalf mocks testing.T Fatalf method.
func (t *TestingFT) Fatalf(format string, args ...interface{}) {
	t.Fatal(fmt.Sprintf(format, args...))
}

// Log mocks testing.T Log method.
func (t *TestingFT) Log(args ...interface{}) {}

// Logf mocks testing.T Logf method.
func (t *TestingFT) Logf(format string, args ...interface{}) {}

// Name mocks testing.T Name method.
func (t *TestingFT) Name() string {
	return ""
}

// Skip mocks testing.T Skip method.
func (t *TestingFT) Skip(args ...interface{}) {}

// SkipNow mocks testing.T SkipNow method.
func (t *TestingFT) SkipNow() {}

// Skipf mocks testing.T Skipf method.
func (t *TestingFT) Skipf(format string, args ...interface{}) {}

// Skipped mocks testing.T Skipped method.
func (t *TestingFT) Skipped() bool {
	return false
}

// Run mocks testing.T Run method.
func (t *TestingFT) Run(name string, f func(t *testing.T)) bool {
	return true
}
