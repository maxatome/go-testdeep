/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew_test

import (
	"fmt"

	"github.com/maxatome/go-testdeep/internal/spew"
)

type Flag int

const (
	flagOne Flag = iota
	flagTwo
)

var flagStrings = map[Flag]string{
	flagOne: "flagOne",
	flagTwo: "flagTwo",
}

func (f Flag) String() string {
	if s, ok := flagStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown flag (%d)", int(f))
}

type Bar struct {
	data uintptr
}

type Foo struct {
	unexportedField Bar
	ExportedField   map[any]any
}

// This example demonstrates how to use Sdump to dump variables to stdout.
func ExampleSdump() {
	// The following package level declarations are assumed for this example:
	/*
		type Flag int

		const (
			flagOne Flag = iota
			flagTwo
		)

		var flagStrings = map[Flag]string{
			flagOne: "flagOne",
			flagTwo: "flagTwo",
		}

		func (f Flag) String() string {
			if s, ok := flagStrings[f]; ok {
				return s
			}
			return fmt.Sprintf("Unknown flag (%d)", int(f))
		}

		type Bar struct {
			data uintptr
		}

		type Foo struct {
			unexportedField Bar
			ExportedField   map[any]any
		}
	*/

	// Setup some sample data structures for the example.
	bar := Bar{uintptr(0x1234)}
	s1 := Foo{bar, map[any]any{"one": true}}
	f := Flag(5)
	b := []byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}

	// Dump!
	fmt.Println(spew.Sdump(s1))
	fmt.Println(spew.Sdump(f))
	fmt.Println(spew.Sdump(b))

	// Output:
	// (spew_test.Foo) {
	//  unexportedField: (spew_test.Bar) {
	//   data: (uintptr) 0x1234
	//  },
	//  ExportedField: (map[interface {}]interface {}) (len=1) {
	//   (string) (len=3) "one": (bool) true
	//  }
	// }
	// (spew_test.Flag) Unknown flag (5)
	// ([]uint8) (len=34) {
	//  00000000  11 12 13 14 15 16 17 18  19 1a 1b 1c 1d 1e 1f 20  |............... |
	//  00000010  21 22 23 24 25 26 27 28  29 2a 2b 2c 2d 2e 2f 30  |!"#$%&'()*+,-./0|
	//  00000020  31 32                                             |12|
	// }
	//
}
