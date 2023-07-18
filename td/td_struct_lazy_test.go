// Copyright (c) 2023, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestStructLazy(t *testing.T) {
	got := struct {
		ValInt int
		ValStr string
	}{0, "foobar"}

	t.Run("Struct OK", func(t *testing.T) {
		got.ValInt = 123

		checkOK(t, got, td.Struct(nil, td.StructFields{
			"ValStr": "foobar",
		}))

		checkOK(t, &got, td.Struct(nil, td.StructFields{
			"ValInt": 123,
			"ValStr": "foobar",
		}))

		checkOK(t, &got, td.Struct(nil, td.StructFields{"=Val*": td.NotZero()}))
	})

	t.Run("SStruct OK", func(t *testing.T) {
		got.ValInt = 0

		checkOK(t, got, td.SStruct(nil, td.StructFields{
			"ValStr": "foobar",
		}))

		checkOK(t, &got, td.SStruct(nil, td.StructFields{
			"ValInt": 0,
			"ValStr": "foobar",
		}))

		got.ValInt = 123
		checkOK(t, &got, td.SStruct(nil, td.StructFields{"=Val*": td.NotZero()}))
	})

	got.ValInt = 666

	ops := []struct {
		name string
		new  func(any, ...td.StructFields) td.TestDeep
	}{
		{"Struct", td.Struct},
		{"SStruct", td.SStruct},
	}
	for _, op := range ops {
		t.Run(op.name+" errors", func(t *testing.T) {
			under := mustContain("under operator " + op.name + " at td_struct_lazy_test.go:")
			badUsage := mustBe("bad usage of " + op.name + " operator")

			checkError(t, got, op.new(nil, td.StructFields{"Zip": 345}),
				expectedError{
					Message: badUsage,
					Path:    mustBe("DATA"),
					Summary: mustBe(`struct { ValInt int; ValStr string } has no field "Zip"`),
					Under:   under,
				})

			checkError(t, got, op.new(nil, td.StructFields{">\tZip": 345}),
				expectedError{
					Message: badUsage,
					Path:    mustBe("DATA"),
					Summary: mustBe(`struct { ValInt int; ValStr string } has no field "Zip" (from ">\tZip")`),
					Under:   under,
				})

			checkError(t, got, op.new(nil, td.StructFields{"ValInt": "zip"}),
				expectedError{
					Message: badUsage,
					Path:    mustBe("DATA"),
					Summary: mustBe("type string of field expected value ValInt differs from struct one (int)"),
				})

			checkError(t, 123,
				op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("bad kind"),
					Path:     mustBe("DATA"),
					Got:      mustContain("int"),
					Expected: mustContain("struct OR *struct"),
					Under:    under,
				})

			n := 123
			checkError(t, &n,
				op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("bad kind"),
					Path:     mustBe("DATA"),
					Got:      mustContain("*int"),
					Expected: mustContain("struct OR *struct"),
					Under:    under,
				})

			type myInt int
			checkError(t, myInt(123),
				op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("bad kind"),
					Path:     mustBe("DATA"),
					Got:      mustContain("int (td_test.myInt type)"),
					Expected: mustContain("struct OR *struct"),
					Under:    under,
				})

			mi := myInt(123)
			checkError(t, &mi,
				op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("bad kind"),
					Path:     mustBe("DATA"),
					Got:      mustContain("*int (*td_test.myInt type)"),
					Expected: mustContain("struct OR *struct"),
					Under:    under,
				})

			checkError(t, nil, op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("DATA"),
					Got:      mustBe("nil"),
					Expected: mustContain("Struct(<any struct type>{})"),
					Under:    under,
				})

			checkError(t, (*struct{ x int })(nil), op.new(nil, td.StructFields{}),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("DATA"),
					Got:      mustBe("(*struct { x int })(<nil>)"),
					Expected: mustContain("non-nil"),
					Under:    under,
				})
		})

		t.Run(op.name+" String", func(t *testing.T) {
			test.EqualStr(t, op.new(nil).String(), op.name+`(<any struct type>{})`)
			test.EqualStr(t,
				op.new(nil,
					td.StructFields{
						"ValBool":     false,
						"= Val*":      td.NotZero(),
						"> Foo":       12,
						"=~Bar[6-9]$": "zip",
					}).String(),
				op.name+`(<any struct type>{
  = Val*:      NotZero()
  =~Bar[6-9]$: "zip"
  > Foo:       12
  ValBool:     false
})`)
		})
	}
}

func TestStructLazyTypeBehind(t *testing.T) {
	equalTypes(t, td.Struct(nil, nil), nil)
	equalTypes(t, td.SStruct(nil, nil), nil)
}
