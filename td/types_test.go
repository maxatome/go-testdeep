// Copyright (c) 2019-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/td"
)

func TestSetlocation(t *testing.T) {
	//nolint: gocritic
//line types_test.go:10
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

	//nolint: gocritic
//line types_test.go:20
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

	//nolint: gocritic
//line types_test.go:30
	tt = &tdutil.T{}
	ok = td.CmpAny(tt,
		12,
		[]any{13, 14, 15})
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

	//nolint: gocritic
//line types_test.go:40
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

	//nolint: gocritic
//line types_test.go:50
	tt = &tdutil.T{}
	ttt = td.NewT(tt)
	ok = ttt.Any(
		12,
		[]any{13, 14, 15})
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

//line /a/full/path/types_test.go:50
	tt = &tdutil.T{}
	ttt = td.NewT(tt)
	ok = ttt.Any(
		12,
		[]any{13, 14, 15})
	if !ok {
		test.EqualStr(t, tt.LogBuf(), `    types_test.go:52: Failed test
        DATA: comparing with Any
        	     got: 12
        	expected: Any(13,
        	              14,
        	              15)
        This is how we got here:
        	TestSetlocation() /a/full/path/types_test.go:52
`) // at least one '/' in file name → "This is how we got here"
	} else {
		t.Error("Cmp returned true!")
	}
}

func TestError(t *testing.T) {
	test.NoError(t, td.Re(`x`).Error())
	test.Error(t, td.Re(123).Error())
}

func TestMarshalJSON(t *testing.T) {
	op := td.String("foo")

	_, err := json.Marshal(op)
	if test.Error(t, err) {
		test.IsTrue(t, strings.HasSuffix(err.Error(), "String TestDeep operator cannot be json.Marshal'led"))
	}
}
