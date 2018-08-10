// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/location"
)

var (
	testDeeper         = reflect.TypeOf((*TestDeep)(nil)).Elem()
	stringerInterface  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	timeType           = reflect.TypeOf(time.Time{})
	intType            = reflect.TypeOf(int(0))
	smuggledGotType    = reflect.TypeOf(SmuggledGot{})
	smuggledGotPtrType = reflect.TypeOf((*SmuggledGot)(nil))
	reflectValueType   = reflect.TypeOf(reflect.Value{})
)

// TestingT is the minimal interface used by CmpDeeply to report
// errors. It is commonly implemented by *testing.T and testing.TB.
type TestingT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Helper()
}

// TestingFT (aka. TestingF<ull>T) is the interface used by T to
// delegate common *testing.T functions to it. Of course, *testing.T
// implements it.
type TestingFT interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	Helper()
	Run(name string, f func(t *testing.T)) bool
}

type testDeepStringer interface {
	_TestDeep()
	String() string
}

// TestDeep is the representation of a testdeep operator. It is not
// intended to be used directly, but through Cmp* functions.
type TestDeep interface {
	testDeepStringer
	Match(ctx Context, got reflect.Value) *Error
	location.GetLocationer
	setLocation(int)
	HandleInvalid() bool
	TypeBehind() reflect.Type
}

// Base is a base type providing some methods needed by the TestDeep
// interface.
type Base struct {
	location location.Location
}

func (t Base) _TestDeep() {}

func (t *Base) setLocation(callDepth int) {
	var ok bool
	t.location, ok = location.New(callDepth)
	if !ok {
		t.location.File = "???"
		t.location.Line = 0
		return
	}

	opDotPos := strings.LastIndex(t.location.Func, ".")

	// Try to go one level upper, to check if it is a CmpXxx function
	cmpLoc, ok := location.New(callDepth + 1)
	if ok {
		cmpDotPos := strings.LastIndex(cmpLoc.Func, ".")

		// Must be in same package as found operator
		if t.location.Func[:opDotPos] == cmpLoc.Func[:cmpDotPos] &&
			strings.HasPrefix(cmpLoc.Func[cmpDotPos+1:], "Cmp") &&
			cmpLoc.Func != "CmpDeeply" {
			t.location = cmpLoc
			opDotPos = cmpDotPos
		}
	}

	t.location.Func = t.location.Func[opDotPos+1:]
}

// GetLocation returns a copy of the location.Location where the TestDeep
// operator has been created.
func (t *Base) GetLocation() location.Location {
	return t.location
}

// HandleInvalid tells testdeep internals that this operator does not
// handle nil values directly.
func (t Base) HandleInvalid() bool {
	return false
}

// TypeBehind returns the type handled by the operator. Only few operators
// knows the type they are handling. If they do not know, nil is
// returned.
func (t Base) TypeBehind() reflect.Type {
	return nil
}

// NewBase returns a new Base struct with location.Location set to the
// "callDepth" depth.
func NewBase(callDepth int) (b Base) {
	b.setLocation(callDepth)
	return
}

// BaseOKNil is a base type providing some methods needed by the TestDeep
// interface, for operators handling nil values.
type BaseOKNil struct {
	Base
}

// HandleInvalid tells testdeep internals that this operator handles
// nil values directly.
func (t BaseOKNil) HandleInvalid() bool {
	return true
}

// NewBaseOKNil returns a new BaseOKNil struct with location.Location set to
// the "callDepth" depth.
func NewBaseOKNil(callDepth int) (b BaseOKNil) {
	b.setLocation(callDepth)
	return
}

// Implements testDeepStringer
type rawString string

func (s rawString) _TestDeep() {}

func (s rawString) String() string {
	return string(s)
}

// Implements testDeepStringer
type rawInt int

func (i rawInt) _TestDeep() {}

func (i rawInt) String() string {
	return strconv.Itoa(int(i))
}
