// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types

import (
	"strconv"
)

type TestDeepStringer interface {
	_TestDeep()
	String() string
}

type TestDeepStamp struct{}

func (_ TestDeepStamp) _TestDeep() {}

// Implements types.TestDeepStringer
type RawString string

func (s RawString) _TestDeep() {}

func (s RawString) String() string {
	return string(s)
}

// Implements types.TestDeepStringer
type RawInt int

func (i RawInt) _TestDeep() {}

func (i RawInt) String() string {
	return strconv.Itoa(int(i))
}
