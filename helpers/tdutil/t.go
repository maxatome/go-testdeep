// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil

import (
	"reflect"
	"testing"
)

// Package tdutil allows to write unit tests for testdeep helpers and
// so provides some helpful functions.
//
// It is not intended to be used in tests outside go-testdeep and its
// helpers perimeter.

// T can be used in tests, to test testing.T behavior as it overrides
// Run() method.
type T struct {
	testing.T
}

// Run is a simplified version of testing.T.Run() method, without edge
// cases.
func (t *T) Run(name string, f func(*testing.T)) bool {
	f(&t.T)
	return !t.Failed()
}

// LogBuf is an ugly hack allowing to access internal testing.T log
// buffer. Keep cool, it is only used for internal unit tests.
func (t *T) LogBuf() string {
	return string(reflect.ValueOf(t.T).FieldByName("output").Bytes()) // nolint: govet
}
