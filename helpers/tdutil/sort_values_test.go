// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdutil_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
)

func TestSortValues(t *testing.T) {
	s := []reflect.Value{
		reflect.ValueOf(4),
		reflect.ValueOf(3),
		reflect.ValueOf(1),
	}
	sort.Sort(tdutil.SortableValues(s))
	if s[0].Int() != 1 || s[1].Int() != 3 || s[2].Int() != 4 {
		t.Errorf("sort error: [ %v, %v, %v ]", s[0].Int(), s[1].Int(), s[2].Int())
	}

	s = []reflect.Value{
		reflect.ValueOf(42),
	}
	sort.Sort(tdutil.SortableValues(s))
	if s[0].Int() != 42 {
		t.Errorf("sort error: [ %v ]", s[0].Int())
	}

	sort.Sort(tdutil.SortableValues(nil))
}
