// Copyright (c) 2021-2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package types_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/internal/types"
)

var _ error = types.OperatorNotJSONMarshallableError("")

func TestOperatorNotJSONMarshallableError(t *testing.T) {
	e := types.OperatorNotJSONMarshallableError("Pipo")

	if e.Error() != "Pipo TestDeep operator cannot be json.Marshal'led" {
		t.Errorf("unexpected %q", e.Error())
	}

	if e.Operator() != "Pipo" {
		t.Errorf("unexpected %q", e.Operator())
	}

	t.Run("AsOperatorNotJSONMarshallableError", func(t *testing.T) {
		ne, ok := types.AsOperatorNotJSONMarshallableError(e)
		if !ok {
			t.Error("AsOperatorNotJSONMarshallableError() returned false")
			return
		}
		if ne != e {
			t.Errorf("AsOperatorNotJSONMarshallableError(): %q ≠ %q",
				ne.Error(), e.Error())
		}

		other := errors.New("Other error")
		_, ok = types.AsOperatorNotJSONMarshallableError(other)
		if ok {
			t.Error("AsOperatorNotJSONMarshallableError() returned true")
			return
		}

		je := &json.MarshalerError{Err: e}
		ne, ok = types.AsOperatorNotJSONMarshallableError(je)
		if !ok {
			t.Error("AsOperatorNotJSONMarshallableError() returned false")
			return
		}
		if ne != e {
			t.Errorf("AsOperatorNotJSONMarshallableError(): %q ≠ %q",
				ne.Error(), e.Error())
		}

		je.Err = other
		_, ok = types.AsOperatorNotJSONMarshallableError(je)
		if ok {
			t.Error("AsOperatorNotJSONMarshallableError() returned true")
			return
		}
	})
}

func TestRawString(t *testing.T) {
	s := types.RawString("foo")
	if str := s.String(); str != "foo" {
		t.Errorf("Very weird, got %s", str)
	}
}

func TestRawInt(t *testing.T) {
	i := types.RawInt(42)
	if str := i.String(); str != "42" {
		t.Errorf("Very weird, got %s", str)
	}
}

func TestRecvKind(t *testing.T) {
	s := types.RecvNothing.String()
	if s != "nothing received on channel" {
		t.Errorf(`got: %q / expected: "nothing received on channel"`, s)
	}

	s = types.RecvClosed.String()
	if s != "channel is closed" {
		t.Errorf(`got: %q / expected: "channel is closed"`, s)
	}
}
