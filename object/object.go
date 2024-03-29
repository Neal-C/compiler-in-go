package object

import (
	"bytes"
	"fmt"
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/code"
	"hash/fnv"
	"log"
	"strings"
)

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	STRING_OBJ            = "STRING"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	HASH_OBJ              = "HASH"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"
	CLOSURE_OBJ           = "CLOSURE"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (self *Null) Type() ObjectType { return NULL_OBJ }
func (self *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (self *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (self *ReturnValue) Inspect() string  { return self.Value.Inspect() }

type Error struct {
	Message string
}

func (self *Error) Type() ObjectType { return ERROR_OBJ }
func (self *Error) Inspect() string  { return "ERROR: " + self.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (self *Function) Type() ObjectType { return FUNCTION_OBJ }
func (self *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, param := range self.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(self.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (self *String) Type() ObjectType { return STRING_OBJ }
func (self *String) Inspect() string  { return self.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (self *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (self *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (self *Array) Type() ObjectType { return ARRAY_OBJ }
func (self *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string

	for _, element := range self.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

func (self *Boolean) HashKey() HashKey {
	var value uint64

	if self.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: self.Type(), Value: value}
}

func (self *Integer) HashKey() HashKey {
	return HashKey{Type: self.Type(), Value: uint64(self.Value)}
}

func (self *String) HashKey() HashKey {
	h := fnv.New64()
	_, err := h.Write([]byte(self.Value))
	if err != nil {
		log.Println("fnv.New64.Write failed: ", err)
	}

	return HashKey{Type: self.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (self *Hash) Type() ObjectType { return HASH_OBJ }

func (self *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string

	for _, pair := range self.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type CompiledFunction struct {
	Instructions       code.Instructions
	NumberOfLocals     int
	NumberOfParameters int
}

func (self *CompiledFunction) Type() ObjectType {
	return COMPILED_FUNCTION_OBJ
}

func (self *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", self)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (self *Closure) Type() ObjectType {
	return CLOSURE_OBJ
}

func (self *Closure) Inspect() string {
	return fmt.Sprintf("CLOSURE[%p]", self)
}
