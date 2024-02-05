package compiler

type SymbolScope string

const (
	GlobalScope  = "GLOBAL"
	LocalScope   = "LOCAL"
	BuiltinScope = "BUILTIN"
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
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
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
		return symbol, ok
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
