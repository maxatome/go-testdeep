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

type TestingT struct {
	LastMessage string
	IsFatal     bool
	HasFailed   bool
}

func (t *TestingT) Error(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = false
	t.HasFailed = true
}

func (t *TestingT) Fatal(args ...interface{}) {
	t.LastMessage = fmt.Sprint(args...)
	t.IsFatal = true
	t.HasFailed = true
}

func (t *TestingT) Helper() {
	// Do nothing
}

type TestingFT struct {
	TestingT
}

func (t *TestingFT) Errorf(format string, args ...interface{}) {
	t.Error(fmt.Sprintf(format, args...))
}

func (t *TestingFT) Fail() {
	t.HasFailed = true
}

func (t *TestingFT) FailNow() {
	t.HasFailed = true
	t.IsFatal = true
}

func (t *TestingFT) Failed() bool {
	return t.HasFailed
}

func (t *TestingFT) Fatalf(format string, args ...interface{}) {
	t.Fatal(fmt.Sprintf(format, args...))
}

func (t *TestingFT) Log(args ...interface{})                 {}
func (t *TestingFT) Logf(format string, args ...interface{}) {}

func (t *TestingFT) Name() string {
	return ""
}

func (t *TestingFT) Skip(args ...interface{})                 {}
func (t *TestingFT) SkipNow()                                 {}
func (t *TestingFT) Skipf(format string, args ...interface{}) {}
func (t *TestingFT) Skipped() bool {
	return false
}

func (t *TestingFT) Run(name string, f func(t *testing.T)) bool {
	return true
}
