// Copyright (c) 2019, Maxime Soulé
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

func TestLax(t *testing.T) {
	checkOK(t, int64(1234), td.Lax(1234))

	type MyInt int32
	checkOK(t, int64(123), td.Lax(MyInt(123)))
	checkOK(t, MyInt(123), td.Lax(int64(123)))

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
		td.Lax(expectedStruct{
			name: "bob",
			age:  42,
		}))
	checkOK(t,
		&gotStruct{
			name: "bob",
			age:  42,
		},
		td.Lax(&expectedStruct{
			name: "bob",
			age:  42,
		}))

	checkError(t, int64(123), td.Between(120, 125),
		expectedError{
			Message: mustBe("type mismatch"),
		})
	checkOK(t, int64(123), td.Lax(td.Between(120, 125)))

	// nil cases
	checkOK(t, nil, td.Lax(nil))
	checkOK(t, (*gotStruct)(nil), td.Lax((*expectedStruct)(nil)))
	checkOK(t, (*gotStruct)(nil), td.Lax(nil))
	checkOK(t, (chan int)(nil), td.Lax(nil))
	checkOK(t, (func())(nil), td.Lax(nil))
	checkOK(t, (map[int]int)(nil), td.Lax(nil))
	checkOK(t, ([]int)(nil), td.Lax(nil))

	//
	// String
	test.EqualStr(t, td.Lax(6).String(), "Lax(6)")
}

func TestLaxTypeBehind(t *testing.T) {
	equalTypes(t, td.Lax(nil), nil)

	type MyBool bool
	equalTypes(t, td.Lax(MyBool(false)), false)
	equalTypes(t, td.Lax(0), int64(0))
	equalTypes(t, td.Lax(uint8(0)), uint64(0))
	equalTypes(t, td.Lax(float32(0)), float64(0))
	equalTypes(t, td.Lax(complex64(complex(1, 1))), complex128(complex(1, 1)))
	type MyString string
	equalTypes(t, td.Lax(MyString("")), "")

	type MyBytes []byte
	equalTypes(t, td.Lax([]byte{}), []byte{})
	equalTypes(t, td.Lax(MyBytes{}), MyBytes{})

	// Another TestDeep operator delegation
	equalTypes(t, td.Lax(td.Struct(MyStruct{}, nil)), MyStruct{})
	equalTypes(t, td.Lax(td.Any(1, 1.2)), nil)
}
