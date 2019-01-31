package ctxerr_test

import (
	"bytes"
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func errorSummaryToString(s ctxerr.ErrorSummary, prefix string) string {
	var buf bytes.Buffer
	s.AppendSummary(&buf, prefix)
	return buf.String()
}

func fchomp(s string) string {
	if s[0] == '\n' {
		return s[1:]
	}
	return s
}

func TestErrorSummary(t *testing.T) {
	//
	// errorSummaryString
	summary := ctxerr.NewSummary("foobar")

	test.EqualStr(t, errorSummaryToString(summary, ""), `foobar`)
	test.EqualStr(t, errorSummaryToString(summary, "----"), `----foobar`)

	summary = ctxerr.NewSummary("foo\nbar")
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----foo
----bar`))

	//
	// ErrorSummaryItem
	summary = ctxerr.ErrorSummaryItem{
		Label: "  the_label: ",
		Value: "foo\nbar",
	}
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----  the_label: foo
----             bar`))

	summary = ctxerr.ErrorSummaryItem{
		Label:       "  the_label: ",
		Value:       "foo\nbar",
		Explanation: "And the\nexplanation...",
	}
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----  the_label: foo
----             bar
----And the
----explanation...`))

	//
	// ErrorSummaryItems
	summary = ctxerr.ErrorSummaryItems{
		{
			Label:       "first label: ",
			Value:       "foo\nbar",
			Explanation: "And the\nexplanation...",
		},
		{
			Label: "  2nd label: ",
			Value: "zip\nzap",
		},
		{
			Label: "  3rd label: ",
			Value: "666",
		},
	}
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----first label: foo
----             bar
----And the
----explanation...
----  2nd label: zip
----             zap
----  3rd label: 666`))

	//
	// NewSummaryReason
	summary = ctxerr.NewSummaryReason(666, "")
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----  value: 666
----it failed but didn't say why`))

	summary = ctxerr.NewSummaryReason(666, "evil number not accepted!")
	test.EqualStr(t, errorSummaryToString(summary, "----"), fchomp(`
----        value: 666
----it failed coz: evil number not accepted!`))
}
