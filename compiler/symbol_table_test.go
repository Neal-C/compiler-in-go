package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")

	if a != expected["a"] {
		t.Errorf("expected a=%+v , got %+v", expected["a"], a)
	}

	b := global.Define("b")

	if b != expected["b"] {
		t.Errorf("expected b=%+v , got %+v", expected["b"], b)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, symbol := range expected {
		result, ok := global.Resolve(symbol.Name)

		if !ok {
			t.Errorf("name %s is not resolvable", symbol.Name)
			continue
		}

		if result != symbol {
			t.Errorf("expected %s to resolve to %+v , got %+v", symbol.Name, symbol, result)
		}
	}

}