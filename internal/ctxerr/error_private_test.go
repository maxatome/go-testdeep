// Copyright (c) 2019 Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr

import (
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
)

func TestColor(t *testing.T) {
	envColorV := os.Getenv("TESTDEEP_COLOR")
	defer func() { os.Setenv("TESTDEEP_COLOR", envColorV) }()

	// off
	for _, flag := range []string{"off", "xxbad"} {
		os.Setenv("TESTDEEP_COLOR", flag)
		os.Setenv("MY_TEST_COLOR", "green")
		light, bold, off := colorFromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "")
		test.EqualStr(t, bold, "")
		test.EqualStr(t, off, "")
	}

	// on
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
	}
}
