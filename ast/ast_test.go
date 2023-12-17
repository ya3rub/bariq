package ast

import (
	"testing"

	"bariq/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: []Stmt{
			&LetStmt{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "hello"},
					Value: "hello",
				},
				Value: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "world"},
					Value: "world",
				},
			},

			&ReturnStmt{
				Token: token.Token{Type: token.LET, Literal: "return"},
				Value: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "world"},
					Value: "world",
				},
			},
		},
	}
	expectedString := "let hello = world;return world;"
	if program.String() != expectedString {
		t.Errorf(
			"program.String() error, expected %s but got %s",
			expectedString,
			program.String(),
		)
	}
}
