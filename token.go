// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

type tokenType int

const (
	ILLEGAL tokenType = iota
	EOL

	literalsBegin
	IDENT  // x
	NUMBER // 3
	HEX    // 0xDEADBEEF
	BINARY // 0b10101101100
	OCTAL  // 0o666
	literalsEnd

	operatorsBegin
	ADD       // +
	SUB       // -
	DIV       // /
	MUL       // *
	POW       // **
	REM       // %
	UNARY_MIN // -

	bitwiseBegin
	AND // &
	OR  // |
	XOR // ^
	LSH // <<
	RSH // >>
	NOT // ~

	assignmentBegin
	AND_EQ // &=
	OR_EQ  // |=
	XOR_EQ // ^=
	LSH_EQ // <<=
	RSH_EQ // >>=
	bitwiseEnd

	EQ     // =
	ADD_EQ // +=
	SUB_EQ // -=
	DIV_EQ // /=
	MUL_EQ // *=
	POW_EQ // **=
	REM_EQ // %=
	assignmentEnd

	BANG_EQ // !=
	EQ_EQ   // ==
	GT      // >
	GT_EQ   // >=
	LT      // <
	LT_EQ   // <=
	operatorsEnd

	LPAREN // (
	RPAREN // )
	COMMA  // ,
)

var tokens = map[tokenType]string{
	ILLEGAL: "illegal",
	EOL:     "end of line",

	IDENT:  "identifier",
	NUMBER: "number",
	HEX:    "hex number",
	BINARY: "binary number",
	OCTAL:  "octal number",

	ADD:       "+",
	SUB:       "-",
	DIV:       "/",
	MUL:       "*",
	POW:       "**",
	REM:       "%",
	UNARY_MIN: "-",

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

	BANG_EQ: "!=",
	EQ_EQ:   "==",
	GT:      ">",
	GT_EQ:   ">=",
	LT:      "<",
	LT_EQ:   "<=",

	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
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
	return tok.Value
}

func (t tokenType) String() string {
	if tok, ok := tokens[t]; ok {
		return tok
	}

	return "???"
}

// Is checks if the token is given token type
func (tok *Token) Is(toktype tokenType) bool {
	return tok.Type == toktype
}

// IsOperator checks if the token is an operator
func (tok *Token) IsOperator() bool {
	return tok.Type > operatorsBegin && tok.Type < operatorsEnd
}

// IsBitwise checks if the token type is a bitwise operator
func (tok *Token) IsBitwise() bool {
	return tok.Type > bitwiseBegin && tok.Type < bitwiseEnd
}

// IsLiteral checks if the token is a literal
func (tok *Token) IsLiteral() bool {
	return tok.Type > literalsBegin && tok.Type < literalsEnd
}

// IsAssignment checks if the token is an assignment operator
func (tok *Token) IsAssignment() bool {
	return tok.Type > assignmentBegin && tok.Type < assignmentEnd
}
