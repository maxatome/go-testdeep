/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Test Summary:
NOTE: For each test, a nil pointer, a single pointer and double pointer to the
base test element are also tested to ensure proper indirection across all types.

- Max int8, int16, int32, int64, int
- Max uint8, uint16, uint32, uint64, uint
- Boolean true and false
- Standard complex64 and complex128
- Array containing standard ints
- Array containing type with custom formatter on pointer receiver only
- Array containing interfaces
- Array containing bytes
- Slice containing standard float32 values
- Slice containing type with custom formatter on pointer receiver only
- Slice containing interfaces
- Slice containing bytes
- Nil slice
- Standard string
- Nil interface
- Sub-interface
- Map with string keys and int vals
- Map with custom formatter type on pointer receiver only keys and vals
- Map with interface keys and values
- Map with nil interface value
- Struct with primitives
- Struct that contains another struct
- Struct that contains custom type with fmt.Stringer pointer interface via both
  exported and unexported fields
- Struct that contains embedded struct and field to same struct
- Uintptr to 0 (null pointer)
- Uintptr address of real variable
- Unsafe.Pointer to 0 (null pointer)
- Unsafe.Pointer to address of real variable
- Nil channel
- Standard int channel
- Function with no params and no returns
- Function with param and no returns
- Function with multiple params and multiple returns
- Struct that is circular through self referencing
- Structs that are circular through cross referencing
- Structs that are indirectly circular
- Type that panics in its fmt.Stringer interface
*/

package spew_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/maxatome/go-testdeep/internal/spew"
)

// dumpTest is used to describe a test to be performed against Sdump function.
type dumpTest struct {
	name  string
	in    any
	wants []string
}

// dumpTests houses all of the tests to be performed against the Dump method.
var dumpTests = make([]dumpTest, 0)

// addSdumpTest is a helper method to append the passed input and desired result
// to dumpTests.
func addSdumpTest(name string, in any, wants ...string) {
	test := dumpTest{name, in, wants}
	dumpTests = append(dumpTests, test)
}

