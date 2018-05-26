// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

// tdSmuggler is the base class of all smuggler TestDeep operators.
type tdSmuggler struct {
	Base
	expectedValue reflect.Value
	isTestDeeper  bool
}

func newSmuggler(val interface{}) (ret tdSmuggler) {
	ret.Base = NewBase(4)

	// Initializes only if TestDeep operator. Other cases are specific.
	if _, ok := val.(TestDeep); ok {
		ret.expectedValue = reflect.ValueOf(val)
		ret.isTestDeeper = true
	}
	return
}
