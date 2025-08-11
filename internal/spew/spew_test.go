/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/spew"
	"github.com/maxatome/go-testdeep/internal/test"
)

// spewTest is used to describe a test to be performed against the public
// functions of the spew package or ConfigState.
type spewTest struct {
	cs   *spew.ConfigState
	in   any
	want string
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var spewTests []spewTest

func initSpewTests() {
	scsNoMethods := &spew.ConfigState{
		Indent:                  " ",
		DisableMethods:          true,
		DisablePointerAddresses: true,
	}
	scsNoPmethods := &spew.ConfigState{
		Indent:                  " ",
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
	}
	scsMaxDepth := &spew.ConfigState{
		Indent:                  " ",
		MaxDepth:                1,
		DisablePointerAddresses: true,
	}
	scsNoPtrAddr := &spew.ConfigState{DisablePointerAddresses: true}
	scsCap := &spew.ConfigState{EnableCapacities: true}

	// Variables for tests on types which implement fmt.Stringer
	// interface with and without a pointer receiver.
	ts := stringer("test")
	tps := pstringer("test")

	type ptrTester struct {
		s *struct{}
	}
	tptr := &ptrTester{s: &struct{}{}}

	// depthTester is used to test max depth handling for structs, array, slices
	// and maps.
	type depthTester struct {
		ic    indirCir1
		arr   [1]string
		slice []string
		m     map[string]int
	}
	dt := depthTester{
		indirCir1{&indirCir2{}},
		[1]string{"arr"},
		[]string{"slice"},
		map[string]int{"one": 1},
	}
	dz := depthTester{}

	spewTests = []spewTest{
		{nil, int8(127), "(int8) 127"},
		{nil, int16(32767), "(int16) 32767"},
		{nil, 2147483647, "(int) 2147483647"},
		{nil, int64(9223372036854775807), "(int64) 9223372036854775807"},
		{nil, uint8(64), "(uint8) 64"},
		{nil, complex(1, 2), "(complex128) (1+2i)"},
		{nil, complex(float64(5), 6), "(complex128) (5+6i)"},
		{nil, float32(3.14), "(float32) 3.14"},
		{nil, float64(6.28), "(float64) 6.28"},
		{nil, true, "(bool) true"},
		{nil, false, "(bool) false"},
		{nil, complex(-10, -20), "(complex128) (-10-20i)"},
		{nil, dz, "(spew_test.depthTester) {\n}"},
		{scsMaxDepth, dt, `(spew_test.depthTester) {
 ic: (spew_test.indirCir1) {
  <max depth reached>
 },
 arr: ([1]string) (len=1) {
  <max depth reached>
 },
 slice: ([]string) (len=1) {
  <max depth reached>
 },
 m: (map[string]int) (len=1) {
  <max depth reached>
 }
}`},
		// scsNoMethods
		{scsNoMethods, ts, `(spew_test.stringer) (len=4) "test"`},
		{scsNoMethods, &ts, `(*spew_test.stringer)((len=4) "test")`},
		{scsNoMethods, tps, `(spew_test.pstringer) (len=4) "test"`},
		{scsNoMethods, &tps, `(*spew_test.pstringer)((len=4) "test")`},
		// scsNoPmethods
		{scsNoPmethods, ts, "(spew_test.stringer) (len=4) stringer test"},
		{scsNoPmethods, &ts, "(*spew_test.stringer)((len=4) stringer test)"},
		{scsNoPmethods, tps, `(spew_test.pstringer) (len=4) "test"`},
		{scsNoPmethods, &tps, "(*spew_test.pstringer)((len=4) stringer test)"},
		{scsNoPtrAddr, tptr, "(*spew_test.ptrTester)({\ns: (*struct {})({\n})\n})"},
		{scsCap, make([]string, 0, 10), "([]string) (cap=10) {\n}"},
		{scsCap, make([]string, 1, 10), "([]string) (len=1 cap=10) {\n(string) \"\"\n}"},
	}
}

// TestSpew executes all of the tests described by spewTests.
func TestSpew(t *testing.T) {
	initSpewTests()

	t.Logf("Running %d tests", len(spewTests))
	for i, tc := range spewTests {
		var cs []spew.ConfigState
		if tc.cs != nil {
			cs = append(cs, *tc.cs)
		}
		str := spew.Sdump(tc.in, cs...)
		test.EqualStr(t, str, tc.want, "#%d", i)
	}
}
