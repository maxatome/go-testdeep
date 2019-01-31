// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/util"
)

type tdSetResultKind uint8

const (
	itemsSetResult tdSetResultKind = iota
	keysSetResult
)

// Implements fmt.Stringer.
func (k tdSetResultKind) String() string {
	switch k {
	case itemsSetResult:
		return "item"
	case keysSetResult:
		return "key"
	default:
		return "?"
	}
}

type tdSetResult struct {
	types.TestDeepStamp
	Missing []reflect.Value
	Extra   []reflect.Value
	Kind    tdSetResultKind
}

func (r tdSetResult) IsEmpty() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0
}

func (r tdSetResult) Summary() ctxerr.ErrorSummary {
	var missing, extra string

	if len(r.Missing) > 0 {
		if len(r.Missing) > 1 {
			missing = fmt.Sprintf("Missing %d %ss: ", len(r.Missing), r.Kind)
		} else {
			missing = fmt.Sprintf("Missing %s: ", r.Kind)
		}
	}

	if len(r.Extra) > 0 {
		if len(r.Extra) > 1 {
			extra = fmt.Sprintf("Extra %d %ss: ", len(r.Extra), r.Kind)
		} else {
			extra = fmt.Sprintf("Extra %s: ", r.Kind)
		}
	}

	if len(missing) != len(extra) && missing != "" && extra != "" {
		if len(missing) > len(extra) {
			extra = strings.Repeat(" ", len(missing)-len(extra)) + extra
		} else {
			missing = strings.Repeat(" ", len(extra)-len(missing)) + missing
		}
	}

	var summary ctxerr.ErrorSummaryItems

	if missing != "" {
		summary = append(summary, ctxerr.ErrorSummaryItem{
			Label: missing,
			Value: util.ToString(r.Missing),
		})
	}

	if extra != "" {
		summary = append(summary, ctxerr.ErrorSummaryItem{
			Label: extra,
			Value: util.ToString(r.Extra),
		})
	}

	return summary
}
