package eparser

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	res, errs := Lex("some_var123 **= (.5 ** (3 + 4 - 2)) <<= 1.23 % 0.3")
	expected := []tokenType{
		IDENT, POW_EQ, LPAREN, FLOAT, POW, LPAREN, INT, ADD, INT, SUB, INT,
		RPAREN, RPAREN, LSH_EQ, FLOAT, REM, FLOAT, EOL,
	}

	fmt.Print(res)

	if errs != nil {
		t.Error("lexer error(s) found")
	}

	for k, v := range res {
		if expected[k] != v.Type {
			t.Errorf("mismatched token: expected %s, got %s", expected[k], v.Type)
		}
	}
}

func TestUTF8(t *testing.T) {
	if !isIdent('Å') || !isIdent('Ś') {
		t.Error("isIdent doesn't recognize unicode characters")
	}
}
