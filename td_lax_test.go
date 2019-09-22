// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestLax(t *testing.T) {
	checkOK(t, int64(1234), testdeep.Lax(1234))

	type MyInt int32
	checkOK(t, int64(123), testdeep.Lax(MyInt(123)))
	checkOK(t, MyInt(123), testdeep.Lax(int64(123)))

	type gotStruct struct {
		name string
		age  int
	}
	type expectedStruct struct {
		name string
		age  int
	}
	checkOK(t,
		gotStruct{
			name: "bob",
			age:  42,
		},
		testdeep.Lax(expectedStruct{
			name: "bob",
			age:  42,
		}))
	checkOK(t,
		&gotStruct{
			name: "bob",
			age:  42,
		},
		testdeep.Lax(&expectedStruct{
			name: "bob",
			age:  42,
		}))

	checkError(t, int64(123), testdeep.Between(120, 125),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkOK(t, int64(123), testdeep.Lax(testdeep.Between(120, 125)))

	// nil cases
	checkOK(t, nil, testdeep.Lax(nil))
	checkOK(t, (*gotStruct)(nil), testdeep.Lax((*expectedStruct)(nil)))
	checkOK(t, (*gotStruct)(nil), testdeep.Lax(nil))
	checkOK(t, (chan int)(nil), testdeep.Lax(nil))
	checkOK(t, (func())(nil), testdeep.Lax(nil))
	checkOK(t, (map[int]int)(nil), testdeep.Lax(nil))
	checkOK(t, ([]int)(nil), testdeep.Lax(nil))

	//
	// String
	test.EqualStr(t, testdeep.Lax(6).String(), "Lax(6)")
}

func TestLaxTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Lax(nil), nil)

	type MyBool bool
	equalTypes(t, testdeep.Lax(MyBool(false)), false)
	equalTypes(t, testdeep.Lax(0), int64(0))
	equalTypes(t, testdeep.Lax(uint8(0)), uint64(0))
	equalTypes(t, testdeep.Lax(float32(0)), float64(0))
	equalTypes(t, testdeep.Lax(complex64(complex(1, 1))), complex128(complex(1, 1)))
	type MyString string
	equalTypes(t, testdeep.Lax(MyString("")), "")

	type MyBytes []byte
	equalTypes(t, testdeep.Lax([]byte{}), []byte{})
	equalTypes(t, testdeep.Lax(MyBytes{}), MyBytes{})

	// Another TestDeep operator delegation
	equalTypes(t, testdeep.Lax(testdeep.Struct(MyStruct{}, nil)), MyStruct{})
	equalTypes(t, testdeep.Lax(testdeep.Any(1, 1.2)), nil)
}
