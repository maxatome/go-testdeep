// Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

// NOTE: Due to the following build constraints, this file will only be compiled
// when both cgo is supported and "-tags testcgo" is added to the go test
// command line.  This means the cgo tests are only added (and hence run) when
// specifically requested.  This configuration is used because spew itself
// does not require cgo to run even though it does handle certain cgo types
// specially.  Rather than forcing all clients to require cgo and an external
// C compiler just to run the tests, this scheme makes them optional.

//go:build cgo && testcgo
// +build cgo,testcgo

package spew_test

import (
	"fmt"

	"github.com/maxatome/go-testdeep/internal/spew/testdata"
)

func addCgoSdumpTests() {
	// C char pointer.
	v := testdata.GetCgoCharPointer()
	nv := testdata.GetCgoNullCharPointer()
	pv := &v
	vcAddr := fmt.Sprintf("%p", v)
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "*testdata._Ctype_char"
	vs := "116"
	addSdumpTest("C char*", v, "("+vt+")("+vcAddr+")("+vs+")")
	addSdumpTest("C char* ptr", pv, "(*"+vt+")("+vAddr+"->"+vcAddr+")("+vs+")")
	addSdumpTest("C char* 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+"->"+vcAddr+")("+vs+")")
	addSdumpTest("C char* nil ptr", nv, "("+vt+")(<nil>)")

	// C char array.
	v2, v2l := testdata.GetCgoCharArray()
	v2Len := fmt.Sprintf("%d", v2l)
	v2t := "[6]testdata._Ctype_char"
	v2s := "(len=" + v2Len + ") " +
		"{\n 00000000  74 65 73 74 32 00                               " +
		"  |test2.|\n}"
	addSdumpTest("C char[]", v2, "("+v2t+") "+v2s)

	// C unsigned char array.
	v3, v3l := testdata.GetCgoUnsignedCharArray()
	v3Len := fmt.Sprintf("%d", v3l)
	v3t := "[6]testdata._Ctype_unsignedchar"
	v3t2 := "[6]testdata._Ctype_uchar"
	v3s := "(len=" + v3Len + ") " +
		"{\n 00000000  74 65 73 74 33 00                               " +
		"  |test3.|\n}"
	addSdumpTest("C unsigned char[]", v3, "("+v3t+") "+v3s+"\n", "("+v3t2+") "+v3s)

	// C signed char array.
	v4, v4l := testdata.GetCgoSignedCharArray()
	v4Len := fmt.Sprintf("%d", v4l)
	v4t := "[6]testdata._Ctype_schar"
	v4t2 := "testdata._Ctype_schar"
	v4s := "(len=" + v4Len + ") " +
		"{\n (" + v4t2 + ") 116,\n (" + v4t2 + ") 101,\n (" + v4t2 +
		") 115,\n (" + v4t2 + ") 116,\n (" + v4t2 + ") 52,\n (" + v4t2 +
		") 0\n}"
	addSdumpTest("C signed char[]", v4, "("+v4t+") "+v4s)

	// C uint8_t array.
	v5, v5l := testdata.GetCgoUint8tArray()
	v5Len := fmt.Sprintf("%d", v5l)
	v5t := "[6]testdata._Ctype_uint8_t"
	v5t2 := "[6]testdata._Ctype_uchar"
	v5s := "(len=" + v5Len + ") " +
		"{\n 00000000  74 65 73 74 35 00                               " +
		"  |test5.|\n}"
	addSdumpTest("C uint8_t[]", v5, "("+v5t+") "+v5s+"\n", "("+v5t2+") "+v5s)

	// C typedefed unsigned char array.
	v6, v6l := testdata.GetCgoTypdefedUnsignedCharArray()
	v6Len := fmt.Sprintf("%d", v6l)
	v6t := "[6]testdata._Ctype_custom_uchar_t"
	v6t2 := "[6]testdata._Ctype_uchar"
	v6s := "(len=" + v6Len + ") " +
		"{\n 00000000  74 65 73 74 36 00                               " +
		"  |test6.|\n}"
	addSdumpTest("C custom_uchar_t[]", v6, "("+v6t+") "+v6s+"\n", "("+v6t2+") "+v6s)
}
