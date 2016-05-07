// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

import (
	"testing"
)

func TestLex(t *testing.T) {
	res, err := Lex("some_var123 **= (.5 ** (3 + 4 - 2)) <<= 1.23 % 0.3")
	expected := []tokenType{
		IDENT, POW_EQ, LPAREN, NUMBER, POW, LPAREN, NUMBER, ADD, NUMBER, SUB, NUMBER,
		RPAREN, RPAREN, LSH_EQ, NUMBER, REM, NUMBER, EOL,
	}

	if err != nil {
		t.Error("lexer error occured")
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
