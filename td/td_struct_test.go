// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
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
		td.Struct(&MyStruct{}, td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			"Ptr":     nil,
		}))

	checkOK(t, &gotStruct,
		td.Struct((*MyStruct)(nil), td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			"Ptr":     nil,
		}))

	checkError(t, 123,
		td.Struct(&MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, &MyStructBase{},
		td.Struct(&MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*td_test.MyStructBase"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, &gotStruct,
		td.Struct(&MyStruct{}, td.StructFields{
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
		td.Struct(&MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, &gotStruct,
		td.Struct(&MyStruct{
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
		td.Struct(&MyStruct{}, td.StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	// nil cases
	checkError(t, nil, td.Struct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), td.Struct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("non-nil"),
		})

	//
	// Without pointer
	checkOK(t, gotStruct,
		td.Struct(MyStruct{}, td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
		}))

	checkError(t, 123, td.Struct(MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("td_test.MyStruct"),
		})

	checkError(t, gotStruct,
		td.Struct(MyStruct{}, td.StructFields{
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
		td.Struct(MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, gotStruct,
		td.Struct(MyStruct{
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
		td.Struct(MyStruct{}, td.StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	// nil cases
	checkError(t, nil, td.Struct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("td_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), td.Struct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MyStruct"),
			Expected: mustBe("td_test.MyStruct"),
		})

	//
	// Be lax...
	type Struct1 struct {
		name string
		age  int
	}
	type Struct2 struct {
		name string
		age  int
	}

	// Without Lax → error
	checkError(t,
		Struct1{name: "Bob", age: 42},
		td.Struct(Struct2{name: "Bob", age: 42}, nil),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	// With Lax → OK
	checkOK(t,
		Struct1{name: "Bob", age: 42},
		td.Lax(td.Struct(Struct2{name: "Bob", age: 42}, nil)))

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.Struct("test", nil) },
		"usage: Struct(STRUCT|&STRUCT, EXPECTED_FIELDS), but received string as 1st parameter")

	i := 12
	dark.CheckFatalizerBarrierErr(t, func() { td.Struct(&i, nil) },
		"usage: Struct(STRUCT|&STRUCT, EXPECTED_FIELDS), but received *int (ptr) as 1st parameter")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Struct(&MyStruct{}, td.StructFields{"UnknownField": 123}) },
		"Struct(): struct td_test.MyStruct has no field `UnknownField'")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Struct(&MyStruct{}, td.StructFields{"ValBool": 123}) },
		"Struct(): type int of field expected value ValBool differs from struct one (bool)")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.Struct(&MyStruct{}, td.StructFields{"ValBool": nil}) },
		"Struct(): expected value of field ValBool cannot be nil as it is a bool")

	dark.CheckFatalizerBarrierErr(t,
		func() {
			td.Struct(&MyStruct{
				MyStructMid: MyStructMid{
					MyStructBase: MyStructBase{
						ValBool: true,
					},
				},
			},
				td.StructFields{"ValBool": false})
		},
		"Struct(): non zero field ValBool in model already exists in expectedFields")

	//
	// String
	test.EqualStr(t,
		td.Struct(MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			td.StructFields{
				"ValBool": false,
			}).String(),
		`Struct(td_test.MyStruct{
  ValBool: false
  ValInt: 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		td.Struct(&MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			td.StructFields{
				"ValBool": false,
			}).String(),
		`Struct(*td_test.MyStruct{
  ValBool: false
  ValInt: 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		td.Struct(&MyStruct{}, td.StructFields{}).String(),
		`Struct(*td_test.MyStruct{})`)
}

func TestStructPrivateFields(t *testing.T) {
	type privateKey struct {
		num  int
		name string
	}

	type privateValue struct {
		value  string
		weight int
	}

	type MyTime time.Time

	type structPrivateFields struct {
		byKey      map[privateKey]*privateValue
		name       string
		nameb      []byte
		err        error
		iface      interface{}
		properties []int
		birth      time.Time
		birth2     MyTime
		next       *structPrivateFields
	}

	d := func(rfc3339Date string) (ret time.Time) {
		var err error
		ret, err = time.Parse(time.RFC3339Nano, rfc3339Date)
		if err != nil {
			panic(err)
		}
		return
	}

	got := structPrivateFields{
		byKey: map[privateKey]*privateValue{
			{num: 1, name: "foo"}: {value: "test", weight: 12},
			{num: 2, name: "bar"}: {value: "tset", weight: 23},
			{num: 3, name: "zip"}: {value: "ttse", weight: 34},
		},
		name:       "foobar",
		nameb:      []byte("foobar"),
		err:        errors.New("the error"),
		iface:      1234,
		properties: []int{20, 22, 23, 21},
		birth:      d("2018-04-01T10:11:12.123456789Z"),
		birth2:     MyTime(d("2018-03-01T09:08:07.987654321Z")),
		next: &structPrivateFields{
			byKey:  map[privateKey]*privateValue{},
			name:   "sub",
			iface:  bytes.NewBufferString("buffer!"),
			birth:  d("2018-04-02T10:11:12.123456789Z"),
			birth2: MyTime(d("2018-03-02T09:08:07.987654321Z")),
		},
	}

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"name": "foobar",
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"name": td.Re("^foo"),
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"nameb": td.Re("^foo"),
		}))

	checkOKOrPanicIfUnsafeDisabled(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"err": td.Re("error"),
		}))

	checkError(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"iface": td.Re("buffer"),
		}),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA.iface"),
			Got:      mustBe("int"),
			Expected: mustBe("string (convertible) OR fmt.Stringer OR error OR []uint8"),
		})

	checkOKOrPanicIfUnsafeDisabled(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"next": td.Struct(&structPrivateFields{}, td.StructFields{
				"iface": td.Re("buffer"),
			}),
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"properties": []int{20, 22, 23, 21},
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"properties": td.ArrayEach(td.Between(20, 23)),
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"byKey": td.MapEach(td.Struct(&privateValue{}, td.StructFields{
				"weight": td.Between(12, 34),
				"value":  td.Any(td.HasPrefix("t"), td.HasSuffix("e")),
			})),
		}))

	checkOK(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"byKey": td.SuperMapOf(
				map[privateKey]*privateValue{
					{num: 3, name: "zip"}: {value: "ttse", weight: 34},
				},
				td.MapEntries{
					privateKey{num: 2, name: "bar"}: &privateValue{value: "tset", weight: 23},
				}),
		}))

	expected := td.Struct(structPrivateFields{}, td.StructFields{
		"birth":  td.TruncTime(d("2018-04-01T10:11:12Z"), time.Second),
		"birth2": td.TruncTime(MyTime(d("2018-03-01T09:08:07Z")), time.Second),
	})
	if !dark.UnsafeDisabled {
		checkOK(t, got, expected)
	} else {
		checkError(t, got, expected,
			expectedError{
				Message: mustBe("cannot compare"),
				Path:    mustBe("DATA.birth"),
				Summary: mustBe("unexported field that cannot be overridden"),
			})
	}

	checkError(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"next": td.Struct(&structPrivateFields{}, td.StructFields{
				"name":  "sub",
				"birth": td.Code(func(t time.Time) bool { return true }),
			}),
		}),
		expectedError{
			Message: mustBe("cannot compare unexported field"),
			Path:    mustBe("DATA.next.birth"),
			Summary: mustBe("use Code() on surrounding struct instead"),
		})

	checkError(t, got,
		td.Struct(structPrivateFields{}, td.StructFields{
			"next": td.Struct(&structPrivateFields{}, td.StructFields{
				"name": "sub",
				"birth": td.Smuggle(
					func(t time.Time) string { return t.String() },
					"2018-04-01T10:11:12.123456789Z"),
			}),
		}),
		expectedError{
			Message: mustBe("cannot smuggle unexported field"),
			Path:    mustBe("DATA.next.birth"),
			Summary: mustBe("work on surrounding struct instead"),
		})
}

