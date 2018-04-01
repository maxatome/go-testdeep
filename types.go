package testdeep

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var testDeeper = reflect.TypeOf((*TestDeep)(nil)).Elem()

var stringerInterface = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

var timeType = reflect.TypeOf(time.Time{})

type testDeepStringer interface {
	___testDeep___()
	String() string
}

type Location struct {
	File string
	Func string
	Line int
}

func (l Location) IsInitialized() bool {
	return l.File != ""
}
func (l Location) String() string {
	return fmt.Sprintf("%s at %s:%d", l.Func, l.File, l.Line)
}

type TestDeep interface {
	testDeepStringer
	Match(ctx Context, got reflect.Value) *Error
	setLocation(int)
	GetLocation() Location
	HandleInvalid() bool
}

type TestDeepBase struct {
	location Location
}

func (t TestDeepBase) ___testDeep___() {}

func (t *TestDeepBase) setLocation(callDepth int) {
	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		if index := strings.LastIndexAny(file, `/\`); index >= 0 {
			file = file[index+1:]
		}
		t.location.File = file
		t.location.Line = line

		// Try to get the involved TestDeep operator
		pc, _, _, ok := runtime.Caller(callDepth - 1)
		if ok {
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				t.location.Func = fn.Name()
				if index := strings.LastIndex(t.location.Func, "."); index >= 0 {
					t.location.Func = t.location.Func[index+1:]
				}
			}
		}
	} else {
		t.location.File = "???"
		t.location.Line = 0
	}
}

func (t *TestDeepBase) GetLocation() Location {
	return t.location
}

func (t TestDeepBase) HandleInvalid() bool {
	return false
}

func NewTestDeepBase(callDepth int) (b TestDeepBase) {
	b.setLocation(callDepth)
	return
}

type TestDeepBaseOKNil struct {
	TestDeepBase
}

func (t TestDeepBaseOKNil) HandleInvalid() bool {
	return true
}

func NewTestDeepBaseOKNil(callDepth int) (b TestDeepBaseOKNil) {
	b.setLocation(callDepth)
	return
}

// Implements testDeepStringer
type rawString string

func (s rawString) ___testDeep___() {}

func (s rawString) String() string {
	return string(s)
}

// Implements testDeepStringer
type rawInt int

func (i rawInt) ___testDeep___() {}

func (i rawInt) String() string {
	return strconv.Itoa(int(i))
}
