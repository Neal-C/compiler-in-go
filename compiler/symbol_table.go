package compiler

type SymbolScope string

const (
	GlobalScope   = "GLOBAL"
	LocalScope    = "LOCAL"
	BuiltinScope  = "BUILTIN"
	FreeScope     = "FREE"
	FunctionScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	OuterTable *SymbolTable

	store               map[string]Symbol
	numberOfDefinitions int
	FreeSymbols         []Symbol
}

func NewSymbolTable() *SymbolTable {
	freeSymbols := []Symbol{}
	return &SymbolTable{
		store:       make(map[string]Symbol),
		FreeSymbols: freeSymbols,
	}
}

func (self *SymbolTable) Define(name string) Symbol {

	symbol := Symbol{Name: name, Index: self.numberOfDefinitions}

	if self.OuterTable == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	self.store[name] = symbol

	self.numberOfDefinitions++

	return symbol
}

func (self *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := self.store[name]

	if !ok && self.OuterTable != nil {

		symbol, ok := self.OuterTable.Resolve(name)

		if !ok {
			return symbol, ok
		}

		if symbol.Scope == GlobalScope || symbol.Scope == BuiltinScope {
			return symbol, ok
		}

		freeSymbol := self.DefineFree(symbol)

		return freeSymbol, true
	}

	return symbol, ok
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()

	symbolTable.OuterTable = outer

	return symbolTable

}

func (self *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	self.store[name] = symbol
	return symbol
}

func (self *SymbolTable) DefineFree(original Symbol) Symbol {

	self.FreeSymbols = append(self.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(self.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	self.store[original.Name] = symbol

	return symbol
}

func (self *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	self.store[name] = symbol
	return symbol
}
