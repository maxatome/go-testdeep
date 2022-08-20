// Copyright (c) 2020-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types

import (
	"testing"
)

// Only for coverage...
func TestTypes(t *testing.T) {
	(TestDeepStamp{})._TestDeep()
	RawString("")._TestDeep()
	RawInt(0)._TestDeep()
	RecvNothing._TestDeep()
}
