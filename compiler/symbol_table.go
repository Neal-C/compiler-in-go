package compiler

type SymbolScope string

const (
	GlobalScope = "GLOBAL"
	LocalScope  = "LOCAL"
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

	symbol := Symbol{Name: name, Index: self.numberOfDefinitions, Scope: GlobalScope}

	self.store[name] = symbol

	self.numberOfDefinitions++

	return symbol
}

func (self *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := self.store[name]

	return symbol, ok
}

func NewEnclosedTable(outer *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()

	symbolTable.OuterTable = outer

	return symbolTable

}
