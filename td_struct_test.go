package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestStruct(t *testing.T) {
	var gotStruct = MyStruct{
		MyStructMid: MyStructMid{
			MyStructBase: MyStructBase{
				ValBool: true,
			},
			ValStr: "foobar",
		},
		ValInt: 123,
	}

	//
	// Using pointer
	checkOK(t, &gotStruct,
		Struct(&MyStruct{}, StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
		}))

	checkError(t, 123, Struct(&MyStruct{}, StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("*testdeep_test.MyStruct"),
		})

	checkError(t, &MyStructBase{}, Struct(&MyStruct{}, StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*testdeep_test.MyStructBase"),
			Expected: mustContain("*testdeep_test.MyStruct"),
		})

	checkError(t, &gotStruct,
		Struct(&MyStruct{}, StructFields{
			"ValBool": false, // ← does not match
			"ValStr":  "foobar",
			"ValInt":  123,
		}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValBool"),
			Got:      mustContain("true"),
			Expected: mustContain("false"),
		})

	checkOK(t, &gotStruct,
		Struct(&MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, &gotStruct,
		Struct(&MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobax", // ← does not match
			},
			ValInt: 123,
		}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValStr"),
			Got:      mustContain("foobar"),
			Expected: mustContain("foobax"),
		})

	// Zero values
	checkOK(t, &MyStruct{},
		Struct(&MyStruct{}, StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	//
	// Without pointer
	checkOK(t, gotStruct,
		Struct(MyStruct{}, StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
		}))

	checkError(t, 123, Struct(MyStruct{}, StructFields{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("int"),
		Expected: mustContain("testdeep_test.MyStruct"),
	})

	checkError(t, gotStruct,
		Struct(MyStruct{}, StructFields{
			"ValBool": false, // ← does not match
			"ValStr":  "foobar",
			"ValInt":  123,
		}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValBool"),
			Got:      mustContain("true"),
			Expected: mustContain("false"),
		})

	checkOK(t, gotStruct,
		Struct(MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, gotStruct,
		Struct(MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobax", // ← does not match
			},
			ValInt: 123,
		}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValStr"),
			Got:      mustContain("foobar"),
			Expected: mustContain("foobax"),
		})

	// Zero values
	checkOK(t, MyStruct{},
		Struct(MyStruct{}, StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	//
	// Bad usage
	checkPanic(t, func() { Struct("test", nil) }, "usage: Struct")

	i := 12
	checkPanic(t, func() { Struct(&i, nil) }, "usage: Struct")

	checkPanic(t,
		func() { Struct(&MyStruct{}, StructFields{"UnknownField": 123}) },
		"struct testdeep_test.MyStruct has no field `UnknownField'")

	checkPanic(t,
		func() { Struct(&MyStruct{}, StructFields{"ValBool": 123}) },
		"type int of field expected value ValBool differs from struct one (bool)")

	checkPanic(t,
		func() {
			Struct(&MyStruct{
				MyStructMid: MyStructMid{
					MyStructBase: MyStructBase{
						ValBool: true,
					},
				},
			},
				StructFields{"ValBool": false})
		},
		"non zero field ValBool in model already exists in expectedFields")

	//
	// String
	equalStr(t, Struct(MyStruct{
		MyStructMid: MyStructMid{
			ValStr: "foobar",
		},
		ValInt: 123,
	},
		StructFields{
			"ValBool": false,
		}).String(),
		`Struct(testdeep_test.MyStruct{
  ValBool: (bool) false
  ValInt: (int) 123
  ValStr: (string) (len=6) "foobar"
})`)

	equalStr(t, Struct(&MyStruct{
		MyStructMid: MyStructMid{
			ValStr: "foobar",
		},
		ValInt: 123,
	},
		StructFields{
			"ValBool": false,
		}).String(),
		`Struct(*testdeep_test.MyStruct{
  ValBool: (bool) false
  ValInt: (int) 123
  ValStr: (string) (len=6) "foobar"
})`)

	equalStr(t, Struct(&MyStruct{}, StructFields{}).String(),
		`Struct(*testdeep_test.MyStruct{})`)
}