func TestStructTypeBehind(t *testing.T) {
	equalTypes(t, td.Struct(MyStruct{}, nil), MyStruct{})
	equalTypes(t, td.Struct(&MyStruct{}, nil), &MyStruct{})
}

func TestSStruct(t *testing.T) {
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
		td.SStruct(&MyStruct{}, td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			// nil Ptr
		}))

	checkOK(t, &gotStruct,
		td.SStruct((*MyStruct)(nil), td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			// nil Ptr
		}))

	checkError(t, 123,
		td.SStruct(&MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, &MyStructBase{},
		td.SStruct(&MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*td_test.MyStructBase"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, &gotStruct,
		td.SStruct(&MyStruct{}, td.StructFields{
			// ValBool false ← does not match
			"ValStr": "foobar",
			"ValInt": 123,
		}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValBool"),
			Got:      mustContain("true"),
			Expected: mustContain("false"),
		})

	checkOK(t, &gotStruct,
		td.SStruct(&MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, &gotStruct,
		td.SStruct(&MyStruct{
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
	checkOK(t, &MyStruct{}, td.SStruct(&MyStruct{}, nil))
	checkOK(t, &MyStruct{}, td.SStruct(&MyStruct{}, td.StructFields{}))

	// nil cases
	checkError(t, nil, td.SStruct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("*td_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), td.SStruct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("non-nil"),
		})

	//
	// Without pointer
	checkOK(t, gotStruct,
		td.SStruct(MyStruct{}, td.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
		}))

	checkError(t, 123, td.SStruct(MyStruct{}, td.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("td_test.MyStruct"),
		})

	checkError(t, gotStruct,
		td.SStruct(MyStruct{}, td.StructFields{
			// "ValBool" false ← does not match
			"ValStr": "foobar",
			"ValInt": 123,
		}),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA.ValBool"),
			Got:      mustContain("true"),
			Expected: mustContain("false"),
		})

	checkOK(t, gotStruct,
		td.SStruct(MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, gotStruct,
		td.SStruct(MyStruct{
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
	checkOK(t, MyStruct{}, td.Struct(MyStruct{}, td.StructFields{}))
	checkOK(t, MyStruct{}, td.Struct(MyStruct{}, nil))

	// nil cases
	checkError(t, nil, td.SStruct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("td_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), td.SStruct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*td_test.MyStruct"),
			Expected: mustBe("td_test.MyStruct"),
		})

	//
	// Be lax...
	type Struct1 struct {
		name string
		age  int
	}
	type Struct2 struct {
		name string
		age  int
	}

	// Without Lax → error
	checkError(t,
		Struct1{name: "Bob", age: 42},
		td.SStruct(Struct2{name: "Bob", age: 42}, nil),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	// With Lax → OK
	checkOK(t,
		Struct1{name: "Bob", age: 42},
		td.Lax(td.SStruct(Struct2{name: "Bob", age: 42}, nil)))

	//
	// Bad usage
	dark.CheckFatalizerBarrierErr(t, func() { td.SStruct("test", nil) },
		"usage: SStruct(STRUCT|&STRUCT, EXPECTED_FIELDS), but received string as 1st parameter")

	i := 12
	dark.CheckFatalizerBarrierErr(t, func() { td.SStruct(&i, nil) },
		"usage: SStruct(STRUCT|&STRUCT, EXPECTED_FIELDS), but received *int (ptr) as 1st parameter")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.SStruct(&MyStruct{}, td.StructFields{"UnknownField": 123}) },
		"SStruct(): struct td_test.MyStruct has no field `UnknownField'")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.SStruct(&MyStruct{}, td.StructFields{"ValBool": 123}) },
		"SStruct(): type int of field expected value ValBool differs from struct one (bool)")

	dark.CheckFatalizerBarrierErr(t,
		func() { td.SStruct(&MyStruct{}, td.StructFields{"ValBool": nil}) },
		"SStruct(): expected value of field ValBool cannot be nil as it is a bool")

	dark.CheckFatalizerBarrierErr(t,
		func() {
			td.SStruct(&MyStruct{
				MyStructMid: MyStructMid{
					MyStructBase: MyStructBase{
						ValBool: true,
					},
				},
			},
				td.StructFields{"ValBool": false})
		},
		"SStruct(): non zero field ValBool in model already exists in expectedFields")

	//
	// String
	test.EqualStr(t,
		td.SStruct(MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			td.StructFields{
				"ValBool": false,
			}).String(),
		`SStruct(td_test.MyStruct{
  Ptr: (*int)(<nil>)
  ValBool: false
  ValInt: 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		td.SStruct(&MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			td.StructFields{
				"ValBool": false,
			}).String(),
		`SStruct(*td_test.MyStruct{
  Ptr: (*int)(<nil>)
  ValBool: false
  ValInt: 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		td.SStruct(&MyStruct{}, td.StructFields{}).String(),
		`SStruct(*td_test.MyStruct{
  Ptr: (*int)(<nil>)
  ValBool: false
  ValInt: 0
  ValStr: ""
})`)
}

func TestSStructTypeBehind(t *testing.T) {
	equalTypes(t, td.SStruct(MyStruct{}, nil), MyStruct{})
	equalTypes(t, td.SStruct(&MyStruct{}, nil), &MyStruct{})
}
