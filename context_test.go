// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"os"
	"testing"
)

func equalStr(t *testing.T, got, expected string) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`Failed test
	     got: %s
	expected: %s`, got, expected)
	return false
}

func equalInt(t *testing.T, got, expected int) bool {
	if got == expected {
		return true
	}

	t.Helper()
	t.Errorf(`Failed test
	     got: %d
	expected: %d`, got, expected)
	return false
}

func TestContext(t *testing.T) {
	equalStr(t, NewContext().path, "DATA")
	equalStr(t, NewBooleanContext().path, "")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddDepth(".foo").
			path,
		"test.foo")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddDepth(".foo").
			AddDepth(".bar").
			path,
		"test.foo.bar")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "*test"}).
			AddDepth(".foo").
			path,
		"(*test).foo")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddArrayIndex(12).path,
		"test[12]")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "*test"}).
			AddArrayIndex(12).path,
		"(*test)[12]")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddPtr(2).
			path,
		"**test")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test.foo"}).
			AddPtr(1).path, "*test.foo")

	equalStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test[3]"}).
			AddPtr(1).path,
		"*test[3]")

	if NewContextWithConfig(ContextConfig{MaxErrors: -1}).CollectError(nil) != nil {
		t.Errorf("ctx.CollectError(nil) should return nil")
	}

	ctx := ContextConfig{}
	if ctx == DefaultContextConfig {
		t.Errorf("Empty ContextConfig should be ≠ from DefaultContextConfig")
	}
	ctx.sanitize()
	if ctx != DefaultContextConfig {
		t.Errorf("Sanitized empty ContextConfig should be = to DefaultContextConfig")
	}
}

func TestGetMaxErrorsFromEnv(t *testing.T) {
	oldEnv := os.Getenv(envMaxErrors)
	defer func() { os.Setenv(envMaxErrors, oldEnv) }()

	os.Setenv(envMaxErrors, "")
	equalInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "aaa")
	equalInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "-8")
	equalInt(t, getMaxErrorsFromEnv(), -8)
}
