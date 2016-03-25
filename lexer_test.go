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
	expr := "(7 * (3 + 4 - 2)) << 1.23 % 0.3"
	l := newLexer(expr)

	res, errs := l.Lex()
	expected := []tokenType{
		LPAREN, INT, MUL, LPAREN, INT, ADD, INT, SUB, INT, RPAREN, RPAREN, LSH,
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

func TestLongToken(t *testing.T) {
	l := newLexer("** **=")
	res, _ := l.Lex()

	expected := []tokenType{POW, POW_EQ, EOL}

	for k, v := range res {
		if expected[k] != v.Type {
			t.Error("misread long token(s)")
		}
	}
}
