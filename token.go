// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

// TokenType represents the type of token
type TokenType int

// Token is an entity in an expression
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// Tokens is a slice of pointers to a token
type Tokens []*Token

const (
	Eol TokenType = iota // end of line

	literalsBegin
	Ident   // x
	Decimal // 3
	Hex     // 0xDEADBEEF
	Binary  // 0b10101101100
	Octal   // 0o666
	literalsEnd

	operatorsBegin
	Add      // +
	Sub      // -
	Div      // /
	Mul      // *
	Pow      // **
	Rem      // %
	UnaryMin // -

	bitwiseBegin
	And // &
	Or  // |
	Xor // ^
	Lsh // <<
	Rsh // >>
	Not // ~

	assignmentBegin
	AndEq // &=
	OrEq  // |=
	XorEq // ^=
	LshEq // <<=
	RshEq // >>=
	bitwiseEnd

	Eq    // =
	AddEq // +=
	SubEq // -=
	DivEq // /=
	MulEq // *=
	PowEq // **=
	RemEq // %=
	assignmentEnd

	NotEq // !=
	EqEq  // ==
	Gt    // >
	GtEq  // >=
	Lt    // <
	LtEq  // <=
	operatorsEnd

	Lparen // (
	Rparen // )
	Comma  // ,
)

var tokens = map[TokenType]string{
	Eol: "end of line",

	Ident:   "identifier",
	Decimal: "decimal number",
	Hex:     "hex number",
	Binary:  "binary number",
	Octal:   "octal number",

	Add:      "+",
	Sub:      "-",
	Div:      "/",
	Mul:      "*",
	Pow:      "**",
	Rem:      "%",
	UnaryMin: "-",

	And: "&",
	Or:  "|",
	Xor: "^",
	Lsh: "<<",
	Rsh: ">>",
	Not: "~",

	Eq:    "=",
	AddEq: "+=",
	SubEq: "-=",
	DivEq: "/=",
	MulEq: "*=",
	PowEq: "**=",
	RemEq: "%=",

	AndEq: "&=",
	OrEq:  "|=",
	XorEq: "^=",
	LshEq: "<<=",
	RshEq: ">>=",

	NotEq: "!=",
	EqEq:  "==",
	Gt:    ">",
	GtEq:  ">=",
	Lt:    "<",
	LtEq:  "<=",

	Lparen: "(",
	Rparen: ")",
	Comma:  ",",
}

func (tok Token) String() string {
	return tok.Value
}

func (t TokenType) String() string {
	if tok, ok := tokens[t]; ok {
		return tok
	}

	return "???"
}

// Is checks if the token is given token type
func (tok Token) Is(toktype TokenType) bool {
	return tok.Type == toktype
}

// IsOperator checks if the token is an operator
func (tok Token) IsOperator() bool {
	return tok.Type > operatorsBegin && tok.Type < operatorsEnd
}

// IsBitwise checks if the token type is a bitwise operator
func (tok Token) IsBitwise() bool {
	return tok.Type > bitwiseBegin && tok.Type < bitwiseEnd
}

// IsLiteral checks if the token is a literal
func (tok Token) IsLiteral() bool {
	return tok.Type > literalsBegin && tok.Type < literalsEnd
}

// IsAssignment checks if the token is an assignment operator
func (tok Token) IsAssignment() bool {
	return tok.Type > assignmentBegin && tok.Type < assignmentEnd
}
