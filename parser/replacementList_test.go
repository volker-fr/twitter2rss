package parser

import (
	"testing"
)

func TestLen(t *testing.T) {
	i := ReplacementList{}
	j := ReplacementList{{from: 2, to: 3, replacement: "b"}}
	k := ReplacementList{{from: 1, to: 2, replacement: "a"}, {from: 2, to: 3, replacement: "b"}}

	if i.Len() != 0 {
		t.Error("Len test failed for one element")
	}
	if j.Len() != 1 {
		t.Error("Len test failed for two elements")
	}
	if k.Len() != 2 {
		t.Error("Len test failed for two elements")
	}
}

func TestSwap(t *testing.T) {
	i := ReplacementList{{from: 2, to: 3, replacement: "a"}, {from: 1, to: 2, replacement: "b"}}
	// We don't have a Clone(), but this works for now
	before := i[0].from
	i.Swap(0, 1)
	after := i[0].from

	if after == before {
		t.Error("Swap didn't work", before, "=", after)
	}
}

func TestLess(t *testing.T) {
	i := ReplacementList{{from: 1, to: 2, replacement: "a"}, {from: 2, to: 3, replacement: "b"}}

	if !i.Less(0, 1) {
		t.Error("Replacement list Less comparison failed")
	}
}
