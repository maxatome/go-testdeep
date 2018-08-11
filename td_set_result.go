// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"reflect"

	"github.com/maxatome/go-testdeep/internal/str"
	"github.com/maxatome/go-testdeep/internal/types"
)

type tdSetResultKind uint8

const (
	itemsSetResult tdSetResultKind = iota
	keysSetResult
)

// Implements fmt.Stringer.
func (k tdSetResultKind) String() string {
	switch k {
	case itemsSetResult:
		return "items"
	case keysSetResult:
		return "keys"
	default:
		return "?"
	}
}

type tdSetResult struct {
	types.TestDeepStamp
	Missing []reflect.Value
	Extra   []reflect.Value
	Kind    tdSetResultKind
}

var _ types.TestDeepStringer = tdSetResult{}

func (r tdSetResult) IsEmpty() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0
}

func (r tdSetResult) String() string {
	buf := &bytes.Buffer{}

	if len(r.Missing) > 0 {
		buf.WriteString("Missing ")
		buf.WriteString(r.Kind.String())
		buf.WriteString(": ")
		str.SliceToBuffer(buf, r.Missing)
	}

	if len(r.Extra) > 0 {
		if buf.Len() > 0 {
			buf.WriteString("\n  ")
		}
		buf.WriteString("Extra ")
		buf.WriteString(r.Kind.String())
		buf.WriteString(": ")
		str.SliceToBuffer(buf, r.Extra)
	}

	return buf.String()
}
