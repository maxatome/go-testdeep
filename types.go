package testdeep

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	testDeeper        = reflect.TypeOf((*TestDeep)(nil)).Elem()
	stringerInterface = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	timeType          = reflect.TypeOf(time.Time{})
)

type testDeepStringer interface {
	_TestDeep()
	String() string
}

type TestDeep interface {
	testDeepStringer
	Match(ctx Context, got reflect.Value) *Error
	setLocation(int)
	GetLocation() Location
	HandleInvalid() bool
}

type Base struct {
	location Location
}

func (t Base) _TestDeep() {}

func (t *Base) setLocation(callDepth int) {
	var ok bool
	t.location, ok = NewLocation(callDepth)
	if !ok {
		t.location.File = "???"
		t.location.Line = 0
		return
	}

	opDotPos := strings.LastIndex(t.location.Func, ".")

	// Try to go one level upper, to check if it is a CmpXxx function
	cmpLoc, ok := NewLocation(callDepth + 1)
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

func (t *Base) GetLocation() Location {
	return t.location
}

func (t Base) HandleInvalid() bool {
	return false
}

func NewBase(callDepth int) (b Base) {
	b.setLocation(callDepth)
	return
}

type BaseOKNil struct {
	Base
}

func (t BaseOKNil) HandleInvalid() bool {
	return true
}

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
