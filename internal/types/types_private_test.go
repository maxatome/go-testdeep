// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types

import (
	"testing"
)

func TestTypes(t *testing.T) {
	s := RawString("foo")
	if str := s.String(); str != "foo" {
		t.Errorf("Very weird, got %s", str)
	}

	i := RawInt(42)
	if str := i.String(); str != "42" {
		t.Errorf("Very weird, got %s", str)
	}

	// Only for coverage...
	(TestDeepStamp{})._TestDeep()
	RawString("")._TestDeep()
	RawInt(0)._TestDeep()
}
