// Copyright (c) 2019-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Work around https://github.com/golang/go/issues/26995 issue
// (corrected in go 1.12).
// +build go1.12

package td_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestSetlocation(t *testing.T) {
//line /go-testdeep/types_test.go:10
	tt := &tdutil.T{}
	ok := td.Cmp(tt, 12, 13)
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:11: Failed test
        DATA: values differ
        	     got: 12
        	expected: 13
`)
	} else {
		t.Error("Cmp returned true!")
	}

//line /go-testdeep/types_test.go:20
	tt = &tdutil.T{}
	ok = td.Cmp(tt,
		12,
		td.Any(13, 14, 15))
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:21: Failed test
        DATA: comparing with Any
        	     got: 12
        	expected: Any(13,
        	              14,
        	              15)
        [under operator Any at types_test.go:23]
`)
	} else {
		t.Error("Cmp returned true!")
	}

//line /go-testdeep/types_test.go:30
	tt = &tdutil.T{}
	ok = td.CmpAny(tt,
		12,
		[]interface{}{13, 14, 15})
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:31: Failed test
        DATA: comparing with Any
        	     got: 12
        	expected: Any(13,
        	              14,
        	              15)
`)
	} else {
		t.Error("CmpAny returned true!")
	}

//line /go-testdeep/types_test.go:40
	tt = &tdutil.T{}
	ttt := td.NewT(tt)
	ok = ttt.Cmp(
		12,
		td.Any(13, 14, 15))
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:42: Failed test
        DATA: comparing with Any
        	     got: 12
        	expected: Any(13,
        	              14,
        	              15)
        [under operator Any at types_test.go:44]
`)
	} else {
		t.Error("Cmp returned true!")
	}

//line /go-testdeep/types_test.go:50
	tt = &tdutil.T{}
	ttt = td.NewT(tt)
	ok = ttt.Any(
		12,
		[]interface{}{13, 14, 15})
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:52: Failed test
        DATA: comparing with Any
        	     got: 12
        	expected: Any(13,
        	              14,
        	              15)
`)
	} else {
		t.Error("Cmp returned true!")
	}
}
