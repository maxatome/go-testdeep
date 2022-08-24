// Copyright (c) 2020-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package anchors_test

import (
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/internal/anchors"
)

func TestAddAnchorableStructType(t *testing.T) {
	oldAnchorableTypes := anchors.AnchorableTypes
	defer func() { anchors.AnchorableTypes = oldAnchorableTypes }()

	type ok struct{ index int }
	type notComparable struct{ s []int } //nolint: unused

	// Usage error cases
	for i, fn := range []any{
		12,
		func(x ...int) {},
		func(x, y int) {},
		func(x int) (int, int) { return 0, 0 },
		func(x byte) int { return 0 },
		func(x int) int { return 0 },
	} {
		err := anchors.AddAnchorableStructType(fn)
		if err == nil {
			t.Fatalf("#%d function should return an error", i)
		}
		if !strings.HasPrefix(err.Error(), "usage: ") {
			t.Errorf("#%d function returned: `%s` instead of usage", i, err)
		}
	}

	// Not comparable struct
	err := anchors.AddAnchorableStructType(func(nextAnchor int) notComparable {
		return notComparable{}
	})
	if err == nil {
		t.Fatal("function should return an error")
	}
	if err.Error() != "type anchors_test.notComparable is not comparable, it cannot be anchorable" {
		t.Errorf("function returned: `%s` instead of not comparable error", err)
	}

	// Comparable struct => OK
	err = anchors.AddAnchorableStructType(func(nextAnchor int) ok {
		return ok{index: 1000 + nextAnchor}
	})
	if err != nil {
		t.Fatalf("AddAnchorableStructType failed: %s", err)
	}
	if len(anchors.AnchorableTypes) != 2 {
		t.Fatalf("Bad number of anchored type: got=%d expected=2",
			len(anchors.AnchorableTypes))
	}
}
