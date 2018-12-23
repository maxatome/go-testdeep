// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
)

func ExampleEqDeeply() {
	type MyStruct struct {
		Name  string
		Num   int
		Items []int
	}

	got := &MyStruct{
		Name:  "Foobar",
		Num:   12,
		Items: []int{4, 5, 9, 3, 8},
	}

	if EqDeeply(got,
		Struct(&MyStruct{},
			StructFields{
				"Name":  Re("^Foo"),
				"Num":   Between(10, 20),
				"Items": ArrayEach(Between(3, 9)),
			})) {
		fmt.Println("Match!")
	} else {
		fmt.Println("NO!")
	}

	// Output:
	// Match!
}

func ExampleEqDeeplyError() {
//line /testdeep/example.go:1
	type MyStruct struct {
		Name  string
		Num   int
		Items []int
	}

	got := &MyStruct{
		Name:  "Foobar",
		Num:   12,
		Items: []int{4, 5, 9, 3, 8},
	}

	err := EqDeeplyError(got,
		Struct(&MyStruct{},
			StructFields{
				"Name":  Re("^Foo"),
				"Num":   Between(10, 20),
				"Items": ArrayEach(Between(3, 8)),
			}))
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// DATA.Items[2]: values differ
	// 	     got: 9
	// 	expected: 3 ≤ got ≤ 8
	// [under TestDeep operator Between at example.go:18]
}
