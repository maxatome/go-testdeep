package testdeep

import (
	"reflect"
	"testing"
)

func TestGetInterface(t *testing.T) {
	type Private struct {
		private *Private // nolint: megacheck
	}

	// Cases not tested by TestEqualOthers()
	s := Private{
		private: &Private{},
	}

	_, ok := getInterface(reflect.ValueOf(s).Field(0), false)
	if ok {
		t.Error("getInterface() should return false for private field")
	}

	_, ok = getInterface(reflect.ValueOf(s).Field(0), true)
	if UnsafeDisabled {
		if ok {
			t.Error("unsafe package is disabled, getInterface should fail")
		}
	} else if !ok {
		t.Error("unsafe package is available, getInterface should succeed")
	}

	//var (
	//	panicked   bool
	//	panicParam interface{}
	//)
	//
	//func() {
	//	defer func() { panicParam = recover() }()
	//	panicked = true
	//	mustGetInterface(reflect.ValueOf(s).Field(0))
	//	panicked = false
	//}()
	//
	//if panicked {
	//	panicStr, ok := panicParam.(string)
	//	if ok {
	//		const expectedPanic = "getInterface() does not handle map kind"
	//		if panicStr != expectedPanic {
	//			t.Errorf("panic() string `%s' ≠ `%s'", panicStr, expectedPanic)
	//		}
	//	} else {
	//		t.Errorf("panic() occurred but recover()d %T type instead of string",
	//			panicParam)
	//	}
	//} else {
	//	t.Error("panic() did not occur")
	//}
}
