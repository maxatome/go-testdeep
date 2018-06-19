// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"testing"
)

type TestTestingT struct {
	LastMessage string
	IsFatal     bool
	HasFailed   bool
}

func (t *TestTestingT) Error(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = false
	t.HasFailed = true
}

func (t *TestTestingT) Fatal(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = true
	t.HasFailed = true
}

func (t *TestTestingT) Helper() {
	// Do nothing
}

type TestTestingFT struct {
	TestTestingT
}

func (t *TestTestingFT) Errorf(format string, args ...interface{}) {
	t.Error(fmt.Sprintf(format, args...))
}

func (t *TestTestingFT) Fail() {
	t.HasFailed = true
}

func (t *TestTestingFT) FailNow() {
	t.HasFailed = true
	t.IsFatal = true
}

func (t *TestTestingFT) Failed() bool {
	return t.HasFailed
}

func (t *TestTestingFT) Fatalf(format string, args ...interface{}) {
	t.Fatal(fmt.Sprintf(format, args...))
}

func (t *TestTestingFT) Log(args ...interface{})                 {}
func (t *TestTestingFT) Logf(format string, args ...interface{}) {}

func (t *TestTestingFT) Name() string {
	return ""
}

func (t *TestTestingFT) Skip(args ...interface{})                 {}
func (t *TestTestingFT) SkipNow()                                 {}
func (t *TestTestingFT) Skipf(format string, args ...interface{}) {}
func (t *TestTestingFT) Skipped() bool {
	return false
}

func (t *TestTestingFT) Run(name string, f func(t *testing.T)) bool {
	return true
}
