package ctxerr_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func errorSummaryToString(s ctxerr.ErrorSummary, prefix string) string {
	var buf bytes.Buffer
	s.AppendSummary(&buf, prefix)
	return buf.String()
}

func TestErrorSummary(t *testing.T) {
	defer ctxerr.SaveColorState()()

	colored := false
	color := func(enable bool) {
		colored = enable
		if enable {
			os.Setenv("TESTDEEP_COLOR", "on")
		} else {
			os.Setenv("TESTDEEP_COLOR", "off")
		}
		ctxerr.InitColors()
	}

	r := func(s string) string {
		if s[0] == '\n' {
			s = s[1:]
		}
		var repl *strings.Replacer
		if colored {
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

	for _, colored = range []bool{false, true} {
		color(colored)

		//
		// errorSummaryString
		summary := ctxerr.NewSummary("foobar")

		test.EqualStr(t, errorSummaryToString(summary, ""),
			r(`+foobar^`))
		test.EqualStr(t, errorSummaryToString(summary, "----"),
			r(`----+foobar^`))

		summary = ctxerr.NewSummary("foo\nbar")
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----+foo
----bar^`))

		//
		// ErrorSummaryItem
		summary = ctxerr.ErrorSummaryItem{
			Label: "the_label",
			Value: "foo\nbar",
		}
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----*the_label: +foo
----~           ~bar^`))

		summary = ctxerr.ErrorSummaryItem{
			Label:       "the_label",
			Value:       "foo\nbar",
			Explanation: "And the\nexplanation...",
		}
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----*the_label: +foo
----~           ~bar
----And the
----explanation...^`))

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
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----*  first label: +foo
----~               ~bar
----And the
----explanation...^
----*    2nd label: +zip
----~               ~zap^
----*3rd big label: +666^`))

		//
		// NewSummaryReason
		summary = ctxerr.NewSummaryReason(666, "")
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----*  value: +666
----~it failed but didn't say why^`))

		summary = ctxerr.NewSummaryReason(666, "evil number not accepted!")
		test.EqualStr(t, errorSummaryToString(summary, "----"), r(`
----*        value: +666^
----*it failed coz: +evil number not accepted!^`))
	}
}
