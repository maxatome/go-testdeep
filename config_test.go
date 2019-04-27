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
	test.EqualStr(t, newContext().Path.String(), "DATA")
	test.EqualStr(t, newBooleanContext().Path.String(), "")

	if newContextWithConfig(ContextConfig{MaxErrors: -1}).CollectError(nil) != nil {
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
	oldEnv, set := os.LookupEnv(envMaxErrors)
	defer func() {
		if set {
			os.Setenv(envMaxErrors, oldEnv)
		} else {
			os.Unsetenv(envMaxErrors)
		}
	}()

	os.Setenv(envMaxErrors, "")
	test.EqualInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "aaa")
	test.EqualInt(t, getMaxErrorsFromEnv(), 10)

	os.Setenv(envMaxErrors, "-8")
	test.EqualInt(t, getMaxErrorsFromEnv(), -8)
}
