// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"bytes"
	"reflect"

	"github.com/maxatome/go-testdeep/internal/util"
)

type tdList struct {
	BaseOKNil
	items []reflect.Value
}

func newList(items ...interface{}) (ret tdList) {
	ret.BaseOKNil = NewBaseOKNil(4)
	ret.items = make([]reflect.Value, len(items))

	for idx, item := range items {
		ret.items[idx] = reflect.ValueOf(item)
	}
	return
}

func (l *tdList) String() string {
	return util.SliceToBuffer(bytes.NewBufferString(l.GetLocation().Func), l.items).
		String()
}
