package eparser

import (
	"testing"
)

var (
	testExpr        = []rune("3+5*8")
	complexTestExpr = []rune("(7 * (3 + 4 - 2)) << 1.23 % 0.3")
)

func TestPeekAndEat(t *testing.T) {
	l := newLexer(testExpr)
	if l.peek() != '3' {
		t.Error("wrong peeked value")
	}

	if l.eat() != '3' && l.ch != '+' {
		t.Error("eating goes wrong")
	}
}

func TestLex(t *testing.T) {
	l := newLexer(complexTestExpr)
	res, errs := l.Lex()
	expected := map[int]tokenType{
		0:  LPAREN,
		1:  INT,
		2:  MUL,
		3:  LPAREN,
		4:  INT,
		5:  ADD,
		6:  INT,
		7:  SUB,
		8:  INT,
		9:  RPAREN,
		10: RPAREN,
		11: LSH,
		12: FLOAT,
		13: REM,
		14: FLOAT,
		15: EOL,
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
