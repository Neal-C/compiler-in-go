package compiler

type SymbolScope string

const (
	GlobalScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store               map[string]Symbol
	numberOfDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

func (self *SymbolTable) Define(name string) Symbol {

	symbol := Symbol{Name: name, Index: self.numberOfDefinitions, Scope: GlobalScope}

	self.store[name] = symbol

	self.numberOfDefinitions++

	return symbol
}

func (self *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := self.store[name]

	return symbol, ok
}
