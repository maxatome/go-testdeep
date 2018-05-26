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

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"name": "foobar",
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"name": Re("^foo"),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"nameb": Re("^foo"),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"err": Re("error"),
	}))

	checkError(t, got,
		Struct(structPrivateFields{}, StructFields{
			"iface": Re("buffer"),
		}),
		expectedError{
			Message:  mustBe("bad type"),
			Path:     mustBe("DATA.iface"),
			Got:      mustBe("interface {}"),
			Expected: mustBe("string (convertible) OR fmt.Stringer OR error OR []uint8"),
		})

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"next": Struct(&structPrivateFields{}, StructFields{
			"iface": Re("buffer"),
		}),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"properties": []int{20, 22, 23, 21},
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"properties": ArrayEach(Between(20, 23)),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"byKey": MapEach(Struct(&privateValue{}, StructFields{
			"weight": Between(12, 34),
			"value":  Any(HasPrefix("t"), HasSuffix("e")),
		})),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"byKey": SuperMapOf(
			map[privateKey]*privateValue{
				{num: 3, name: "zip"}: {value: "ttse", weight: 34},
			},
			MapEntries{
				privateKey{num: 2, name: "bar"}: &privateValue{value: "tset", weight: 23},
			}),
	}))

	checkOK(t, got, Struct(structPrivateFields{}, StructFields{
		"birth":  TruncTime(d("2018-04-01T10:11:12Z"), time.Second),
		"birth2": TruncTime(MyTime(d("2018-03-01T09:08:07Z")), time.Second),
	}))

	checkError(t, got,
		Struct(structPrivateFields{}, StructFields{
			"next": Struct(&structPrivateFields{}, StructFields{
				"name":  "sub",
				"birth": Code(func(t time.Time) bool { return true }),
			}),
		}),
		expectedError{
			Message: mustBe("cannot compare unexported field"),
			Path:    mustBe("DATA.next.birth"),
			Summary: mustBe("use Code() on surrounding struct instead"),
		})
}
