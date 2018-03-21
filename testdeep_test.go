package testdeep

import (
	"testing"
)

func TestBase(t *testing.T) {
	//
	// Ignore
	if !EqDeeply("test", Ignore()) {
		t.Error("FAIL!")
	}
}
