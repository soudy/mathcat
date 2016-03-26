package eparser

import (
	"testing"
)

func TestPeekAndEat(t *testing.T) {
	l := newLexer("3+5*8")
	if l.peek() != '3' {
		t.Error("wrong peeked value")
	}

	if l.eat() != '3' && l.ch != '+' {
		t.Error("eating goes wrong")
	}
}

func TestLex(t *testing.T) {
	expr := "a **= (7 ** (3 + 4 - 2)) << 1.23 % 0.3"
	l := newLexer(expr)

	res, errs := l.Lex()
	expected := []tokenType{
		IDENT, POW_EQ, LPAREN, INT, POW, LPAREN, INT, ADD, INT, SUB, INT, RPAREN, RPAREN, LSH,
		FLOAT, REM, FLOAT, EOL,
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
