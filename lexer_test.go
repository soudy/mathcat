package eparser

import (
	"testing"
)

func TestLex(t *testing.T) {
	l := newLexer()
	res, errs := l.Lex("a **= (7 ** (3 + 4 - 2)) << 1.23 % 0.3")
	expected := []tokenType{
		IDENT, POW_EQ, LPAREN, INT, POW, LPAREN, INT, ADD, INT, SUB, INT,
		RPAREN, RPAREN, LSH, FLOAT, REM, FLOAT, EOL,
	}

	if errs != nil {
		t.Error("lexer error(s) found")
	}

	for k, v := range res {
		if expected[k] != v.Type {
			t.Error("mismatched token")
		}
	}
}

func TestUTF8(t *testing.T) {
	if !isIdent('Å') || !isIdent('Ś') {
		t.Error("isIdent doesn't recognize unicode characters")
	}
}
