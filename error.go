package testdeep

import (
	"bytes"
	"strings"
)

type Error struct {
	Context Context
	// Message describes the error
	Message string
	// Got value
	Got interface{}
	// Expected value
	Expected interface{}
	// If not nil, Summary is used to display summary instead of using
	// Got + Expected fields
	Summary interface{}
	// If initialized, location of TestDeep operator originator of the error
	Location Location
	// If defined, the current Error comes from this Error
	Origin *Error
}

var booleanError = &Error{}

func (e *Error) Error() string {
	if e == booleanError {
		return ""
	}

	buf := &bytes.Buffer{}

	if pos := strings.Index(e.Message, "%%"); pos >= 0 {
		buf.WriteString(e.Message[:pos])
		buf.WriteString(e.Context.Path)
		buf.WriteString(e.Message[pos+2:])
	} else {
		buf.WriteString(e.Context.Path)
		buf.WriteString(": ")
		buf.WriteString(e.Message)
	}

	buf.WriteByte('\n')

	if e.Summary != nil {
		buf.WriteByte('\t')
		buf.WriteString(indentString(e.SummaryString(), "\t"))
	} else {
		buf.WriteString("\t     got: ")
		buf.WriteString(indentString(e.GotString(), "\t          "))
		buf.WriteString("\n\texpected: ")
		buf.WriteString(indentString(e.ExpectedString(), "\t          "))
	}

	if e.Location.IsInitialized() {
		buf.WriteString("\n[under TestDeep operator ")
		buf.WriteString(e.Location.String())
		buf.WriteByte(']')
	}

	// This error comes from another one
	if e.Origin != nil {
		buf.WriteString("\nOriginates from following error:\n\t")
		buf.WriteString(indentString(e.Origin.Error(), "\t"))
	}

	return buf.String()
}

func (e *Error) GotString() string {
	if e.Summary != nil {
		return ""
	}
	return toString(e.Got)
}

func (e *Error) ExpectedString() string {
	if e.Summary != nil {
		return ""
	}
	return toString(e.Expected)
}

func (e *Error) SummaryString() string {
	if e.Summary == nil {
		return ""
	}
	return toString(e.Summary)
}

func (e *Error) SetLocationIfMissing(t TestDeep) *Error {
	if e != nil && !e.Location.IsInitialized() {
		e.Location = t.GetLocation()
	}
	return e
}
