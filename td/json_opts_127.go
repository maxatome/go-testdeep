// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
)

type jsonv2Options = json.Options

func joinOptions(a, b jsonv2Options) jsonv2Options {
	return json.JoinOptions(a, b)
}

func jsonMarshal(v any, opts jsonv2Options) ([]byte, error) {
	if opts == nil {
		opts = jsonv1.DefaultOptionsV1()
	}
	return json.Marshal(v, opts)
}

func jsonUnmarshal(data []byte, v any, opts jsonv2Options) error {
	if opts == nil {
		opts = jsonv1.DefaultOptionsV1()
	}
	return json.Unmarshal(data, v, opts)
}
