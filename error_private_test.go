package testdeep

import (
	"testing"
)

func TestErrorPrivate(t *testing.T) {
	if booleanError.Error() != "" {
		t.Errorf("booleanError should stringify to empty string, not `%s'",
			booleanError.Error())
	}
}
