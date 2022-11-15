// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestContext(t *testing.T) {
	nctx := newContext(nil)
	test.EqualStr(t, nctx.Path.String(), "DATA")
	if nctx.OriginalTB != nil {
		t.Error("OriginalTB should be nil")
	}

	nctx = newContext(t)
	test.EqualStr(t, nctx.Path.String(), "DATA")
	if nctxt, ok := nctx.OriginalTB.(*testing.T); test.IsTrue(t, ok, "%T", nctx.OriginalTB) {
		if nctxt != t {
			t.Errorf("OriginalTB, got=%p expected=%p", nctxt, t)
		}
	}

	nctx = newContext(Require(t).UseEqual().TestDeepInGotOK())
	_, ok := nctx.OriginalTB.(*T)
	test.IsTrue(t, ok)
	test.IsTrue(t, nctx.FailureIsFatal)
	test.IsTrue(t, nctx.UseEqual)
	test.IsTrue(t, nctx.TestDeepInGotOK)
	test.EqualStr(t, nctx.Path.String(), "DATA")

	nctx = newBooleanContext()
	test.EqualStr(t, nctx.Path.String(), "")
	if nctx.OriginalTB != nil {
		t.Error("OriginalTB should be nil")
	}

	if newContextWithConfig(nil, ContextConfig{MaxErrors: -1}).CollectError(nil) != nil {
		t.Errorf("ctx.CollectError(nil) should return nil")
	}

	ctx := ContextConfig{}
	if ctx.Equal(DefaultContextConfig) {
		t.Errorf("Empty ContextConfig should be ≠ from DefaultContextConfig")
	}
	ctx.sanitize()
	if !ctx.Equal(DefaultContextConfig) {
		t.Errorf("Sanitized empty ContextConfig should be = to DefaultContextConfig")
	}

	ctx.RootName = "PIPO"
	test.EqualStr(t, ctx.OriginalPath(), "PIPO")

	nctx = newContext(t)
	nctx.Path = ctxerr.NewPath("BINGO[0].Zip")
	ctx.forkedFromCtx = &nctx
	test.EqualStr(t, ctx.OriginalPath(), "BINGO[0].Zip")
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
