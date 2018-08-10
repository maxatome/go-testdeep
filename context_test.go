// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestContext(t *testing.T) {
	test.EqualStr(t, NewContext().path, "DATA")
	test.EqualStr(t, NewBooleanContext().path, "")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddDepth(".foo").
			path,
		"test.foo")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddDepth(".foo").
			AddDepth(".bar").
			path,
		"test.foo.bar")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "*test"}).
			AddDepth(".foo").
			path,
		"(*test).foo")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddArrayIndex(12).path,
		"test[12]")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "*test"}).
			AddArrayIndex(12).path,
		"(*test)[12]")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test"}).
			AddPtr(2).
			path,
		"**test")

	test.EqualStr(t,
		NewContextWithConfig(ContextConfig{RootName: "test.foo"}).
			AddPtr(1).path, "*test.foo")

	test.EqualStr(t,
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
	test.EqualInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "aaa")
	test.EqualInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "-8")
	test.EqualInt(t, getMaxErrorsFromEnv(), -8)
}
