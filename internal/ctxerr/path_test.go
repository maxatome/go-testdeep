// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package ctxerr_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestPath(t *testing.T) {
	for i, testCase := range []struct {
		Path     ctxerr.Path
		Expected string
	}{
		{
			Path:     ctxerr.Path(nil),
			Expected: "",
		},
		{
			Path:     ctxerr.Path{},
			Expected: "",
		},
		{
			Path:     ctxerr.NewPath("DATA"),
			Expected: "DATA",
		},
		{
			Path:     ctxerr.NewPath("DATA").AddField("field"),
			Expected: "DATA.field",
		},
		{
			Path:     ctxerr.NewPath("DATA").AddPtr(1),
			Expected: "*DATA",
		},
		{
			Path:     ctxerr.NewPath("DATA").AddPtr(2),
			Expected: "**DATA",
		},
		{
			Path:     ctxerr.NewPath("DATA").AddPtr(1).AddField("field"),
			Expected: "DATA.field",
		},
		{
			Path:     ctxerr.NewPath("DATA").AddPtr(2).AddField("field"),
			Expected: "(*DATA).field",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddField("field").
				AddArrayIndex(42),
			Expected: "DATA.field[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddArrayIndex(42),
			Expected: "(*DATA)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddField("field").
				AddPtr(1).
				AddArrayIndex(42),
			Expected: "(*DATA.field)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddField("field1").
				AddPtr(1).
				AddField("field2").
				AddPtr(1).
				AddArrayIndex(42),
			Expected: "(*DATA.field1.field2)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddArrayIndex(42).
				AddPtr(1),
			Expected: "*(*DATA)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddArrayIndex(42).
				AddPtr(1).
				AddField("field"),
			Expected: "(*DATA)[42].field",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddMapKey("key").
				AddPtr(1),
			Expected: `*(*DATA)["key"]`,
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(1).
				AddMapKey("key").
				AddPtr(1).
				AddField("field"),
			Expected: `(*DATA)["key"].field`,
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(2).
				AddField("field1").
				AddPtr(3).
				AddField("field2").
				AddPtr(1).
				AddArrayIndex(42),
			Expected: "(*(**(*DATA).field1).field2)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(2).
				AddField("field1").
				AddPtr(3).
				AddFunctionCall("FUNC").
				AddArrayIndex(42),
			Expected: "FUNC(***(*DATA).field1)[42]",
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(2).
				AddField("field1").
				AddPtr(3).
				AddFunctionCall("FUNC").
				AddArrayIndex(42).
				AddMapKey("key"),
			Expected: `FUNC(***(*DATA).field1)[42]["key"]`,
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(2).
				AddField("field1").
				AddPtr(3).
				AddFunctionCall("FUNC").
				AddArrayIndex(42).
				AddMapKey("key").
				AddCustomLevel("→panic"),
			Expected: `FUNC(***(*DATA).field1)[42]["key"]→panic`,
		},
		{
			Path: ctxerr.NewPath("DATA").
				AddPtr(2).
				AddField("field1").
				AddPtr(3).
				AddFunctionCall("FUNC").
				AddPtr(2).
				AddArrayIndex(42).
				AddMapKey("key").
				AddCustomLevel("→panic"),
			Expected: `(**FUNC(***(*DATA).field1))[42]["key"]→panic`,
		},
	} {
		test.EqualStr(t,
			testCase.Path.String(), testCase.Expected,
			"test case #%d", i)
	}

	var nilPath ctxerr.Path
	for i, newPath := range []ctxerr.Path{
		nilPath.Copy(),
		nilPath.AddField("foo"),
		nilPath.AddArrayIndex(42),
		nilPath.AddMapKey("bar"),
		nilPath.AddPtr(12),
		nilPath.AddFunctionCall("zip"),
		nilPath.AddCustomLevel("custom"),
	} {
		if newPath != nil {
			t.Errorf("at #%d, got=%p expected=nil", i, newPath)
		}
	}
}

func TestEqual(t *testing.T) {
	path := ctxerr.NewPath("DATA").
		AddPtr(2).
		AddField("field1")
	test.EqualInt(t, path.Len(), 2)

	test.IsTrue(t, path.Equal(ctxerr.NewPath("DATA").AddPtr(1).AddPtr(1).AddField("field1")))

	test.IsFalse(t, path.Equal(ctxerr.NewPath("DATA")))
	test.IsFalse(t, path.Equal(ctxerr.NewPath("DATA").AddPtr(2).AddField("field2")))
}

/*
func BenchmarkStringString(b *testing.B) {
	path := ctxerr.NewPath("DATA").
		AddField("field1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path.String()
	}
}

func BenchmarkStringByte(b *testing.B) {
	path := ctxerr.NewPath("DATA").
		AddField("field1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path.Stringx()
	}
}
*/
