package ast

import (
	"github.com/Neal-C/compiler-in-go/token"
	"testing"
)

func TestString(t *testing.T) {
	input := "let myVar = anotherVar;"

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != input {
		t.Errorf("program.String() wront, got %q", program.String())
	}
}
