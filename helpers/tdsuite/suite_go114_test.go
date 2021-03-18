// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build go1.14

package tdsuite_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
)

// FullCleanup has tests and all possible hooks.
type FullCleanup struct{ base }

func (f *FullCleanup) Setup(t *td.T) error { f.rec(); return nil }
func (f *FullCleanup) PreTest(t *td.T, tn string) error {
	f.rec(tn)
	t.Cleanup(func() { f.rec(tn) })
	return nil
}
func (f *FullCleanup) PostTest(t *td.T, tn string) error {
	f.rec(tn)
	t.Cleanup(func() { f.rec(tn) })
	return nil
}
func (f *FullCleanup) BetweenTests(t *td.T, prev, next string) error {
	f.rec(prev, next)
	return nil
}
func (f *FullCleanup) Destroy(t *td.T) error { f.rec(); return nil }

func (f *FullCleanup) Test1(t *td.T) {
	f.rec()
	t.Cleanup(func() { f.rec() })
}
func (f *FullCleanup) Test2(assert *td.T, require *td.T) {
	f.rec()
	assert.Cleanup(func() { f.rec() })
}
func (f *FullCleanup) Test3(t *td.T) {
	f.rec()
	t.Cleanup(func() { f.rec() })
}
func (f *FullCleanup) Testimony(t *td.T) {} // not a test method

var (
	_ tdsuite.Setup        = (*FullCleanup)(nil)
	_ tdsuite.PreTest      = (*FullCleanup)(nil)
	_ tdsuite.PostTest     = (*FullCleanup)(nil)
	_ tdsuite.BetweenTests = (*FullCleanup)(nil)
	_ tdsuite.Destroy      = (*FullCleanup)(nil)
)

func TestRunCleanup(t *testing.T) {
	t.Run("Full", func(t *testing.T) {
		suite := FullCleanup{}
		td.CmpTrue(t, tdsuite.Run(t, &suite))
		ok := td.Cmp(t, suite.calls, []string{
			"Setup",
			/**/ "PreTest+Test1",
			/**/ "Test1",
			/**/ "PostTest+Test1",
			/**/ "PostTest.Cleanup+Test1",
			/**/ "Test1.Cleanup",
			/**/ "PreTest.Cleanup+Test1",
			"BetweenTests+Test1+Test2",
			/**/ "PreTest+Test2",
			/**/ "Test2",
			/**/ "PostTest+Test2",
			/**/ "PostTest.Cleanup+Test2",
			/**/ "Test2.Cleanup",
			/**/ "PreTest.Cleanup+Test2",
			"BetweenTests+Test2+Test3",
			/**/ "PreTest+Test3",
			/**/ "Test3",
			/**/ "PostTest+Test3",
			/**/ "PostTest.Cleanup+Test3",
			/**/ "Test3.Cleanup",
			/**/ "PreTest.Cleanup+Test3",
			"Destroy",
		})
		if !ok {
			for _, c := range suite.calls {
				switch c[0] {
				case 'S', 'B', 'D':
					t.Log(c)
				default:
					t.Log("  ", c)
				}
			}
		}
	})
}
