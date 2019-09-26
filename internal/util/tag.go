// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util

import (
	"errors"
	"unicode"
)

var (
	ErrTagEmpty   = errors.New("A tag cannot be empty")
	ErrTagInvalid = errors.New("Invalid tag, should match (Letter|_)(Letter|_|Number)*")
)

func CheckTag(tag string) error {
	if tag == "" {
		return ErrTagEmpty
	}

	for i, r := range tag {
		if !(unicode.IsLetter(r) || r == '_' || (i > 0 && unicode.IsNumber(r))) {
			return ErrTagInvalid
		}
	}

	return nil
}
