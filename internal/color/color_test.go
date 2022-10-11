// Copyright (c) 2019, 2020 Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package color_test

import (
	"os"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestColor(t *testing.T) {
	defer color.SaveState()()

	// off
	for _, flag := range []string{"off", "xxbad"} {
		os.Setenv("TESTDEEP_COLOR", flag)
		os.Setenv("MY_TEST_COLOR", "green")
		light, bold, off := color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "")
		test.EqualStr(t, bold, "")
		test.EqualStr(t, off, "")

		var b strings.Builder
		color.AppendTestNameOn(&b)
		test.EqualInt(t, b.Len(), 0)
		color.AppendTestNameOff(&b)
		test.EqualInt(t, b.Len(), 0)
	}

	// on
	colorTestNameOnSave, colorTestNameOffSave := color.TestNameOn, color.TestNameOff
	defer func() {
		color.TestNameOn, color.TestNameOff = colorTestNameOnSave, colorTestNameOffSave
	}()
	for _, flag := range []string{"on", ""} {
		os.Setenv("TESTDEEP_COLOR", flag)
		os.Setenv("MY_TEST_COLOR", "")
		light, bold, off := color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;31m")
		test.EqualStr(t, bold, "\x1b[1;31m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override
		os.Setenv("MY_TEST_COLOR", "green")
		light, bold, off = color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;32m")
		test.EqualStr(t, bold, "\x1b[1;32m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override including background
		os.Setenv("MY_TEST_COLOR", "green:magenta")
		light, bold, off = color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;32m\x1b[45m")
		test.EqualStr(t, bold, "\x1b[1;32m\x1b[45m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + override including background only
		os.Setenv("MY_TEST_COLOR", ":magenta")
		light, bold, off = color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[45m")
		test.EqualStr(t, bold, "\x1b[45m")
		test.EqualStr(t, off, "\x1b[0m")

		// on + bad colors
		os.Setenv("MY_TEST_COLOR", "foo:bar")
		light, bold, off = color.FromEnv("MY_TEST_COLOR", "red")
		test.EqualStr(t, light, "\x1b[0;31m") // red
		test.EqualStr(t, bold, "\x1b[1;31m")  // bold red
		test.EqualStr(t, off, "\x1b[0m")

		// Color test name
		_, color.TestNameOn, color.TestNameOff = color.FromEnv(color.EnvColorTitle, "yellow")
		var b strings.Builder
		color.AppendTestNameOn(&b)
		test.EqualStr(t, b.String(), "\x1b[1;33m")
		color.AppendTestNameOff(&b)
		test.EqualStr(t, b.String(), "\x1b[1;33m\x1b[0m")
	}
}

func TestSaveState(t *testing.T) {
	check := func(expected string) {
		t.Helper()
		test.EqualStr(t, os.Getenv("TESTDEEP_COLOR"), expected)
	}

	defer color.SaveState()()
	check("off")

	func() {
		defer color.SaveState(true)()
		check("on")
	}()
	check("off")

	func() {
		defer color.SaveState(false)()
		check("off")
	}()
	check("off")

	os.Unsetenv("TESTDEEP_COLOR")
	checkDoesNotExist := func() {
		t.Helper()
		_, exists := os.LookupEnv("TESTDEEP_COLOR")
		test.IsFalse(t, exists)
	}

	func() {
		defer color.SaveState(true)()
		check("on")
	}()
	checkDoesNotExist()

	func() {
		defer color.SaveState(false)()
		check("off")
	}()
	checkDoesNotExist()
}

func TestBad(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t, color.Bad("test"), "test")
	test.EqualStr(t, color.Bad("test %d", 123), "test 123")
}

func TestBadUsage(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", nil, 1, true),
		"usage: Zzz(STRING), but received nil as 1st parameter")

	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", 42, 1, true),
		"usage: Zzz(STRING), but received int as 1st parameter")

	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", []int{}, 1, true),
		"usage: Zzz(STRING), but received []int (slice) as 1st parameter")
	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", []int{}, 1, false),
		"usage: Zzz(STRING), but received []int as 1st parameter")

	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", nil, 1, true),
		"usage: Zzz(STRING), but received nil as 1st parameter")
	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", nil, 2, true),
		"usage: Zzz(STRING), but received nil as 2nd parameter")
	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", nil, 3, true),
		"usage: Zzz(STRING), but received nil as 3rd parameter")
	test.EqualStr(t,
		color.BadUsage("Zzz(STRING)", nil, 4, true),
		"usage: Zzz(STRING), but received nil as 4th parameter")
}

func TestTooManyParams(t *testing.T) {
	defer color.SaveState()()

	test.EqualStr(t, color.TooManyParams("Zzz(PARAM)"),
		"usage: Zzz(PARAM), too many parameters")
}

func TestUnBad(t *testing.T) {
	defer color.SaveState(true)()

	const mesg = "test"
	s := color.Bad(mesg)
	if s == mesg {
		t.Errorf("Bad should produce colored output: %s ≠ %s", s, mesg)
	}
	test.EqualStr(t, color.UnBad(s), mesg)
}
