// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/color"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func errorSummaryToString(s ctxerr.ErrorSummary, prefix string, uncolorized bool) string {
	var buf strings.Builder
	s.AppendSummary(&buf, prefix, !uncolorized)
	return buf.String()
}

func TestErrorSummary(t *testing.T) {
	defer color.SaveState()()

	var expectedColorized bool
	r := func(s string) string {
		if s[0] == '\n' {
			s = s[1:]
		}
		var repl *strings.Replacer
		if expectedColorized {
			repl = strings.NewReplacer(
				"*", "\x1b[1;31m", // bold red
				"+", "\x1b[0;31m", // red light
				"^", "\x1b[0m", // red off
				"~", "", // just ignore, for vertical alignment purpose
			)
		} else {
			repl = strings.NewReplacer(
				"*", "", // bold red
				"+", "", // red light
				"^", "", // red off
				"~", "", // just ignore, for vertical alignment purpose
			)
		}
		return repl.Replace(s)
	}

	testCases := []struct {
		name              string
		envColorized      bool
		forceUncolorized  bool
		expectedColorized bool
	}{
		{
			name:              "no color via env",
			envColorized:      false,
			expectedColorized: false,
		},
		{
			name:              "colorized",
			envColorized:      true,
			expectedColorized: true,
		},
		{
			name:              "colorized, but force uncolorized",
			envColorized:      true,
			forceUncolorized:  true,
			expectedColorized: false,
		},
		{
			name:              "no color via env and force uncolorized",
			envColorized:      false,
			forceUncolorized:  true,
			expectedColorized: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedColorized = tc.expectedColorized

			color.SaveState(tc.envColorized)

			t.Logf("colorized=%t force=%t expected=%t", tc.envColorized, tc.forceUncolorized, tc.expectedColorized)

			//
			// errorSummaryString
			summary := ctxerr.NewSummary("foobar")

			test.EqualStr(t, errorSummaryToString(summary, "", tc.forceUncolorized),
				r(`+foobar^`))
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized),
				r(`----+foobar^`))

			summary = ctxerr.NewSummary("foo\nbar")
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----+foo^
----+bar^`))

			//
			// ErrorSummaryItem
			summary = ctxerr.ErrorSummaryItem{
				Label: "the_label",
				Value: "foo\nbar",
			}
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----*the_label: +foo^
----~           +bar^`))

			summary = ctxerr.ErrorSummaryItem{
				Label:       "the_label",
				Value:       "foo\nbar",
				Explanation: "And the\nexplanation...",
			}
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----*the_label: +foo^
----~           +bar^
----+And the^
----+explanation...^`))

			//
			// ErrorSummaryItems
			summary = ctxerr.ErrorSummaryItems{
				{
					Label:       "first label",
					Value:       "foo\nbar",
					Explanation: "And the\nexplanation...",
				},
				{
					Label: "2nd label",
					Value: "zip\nzap",
				},
				{
					Label: "3rd big label",
					Value: "666",
				},
			}
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----*  first label: +foo^
----~               +bar^
----+And the^
----+explanation...^
----*    2nd label: +zip^
----~               +zap^
----*3rd big label: +666^`))

			//
			// NewSummaryReason
			summary = ctxerr.NewSummaryReason(666, "")
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----*  value: +666^
----+it failed but didn't say why^`))

			summary = ctxerr.NewSummaryReason(666, "evil number not accepted!")
			test.EqualStr(t, errorSummaryToString(summary, "----", tc.forceUncolorized), r(`
----*        value: +666^
----*it failed coz: +evil number not accepted!^`))
		})
	}
}
