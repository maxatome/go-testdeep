// Copyright (c) 2019 Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr

import (
	"bytes"
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestColor(t *testing.T) {
	defer SaveColorState()()

	// off
	for _, flag := range []string{"off", "xxbad"} {
		os.Setenv("TESTDEEP_COLOR", flag)
		os.Setenv("MY_TEST_COLOR", "green")
		light, bold, off := colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "")
		test.EqualStr(t, bold, "")
		test.EqualStr(t, off, "")

		var b bytes.Buffer
		ColorizeTestNameOn(&b)
		test.EqualInt(t, b.Len(), 0)
		ColorizeTestNameOff(&b)
		test.EqualInt(t, b.Len(), 0)
	}

	// on
	colorTestNameOnSave, colorTestNameOffSave := colorTestNameOn, colorTestNameOff
	defer func() {
		colorTestNameOn, colorTestNameOff = colorTestNameOnSave, colorTestNameOffSave
	}()
	for _, flag := range []string{"on", ""} {
		os.Setenv("TESTDEEP_COLOR", flag)
		os.Setenv("MY_TEST_COLOR", "")
		light, bold, off := colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;31m")
		test.EqualStr(t, bold, "\x1b[1;31m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override
		os.Setenv("MY_TEST_COLOR", "green")
		light, bold, off = colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;32m")
		test.EqualStr(t, bold, "\x1b[1;32m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override including background
		os.Setenv("MY_TEST_COLOR", "green:magenta")
		light, bold, off = colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;32m\x1b[45m")
		test.EqualStr(t, bold, "\x1b[1;32m\x1b[45m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override including background only
		os.Setenv("MY_TEST_COLOR", ":magenta")
		light, bold, off = colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[45m")
		test.EqualStr(t, bold, "\x1b[45m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + bad colors
		os.Setenv("MY_TEST_COLOR", "foo:bar")
		light, bold, off = colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;31m") // red
		test.EqualStr(t, bold, "\x1b[1;31m")  // bold red
		test.EqualStr(t, off, "\x1b[0m")

		// Color test name
		_, colorTestNameOn, colorTestNameOff = colorFromEnv(envColorTitle, "yellow")
		var b bytes.Buffer
		ColorizeTestNameOn(&b)
		test.EqualStr(t, b.String(), "\x1b[1;33m")
		ColorizeTestNameOff(&b)
		test.EqualStr(t, b.String(), "\x1b[1;33m\x1b[0m")
	}
}
