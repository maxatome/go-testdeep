// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package nocolor is only intended to easily disable coloring output
// of failure reports, typically useful in golang playground.
//
// Simply import it, and nothing else:
//
//	import _ "github.com/maxatome/go-testdeep/helpers/nocolor"
package nocolor

import "os"

func init() {
	os.Setenv("TESTDEEP_COLOR", "off")
}
