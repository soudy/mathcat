// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import "testing"

func TestLex(t *testing.T) {
	res, err := Lex("some_var123 **= (.5 ** (3 + 4 - 2)) <<= 1.23 % -0.3")
	expected := []TokenType{
		Ident, PowEq, Lparen, Decimal, Pow, Lparen, Decimal, Add, Decimal, Sub,
		Decimal, Rparen, Rparen, LshEq, Decimal, Rem, UnaryMin, Decimal, Eol,
	}

	if err != nil {
		t.Errorf("unexpected lexer error occured: %s", err)
	}

	for k, v := range res {
		if expected[k] != v.Type {
			t.Errorf("mismatched token: expected %s, got %s", expected[k], v.Type)
		}
	}
}

func TestLiterals(t *testing.T) {
	res, err := Lex("0xBEEF 0b10101010 0o111762 12.23 .33 2e-10 2E10 0XBBA")
	expected := []TokenType{
		Hex, Binary, Octal, Decimal, Decimal, Decimal, Decimal, Hex, Eol,
	}

	if err != nil {
		t.Errorf("unexpected lexer error occured: %s", err)
	}

	for k, v := range res {
		if expected[k] != v.Type {
			t.Errorf("mismatched token: expected %s, got %s", expected[k], v.Type)
		}
	}
}

func TestOperators(t *testing.T) {
	// We add a number before - sign so it doesn't see it as unary
	res, err := Lex(`= += -= /= *= **= %= &= |=  ^= <<= >>= == != > >= < <= | ^
	& << >> ~ + 5 - * / ** % -`)
	expected := []TokenType{
		Eq, AddEq, SubEq, DivEq, MulEq, PowEq, RemEq, AndEq, OrEq,
		XorEq, LshEq, RshEq, EqEq, NotEq, Gt, GtEq, Lt, LtEq, Or, Xor,
		And, Lsh, Rsh, Not, Add, Decimal, Sub, Mul, Div, Pow, Rem, UnaryMin, Eol,
	}

	if err != nil {
		t.Errorf("unexpected lexer error occured: %s", err)
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
