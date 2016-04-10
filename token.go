// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package eparser

import "fmt"

type tokenType int

const (
	ILLEGAL tokenType = iota
	EOL

	literalsBegin
	IDENT // x
	INT   // 3
	FLOAT // 3.14
	literalsEnd

	operatorsBegin
	ADD // +
	SUB // -
	DIV // /
	MUL // *
	POW // **
	REM // %

	AND // &
	OR  // |
	XOR // ^
	LSH // <<
	RSH // >>
	NOT // ~

	EQ     // =
	ADD_EQ // +=
	SUB_EQ // -=
	DIV_EQ // /=
	MUL_EQ // *=
	POW_EQ // **=
	REM_EQ // %=

	AND_EQ // &=
	OR_EQ  // |=
	XOR_EQ // ^=
	LSH_EQ // <<=
	RSH_EQ // >>=
	operatorsEnd

	LPAREN // (
	RPAREN // )

	TOKEN_COUNT
)

var tokens = [...]string{
	ILLEGAL: "illegal",
	EOL:     "end of line",

	IDENT: "identifier",
	INT:   "integer",
	FLOAT: "float",

	ADD: "+",
	SUB: "-",
	DIV: "/",
	MUL: "*",
	POW: "**",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	LSH: "<<",
	RSH: ">>",
	NOT: "~",

	EQ:     "=",
	ADD_EQ: "+=",
	SUB_EQ: "-=",
	DIV_EQ: "/=",
	MUL_EQ: "*=",
	POW_EQ: "**=",
	REM_EQ: "%=",

	AND_EQ: "&=",
	OR_EQ:  "|=",
	XOR_EQ: "^=",
	LSH_EQ: "<<=",
	RSH_EQ: ">>=",

	LPAREN: "(",
	RPAREN: ")",
}

// Token is an entity in an expression
type Token struct {
	Type  tokenType
	Value string
	Pos   int
}

func newToken(toktype tokenType, val string, pos int) *Token {
	return &Token{
		Type:  toktype,
		Value: val,
		Pos:   pos,
	}
}

func (tok *Token) String() string {
	if tok.Type < TOKEN_COUNT {
		return fmt.Sprintf("%d: '%s' ( %s )", tok.Pos, tok.Value, tok.Type)
	}

	return fmt.Sprintf("%d: '%s' ( ??? )", tok.Pos, tok.Value)
}

func (t tokenType) String() string {
	if t < TOKEN_COUNT {
		return tokens[t]
	}

	return "???"
}

// IsOperator checks if the token is an operator
func (tok *Token) IsOperator() bool {
	return tok.Type > operatorsBegin && tok.Type < operatorsEnd
}

// IsLiteral checks if the token is a literal
func (tok *Token) IsLiteral() bool {
	return tok.Type > literalsBegin && tok.Type < literalsEnd
}