func addIntSdumpTests() {
	// Max int8.
	v := int8(127)
	nv := (*int8)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "int8"
	vs := "127"
	addSdumpTest("int8", v, "("+vt+") "+vs)
	addSdumpTest("int8 ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("int8 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("int8 nil ptr", nv, "(*"+vt+")(<nil>)")

	// Max int16.
	v2 := int16(32767)
	nv2 := (*int16)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "int16"
	v2s := "32767"
	addSdumpTest("int16", v2, "("+v2t+") "+v2s)
	addSdumpTest("int16 ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("int16 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("int16 nil ptr", nv2, "(*"+v2t+")(<nil>)")

	// Max int32.
	v3 := int32(2147483647)
	nv3 := (*int32)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "int32"
	v3s := "2147483647"
	addSdumpTest("int32", v3, "("+v3t+") "+v3s)
	addSdumpTest("int32 ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")")
	addSdumpTest("int32 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")")
	addSdumpTest("int32 nil ptr", nv3, "(*"+v3t+")(<nil>)")

	// Max int64.
	v4 := int64(9223372036854775807)
	nv4 := (*int64)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "int64"
	v4s := "9223372036854775807"
	addSdumpTest("int64", v4, "("+v4t+") "+v4s)
	addSdumpTest("int64 ptr", pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")")
	addSdumpTest("int64 2ptr", &pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")")
	addSdumpTest("int64 nil ptr", nv4, "(*"+v4t+")(<nil>)")

	// Max int.
	v5 := int(2147483647)
	nv5 := (*int)(nil)
	pv5 := &v5
	v5Addr := fmt.Sprintf("%p", pv5)
	pv5Addr := fmt.Sprintf("%p", &pv5)
	v5t := "int"
	v5s := "2147483647"
	addSdumpTest("int", v5, "("+v5t+") "+v5s)
	addSdumpTest("int ptr", pv5, "(*"+v5t+")("+v5Addr+")("+v5s+")")
	addSdumpTest("int 2ptr", &pv5, "(**"+v5t+")("+pv5Addr+"->"+v5Addr+")("+v5s+")")
	addSdumpTest("int nil ptr", nv5, "(*"+v5t+")(<nil>)")
}

func addUintSdumpTests() {
	// Max uint8.
	v := uint8(255)
	nv := (*uint8)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "uint8"
	vs := "255"
	addSdumpTest("uint8", v, "("+vt+") "+vs)
	addSdumpTest("uint8 ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("uint8 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("uint8 nil ptr", nv, "(*"+vt+")(<nil>)")

	// Max uint16.
	v2 := uint16(65535)
	nv2 := (*uint16)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "uint16"
	v2s := "65535"
	addSdumpTest("uint16", v2, "("+v2t+") "+v2s)
	addSdumpTest("uint16 ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("uint16 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("uint16 nil ptr", nv2, "(*"+v2t+")(<nil>)")

	// Max uint32.
	v3 := uint32(4294967295)
	nv3 := (*uint32)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "uint32"
	v3s := "4294967295"
	addSdumpTest("uint32", v3, "("+v3t+") "+v3s)
	addSdumpTest("uint32 ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")")
	addSdumpTest("uint32 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")")
	addSdumpTest("uint32 nil ptr", nv3, "(*"+v3t+")(<nil>)")

	// Max uint64.
	v4 := uint64(18446744073709551615)
	nv4 := (*uint64)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "uint64"
	v4s := "18446744073709551615"
	addSdumpTest("uint64", v4, "("+v4t+") "+v4s)
	addSdumpTest("uint64 ptr", pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")")
	addSdumpTest("uint64 2ptr", &pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")")
	addSdumpTest("uint64 nil ptr", nv4, "(*"+v4t+")(<nil>)")

	// Max uint.
	v5 := uint(4294967295)
	nv5 := (*uint)(nil)
	pv5 := &v5
	v5Addr := fmt.Sprintf("%p", pv5)
	pv5Addr := fmt.Sprintf("%p", &pv5)
	v5t := "uint"
	v5s := "4294967295"
	addSdumpTest("uint", v5, "("+v5t+") "+v5s)
	addSdumpTest("uint ptr", pv5, "(*"+v5t+")("+v5Addr+")("+v5s+")")
	addSdumpTest("uint 2ptr", &pv5, "(**"+v5t+")("+pv5Addr+"->"+v5Addr+")("+v5s+")")
	addSdumpTest("uint nil ptr", nv5, "(*"+v5t+")(<nil>)")
}

func addBoolSdumpTests() {
	// Boolean true.
	v := bool(true)
	nv := (*bool)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "bool"
	vs := "true"
	addSdumpTest("true", v, "("+vt+") "+vs)
	addSdumpTest("true ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("true 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("true nil ptr", nv, "(*"+vt+")(<nil>)")

	// Boolean false.
	v2 := bool(false)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "bool"
	v2s := "false"
	addSdumpTest("false", v2, "("+v2t+") "+v2s)
	addSdumpTest("false ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("false 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
}

func addFloatSdumpTests() {
	// Standard float32.
	v := float32(3.1415)
	nv := (*float32)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "float32"
	vs := "3.1415"
	addSdumpTest("float32", v, "("+vt+") "+vs)
	addSdumpTest("float32 ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("float32 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("float32 nil ptr", nv, "(*"+vt+")(<nil>)")

	// Standard float64.
	v2 := float64(3.1415926)
	nv2 := (*float64)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "float64"
	v2s := "3.1415926"
	addSdumpTest("float64", v2, "("+v2t+") "+v2s)
	addSdumpTest("float64 ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("float64 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("float64 nil ptr", nv2, "(*"+v2t+")(<nil>)")
}

func addComplexSdumpTests() {
	// Standard complex64.
	v := complex(float32(6), -2)
	nv := (*complex64)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "complex64"
	vs := "(6-2i)"
	addSdumpTest("complex64", v, "("+vt+") "+vs)
	addSdumpTest("complex64 ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("complex64 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("complex64 nil ptr", nv, "(*"+vt+")(<nil>)")

	// Standard complex128.
	v2 := complex(float64(-6), 2)
	nv2 := (*complex128)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "complex128"
	v2s := "(-6+2i)"
	addSdumpTest("complex128", v2, "("+v2t+") "+v2s)
	addSdumpTest("complex128 ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("complex128 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("complex128 nil ptr", nv2, "(*"+v2t+")(<nil>)")
}

func addArraySdumpTests() {
	// Array containing standard ints.
	v := [3]int{1, 2, 3}
	vLen := fmt.Sprintf("%d", len(v))
	nv := (*[3]int)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "int"
	vs := "(len=" + vLen + ") {\n (" + vt + ") 1,\n (" +
		vt + ") 2,\n (" + vt + ") 3\n}"
	addSdumpTest("array of int", v, "([3]"+vt+") "+vs)
	addSdumpTest("array of int ptr", pv, "(*[3]"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("array of int 2ptr", &pv, "(**[3]"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("array of int nil ptr", nv, "(*[3]"+vt+")(<nil>)")

	// Array containing type with fmt.Stringer on pointer receiver only.
	v2i0 := pstringer("1")
	v2i1 := pstringer("2")
	v2i2 := pstringer("3")
	v2 := [3]pstringer{v2i0, v2i1, v2i2}
	v2i0Len := fmt.Sprintf("%d", len(v2i0))
	v2i1Len := fmt.Sprintf("%d", len(v2i1))
	v2i2Len := fmt.Sprintf("%d", len(v2i2))
	v2Len := fmt.Sprintf("%d", len(v2))
	nv2 := (*[3]pstringer)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.pstringer"
	v2sp := "(len=" + v2Len + ") {\n (" + v2t +
		") (len=" + v2i0Len + ") stringer 1,\n (" + v2t +
		") (len=" + v2i1Len + ") stringer 2,\n (" + v2t +
		") (len=" + v2i2Len + ") stringer 3\n}"
	v2s := v2sp
	if spew.UnsafeDisabled {
		v2s = "(len=" + v2Len + ") {\n (" + v2t +
			") (len=" + v2i0Len + ") \"1\",\n (" + v2t +
			") (len=" + v2i1Len + ") \"2\",\n (" + v2t +
			") (len=" + v2i2Len + ") " + "\"3\"\n}"
	}
	addSdumpTest("array of stringer", v2, "([3]"+v2t+") "+v2s)
	addSdumpTest("array of stringer ptr", pv2, "(*[3]"+v2t+")("+v2Addr+")("+v2sp+")")
	addSdumpTest("array of stringer 2ptr", &pv2, "(**[3]"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2sp+")")
	addSdumpTest("array of stringer nil ptr", nv2, "(*[3]"+v2t+")(<nil>)")

	// Array containing interfaces.
	v3i0 := "one"
	v3 := [3]any{v3i0, int(2), uint(3)}
	v3i0Len := fmt.Sprintf("%d", len(v3i0))
	v3Len := fmt.Sprintf("%d", len(v3))
	nv3 := (*[3]any)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "[3]interface {}"
	v3t2 := "string"
	v3t3 := "int"
	v3t4 := "uint"
	v3s := "(len=" + v3Len + ") {\n (" + v3t2 + ") " +
		"(len=" + v3i0Len + ") \"one\",\n (" + v3t3 + ") 2,\n (" +
		v3t4 + ") 3\n}"
	addSdumpTest("array of iface", v3, "("+v3t+") "+v3s)
	addSdumpTest("array of iface ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")")
	addSdumpTest("array of iface 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")")
	addSdumpTest("array of iface nil ptr", nv3, "(*"+v3t+")(<nil>)")

	// Array containing bytes.
	v4 := [34]byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}
	v4Len := fmt.Sprintf("%d", len(v4))
	nv4 := (*[34]byte)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "[34]uint8"
	v4s := "(len=" + v4Len + ") " +
		"{\n 00000000  11 12 13 14 15 16 17 18  19 1a 1b 1c 1d 1e 1f 20" +
		"  |............... |\n" +
		" 00000010  21 22 23 24 25 26 27 28  29 2a 2b 2c 2d 2e 2f 30" +
		"  |!\"#$%&'()*+,-./0|\n" +
		" 00000020  31 32                                           " +
		"  |12|\n}"
	addSdumpTest("array of byte", v4, "("+v4t+") "+v4s)
	addSdumpTest("array of byte ptr", pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")")
	addSdumpTest("array of byte 2ptr", &pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")")
	addSdumpTest("array of byte nil ptr", nv4, "(*"+v4t+")(<nil>)")
}

func addSliceSdumpTests() {
	// Slice containing standard float32 values.
	v := []float32{3.14, 6.28, 12.56}
	vLen := fmt.Sprintf("%d", len(v))
	nv := (*[]float32)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "float32"
	vs := "(len=" + vLen + ") {\n (" + vt + ") 3.14,\n (" +
		vt + ") 6.28,\n (" + vt + ") 12.56\n}"
	addSdumpTest("slice of float32", v, "([]"+vt+") "+vs)
	addSdumpTest("slice of float32 ptr", pv, "(*[]"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("slice of float32 2ptr", &pv, "(**[]"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("slice of float32 nil ptr", nv, "(*[]"+vt+")(<nil>)")

	// Slice containing type with custom formatter on pointer receiver only.
	v2i0 := pstringer("1")
	v2i1 := pstringer("2")
	v2i2 := pstringer("3")
	v2 := []pstringer{v2i0, v2i1, v2i2}
	v2i0Len := fmt.Sprintf("%d", len(v2i0))
	v2i1Len := fmt.Sprintf("%d", len(v2i1))
	v2i2Len := fmt.Sprintf("%d", len(v2i2))
	v2Len := fmt.Sprintf("%d", len(v2))
	nv2 := (*[]pstringer)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.pstringer"
	v2s := "(len=" + v2Len + ") {\n (" + v2t + ") (len=" +
		v2i0Len + ") stringer 1,\n (" + v2t + ") (len=" + v2i1Len +
		") stringer 2,\n (" + v2t + ") (len=" + v2i2Len + ") " +
		"stringer 3\n}"
	addSdumpTest("slice of stringer", v2, "([]"+v2t+") "+v2s)
	addSdumpTest("slice of stringer ptr", pv2, "(*[]"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("slice of stringer 2ptr", &pv2, "(**[]"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("slice of stringer nil ptr", nv2, "(*[]"+v2t+")(<nil>)")

	// Slice containing interfaces.
	v3i0 := "one"
	v3 := []any{v3i0, int(2), uint(3), nil}
	v3i0Len := fmt.Sprintf("%d", len(v3i0))
	v3Len := fmt.Sprintf("%d", len(v3))
	nv3 := (*[]any)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "[]interface {}"
	v3t2 := "string"
	v3t3 := "int"
	v3t4 := "uint"
	v3t5 := "interface {}"
	v3s := "(len=" + v3Len + ") {\n (" + v3t2 + ") " +
		"(len=" + v3i0Len + ") \"one\",\n (" + v3t3 + ") 2,\n (" +
		v3t4 + ") 3,\n (" + v3t5 + ") <nil>\n}"
	addSdumpTest("slice of iface", v3, "("+v3t+") "+v3s)
	addSdumpTest("slice of iface ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")")
	addSdumpTest("slice of iface 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")")
	addSdumpTest("slice of iface nil ptr", nv3, "(*"+v3t+")(<nil>)")

	// Slice containing bytes.
	v4 := []byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}
	v4Len := fmt.Sprintf("%d", len(v4))
	nv4 := (*[]byte)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "[]uint8"
	v4s := "(len=" + v4Len + ") " +
		"{\n 00000000  11 12 13 14 15 16 17 18  19 1a 1b 1c 1d 1e 1f 20" +
		"  |............... |\n" +
		" 00000010  21 22 23 24 25 26 27 28  29 2a 2b 2c 2d 2e 2f 30" +
		"  |!\"#$%&'()*+,-./0|\n" +
		" 00000020  31 32                                           " +
		"  |12|\n}"
	addSdumpTest("slice of byte", v4, "("+v4t+") "+v4s)
	addSdumpTest("slice of byte ptr", pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")")
	addSdumpTest("slice of byte 2ptr", &pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")")
	addSdumpTest("slice of byte nil ptr", nv4, "(*"+v4t+")(<nil>)")

	type myUint8 uint8
	v4b := []myUint8{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}
	addSdumpTest("slice of myUint8", v4b, "([]spew_test.myUint8) "+v4s)

	// Nil slice.
	v5 := []int(nil)
	nv5 := (*[]int)(nil)
	pv5 := &v5
	v5Addr := fmt.Sprintf("%p", pv5)
	pv5Addr := fmt.Sprintf("%p", &pv5)
	v5t := "[]int"
	v5s := "<nil>"
	addSdumpTest("nil slice", v5, "("+v5t+") "+v5s)
	addSdumpTest("nil slice ptr", pv5, "(*"+v5t+")("+v5Addr+")("+v5s+")")
	addSdumpTest("nil slice 2ptr", &pv5, "(**"+v5t+")("+pv5Addr+"->"+v5Addr+")("+v5s+")")
	addSdumpTest("nil slice nil ptr", nv5, "(*"+v5t+")(<nil>)")
}

func addStringSdumpTests() {
	// Standard string.
	v := "test"
	vLen := fmt.Sprintf("%d", len(v))
	nv := (*string)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "string"
	vs := "(len=" + vLen + ") \"test\""
	addSdumpTest("string", v, "("+vt+") "+vs)
	addSdumpTest("string ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("string 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("string nil ptr", nv, "(*"+vt+")(<nil>)")
}

func addInterfaceSdumpTests() {
	// Nil interface.
	var v any
	nv := (*any)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "interface {}"
	vs := "<nil>"
	addSdumpTest("nil iface", v, "("+vt+") "+vs)
	addSdumpTest("nil iface ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("nil iface 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("nil iface nil ptr", nv, "(*"+vt+")(<nil>)")

	// Sub-interface.
	v2 := any(uint16(65535))
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "uint16"
	v2s := "65535"
	addSdumpTest("iface", v2, "("+v2t+") "+v2s)
	addSdumpTest("iface ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("iface 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
}

func addMapSdumpTests() {
	// Map with string keys and int vals.
	k := "one"
	kk := "two"
	m := map[string]int{k: 1, kk: 2}
	klen := fmt.Sprintf("%d", len(k)) // not kLen to shut golint up
	kkLen := fmt.Sprintf("%d", len(kk))
	mLen := fmt.Sprintf("%d", len(m))
	nilMap := map[string]int(nil)
	nm := (*map[string]int)(nil)
	pm := &m
	mAddr := fmt.Sprintf("%p", pm)
	pmAddr := fmt.Sprintf("%p", &pm)
	mt := "map[string]int"
	mt1 := "string"
	mt2 := "int"
	ms := "(len=" + mLen + ") {\n (" + mt1 + ") (len=" + klen + ") " +
		"\"one\": (" + mt2 + ") 1,\n (" + mt1 + ") (len=" + kkLen +
		") \"two\": (" + mt2 + ") 2\n}"
	ms2 := "(len=" + mLen + ") {\n (" + mt1 + ") (len=" + kkLen + ") " +
		"\"two\": (" + mt2 + ") 2,\n (" + mt1 + ") (len=" + klen +
		") \"one\": (" + mt2 + ") 1\n}"
	addSdumpTest("map string-int", m,
		"("+mt+") "+ms,
		"("+mt+") "+ms2)
	addSdumpTest("map string-int ptr", pm,
		"(*"+mt+")("+mAddr+")("+ms+")",
		"(*"+mt+")("+mAddr+")("+ms2+")")
	addSdumpTest("map string-int 2ptr", &pm,
		"(**"+mt+")("+pmAddr+"->"+mAddr+")("+ms+")",
		"(**"+mt+")("+pmAddr+"->"+mAddr+")("+ms2+")")
	addSdumpTest("map string-int nil ptr", nm, "(*"+mt+")(<nil>)")
	addSdumpTest("nil map string-int", nilMap, "("+mt+") <nil>")

	// Map with fmt.Stringer on pointer receiver only keys and vals.
	k2 := pstringer("one")
	v2 := pstringer("1")
	m2 := map[pstringer]pstringer{k2: v2}
	k2Len := fmt.Sprintf("%d", len(k2))
	v2Len := fmt.Sprintf("%d", len(v2))
	m2Len := fmt.Sprintf("%d", len(m2))
	nilMap2 := map[pstringer]pstringer(nil)
	nm2 := (*map[pstringer]pstringer)(nil)
	pm2 := &m2
	m2Addr := fmt.Sprintf("%p", pm2)
	pm2Addr := fmt.Sprintf("%p", &pm2)
	m2t := "map[spew_test.pstringer]spew_test.pstringer"
	m2t1 := "spew_test.pstringer"
	m2t2 := "spew_test.pstringer"
	m2s := "(len=" + m2Len + ") {\n (" + m2t1 + ") (len=" + k2Len + ") " +
		"stringer one: (" + m2t2 + ") (len=" + v2Len + ") stringer 1\n}"
	if spew.UnsafeDisabled {
		m2s = "(len=" + m2Len + ") {\n (" + m2t1 + ") (len=" + k2Len + ") " +
			"\"one\": (" + m2t2 + ") (len=" + v2Len + ") \"1\"\n}"
	}
	addSdumpTest("map pstringer-pstringer", m2, "("+m2t+") "+m2s)
	addSdumpTest("map pstringer-pstringer ptr", pm2, "(*"+m2t+")("+m2Addr+")("+m2s+")")
	addSdumpTest("map pstringer-pstringer 2ptr", &pm2, "(**"+m2t+")("+pm2Addr+"->"+m2Addr+")("+m2s+")")
	addSdumpTest("map pstringer-pstringer nil ptr", nm2, "(*"+m2t+")(<nil>)")
	addSdumpTest("nil map pstringer-pstringer", nilMap2, "("+m2t+") <nil>")

	// Map with interface keys and values.
	k3 := "one"
	k3Len := fmt.Sprintf("%d", len(k3))
	m3 := map[any]any{k3: 1}
	m3Len := fmt.Sprintf("%d", len(m3))
	nilMap3 := map[any]any(nil)
	nm3 := (*map[any]any)(nil)
	pm3 := &m3
	m3Addr := fmt.Sprintf("%p", pm3)
	pm3Addr := fmt.Sprintf("%p", &pm3)
	m3t := "map[interface {}]interface {}"
	m3t1 := "string"
	m3t2 := "int"
	m3s := "(len=" + m3Len + ") {\n (" + m3t1 + ") (len=" + k3Len + ") " +
		"\"one\": (" + m3t2 + ") 1\n}"
	addSdumpTest("map any-any", m3, "("+m3t+") "+m3s)
	addSdumpTest("map any-any ptr", pm3, "(*"+m3t+")("+m3Addr+")("+m3s+")")
	addSdumpTest("map any-any 2ptr", &pm3, "(**"+m3t+")("+pm3Addr+"->"+m3Addr+")("+m3s+")")
	addSdumpTest("map any-any nil ptr", nm3, "(*"+m3t+")(<nil>)")
	addSdumpTest("nil map any-any", nilMap3, "("+m3t+") <nil>")

	// Map with nil interface value.
	k4 := "nil"
	k4Len := fmt.Sprintf("%d", len(k4))
	m4 := map[string]any{k4: nil}
	m4Len := fmt.Sprintf("%d", len(m4))
	nilMap4 := map[string]any(nil)
	nm4 := (*map[string]any)(nil)
	pm4 := &m4
	m4Addr := fmt.Sprintf("%p", pm4)
	pm4Addr := fmt.Sprintf("%p", &pm4)
	m4t := "map[string]interface {}"
	m4t1 := "string"
	m4t2 := "interface {}"
	m4s := "(len=" + m4Len + ") {\n (" + m4t1 + ") (len=" + k4Len + ")" +
		" \"nil\": (" + m4t2 + ") <nil>\n}"
	addSdumpTest("map string-any", m4, "("+m4t+") "+m4s)
	addSdumpTest("map string-any ptr", pm4, "(*"+m4t+")("+m4Addr+")("+m4s+")")
	addSdumpTest("map string-any 2ptr", &pm4, "(**"+m4t+")("+pm4Addr+"->"+m4Addr+")("+m4s+")")
	addSdumpTest("map string-any nil ptr", nm4, "(*"+m4t+")(<nil>)")
	addSdumpTest("nil map string-any", nilMap4, "("+m4t+") <nil>")
}

func addStructSdumpTests() {
	// Struct with primitives.
	type s1 struct {
		a int8
		b uint8
	}
	v := s1{127, 255}
	nv := (*s1)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.s1"
	vt2 := "int8"
	vt3 := "uint8"
	vs := "{\n a: (" + vt2 + ") 127,\n b: (" + vt3 + ") 255\n}"
	addSdumpTest("struct", v, "("+vt+") "+vs)
	addSdumpTest("struct ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("struct 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("struct nil ptr", nv, "(*"+vt+")(<nil>)")

	// Struct that contains another struct.
	type s2 struct {
		s1 s1
		b  bool
	}
	v2 := s2{s1{127, 255}, true}
	nv2 := (*s2)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.s2"
	v2t2 := "spew_test.s1"
	v2t3 := "int8"
	v2t4 := "uint8"
	v2t5 := "bool"
	v2s := "{\n s1: (" + v2t2 + ") {\n  a: (" + v2t3 + ") 127,\n  b: (" +
		v2t4 + ") 255\n },\n b: (" + v2t5 + ") true\n}"
	addSdumpTest("struct of struct", v2, "("+v2t+") "+v2s)
	addSdumpTest("struct of struct ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("struct of struct 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("struct of struct nil ptr", nv2, "(*"+v2t+")(<nil>)")

	// Struct that contains custom type with fmt.Stringer pointer
	// interface via both exported and unexported fields.
	type s3 struct {
		s pstringer
		S pstringer
	}
	v3 := s3{"test", "test2"}
	nv3 := (*s3)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "spew_test.s3"
	v3t2 := "spew_test.pstringer"
	v3s := "{\n s: (" + v3t2 + ") (len=4) stringer test,\n S: (" + v3t2 +
		") (len=5) stringer test2\n}"
	v3sp := v3s
	if spew.UnsafeDisabled {
		v3s = "{\n s: (" + v3t2 + ") (len=4) \"test\",\n S: (" + v3t2 +
			") (len=5) \"test2\"\n}"
		v3sp = "{\n s: (" + v3t2 + ") (len=4) \"test\",\n S: (" + v3t2 +
			") (len=5) stringer test2\n}"
	}
	addSdumpTest("struct of pstringer", v3, "("+v3t+") "+v3s)
	addSdumpTest("struct of pstringer ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3sp+")")
	addSdumpTest("struct of pstringer 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3sp+")")
	addSdumpTest("struct of pstringer nil ptr", nv3, "(*"+v3t+")(<nil>)")

	// Struct that contains embedded struct and field to same struct.
	e := embed{"embedstr"}
	eLen := fmt.Sprintf("%d", len("embedstr"))
	v4 := embedwrap{embed: &e, e: &e}
	nv4 := (*embedwrap)(nil)
	pv4 := &v4
	eAddr := fmt.Sprintf("%p", &e)
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "spew_test.embedwrap"
	v4t2 := "spew_test.embed"
	v4t3 := "string"
	v4s := "{\n embed: (*" + v4t2 + ")(" + eAddr + ")({\n  a: (" + v4t3 +
		") (len=" + eLen + ") \"embedstr\"\n }),\n e: (*" + v4t2 +
		")(" + eAddr + ")({\n  a: (" + v4t3 + ") (len=" + eLen + ")" +
		" \"embedstr\"\n })\n}"
	addSdumpTest("embedded struct", v4, "("+v4t+") "+v4s)
	addSdumpTest("embedded struct ptr", pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")")
	addSdumpTest("embedded struct 2ptr", &pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")")
	addSdumpTest("embedded struct nil ptr", nv4, "(*"+v4t+")(<nil>)")
}

func addUintptrSdumpTests() {
	// Null pointer.
	v := uintptr(0)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "uintptr"
	vs := "<nil>"
	addSdumpTest("uintptr", v, "("+vt+") "+vs)
	addSdumpTest("uintptr ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("uintptr 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")

	// Address of real variable.
	i := 1
	v2 := uintptr(unsafe.Pointer(&i))
	nv2 := (*uintptr)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "uintptr"
	v2s := fmt.Sprintf("%p", &i)
	addSdumpTest("real uintptr", v2, "("+v2t+") "+v2s)
	addSdumpTest("real uintptr ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("real uintptr 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("real uintptr nil ptr", nv2, "(*"+v2t+")(<nil>)")
}

func addUnsafePointerSdumpTests() {
	// Null pointer.
	v := unsafe.Pointer(nil)
	nv := (*unsafe.Pointer)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "unsafe.Pointer"
	vs := "<nil>"
	addSdumpTest("unsafe.Pointer", v, "("+vt+") "+vs)
	addSdumpTest("unsafe.Pointer ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("unsafe.Pointer 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("unsafe.Pointer nil ptr", nv, "(*"+vt+")(<nil>)")

	// Address of real variable.
	i := 1
	v2 := unsafe.Pointer(&i)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "unsafe.Pointer"
	v2s := fmt.Sprintf("%p", &i)
	addSdumpTest("real unsafe.Pointer", v2, "("+v2t+") "+v2s)
	addSdumpTest("real unsafe.Pointer ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("real unsafe.Pointer 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("real unsafe.Pointer nil ptr", nv, "(*"+vt+")(<nil>)")
}

func addChanSdumpTests() {
	// Nil channel.
	var v chan int
	pv := &v
	nv := (*chan int)(nil)
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "chan int"
	vs := "<nil>"
	addSdumpTest("nil chan", v, "("+vt+") "+vs)
	addSdumpTest("nil chan ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("nil chan 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("nil chan nil ptr", nv, "(*"+vt+")(<nil>)")

	// Real channel.
	v2 := make(chan int)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "chan int"
	v2s := fmt.Sprintf("%p", v2)
	addSdumpTest("chan", v2, "("+v2t+") "+v2s)
	addSdumpTest("chan ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("chan 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
}

func addFuncSdumpTests() {
	// Function with no params and no returns.
	v := addIntSdumpTests
	nv := (*func())(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "func()"
	vs := fmt.Sprintf("%p", v)
	addSdumpTest("func()", v, "("+vt+") "+vs)
	addSdumpTest("func() ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("func() 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("func() nil ptr", nv, "(*"+vt+")(<nil>)")

	// Function with param and no returns.
	v2 := TestSdump
	nv2 := (*func(*testing.T))(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "func(*testing.T)"
	v2s := fmt.Sprintf("%p", v2)
	addSdumpTest("func(P)", v2, "("+v2t+") "+v2s)
	addSdumpTest("func(P) ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")")
	addSdumpTest("func(P) 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")")
	addSdumpTest("func(P) nil ptr", nv2, "(*"+v2t+")(<nil>)")

	// Function with multiple params and multiple returns.
	v3 := func(i int, s string) (b bool, err error) {
		return true, nil
	}
	nv3 := (*func(int, string) (bool, error))(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "func(int, string) (bool, error)"
	v3s := fmt.Sprintf("%p", v3)
	addSdumpTest("func(P)R", v3, "("+v3t+") "+v3s)
	addSdumpTest("func(P)R ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")")
	addSdumpTest("func(P)R 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")")
	addSdumpTest("func(P)R nil ptr", nv3, "(*"+v3t+")(<nil>)")
}

func addCircularSdumpTests() {
	// Struct that is circular through self referencing.
	type circular struct {
		c *circular
	}
	v := circular{nil}
	v.c = &v
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.circular"
	vs := "{\n c: (*" + vt + ")(" + vAddr + ")({\n  c: (*" + vt + ")(" +
		vAddr + ")(<already shown>)\n })\n}"
	vs2 := "{\n c: (*" + vt + ")(" + vAddr + ")(<already shown>)\n}"
	addSdumpTest("circ struct", v, "("+vt+") "+vs)
	addSdumpTest("circ struct ptr", pv, "(*"+vt+")("+vAddr+")("+vs2+")")
	addSdumpTest("circ struct 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs2+")")

	// Structs that are circular through cross referencing.
	v2 := xref1{nil}
	ts2 := xref2{&v2}
	v2.ps2 = &ts2
	pv2 := &v2
	ts2Addr := fmt.Sprintf("%p", &ts2)
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.xref1"
	v2t2 := "spew_test.xref2"
	v2s := "{\n ps2: (*" + v2t2 + ")(" + ts2Addr + ")({\n  ps1: (*" + v2t +
		")(" + v2Addr + ")({\n   ps2: (*" + v2t2 + ")(" + ts2Addr +
		")(<already shown>)\n  })\n })\n}"
	v2s2 := "{\n ps2: (*" + v2t2 + ")(" + ts2Addr + ")({\n  ps1: (*" + v2t +
		")(" + v2Addr + ")(<already shown>)\n })\n}"
	addSdumpTest("circ X struct", v2, "("+v2t+") "+v2s)
	addSdumpTest("circ X struct ptr", pv2, "(*"+v2t+")("+v2Addr+")("+v2s2+")")
	addSdumpTest("circ X struct 2ptr", &pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s2+")")

	// Structs that are indirectly circular.
	v3 := indirCir1{nil}
	tic2 := indirCir2{nil}
	tic3 := indirCir3{&v3}
	tic2.ps3 = &tic3
	v3.ps2 = &tic2
	pv3 := &v3
	tic2Addr := fmt.Sprintf("%p", &tic2)
	tic3Addr := fmt.Sprintf("%p", &tic3)
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "spew_test.indirCir1"
	v3t2 := "spew_test.indirCir2"
	v3t3 := "spew_test.indirCir3"
	v3s := "{\n ps2: (*" + v3t2 + ")(" + tic2Addr + ")({\n  ps3: (*" + v3t3 +
		")(" + tic3Addr + ")({\n   ps1: (*" + v3t + ")(" + v3Addr +
		")({\n    ps2: (*" + v3t2 + ")(" + tic2Addr +
		")(<already shown>)\n   })\n  })\n })\n}"
	v3s2 := "{\n ps2: (*" + v3t2 + ")(" + tic2Addr + ")({\n  ps3: (*" + v3t3 +
		")(" + tic3Addr + ")({\n   ps1: (*" + v3t + ")(" + v3Addr +
		")(<already shown>)\n  })\n })\n}"
	addSdumpTest("circ indir struct", v3, "("+v3t+") "+v3s)
	addSdumpTest("circ indir struct ptr", pv3, "(*"+v3t+")("+v3Addr+")("+v3s2+")")
	addSdumpTest("circ indir struct 2ptr", &pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s2+")")
}

func addPanicSdumpTests() {
	// Type that panics in its fmt.Stringer interface.
	v := panicer(127)
	nv := (*panicer)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.panicer"
	vs := "(PANIC=test panic)127"
	addSdumpTest("panicer", v, "("+vt+") "+vs)
	addSdumpTest("panicer ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("panicer 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("panicer nil ptr", nv, "(*"+vt+")(<nil>)")
}

func addErrorSdumpTests() {
	// Type that has a custom Error interface.
	v := customError(127)
	nv := (*customError)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.customError"
	vs := "error: 127"
	addSdumpTest("customError", v, "("+vt+") "+vs)
	addSdumpTest("customError ptr", pv, "(*"+vt+")("+vAddr+")("+vs+")")
	addSdumpTest("customError 2ptr", &pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")")
	addSdumpTest("customError nil ptr", nv, "(*"+vt+")(<nil>)")
}

// TestSdump executes all of the tests described by dumpTests.
func TestSdump(t *testing.T) {
	// Setup tests.
	addIntSdumpTests()
	addUintSdumpTests()
	addBoolSdumpTests()
	addFloatSdumpTests()
	addComplexSdumpTests()
	addArraySdumpTests()
	addSliceSdumpTests()
	addStringSdumpTests()
	addInterfaceSdumpTests()
	addMapSdumpTests()
	addStructSdumpTests()
	addUintptrSdumpTests()
	addUnsafePointerSdumpTests()
	addChanSdumpTests()
	addFuncSdumpTests()
	addCircularSdumpTests()
	addPanicSdumpTests()
	addErrorSdumpTests()
	addCgoSdumpTests()

	t.Logf("Running %d tests", len(dumpTests))
	for _, test := range dumpTests {
		t.Run(test.name, func(t *testing.T) {
			s := spew.Sdump(test.in)
			if testFailed(s, test.wants) {
				t.Errorf("%s\n got: %s\n%s", test.name, s, stringizeWants(test.wants))
			}
		})
	}
}
