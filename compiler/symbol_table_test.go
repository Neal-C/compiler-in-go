package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	gloabal := NewSymbolTable()

	a := gloabal.Define("a")

	if a != expected["a"] {
		t.Errorf("expected a=%+v , got %+v", expected["a"], a)
	}

	b := gloabal.Define("b")

	if b != expected["b"] {
		t.Errorf("expected b=%+v , got %+v", expected["b"], b)
	}
}
