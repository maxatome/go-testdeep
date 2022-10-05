// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp/internal"
	"github.com/maxatome/go-testdeep/td"
)

func TestResponseList(t *testing.T) {
	assert, require := td.AssertRequire(t)

	rList := internal.NewResponseList()
	require.NotNil(rList)

	assert.Nil(rList.Last())
	assert.Nil(rList.Get("unknown"))

	assert.String(rList.RecordLast("first"), "no last response to record")

	r1 := internal.NewResponse(newResponseRecorder("12345"))
	require.NotNil(r1)
	rList.SetLast(r1)
	assert.Shallow(rList.Last(), r1)

	assert.CmpNoError(rList.RecordLast("first"))
	assert.Shallow(rList.Get("first"), r1)

	assert.String(rList.RecordLast("again"), `last response is already recorded as "first"`)

	r2 := internal.NewResponse(newResponseRecorder("body"))
	require.NotNil(r2)
	rList.SetLast(r2)
	assert.Shallow(rList.Last(), r2)

	assert.CmpNoError(rList.RecordLast("second"))
	assert.Shallow(rList.Get("second"), r2)
	assert.Shallow(rList.Get("first"), r1)

	rList.Reset()
	assert.Nil(rList.Last())
	assert.Nil(rList.Get("first"))
	assert.Nil(rList.Get("second"))
}
