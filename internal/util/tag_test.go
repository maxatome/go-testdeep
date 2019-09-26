// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/internal/util"
)

func TestCheckTag(t *testing.T) {
	tags := []string{
		"tag12",
		"_1é",
		"a9",
		"a",
		"é൫",
		"é",
		"_",
	}
	for _, tag := range tags {
		if err := util.CheckTag(tag); err != nil {
			t.Errorf("check(%s) failed: %s", tag, err)
		}
	}

	tagsInfo := []struct {
		tag string
		err error
	}{
		{tag: "", err: util.ErrTagEmpty},
		{tag: "൫a", err: util.ErrTagInvalid},
		{tag: "9a", err: util.ErrTagInvalid},
		{tag: "é ", err: util.ErrTagInvalid},
	}
	for _, info := range tagsInfo {
		err := util.CheckTag(info.tag)
		if err == nil {
			t.Errorf("check(%s) should not succeed", info.tag)
		} else if err != info.err {
			t.Errorf(`check(%s) returned "%s" intead of expected "%s"`,
				info.tag, err, info.err)
		}
	}
}
