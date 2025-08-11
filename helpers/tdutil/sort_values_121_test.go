// Copyright (c) 2025, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

//go:build go1.21
// +build go1.21

package tdutil_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestCmpValuesFunc(t *testing.T) {
	s := []reflect.Value{
		reflect.ValueOf(4),
		reflect.ValueOf(3),
		reflect.ValueOf(1),
	}
	slices.SortFunc(s, tdutil.CmpValuesFunc())
	if s[0].Int() != 1 || s[1].Int() != 3 || s[2].Int() != 4 {
		t.Errorf("sort error: [ %v, %v, %v ]", s[0].Int(), s[1].Int(), s[2].Int())
	}
}
