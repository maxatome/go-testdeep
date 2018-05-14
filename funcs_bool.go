package testdeep

import (
	"testing"
)

// CmpTrue is a shortcut for:
//
//   CmpDeeply(t, got, true, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpTrue(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, true, args...)
}

// CmpFalse is a shortcut for:
//
//   CmpDeeply(t, got, false, args...)
//
// Returns true if the test is OK, false if it fails.
func CmpFalse(t *testing.T, got interface{}, args ...interface{}) bool {
	return CmpDeeply(t, got, false, args...)
}
