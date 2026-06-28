// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build !go1.27
// +build !go1.27

package td

import (
	"encoding/json"
)

// Before go1.27 this type cannot be reached from outside.
type jsonv2Options interface{ private() }

func joinOptions(a, b jsonv2Options) jsonv2Options {
	return nil
}

func jsonMarshal(v any, opts jsonv2Options) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v any, opts jsonv2Options) error {
	return json.Unmarshal(data, v)
}
