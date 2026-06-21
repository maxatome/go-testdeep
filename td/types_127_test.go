// Copyright (c) 2026, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.27
// +build go1.27

package td_test

import "fmt"

func init() {
	// Starting go1.27, the first file of the stack trace is not
	// types_test.go but module_path/td/types_test.go. So as a "/" is
	// detected in file name, go-testdeep automatically add "This is how
	// we got here" paragraph.
	thisIsHowWeGotHere = func(line int) string {
		return fmt.Sprintf(`        This is how we got here:
        	TestSetlocation() td/types_test.go:%d
`, line)
	}
}
