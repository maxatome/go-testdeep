// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/maxatome/go-testdeep/internal/location"
	"github.com/maxatome/go-testdeep/internal/util"
)

const (
	envColor         = "TESTDEEP_COLOR"
	envColorTestName = "TESTDEEP_COLOR_TEST_NAME"
	envColorTitle    = "TESTDEEP_COLOR_TITLE"
	envColorOK       = "TESTDEEP_COLOR_OK"
	envColorBad      = "TESTDEEP_COLOR_BAD"
)

var (
	colorTestNameOn, colorTestNameOff       string
	colorTitleOn, colorTitleOff             string
	colorOKOn, colorOKOnBold, colorOKOff    string
	colorBadOn, colorBadOnBold, colorBadOff string
)

var colorsInitOnce sync.Once

func colorsInit() {
	colorsInitOnce.Do(func() {
		_, colorTestNameOn, colorTestNameOff = colorFromEnv(envColorTestName, "yellow")
		_, colorTitleOn, colorTitleOff = colorFromEnv(envColorTitle, "cyan")
		colorOKOn, colorOKOnBold, colorOKOff = colorFromEnv(envColorOK, "green")
		colorBadOn, colorBadOnBold, colorBadOff = colorFromEnv(envColorBad, "red")
	})
}

// SaveColorState saves the "TESTDEEP_COLOR" environment variable
// value, sets it to "on" (if true passed as on) or "false" (if on not
// passed or set to false), resets the colors initialization and
// returns a function to be called in a defer statement. Only intended
// to be used in tests like:
//
//   defer ctxerr.SaveColorState()()
//
// It is not thread-safe.
func SaveColorState(on ...bool) func() {
	colorState, set := os.LookupEnv(envColor)
	if len(on) == 0 || !on[0] {
		os.Setenv(envColor, "off") // nolint: errcheck
	} else {
		os.Setenv(envColor, "on") // nolint: errcheck
	}
	colorsInitOnce = sync.Once{}
	return func() {
		if set {
			os.Setenv(envColor, colorState) // nolint: errcheck
		} else {
			os.Unsetenv(envColor) // nolint: errcheck
		}
		colorsInitOnce = sync.Once{}
	}
}

var colors = map[string]byte{
	"black":   '0',
	"red":     '1',
	"green":   '2',
	"yellow":  '3',
	"blue":    '4',
	"magenta": '5',
	"cyan":    '6',
	"white":   '7',
	"gray":    '7',
}

func colorFromEnv(env, defaultColor string) (string, string, string) {
	var color string
	switch os.Getenv(envColor) {
	case "on", "":
		if curColor := os.Getenv(env); curColor != "" {
			color = curColor
		} else {
			color = defaultColor
		}
	default: // "off" or any other value
		color = ""
	}

	if color == "" {
		return "", "", ""
	}

	names := strings.SplitN(color, ":", 2)

	light := [...]byte{
		//   0    1    2    4    4    5    6
		'\x1b', '[', '0', ';', '3', 'y', 'm', // foreground
		//   7    8    9   10   11
		'\x1b', '[', '4', 'z', 'm', // background
	}
	bold := [...]byte{
		//   0    1    2    4    4    5    6
		'\x1b', '[', '1', ';', '3', 'y', 'm', // foreground
		//   7    8    9   10   11
		'\x1b', '[', '4', 'z', 'm', // background
	}

	var start, end int

	// Foreground
	if names[0] != "" {
		c := colors[names[0]]
		if c == 0 {
			c = colors[defaultColor]
		}

		light[5] = c
		bold[5] = c

		end = 7
	} else {
		start = 7
	}

	// Background
	if len(names) > 1 && names[1] != "" {
		c := colors[names[1]]
		if c != 0 {
			light[10] = c
			bold[10] = c

			end = 12
		}
	}

	return string(light[start:end]), string(bold[start:end]), "\x1b[0m"
}

// ColorizeTestNameOn enable test name color in b.
func ColorizeTestNameOn(b *bytes.Buffer) {
	colorsInit()
	b.WriteString(colorTestNameOn)
}

// ColorizeTestNameOff disable test name color in b.
func ColorizeTestNameOff(b *bytes.Buffer) {
	colorsInit()
	b.WriteString(colorTestNameOff)
}

// Bad returns a string surrounded by BAD color. If len(args) is > 0,
// s and args are given to fmt.Sprintf.
//
// Typically used in panic() when the user made a mistake.
func Bad(s string, args ...interface{}) string {
	colorsInit()
	if len(args) == 0 {
		return colorBadOnBold + s + colorBadOff
	}
	return fmt.Sprintf(colorBadOnBold+s+colorBadOff, args...)
}

// BadUsage returns a string surrounded by BAD color to notice the
// user he passes a bad parameter to a function. Typically used in a
// panic().
func BadUsage(usage string, param interface{}, pos int, kind bool) string {
	colorsInit()

	var b bytes.Buffer
	fmt.Fprintf(&b, "%susage: %s, but received ", colorBadOnBold, usage)

	if param == nil {
		b.WriteString("nil")
	} else {
		t := reflect.TypeOf(param)
		if kind && t.String() != t.Kind().String() {
			fmt.Fprintf(&b, "%s (%s)", t, t.Kind())
		} else {
			b.WriteString(t.String())
		}
	}

	b.WriteString(" as ")
	switch pos {
	case 1:
		b.WriteString("1st")
	case 2:
		b.WriteString("2nd")
	case 3:
		b.WriteString("3rd")
	default:
		fmt.Fprintf(&b, "%dth", pos)
	}
	b.WriteString(" parameter")
	b.WriteString(colorBadOff)
	return b.String()
}

