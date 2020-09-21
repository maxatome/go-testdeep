// Copyright (c) 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types

import (
	"fmt"
	"reflect"
	"time"
)

var (
	Bool        = reflect.TypeOf(false)
	Interface   = reflect.TypeOf((*interface{})(nil)).Elem()
	FmtStringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	Error       = reflect.TypeOf((*error)(nil)).Elem()
	Time        = reflect.TypeOf(time.Time{})
	Int         = reflect.TypeOf(int(0))
	Uint8       = reflect.TypeOf(uint8(0))
	Rune        = reflect.TypeOf(rune(0))
	String      = reflect.TypeOf("")
)
