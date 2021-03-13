// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build js appengine safe disableunsafe race

package dark

import (
	"errors"
)

func GetFatalizer() Fatalizer {
	return FatalPanic("")
}

func FatalizerBarrier(fn func()) (err error) {
	defer func() {
		if x := recover(); x != nil {
			s, ok := x.(FatalPanic)
			if !ok {
				panic(x) // rethrow
			}
			err = errors.New(string(s))
		}
	}()
	fn()
	return
}
