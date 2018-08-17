// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
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
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			"Ptr":     nil,
		}))

	checkOK(t, &gotStruct,
		testdeep.Struct((*MyStruct)(nil), testdeep.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
			"Ptr":     nil,
		}))

	checkError(t, 123,
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("*testdeep_test.MyStruct"),
		})

	checkError(t, &MyStructBase{},
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("*testdeep_test.MyStructBase"),
			Expected: mustContain("*testdeep_test.MyStruct"),
		})

	checkError(t, &gotStruct,
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{
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
		testdeep.Struct(&MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, &gotStruct,
		testdeep.Struct(&MyStruct{
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
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	// nil cases
	checkError(t, nil, testdeep.Struct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("*testdeep_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), testdeep.Struct(&MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("non-nil"),
		})

	//
	// Without pointer
	checkOK(t, gotStruct,
		testdeep.Struct(MyStruct{}, testdeep.StructFields{
			"ValBool": true,
			"ValStr":  "foobar",
			"ValInt":  123,
		}))

	checkError(t, 123, testdeep.Struct(MyStruct{}, testdeep.StructFields{}),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustContain("int"),
			Expected: mustContain("testdeep_test.MyStruct"),
		})

	checkError(t, gotStruct,
		testdeep.Struct(MyStruct{}, testdeep.StructFields{
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
		testdeep.Struct(MyStruct{
			MyStructMid: MyStructMid{
				MyStructBase: MyStructBase{
					ValBool: true,
				},
				ValStr: "foobar",
			},
			ValInt: 123,
		}, nil))

	checkError(t, gotStruct,
		testdeep.Struct(MyStruct{
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
		testdeep.Struct(MyStruct{}, testdeep.StructFields{
			"ValBool": false,
			"ValStr":  "",
			"ValInt":  0,
		}))

	// nil cases
	checkError(t, nil, testdeep.Struct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("values differ"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustContain("testdeep_test.MyStruct"),
		})

	checkError(t, (*MyStruct)(nil), testdeep.Struct(MyStruct{}, nil),
		expectedError{
			Message:  mustBe("type mismatch"),
			Path:     mustBe("DATA"),
			Got:      mustBe("*testdeep_test.MyStruct"),
			Expected: mustBe("testdeep_test.MyStruct"),
		})

	//
	// Bad usage
	test.CheckPanic(t, func() { testdeep.Struct("test", nil) }, "usage: Struct")

	i := 12
	test.CheckPanic(t, func() { testdeep.Struct(&i, nil) }, "usage: Struct")

	test.CheckPanic(t,
		func() { testdeep.Struct(&MyStruct{}, testdeep.StructFields{"UnknownField": 123}) },
		"struct testdeep_test.MyStruct has no field `UnknownField'")

	test.CheckPanic(t,
		func() { testdeep.Struct(&MyStruct{}, testdeep.StructFields{"ValBool": 123}) },
		"type int of field expected value ValBool differs from struct one (bool)")

	test.CheckPanic(t,
		func() { testdeep.Struct(&MyStruct{}, testdeep.StructFields{"ValBool": nil}) },
		"expected value of field ValBool cannot be nil as it is a bool")

	test.CheckPanic(t,
		func() {
			testdeep.Struct(&MyStruct{
				MyStructMid: MyStructMid{
					MyStructBase: MyStructBase{
						ValBool: true,
					},
				},
			},
				testdeep.StructFields{"ValBool": false})
		},
		"non zero field ValBool in model already exists in expectedFields")

	//
	// String
	test.EqualStr(t,
		testdeep.Struct(MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			testdeep.StructFields{
				"ValBool": false,
			}).String(),
		`Struct(testdeep_test.MyStruct{
  ValBool: (bool) false
  ValInt: (int) 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		testdeep.Struct(&MyStruct{
			MyStructMid: MyStructMid{
				ValStr: "foobar",
			},
			ValInt: 123,
		},
			testdeep.StructFields{
				"ValBool": false,
			}).String(),
		`Struct(*testdeep_test.MyStruct{
  ValBool: (bool) false
  ValInt: (int) 123
  ValStr: "foobar"
})`)

	test.EqualStr(t,
		testdeep.Struct(&MyStruct{}, testdeep.StructFields{}).String(),
		`Struct(*testdeep_test.MyStruct{})`)
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
		err:        errors.New("the error!"),
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
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"name": "foobar",
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"name": testdeep.Re("^foo"),
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"nameb": testdeep.Re("^foo"),
		}))

	checkOKOrPanicIfUnsafeDisabled(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"err": testdeep.Re("error"),
		}))

	checkError(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"iface": testdeep.Re("buffer"),
		}),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA.iface"),
			Got:      mustBe("interface {}"),
			Expected: mustBe("string (convertible) OR fmt.Stringer OR error OR []uint8"),
		})

	checkOKOrPanicIfUnsafeDisabled(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"next": testdeep.Struct(&structPrivateFields{}, testdeep.StructFields{
				"iface": testdeep.Re("buffer"),
			}),
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"properties": []int{20, 22, 23, 21},
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"properties": testdeep.ArrayEach(testdeep.Between(20, 23)),
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"byKey": testdeep.MapEach(testdeep.Struct(&privateValue{}, testdeep.StructFields{
				"weight": testdeep.Between(12, 34),
				"value":  testdeep.Any(testdeep.HasPrefix("t"), testdeep.HasSuffix("e")),
			})),
		}))

	checkOK(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"byKey": testdeep.SuperMapOf(
				map[privateKey]*privateValue{
					{num: 3, name: "zip"}: {value: "ttse", weight: 34},
				},
				testdeep.MapEntries{
					privateKey{num: 2, name: "bar"}: &privateValue{value: "tset", weight: 23},
				}),
		}))

	expected := testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
		"birth":  testdeep.TruncTime(d("2018-04-01T10:11:12Z"), time.Second),
		"birth2": testdeep.TruncTime(MyTime(d("2018-03-01T09:08:07Z")), time.Second),
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
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"next": testdeep.Struct(&structPrivateFields{}, testdeep.StructFields{
				"name":  "sub",
				"birth": testdeep.Code(func(t time.Time) bool { return true }),
			}),
		}),
		expectedError{
			Message: mustBe("cannot compare unexported field"),
			Path:    mustBe("DATA.next.birth"),
			Summary: mustBe("use Code() on surrounding struct instead"),
		})

	checkError(t, got,
		testdeep.Struct(structPrivateFields{}, testdeep.StructFields{
			"next": testdeep.Struct(&structPrivateFields{}, testdeep.StructFields{
				"name": "sub",
				"birth": testdeep.Smuggle(
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
	equalTypes(t, testdeep.Struct(MyStruct{}, nil), MyStruct{})
	equalTypes(t, testdeep.Struct(&MyStruct{}, nil), &MyStruct{})
}
