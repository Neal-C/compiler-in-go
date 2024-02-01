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

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		Symbol{Name: "b", Scope: GlobalScope, Index: 1},
		Symbol{Name: "c", Scope: LocalScope, Index: 0},
		Symbol{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, symbl := range expected {
		result, ok := local.Resolve(symbl.Name)

		if !ok {
			t.Errorf("Name is not resolvable: %s", symbl.Name)
			continue
		}

		if result != symbl {
			t.Errorf("expected %s to resolve to %+v , got = %+v", symbl.Name, symbl, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	testTable := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			table: firstLocal,
			expectedSymbols: []Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			table: secondLocal,
			expectedSymbols: []Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "e", Scope: LocalScope, Index: 0},
				Symbol{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range testTable {
		for _, symbl := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(symbl.Name)

			if !ok {
				t.Errorf("name %s not resolvable", symbl.Name)
				continue
			}

			if result != symbl {
				t.Errorf("expected %s to resolve to %+v, got = %+v", result.Name, symbl, result)
			}
		}
	}
}
