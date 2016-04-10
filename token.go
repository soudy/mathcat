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
)

var tokens = map[tokenType]string{
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

// token is an entity in an expression
type token struct {
	Type  tokenType
	Value string
	Pos   int
}

func newToken(toktype tokenType, val string, pos int) *token {
	return &token{
		Type:  toktype,
		Value: val,
		Pos:   pos,
	}
}

func (tok *token) String() string {
	if _, ok := tokens[tok.Type]; ok {
		return fmt.Sprintf("%d: '%s' ( %s )\n", tok.Pos, tok.Value, tok.Type)
	}

	return fmt.Sprintf("%d: '%s' ( ??? )\n", tok.Pos, tok.Value)
}

func (t tokenType) String() string {
	if _, ok := tokens[t]; ok {
		return tokens[t]
	}

	return "???"
}

// IsOperator checks if the token is an operator
func (tok *token) IsOperator() bool {
	return tok.Type > operatorsBegin && tok.Type < operatorsEnd
}

// IsLiteral checks if the token is a literal
func (tok *token) IsLiteral() bool {
	return tok.Type > literalsBegin && tok.Type < literalsEnd
}
