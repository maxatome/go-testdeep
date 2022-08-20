// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/internal/test"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/td"
)

func TestRecv(t *testing.T) {
	fillCh := func(ch chan int, val int) {
		ch <- val // td.Cmp
		ch <- val // EqDeeply aka boolean context
		ch <- val // EqDeeplyError
		ch <- val // interface + td.Cmp
		ch <- val // interface + EqDeeply aka boolean context
		ch <- val // interface + EqDeeplyError
	}
	mkCh := func(val int) chan int {
		ch := make(chan int, 6)
		fillCh(ch, val)
		close(ch)
		return ch
	}

	t.Run("all good", func(t *testing.T) {
		ch := mkCh(1)
		checkOK(t, ch, td.Recv(1))
		checkOK(t, ch, td.Recv(td.RecvClosed, 10*time.Microsecond))

		ch = mkCh(42)
		checkOK(t, ch, td.Recv(td.Between(40, 45)))
		checkOK(t, ch, td.Recv(td.RecvClosed))
	})

	t.Run("complete cycle", func(t *testing.T) {
		ch := make(chan int, 6)

		t.Run("empty", func(t *testing.T) {
			checkOK(t, ch, td.Recv(td.RecvNothing))
			checkOK(t, ch, td.Recv(td.RecvNothing, 10*time.Microsecond))

			checkOK(t, &ch, td.Recv(td.RecvNothing))
			checkOK(t, &ch, td.Recv(td.RecvNothing, 10*time.Microsecond))
		})

		t.Run("just filled", func(t *testing.T) {
			fillCh(ch, 33)
			checkOK(t, ch, td.Recv(33))

			fillCh(ch, 34)
			checkOK(t, &ch, td.Recv(34))
		})

		t.Run("nothing to recv on channel", func(t *testing.T) {
			checkError(t, ch, td.Recv(td.RecvClosed),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("nothing received on channel"),
					Expected: mustBe("channel is closed"),
				})
			checkError(t, &ch, td.Recv(td.RecvClosed),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("nothing received on channel"),
					Expected: mustBe("channel is closed"),
				})

			checkError(t, ch, td.Recv(42),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("nothing received on channel"),
					Expected: mustBe("42"),
				})
			checkError(t, &ch, td.Recv(42),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("nothing received on channel"),
					Expected: mustBe("42"),
				})
		})

		close(ch)

		t.Run("closed channel", func(t *testing.T) {
			checkError(t, ch, td.Recv(td.RecvNothing),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("channel is closed"),
					Expected: mustBe("nothing received on channel"),
				})
			checkError(t, &ch, td.Recv(td.RecvNothing),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("channel is closed"),
					Expected: mustBe("nothing received on channel"),
				})

			checkError(t, ch, td.Recv(42),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("channel is closed"),
					Expected: mustBe("42"),
				})
			checkError(t, &ch, td.Recv(42),
				expectedError{
					Message:  mustBe("values differ"),
					Path:     mustBe("recv(DATA)"),
					Got:      mustBe("channel is closed"),
					Expected: mustBe("42"),
				})
		})
	})

	t.Run("nil channel", func(t *testing.T) {
		var ch chan int
		checkError(t, ch, td.Recv(42),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("recv(DATA)"),
				Got:      mustBe("nothing received on channel"),
				Expected: mustBe("42"),
			})
		checkError(t, &ch, td.Recv(42),
			expectedError{
				Message:  mustBe("values differ"),
				Path:     mustBe("recv(DATA)"),
				Got:      mustBe("nothing received on channel"),
				Expected: mustBe("42"),
			})
	})

	t.Run("nil pointer", func(t *testing.T) {
		checkError(t, (*chan int)(nil), td.Recv(42),
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *chan (*chan int type)"),
				Expected: mustBe("non-nil *chan"),
			})
	})

	t.Run("chan any", func(t *testing.T) {
		ch := make(chan any, 6)
		fillCh := func(val any) {
			ch <- val // td.Cmp
			ch <- val // EqDeeply aka boolean context
			ch <- val // EqDeeplyError
			ch <- val // interface + td.Cmp
			ch <- val // interface + EqDeeply aka boolean context
			ch <- val // interface + EqDeeplyError
		}

		fillCh(1)
		checkOK(t, ch, td.Recv(1))

		fillCh(nil)
		checkOK(t, ch, td.Recv(nil))

		close(ch)
		checkOK(t, ch, td.Recv(td.RecvClosed))
	})

	t.Run("errors", func(t *testing.T) {
		checkError(t, "never tested",
			td.Recv(23, time.Second, time.Second),
			expectedError{
				Message: mustBe("bad usage of Recv operator"),
				Path:    mustBe("DATA"),
				Summary: mustBe("usage: Recv(EXPECTED[, TIMEOUT]), too many parameters"),
			})

		checkError(t, 42, td.Recv(33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("chan OR *chan"),
			})

		checkError(t, &struct{}{}, td.Recv(33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*struct (*struct {} type)"),
				Expected: mustBe("chan OR *chan"),
			})

		checkError(t, nil, td.Recv(33),
			expectedError{
				Message:  mustBe("bad kind"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("chan OR *chan"),
			})
	})
}

func TestRecvString(t *testing.T) {
	test.EqualStr(t, td.Recv(3).String(), "recv=3")
	test.EqualStr(t, td.Recv(td.Between(3, 8)).String(), "recv: 3 ≤ got ≤ 8")
	test.EqualStr(t, td.Recv(td.Gt(8)).String(), "recv: > 8")

	// Erroneous op
	test.EqualStr(t, td.Recv(3, 0, 0).String(), "Recv(<ERROR>)")
}

func TestRecvTypeBehind(t *testing.T) {
	equalTypes(t, td.Recv(3), 0)
	equalTypes(t, td.Recv(td.Between(3, 4)), 0)

	// Erroneous op
	equalTypes(t, td.Recv(3, 0, 0), nil)
}

func TestRecvKind(t *testing.T) {
	test.IsTrue(t, td.RecvNothing == types.RecvNothing)
	test.IsTrue(t, td.RecvClosed == types.RecvClosed)
}
