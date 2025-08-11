/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
This test file is part of the spew package rather than than the spew_test
package because it needs access to internals to properly test certain cases
which are not possible via the public interface since they should never happen.
*/

package spew

import (
	"reflect"
	"strings"
	"testing"
)

// TestInvalidReflectValue ensures the dump and formatter code handles an
// invalid reflect value properly.  This needs access to internal state since it
// should never happen in real code and therefore can't be tested via the public
// API.
func TestInvalidReflectValue(t *testing.T) {
	// Dump invalid reflect value.
	v := new(reflect.Value)
	var buf strings.Builder
	d := dumpState{w: &buf, cs: &Config}
	d.dump(*v)
	s := buf.String()
	want := "<invalid>"
	if s != want {
		t.Errorf("InvalidReflectValue\n got: %s want: %s", s, want)
	}
}
