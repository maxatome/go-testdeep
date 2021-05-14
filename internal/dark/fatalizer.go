package dark

import (
	"fmt"
	"strings"
	"testing"
)

// Fatalizer is an interface used to raise a fatal error. It is the
// implementers responsibility that Fatal() never returns.
type Fatalizer interface {
	Helper()
	Fatal(args ...interface{})
}

// FatalPanic implements Fatalizer using panic().
type FatalPanic string

func (p FatalPanic) Helper() {}
func (p FatalPanic) Fatal(args ...interface{}) {
	panic(FatalPanic(fmt.Sprint(args...)))
}
func (p FatalPanic) String() string {
	return string(p)
}

// CheckFatalizerBarrierErr checks that "fn" called Fatal() on a
// Fatalizer and that the fatal message contains "contains".
func CheckFatalizerBarrierErr(t testing.TB, fn func(), contains string) bool {
	t.Helper()

	err := FatalizerBarrier(fn)
	if err == nil {
		t.Errorf("dark.FatalizerBarrier() did not return an error")
		return false
	}

	if !strings.Contains(err.Error(), contains) {
		t.Errorf("dark.FatalizerBarrier() error %q\ndoes not contain %q",
			err.Error(), contains)
		return false
	}
	return true
}

// Fatal calls Fatal of "f" followed by a panic in case "f" does not
// correctly die during the Fatal call.
func Fatal(f Fatalizer, args ...interface{}) {
	f.Helper()
	f.Fatal(args...)
	// Should not be reached if f, as a good Fatalizer, really dies during Fatal()
	panic(fmt.Sprint(args...))
}
