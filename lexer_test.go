// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import "testing"

func TestLex(t *testing.T) {
	res, err := Lex("some_var123 **= (.5 ** (3 + 4 - 2)) <<= 1.23 % -0.3")
	expected := []TokenType{
		IDENT, POW_EQ, LPAREN, NUMBER, POW, LPAREN, NUMBER, ADD, NUMBER, SUB,
		NUMBER, RPAREN, RPAREN, LSH_EQ, NUMBER, REM, UNARY_MIN, NUMBER, EOL,
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
		HEX, BINARY, OCTAL, NUMBER, NUMBER, NUMBER, NUMBER, HEX, EOL,
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
		EQ, ADD_EQ, SUB_EQ, DIV_EQ, MUL_EQ, POW_EQ, REM_EQ, AND_EQ, OR_EQ,
		XOR_EQ, LSH_EQ, RSH_EQ, EQ_EQ, BANG_EQ, GT, GT_EQ, LT, LT_EQ, OR, XOR,
		AND, LSH, RSH, NOT, ADD, NUMBER, SUB, MUL, DIV, POW, REM, UNARY_MIN, EOL,
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
