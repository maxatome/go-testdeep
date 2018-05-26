// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

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