// TooManyParams returns a string surrounded by BAD color to notice
// the user he called a variadic function with too many
// parameters. Typically used in a panic().
func TooManyParams(usage string) string {
	colorsInit()
	return colorBadOnBold + "usage: " + usage + ", too many parameters" + colorBadOff
}

// Error represents errors generated by td (go-testdeep) functions.
type Error struct {
	// Context when the error occurred
	Context Context
	// Message describes the error
	Message string
	// Got value
	Got interface{}
	// Expected value
	Expected interface{}
	// If not nil, Summary is used to display summary instead of using
	// Got + Expected fields
	Summary ErrorSummary
	// If initialized, location of TestDeep operator originator of the error
	Location location.Location
	// If defined, the current Error comes from this Error
	Origin *Error
	// If defined, points to the next Error
	Next *Error
}

var (
	// BooleanError is the *Error returned when an error occurs in a
	// boolean context.
	BooleanError = &Error{}

	// ErrTooManyErrors is chained to the last error encountered when
	// the maximum number of errors has been reached.
	ErrTooManyErrors = &Error{
		Message: "Too many errors (use TESTDEEP_MAX_ERRORS=-1 to see all)",
	}
)

// Error implements error interface.
func (e *Error) Error() string {
	buf := bytes.Buffer{}

	e.Append(&buf, "")

	return buf.String()
}

// Append appends the Error contents to "buf" using prefix "prefix"
// for each line.
func (e *Error) Append(buf *bytes.Buffer, prefix string) {
	if e == BooleanError {
		return
	}

	colorsInit()

	var writeEolPrefix func()
	if prefix != "" {
		eolPrefix := make([]byte, 1+len(prefix))
		eolPrefix[0] = '\n'
		copy(eolPrefix[1:], prefix)

		writeEolPrefix = func() {
			buf.Write(eolPrefix)
		}
		buf.WriteString(prefix)
	} else {
		writeEolPrefix = func() {
			buf.WriteByte('\n')
		}
	}

	if e == ErrTooManyErrors {
		buf.WriteString(colorTitleOn)
		buf.WriteString(e.Message)
		buf.WriteString(colorTitleOff)
		return
	}

	buf.WriteString(colorTitleOn)
	if pos := strings.Index(e.Message, "%%"); pos >= 0 {
		buf.WriteString(e.Message[:pos])
		buf.WriteString(e.Context.Path.String())
		buf.WriteString(e.Message[pos+2:])
	} else {
		buf.WriteString(e.Context.Path.String())
		buf.WriteString(": ")
		buf.WriteString(e.Message)
	}
	buf.WriteString(colorTitleOff)

	if e.Summary != nil {
		buf.WriteByte('\n')
		e.Summary.AppendSummary(buf, prefix+"\t")
	} else {
		writeEolPrefix()
		buf.WriteString(colorBadOnBold)
		buf.WriteString("\t     got: ")
		buf.WriteString(colorBadOn)
		util.IndentStringIn(buf, e.GotString(), prefix+"\t          ")
		buf.WriteString(colorBadOff)
		writeEolPrefix()
		buf.WriteString(colorOKOnBold)
		buf.WriteString("\texpected: ")
		buf.WriteString(colorOKOn)
		util.IndentStringIn(buf, e.ExpectedString(), prefix+"\t          ")
		buf.WriteString(colorOKOff)
	}

	// This error comes from another one
	if e.Origin != nil {
		writeEolPrefix()
		buf.WriteString("Originates from following error:\n")

		e.Origin.Append(buf, prefix+"\t")
	}

	if e.Location.IsInitialized() &&
		!e.Location.BehindCmp && // no need to log Cmp* func
		(e.Next == nil || e.Next.Location != e.Location) {
		writeEolPrefix()
		buf.WriteString("[under TestDeep operator ")
		buf.WriteString(e.Location.String())
		buf.WriteByte(']')
	}

	if e.Next != nil {
		buf.WriteByte('\n')
		e.Next.Append(buf, prefix) // next error at same level
	}
}

// GotString returns the string corresponding to the Got
// field. Returns the empty string if the Error Summary field is not
// nil.
func (e *Error) GotString() string {
	if e.Summary != nil {
		return ""
	}
	return util.ToString(e.Got)
}

// ExpectedString returns the string corresponding to the Expected
// field. Returns the empty string if the Error Summary field is not
// nil.
func (e *Error) ExpectedString() string {
	if e.Summary != nil {
		return ""
	}
	return util.ToString(e.Expected)
}

// SummaryString returns the string corresponding to the Summary
// field. Returns the empty string if the Error Summary field is nil.
func (e *Error) SummaryString() string {
	if e.Summary == nil {
		return ""
	}

	var buf bytes.Buffer
	e.Summary.AppendSummary(&buf, "")
	return buf.String()
}
