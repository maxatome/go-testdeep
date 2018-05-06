package testdeep

import (
	"fmt"
	"runtime"
	"strings"
)

type Location struct {
	File string
	Func string
	Line int
}

func NewLocation(callDepth int) (loc Location, ok bool) {
	_, loc.File, loc.Line, ok = runtime.Caller(callDepth + 1)
	if !ok {
		return
	}

	if index := strings.LastIndexAny(loc.File, `/\`); index >= 0 {
		loc.File = loc.File[index+1:]
	}

	pc, _, _, ok := runtime.Caller(callDepth)
	if !ok {
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn != nil {
		loc.Func = fn.Name()
	} else {
		ok = false
	}
	return
}

func (l Location) IsInitialized() bool {
	return l.File != ""
}
func (l Location) String() string {
	return fmt.Sprintf("%s at %s:%d", l.Func, l.File, l.Line)
}
