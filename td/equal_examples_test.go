// Copyright (c) 2018-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"fmt"

	"github.com/maxatome/go-testdeep/td"
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

	if td.EqDeeply(got,
		td.Struct(&MyStruct{},
			td.StructFields{
				"Name":  td.Re("^Foo"),
				"Num":   td.Between(10, 20),
				"Items": td.ArrayEach(td.Between(3, 9)),
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

	err := td.EqDeeplyError(got,
		td.Struct(&MyStruct{},
			td.StructFields{
				"Name":  td.Re("^Foo"),
				"Num":   td.Between(10, 20),
				"Items": td.ArrayEach(td.Between(3, 8)),
			}))
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// DATA.Items[2]: values differ
	// 	     got: 9
	// 	expected: 3 ≤ got ≤ 8
	// [under operator Between at example.go:18]
}
