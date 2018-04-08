package testdeep

import (
	"bytes"
	"reflect"
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
		return "items"
	case keysSetResult:
		return "keys"
	default:
		return "?"
	}
}

type tdSetResult struct {
	Missing []reflect.Value
	Extra   []reflect.Value
	Kind    tdSetResultKind
}

var _ testDeepStringer = tdSetResult{}

func (r tdSetResult) _TestDeep() {}

func (r tdSetResult) IsEmpty() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0
}

func (r tdSetResult) String() string {
	buf := &bytes.Buffer{}

	if len(r.Missing) > 0 {
		buf.WriteString("Missing ")
		buf.WriteString(r.Kind.String())
		buf.WriteString(": ")
		sliceToBuffer(buf, r.Missing)
	}

	if len(r.Extra) > 0 {
		if buf.Len() > 0 {
			buf.WriteString("\n  ")
		}
		buf.WriteString("Extra ")
		buf.WriteString(r.Kind.String())
		buf.WriteString(": ")
		sliceToBuffer(buf, r.Extra)
	}

	return buf.String()
}
