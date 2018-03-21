package testdeep

import (
	"fmt"
	"reflect"
	"strings"
)

type stringKind uint8

const (
	exactString stringKind = iota
	prefixString
	suffixString
	containString
)

type tdStringBase struct {
	TestDeepBase
	expected string
}

func newStringBase(expected string) tdStringBase {
	return tdStringBase{
		TestDeepBase: NewTestDeepBase(4),
		expected:     expected,
	}
}

func getString(ctx Context, got reflect.Value) (string, *Error) {
	switch got.Kind() {
	case reflect.String:
		return got.String(), nil

	default:
		if got.CanInterface() {
			switch iface := got.Interface().(type) {
			case fmt.Stringer:
				return iface.String(), nil
			case error:
				return iface.Error(), nil
			}
		}
	}

	if ctx.booleanError {
		return "", booleanError
	}
	return "", &Error{
		Context: ctx,
		Message: "bad type",
		Got:     rawString(got.Type().String()),
		Expected: rawString(
			"string (convertible) OR fmt.Stringer OR error"),
	}
}

type tdString struct {
	tdStringBase
}

var _ TestDeep = &tdString{}

func String(expected string) TestDeep {
	return &tdString{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdString) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		err.Location = s.GetLocation()
		return err
	}

	if str == s.expected {
		return nil
	}
	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "does not match",
		Got:      str,
		Expected: s,
		Location: s.GetLocation(),
	}
}

func (s *tdString) String() string {
	return toString(s.expected)
}

type tdHasPrefix struct {
	tdStringBase
}

var _ TestDeep = &tdHasPrefix{}

func HasPrefix(expected string) TestDeep {
	return &tdHasPrefix{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdHasPrefix) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		err.Location = s.GetLocation()
		return err
	}

	if strings.HasPrefix(str, s.expected) {
		return nil
	}
	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "has not prefix",
		Got:      str,
		Expected: s,
		Location: s.GetLocation(),
	}
}

func (s *tdHasPrefix) String() string {
	return "HasPrefix(" + toString(s.expected) + ")"
}

type tdHasSuffix struct {
	tdStringBase
}

var _ TestDeep = &tdHasSuffix{}

func HasSuffix(expected string) TestDeep {
	return &tdHasSuffix{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdHasSuffix) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		err.Location = s.GetLocation()
		return err
	}

	if strings.HasSuffix(str, s.expected) {
		return nil
	}
	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "has not suffix",
		Got:      str,
		Expected: s,
		Location: s.GetLocation(),
	}
}

func (s *tdHasSuffix) String() string {
	return "HasSuffix(" + toString(s.expected) + ")"
}

type tdContains struct {
	tdStringBase
}

var _ TestDeep = &tdContains{}

func Contains(expected string) TestDeep {
	return &tdContains{
		tdStringBase: newStringBase(expected),
	}
}

func (s *tdContains) Match(ctx Context, got reflect.Value) *Error {
	str, err := getString(ctx, got)
	if err != nil {
		err.Location = s.GetLocation()
		return err
	}

	if strings.Contains(str, s.expected) {
		return nil
	}
	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "does not contain",
		Got:      str,
		Expected: s,
		Location: s.GetLocation(),
	}
}

func (s *tdContains) String() string {
	return "Contains(" + toString(s.expected) + ")"
}
